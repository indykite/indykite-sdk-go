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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/onsi/gomega/types"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/test/testutil"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// doErr performs a call against a handler and returns the resulting error.
func doErr(ctx context.Context, handler http.HandlerFunc, opts ...transport.Option) error {
	GinkgoHelper()
	srv := testutil.NewServer(handler)
	c := newClient(auth.PlaneRuntime, srv.URL, opts...)
	return c.Do(ctx, http.MethodGet, "/x", nil, nil)
}

// mustAPIError asserts err is an *APIError and returns it.
func mustAPIError(err error) *transport.APIError {
	GinkgoHelper()
	apiErr, ok := transport.AsAPIError(err)
	Expect(ok).To(BeTrue(), "expected *APIError, got %T: %v", err, err)
	return apiErr
}

var _ = Describe("APIError", func() {
	DescribeTable("treats 401 and 403 as unauthorized",
		func(ctx SpecContext, status int) {
			apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"error":"denied"}`))
			}))
			Expect(apiErr.IsUnauthorized()).To(BeTrue(), "IsUnauthorized")
			Expect(apiErr.IsNotFound()).To(BeFalse(), "IsNotFound")
			Expect(apiErr.Message).To(Equal("denied"), "Message (from error field)")
		},
		Entry("401", http.StatusUnauthorized),
		Entry("403", http.StatusForbidden),
	)

	It("keeps a non-JSON body and falls back to the status text", func(ctx SpecContext) {
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`upstream exploded`))
		}, transport.WithRetry(transport.RetryConfig{MaxAttempts: 1})))

		Expect(apiErr.Message).To(BeEmpty())
		Expect(string(apiErr.Body)).To(Equal("upstream exploded"))
		// Error() falls back to the HTTP status text when the body has no message.
		Expect(apiErr.Error()).To(Equal("indykite: 500 Internal Server Error"))
		Expect(apiErr.RetryAfter()).To(BeZero())
	})

	It("includes the request ID in the error string", func(ctx SpecContext) {
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("X-Request-ID", "req-9")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"bad input"}`))
		}))
		Expect(apiErr.Error()).To(Equal("indykite: 400 bad input (request_id=req-9)"))
	})

	It("is found through wrapping and not on foreign errors", func(ctx SpecContext) {
		err := doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		wrapped := fmt.Errorf("calling config: %w", err)
		apiErr, ok := transport.AsAPIError(wrapped)
		Expect(ok).To(BeTrue(), "AsAPIError(wrapped)")
		Expect(apiErr.IsNotFound()).To(BeTrue())

		_, ok = transport.AsAPIError(errors.New("plain"))
		Expect(ok).To(BeFalse(), "AsAPIError should be false for a non-API error")
		_, ok = transport.AsAPIError(nil)
		Expect(ok).To(BeFalse(), "AsAPIError should be false for nil")
	})

	It("does not populate Code/Message from a non-JSON content type", func(ctx SpecContext) {
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`<html>bad gateway</html>`))
		}, transport.WithRetry(transport.RetryConfig{MaxAttempts: 1})))

		Expect(apiErr.Code).To(BeEmpty(), "HTML body should not populate Code")
		Expect(apiErr.Message).To(BeEmpty(), "HTML body should not populate Message")
		Expect(string(apiErr.Body)).To(ContainSubstring("bad gateway"))
	})

	// The platform's DetailedError carries the useful part in the "errors" array.
	It("reads details from the errors array", func(ctx SpecContext) {
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = io.WriteString(w,
				`{"message":"Unprocessable Entity","errors":["Key: 'Location' Error: gid=PROJECT validation failed"]}`)
		}))
		Expect(apiErr.Details).To(HaveExactElements(ContainSubstring("gid=PROJECT")))
		Expect(apiErr.Error()).To(ContainSubstring("gid=PROJECT validation failed"),
			"Error() should include the validation detail")
	})
})

var _ = Describe("Retry", func() {
	It("honors Retry-After and gives up after max attempts", func(ctx SpecContext) {
		var calls int
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			calls++
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"slow down"}`))
		}, transport.WithRetry(transport.RetryConfig{
			// MaxDelay caps the 2s Retry-After hint so the test stays fast while
			// still exercising the hint-driven backoff path.
			MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
		})))

		Expect(apiErr.StatusCode).To(Equal(http.StatusTooManyRequests))
		Expect(apiErr.RetryAfter()).To(Equal(2 * time.Second))
		Expect(calls).To(Equal(3), "server calls (gave up after max attempts)")
	})

	It("does not retry a plain 4xx", func(ctx SpecContext) {
		var calls int
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			calls++
			w.WriteHeader(http.StatusBadRequest)
		}, transport.WithRetry(transport.RetryConfig{
			MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
		})))

		Expect(apiErr.StatusCode).To(Equal(http.StatusBadRequest))
		Expect(calls).To(Equal(1), "server calls (4xx must not retry)")
	})

	It("exhausts attempts on 5xx with exponential backoff", func(ctx SpecContext) {
		var calls int
		apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
			calls++
			w.WriteHeader(http.StatusBadGateway)
		}, transport.WithRetry(transport.RetryConfig{
			// BaseDelay doubling past MaxDelay exercises the cap branch.
			MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 2 * time.Millisecond,
		})))

		Expect(apiErr.StatusCode).To(Equal(http.StatusBadGateway))
		Expect(calls).To(Equal(4), "server calls")
	})

	It("retries a network error and then fails", func(ctx SpecContext) {
		// A server that is already closed produces a connection error (a net.Error),
		// which the client should retry before giving up.
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{}`))
		}))
		deadURL := srv.URL
		srv.Close()

		c := newClient(auth.PlaneRuntime, deadURL, transport.WithRetry(transport.RetryConfig{
			MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
		}))
		err := c.Do(ctx, http.MethodGet, "/x", nil, nil)
		Expect(err).To(HaveOccurred(), "expected a network error")
		_, ok := transport.AsAPIError(err)
		Expect(ok).To(BeFalse(), "network failure should not be an APIError: %v", err)
	})

	It("aborts when the context is cancelled during backoff", func(ctx SpecContext) {
		srv := testutil.NewServer(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Retry-After", "5")
			w.WriteHeader(http.StatusServiceUnavailable)
		})

		c := newClient(auth.PlaneRuntime, srv.URL, transport.WithRetry(transport.RetryConfig{
			MaxAttempts: 3, BaseDelay: time.Second, MaxDelay: 10 * time.Second,
		}))
		tctx, cancel := context.WithTimeout(ctx, 30*time.Millisecond)
		defer cancel()

		err := c.Do(tctx, http.MethodGet, "/x", nil, nil)
		Expect(err).To(MatchError(context.DeadlineExceeded))
	})

	DescribeTable("parses Retry-After HTTP-date variants",
		func(ctx SpecContext, header func() string, want types.GomegaMatcher) {
			h := header()
			apiErr := mustAPIError(doErr(ctx, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Retry-After", h)
				w.WriteHeader(http.StatusTooManyRequests)
			}, transport.WithRetry(transport.RetryConfig{MaxAttempts: 1})))
			Expect(apiErr.RetryAfter()).To(want, "RetryAfter for header %q", h)
		},
		Entry("future http date",
			func() string { return time.Now().Add(time.Hour).UTC().Format(http.TimeFormat) },
			BeNumerically(">", 50*time.Minute)),
		Entry("past http date",
			func() string { return time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat) },
			BeZero()),
		Entry("garbage", func() string { return "not-a-time" }, BeZero()),
	)

	It("has the documented defaults", func() {
		r := transport.DefaultRetryConfig()
		Expect(r.MaxAttempts).To(Equal(4))
		Expect(r.BaseDelay).To(Equal(200 * time.Millisecond))
		Expect(r.MaxDelay).To(Equal(5 * time.Second))
	})
})
