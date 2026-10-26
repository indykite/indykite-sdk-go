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

	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iterator", func() {
	It("yields items through manual Next/Item", func(ctx SpecContext) {
		fetch := func(_ context.Context, token string) (transport.Page[string], error) {
			if token == "" {
				return transport.Page[string]{Items: []string{"a", "b"}, NextToken: ""}, nil
			}
			return transport.Page[string]{}, fmt.Errorf("unexpected token %q", token)
		}

		it := transport.NewIterator(fetch)

		var got []string
		for it.Next(ctx) {
			got = append(got, it.Item())
		}
		Expect(it.Err()).NotTo(HaveOccurred())
		Expect(got).To(Equal([]string{"a", "b"}))
		// Exhausted iterators keep returning false without refetching.
		Expect(it.Next(ctx)).To(BeFalse(), "Next after exhaustion")
	})

	It("propagates fetch errors", func(ctx SpecContext) {
		boom := errors.New("boom")
		fetch := func(_ context.Context, token string) (transport.Page[int], error) {
			if token == "" {
				return transport.Page[int]{Items: []int{1}, NextToken: "p2"}, nil
			}
			return transport.Page[int]{}, boom
		}

		it := transport.NewIterator(fetch)
		got, err := it.Collect(ctx)
		Expect(err).To(MatchError(boom))
		Expect(got).To(Equal([]int{1}), "partial items")
		// After an error, Next stays false and Err stays set.
		Expect(it.Next(ctx)).To(BeFalse(), "Next after error")
		Expect(it.Err()).To(MatchError(boom))
	})

	It("handles an empty first page", func(ctx SpecContext) {
		fetch := func(_ context.Context, _ string) (transport.Page[int], error) {
			return transport.Page[int]{}, nil
		}
		got, err := transport.NewIterator(fetch).Collect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(BeEmpty())
	})

	It("skips an empty middle page", func(ctx SpecContext) {
		fetch := func(_ context.Context, token string) (transport.Page[int], error) {
			switch token {
			case "":
				return transport.Page[int]{Items: nil, NextToken: "p2"}, nil
			case "p2":
				return transport.Page[int]{Items: []int{7, 8}, NextToken: ""}, nil
			default:
				return transport.Page[int]{}, fmt.Errorf("unexpected token %q", token)
			}
		}
		got, err := transport.NewIterator(fetch).Collect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal([]int{7, 8}))
	})
})
