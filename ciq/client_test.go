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

package ciq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/transport"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

func newClient(t *testing.T, h http.HandlerFunc) *ciq.Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return ciq.NewClient(tc)
}

func TestExecuteSinglePage(t *testing.T) {
	var gotPath string
	var gotBody map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		_, _ = io.WriteString(w, `{"data":[{"nodes":{"s":{"id":"gpu-7"}},"aggregate_values":{"count":2}}]}`)
	})

	resp, err := c.Execute(context.Background(), ciq.ExecuteRequest{
		ID:          "get-servers",
		InputParams: map[string]any{"region": "eu"},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if gotPath != "/contx-iq/v1/execute" {
		t.Errorf("path = %q", gotPath)
	}
	if gotBody["id"] != "get-servers" {
		t.Errorf("id = %v", gotBody["id"])
	}
	ip := gotBody["input_params"].(map[string]any)
	if ip["region"] != "eu" {
		t.Errorf("input_params = %v", ip)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("data len = %d", len(resp.Data))
	}
	if resp.Data[0].Nodes["s"].(map[string]any)["id"] != "gpu-7" {
		t.Errorf("nodes = %v", resp.Data[0].Nodes)
	}
	if resp.Data[0].AggregateValues["count"].(float64) != 2 {
		t.Errorf("aggregate_values = %v", resp.Data[0].AggregateValues)
	}
}

// TestPaginationStopsOnShortPage drives the size-based paginator: with PageSize
// 2, a full first page must trigger a second request, and a short second page
// must end iteration.
func TestPaginationStopsOnShortPage(t *testing.T) {
	var requestedTokens []int
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body ciq.ExecuteRequest
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		requestedTokens = append(requestedTokens, body.PageToken)

		switch body.PageToken {
		case 1:
			_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"a"}}},{"nodes":{"n":{"id":"b"}}}]}`)
		case 2:
			_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"c"}}}]}`) // short page -> last
		default:
			_, _ = io.WriteString(w, `{"data":[]}`)
		}
	})

	rows, err := c.All(context.Background(), ciq.ExecuteRequest{ID: "q", PageSize: 2})
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("collected %d rows, want 3", len(rows))
	}
	ids := []string{
		rows[0].Nodes["n"].(map[string]any)["id"].(string),
		rows[1].Nodes["n"].(map[string]any)["id"].(string),
		rows[2].Nodes["n"].(map[string]any)["id"].(string),
	}
	if fmt.Sprint(ids) != "[a b c]" {
		t.Errorf("ids = %v", ids)
	}
	if fmt.Sprint(requestedTokens) != "[1 2]" {
		t.Errorf("requested page tokens = %v, want [1 2]", requestedTokens)
	}
}

// TestPaginationExactMultiple verifies the one extra (empty) fetch when the total
// is an exact multiple of PageSize.
func TestPaginationExactMultiple(t *testing.T) {
	var tokens []int
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body ciq.ExecuteRequest
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		tokens = append(tokens, body.PageToken)
		if body.PageToken == 1 {
			_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"a"}}},{"nodes":{"n":{"id":"b"}}}]}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[]}`) // page 2 empty -> stop
	})

	rows, err := c.All(context.Background(), ciq.ExecuteRequest{ID: "q", PageSize: 2})
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("rows = %d, want 2", len(rows))
	}
	if fmt.Sprint(tokens) != "[1 2]" {
		t.Errorf("tokens = %v, want [1 2]", tokens)
	}
}

// PreprocessParams (CIQ v2.0) are forwarded when set and omitted otherwise.
func TestPreprocessParams(t *testing.T) {
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"data":[]}`)
	})

	if _, err := c.Execute(context.Background(), ciq.ExecuteRequest{
		ID:               "gid:q",
		PreprocessParams: map[string]string{"mode": "strict"},
	}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	pp, ok := body["preprocess_params"].(map[string]any)
	if !ok || pp["mode"] != "strict" {
		t.Errorf("preprocess_params = %v", body["preprocess_params"])
	}

	body = nil
	if _, err := c.Execute(context.Background(), ciq.ExecuteRequest{ID: "gid:q"}); err != nil {
		t.Fatalf("Execute (no params): %v", err)
	}
	if _, present := body["preprocess_params"]; present {
		t.Error("preprocess_params must be omitted when unset")
	}
}
