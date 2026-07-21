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

package transport_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/transport"
)

// stubProvider is a TokenProvider that returns a fixed token (no signing key
// needed for transport-level tests).
type stubProvider struct{ tok string }

func (s stubProvider) Token(context.Context) (string, error) { return s.tok, nil }

func newClient(t *testing.T, plane auth.Plane, baseURL string, opts ...transport.Option) *transport.Client {
	t.Helper()
	a := auth.NewWithProvider(plane, stubProvider{tok: "test-token"})
	all := append([]transport.Option{transport.WithBaseURL(baseURL)}, opts...)
	c, err := transport.NewClient(a, all...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestRuntimePlaneHeaderAndJSONRoundTrip(t *testing.T) {
	var gotKey, gotAuth, gotReqID, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
		gotAuth = r.Header.Get("Authorization")
		gotReqID = r.Header.Get("X-Request-ID")
		buf := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buf)
		gotBody = string(buf)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"decision":true}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL)
	var out struct {
		Decision bool `json:"decision"`
	}
	in := map[string]any{"action": "READ"}
	if err := c.Do(context.Background(), http.MethodPost, "/access/v1/evaluation", in, &out); err != nil {
		t.Fatalf("Do: %v", err)
	}

	if gotKey != "test-token" {
		t.Errorf("X-IK-ClientKey = %q, want test-token", gotKey)
	}
	if gotAuth != "" {
		t.Errorf("Authorization should be empty on runtime plane, got %q", gotAuth)
	}
	if gotReqID == "" {
		t.Error("X-Request-ID not set")
	}
	if gotBody != `{"action":"READ"}` {
		t.Errorf("request body = %q", gotBody)
	}
	if !out.Decision {
		t.Error("decoded decision = false, want true")
	}
}

func TestControlPlaneUsesBearer(t *testing.T) {
	var gotAuth, gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotKey = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneControl, srv.URL)
	if err := c.Do(context.Background(), http.MethodGet, "/configs/v1/projects", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotAuth != "Bearer test-token" {
		t.Errorf("Authorization = %q, want 'Bearer test-token'", gotAuth)
	}
	if gotKey != "" {
		t.Errorf("X-IK-ClientKey should be empty on control plane, got %q", gotKey)
	}
}

func TestRetryOn503ThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL, transport.WithRetry(transport.RetryConfig{
		MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
	}))
	if err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if got := calls.Load(); got != 2 {
		t.Errorf("server called %d times, want 2 (one retry)", got)
	}
}

func TestAPIErrorOn404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"NOT_FOUND","message":"no such thing"}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL)
	err := c.Do(context.Background(), http.MethodGet, "/missing", nil, nil)
	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if !apiErr.IsNotFound() {
		t.Errorf("IsNotFound = false, status=%d", apiErr.StatusCode)
	}
	if apiErr.Code != "NOT_FOUND" || apiErr.Message != "no such thing" || apiErr.RequestID != "req-123" {
		t.Errorf("APIError fields = %+v", apiErr)
	}
}

func TestPaginator(t *testing.T) {
	// Two pages of two items each, then empty token.
	fetch := func(_ context.Context, token string) (transport.Page[int], error) {
		switch token {
		case "":
			return transport.Page[int]{Items: []int{1, 2}, NextToken: "p2"}, nil
		case "p2":
			return transport.Page[int]{Items: []int{3, 4}, NextToken: ""}, nil
		default:
			return transport.Page[int]{}, fmt.Errorf("unexpected token %q", token)
		}
	}

	it := transport.NewIterator(fetch)
	got, err := it.Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	want := []int{1, 2, 3, 4}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("paginator collected %v, want %v", got, want)
	}
}
