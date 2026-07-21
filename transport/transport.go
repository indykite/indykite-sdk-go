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

// Package transport is the shared HTTP core of the REST SDK: base-URL/region
// resolution, authentication, JSON encode/decode, retries with backoff,
// request-id propagation, optional OpenTelemetry tracing and typed API errors.
//
// Domain packages (authzen, ciq, capture, entitymatching, config) are thin
// facades built on top of a *transport.Client.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/indykite/indykite-sdk-go/auth"
)

const (
	defaultUserAgent  = "indykite-sdk-go"
	defaultTimeout    = 30 * time.Second
	regionURLTemplate = "https://%s.api.indykite.com"

	headerRequestID   = "X-Request-ID"
	headerContentType = "Content-Type"
	headerAccept      = "Accept"
	headerUserAgent   = "User-Agent"
	contentTypeJSON   = "application/json"
)

// Client is the shared HTTP client used by every domain package.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	auth       *auth.Authenticator
	userAgent  string
	retry      RetryConfig
}

// Option configures a Client.
type Option func(*clientOptions) error

type clientOptions struct {
	httpClient *http.Client
	retry      *RetryConfig
	baseURL    string
	region     string
	userAgent  string
	tracing    bool
}

// WithBaseURL sets an explicit base URL (highest precedence).
func WithBaseURL(raw string) Option {
	return func(o *clientOptions) error { o.baseURL = raw; return nil }
}

// WithRegion resolves the base URL as https://<region>.api.indykite.com.
func WithRegion(region string) Option {
	return func(o *clientOptions) error { o.region = region; return nil }
}

// WithHTTPClient supplies a custom *http.Client.
func WithHTTPClient(c *http.Client) Option {
	return func(o *clientOptions) error { o.httpClient = c; return nil }
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(o *clientOptions) error { o.userAgent = ua; return nil }
}

// WithRetry overrides the retry policy.
func WithRetry(r RetryConfig) Option {
	return func(o *clientOptions) error { o.retry = &r; return nil }
}

// WithTracing wraps the HTTP transport with OpenTelemetry instrumentation.
func WithTracing() Option {
	return func(o *clientOptions) error { o.tracing = true; return nil }
}

// NewClient builds a Client. Base URL precedence:
// WithBaseURL > WithRegion > authenticator credential BaseURL.
func NewClient(authenticator *auth.Authenticator, opts ...Option) (*Client, error) {
	if authenticator == nil {
		return nil, errors.New("transport: authenticator is required")
	}

	o := &clientOptions{
		userAgent: defaultUserAgent,
		retry:     nil,
	}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	rawBase, err := resolveBaseURL(o, authenticator)
	if err != nil {
		return nil, err
	}
	base, err := url.Parse(rawBase)
	if err != nil {
		return nil, fmt.Errorf("transport: invalid base URL %q: %w", rawBase, err)
	}

	httpClient := o.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	if o.tracing {
		httpClient = instrumentHTTPClient(httpClient)
	}

	retry := DefaultRetryConfig()
	if o.retry != nil {
		retry = *o.retry
	}

	return &Client{
		baseURL:    base,
		httpClient: httpClient,
		auth:       authenticator,
		userAgent:  o.userAgent,
		retry:      retry,
	}, nil
}

func resolveBaseURL(o *clientOptions, a *auth.Authenticator) (string, error) {
	switch {
	case o.baseURL != "":
		return o.baseURL, nil
	case o.region != "":
		return fmt.Sprintf(regionURLTemplate, o.region), nil
	case a.BaseURL != "":
		return a.BaseURL, nil
	default:
		return "", errors.New("transport: no base URL — set WithBaseURL/WithRegion or provide a credential baseUrl")
	}
}

// CallOption mutates the outgoing request before it is sent — used for per-call
// headers such as If-Match (optimistic concurrency). Options run after the
// standard/auth headers, so they can override them.
type CallOption func(*http.Request)

// WithHeader sets a request header on a single call.
func WithHeader(key, value string) CallOption {
	return func(r *http.Request) { r.Header.Set(key, value) }
}

// Response carries response metadata for callers that need more than the decoded
// body — e.g. the ETag header for control-plane optimistic concurrency.
type Response struct {
	Header     http.Header
	RequestID  string
	StatusCode int
}

// Do performs a single API call: it encodes body as JSON (unless body is nil,
// []byte or io.Reader), authenticates, retries per the policy, and decodes a 2xx
// response into out (out may be nil to discard the body). Non-2xx responses are
// returned as *APIError.
func (c *Client) Do(ctx context.Context, method, path string, body, out any, opts ...CallOption) error {
	_, err := c.DoResp(ctx, method, path, body, out, opts...)
	return err
}

// DoResp is like Do but also returns response metadata (status, headers, request
// id). Use it when you need a response header such as ETag.
func (c *Client) DoResp(
	ctx context.Context,
	method, path string,
	body, out any,
	opts ...CallOption,
) (*Response, error) {
	bodyBytes, err := encodeBody(body)
	if err != nil {
		return nil, err
	}

	endpoint := c.endpoint(path)

	var lastErr error
	for attempt := range c.retry.MaxAttempts {
		if attempt > 0 {
			if waitErr := sleep(ctx, c.retry.backoff(attempt, lastErr)); waitErr != nil {
				return nil, waitErr
			}
		}

		resp, doErr := c.attempt(ctx, method, endpoint, bodyBytes, opts)
		if doErr != nil {
			lastErr = doErr
			if isRetryableErr(doErr) {
				continue
			}
			return nil, doErr
		}

		// Retryable status: remember it and loop (body already drained in attempt()).
		if isRetryableStatus(resp.statusCode) && attempt < c.retry.MaxAttempts-1 {
			lastErr = resp.toError()
			continue
		}

		if resp.statusCode < 200 || resp.statusCode >= 300 {
			return nil, resp.toError()
		}
		if decErr := decodeInto(resp.body, out); decErr != nil {
			return nil, decErr
		}
		return &Response{StatusCode: resp.statusCode, Header: resp.header, RequestID: resp.requestID}, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("transport: request failed after %d attempts", c.retry.MaxAttempts)
}

// endpoint joins the base URL with path, preserving any query string in path.
func (c *Client) endpoint(path string) string {
	p, rawQuery, _ := strings.Cut(path, "?")
	u := c.baseURL.JoinPath(p)
	u.RawQuery = rawQuery
	return u.String()
}

type rawResponse struct {
	header     http.Header
	requestID  string
	body       []byte
	statusCode int
	retryAfter time.Duration
}

func (r *rawResponse) toError() *APIError {
	e := newAPIError(r.statusCode, r.requestID, r.body)
	e.retryAfter = r.retryAfter
	return e
}

func (c *Client) attempt(
	ctx context.Context,
	method, endpoint string,
	body []byte,
	opts []CallOption,
) (*rawResponse, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("transport: build request: %w", err)
	}

	if body != nil {
		req.Header.Set(headerContentType, contentTypeJSON)
	}
	req.Header.Set(headerAccept, contentTypeJSON)
	req.Header.Set(headerUserAgent, c.userAgent)
	if req.Header.Get(headerRequestID) == "" {
		req.Header.Set(headerRequestID, uuid.NewString())
	}
	if err = c.auth.Apply(ctx, req); err != nil {
		return nil, fmt.Errorf("transport: authenticate: %w", err)
	}
	// Per-call options last so they can override standard headers.
	for _, o := range opts {
		o(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("transport: read response: %w", err)
	}

	return &rawResponse{
		statusCode: resp.StatusCode,
		requestID:  responseRequestID(resp.Header),
		retryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		header:     resp.Header,
		body:       data,
	}, nil
}

func encodeBody(body any) ([]byte, error) {
	switch b := body.(type) {
	case nil:
		return nil, nil
	case []byte:
		return b, nil
	case io.Reader:
		return io.ReadAll(b)
	default:
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("transport: encode request body: %w", err)
		}
		return data, nil
	}
}

func decodeInto(body []byte, out any) error {
	if out == nil || len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("transport: decode response body: %w", err)
	}
	return nil
}

func responseRequestID(h http.Header) string {
	for _, k := range []string{headerRequestID, "X-Request-Id", "X-Indykite-Request-Id"} {
		if v := h.Get(k); v != "" {
			return v
		}
	}
	return ""
}

func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := time.ParseDuration(strings.TrimSpace(v) + "s"); err == nil {
		return secs
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
