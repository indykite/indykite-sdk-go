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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/test/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// testJWK is the private-key JWK from appAgentCredJSON as a raw JSON value, for
// building in-memory Credentials.
const testJWK = `{"kty":"EC","d":"2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s","use":"sig","crv":"P-256",` +
	`"x":"Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA","y":"DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4",` +
	`"alg":"ES256"}`

func clearCredentialEnv() {
	GinkgoHelper()
	for _, env := range []string{
		auth.EnvAppAgentCredentials, auth.EnvAppAgentCredentials + "_FILE",
		auth.EnvServiceAccountCredentials, auth.EnvServiceAccountCredentials + "_FILE",
	} {
		testutil.SetEnv(env, "")
	}
}

// writeTempFile writes data to name inside a per-spec temp directory and
// returns its path.
func writeTempFile(name string, data []byte) string {
	GinkgoHelper()
	path := filepath.Join(GinkgoT().TempDir(), name)
	Expect(os.WriteFile(path, data, 0o600)).To(Succeed())
	return path
}

// applyToRequest applies a to a fresh request and returns its headers.
func applyToRequest(a *auth.Authenticator) http.Header {
	GinkgoHelper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	Expect(err).NotTo(HaveOccurred())
	Expect(a.Apply(context.Background(), req)).To(Succeed())
	return req.Header
}

// fakeJWT builds an unsigned JWT-shaped token around payload; signature content
// is irrelevant, only the exp claim is inspected.
func fakeJWT(payload string) string {
	return "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9." +
		base64.RawURLEncoding.EncodeToString([]byte(payload)) + ".c2ln"
}

var _ = Describe("Plane", func() {
	It("stringifies", func() {
		Expect(auth.PlaneRuntime.String()).To(Equal("runtime"))
		Expect(auth.PlaneControl.String()).To(Equal("control"))
	})
})

var _ = Describe("Static credentials", func() {
	It("rejects invalid JSON", func() {
		_, err := auth.ServiceAccountFromJSON(context.Background(), []byte(`{not json`))
		Expect(err).To(MatchError(ContainSubstring("invalid credentials JSON")))
	})

	It("falls through to ErrNoCredentials on empty JSON", func() {
		_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticJSON(nil))
		Expect(err).To(MatchError(auth.ErrNoCredentials))
	})

	It("loads StaticCredentials and prefers baseUrl over endpoint", func() {
		creds := &auth.Credentials{
			AppAgentID:    "agent-1",
			BaseURL:       "https://base.example",
			Endpoint:      "endpoint.example:8080",
			PrivateKeyJWK: json.RawMessage(testJWK),
		}
		a, err := auth.New(context.Background(), auth.PlaneRuntime, nil, auth.StaticCredentials(creds))
		Expect(err).NotTo(HaveOccurred())
		Expect(a.BaseURL).To(Equal("https://base.example"), "baseUrl should win over endpoint")
	})

	It("falls back to the endpoint when there is no baseUrl", func() {
		creds := &auth.Credentials{
			AppAgentID:    "agent-1",
			Endpoint:      "endpoint.example:8080",
			PrivateKeyJWK: json.RawMessage(testJWK),
		}
		a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
		Expect(err).NotTo(HaveOccurred())
		Expect(a.BaseURL).To(Equal("endpoint.example:8080"))
	})

	It("fails when the client ID is missing", func() {
		creds := &auth.Credentials{PrivateKeyJWK: json.RawMessage(testJWK)}
		_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
		Expect(err).To(MatchError(ContainSubstring("missing client ID")))
	})
})

var _ = Describe("Environment credentials", Serial, func() {
	BeforeEach(func() {
		clearCredentialEnv()
	})

	// The App Agent credential is the token itself; JSON there is an error.
	It("rejects JSON as the App Agent credential", func() {
		testutil.SetEnv(auth.EnvAppAgentCredentials, appAgentCredJSON)

		_, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).To(HaveOccurred(), "app agent credential must be the raw token, not JSON")
	})

	It("loads the App Agent token", func() {
		testutil.SetEnv(auth.EnvAppAgentCredentials, "env-agent-token")

		a, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Plane()).To(Equal(auth.PlaneRuntime))
		Expect(a.BaseURL).To(BeEmpty(), "a token carries no URL hint")
	})

	It("loads a service account from the env value", func() {
		saCred := strings.Replace(appAgentCredJSON, `"appAgentId"`, `"serviceAccountId"`, 1)
		testutil.SetEnv(auth.EnvServiceAccountCredentials, saCred)

		a, err := auth.ServiceAccountFromEnv(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Plane()).To(Equal(auth.PlaneControl))
	})

	It("loads a service account from the _FILE env", func() {
		saCred := strings.Replace(appAgentCredJSON, `"appAgentId"`, `"serviceAccountId"`, 1)
		testutil.SetEnv(auth.EnvServiceAccountCredentials+"_FILE", writeTempFile("sa.json", []byte(saCred)))

		_, err := auth.ServiceAccountFromEnv(context.Background())
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns ErrNoCredentials when the env is absent", func() {
		_, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).To(MatchError(auth.ErrNoCredentials))
		_, err = auth.ServiceAccountFromEnv(context.Background())
		Expect(err).To(MatchError(auth.ErrNoCredentials))
	})

	It("rejects invalid JSON in env", func() {
		testutil.SetEnv(auth.EnvServiceAccountCredentials, `{broken`)
		_, err := auth.ServiceAccountFromEnv(context.Background())
		Expect(err).To(HaveOccurred())
	})

	It("fails on a missing credentials file", func() {
		testutil.SetEnv(auth.EnvAppAgentCredentials+"_FILE", filepath.Join(GinkgoT().TempDir(), "nope.json"))
		_, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).To(HaveOccurred())
	})

	It("enforces the credentials file size limit", func() {
		oversized := writeTempFile("big.json", bytes.Repeat([]byte("x"), 11*1024))
		testutil.SetEnv(auth.EnvAppAgentCredentials+"_FILE", oversized)
		_, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).To(MatchError(ContainSubstring("exceeds")))
	})

	// A credential file may hold only the pre-issued token (the raw X-IK-ClientKey
	// value): no JSON, no private key.
	It("sends a trimmed raw token credential", func() {
		testutil.SetEnv(auth.EnvAppAgentCredentials, "  raw-agent-token-123\n")

		a, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(applyToRequest(a).Get("X-Ik-Clientkey")).To(Equal("raw-agent-token-123"))
	})

	It("sends a raw token credential read from a file", func() {
		testutil.SetEnv(auth.EnvAppAgentCredentials+"_FILE", writeTempFile("token", []byte("file-token-456\n")))

		a, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(applyToRequest(a).Get("X-Ik-Clientkey")).To(Equal("file-token-456"))
	})

	// An expired token WITHOUT a key is still sent (the platform reports the 401).
	It("sends an expired token when there is no key", func() {
		expired := fakeJWT(`{"exp":1600000000}`)
		testutil.SetEnv(auth.EnvAppAgentCredentials, expired)

		a, err := auth.AppAgentFromEnv(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(applyToRequest(a).Get("X-Ik-Clientkey")).To(Equal(expired), "token should be passed through")
	})
})

var _ = Describe("Pre-issued tokens", Serial, func() {
	BeforeEach(func() {
		clearCredentialEnv()
	})

	// A JSON credential carrying both a valid pre-issued token and a private key
	// must send the token verbatim (no minting).
	It("prefers a valid pre-issued token over the key", func() {
		cred := `{"serviceAccountId":"gid:sa","token":"opaque-pre-issued","privateKeyJWK":` + testJWK + `}`
		a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(cred))
		Expect(err).NotTo(HaveOccurred())
		Expect(applyToRequest(a).Get("Authorization")).To(Equal("Bearer opaque-pre-issued"))
	})

	// An expired pre-issued JWT with a private key available falls back to minting.
	It("falls back to minting when the pre-issued token expired", func() {
		// exp in this JWT payload is 1600000000 (2020) — long expired; signature
		// content is irrelevant, only the exp claim is inspected.
		expired := fakeJWT(`{"exp":1600000000,"sub":"gid:sa"}`)
		cred := `{"serviceAccountId":"gid:sa","token":"` + expired + `","privateKeyJWK":` + testJWK + `}`
		a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(cred))
		Expect(err).NotTo(HaveOccurred())

		got := applyToRequest(a).Get("Authorization")
		Expect(got).NotTo(Equal("Bearer "+expired),
			"expired pre-issued token was sent; expected a freshly minted JWT")
		Expect(got).To(HavePrefix("Bearer ey"), "want a minted JWT")
	})

	// A pre-issued JWT within the refresh margin (expires in under a minute) is
	// re-minted from the key rather than served, so a long-running client refreshes
	// before the platform would reject the token.
	It("re-mints a pre-issued token within the refresh margin", func() {
		soonJWT := fakeJWT(fmt.Sprintf(`{"exp":%d,"sub":"gid:sa"}`, time.Now().Add(30*time.Second).Unix()))
		cred := `{"serviceAccountId":"gid:sa","token":"` + soonJWT + `","privateKeyJWK":` + testJWK + `}`

		a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(cred))
		Expect(err).NotTo(HaveOccurred())

		got := applyToRequest(a).Get("Authorization")
		Expect(got).NotTo(Equal("Bearer "+soonJWT),
			"token within the refresh margin was served; expected a freshly minted JWT")
		Expect(got).To(HavePrefix("Bearer ey"), "want a minted JWT")
	})
})
