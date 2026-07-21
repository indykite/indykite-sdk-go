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

package entitymatching_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/entitymatching"
	"github.com/indykite/indykite-sdk-go/transport"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

func newClient(t *testing.T, h http.HandlerFunc) *entitymatching.Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return entitymatching.NewClient(tc)
}

func TestRunAccepted(t *testing.T) {
	var path, key string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		key = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusAccepted) // 202
		_, _ = io.WriteString(w, `{"id":"gid:em1","last_run_time":"t0","etag":"\"v1\""}`)
	})

	run, err := c.Run(context.Background(), "gid:em1", entitymatching.RunRequest{
		SimilarityScoreCutoff: 0.8,
		CustomPropertyMappings: []entitymatching.CustomPropertyMapping{
			{SourceNodeProperty: "email", TargetNodeProperty: "email"},
		},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if path != "/entity-matching/v1/pipelines/gid:em1/runs" {
		t.Errorf("path = %q", path)
	}
	if key != "tok" {
		t.Errorf("expected runtime-plane X-IK-ClientKey, got %q", key)
	}
	if body["similarity_score_cutoff"].(float64) != 0.8 {
		t.Errorf("cutoff = %v", body["similarity_score_cutoff"])
	}
	cpm := body["custom_property_mappings"].([]any)[0].(map[string]any)
	if cpm["source_node_property"] != "email" {
		t.Errorf("custom_property_mappings = %v", cpm)
	}
	if run.ID != "gid:em1" || run.Etag != `"v1"` {
		t.Errorf("run = %+v", run)
	}
}

func TestSuggestedMappingsReady(t *testing.T) {
	var path string
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_, _ = io.WriteString(w, `{"id":"gid:em1","suggested_property_mappings":[
			{"source_node_type":"Person","source_node_property":"email",
			 "target_node_type":"Lead","target_node_property":"email","similarity_score_cutoff":0.9}]}`)
	})

	m, err := c.SuggestedPropertyMappings(context.Background(), "gid:em1")
	if err != nil {
		t.Fatalf("SuggestedPropertyMappings: %v", err)
	}
	if path != "/entity-matching/v1/pipelines/gid:em1/property-mappings" {
		t.Errorf("path = %q", path)
	}
	if len(m.SuggestedPropertyMappings) != 1 {
		t.Fatalf("mappings = %d", len(m.SuggestedPropertyMappings))
	}
	sm := m.SuggestedPropertyMappings[0]
	if sm.SourceNodeType != "Person" || sm.SimilarityScoreCutoff != 0.9 {
		t.Errorf("mapping = %+v", sm)
	}
}

func TestSuggestedMappingsNotReady(t *testing.T) {
	c := newClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted) // 202 -> still computing
		_, _ = io.WriteString(w, `{"message":"mapping in progress"}`)
	})

	_, err := c.SuggestedPropertyMappings(context.Background(), "gid:em1")
	if !errors.Is(err, entitymatching.ErrMappingNotReady) {
		t.Errorf("err = %v, want ErrMappingNotReady", err)
	}
}

func TestStatus(t *testing.T) {
	var method, path string
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_, _ = io.WriteString(w, `{"id":"gid:em1","property_mapping_status":"DONE","entity_matching_status":"RUNNING"}`)
	})

	st, err := c.Status(context.Background(), "gid:em1")
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if method != http.MethodGet || path != "/entity-matching/v1/pipelines/gid:em1/status" {
		t.Errorf("method/path = %s %s", method, path)
	}
	if st.PropertyMappingStatus != "DONE" || st.EntityMatchingStatus != "RUNNING" {
		t.Errorf("status = %+v", st)
	}
}
