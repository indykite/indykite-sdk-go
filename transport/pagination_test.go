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
	"testing"

	"github.com/indykite/indykite-sdk-go/transport"
)

func TestIteratorManualNextItem(t *testing.T) {
	fetch := func(_ context.Context, token string) (transport.Page[string], error) {
		if token == "" {
			return transport.Page[string]{Items: []string{"a", "b"}, NextToken: ""}, nil
		}
		return transport.Page[string]{}, fmt.Errorf("unexpected token %q", token)
	}

	it := transport.NewIterator(fetch)
	ctx := context.Background()

	var got []string
	for it.Next(ctx) {
		got = append(got, it.Item())
	}
	if err := it.Err(); err != nil {
		t.Fatalf("Err: %v", err)
	}
	if fmt.Sprint(got) != "[a b]" {
		t.Errorf("items = %v, want [a b]", got)
	}
	// Exhausted iterators keep returning false without refetching.
	if it.Next(ctx) {
		t.Error("Next after exhaustion = true, want false")
	}
}

func TestIteratorErrorPropagation(t *testing.T) {
	boom := errors.New("boom")
	fetch := func(_ context.Context, token string) (transport.Page[int], error) {
		if token == "" {
			return transport.Page[int]{Items: []int{1}, NextToken: "p2"}, nil
		}
		return transport.Page[int]{}, boom
	}

	it := transport.NewIterator(fetch)
	got, err := it.Collect(context.Background())
	if !errors.Is(err, boom) {
		t.Fatalf("Collect err = %v, want boom", err)
	}
	if len(got) != 1 || got[0] != 1 {
		t.Errorf("partial items = %v, want [1]", got)
	}
	// After an error, Next stays false and Err stays set.
	if it.Next(context.Background()) {
		t.Error("Next after error = true, want false")
	}
	if !errors.Is(it.Err(), boom) {
		t.Errorf("Err = %v, want boom", it.Err())
	}
}

func TestIteratorEmptyFirstPage(t *testing.T) {
	fetch := func(_ context.Context, _ string) (transport.Page[int], error) {
		return transport.Page[int]{}, nil
	}
	got, err := transport.NewIterator(fetch).Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("items = %v, want empty", got)
	}
}

func TestIteratorSkipsEmptyMiddlePage(t *testing.T) {
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
	got, err := transport.NewIterator(fetch).Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	if fmt.Sprint(got) != "[7 8]" {
		t.Errorf("items = %v, want [7 8]", got)
	}
}
