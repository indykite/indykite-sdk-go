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
	"strings"
	"testing"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/transport"
)

// doErr performs a call against a handler and returns the resulting error.
func doErr(t *testing.T, handler http.HandlerFunc, opts ...transport.Option) error {
	t.Helper()
	srv := httptest.NewServer(handler)
	defer srv.Close()
	c := newClient(t, auth.PlaneRuntime, srv.URL, opts...)
	return c.Do(context.Background(), http.MethodGet, "/x", nil, nil)
}

func TestAPIErrorUnauthorizedAndForbidden(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"denied"}`))
		})
		apiErr, ok := transport.AsAPIError(err)
		if !ok {
			t.Fatalf("status %d: expected *APIError, got %v", status, err)
		}
		if !apiErr.IsUnauthorized() {
			t.Errorf("status %d: IsUnauthorized = false", status)
		}
		if apiErr.IsNotFound() {
			t.Errorf("status %d: IsNotFound = true", status)
		}
		if apiErr.Message != "denied" {
			t.Errorf("status %d: Message = %q, want denied (from error field)", status, apiErr.Message)
		}
	}
}

func TestAPIErrorNonJSONBodyAndErrorString(t *testing.T) {
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`upstream exploded`))
	}, transport.WithRetry(transport.RetryConfig{MaxAttempts: 1}))

	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %v", err)
	}
	if apiErr.Message != "" || string(apiErr.Body) != "upstream exploded" {
		t.Errorf("Message = %q Body = %q", apiErr.Message, apiErr.Body)
	}
	// Error() falls back to the HTTP status text when the body has no message.
	if got := apiErr.Error(); got != "indykite: 500 Internal Server Error" {
		t.Errorf("Error() = %q", got)
	}
	if apiErr.RetryAfter() != 0 {
		t.Errorf("RetryAfter = %v, want 0", apiErr.RetryAfter())
	}
}

func TestAPIErrorStringWithRequestID(t *testing.T) {
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Request-ID", "req-9")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad input"}`))
	})
	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %v", err)
	}
	want := "indykite: 400 bad input (request_id=req-9)"
	if apiErr.Error() != want {
		t.Errorf("Error() = %q, want %q", apiErr.Error(), want)
	}
}

func TestAsAPIErrorOnWrappedAndForeignErrors(t *testing.T) {
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	wrapped := fmt.Errorf("calling config: %w", err)
	apiErr, ok := transport.AsAPIError(wrapped)
	if !ok || !apiErr.IsNotFound() {
		t.Errorf("AsAPIError(wrapped) = %v, %v", apiErr, ok)
	}

	if _, ok = transport.AsAPIError(errors.New("plain")); ok {
		t.Error("AsAPIError should be false for a non-API error")
	}
	if _, ok = transport.AsAPIError(nil); ok {
		t.Error("AsAPIError should be false for nil")
	}
}

func TestRetryHonorsRetryAfterAndGivesUp(t *testing.T) {
	var calls int
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.Header().Set("Retry-After", "2")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"slow down"}`))
	}, transport.WithRetry(transport.RetryConfig{
		// MaxDelay caps the 2s Retry-After hint so the test stays fast while
		// still exercising the hint-driven backoff path.
		MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
	}))

	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %v", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
	if apiErr.RetryAfter() != 2*time.Second {
		t.Errorf("RetryAfter = %v, want 2s", apiErr.RetryAfter())
	}
	if calls != 3 {
		t.Errorf("server called %d times, want 3 (gave up after max attempts)", calls)
	}
}

func TestNoRetryOnPlain4xx(t *testing.T) {
	var calls int
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
	}, transport.WithRetry(transport.RetryConfig{
		MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
	}))

	apiErr, ok := transport.AsAPIError(err)
	if !ok || apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 APIError, got %v", err)
	}
	if calls != 1 {
		t.Errorf("server called %d times, want 1 (4xx must not retry)", calls)
	}
}

func TestRetryExhaustsOn5xxWithExponentialBackoff(t *testing.T) {
	var calls int
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadGateway)
	}, transport.WithRetry(transport.RetryConfig{
		// BaseDelay doubling past MaxDelay exercises the cap branch.
		MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: 2 * time.Millisecond,
	}))

	apiErr, ok := transport.AsAPIError(err)
	if !ok || apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502 APIError, got %v", err)
	}
	if calls != 4 {
		t.Errorf("server called %d times, want 4", calls)
	}
}

func TestRetryOnNetworkErrorThenFail(t *testing.T) {
	// A server that is already closed produces a connection error (a net.Error),
	// which the client should retry before giving up.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	deadURL := srv.URL
	srv.Close()

	c := newClient(t, auth.PlaneRuntime, deadURL, transport.WithRetry(transport.RetryConfig{
		MaxAttempts: 2, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond,
	}))
	err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil)
	if err == nil {
		t.Fatal("expected a network error")
	}
	if _, ok := transport.AsAPIError(err); ok {
		t.Errorf("network failure should not be an APIError: %v", err)
	}
}

func TestRetryAbortsWhenContextCancelledDuringBackoff(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	c := newClient(t, auth.PlaneRuntime, srv.URL, transport.WithRetry(transport.RetryConfig{
		MaxAttempts: 3, BaseDelay: time.Second, MaxDelay: 10 * time.Second,
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := c.Do(ctx, http.MethodGet, "/x", nil, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want context.DeadlineExceeded", err)
	}
}

func TestRetryAfterHTTPDateVariants(t *testing.T) {
	tests := []struct {
		name   string
		header string
		check  func(d time.Duration) bool
	}{
		{"future http date", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat),
			func(d time.Duration) bool { return d > 50*time.Minute }},
		{"past http date", time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat),
			func(d time.Duration) bool { return d == 0 }},
		{"garbage", "not-a-time", func(d time.Duration) bool { return d == 0 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Retry-After", tc.header)
				w.WriteHeader(http.StatusTooManyRequests)
			}, transport.WithRetry(transport.RetryConfig{MaxAttempts: 1}))
			apiErr, ok := transport.AsAPIError(err)
			if !ok {
				t.Fatalf("expected *APIError, got %v", err)
			}
			if !tc.check(apiErr.RetryAfter()) {
				t.Errorf("RetryAfter = %v for header %q", apiErr.RetryAfter(), tc.header)
			}
		})
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	r := transport.DefaultRetryConfig()
	if r.MaxAttempts != 4 || r.BaseDelay != 200*time.Millisecond || r.MaxDelay != 5*time.Second {
		t.Errorf("DefaultRetryConfig = %+v", r)
	}
}

func TestNonJSONErrorContentType(t *testing.T) {
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`<html>bad gateway</html>`))
	}, transport.WithRetry(transport.RetryConfig{MaxAttempts: 1}))

	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %v", err)
	}
	if apiErr.Code != "" || apiErr.Message != "" {
		t.Errorf("HTML body should not populate Code/Message: %+v", apiErr)
	}
	if !strings.Contains(string(apiErr.Body), "bad gateway") {
		t.Errorf("Body = %q", apiErr.Body)
	}
}

// The platform's DetailedError carries the useful part in the "errors" array.
func TestAPIErrorDetailsFromErrorsArray(t *testing.T) {
	err := doErr(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = io.WriteString(w,
			`{"message":"Unprocessable Entity","errors":["Key: 'Location' Error: gid=PROJECT validation failed"]}`)
	})
	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("error is %T, want APIError", err)
	}
	if len(apiErr.Details) != 1 || !strings.Contains(apiErr.Details[0], "gid=PROJECT") {
		t.Errorf("Details = %v", apiErr.Details)
	}
	if !strings.Contains(apiErr.Error(), "gid=PROJECT validation failed") {
		t.Errorf("Error() = %q, want the validation detail included", apiErr.Error())
	}
}
