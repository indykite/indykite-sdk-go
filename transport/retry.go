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
	"math/rand/v2"
	"time"
)

// RetryConfig controls retry behaviour for transient failures (network errors,
// 429 and 5xx). Retries apply to all methods because the platform's runtime
// endpoints (AuthZEN evaluations, CIQ reads) are effectively idempotent and
// capture upserts are keyed/idempotent.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (>=1). 1 disables retries.
	MaxAttempts int
	// BaseDelay is the first backoff delay.
	BaseDelay time.Duration
	// MaxDelay caps the backoff delay.
	MaxDelay time.Duration
}

// DefaultRetryConfig returns sensible defaults: up to 4 attempts with
// exponential backoff and full jitter.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 4,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    5 * time.Second,
	}
}

// backoff returns the delay before the given attempt (attempt is 1-based for the
// first retry). If the last error carried a Retry-After hint it takes priority.
func (r RetryConfig) backoff(attempt int, lastErr error) time.Duration {
	if d := retryAfterFromError(lastErr); d > 0 {
		return capDelay(d, r.MaxDelay)
	}

	// Exponential: base * 2^(attempt-1), capped, then full jitter.
	delay := r.BaseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= r.MaxDelay {
			delay = r.MaxDelay
			break
		}
	}
	delay = capDelay(delay, r.MaxDelay)
	if delay <= 0 {
		return 0
	}
	return time.Duration(rand.Int64N(int64(delay)) + 1) //nolint:gosec // skipcq: GSC-G404 -- jitter
}

func retryAfterFromError(err error) time.Duration {
	apiErr, ok := AsAPIError(err)
	if !ok {
		return 0
	}
	return apiErr.RetryAfter()
}

func capDelay(d, maxDelay time.Duration) time.Duration {
	if maxDelay > 0 && d > maxDelay {
		return maxDelay
	}
	return d
}
