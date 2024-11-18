// Copyright (c) 2024 IndyKite
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

package knowledge

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/indykite/indykite-sdk-go/errors"
	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
)

// IdentityKnowledgeRead sends a READ operation to the Identity Knowledge API,
// with the desired query, inputParams and returns.
func (c *Client) IdentityKnowledgeRead(
	ctx context.Context,
	query string,
	inputParams map[string]*objects.Value,
	returns []*knowledgepb.Return,
	opts ...grpc.CallOption,
) (*knowledgepb.IdentityKnowledgeReadResponse, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.IdentityKnowledgeRead(ctx, &knowledgepb.IdentityKnowledgeReadRequest{
		Query:       query,
		InputParams: inputParams,
		Returns:     returns,
	}, opts...)
	return resp, err
}

// GetIdentityByID is a helper function that queries for a Identity node by its id.
func (c *Client) GetIdentityByID(
	ctx context.Context,
	id string,
	opts ...grpc.CallOption,
) (*knowledgeobjects.Node, error) {
	return c.GetNodeByID(ctx, id, true, opts...)
}

// GetIdentityByIdentifier is a helper function that queries for a Identity node
// by its externalID + type combination.
func (c *Client) GetIdentityByIdentifier(
	ctx context.Context,
	identifier *Identifier,
	opts ...grpc.CallOption,
) (*knowledgeobjects.Node, error) {
	return c.GetNodeByIdentifier(ctx, identifier, true, opts...)
}

// ListIdentitiesByProperty is a helper function that lists all Identity nodes
// that have the specified property.
func (c *Client) ListIdentitiesByProperty(
	ctx context.Context,
	property *knowledgeobjects.Property,
	opts ...grpc.CallOption) ([]*knowledgeobjects.Node, error) {
	return c.ListNodesByProperty(ctx, property, true, opts...)
}

// ListIdentities is a helper function that lists all Identity nodes.
func (c *Client) ListIdentities(
	ctx context.Context,
	opts ...grpc.CallOption) ([]*knowledgeobjects.Node, error) {
	return c.ListNodes(ctx, "DigitalTwin", opts...)
}

// ListNodes is a helper function that lists all nodes by given node type, regardless of whether they are Identities
// or not.
func (c *Client) ListNodes(
	ctx context.Context,
	nodeType string,
	opts ...grpc.CallOption,
) ([]*knowledgeobjects.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	query := fmt.Sprintf("MATCH (n:%s)", nodeType)
	params := map[string]*objects.Value{}
	resp, err := c.client.IdentityKnowledgeRead(ctx, &knowledgepb.IdentityKnowledgeReadRequest{
		Query:       query,
		InputParams: params,
		Returns: []*knowledgepb.Return{
			{
				Variable: "n",
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp.GetNodes(), nil
}

// ListNodesByProperty is a helper function that lists all nodes that have the specified type and property.
func (c *Client) ListNodesByProperty(
	ctx context.Context,
	property *knowledgeobjects.Property,
	isIdentity bool,
	opts ...grpc.CallOption,
) ([]*knowledgeobjects.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	query := "MATCH (n:Resource)-[:HAS]->(p:Property)"
	if isIdentity {
		query = "MATCH (n:DigitalTwin)-[:HAS]->(p:Property)"
	}
	query = fmt.Sprintf(
		"%s WHERE p.type='%s' and p.value=$%s",
		query,
		property.Type,
		property.Type,
	)

	params := make(map[string]*objects.Value)
	switch property.Value.Type.(type) {
	case *objects.Value_IntegerValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_IntegerValue{
				IntegerValue: property.GetValue().GetIntegerValue()}}
	case *objects.Value_StringValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_StringValue{
				StringValue: property.GetValue().GetStringValue()}}
	case *objects.Value_BoolValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_BoolValue{
				BoolValue: property.GetValue().GetBoolValue()}}
	case *objects.Value_DoubleValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_DoubleValue{
				DoubleValue: property.GetValue().GetDoubleValue()}}
	case *objects.Value_TimeValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_TimeValue{
				TimeValue: property.GetValue().GetTimeValue()}}
	case *objects.Value_DurationValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_DurationValue{
				DurationValue: property.GetValue().GetDurationValue()}}
	case *objects.Value_ArrayValue:
		params[property.Type] = &objects.Value{
			Type: &objects.Value_ArrayValue{
				ArrayValue: property.GetValue().GetArrayValue()}}
	default:
		return nil, errors.New(codes.InvalidArgument, "this type cannot be used for queries")
	}
	resp, err := c.client.IdentityKnowledgeRead(ctx, &knowledgepb.IdentityKnowledgeReadRequest{
		Query:       query,
		InputParams: params,
		Returns: []*knowledgepb.Return{
			{
				Variable: "n",
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp.GetNodes(), nil
}

func (c *Client) GetNodeByID(
	ctx context.Context,
	id string,
	isIdentity bool,
	opts ...grpc.CallOption,
) (*knowledgeobjects.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	query := "MATCH (n:Resource)"
	if isIdentity {
		query = "MATCH (n:DigitalTwin)"
	}
	query = fmt.Sprintf(
		"%s WHERE n.%s=$%s",
		query,
		ID,
		ID,
	)
	params := map[string]*objects.Value{
		ID: {
			Type: &objects.Value_StringValue{StringValue: id},
		},
	}
	resp, err := c.client.IdentityKnowledgeRead(ctx, &knowledgepb.IdentityKnowledgeReadRequest{
		Query:       query,
		InputParams: params,
		Returns: []*knowledgepb.Return{
			{
				Variable: "n",
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	return ParseSingleNodeFromNodes(resp.GetNodes())
}

func (c *Client) GetNodeByIdentifier(ctx context.Context,
	identifier *Identifier,
	isIdentity bool,
	opts ...grpc.CallOption,
) (*knowledgeobjects.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	query := "MATCH (n:Resource)"
	if isIdentity {
		query = "MATCH (n:DigitalTwin)"
	}
	query = fmt.Sprintf(
		"%s WHERE n.%s=$%s AND n.%s=$%s",
		query,
		ExternalID,
		ExternalID,
		Type,
		Type,
	)
	params := map[string]*objects.Value{
		ExternalID: {
			Type: &objects.Value_StringValue{StringValue: identifier.ExternalID},
		},
		Type: {
			Type: &objects.Value_StringValue{StringValue: identifier.Type},
		},
	}
	resp, err := c.client.IdentityKnowledgeRead(ctx, &knowledgepb.IdentityKnowledgeReadRequest{
		Query:       query,
		InputParams: params,
		Returns: []*knowledgepb.Return{
			{
				Variable: "n",
			},
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	return ParseSingleNodeFromNodes(resp.GetNodes())
}

func ParseSingleNodeFromNodes(nodes []*knowledgeobjects.Node) (*knowledgeobjects.Node, error) {
	switch len(nodes) {
	case 0:
		return nil, errors.New(codes.NotFound, "node not found")
	case 1:
		return nodes[0], nil
	default:
		// if this happens, a uniqueness constraint in the DB has been violated, this should not happen
		return nil, errors.New(codes.Internal, "unable to complete request")
	}
}

// Identifier is the combination of ExternalID and Type
// which uniquely identifies a node in the Identity Knowledge Graph.
type Identifier struct {
	ExternalID string
	Type       string
}

const (
	ExternalID = "external_id"
	Type       = "type"
	ID         = "id"
)
