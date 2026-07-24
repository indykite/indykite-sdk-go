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

import "context"

// Page is one page of results plus the token for the next page. Paginated
// domain endpoints (e.g. ciq.Execute) return this from their fetch func.
type Page[T any] struct {
	NextToken string
	Items     []T
}

// FetchPage retrieves a single page given a page token ("" for the first page).
type FetchPage[T any] func(ctx context.Context, pageToken string) (Page[T], error)

// Iterator lazily walks page_token / page_size paginated endpoints.
//
//	it := transport.NewIterator(fetch)
//	for it.Next(ctx) {
//	    item := it.Item()
//	}
//	if err := it.Err(); err != nil { ... }
type Iterator[T any] struct {
	fetch     FetchPage[T]
	err       error
	nextToken string
	current   T
	buf       []T
	started   bool
	done      bool
}

// NewIterator builds an Iterator over the given fetch function.
func NewIterator[T any](fetch FetchPage[T]) *Iterator[T] {
	return &Iterator[T]{fetch: fetch}
}

// Next advances to the next item, fetching a new page when the buffer drains.
// It returns false when the stream is exhausted or an error occurred (check Err).
func (it *Iterator[T]) Next(ctx context.Context) bool {
	if it.err != nil || it.done {
		return false
	}

	for len(it.buf) == 0 {
		if it.started && it.nextToken == "" {
			it.done = true
			return false
		}
		page, err := it.fetch(ctx, it.nextToken)
		if err != nil {
			it.err = err
			return false
		}
		it.started = true
		it.buf = page.Items
		it.nextToken = page.NextToken
		if len(it.buf) == 0 && it.nextToken == "" {
			it.done = true
			return false
		}
	}

	it.current = it.buf[0]
	it.buf = it.buf[1:]
	return true
}

// Item returns the current item (valid after Next returns true).
func (it *Iterator[T]) Item() T { return it.current }

// Err returns the first error encountered while iterating, if any.
func (it *Iterator[T]) Err() error { return it.err }

// Collect drains the iterator into a slice.
func (it *Iterator[T]) Collect(ctx context.Context) ([]T, error) {
	var out []T
	for it.Next(ctx) {
		out = append(out, it.Item())
	}
	return out, it.Err()
}
