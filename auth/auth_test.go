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

package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/indykite/indykite-sdk-go/auth"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// A valid EC P-256 App Agent credential (test key, not a real secret).
//
//nolint:gosec // G101: test signing key, not a real secret.
const appAgentCredJSON = `{
  "endpoint": "hera:8084",
  "appAgentId": "fa50a80e-4840-4fc0-8958-982b84827f83",
  "privateKeyJWK": {
    "kty": "EC",
    "d": "2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s",
    "use": "sig",
    "crv": "P-256",
    "kid": "vDUXHBZcRw1KyFPyB0EI2XLBzyP9iGyfvaSX3MNtUlk",
    "x": "Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA",
    "y": "DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4",
    "alg": "ES256"
  }
}`

var _ = Describe("AppAgentFromToken", func() {
	It("puts the token into the X-IK-ClientKey header", func() {
		a, err := auth.AppAgentFromToken("agent-token-abc")
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Plane()).To(Equal(auth.PlaneRuntime))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/eval", http.NoBody)
		Expect(a.Apply(context.Background(), req)).To(Succeed())

		Expect(req.Header.Get(auth.HeaderClientKey)).To(Equal("agent-token-abc"),
			"X-IK-ClientKey should carry the token verbatim")
		Expect(req.Header.Get(auth.HeaderAuthorization)).To(BeEmpty(),
			"Authorization header should be empty on the runtime plane")
	})

	It("rejects JSON", func() {
		_, err := auth.AppAgentFromToken(`{"appAgentId":"x"}`)
		Expect(err).To(HaveOccurred(),
			"the App Agent credential must be the raw token, not JSON")
	})
})

var _ = Describe("ServiceAccountFromJSON", func() {
	It("uses a Bearer Authorization header", func() {
		// Same key shape, but as a service account credential.
		saCred := strings.Replace(appAgentCredJSON,
			`"appAgentId": "fa50a80e-4840-4fc0-8958-982b84827f83",`,
			`"serviceAccountId": "fa50a80e-4840-4fc0-8958-982b84827f83",`, 1)

		a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(saCred))
		Expect(err).NotTo(HaveOccurred())

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/cfg", http.NoBody)
		Expect(a.Apply(context.Background(), req)).To(Succeed())

		Expect(req.Header.Get(auth.HeaderAuthorization)).To(HavePrefix("Bearer "))
		Expect(req.Header.Get(auth.HeaderClientKey)).To(BeEmpty(),
			"X-IK-ClientKey should be empty on the control plane")
	})
})

var _ = Describe("New", func() {
	It("fails when no loader yields credentials", func() {
		_, err := auth.New(context.Background(), auth.PlaneRuntime)
		Expect(err).To(HaveOccurred())
	})
})
