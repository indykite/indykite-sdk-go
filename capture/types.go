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

// Node identifies a node by type + external id. On a composite IKG, Location
// optionally routes the node to a specific constituent.
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
type BatchResults struct {
	Results []Result `json:"results"`
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
	UseGlobalDB   bool           `json:"use_global_db,omitempty"`
}

type deleteRelationshipsRequest struct {
	Relationships []Relationship `json:"relationships"`
	UseGlobalDB   bool           `json:"use_global_db,omitempty"`
}

type deleteRelationshipPropertiesRequest struct {
	Relationships []DeleteRelationshipProperties `json:"relationships"`
	UseGlobalDB   bool                           `json:"use_global_db,omitempty"`
}
