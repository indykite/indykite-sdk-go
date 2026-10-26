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

package config_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AdminClient", func() {
	It("authenticates the control plane with a Bearer token and no client key", func(ctx SpecContext) {
		var gotAuth, gotKey string
		admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
			gotAuth = r.Header.Get("Authorization")
			gotKey = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
			_, _ = io.WriteString(w, `{"data":[]}`)
		})

		_, err := admin.AuthorizationPolicies().List(ctx, "gid:proj", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(gotAuth).To(Equal("Bearer sa-token"))
		Expect(gotKey).To(BeEmpty(), "X-IK-ClientKey must be empty on control plane")
	})

	It("sends list query parameters and decodes the items", func(ctx SpecContext) {
		var query url.Values
		admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
			query = r.URL.Query()
			_, _ = io.WriteString(w, `{"data":[{"id":"gid:p1","name":"can-read","status":"ACTIVE","policy":"{}"}]}`)
		})

		pols, err := admin.AuthorizationPolicies().List(ctx, "gid:proj", "kbac", config.WithSearch("read"))
		Expect(err).NotTo(HaveOccurred())
		Expect(query.Get("project_id")).To(Equal("gid:proj"))
		Expect(query.Get("type")).To(Equal("kbac"))
		Expect(query.Get("search")).To(Equal("read"))
		Expect(pols).To(HaveLen(1))
		Expect(pols[0].Name).To(Equal("can-read"))
		Expect(pols[0].Status).To(Equal("ACTIVE"))
	})

	It("captures the ETag of a create", func(ctx SpecContext) {
		var gotBody map[string]any
		admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &gotBody)
			w.Header().Set("ETag", `"v1"`)
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"gid:newpol","create_time":"t0"}`)
		})

		res, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
			ProjectID: "gid:proj", Name: "can-provision", Policy: `{"x":1}`, Status: config.StatusActive,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(gotBody).To(HaveKeyWithValue("project_id", "gid:proj"))
		Expect(gotBody).To(HaveKeyWithValue("name", "can-provision"))
		Expect(gotBody).To(HaveKeyWithValue("status", "ACTIVE"))
		Expect(res.ID).To(Equal("gid:newpol"))
		Expect(res.ETag).To(Equal(`"v1"`))
	})

	// The optimistic-concurrency round trip: Read returns an ETag, Update must
	// echo it in If-Match.
	It("echoes the read ETag as If-Match on update", func(ctx SpecContext) {
		var ifMatch string
		admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("ETag", `"v7"`)
				_, _ = io.WriteString(w, `{"id":"gid:p","name":"can-read","status":"ACTIVE","policy":"{}"}`)
			case http.MethodPut:
				ifMatch = r.Header.Get("If-Match")
				w.Header().Set("ETag", `"v8"`)
				_, _ = io.WriteString(w, `{"id":"gid:p","update_time":"t1"}`)
			}
		})

		api := admin.AuthorizationPolicies()
		pol, err := api.Read(ctx, "gid:p")
		Expect(err).NotTo(HaveOccurred())
		Expect(pol.ETag).To(Equal(`"v7"`))

		desc := "updated"
		res, err := api.Update(ctx, "gid:p", pol.ETag, &config.UpdateAuthorizationPolicy{
			Description: &desc, Policy: `{"y":2}`, Status: config.StatusActive,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(ifMatch).To(Equal(`"v7"`))
		Expect(res.ETag).To(Equal(`"v8"`))
	})

	It("omits If-Match on a delete without an ETag", func(ctx SpecContext) {
		var hadIfMatch bool
		var method, path string
		admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
			_, hadIfMatch = r.Header["If-Match"]
			method, path = r.Method, r.URL.Path
			_, _ = io.WriteString(w, `{"id":"gid:p"}`)
		})

		Expect(admin.AppAgents().Delete(ctx, "gid:p", "")).To(Succeed())
		Expect(hadIfMatch).To(BeFalse(), "If-Match should be absent when no etag is given")
		Expect(method).To(Equal(http.MethodDelete))
		Expect(path).To(Equal("/configs/v1/application-agents/gid:p"))
	})

	Describe("NewAdminClientFromCredentials", func() {
		It("signs requests with a Bearer token minted from the credential", func(ctx SpecContext) {
			// A valid EC service-account credential (test key).
			cred := `{
      "serviceAccountId": "fa50a80e-4840-4fc0-8958-982b84827f83",
      "privateKeyJWK": {"kty":"EC","d":"2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s","use":"sig","crv":"P-256",
        "x":"Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA",
        "y":"DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4","alg":"ES256"}
    }`
			var gotAuth string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotAuth = r.Header.Get("Authorization")
				_, _ = io.WriteString(w, `{"data":[]}`)
			}))
			DeferCleanup(srv.Close)

			admin, err := config.NewAdminClientFromCredentials(ctx, []byte(cred), transport.WithBaseURL(srv.URL))
			Expect(err).NotTo(HaveOccurred())
			_, err = admin.AppAgents().List(ctx, "gid:proj")
			Expect(err).NotTo(HaveOccurred())
			Expect(gotAuth).To(HavePrefix("Bearer "))
		})

		It("rejects a credential that is not JSON", func(ctx SpecContext) {
			_, err := config.NewAdminClientFromCredentials(ctx, []byte(`not-json`))
			Expect(err).To(HaveOccurred())
		})
	})
})
