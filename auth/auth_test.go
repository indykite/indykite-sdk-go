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
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
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

func TestAppAgentTokenIntoClientKeyHeader(t *testing.T) {
	a, err := auth.AppAgentFromToken("agent-token-abc")
	if err != nil {
		t.Fatalf("AppAgentFromToken: %v", err)
	}
	if a.Plane() != auth.PlaneRuntime {
		t.Errorf("plane = %v, want runtime", a.Plane())
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/eval", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if tok := req.Header.Get(auth.HeaderClientKey); tok != "agent-token-abc" {
		t.Errorf("X-IK-ClientKey = %q, want the token verbatim", tok)
	}
	if req.Header.Get(auth.HeaderAuthorization) != "" {
		t.Error("Authorization header should be empty on the runtime plane")
	}
}

func TestAppAgentTokenRejectsJSON(t *testing.T) {
	if _, err := auth.AppAgentFromToken(`{"appAgentId":"x"}`); err == nil {
		t.Error("expected error: the App Agent credential must be the raw token, not JSON")
	}
}

func TestServiceAccountUsesBearer(t *testing.T) {
	// Same key shape, but as a service account credential.
	saCred := strings.Replace(appAgentCredJSON,
		`"appAgentId": "fa50a80e-4840-4fc0-8958-982b84827f83",`,
		`"serviceAccountId": "fa50a80e-4840-4fc0-8958-982b84827f83",`, 1)

	a, err := auth.ServiceAccountFromJSON(context.Background(), []byte(saCred))
	if err != nil {
		t.Fatalf("ServiceAccountFromJSON: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/cfg", http.NoBody)
	if err = a.Apply(context.Background(), req); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	authz := req.Header.Get(auth.HeaderAuthorization)
	if !strings.HasPrefix(authz, "Bearer ") {
		t.Errorf("Authorization = %q, want Bearer prefix", authz)
	}
	if req.Header.Get(auth.HeaderClientKey) != "" {
		t.Error("X-IK-ClientKey should be empty on the control plane")
	}
}

func TestMissingCredentials(t *testing.T) {
	_, err := auth.New(context.Background(), auth.PlaneRuntime)
	if err == nil {
		t.Fatal("expected error when no loader yields credentials")
	}
}
