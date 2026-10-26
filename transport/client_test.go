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
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/test/testutil"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

var _ = Describe("Client", func() {
	Context("base URL resolution", func() {
		It("resolves the URL from WithRegion", func(ctx SpecContext) {
			rt := &captureRT{}
			a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})
			c, err := transport.NewClient(a,
				transport.WithRegion("eu"),
				transport.WithHTTPClient(&http.Client{Transport: rt}))
			Expect(err).NotTo(HaveOccurred())
			Expect(c.Do(ctx, http.MethodGet, "/ping", nil, nil)).To(Succeed())
			Expect(rt.url.Scheme).To(Equal("https"))
			Expect(rt.url.Host).To(Equal("eu.api.indykite.com"))
		})

		It("gives WithBaseURL precedence over WithRegion and the credential", func(ctx SpecContext) {
			rt := &captureRT{}
			a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})
			a.BaseURL = "https://from-credential.example"
			c, err := transport.NewClient(a,
				transport.WithBaseURL("https://explicit.example"),
				transport.WithRegion("eu"),
				transport.WithHTTPClient(&http.Client{Transport: rt}))
			Expect(err).NotTo(HaveOccurred())
			Expect(c.Do(ctx, http.MethodGet, "/ping", nil, nil)).To(Succeed())
			Expect(rt.url.Host).To(Equal("explicit.example"), "WithBaseURL precedence")
		})

		It("falls back to the credential's base URL", func(ctx SpecContext) {
			rt := &captureRT{}
			a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})
			a.BaseURL = "https://from-credential.example"
			c, err := transport.NewClient(a, transport.WithHTTPClient(&http.Client{Transport: rt}))
			Expect(err).NotTo(HaveOccurred())
			Expect(c.Do(ctx, http.MethodGet, "/ping", nil, nil)).To(Succeed())
			Expect(rt.url.Host).To(Equal("from-credential.example"))
		})
	})

	It("rejects invalid construction", func() {
		a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{tok: "t"})

		_, err := transport.NewClient(nil)
		Expect(err).To(HaveOccurred(), "expected error for nil authenticator")
		_, err = transport.NewClient(a)
		Expect(err).To(HaveOccurred(), "expected error when no base URL is configured anywhere")
		_, err = transport.NewClient(a, transport.WithBaseURL("://bad"))
		Expect(err).To(HaveOccurred(), "expected error for unparsable base URL")
	})

	It("sends a custom or default User-Agent", func(ctx SpecContext) {
		var gotUA string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			gotUA = r.Header.Get("User-Agent")
			_, _ = w.Write([]byte(`{}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL, transport.WithUserAgent("custom-agent/1.0"))
		Expect(c.Do(ctx, http.MethodGet, "/x", nil, nil)).To(Succeed())
		Expect(gotUA).To(Equal("custom-agent/1.0"))

		c = newClient(auth.PlaneRuntime, srv.URL)
		Expect(c.Do(ctx, http.MethodGet, "/x", nil, nil)).To(Succeed())
		Expect(gotUA).To(Equal("indykite-sdk-go"), "default User-Agent")
	})

	It("still round-trips with tracing enabled", func(ctx SpecContext) {
		srv := testutil.NewServer(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"ok":true}`))
		})

		var out struct {
			OK bool `json:"ok"`
		}
		c := newClient(auth.PlaneRuntime, srv.URL, transport.WithTracing())
		Expect(c.Do(ctx, http.MethodGet, "/traced", nil, &out)).To(Succeed(), "Do with tracing")
		Expect(out.OK).To(BeTrue())

		// Tracing must also work when the client has an explicit transport.
		c = newClient(auth.PlaneRuntime, srv.URL,
			transport.WithHTTPClient(&http.Client{Transport: http.DefaultTransport}),
			transport.WithTracing())
		Expect(c.Do(ctx, http.MethodGet, "/traced", nil, &out)).To(Succeed(),
			"Do with tracing over custom transport")
	})

	It("lets WithHeader call options override headers", func(ctx SpecContext) {
		var gotIfMatch, gotUA string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			gotIfMatch = r.Header.Get("If-Match")
			gotUA = r.Header.Get("User-Agent")
			_, _ = w.Write([]byte(`{}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL)
		err := c.Do(ctx, http.MethodPut, "/cfg", nil, nil,
			transport.WithHeader("If-Match", `"etag-1"`),
			transport.WithHeader("User-Agent", "override/2.0"))
		Expect(err).NotTo(HaveOccurred())
		Expect(gotIfMatch).To(Equal(`"etag-1"`))
		Expect(gotUA).To(Equal("override/2.0"), "call options run last")
	})

	It("exposes headers and request ID through DoResp", func(ctx SpecContext) {
		srv := testutil.NewServer(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("ETag", `"v42"`)
			w.Header().Set("X-Request-ID", "req-77")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"gid:x"}`))
		})

		var out struct {
			ID string `json:"id"`
		}
		c := newClient(auth.PlaneRuntime, srv.URL)
		resp, err := c.DoResp(ctx, http.MethodPost, "/configs/v1/things", map[string]any{"n": 1}, &out)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		Expect(resp.Header.Get("ETag")).To(Equal(`"v42"`))
		Expect(resp.RequestID).To(Equal("req-77"))
		Expect(out.ID).To(Equal("gid:x"))
	})

	It("encodes each supported body variant", func(ctx SpecContext) {
		var gotBody, gotCT string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			data, _ := io.ReadAll(r.Body)
			gotBody = string(data)
			gotCT = r.Header.Get("Content-Type")
			_, _ = w.Write([]byte(`{}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL)

		// []byte passes through unmodified.
		Expect(c.Do(ctx, http.MethodPost, "/raw", []byte(`{"raw":1}`), nil)).To(Succeed())
		Expect(gotBody).To(Equal(`{"raw":1}`))
		Expect(gotCT).To(Equal("application/json"))

		// io.Reader is read fully.
		Expect(c.Do(ctx, http.MethodPost, "/rdr", strings.NewReader(`{"r":2}`), nil)).To(Succeed())
		Expect(gotBody).To(Equal(`{"r":2}`))

		// Unmarshalable values fail before any request is made.
		err := c.Do(ctx, http.MethodPost, "/bad", make(chan int), nil)
		Expect(err).To(MatchError(ContainSubstring("encode request body")))
	})

	It("fails to decode a malformed response body", func(ctx SpecContext) {
		srv := testutil.NewServer(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`not-json`))
		})

		var out map[string]any
		c := newClient(auth.PlaneRuntime, srv.URL)
		err := c.Do(ctx, http.MethodGet, "/x", nil, &out)
		Expect(err).To(MatchError(ContainSubstring("decode response body")))
	})

	It("preserves the endpoint's query string", func(ctx SpecContext) {
		var gotPath, gotQuery string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL)
		Expect(c.Do(ctx, http.MethodGet, "/list?page_size=5&page_token=t1", nil, nil)).To(Succeed())
		Expect(gotPath).To(Equal("/list"))
		Expect(gotQuery).To(Equal("page_size=5&page_token=t1"))
	})
})
