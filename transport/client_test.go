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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/transport"
)

// captureRT records the request URL and returns a canned 200 response, so the
// resolved base URL can be observed without real network access.
type captureRT struct{ url *url.URL }

func (c *captureRT) RoundTrip(r *http.Request) (*http.Response, error) {
	c.url = r.URL
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(`{}`)),
	}, nil
}

func TestWithRegionResolvesURL(t *testing.T) {
	rt := &captureRT{}
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})
	c, err := transport.NewClient(a,
		transport.WithRegion("eu"),
		transport.WithHTTPClient(&http.Client{Transport: rt}))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if err = c.Do(context.Background(), http.MethodGet, "/ping", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if rt.url.Host != "eu.api.indykite.com" || rt.url.Scheme != "https" {
		t.Errorf("resolved URL = %v, want https://eu.api.indykite.com", rt.url)
	}
}

func TestWithBaseURLTakesPrecedenceOverRegion(t *testing.T) {
	rt := &captureRT{}
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})
	a.BaseURL = "https://from-credential.example"
	c, err := transport.NewClient(a,
		transport.WithBaseURL("https://explicit.example"),
		transport.WithRegion("eu"),
		transport.WithHTTPClient(&http.Client{Transport: rt}))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if err = c.Do(context.Background(), http.MethodGet, "/ping", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if rt.url.Host != "explicit.example" {
		t.Errorf("host = %q, want explicit.example (WithBaseURL precedence)", rt.url.Host)
	}
}

func TestBaseURLFallsBackToCredential(t *testing.T) {
	rt := &captureRT{}
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})
	a.BaseURL = "https://from-credential.example"
	c, err := transport.NewClient(a, transport.WithHTTPClient(&http.Client{Transport: rt}))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if err = c.Do(context.Background(), http.MethodGet, "/ping", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if rt.url.Host != "from-credential.example" {
		t.Errorf("host = %q, want from-credential.example", rt.url.Host)
	}
}

func TestNewClientErrors(t *testing.T) {
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})

	if _, err := transport.NewClient(nil); err == nil {
		t.Error("expected error for nil authenticator")
	}
	if _, err := transport.NewClient(a); err == nil {
		t.Error("expected error when no base URL is configured anywhere")
	}
	if _, err := transport.NewClient(a, transport.WithBaseURL("://bad")); err == nil {
		t.Error("expected error for unparsable base URL")
	}
}

func TestWithUserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL, transport.WithUserAgent("custom-agent/1.0"))
	if err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotUA != "custom-agent/1.0" {
		t.Errorf("User-Agent = %q, want custom-agent/1.0", gotUA)
	}

	c = newClient(t, auth.PlaneRuntime, srv.URL)
	if err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotUA != "indykite-sdk-go" {
		t.Errorf("default User-Agent = %q, want indykite-sdk-go", gotUA)
	}
}

func TestWithTracingStillRoundTrips(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	var out struct {
		OK bool `json:"ok"`
	}
	c := newClient(t, auth.PlaneRuntime, srv.URL, transport.WithTracing())
	if err := c.Do(context.Background(), http.MethodGet, "/traced", nil, &out); err != nil {
		t.Fatalf("Do with tracing: %v", err)
	}
	if !out.OK {
		t.Error("decoded ok = false, want true")
	}

	// Tracing must also work when the client has an explicit transport.
	c = newClient(t, auth.PlaneRuntime, srv.URL,
		transport.WithHTTPClient(&http.Client{Transport: http.DefaultTransport}),
		transport.WithTracing())
	if err := c.Do(context.Background(), http.MethodGet, "/traced", nil, &out); err != nil {
		t.Fatalf("Do with tracing over custom transport: %v", err)
	}
}

func TestWithHeaderCallOptionOverrides(t *testing.T) {
	var gotIfMatch, gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotIfMatch = r.Header.Get("If-Match")
		gotUA = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL)
	err := c.Do(context.Background(), http.MethodPut, "/cfg", nil, nil,
		transport.WithHeader("If-Match", `"etag-1"`),
		transport.WithHeader("User-Agent", "override/2.0"))
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotIfMatch != `"etag-1"` {
		t.Errorf("If-Match = %q, want %q", gotIfMatch, `"etag-1"`)
	}
	if gotUA != "override/2.0" {
		t.Errorf("User-Agent = %q, want override/2.0 (call options run last)", gotUA)
	}
}

func TestDoRespExposesHeadersAndRequestID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("ETag", `"v42"`)
		w.Header().Set("X-Request-ID", "req-77")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"gid:x"}`))
	}))
	defer srv.Close()

	var out struct {
		ID string `json:"id"`
	}
	c := newClient(t, auth.PlaneRuntime, srv.URL)
	resp, err := c.DoResp(context.Background(), http.MethodPost, "/configs/v1/things", map[string]any{"n": 1}, &out)
	if err != nil {
		t.Fatalf("DoResp: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want 201", resp.StatusCode)
	}
	if resp.Header.Get("ETag") != `"v42"` {
		t.Errorf("ETag = %q, want %q", resp.Header.Get("ETag"), `"v42"`)
	}
	if resp.RequestID != "req-77" {
		t.Errorf("RequestID = %q, want req-77", resp.RequestID)
	}
	if out.ID != "gid:x" {
		t.Errorf("decoded id = %q, want gid:x", out.ID)
	}
}

func TestBodyEncodingVariants(t *testing.T) {
	var gotBody, gotCT string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		gotBody = string(data)
		gotCT = r.Header.Get("Content-Type")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL)
	ctx := context.Background()

	// []byte passes through unmodified.
	if err := c.Do(ctx, http.MethodPost, "/raw", []byte(`{"raw":1}`), nil); err != nil {
		t.Fatalf("Do []byte: %v", err)
	}
	if gotBody != `{"raw":1}` || gotCT != "application/json" {
		t.Errorf("[]byte body = %q ct = %q", gotBody, gotCT)
	}

	// io.Reader is read fully.
	if err := c.Do(ctx, http.MethodPost, "/rdr", strings.NewReader(`{"r":2}`), nil); err != nil {
		t.Fatalf("Do reader: %v", err)
	}
	if gotBody != `{"r":2}` {
		t.Errorf("reader body = %q", gotBody)
	}

	// Unmarshalable values fail before any request is made.
	if err := c.Do(ctx, http.MethodPost, "/bad", make(chan int), nil); err == nil {
		t.Error("expected encode error for chan body")
	} else if !strings.Contains(err.Error(), "encode request body") {
		t.Errorf("error = %v, want encode request body", err)
	}
}

func TestDecodeErrorOnMalformedBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer srv.Close()

	var out map[string]any
	c := newClient(t, auth.PlaneRuntime, srv.URL)
	err := c.Do(context.Background(), http.MethodGet, "/x", nil, &out)
	if err == nil || !strings.Contains(err.Error(), "decode response body") {
		t.Errorf("error = %v, want decode response body", err)
	}
}

func TestEndpointPreservesQueryString(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL)
	if err := c.Do(context.Background(), http.MethodGet, "/list?page_size=5&page_token=t1", nil, nil); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotPath != "/list" || gotQuery != "page_size=5&page_token=t1" {
		t.Errorf("path = %q query = %q", gotPath, gotQuery)
	}
}
