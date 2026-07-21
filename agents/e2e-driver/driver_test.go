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

package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/config"
)

const ecKey = `"privateKeyJWK":{"kty":"EC","d":"2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s","use":"sig","crv":"P-256",
  "x":"Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA","y":"DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4","alg":"ES256"}`

func appAgentCred() string {
	return "e2e-test-app-agent-token"
}
func serviceAccount() []byte {
	return []byte(`{"serviceAccountId":"fa50a80e-4840-4fc0-8958-982b84827f83",` + ecKey + `}`)
}

// recorder captures the methods+paths the driver hits.
type recorder struct {
	mu   sync.Mutex
	hits []string
}

func (r *recorder) add(s string) { r.mu.Lock(); r.hits = append(r.hits, s); r.mu.Unlock() }
func (r *recorder) has(s string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Contains(r.hits, s)
}

func mockPlatform(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.add(r.Method + " " + r.URL.Path)
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/configs/v1/authorization-policies":
			w.Header().Set("ETag", `"v1"`)
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"gid:pol1","create_time":"t0"}`)
		case r.URL.Path == "/capture/v1/nodes":
			_, _ = io.WriteString(w, `{"results":[{"id":"gid:alice"},{"id":"gid:doc"}]}`)
		case r.URL.Path == "/capture/v1/relationships":
			_, _ = io.WriteString(w, `{"results":[{"id":"gid:owns"}]}`)
		case r.URL.Path == "/access/v1/evaluation":
			_, _ = io.WriteString(w, `{"decision":true}`)
		case r.URL.Path == "/contx-iq/v1/execute":
			_, _ = io.WriteString(w, `{"data":[{"nodes":{"d":{"id":"doc-1"}}}]}`)
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/configs/v1/authorization-policies/"):
			_, _ = io.WriteString(w, `{"id":"gid:pol1"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func buildDriver(t *testing.T, baseURL string) *Driver {
	t.Helper()
	admin, err := config.NewAdminClientFromCredentials(context.Background(), serviceAccount(),
		indykite.WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("admin: %v", err)
	}
	ik, err := indykite.NewClient(context.Background(), appAgentCred(), indykite.WithBaseURL(baseURL))
	if err != nil {
		t.Fatalf("runtime: %v", err)
	}
	return NewDriver(admin, ik)
}

func sampleScenario() Scenario {
	return Scenario{
		Policy: config.CreateAuthorizationPolicy{
			ProjectID: "gid:proj", Name: "can-read", Status: config.StatusActive, Policy: `{}`,
		},
		Nodes: []capture.UpsertNode{
			{Node: capture.Node{ExternalID: "alice", Type: "Person"}},
			{Node: capture.Node{ExternalID: "doc-1", Type: "Document"}},
		},
		Relationship: &capture.Relationship{
			Type:   "OWNS",
			Source: &capture.Node{ExternalID: "alice", Type: "Person"},
			Target: &capture.Node{ExternalID: "doc-1", Type: "Document"},
		},
		Subject: authzen.NewNode("Person", "alice"), Action: "READ",
		Resource: authzen.NewNode("Document", "doc-1"), WantAllow: true,
		QueryID: "get-doc",
	}
}

func TestDriverFullScenario(t *testing.T) {
	rec := &recorder{}
	srv := mockPlatform(t, rec)
	defer srv.Close()

	sc := sampleScenario()
	steps, err := buildDriver(t, srv.URL).Run(context.Background(), &sc)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !AllOK(steps) {
		t.Fatalf("not all steps passed: %+v", steps)
	}

	// Every plane/endpoint must have been exercised, including cleanup.
	for _, want := range []string{
		"POST /configs/v1/authorization-policies",
		"POST /capture/v1/nodes",
		"POST /capture/v1/relationships",
		"POST /access/v1/evaluation",
		"POST /contx-iq/v1/execute",
		"DELETE /configs/v1/authorization-policies/gid:pol1",
	} {
		if !rec.has(want) {
			t.Errorf("driver never issued %q; hits=%v", want, rec.hits)
		}
	}
	// Last step is always cleanup.
	if steps[len(steps)-1].Name != "cleanup-policy" || !steps[len(steps)-1].OK {
		t.Errorf("cleanup step missing or failed: %+v", steps)
	}
}

func TestDriverDecisionMismatchReportsFail(t *testing.T) {
	rec := &recorder{}
	srv := mockPlatform(t, rec) // evaluation returns decision=true
	defer srv.Close()

	sc := sampleScenario()
	sc.WantAllow = false // expect deny, platform says allow -> step must fail
	sc.QueryID = ""      // skip CIQ for this check

	steps, err := buildDriver(t, srv.URL).Run(context.Background(), &sc)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if AllOK(steps) {
		t.Fatal("expected the decision step to fail on mismatch")
	}
	// Cleanup must still have run despite the failed assertion.
	if !rec.has("DELETE /configs/v1/authorization-policies/gid:pol1") {
		t.Error("cleanup did not run after a failed assertion")
	}
}
