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
	"sync/atomic"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/test/testutil"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// stubProvider is a TokenProvider that returns a fixed token (no signing key
// needed for transport-level tests).
type stubProvider struct{ tok string }

func (s stubProvider) Token(context.Context) (string, error) { return s.tok, nil }

func newClient(plane auth.Plane, baseURL string, opts ...transport.Option) *transport.Client {
	GinkgoHelper()
	a := auth.NewWithProvider(plane, stubProvider{tok: "test-token"})
	all := append([]transport.Option{transport.WithBaseURL(baseURL)}, opts...)
	c, err := transport.NewClient(a, all...)
	Expect(err).NotTo(HaveOccurred(), "NewClient")
	return c
}

var _ = Describe("Smoke", func() {
	It("sends the runtime-plane client key and round-trips JSON", func(ctx SpecContext) {
		var gotKey, gotAuth, gotReqID, gotBody string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			gotKey = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
			gotAuth = r.Header.Get("Authorization")
			gotReqID = r.Header.Get("X-Request-ID")
			buf := make([]byte, r.ContentLength)
			_, _ = r.Body.Read(buf)
			gotBody = string(buf)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"decision":true}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL)
		var out struct {
			Decision bool `json:"decision"`
		}
		in := map[string]any{"action": "READ"}
		Expect(c.Do(ctx, http.MethodPost, "/access/v1/evaluation", in, &out)).To(Succeed())

		Expect(gotKey).To(Equal("test-token"), "X-IK-ClientKey")
		Expect(gotAuth).To(BeEmpty(), "Authorization should be empty on runtime plane")
		Expect(gotReqID).NotTo(BeEmpty(), "X-Request-ID not set")
		Expect(gotBody).To(Equal(`{"action":"READ"}`))
		Expect(out.Decision).To(BeTrue())
	})

	It("uses a bearer token on the control plane", func(ctx SpecContext) {
		var gotAuth, gotKey string
		srv := testutil.NewServer(func(w http.ResponseWriter, r *http.Request) {
			gotAuth = r.Header.Get("Authorization")
			gotKey = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
			_, _ = w.Write([]byte(`{}`))
		})

		c := newClient(auth.PlaneControl, srv.URL)
		Expect(c.Do(ctx, http.MethodGet, "/configs/v1/projects", nil, nil)).To(Succeed())
		Expect(gotAuth).To(Equal("Bearer test-token"))
		Expect(gotKey).To(BeEmpty(), "X-IK-ClientKey should be empty on control plane")
	})

	It("retries a 503 and then succeeds", func(ctx SpecContext) {
		var calls atomic.Int32
		srv := testutil.NewServer(func(w http.ResponseWriter, _ *http.Request) {
			if calls.Add(1) == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL, transport.WithRetry(transport.RetryConfig{
			MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
		}))
		Expect(c.Do(ctx, http.MethodGet, "/x", nil, nil)).To(Succeed())
		Expect(calls.Load()).To(Equal(int32(2)), "server calls (one retry)")
	})

	It("returns an APIError on 404", func(ctx SpecContext) {
		srv := testutil.NewServer(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Request-ID", "req-123")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code":"NOT_FOUND","message":"no such thing"}`))
		})

		c := newClient(auth.PlaneRuntime, srv.URL)
		err := c.Do(ctx, http.MethodGet, "/missing", nil, nil)
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "expected *APIError, got %T: %v", err, err)
		Expect(apiErr.IsNotFound()).To(BeTrue(), "IsNotFound, status=%d", apiErr.StatusCode)
		Expect(apiErr.Code).To(Equal("NOT_FOUND"))
		Expect(apiErr.Message).To(Equal("no such thing"))
		Expect(apiErr.RequestID).To(Equal("req-123"))
	})

	It("collects every page through the paginator", func(ctx SpecContext) {
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
		got, err := it.Collect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal([]int{1, 2, 3, 4}))
	})
})
