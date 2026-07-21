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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/ciq"
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

func TestRuntimeFacadeRoutesAndAuth(t *testing.T) {
	var paths []string
	var clientKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer srv.Close()

	cli, err := indykite.NewClient(context.Background(), appAgentCred(), indykite.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ok, err := cli.AuthZEN().Allowed(context.Background(),
		authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
	if err != nil || !ok {
		t.Fatalf("Allowed: ok=%v err=%v", ok, err)
	}

	rows, err := cli.CIQ().All(context.Background(), ciq.ExecuteRequest{ID: "q", PageSize: 100})
	if err != nil {
		t.Fatalf("CIQ.All: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("rows = %d, want 1", len(rows))
	}

	if clientKey == "" {
		t.Error("runtime plane should send the App Agent client key")
	}
	joined := strings.Join(paths, ",")
	if !strings.Contains(joined, "/access/v1/evaluation") || !strings.Contains(joined, "/contx-iq/v1/execute") {
		t.Errorf("paths routed = %v", paths)
	}
}

func TestAdminFacadeUsesBearer(t *testing.T) {
	var authz, path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz = r.Header.Get("Authorization")
		path = r.URL.Path
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	defer srv.Close()

	admin, err := indykite.NewAdmin(context.Background(), serviceAccountCred(), indykite.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewAdmin: %v", err)
	}
	if _, err = admin.AuthorizationPolicies().List(context.Background(), "gid:proj", ""); err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.HasPrefix(authz, "Bearer ") {
		t.Errorf("Authorization = %q, want Bearer prefix", authz)
	}
	if path != "/configs/v1/authorization-policies" {
		t.Errorf("path = %q", path)
	}
}

func TestMissingCredentialsFails(t *testing.T) {
	if _, err := indykite.NewClient(context.Background(), "", indykite.WithRegion("eu")); err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestNewClientErrorPaths(t *testing.T) {
	ctx := context.Background()
	// A JSON document is not a valid App Agent credential token.
	if _, err := indykite.NewClient(ctx, `{"appAgentId":"x"}`, indykite.WithRegion("eu")); err == nil {
		t.Error("expected error for JSON in place of the credential token")
	}
	// The token carries no base URL hint; without an option the transport errors.
	if _, err := indykite.NewClient(ctx, appAgentCred()); err == nil {
		t.Error("expected error when no base URL can be resolved")
	}
}

func TestAccessorsNonNil(t *testing.T) {
	cli, err := indykite.NewClient(context.Background(), appAgentCred(), indykite.WithRegion("eu"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if cli.AuthZEN() == nil || cli.CIQ() == nil || cli.Capture() == nil || cli.EntityMatching() == nil {
		t.Error("all service accessors must be non-nil")
	}
}

func TestNewClientFromEnv(t *testing.T) {
	t.Setenv(envAppCreds, appAgentCred())
	cli, err := indykite.NewClientFromEnv(context.Background(), indykite.WithRegion("eu"))
	if err != nil {
		t.Fatalf("NewClientFromEnv: %v", err)
	}
	if cli.AuthZEN() == nil {
		t.Error("AuthZEN accessor is nil")
	}

	t.Setenv(envAppCreds, `{broken`)
	if _, err = indykite.NewClientFromEnv(context.Background(), indykite.WithRegion("eu")); err == nil {
		t.Error("expected error for invalid credentials in env")
	}
}

func TestNewAdminErrorPaths(t *testing.T) {
	if _, err := indykite.NewAdmin(context.Background(), []byte(`{broken`), indykite.WithRegion("eu")); err == nil {
		t.Error("expected error for invalid credentials JSON")
	}
}

func TestNewAdminFromEnv(t *testing.T) {
	t.Setenv(envSACreds, string(serviceAccountCred()))
	admin, err := indykite.NewAdminFromEnv(context.Background(), indykite.WithRegion("eu"))
	if err != nil {
		t.Fatalf("NewAdminFromEnv: %v", err)
	}
	if admin == nil {
		t.Fatal("admin client is nil")
	}

	// Credential without baseUrl and no option -> transport error.
	if _, err = indykite.NewAdminFromEnv(context.Background()); err == nil {
		t.Error("expected error when no base URL can be resolved")
	}

	t.Setenv(envSACreds, `{broken`)
	if _, err = indykite.NewAdminFromEnv(context.Background(), indykite.WithRegion("eu")); err == nil {
		t.Error("expected error for invalid credentials in env")
	}
}

func TestOptionsPassedThroughEndToEnd(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = io.WriteString(w, `{"decision":true}`)
	}))
	defer srv.Close()

	cli, err := indykite.NewClient(context.Background(), appAgentCred(),
		indykite.WithBaseURL(srv.URL),
		indykite.WithUserAgent("facade-test/1.0"),
		indykite.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
		indykite.WithRetry(indykite.RetryConfig{
			MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond,
		}),
		indykite.WithTracing())
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ok, err := cli.AuthZEN().Allowed(context.Background(),
		authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
	if err != nil || !ok {
		t.Fatalf("Allowed: ok=%v err=%v", ok, err)
	}
	if gotUA != "facade-test/1.0" {
		t.Errorf("User-Agent = %q, want facade-test/1.0 (WithUserAgent passed through)", gotUA)
	}
}
