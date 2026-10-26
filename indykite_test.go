// Copyright (c) 2026 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package indykite_test

import (
	"io"
	"net/http"
	"time"

	"github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/test/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//nolint:gosec // environment variable names, not credentials.
const (
	envAppCreds = "INDYKITE_APPLICATION_CREDENTIALS"
	envSACreds  = "INDYKITE_SERVICE_ACCOUNT_CREDENTIALS"
)

// A valid EC P-256 credential (test key, not a real secret).
const ecKeyJWK = `"privateKeyJWK": {"kty":"EC","d":"2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s","use":"sig",
  "crv":"P-256","x":"Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA",
  "y":"DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4","alg":"ES256"}`

// appAgentCred is the App Agent credential: the raw token itself.
func appAgentCred() string {
	return "test-app-agent-token"
}

func serviceAccountCred() []byte {
	return []byte(`{"serviceAccountId":"fa50a80e-4840-4fc0-8958-982b84827f83",` + ecKeyJWK + `}`)
}

var _ = Describe("Runtime facade", func() {
	It("routes AuthZEN and CIQ calls and authenticates with the App Agent client key", func(ctx SpecContext) {
		var paths []string
		var clientKey string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			paths = append(paths, r.URL.Path)
			if v := r.Header.Get("X-Ik-Clientkey"); v != "" { // canonical form of X-IK-ClientKey
				clientKey = v
			}
			switch r.URL.Path {
			case "/access/v1/evaluation":
				_, _ = io.WriteString(w, `{"decision":true}`)
			case "/contx-iq/v1/execute":
				_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"x"}}}]}`)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		})

		cli, err := indykite.NewClient(ctx, appAgentCred(), indykite.WithBaseURL(srv.URL))
		Expect(err).NotTo(HaveOccurred())

		ok, err := cli.AuthZEN().Allowed(ctx,
			authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		rows, err := cli.CIQ().All(ctx, ciq.ExecuteRequest{ID: "q", PageSize: 100})
		Expect(err).NotTo(HaveOccurred())
		Expect(rows).To(HaveLen(1))

		Expect(clientKey).NotTo(BeEmpty(), "runtime plane should send the App Agent client key")
		Expect(paths).To(ContainElements("/access/v1/evaluation", "/contx-iq/v1/execute"))
	})

	It("exposes non-nil service accessors", func(ctx SpecContext) {
		cli, err := indykite.NewClient(ctx, appAgentCred(), indykite.WithRegion("eu"))
		Expect(err).NotTo(HaveOccurred())
		Expect(cli.AuthZEN()).NotTo(BeNil())
		Expect(cli.CIQ()).NotTo(BeNil())
		Expect(cli.Capture()).NotTo(BeNil())
		Expect(cli.EntityMatching()).NotTo(BeNil())
		Expect(cli.Audit()).NotTo(BeNil())
	})

	It("passes options through end to end", func(ctx SpecContext) {
		var gotUA string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			gotUA = r.Header.Get("User-Agent")
			_, _ = io.WriteString(w, `{"decision":true}`)
		})

		cli, err := indykite.NewClient(ctx, appAgentCred(),
			indykite.WithBaseURL(srv.URL),
			indykite.WithUserAgent("facade-test/1.0"),
			indykite.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
			indykite.WithRetry(indykite.RetryConfig{
				MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond,
			}),
			indykite.WithTracing())
		Expect(err).NotTo(HaveOccurred())

		ok, err := cli.AuthZEN().Allowed(ctx,
			authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
		Expect(gotUA).To(Equal("facade-test/1.0"), "WithUserAgent should be passed through")
	})
})

var _ = Describe("NewClient", func() {
	It("fails for empty credentials", func(ctx SpecContext) {
		_, err := indykite.NewClient(ctx, "", indykite.WithRegion("eu"))
		Expect(err).To(HaveOccurred())
	})

	It("rejects a JSON document in place of the credential token", func(ctx SpecContext) {
		// A JSON document is not a valid App Agent credential token.
		_, err := indykite.NewClient(ctx, `{"appAgentId":"x"}`, indykite.WithRegion("eu"))
		Expect(err).To(HaveOccurred())
	})

	It("fails when no base URL can be resolved", func(ctx SpecContext) {
		// The token carries no base URL hint; without an option the transport errors.
		_, err := indykite.NewClient(ctx, appAgentCred())
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("NewClientFromEnv", Serial, func() {
	It("builds a client from the environment and rejects invalid credentials", func(ctx SpecContext) {
		testutil.SetEnv(envAppCreds, appAgentCred())
		cli, err := indykite.NewClientFromEnv(ctx, indykite.WithRegion("eu"))
		Expect(err).NotTo(HaveOccurred())
		Expect(cli.AuthZEN()).NotTo(BeNil())

		testutil.SetEnv(envAppCreds, `{broken`)
		_, err = indykite.NewClientFromEnv(ctx, indykite.WithRegion("eu"))
		Expect(err).To(HaveOccurred(), "expected error for invalid credentials in env")
	})
})

var _ = Describe("Admin facade", func() {
	It("authenticates config calls with a Bearer token", func(ctx SpecContext) {
		var authz, path string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			authz = r.Header.Get("Authorization")
			path = r.URL.Path
			_, _ = io.WriteString(w, `{"data":[]}`)
		})

		admin, err := indykite.NewAdmin(ctx, serviceAccountCred(), indykite.WithBaseURL(srv.URL))
		Expect(err).NotTo(HaveOccurred())
		_, err = admin.AuthorizationPolicies().List(ctx, "gid:proj", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(authz).To(HavePrefix("Bearer "))
		Expect(path).To(Equal("/configs/v1/authorization-policies"))
	})

	It("fails for invalid credentials JSON", func(ctx SpecContext) {
		_, err := indykite.NewAdmin(ctx, []byte(`{broken`), indykite.WithRegion("eu"))
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("NewAdminFromEnv", Serial, func() {
	It("builds an admin client from the environment and reports resolution errors", func(ctx SpecContext) {
		testutil.SetEnv(envSACreds, string(serviceAccountCred()))
		admin, err := indykite.NewAdminFromEnv(ctx, indykite.WithRegion("eu"))
		Expect(err).NotTo(HaveOccurred())
		Expect(admin).NotTo(BeNil())

		// Credential without baseUrl and no option -> transport error.
		_, err = indykite.NewAdminFromEnv(ctx)
		Expect(err).To(HaveOccurred(), "expected error when no base URL can be resolved")

		testutil.SetEnv(envSACreds, `{broken`)
		_, err = indykite.NewAdminFromEnv(ctx, indykite.WithRegion("eu"))
		Expect(err).To(HaveOccurred(), "expected error for invalid credentials in env")
	})
})
