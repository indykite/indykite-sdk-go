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

package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// APIError is a non-2xx response from a platform REST endpoint. It is the
// gRPC-free replacement for the SDK's status-code based errors.
type APIError struct {
	// Code is the platform error code, when present in the body.
	Code string
	// Message is the human-readable error message.
	Message string
	// RequestID correlates the call with server-side logs/traces.
	RequestID string
	// Details are the per-field validation messages from the platform's
	// "errors" array, when present.
	Details []string
	// Body is the raw response body for debugging.
	Body []byte
	// retryAfter is the server-provided Retry-After hint, if any.
	retryAfter time.Duration
	// StatusCode is the HTTP status.
	StatusCode int
}

// RetryAfter returns the server's Retry-After hint, or 0 if none was provided.
func (e *APIError) RetryAfter() time.Duration { return e.retryAfter }

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}
	if len(e.Details) > 0 {
		msg += ": " + strings.Join(e.Details, "; ")
	}
	if e.RequestID != "" {
		return fmt.Sprintf("indykite: %d %s (request_id=%s)", e.StatusCode, msg, e.RequestID)
	}
	return fmt.Sprintf("indykite: %d %s", e.StatusCode, msg)
}

// IsNotFound reports whether the error is a 404.
func (e *APIError) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

// IsUnauthorized reports whether the error is a 401/403.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

// AsAPIError extracts an *APIError from err, if present.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// errorBody is the best-effort shape of a platform error response. The platform
// uses a couple of conventions, so we probe the common fields.
type errorBody struct {
	Code    string   `json:"code"`
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Detail  string   `json:"detail"`
	Errors  []string `json:"errors"`
}

func newAPIError(status int, requestID string, body []byte) *APIError {
	e := &APIError{StatusCode: status, RequestID: requestID, Body: body}
	var parsed errorBody
	if len(body) > 0 && json.Unmarshal(body, &parsed) == nil {
		e.Code = parsed.Code
		e.Message = firstNonEmpty(parsed.Message, parsed.Error, parsed.Detail)
		e.Details = parsed.Errors
	}
	return e
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// isRetryableStatus reports whether an HTTP status warrants a retry.
func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests, // 429
		http.StatusBadGateway,         // 502
		http.StatusServiceUnavailable, // 503
		http.StatusGatewayTimeout:     // 504
		return true
	default:
		return false
	}
}

// isRetryableErr reports whether a transport-level error warrants a retry.
func isRetryableErr(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}
