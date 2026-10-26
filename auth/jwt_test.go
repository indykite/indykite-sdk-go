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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/indykite/indykite-sdk-go/auth"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// pemPKCS8Key generates a fresh EC P-256 private key encoded as PKCS8 PEM.
func pemPKCS8Key() string {
	GinkgoHelper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Expect(err).NotTo(HaveOccurred(), "generate key")
	der, err := x509.MarshalPKCS8PrivateKey(key)
	Expect(err).NotTo(HaveOccurred(), "marshal PKCS8")
	var b strings.Builder
	Expect(pem.Encode(&b, &pem.Block{Type: "PRIVATE KEY", Bytes: der})).To(Succeed(), "encode PEM")
	return b.String()
}

// mintToken applies the authenticator to a request and returns the runtime-plane
// token it produced.
func mintToken(a *auth.Authenticator) string {
	GinkgoHelper()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", http.NoBody)
	Expect(a.Apply(context.Background(), req)).To(Succeed())
	tok := req.Header.Get(auth.HeaderClientKey)
	Expect(tok).NotTo(BeEmpty(), "no token minted")
	return tok
}

// PEM PKCS8 keys carry no "alg" hint, so the authenticator builds fine but the
// current signer refuses to mint until an algorithm is known. Cover the loading
// branches and the sign-time error propagation through Apply.
func applyPKCS8(creds *auth.Credentials) error {
	GinkgoHelper()
	a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	Expect(err).NotTo(HaveOccurred(), "New with PKCS8 key")
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", http.NoBody)
	return a.Apply(context.Background(), req)
}

// expectPKCS8Outcome accepts either a minted token or a signing-algorithm error.
func expectPKCS8Outcome(err error, kind string) {
	GinkgoHelper()
	if err == nil {
		GinkgoWriter.Printf("PKCS8 %s key minted a token\n", kind)
		return
	}
	Expect(err).To(MatchError(ContainSubstring("algorithm")), "want a signing-algorithm error")
}

type errProvider struct{}

func (errProvider) Token(context.Context) (string, error) {
	return "", errors.New("token backend down")
}

var _ = Describe("JWT minting", func() {
	It("mints a JWT with the expected payload claims and reuses it", func() {
		a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticJSON([]byte(appAgentCredJSON)))
		Expect(err).NotTo(HaveOccurred())
		tok := mintToken(a)

		parts := strings.Split(tok, ".")
		Expect(parts).To(HaveLen(3))
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		Expect(err).NotTo(HaveOccurred(), "decode payload")
		var claims struct {
			Iss string  `json:"iss"`
			Sub string  `json:"sub"`
			Jti string  `json:"jti"`
			Iat float64 `json:"iat"`
			Exp float64 `json:"exp"`
		}
		Expect(json.Unmarshal(payload, &claims)).To(Succeed())
		const wantID = "fa50a80e-4840-4fc0-8958-982b84827f83"
		Expect(claims.Iss).To(Equal(wantID))
		Expect(claims.Sub).To(Equal(wantID))
		Expect(claims.Jti).NotTo(BeEmpty(), "jti claim missing")
		// Default lifetime is one hour when none is configured.
		Expect(claims.Exp - claims.Iat).To(BeEquivalentTo(3600))

		// A second Apply must reuse the cached token (oauth2.ReuseTokenSource).
		Expect(mintToken(a)).To(Equal(tok), "token was re-minted instead of reused before expiry")
	})

	It("honours a configured token lifetime", func() {
		creds := &auth.Credentials{
			AppAgentID:    "agent-1",
			TokenLifetime: "30m",
			PrivateKeyJWK: json.RawMessage(testJWK),
		}
		a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
		Expect(err).NotTo(HaveOccurred())
		parts := strings.Split(mintToken(a), ".")
		payload, _ := base64.RawURLEncoding.DecodeString(parts[1])
		var claims struct {
			Iat float64 `json:"iat"`
			Exp float64 `json:"exp"`
		}
		Expect(json.Unmarshal(payload, &claims)).To(Succeed())
		Expect(claims.Exp - claims.Iat).To(BeEquivalentTo(1800))
	})

	It("rejects an invalid token lifetime", func() {
		creds := &auth.Credentials{
			AppAgentID:    "agent-1",
			TokenLifetime: "not-a-duration",
			PrivateKeyJWK: json.RawMessage(testJWK),
		}
		_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
		Expect(err).To(MatchError(ContainSubstring("tokenLifetime")))
	})
})

var _ = Describe("Private keys", func() {
	It("loads a PKCS8 PEM key", func() {
		creds := &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8: pemPKCS8Key()}
		expectPKCS8Outcome(applyPKCS8(creds), "PEM")
	})

	It("loads a base64 PKCS8 key", func() {
		encoded := base64.StdEncoding.EncodeToString([]byte(pemPKCS8Key()))
		creds := &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8Base64: encoded}
		expectPKCS8Outcome(applyPKCS8(creds), "base64")
	})

	DescribeTable("rejects unusable keys",
		func(creds *auth.Credentials) {
			_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
			Expect(err).To(HaveOccurred())
		},
		Entry("no key at all", &auth.Credentials{AppAgentID: "agent-1"}),
		Entry("invalid base64", &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8Base64: "%%%not-base64%%%"}),
		Entry("invalid PEM", &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8: "not a pem block"}),
		Entry("JWK not a key", &auth.Credentials{AppAgentID: "agent-1",
			PrivateKeyJWK: json.RawMessage(`{"kty":"nope"}`)}),
		Entry("JWK string not JSON", &auth.Credentials{AppAgentID: "agent-1",
			PrivateKeyJWK: json.RawMessage(`"unterminated`)}),
		Entry("JWK string not a key", &auth.Credentials{AppAgentID: "agent-1",
			PrivateKeyJWK: json.RawMessage(`"not a key"`)}),
	)

	It("accepts a JWK delivered as a JSON string", func() {
		// The JWK may be delivered as a JSON string containing the key document.
		quoted, err := json.Marshal(testJWK)
		Expect(err).NotTo(HaveOccurred(), "quote JWK")
		creds := &auth.Credentials{AppAgentID: "agent-1", PrivateKeyJWK: quoted}
		a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Split(mintToken(a), ".")).To(HaveLen(3), "token is not a JWT")
	})
})

var _ = Describe("Providers and loaders", func() {
	It("propagates the provider error through Apply", func() {
		a := auth.NewWithProvider(auth.PlaneControl, errProvider{})
		Expect(a.Plane()).To(Equal(auth.PlaneControl))
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", http.NoBody)
		Expect(a.Apply(context.Background(), req)).To(MatchError(ContainSubstring("token backend down")))
	})

	It("stops the loader chain on a loader error", func() {
		boom := errors.New("loader failed")
		failing := func(context.Context) (*auth.Credentials, error) { return nil, boom }
		_, err := auth.New(context.Background(), auth.PlaneRuntime, failing, auth.StaticJSON([]byte(appAgentCredJSON)))
		Expect(err).To(MatchError(boom))
	})
})
