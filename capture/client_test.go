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

package capture_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/transport"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

func newClient(t *testing.T, h http.HandlerFunc) *capture.Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return capture.NewClient(tc)
}

func TestUpsertNodesBodyAndPath(t *testing.T) {
	var path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:abc"}]}`)
	})

	res, err := c.UpsertNodes(context.Background(), capture.UpsertNode{
		Node:       capture.Node{ExternalID: "ada", Type: "Person"},
		Labels:     []string{"Employee"},
		IsIdentity: true,
		Properties: []capture.Property{{
			BaseProperty: capture.BaseProperty{Type: "email", Value: "ada@x.io"},
		}},
	})
	if err != nil {
		t.Fatalf("UpsertNodes: %v", err)
	}
	if path != "/capture/v1/nodes" {
		t.Errorf("path = %q", path)
	}

	nodes := body["nodes"].([]any)
	n0 := nodes[0].(map[string]any)
	// Embedded Node fields must be promoted to the top level of the node object.
	if n0["external_id"] != "ada" || n0["type"] != "Person" {
		t.Errorf("node identity = %v", n0)
	}
	if isIdentity, _ := n0["is_identity"].(bool); !isIdentity {
		t.Errorf("is_identity = %v", n0["is_identity"])
	}
	prop := n0["properties"].([]any)[0].(map[string]any)
	if prop["type"] != "email" || prop["value"] != "ada@x.io" {
		t.Errorf("property = %v", prop)
	}
	// external_value omitted when Value is set.
	if _, present := prop["external_value"]; present {
		t.Error("external_value should be omitted when value is set")
	}
	if len(res.Results) != 1 || res.Results[0].ID != "gid:abc" {
		t.Errorf("results = %+v", res.Results)
	}
}

func TestUpsertRelationships(t *testing.T) {
	var path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:rel"}]}`)
	})

	_, err := c.UpsertRelationships(context.Background(), capture.Relationship{
		Type:   "WORKS_AT",
		Source: &capture.Node{ExternalID: "ada", Type: "Person"},
		Target: &capture.Node{ExternalID: "acme", Type: "Org"},
	})
	if err != nil {
		t.Fatalf("UpsertRelationships: %v", err)
	}
	if path != "/capture/v1/relationships" {
		t.Errorf("path = %q", path)
	}
	rel := body["relationships"].([]any)[0].(map[string]any)
	src := rel["source"].(map[string]any)
	if rel["type"] != "WORKS_AT" || src["external_id"] != "ada" {
		t.Errorf("relationship = %v", rel)
	}
}

func TestDeleteNodesPath(t *testing.T) {
	var path string
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_, _ = io.WriteString(w, `{"results":[]}`)
	})
	if _, err := c.DeleteNodes(context.Background(),
		capture.Node{ExternalID: "ada", Type: "Person"}); err != nil {
		t.Fatalf("DeleteNodes: %v", err)
	}
	if path != "/capture/v1/nodes/delete" {
		t.Errorf("path = %q", path)
	}
}

func TestReadDataSchema(t *testing.T) {
	var method, path string
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		_, _ = io.WriteString(w, `{"node_types":["Person","Org"]}`)
	})
	schema, err := c.ReadDataSchema(context.Background())
	if err != nil {
		t.Fatalf("ReadDataSchema: %v", err)
	}
	if method != http.MethodGet || path != "/data-schema/v1" {
		t.Errorf("method/path = %s %s", method, path)
	}
	if _, ok := schema["node_types"]; !ok {
		t.Errorf("schema = %v", schema)
	}
}

// TestUpsertNodesChunked verifies the streaming->batch replacement: 5 nodes with
// chunk size 2 must produce 3 requests (2+2+1) and aggregate all 5 results.
func TestUpsertNodesChunked(t *testing.T) {
	var batchSizes []int
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Nodes []json.RawMessage `json:"nodes"`
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		batchSizes = append(batchSizes, len(body.Nodes))
		// Echo one result per node.
		results := make([]string, len(body.Nodes))
		for i := range results {
			results[i] = `{"id":"gid:n"}`
		}
		_, _ = io.WriteString(w, `{"results":[`+strings.Join(results, ",")+`]}`)
	})

	nodes := make([]capture.UpsertNode, 5)
	for i := range nodes {
		nodes[i] = capture.UpsertNode{Node: capture.Node{ExternalID: id(i), Type: "Person"}}
	}

	res, err := c.UpsertNodesChunked(context.Background(), nodes, 2)
	if err != nil {
		t.Fatalf("UpsertNodesChunked: %v", err)
	}
	if got := sum(batchSizes); got != 5 {
		t.Errorf("total nodes sent = %d, want 5", got)
	}
	if len(batchSizes) != 3 || batchSizes[0] != 2 || batchSizes[2] != 1 {
		t.Errorf("batch sizes = %v, want [2 2 1]", batchSizes)
	}
	if len(res.Results) != 5 {
		t.Errorf("aggregated results = %d, want 5", len(res.Results))
	}
}

func TestEmptyBatchRejected(t *testing.T) {
	c := newClient(t, func(http.ResponseWriter, *http.Request) {
		t.Error("server should not be called for empty batch")
	})
	_, err := c.UpsertNodes(context.Background())
	if !errors.Is(err, capture.ErrEmptyBatch) {
		t.Errorf("err = %v, want ErrEmptyBatch", err)
	}
}

func TestDeleteNodeProperties(t *testing.T) {
	var method, path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:abc"}]}`)
	})

	res, err := c.DeleteNodeProperties(context.Background(), capture.DeleteNodeProperties{
		Node:          capture.Node{ExternalID: "ada", Type: "Person"},
		PropertyTypes: []string{"email", "phone"},
	})
	if err != nil {
		t.Fatalf("DeleteNodeProperties: %v", err)
	}
	if method != http.MethodPost || path != "/capture/v1/nodes/properties/delete" {
		t.Errorf("method/path = %s %s", method, path)
	}
	n0 := body["nodes"].([]any)[0].(map[string]any)
	if n0["external_id"] != "ada" || n0["type"] != "Person" {
		t.Errorf("node identity = %v", n0)
	}
	props := n0["property_types"].([]any)
	if len(props) != 2 || props[0] != "email" || props[1] != "phone" {
		t.Errorf("property_types = %v", props)
	}
	if len(res.Results) != 1 || res.Results[0].ID != "gid:abc" {
		t.Errorf("results = %+v", res.Results)
	}
}

func TestDeleteNodePropertyMetadata(t *testing.T) {
	var method, path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:meta"}]}`)
	})

	res, err := c.DeleteNodePropertyMetadata(context.Background(), capture.DeleteNodePropertyMetadata{
		Node:           capture.Node{ExternalID: "ada", Type: "Person"},
		PropertyType:   "email",
		MetadataFields: []string{"source", "assurance_level"},
	})
	if err != nil {
		t.Fatalf("DeleteNodePropertyMetadata: %v", err)
	}
	if method != http.MethodPost || path != "/capture/v1/nodes/properties/metadata/delete" {
		t.Errorf("method/path = %s %s", method, path)
	}
	n0 := body["nodes"].([]any)[0].(map[string]any)
	if n0["external_id"] != "ada" || n0["property_type"] != "email" {
		t.Errorf("node = %v", n0)
	}
	fields := n0["metadata_fields"].([]any)
	if len(fields) != 2 || fields[0] != "source" || fields[1] != "assurance_level" {
		t.Errorf("metadata_fields = %v", fields)
	}
	if len(res.Results) != 1 || res.Results[0].ID != "gid:meta" {
		t.Errorf("results = %+v", res.Results)
	}
}

func TestDeleteRelationships(t *testing.T) {
	var method, path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:rel"}]}`)
	})

	res, err := c.DeleteRelationships(context.Background(), capture.Relationship{
		Type:   "WORKS_AT",
		Source: &capture.Node{ExternalID: "ada", Type: "Person"},
		Target: &capture.Node{ExternalID: "acme", Type: "Org"},
	})
	if err != nil {
		t.Fatalf("DeleteRelationships: %v", err)
	}
	if method != http.MethodPost || path != "/capture/v1/relationships/delete" {
		t.Errorf("method/path = %s %s", method, path)
	}
	rel := body["relationships"].([]any)[0].(map[string]any)
	src := rel["source"].(map[string]any)
	tgt := rel["target"].(map[string]any)
	if rel["type"] != "WORKS_AT" || src["external_id"] != "ada" || tgt["external_id"] != "acme" {
		t.Errorf("relationship = %v", rel)
	}
	if len(res.Results) != 1 || res.Results[0].ID != "gid:rel" {
		t.Errorf("results = %+v", res.Results)
	}
}

func TestDeleteRelationshipProperties(t *testing.T) {
	var method, path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:rel"}]}`)
	})

	_, err := c.DeleteRelationshipProperties(context.Background(), capture.DeleteRelationshipProperties{
		Relationship: capture.Relationship{
			Type:   "OWNS",
			Source: &capture.Node{ExternalID: "ada", Type: "Person"},
			Target: &capture.Node{ExternalID: "kitt", Type: "Car"},
		},
		PropertyTypes: []string{"status"},
	})
	if err != nil {
		t.Fatalf("DeleteRelationshipProperties: %v", err)
	}
	if method != http.MethodPost || path != "/capture/v1/relationships/properties/delete" {
		t.Errorf("method/path = %s %s", method, path)
	}
	rel := body["relationships"].([]any)[0].(map[string]any)
	src := rel["source"].(map[string]any)
	if rel["type"] != "OWNS" || src["external_id"] != "ada" {
		t.Errorf("relationship = %v", rel)
	}
	props := rel["property_types"].([]any)
	if len(props) != 1 || props[0] != "status" {
		t.Errorf("property_types = %v", props)
	}
}

// TestUpsertRelationshipsChunked mirrors the node variant: 5 relationships with
// chunk size 2 must produce 3 requests (2+2+1) and aggregate all 5 results.
func TestUpsertRelationshipsChunked(t *testing.T) {
	var batchSizes []int
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Relationships []json.RawMessage `json:"relationships"`
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		batchSizes = append(batchSizes, len(body.Relationships))
		results := make([]string, len(body.Relationships))
		for i := range results {
			results[i] = `{"id":"gid:r"}`
		}
		_, _ = io.WriteString(w, `{"results":[`+strings.Join(results, ",")+`]}`)
	})

	rels := make([]capture.Relationship, 5)
	for i := range rels {
		rels[i] = capture.Relationship{
			Type:   "KNOWS",
			Source: &capture.Node{ExternalID: id(i), Type: "Person"},
			Target: &capture.Node{ExternalID: "z", Type: "Person"},
		}
	}

	res, err := c.UpsertRelationshipsChunked(context.Background(), rels, 2)
	if err != nil {
		t.Fatalf("UpsertRelationshipsChunked: %v", err)
	}
	if got := sum(batchSizes); got != 5 {
		t.Errorf("total relationships sent = %d, want 5", got)
	}
	if len(batchSizes) != 3 || batchSizes[0] != 2 || batchSizes[1] != 2 || batchSizes[2] != 1 {
		t.Errorf("batch sizes = %v, want [2 2 1]", batchSizes)
	}
	if len(res.Results) != 5 {
		t.Errorf("aggregated results = %d, want 5", len(res.Results))
	}
}

// A chunk size <= 0 or > MaxBatchSize is normalized to MaxBatchSize, so a small
// slice goes out as a single request.
func TestChunkedNormalizesChunkSize(t *testing.T) {
	var requests int
	c := newClient(t, func(w http.ResponseWriter, _ *http.Request) {
		requests++
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:n"}]}`)
	})

	nodes := []capture.UpsertNode{{Node: capture.Node{ExternalID: "a", Type: "Person"}}}
	if _, err := c.UpsertNodesChunked(context.Background(), nodes, 0); err != nil {
		t.Fatalf("UpsertNodesChunked(size 0): %v", err)
	}
	if _, err := c.UpsertNodesChunked(context.Background(), nodes, capture.MaxBatchSize+1); err != nil {
		t.Fatalf("UpsertNodesChunked(size max+1): %v", err)
	}
	if requests != 2 {
		t.Errorf("requests = %d, want 2 (one per call)", requests)
	}
}

// A failing batch stops the chunked upsert and surfaces the API error.
func TestChunkedStopsOnError(t *testing.T) {
	var requests int
	c := newClient(t, func(w http.ResponseWriter, _ *http.Request) {
		requests++
		if requests == 2 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `{"message":"denied"}`)
			return
		}
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:n"},{"id":"gid:n"}]}`)
	})

	nodes := make([]capture.UpsertNode, 5)
	for i := range nodes {
		nodes[i] = capture.UpsertNode{Node: capture.Node{ExternalID: id(i), Type: "Person"}}
	}
	res, err := c.UpsertNodesChunked(context.Background(), nodes, 2)
	if err == nil {
		t.Fatal("expected error from failing batch")
	}
	apiErr, ok := transport.AsAPIError(err)
	if !ok || apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("err = %v, want 403 APIError", err)
	}
	if requests != 2 {
		t.Errorf("requests = %d, want 2 (stop after failure)", requests)
	}
	if len(res.Results) != 2 {
		t.Errorf("partial results = %d, want 2", len(res.Results))
	}
}

// Every batch write must reject an empty batch without touching the server.
func TestEmptyBatchRejectedEverywhere(t *testing.T) {
	c := newClient(t, func(http.ResponseWriter, *http.Request) {
		t.Error("server should not be called for empty batch")
	})
	ctx := context.Background()
	calls := map[string]func() error{
		"DeleteNodes": func() error {
			_, err := c.DeleteNodes(ctx)
			return err
		},
		"DeleteNodeProperties": func() error {
			_, err := c.DeleteNodeProperties(ctx)
			return err
		},
		"DeleteNodePropertyMetadata": func() error {
			_, err := c.DeleteNodePropertyMetadata(ctx)
			return err
		},
		"UpsertRelationships": func() error {
			_, err := c.UpsertRelationships(ctx)
			return err
		},
		"DeleteRelationships": func() error {
			_, err := c.DeleteRelationships(ctx)
			return err
		},
		"DeleteRelationshipProperties": func() error {
			_, err := c.DeleteRelationshipProperties(ctx)
			return err
		},
		"UpsertNodesChunked": func() error {
			_, err := c.UpsertNodesChunked(ctx, nil, 2)
			return err
		},
		"UpsertRelationshipsChunked": func() error {
			_, err := c.UpsertRelationshipsChunked(ctx, nil, 2)
			return err
		},
	}
	for name, call := range calls {
		if err := call(); !errors.Is(err, capture.ErrEmptyBatch) {
			t.Errorf("%s: err = %v, want ErrEmptyBatch", name, err)
		}
	}
}

// A non-2xx response surfaces as *transport.APIError on every read/write path.
func TestReadDataSchemaError(t *testing.T) {
	c := newClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"message":"denied"}`)
	})
	if _, err := c.ReadDataSchema(context.Background()); err == nil {
		t.Fatal("expected error")
	} else if apiErr, ok := transport.AsAPIError(err); !ok || apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("err = %v, want 403 APIError", err)
	}
}

// Small test helpers.
func sum(xs []int) int {
	t := 0
	for _, x := range xs {
		t += x
	}
	return t
}

func id(i int) string { return string(rune('a' + i)) }

// GlobalDB routes relationship operations with use_global_db; node identities
// carry an optional location for composite IKGs.
func TestGlobalDBAndNodeLocation(t *testing.T) {
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:r"}]}`)
	})

	rel := capture.Relationship{
		Type:   "OWNS",
		Source: &capture.Node{ExternalID: "p1", Type: "Person", Location: "eu-db"},
		Target: &capture.Node{ExternalID: "c1", Type: "Car"},
	}
	if _, err := c.GlobalDB().UpsertRelationships(context.Background(), rel); err != nil {
		t.Fatalf("UpsertRelationships: %v", err)
	}
	if v, _ := body["use_global_db"].(bool); !v {
		t.Errorf("use_global_db = %v, want true", body["use_global_db"])
	}
	src := body["relationships"].([]any)[0].(map[string]any)["source"].(map[string]any)
	if src["location"] != "eu-db" {
		t.Errorf("source.location = %v, want eu-db", src["location"])
	}

	// The plain client must not send the flag at all.
	body = nil
	if _, err := c.UpsertRelationships(context.Background(), rel); err != nil {
		t.Fatalf("UpsertRelationships (plain): %v", err)
	}
	if _, present := body["use_global_db"]; present {
		t.Error("use_global_db must be omitted on the non-global client")
	}

	if _, err := c.GlobalDB().DeleteRelationships(context.Background(), rel); err != nil {
		t.Fatalf("DeleteRelationships: %v", err)
	}
	if v, _ := body["use_global_db"].(bool); !v {
		t.Errorf("delete use_global_db = %v, want true", body["use_global_db"])
	}
}

// UpsertSingleNode addresses one node via the path id; the body carries only
// labels/properties/is_identity.
func TestUpsertSingleNode(t *testing.T) {
	var method, path string
	var body map[string]any
	c := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"results":[{"id":"gid:n1"}]}`)
	})

	res, err := c.UpsertSingleNode(context.Background(), "Person:ada", capture.PutNode{
		IsIdentity: true,
		Properties: []capture.Property{{
			BaseProperty: capture.BaseProperty{Type: "email", Value: "ada@x.io"},
		}},
	})
	if err != nil {
		t.Fatalf("UpsertSingleNode: %v", err)
	}
	if method != http.MethodPut || path != "/capture/v1/nodes/Person:ada" {
		t.Errorf("method/path = %s %s", method, path)
	}
	if _, present := body["external_id"]; present {
		t.Error("body must not carry external_id; identity comes from the path")
	}
	if isIdentity, _ := body["is_identity"].(bool); !isIdentity {
		t.Errorf("is_identity = %v", body["is_identity"])
	}
	if len(res.Results) != 1 || res.Results[0].ID != "gid:n1" {
		t.Errorf("results = %+v", res.Results)
	}
}
