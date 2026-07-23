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

// Package capture is the client for the IndyKite Capture (IKG ingest) API
// (/capture/v1/*) and the data-schema read (/data-schema/v1). It runs on the
// runtime plane (App Agent token) and is a thin facade over a *transport.Client.
//
//	cap := capture.NewClient(client)
//	_, err := cap.UpsertNodes(ctx, capture.UpsertNode{
//	    Node: capture.Node{ExternalID: "ada", Type: "Person"},
//	})
//
// Capture caps a batch at MaxBatchSize entries; UpsertNodesChunked /
// UpsertRelationshipsChunked split larger slices automatically.
package capture

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/indykite/indykite-sdk-go/transport"
)

const (
	pathNodes              = "/capture/v1/nodes"
	pathNodesDelete        = "/capture/v1/nodes/delete"
	pathNodePropsDelete    = "/capture/v1/nodes/properties/delete"
	pathNodePropMetaDelete = "/capture/v1/nodes/properties/metadata/delete"
	pathRelationships      = "/capture/v1/relationships"
	pathRelDelete          = "/capture/v1/relationships/delete"
	pathRelPropsDelete     = "/capture/v1/relationships/properties/delete"
	pathDataSchema         = "/data-schema/v1"

	// MaxBatchSize is the maximum number of entries the Capture API accepts per
	// request.
	MaxBatchSize = 250
)

// ErrEmptyBatch is returned when a write is called with no entries.
var ErrEmptyBatch = errors.New("capture: batch must contain at least one entry")

// Client calls the Capture and data-schema APIs.
type Client struct {
	t *transport.Client
	// useGlobalDB routes relationship operations to the global constituent of
	// a composite IKG (sent as use_global_db).
	useGlobalDB bool
}

// NewClient builds a Capture client over the shared transport.
func NewClient(t *transport.Client) *Client {
	return &Client{t: t}
}

// GlobalDB returns a view of the client whose relationship operations set
// use_global_db, targeting the global constituent of a composite IKG. Node
// operations route per node via Node.Location instead.
func (c *Client) GlobalDB() *Client {
	c2 := *c
	c2.useGlobalDB = true
	return &c2
}

func (c *Client) post(ctx context.Context, path string, body any) (*BatchResults, error) {
	var out BatchResults
	if err := c.t.Do(ctx, http.MethodPost, path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpsertNodes creates or updates up to MaxBatchSize nodes.
func (c *Client) UpsertNodes(ctx context.Context, nodes ...UpsertNode) (*BatchResults, error) {
	if len(nodes) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathNodes, upsertNodesRequest{Nodes: nodes})
}

// UpsertSingleNode creates or updates one node addressed by id, which is
// either a gid ("gid:...") or "<Type>:<external_id>". Unlike UpsertNodes, the
// node's identity comes from the id, not the body.
func (c *Client) UpsertSingleNode(ctx context.Context, id string, node PutNode) (*BatchResults, error) {
	var out BatchResults
	if err := c.t.Do(ctx, http.MethodPut, pathNodes+"/"+url.PathEscape(id), node, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNodes deletes up to MaxBatchSize nodes.
func (c *Client) DeleteNodes(ctx context.Context, nodes ...Node) (*BatchResults, error) {
	if len(nodes) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathNodesDelete, deleteNodesRequest{Nodes: nodes})
}

// DeleteNodeProperties removes named property types from nodes.
func (c *Client) DeleteNodeProperties(ctx context.Context, items ...DeleteNodeProperties) (*BatchResults, error) {
	if len(items) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathNodePropsDelete, deleteNodePropertiesRequest{Nodes: items})
}

// DeleteNodePropertyMetadata removes named metadata fields from node properties.
func (c *Client) DeleteNodePropertyMetadata(
	ctx context.Context,
	items ...DeleteNodePropertyMetadata,
) (*BatchResults, error) {
	if len(items) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathNodePropMetaDelete, deleteNodePropertyMetadataRequest{Nodes: items})
}

// UpsertRelationships creates or updates up to MaxBatchSize relationships.
func (c *Client) UpsertRelationships(ctx context.Context, rels ...Relationship) (*BatchResults, error) {
	if len(rels) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathRelationships, upsertRelationshipsRequest{Relationships: rels, UseGlobalDB: c.useGlobalDB})
}

// DeleteRelationships deletes up to MaxBatchSize relationships.
func (c *Client) DeleteRelationships(ctx context.Context, rels ...Relationship) (*BatchResults, error) {
	if len(rels) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathRelDelete, deleteRelationshipsRequest{Relationships: rels, UseGlobalDB: c.useGlobalDB})
}

// DeleteRelationshipProperties removes named property types from relationships.
func (c *Client) DeleteRelationshipProperties(
	ctx context.Context,
	items ...DeleteRelationshipProperties,
) (*BatchResults, error) {
	if len(items) == 0 {
		return nil, ErrEmptyBatch
	}
	return c.post(ctx, pathRelPropsDelete, deleteRelationshipPropertiesRequest{
		Relationships: items, UseGlobalDB: c.useGlobalDB,
	})
}

// ReadDataSchema returns the current data schema as a decoded JSON document.
func (c *Client) ReadDataSchema(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	if err := c.t.Do(ctx, http.MethodGet, pathDataSchema, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpsertNodesChunked upserts any number of nodes, splitting into batches of at
// most chunkSize (clamped to 1..MaxBatchSize; 0 means MaxBatchSize) and
// aggregating the results.
func (c *Client) UpsertNodesChunked(ctx context.Context, nodes []UpsertNode, chunkSize int) (*BatchResults, error) {
	if len(nodes) == 0 {
		return nil, ErrEmptyBatch
	}
	agg := &BatchResults{}
	for _, batch := range chunk(nodes, normalizeChunk(chunkSize)) {
		res, err := c.UpsertNodes(ctx, batch...)
		if err != nil {
			return agg, err
		}
		agg.Results = append(agg.Results, res.Results...)
	}
	return agg, nil
}

// UpsertRelationshipsChunked upserts any number of relationships in batches.
func (c *Client) UpsertRelationshipsChunked(
	ctx context.Context,
	rels []Relationship,
	chunkSize int,
) (*BatchResults, error) {
	if len(rels) == 0 {
		return nil, ErrEmptyBatch
	}
	agg := &BatchResults{}
	for _, batch := range chunk(rels, normalizeChunk(chunkSize)) {
		res, err := c.UpsertRelationships(ctx, batch...)
		if err != nil {
			return agg, err
		}
		agg.Results = append(agg.Results, res.Results...)
	}
	return agg, nil
}

func normalizeChunk(size int) int {
	if size <= 0 || size > MaxBatchSize {
		return MaxBatchSize
	}
	return size
}

func chunk[T any](items []T, size int) [][]T {
	var out [][]T
	for start := 0; start < len(items); start += size {
		end := min(start+size, len(items))
		out = append(out, items[start:end])
	}
	return out
}
