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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

// newClient wires a capture.Client to an httptest server running h; the server
// is torn down when the spec ends.
func newClient(h http.HandlerFunc) *capture.Client {
	GinkgoHelper()
	srv := httptest.NewServer(h)
	DeferCleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	Expect(err).NotTo(HaveOccurred())
	return capture.NewClient(tc)
}

// recorded is what the fake server saw on the last call.
type recorded struct {
	body   map[string]any
	method string
	path   string
}

// recordingClient returns a client whose server records the method, path and
// decoded JSON body of the last request and answers with reply.
func recordingClient(reply string) (*capture.Client, *recorded) {
	GinkgoHelper()
	rec := &recorded{}
	c := newClient(func(w http.ResponseWriter, r *http.Request) {
		rec.method, rec.path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		rec.body = nil
		_ = json.Unmarshal(raw, &rec.body)
		_, _ = io.WriteString(w, reply)
	})
	return c, rec
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

var _ = Describe("Nodes", func() {
	It("upserts nodes with promoted identity fields and typed properties", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:abc"}]}`)

		res, err := c.UpsertNodes(ctx, capture.UpsertNode{
			Node:       capture.Node{ExternalID: "ada", Type: "Person"},
			Labels:     []string{"Employee"},
			IsIdentity: true,
			Properties: []capture.Property{{
				BaseProperty: capture.BaseProperty{Type: "email", Value: "ada@x.io"},
			}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/capture/v1/nodes"))

		nodes := rec.body["nodes"].([]any)
		n0 := nodes[0].(map[string]any)
		// Embedded Node fields must be promoted to the top level of the node object.
		Expect(n0).To(HaveKeyWithValue("external_id", "ada"))
		Expect(n0).To(HaveKeyWithValue("type", "Person"))
		Expect(n0).To(HaveKeyWithValue("is_identity", true))
		prop := n0["properties"].([]any)[0].(map[string]any)
		Expect(prop).To(HaveKeyWithValue("type", "email"))
		Expect(prop).To(HaveKeyWithValue("value", "ada@x.io"))
		// external_value omitted when Value is set.
		Expect(prop).NotTo(HaveKey("external_value"), "external_value should be omitted when value is set")
		Expect(res.Results).To(HaveLen(1))
		Expect(res.Results[0].ID).To(Equal("gid:abc"))
	})

	It("deletes nodes via the delete path", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[]}`)

		_, err := c.DeleteNodes(ctx, capture.Node{ExternalID: "ada", Type: "Person"})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/capture/v1/nodes/delete"))
	})

	It("deletes node properties", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:abc"}]}`)

		res, err := c.DeleteNodeProperties(ctx, capture.DeleteNodeProperties{
			Node:          capture.Node{ExternalID: "ada", Type: "Person"},
			PropertyTypes: []string{"email", "phone"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.method).To(Equal(http.MethodPost))
		Expect(rec.path).To(Equal("/capture/v1/nodes/properties/delete"))
		n0 := rec.body["nodes"].([]any)[0].(map[string]any)
		Expect(n0).To(HaveKeyWithValue("external_id", "ada"))
		Expect(n0).To(HaveKeyWithValue("type", "Person"))
		Expect(n0["property_types"]).To(Equal([]any{"email", "phone"}))
		Expect(res.Results).To(HaveLen(1))
		Expect(res.Results[0].ID).To(Equal("gid:abc"))
	})

	It("deletes node property metadata", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:meta"}]}`)

		res, err := c.DeleteNodePropertyMetadata(ctx, capture.DeleteNodePropertyMetadata{
			Node:           capture.Node{ExternalID: "ada", Type: "Person"},
			PropertyType:   "email",
			MetadataFields: []string{"source", "assurance_level"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.method).To(Equal(http.MethodPost))
		Expect(rec.path).To(Equal("/capture/v1/nodes/properties/metadata/delete"))
		n0 := rec.body["nodes"].([]any)[0].(map[string]any)
		Expect(n0).To(HaveKeyWithValue("external_id", "ada"))
		Expect(n0).To(HaveKeyWithValue("property_type", "email"))
		Expect(n0["metadata_fields"]).To(Equal([]any{"source", "assurance_level"}))
		Expect(res.Results).To(HaveLen(1))
		Expect(res.Results[0].ID).To(Equal("gid:meta"))
	})

	// UpsertSingleNode addresses one node via the path id; the body carries only
	// labels/properties/is_identity.
	It("upserts a single node addressed by the path id", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:n1"}]}`)

		res, err := c.UpsertSingleNode(ctx, "Person:ada", capture.PutNode{
			IsIdentity: true,
			Properties: []capture.Property{{
				BaseProperty: capture.BaseProperty{Type: "email", Value: "ada@x.io"},
			}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.method).To(Equal(http.MethodPut))
		Expect(rec.path).To(Equal("/capture/v1/nodes/Person:ada"))
		Expect(rec.body).NotTo(HaveKey("external_id"), "body must not carry external_id; identity comes from the path")
		Expect(rec.body).To(HaveKeyWithValue("is_identity", true))
		Expect(res.Results).To(HaveLen(1))
		Expect(res.Results[0].ID).To(Equal("gid:n1"))
	})
})

var _ = Describe("Relationships", func() {
	It("upserts relationships", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:rel"}]}`)

		_, err := c.UpsertRelationships(ctx, capture.Relationship{
			Type:   "WORKS_AT",
			Source: &capture.Node{ExternalID: "ada", Type: "Person"},
			Target: &capture.Node{ExternalID: "acme", Type: "Org"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/capture/v1/relationships"))
		rel := rec.body["relationships"].([]any)[0].(map[string]any)
		Expect(rel).To(HaveKeyWithValue("type", "WORKS_AT"))
		Expect(rel["source"]).To(HaveKeyWithValue("external_id", "ada"))
	})

	It("deletes relationships", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:rel"}]}`)

		res, err := c.DeleteRelationships(ctx, capture.Relationship{
			Type:   "WORKS_AT",
			Source: &capture.Node{ExternalID: "ada", Type: "Person"},
			Target: &capture.Node{ExternalID: "acme", Type: "Org"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.method).To(Equal(http.MethodPost))
		Expect(rec.path).To(Equal("/capture/v1/relationships/delete"))
		rel := rec.body["relationships"].([]any)[0].(map[string]any)
		Expect(rel).To(HaveKeyWithValue("type", "WORKS_AT"))
		Expect(rel["source"]).To(HaveKeyWithValue("external_id", "ada"))
		Expect(rel["target"]).To(HaveKeyWithValue("external_id", "acme"))
		Expect(res.Results).To(HaveLen(1))
		Expect(res.Results[0].ID).To(Equal("gid:rel"))
	})

	It("deletes relationship properties", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:rel"}]}`)

		_, err := c.DeleteRelationshipProperties(ctx, capture.DeleteRelationshipProperties{
			Relationship: capture.Relationship{
				Type:   "OWNS",
				Source: &capture.Node{ExternalID: "ada", Type: "Person"},
				Target: &capture.Node{ExternalID: "kitt", Type: "Car"},
			},
			PropertyTypes: []string{"status"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.method).To(Equal(http.MethodPost))
		Expect(rec.path).To(Equal("/capture/v1/relationships/properties/delete"))
		rel := rec.body["relationships"].([]any)[0].(map[string]any)
		Expect(rel).To(HaveKeyWithValue("type", "OWNS"))
		Expect(rel["source"]).To(HaveKeyWithValue("external_id", "ada"))
		Expect(rel["property_types"]).To(Equal([]any{"status"}))
	})

	// GlobalDB routes relationship operations with use_global_db; node identities
	// carry an optional location for composite IKGs.
	It("routes GlobalDB operations with use_global_db and carries node location", func(ctx SpecContext) {
		c, rec := recordingClient(`{"results":[{"id":"gid:r"}]}`)

		rel := capture.Relationship{
			Type:   "OWNS",
			Source: &capture.Node{ExternalID: "p1", Type: "Person", Location: "eu-db"},
			Target: &capture.Node{ExternalID: "c1", Type: "Car"},
		}
		_, err := c.GlobalDB().UpsertRelationships(ctx, rel)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(HaveKeyWithValue("use_global_db", true))
		src := rec.body["relationships"].([]any)[0].(map[string]any)["source"]
		Expect(src).To(HaveKeyWithValue("location", "eu-db"))

		// The plain client must not send the flag at all.
		_, err = c.UpsertRelationships(ctx, rel)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).NotTo(HaveKey("use_global_db"), "use_global_db must be omitted on the non-global client")

		_, err = c.GlobalDB().DeleteRelationships(ctx, rel)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(HaveKeyWithValue("use_global_db", true), "delete use_global_db")
	})
})

var _ = Describe("ReadDataSchema", func() {
	It("GETs the data schema", func(ctx SpecContext) {
		c, rec := recordingClient(`{"node_types":["Person","Org"]}`)

		schema, err := c.ReadDataSchema(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.method).To(Equal(http.MethodGet))
		Expect(rec.path).To(Equal("/data-schema/v1"))
		Expect(schema).To(HaveKey("node_types"))
	})

	// A non-2xx response surfaces as *transport.APIError on every read/write path.
	It("surfaces a non-2xx response as an APIError", func(ctx SpecContext) {
		c := newClient(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `{"message":"denied"}`)
		})

		_, err := c.ReadDataSchema(ctx)
		Expect(err).To(HaveOccurred())
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
		Expect(apiErr.StatusCode).To(Equal(http.StatusForbidden))
	})
})

var _ = Describe("Chunked upserts", func() {
	// batchRecorder answers each request with one result per item in the
	// collection named key and records the batch sizes it saw.
	batchRecorder := func(key, resultID string) (*capture.Client, *[]int) {
		GinkgoHelper()
		var batchSizes []int
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			var body map[string][]json.RawMessage
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &body)
			batchSizes = append(batchSizes, len(body[key]))
			// Echo one result per item.
			results := make([]string, len(body[key]))
			for i := range results {
				results[i] = `{"id":"` + resultID + `"}`
			}
			_, _ = io.WriteString(w, `{"results":[`+strings.Join(results, ",")+`]}`)
		})
		return c, &batchSizes
	}

	// Verifies the streaming->batch replacement: 5 nodes with chunk size 2 must
	// produce 3 requests (2+2+1) and aggregate all 5 results.
	It("splits nodes into chunks and aggregates the results", func(ctx SpecContext) {
		c, batchSizes := batchRecorder("nodes", "gid:n")

		nodes := make([]capture.UpsertNode, 5)
		for i := range nodes {
			nodes[i] = capture.UpsertNode{Node: capture.Node{ExternalID: id(i), Type: "Person"}}
		}

		res, err := c.UpsertNodesChunked(ctx, nodes, 2)
		Expect(err).NotTo(HaveOccurred())
		Expect(sum(*batchSizes)).To(Equal(5), "total nodes sent")
		Expect(*batchSizes).To(Equal([]int{2, 2, 1}))
		Expect(res.Results).To(HaveLen(5), "aggregated results")
	})

	// Mirrors the node variant: 5 relationships with chunk size 2 must produce
	// 3 requests (2+2+1) and aggregate all 5 results.
	It("splits relationships into chunks and aggregates the results", func(ctx SpecContext) {
		c, batchSizes := batchRecorder("relationships", "gid:r")

		rels := make([]capture.Relationship, 5)
		for i := range rels {
			rels[i] = capture.Relationship{
				Type:   "KNOWS",
				Source: &capture.Node{ExternalID: id(i), Type: "Person"},
				Target: &capture.Node{ExternalID: "z", Type: "Person"},
			}
		}

		res, err := c.UpsertRelationshipsChunked(ctx, rels, 2)
		Expect(err).NotTo(HaveOccurred())
		Expect(sum(*batchSizes)).To(Equal(5), "total relationships sent")
		Expect(*batchSizes).To(Equal([]int{2, 2, 1}))
		Expect(res.Results).To(HaveLen(5), "aggregated results")
	})

	// A chunk size <= 0 or > MaxBatchSize is normalized to MaxBatchSize, so a
	// small slice goes out as a single request.
	It("normalizes an out-of-range chunk size to MaxBatchSize", func(ctx SpecContext) {
		var requests int
		c := newClient(func(w http.ResponseWriter, _ *http.Request) {
			requests++
			_, _ = io.WriteString(w, `{"results":[{"id":"gid:n"}]}`)
		})

		nodes := []capture.UpsertNode{{Node: capture.Node{ExternalID: "a", Type: "Person"}}}
		_, err := c.UpsertNodesChunked(ctx, nodes, 0)
		Expect(err).NotTo(HaveOccurred(), "size 0")
		_, err = c.UpsertNodesChunked(ctx, nodes, capture.MaxBatchSize+1)
		Expect(err).NotTo(HaveOccurred(), "size max+1")
		Expect(requests).To(Equal(2), "one request per call")
	})

	// A failing batch stops the chunked upsert and surfaces the API error.
	It("stops on the first failing batch and keeps the partial results", func(ctx SpecContext) {
		var requests int
		c := newClient(func(w http.ResponseWriter, _ *http.Request) {
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
		res, err := c.UpsertNodesChunked(ctx, nodes, 2)
		Expect(err).To(HaveOccurred(), "expected error from failing batch")
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
		Expect(apiErr.StatusCode).To(Equal(http.StatusForbidden))
		Expect(requests).To(Equal(2), "stop after failure")
		Expect(res.Results).To(HaveLen(2), "partial results")
	})
})

// Every batch write must reject an empty batch without touching the server.
var _ = DescribeTable("rejects an empty batch without touching the server",
	func(ctx SpecContext, call func(context.Context, *capture.Client) error) {
		var called bool
		c := newClient(func(http.ResponseWriter, *http.Request) { called = true })

		Expect(call(ctx, c)).To(MatchError(capture.ErrEmptyBatch))
		Expect(called).To(BeFalse(), "server should not be called for empty batch")
	},
	Entry("UpsertNodes", func(ctx context.Context, c *capture.Client) error {
		_, err := c.UpsertNodes(ctx)
		return err
	}),
	Entry("DeleteNodes", func(ctx context.Context, c *capture.Client) error {
		_, err := c.DeleteNodes(ctx)
		return err
	}),
	Entry("DeleteNodeProperties", func(ctx context.Context, c *capture.Client) error {
		_, err := c.DeleteNodeProperties(ctx)
		return err
	}),
	Entry("DeleteNodePropertyMetadata", func(ctx context.Context, c *capture.Client) error {
		_, err := c.DeleteNodePropertyMetadata(ctx)
		return err
	}),
	Entry("UpsertRelationships", func(ctx context.Context, c *capture.Client) error {
		_, err := c.UpsertRelationships(ctx)
		return err
	}),
	Entry("DeleteRelationships", func(ctx context.Context, c *capture.Client) error {
		_, err := c.DeleteRelationships(ctx)
		return err
	}),
	Entry("DeleteRelationshipProperties", func(ctx context.Context, c *capture.Client) error {
		_, err := c.DeleteRelationshipProperties(ctx)
		return err
	}),
	Entry("UpsertNodesChunked", func(ctx context.Context, c *capture.Client) error {
		_, err := c.UpsertNodesChunked(ctx, nil, 2)
		return err
	}),
	Entry("UpsertRelationshipsChunked", func(ctx context.Context, c *capture.Client) error {
		_, err := c.UpsertRelationshipsChunked(ctx, nil, 2)
		return err
	}),
)
