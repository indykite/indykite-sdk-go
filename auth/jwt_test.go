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
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
)

// pemPKCS8Key generates a fresh EC P-256 private key encoded as PKCS8 PEM.
func pemPKCS8Key(t *testing.T) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal PKCS8: %v", err)
	}
	var b strings.Builder
	if err = pem.Encode(&b, &pem.Block{Type: "PRIVATE KEY", Bytes: der}); err != nil {
		t.Fatalf("encode PEM: %v", err)
	}
	return b.String()
}

// mintToken applies the authenticator to a request and returns the runtime-plane
// token it produced.
func mintToken(t *testing.T, a *auth.Authenticator) string {
	t.Helper()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", http.NoBody)
	if err := a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	tok := req.Header.Get(auth.HeaderClientKey)
	if tok == "" {
		t.Fatal("no token minted")
	}
	return tok
}

func TestMintedJWTPayloadClaims(t *testing.T) {
	a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticJSON([]byte(appAgentCredJSON)))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	tok := mintToken(t, a)

	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("token has %d segments, want 3", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	var claims struct {
		Iss string  `json:"iss"`
		Sub string  `json:"sub"`
		Jti string  `json:"jti"`
		Iat float64 `json:"iat"`
		Exp float64 `json:"exp"`
	}
	if err = json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}
	const wantID = "fa50a80e-4840-4fc0-8958-982b84827f83"
	if claims.Iss != wantID || claims.Sub != wantID {
		t.Errorf("iss = %q sub = %q, want %q", claims.Iss, claims.Sub, wantID)
	}
	if claims.Jti == "" {
		t.Error("jti claim missing")
	}
	// Default lifetime is one hour when none is configured.
	if got := claims.Exp - claims.Iat; got != 3600 {
		t.Errorf("exp-iat = %v, want 3600", got)
	}

	// A second Apply must reuse the cached token (oauth2.ReuseTokenSource).
	if again := mintToken(t, a); again != tok {
		t.Error("token was re-minted instead of reused before expiry")
	}
}

func TestTokenLifetimeConfigurable(t *testing.T) {
	creds := &auth.Credentials{
		AppAgentID:    "agent-1",
		TokenLifetime: "30m",
		PrivateKeyJWK: json.RawMessage(testJWK),
	}
	a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	parts := strings.Split(mintToken(t, a), ".")
	payload, _ := base64.RawURLEncoding.DecodeString(parts[1])
	var claims struct {
		Iat float64 `json:"iat"`
		Exp float64 `json:"exp"`
	}
	if err = json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}
	if got := claims.Exp - claims.Iat; got != 1800 {
		t.Errorf("exp-iat = %v, want 1800", got)
	}
}

func TestTokenLifetimeInvalid(t *testing.T) {
	creds := &auth.Credentials{
		AppAgentID:    "agent-1",
		TokenLifetime: "not-a-duration",
		PrivateKeyJWK: json.RawMessage(testJWK),
	}
	_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	if err == nil || !strings.Contains(err.Error(), "tokenLifetime") {
		t.Errorf("error = %v, want tokenLifetime parse error", err)
	}
}

// PEM PKCS8 keys carry no "alg" hint, so the authenticator builds fine but the
// current signer refuses to mint until an algorithm is known. Cover the loading
// branches and the sign-time error propagation through Apply.
func applyPKCS8(t *testing.T, creds *auth.Credentials) error {
	t.Helper()
	a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	if err != nil {
		t.Fatalf("New with PKCS8 key: %v", err)
	}
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", http.NoBody)
	return a.Apply(context.Background(), req)
}

func TestPrivateKeyPKCS8PEM(t *testing.T) {
	creds := &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8: pemPKCS8Key(t)}
	if err := applyPKCS8(t, creds); err == nil {
		t.Log("PKCS8 PEM key minted a token")
	} else if !strings.Contains(err.Error(), "algorithm") {
		t.Errorf("Apply error = %v, want a signing-algorithm error", err)
	}
}

func TestPrivateKeyPKCS8Base64(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte(pemPKCS8Key(t)))
	creds := &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8Base64: encoded}
	if err := applyPKCS8(t, creds); err == nil {
		t.Log("PKCS8 base64 key minted a token")
	} else if !strings.Contains(err.Error(), "algorithm") {
		t.Errorf("Apply error = %v, want a signing-algorithm error", err)
	}
}

func TestPrivateKeyErrors(t *testing.T) {
	tests := []struct {
		name  string
		creds *auth.Credentials
	}{
		{"no key at all", &auth.Credentials{AppAgentID: "agent-1"}},
		{"invalid base64", &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8Base64: "%%%not-base64%%%"}},
		{"invalid PEM", &auth.Credentials{AppAgentID: "agent-1", PrivateKeyPKCS8: "not a pem block"}},
		{"JWK not a key", &auth.Credentials{AppAgentID: "agent-1",
			PrivateKeyJWK: json.RawMessage(`{"kty":"nope"}`)}},
		{"JWK string not JSON", &auth.Credentials{AppAgentID: "agent-1",
			PrivateKeyJWK: json.RawMessage(`"unterminated`)}},
		{"JWK string not a key", &auth.Credentials{AppAgentID: "agent-1",
			PrivateKeyJWK: json.RawMessage(`"not a key"`)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(tc.creds))
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestPrivateKeyJWKAsJSONString(t *testing.T) {
	// The JWK may be delivered as a JSON string containing the key document.
	quoted, err := json.Marshal(testJWK)
	if err != nil {
		t.Fatalf("quote JWK: %v", err)
	}
	creds := &auth.Credentials{AppAgentID: "agent-1", PrivateKeyJWK: quoted}
	a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	if err != nil {
		t.Fatalf("New with quoted JWK: %v", err)
	}
	if parts := strings.Split(mintToken(t, a), "."); len(parts) != 3 {
		t.Errorf("token is not a JWT (%d segments)", len(parts))
	}
}

type errProvider struct{}

func (errProvider) Token(context.Context) (string, error) {
	return "", errors.New("token backend down")
}

func TestApplyPropagatesProviderError(t *testing.T) {
	a := auth.NewWithProvider(auth.PlaneControl, errProvider{})
	if a.Plane() != auth.PlaneControl {
		t.Errorf("plane = %v, want control", a.Plane())
	}
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", http.NoBody)
	err := a.Apply(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "token backend down") {
		t.Errorf("Apply error = %v, want provider error", err)
	}
}

func TestLoaderErrorStopsChain(t *testing.T) {
	boom := errors.New("loader failed")
	failing := func(context.Context) (*auth.Credentials, error) { return nil, boom }
	_, err := auth.New(context.Background(), auth.PlaneRuntime, failing, auth.StaticJSON([]byte(appAgentCredJSON)))
	if !errors.Is(err, boom) {
		t.Errorf("error = %v, want the loader error", err)
	}
}
