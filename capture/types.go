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

package capture

import "encoding/json"

const (
	// LocationDefault is the BatchResults.ByLocation key for entries without a
	// location. It is reserved: a Node.Location must not use it.
	LocationDefault = "__default"
	// LocationGlobal is the location of the global constituent of a composite
	// IKG. In BatchResults.ByLocation it holds the global nodes, the node
	// proxies, and the relationships between nodes of different locations.
	LocationGlobal = "global"
)

// Node identifies a node by type + external id. On a composite IKG, Location
// optionally routes the node to a specific constituent.
//
// A relationship has no location of its own: it lives in the constituent of
// its source and target nodes when both have the same Location, and in the
// global constituent when they differ.
type Node struct {
	ExternalID string `json:"external_id"`
	Type       string `json:"type"`
	Location   string `json:"location,omitempty"`
}

// BaseProperty is a property value. Exactly one of Value or ExternalValue is set.
type BaseProperty struct {
	// Value is a string, number, boolean, or array of those. Mutually exclusive
	// with ExternalValue.
	Value any `json:"value,omitempty"`
	// Type is the property type.
	Type string `json:"type"`
	// ExternalValue references an external data resolver. Mutually exclusive
	// with Value.
	ExternalValue string `json:"external_value,omitempty"`
}

// Metadata is optional property metadata.
type Metadata struct {
	CustomMetadata map[string]any `json:"custom_metadata,omitempty"`
	Source         string         `json:"source,omitempty"`
	VerifiedTime   string         `json:"verified_time,omitempty"`
	AssuranceLevel int            `json:"assurance_level,omitempty"`
}

// Property is a node property: a BaseProperty plus optional Metadata.
type Property struct {
	Metadata *Metadata `json:"metadata,omitempty"`
	BaseProperty
}

// Relationship is a typed edge between two nodes, with optional properties.
type Relationship struct {
	Type       string         `json:"type"`
	Source     *Node          `json:"source"`
	Target     *Node          `json:"target"`
	Properties []BaseProperty `json:"properties,omitempty"`
}

// UpsertNode is a node to create or update.
type UpsertNode struct {
	Node
	Labels     []string   `json:"labels,omitempty"`
	Properties []Property `json:"properties,omitempty"`
	IsIdentity bool       `json:"is_identity"`
	// CreateProxy asks for a proxy of the node (external id, type and location
	// only) in the global constituent of a composite IKG. It is ignored for a
	// node in the global location and when no node of the request has a
	// location.
	CreateProxy bool `json:"create_proxy,omitempty"`
}

// PutNode is the body of a single-node upsert (PUT /capture/v1/nodes/{id});
// the node's identity comes from the path id, not the body.
type PutNode struct {
	Labels     []string   `json:"labels,omitempty"`
	Properties []Property `json:"properties,omitempty"`
	IsIdentity bool       `json:"is_identity"`
}

// DeleteNodeProperties names property types to remove from a node.
type DeleteNodeProperties struct {
	Node
	PropertyTypes []string `json:"property_types"`
}

// DeleteNodePropertyMetadata names metadata fields to remove from a node property.
type DeleteNodePropertyMetadata struct {
	Node
	PropertyType   string   `json:"property_type"`
	MetadataFields []string `json:"metadata_fields"`
}

// DeleteRelationshipProperties names property types to remove from a relationship.
type DeleteRelationshipProperties struct {
	Relationship
	PropertyTypes []string `json:"property_types"`
}

// Result is a single operation result (the affected node/relationship gid).
type Result struct {
	ID string `json:"id"`
}

// BatchResults is the response for batch capture operations.
//
// Exactly one of the two fields is filled, depending on the request:
//   - Results, when no node of the request has a location.
//   - ByLocation, when at least one node has a location (composite IKG).
type BatchResults struct {
	// ByLocation lists the affected node or relationship gids per location, in
	// request order. Entries without a location are under LocationDefault.
	// A location where nothing was found has an empty list. For a node upsert
	// or delete, LocationGlobal holds the global nodes followed by the node
	// proxies; the chunked upserts keep that order per chunk only.
	//
	// The key "results" cannot be told apart from Results when its list is
	// empty, so avoid it as a location name.
	ByLocation map[string][]string `json:"-"`
	// Results lists the affected node or relationship gids.
	Results []Result `json:"results"`
}

// UnmarshalJSON decodes both response shapes of the batch operations: the
// {"results":[{"id":...}]} object and the per-location {"<location>":[gid,...]}
// object returned for a composite IKG request.
func (b *BatchResults) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	*b = BatchResults{}
	if fields == nil {
		return nil
	}

	// A "results" list of {"id":...} objects is the plain shape; other keys are
	// ignored so that an added response field does not fail a successful write.
	// A "results" list of gids is a location that happens to have that name; an
	// empty one is only read as the plain shape when it is the sole key.
	if raw, ok := fields["results"]; ok {
		var results []Result
		if err := json.Unmarshal(raw, &results); err == nil && (len(results) > 0 || len(fields) == 1) {
			b.Results = results
			return nil
		}
	}

	b.ByLocation = make(map[string][]string, len(fields))
	for location, raw := range fields {
		var gids []string
		if err := json.Unmarshal(raw, &gids); err != nil {
			return err
		}
		b.ByLocation[location] = gids
	}
	return nil
}

// MarshalJSON encodes the shape UnmarshalJSON decodes, so a BatchResults
// survives a round trip.
func (b BatchResults) MarshalJSON() ([]byte, error) {
	if b.ByLocation != nil {
		return json.Marshal(b.ByLocation)
	}
	results := b.Results
	if results == nil {
		results = []Result{}
	}
	return json.Marshal(struct {
		Results []Result `json:"results"`
	}{Results: results})
}

// merge appends the results of other, keeping the per-location grouping.
//
// The platform picks the response shape per request, so the chunks of one
// upsert can mix both. Once any chunk is per-location, the plain results move
// under LocationDefault, where the platform itself files entries without a
// location, so that still only one of the two fields is filled.
func (b *BatchResults) merge(other *BatchResults) {
	if b.ByLocation == nil && other.ByLocation == nil {
		b.Results = append(b.Results, other.Results...)
		return
	}
	if b.ByLocation == nil {
		b.ByLocation = make(map[string][]string, len(other.ByLocation)+1)
	}
	for _, plain := range [][]Result{b.Results, other.Results} {
		for _, result := range plain {
			b.ByLocation[LocationDefault] = append(b.ByLocation[LocationDefault], result.ID)
		}
	}
	b.Results = nil
	for location, gids := range other.ByLocation {
		merged := b.ByLocation[location]
		if merged == nil {
			// Keep a location where nothing was found as an empty list.
			merged = []string{}
		}
		b.ByLocation[location] = append(merged, gids...)
	}
}

// Request envelopes (match the platform's {"nodes":[...]} / {"relationships":[...]} bodies).

type upsertNodesRequest struct {
	Nodes []UpsertNode `json:"nodes"`
}

type deleteNodesRequest struct {
	Nodes []Node `json:"nodes"`
}

type deleteNodePropertiesRequest struct {
	Nodes []DeleteNodeProperties `json:"nodes"`
}

type deleteNodePropertyMetadataRequest struct {
	Nodes []DeleteNodePropertyMetadata `json:"nodes"`
}

type upsertRelationshipsRequest struct {
	Relationships []Relationship `json:"relationships"`
}

type deleteRelationshipsRequest struct {
	Relationships []Relationship `json:"relationships"`
}

type deleteRelationshipPropertiesRequest struct {
	Relationships []DeleteRelationshipProperties `json:"relationships"`
}
