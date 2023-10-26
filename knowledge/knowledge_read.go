// Copyright (c) 2023 IndyKite
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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/indykite/indykite-sdk-go/errors"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta2"
	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
)

// Read sends a READ operation to the Identity Knowledge API, with the desired path and optional conditions.
func (c *Client) Read(
	ctx context.Context,
	path string,
	conditions string,
	inputParams map[string]*knowledgepb.InputParam,
	opts ...grpc.CallOption,
) (*knowledgepb.IdentityKnowledgeResponse, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.IdentityKnowledge(ctx, &knowledgepb.IdentityKnowledgeRequest{
		Path:        path,
		Conditions:  conditions,
		InputParams: inputParams,
		Operation:   knowledgepb.Operation_OPERATION_READ,
	}, opts...)
	return resp, err
}

// GetDigitalTwinByID is a helper function that queries for a DigitalTwin node by its id.
func (c *Client) GetDigitalTwinByID(
	ctx context.Context,
	id string,
	opts ...grpc.CallOption,
) (*knowledgepb.Node, error) {
	return c.getNodeByID(ctx, id, DigitalTwin, opts...)
}

// GetDigitalTwinByIdentifier is a helper function that queries for a DigitalTwin node
// by its externalID + type combination.
func (c *Client) GetDigitalTwinByIdentifier(
	ctx context.Context,
	identifier *Identifier,
	opts ...grpc.CallOption,
) (*knowledgepb.Node, error) {
	return c.getNodeByIdentifier(ctx, DigitalTwin, identifier, opts...)
}

// GetResourceByID is a helper function that queries for a Resource node by its id.
func (c *Client) GetResourceByID(
	ctx context.Context,
	id string,
	opts ...grpc.CallOption,
) (*knowledgepb.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	return c.getNodeByID(ctx, id, Resource, opts...)
}

// GetResourceByIdentifier is a helper function that queries for a Resource node
// by its externalID + type combination.
func (c *Client) GetResourceByIdentifier(
	ctx context.Context,
	identifier *Identifier,
	opts ...grpc.CallOption,
) (*knowledgepb.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	return c.getNodeByIdentifier(ctx, Resource, identifier, opts...)
}

// ListDigitalTwinsByProperty is a helper function that lists all DigitalTwin nodes
// that have the specified property.
func (c *Client) ListDigitalTwinsByProperty(
	ctx context.Context,
	property *knowledgepb.Property,
	opts ...grpc.CallOption) ([]*knowledgepb.Node, error) {
	return c.ListNodesByProperty(ctx, DigitalTwin, property, opts...)
}

// ListResourcesByProperty is a helper function that lists all Resource nodes.
// that have the specified property.
func (c *Client) ListResourcesByProperty(
	ctx context.Context,
	property *knowledgepb.Property,
	opts ...grpc.CallOption) ([]*knowledgepb.Node, error) {
	return c.ListNodesByProperty(ctx, Resource, property, opts...)
}

// ListDigitalTwins is a helper function that lists all DigitalTwin nodes.
func (c *Client) ListDigitalTwins(
	ctx context.Context,
	opts ...grpc.CallOption) ([]*knowledgepb.Node, error) {
	return c.ListNodes(ctx, DigitalTwin, opts...)
}

// ListResources is a helper function that lists all Resource nodes.
// that have the specified property.
func (c *Client) ListResources(
	ctx context.Context,
	opts ...grpc.CallOption) ([]*knowledgepb.Node, error) {
	return c.ListNodes(ctx, Resource, opts...)
}

// ListNodes is a helper function that lists all nodes by type, regardless of whether they are DigitalTwins
// or Resources. The nodeType argument should be in PascalCase.
func (c *Client) ListNodes(
	ctx context.Context,
	nodeType string,
	opts ...grpc.CallOption,
) ([]*knowledgepb.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	path := fmt.Sprintf("(:%s)", nodeType)
	resp, err := c.client.IdentityKnowledge(ctx, &knowledgepb.IdentityKnowledgeRequest{
		Path:      path,
		Operation: knowledgepb.Operation_OPERATION_READ,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return parseMultipleNodesFromPaths(resp.GetPaths())
}

// ListNodesByProperty is a helper function that lists all nodes that have the specified type and property.
func (c *Client) ListNodesByProperty(
	ctx context.Context,
	nodeType string,
	property *knowledgepb.Property,
	opts ...grpc.CallOption,
) ([]*knowledgepb.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	path := fmt.Sprintf("(n:%s)", nodeType)
	conditions := fmt.Sprintf("WHERE n.%s = $%s", property.Key, property.Key)
	params := make(map[string]*knowledgepb.InputParam)
	switch property.Value.Value.(type) {
	case *objects.Value_IntegerValue:
		params[property.Key] = &knowledgepb.InputParam{
			Value: &knowledgepb.InputParam_IntegerValue{
				IntegerValue: property.GetValue().GetIntegerValue()}}
	case *objects.Value_StringValue:
		params[property.Key] = &knowledgepb.InputParam{
			Value: &knowledgepb.InputParam_StringValue{
				StringValue: property.GetValue().GetStringValue()}}
	default:
		return nil, errors.New(codes.InvalidArgument, "only string or integer properties can be used for queries")
	}
	resp, err := c.client.IdentityKnowledge(ctx, &knowledgepb.IdentityKnowledgeRequest{
		Path:        path,
		Conditions:  conditions,
		InputParams: params,
		Operation:   knowledgepb.Operation_OPERATION_READ,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return parseMultipleNodesFromPaths(resp.GetPaths())
}

func (c *Client) getNodeByID(
	ctx context.Context,
	id string,
	nodeType string,
	opts ...grpc.CallOption,
) (*knowledgepb.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	path := fmt.Sprintf("(n:%s)", nodeType)
	conditions := fmt.Sprintf("WHERE n.%s = $%s", ID, ID)
	params := map[string]*knowledgepb.InputParam{
		ID: {
			Value: &knowledgepb.InputParam_StringValue{StringValue: id},
		},
	}
	resp, err := c.client.IdentityKnowledge(ctx, &knowledgepb.IdentityKnowledgeRequest{
		Path:        path,
		Conditions:  conditions,
		InputParams: params,
		Operation:   knowledgepb.Operation_OPERATION_READ,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return parseSingleNodeFromPaths(resp.GetPaths())
}

func (c *Client) getNodeByIdentifier(ctx context.Context,
	nodeType string,
	identifier *Identifier,
	opts ...grpc.CallOption,
) (*knowledgepb.Node, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	path := fmt.Sprintf("(n:%s)", nodeType)
	conditions := fmt.Sprintf(
		"WHERE n.%s = $%s AND n.%s = $%s",
		ExternalID,
		ExternalID,
		Type,
		Type,
	)
	params := map[string]*knowledgepb.InputParam{
		ExternalID: {
			Value: &knowledgepb.InputParam_StringValue{StringValue: identifier.ExternalID},
		},
		Type: {
			Value: &knowledgepb.InputParam_StringValue{StringValue: strings.ToLower(identifier.Type)},
		},
	}
	resp, err := c.client.IdentityKnowledge(ctx, &knowledgepb.IdentityKnowledgeRequest{
		Path:        path,
		Conditions:  conditions,
		InputParams: params,
		Operation:   knowledgepb.Operation_OPERATION_READ,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return parseSingleNodeFromPaths(resp.GetPaths())
}

func parseSingleNodeFromPaths(paths []*knowledgepb.Path) (*knowledgepb.Node, error) {
	switch len(paths) {
	case 0:
		return nil, errors.New(codes.NotFound, "node not found")
	case 1:
		return paths[0].GetNodes()[0], nil
	default:
		// if this happens, a uniqueness constraint in the DB has been violated, this should not happen
		return nil, errors.New(codes.Internal, "unable to complete request")
	}
}

func parseMultipleNodesFromPaths(paths []*knowledgepb.Path) ([]*knowledgepb.Node, error) {
	if len(paths) == 0 {
		return nil, errors.New(codes.NotFound, "no nodes found")
	}
	nodes := make([]*knowledgepb.Node, len(paths))
	for i, p := range paths {
		nodes[i] = p.GetNodes()[0]
	}
	return nodes, nil
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// NodesRecordsWithTypeNode is a helper function that get all nodes of specific type either DigitalTwin
// or Resource and transform them into records.
func (c *Client) NodesRecordsWithTypeNode(
	ctx context.Context,
	nodeType string,
	opts ...grpc.CallOption,
) ([]*ingestpb.Record, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	path := fmt.Sprintf("(:%s)", nodeType)
	resp, err := c.client.IdentityKnowledge(ctx, &knowledgepb.IdentityKnowledgeRequest{
		Path:      path,
		Operation: knowledgepb.Operation_OPERATION_READ,
	}, opts...)
	if err != nil {
		return nil, err
	}
	nodes, err := parseMultipleNodesFromPaths(resp.GetPaths())
	if err != nil {
		return nil, err
	}
	records := []*ingestpb.Record{}
	for _, node := range nodes {
		caser := cases.Title(language.English)
		record := &ingestpb.Record{
			Id: GenerateRandomString(12),
			Operation: &ingestpb.Record_Delete{
				Delete: &ingestpb.DeleteData{
					Data: &ingestpb.DeleteData_Node{
						Node: &ingestpb.NodeMatch{
							ExternalId: node.ExternalId,
							Type:       caser.String(node.Type),
						},
					},
				},
			}, // lint:file-ignore U1000 Ignore report
		}
		records = append(records, record)
	}
	return records, nil
}

// Identifier is the combination of ExternalID and Type
// which uniquely identifies a node in the Identity Knowledge Graph.
type Identifier struct {
	ExternalID string
	Type       string
}

const (
	DigitalTwin = "DigitalTwin"
	Resource    = "Resource"
	ExternalID  = "external_id"
	Type        = "type"
	ID          = "id"
)
