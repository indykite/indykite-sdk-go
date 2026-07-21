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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "sa-token", nil }

// adminTo wires an AdminClient to the given handler using a CONTROL-PLANE
// authenticator (so requests must use Bearer auth).
func adminTo(t *testing.T, h http.HandlerFunc) *config.AdminClient {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneControl, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return config.NewAdminClient(tc)
}

func TestControlPlaneUsesBearer(t *testing.T) {
	var gotAuth, gotKey string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotKey = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
		_, _ = io.WriteString(w, `{"data":[]}`)
	})

	if _, err := admin.AuthorizationPolicies().List(context.Background(), "gid:proj", ""); err != nil {
		t.Fatalf("List: %v", err)
	}
	if gotAuth != "Bearer sa-token" {
		t.Errorf("Authorization = %q, want 'Bearer sa-token'", gotAuth)
	}
	if gotKey != "" {
		t.Errorf("X-IK-ClientKey must be empty on control plane, got %q", gotKey)
	}
}

func TestListWithQueryParams(t *testing.T) {
	var query url.Values
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		_, _ = io.WriteString(w, `{"data":[{"id":"gid:p1","name":"can-read","status":"ACTIVE","policy":"{}"}]}`)
	})

	pols, err := admin.AuthorizationPolicies().List(context.Background(), "gid:proj", "kbac", config.WithSearch("read"))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if query.Get("project_id") != "gid:proj" || query.Get("type") != "kbac" || query.Get("search") != "read" {
		t.Errorf("query = %v", query)
	}
	if len(pols) != 1 || pols[0].Name != "can-read" || pols[0].Status != "ACTIVE" {
		t.Errorf("policies = %+v", pols)
	}
}

func TestCreateCapturesETag(t *testing.T) {
	var gotBody map[string]any
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.Header().Set("ETag", `"v1"`)
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"id":"gid:newpol","create_time":"t0"}`)
	})

	res, err := admin.AuthorizationPolicies().Create(context.Background(), &config.CreateAuthorizationPolicy{
		ProjectID: "gid:proj", Name: "can-provision", Policy: `{"x":1}`, Status: config.StatusActive,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if gotBody["project_id"] != "gid:proj" || gotBody["name"] != "can-provision" || gotBody["status"] != "ACTIVE" {
		t.Errorf("body = %v", gotBody)
	}
	if res.ID != "gid:newpol" || res.ETag != `"v1"` {
		t.Errorf("result = %+v", res)
	}
}

// TestReadThenUpdateSendsIfMatch is the optimistic-concurrency round trip: Read
// returns an ETag, Update must echo it in If-Match.
func TestReadThenUpdateSendsIfMatch(t *testing.T) {
	var ifMatch string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
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
	pol, err := api.Read(context.Background(), "gid:p")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if pol.ETag != `"v7"` {
		t.Fatalf("read ETag = %q", pol.ETag)
	}

	desc := "updated"
	res, err := api.Update(context.Background(), "gid:p", pol.ETag, &config.UpdateAuthorizationPolicy{
		Description: &desc, Policy: `{"y":2}`, Status: config.StatusActive,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if ifMatch != `"v7"` {
		t.Errorf("If-Match = %q, want \"v7\"", ifMatch)
	}
	if res.ETag != `"v8"` {
		t.Errorf("update ETag = %q", res.ETag)
	}
}

func TestDeleteWithoutETagOmitsIfMatch(t *testing.T) {
	var hadIfMatch bool
	var method, path string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		_, hadIfMatch = r.Header["If-Match"]
		method, path = r.Method, r.URL.Path
		_, _ = io.WriteString(w, `{"id":"gid:p"}`)
	})

	if err := admin.AppAgents().Delete(context.Background(), "gid:p", ""); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if hadIfMatch {
		t.Error("If-Match should be absent when no etag is given")
	}
	if method != http.MethodDelete || path != "/configs/v1/application-agents/gid:p" {
		t.Errorf("method/path = %s %s", method, path)
	}
}

func TestAdminClientFromCredentials(t *testing.T) {
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
	defer srv.Close()

	admin, err := config.NewAdminClientFromCredentials(context.Background(), []byte(cred),
		transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewAdminClientFromCredentials: %v", err)
	}
	if _, err = admin.AppAgents().List(context.Background(), "gid:proj"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.HasPrefix(gotAuth, "Bearer ") {
		t.Errorf("Authorization = %q, want Bearer prefix", gotAuth)
	}
}
