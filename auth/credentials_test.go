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
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
)

// testJWK is the private-key JWK from appAgentCredJSON as a raw JSON value, for
// building in-memory Credentials.
const testJWK = `{"kty":"EC","d":"2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s","use":"sig","crv":"P-256",` +
	`"x":"Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA","y":"DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4",` +
	`"alg":"ES256"}`

func clearCredentialEnv(t *testing.T) {
	t.Helper()
	for _, env := range []string{
		auth.EnvAppAgentCredentials, auth.EnvAppAgentCredentials + "_FILE",
		auth.EnvServiceAccountCredentials, auth.EnvServiceAccountCredentials + "_FILE",
	} {
		t.Setenv(env, "")
	}
}

func TestPlaneString(t *testing.T) {
	if got := auth.PlaneRuntime.String(); got != "runtime" {
		t.Errorf("PlaneRuntime.String() = %q, want runtime", got)
	}
	if got := auth.PlaneControl.String(); got != "control" {
		t.Errorf("PlaneControl.String() = %q, want control", got)
	}
}

func TestStaticJSONInvalid(t *testing.T) {
	_, err := auth.ServiceAccountFromJSON(context.Background(), []byte(`{not json`))
	if err == nil || !strings.Contains(err.Error(), "invalid credentials JSON") {
		t.Errorf("error = %v, want invalid credentials JSON", err)
	}
}

func TestStaticJSONEmptyFallsThroughToNoCredentials(t *testing.T) {
	_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticJSON(nil))
	if !errors.Is(err, auth.ErrNoCredentials) {
		t.Errorf("error = %v, want ErrNoCredentials", err)
	}
}

func TestStaticCredentialsLoaderAndBaseURL(t *testing.T) {
	creds := &auth.Credentials{
		AppAgentID:    "agent-1",
		BaseURL:       "https://base.example",
		Endpoint:      "endpoint.example:8080",
		PrivateKeyJWK: json.RawMessage(testJWK),
	}
	a, err := auth.New(context.Background(), auth.PlaneRuntime, nil, auth.StaticCredentials(creds))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if a.BaseURL != "https://base.example" {
		t.Errorf("BaseURL = %q, want baseUrl to win over endpoint", a.BaseURL)
	}
}

func TestEndpointFallbackWhenNoBaseURL(t *testing.T) {
	creds := &auth.Credentials{
		AppAgentID:    "agent-1",
		Endpoint:      "endpoint.example:8080",
		PrivateKeyJWK: json.RawMessage(testJWK),
	}
	a, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if a.BaseURL != "endpoint.example:8080" {
		t.Errorf("BaseURL = %q, want the endpoint fallback", a.BaseURL)
	}
}

func TestClientIDMissing(t *testing.T) {
	creds := &auth.Credentials{PrivateKeyJWK: json.RawMessage(testJWK)}
	_, err := auth.New(context.Background(), auth.PlaneRuntime, auth.StaticCredentials(creds))
	if err == nil || !strings.Contains(err.Error(), "missing client ID") {
		t.Errorf("error = %v, want missing client ID", err)
	}
}

// The App Agent credential is the token itself; JSON there is an error.
func TestAppAgentFromEnvRejectsJSON(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv(auth.EnvAppAgentCredentials, appAgentCredJSON)

	if _, err := auth.AppAgentFromEnv(context.Background()); err == nil {
		t.Error("expected error: app agent credential must be the raw token, not JSON")
	}
}

func TestAppAgentFromEnvToken(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv(auth.EnvAppAgentCredentials, "env-agent-token")

	a, err := auth.AppAgentFromEnv(context.Background())
	if err != nil {
		t.Fatalf("AppAgentFromEnv: %v", err)
	}
	if a.Plane() != auth.PlaneRuntime {
		t.Errorf("plane = %v, want runtime", a.Plane())
	}
	if a.BaseURL != "" {
		t.Errorf("BaseURL = %q, want empty (a token carries no URL hint)", a.BaseURL)
	}
}

func TestServiceAccountFromEnvValue(t *testing.T) {
	clearCredentialEnv(t)
	saCred := strings.Replace(appAgentCredJSON, `"appAgentId"`, `"serviceAccountId"`, 1)
	t.Setenv(auth.EnvServiceAccountCredentials, saCred)

	a, err := auth.ServiceAccountFromEnv(context.Background())
	if err != nil {
		t.Fatalf("ServiceAccountFromEnv: %v", err)
	}
	if a.Plane() != auth.PlaneControl {
		t.Errorf("plane = %v, want control", a.Plane())
	}
}

func TestServiceAccountFromEnvFile(t *testing.T) {
	clearCredentialEnv(t)
	saCred := strings.Replace(appAgentCredJSON, `"appAgentId"`, `"serviceAccountId"`, 1)
	path := filepath.Join(t.TempDir(), "sa.json")
	if err := os.WriteFile(path, []byte(saCred), 0o600); err != nil {
		t.Fatalf("write creds file: %v", err)
	}
	t.Setenv(auth.EnvServiceAccountCredentials+"_FILE", path)

	if _, err := auth.ServiceAccountFromEnv(context.Background()); err != nil {
		t.Fatalf("ServiceAccountFromEnv (file): %v", err)
	}
}

func TestEnvAbsentReturnsNoCredentials(t *testing.T) {
	clearCredentialEnv(t)
	if _, err := auth.AppAgentFromEnv(context.Background()); !errors.Is(err, auth.ErrNoCredentials) {
		t.Errorf("AppAgentFromEnv = %v, want ErrNoCredentials", err)
	}
	if _, err := auth.ServiceAccountFromEnv(context.Background()); !errors.Is(err, auth.ErrNoCredentials) {
		t.Errorf("ServiceAccountFromEnv = %v, want ErrNoCredentials", err)
	}
}

func TestEnvInvalidJSON(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv(auth.EnvServiceAccountCredentials, `{broken`)
	if _, err := auth.ServiceAccountFromEnv(context.Background()); err == nil {
		t.Error("expected error for invalid JSON in env")
	}
}

func TestEnvFileMissing(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv(auth.EnvAppAgentCredentials+"_FILE", filepath.Join(t.TempDir(), "nope.json"))
	if _, err := auth.AppAgentFromEnv(context.Background()); err == nil {
		t.Error("expected error for missing credentials file")
	}
}

func TestEnvFileTooLarge(t *testing.T) {
	clearCredentialEnv(t)
	path := filepath.Join(t.TempDir(), "big.json")
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), 11*1024), 0o600); err != nil {
		t.Fatalf("write big file: %v", err)
	}
	t.Setenv(auth.EnvAppAgentCredentials+"_FILE", path)
	_, err := auth.AppAgentFromEnv(context.Background())
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("error = %v, want size-limit error", err)
	}
}

// A credential file may hold only the pre-issued token (the raw X-IK-ClientKey
// value): no JSON, no private key.
func TestRawTokenCredential(t *testing.T) {
	clearCredentialEnv(t)
	t.Setenv(auth.EnvAppAgentCredentials, "  raw-agent-token-123\n")

	a, err := auth.AppAgentFromEnv(context.Background())
	if err != nil {
		t.Fatalf("AppAgentFromEnv: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got := req.Header.Get("X-Ik-Clientkey"); got != "raw-agent-token-123" {
		t.Errorf("X-IK-ClientKey = %q, want the trimmed raw token", got)
	}
}

func TestRawTokenCredentialFromFile(t *testing.T) {
	clearCredentialEnv(t)
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte("file-token-456\n"), 0o600); err != nil {
		t.Fatalf("write token file: %v", err)
	}
	t.Setenv(auth.EnvAppAgentCredentials+"_FILE", path)

	a, err := auth.AppAgentFromEnv(context.Background())
	if err != nil {
		t.Fatalf("AppAgentFromEnv: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got := req.Header.Get("X-Ik-Clientkey"); got != "file-token-456" {
		t.Errorf("X-IK-ClientKey = %q, want the file token", got)
	}
}

// A JSON credential carrying both a valid pre-issued token and a private key
// must send the token verbatim (no minting).
func TestPreIssuedTokenPreferredOverKey(t *testing.T) {
	clearCredentialEnv(t)
	cred := `{"serviceAccountId":"gid:sa","token":"opaque-pre-issued","privateKeyJWK":` + testJWK + `}`
	a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(cred))
	if err != nil {
		t.Fatalf("ServiceAccountFromJSON: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer opaque-pre-issued" {
		t.Errorf("Authorization = %q, want the pre-issued token", got)
	}
}

// An expired pre-issued JWT with a private key available falls back to minting.
func TestExpiredPreIssuedTokenFallsBackToMinting(t *testing.T) {
	clearCredentialEnv(t)
	// exp in this JWT payload is 1600000000 (2020) — long expired; signature
	// content is irrelevant, only the exp claim is inspected.
	expired := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9." +
		base64.RawURLEncoding.EncodeToString([]byte(`{"exp":1600000000,"sub":"gid:sa"}`)) +
		".c2ln"
	cred := `{"serviceAccountId":"gid:sa","token":"` + expired + `","privateKeyJWK":` + testJWK + `}`
	a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(cred))
	if err != nil {
		t.Fatalf("ServiceAccountFromJSON: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	got := req.Header.Get("Authorization")
	if got == "Bearer "+expired {
		t.Error("expired pre-issued token was sent; expected a freshly minted JWT")
	}
	if !strings.HasPrefix(got, "Bearer ey") {
		t.Errorf("Authorization = %q, want a minted JWT", got)
	}
}

// A pre-issued JWT within the refresh margin (expires in under a minute) is
// re-minted from the key rather than served, so a long-running client refreshes
// before the platform would reject the token.
func TestPreIssuedTokenWithinMarginRefreshes(t *testing.T) {
	clearCredentialEnv(t)
	soon := fmt.Sprintf(`{"exp":%d,"sub":"gid:sa"}`, time.Now().Add(30*time.Second).Unix())
	soonJWT := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9." +
		base64.RawURLEncoding.EncodeToString([]byte(soon)) + ".c2ln"
	cred := `{"serviceAccountId":"gid:sa","token":"` + soonJWT + `","privateKeyJWK":` + testJWK + `}`

	a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(cred))
	if err != nil {
		t.Fatalf("ServiceAccountFromJSON: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	got := req.Header.Get("Authorization")
	if got == "Bearer "+soonJWT {
		t.Error("token within the refresh margin was served; expected a freshly minted JWT")
	}
	if !strings.HasPrefix(got, "Bearer ey") {
		t.Errorf("Authorization = %q, want a minted JWT", got)
	}
}

// An expired token WITHOUT a key is still sent (the platform reports the 401).
func TestExpiredTokenWithoutKeyIsSent(t *testing.T) {
	clearCredentialEnv(t)
	expired := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9." +
		base64.RawURLEncoding.EncodeToString([]byte(`{"exp":1600000000}`)) +
		".c2ln"
	t.Setenv(auth.EnvAppAgentCredentials, expired)

	a, err := auth.AppAgentFromEnv(context.Background())
	if err != nil {
		t.Fatalf("AppAgentFromEnv: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://x/y", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got := req.Header.Get("X-Ik-Clientkey"); got != expired {
		t.Errorf("X-IK-ClientKey = %q, want the token passed through", got)
	}
}
