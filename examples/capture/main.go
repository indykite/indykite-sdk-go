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

// Command capture demonstrates the Capture (IKG ingest) API: batch upserts and
// deletes of nodes, relationships, properties and property metadata.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

const subcommands = "upsert-nodes upsert-relationships delete-nodes delete-node-properties " +
	"delete-node-property-metadata delete-relationships delete-relationship-properties chunked-upsert"

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("capture", subcommands)
	}

	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	capc := cli.Capture()

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	externalID := fs.String("external-id", "person-1", "node external id")
	nodeType := fs.String("type", "Person", "node type")
	_ = fs.Parse(os.Args[2:])

	person := capture.Node{ExternalID: *externalID, Type: *nodeType}
	car := capture.Node{ExternalID: *externalID + "-car", Type: "Car"}
	owns := capture.Relationship{Type: "OWNS", Source: &person, Target: &car}

	var res *capture.BatchResults
	switch os.Args[1] {
	case "upsert-nodes":
		res, err = capc.UpsertNodes(ctx,
			capture.UpsertNode{
				Node:       person,
				IsIdentity: true,
				Properties: []capture.Property{{
					BaseProperty: capture.BaseProperty{Type: "email", Value: *externalID + "@example.com"},
					Metadata:     &capture.Metadata{Source: "example", AssuranceLevel: 1},
				}},
			},
			capture.UpsertNode{Node: car},
		)

	case "upsert-relationships":
		res, err = capc.UpsertRelationships(ctx, owns)

	case "delete-nodes":
		res, err = capc.DeleteNodes(ctx, person, car)

	case "delete-node-properties":
		res, err = capc.DeleteNodeProperties(ctx, capture.DeleteNodeProperties{
			Node: person, PropertyTypes: []string{"email"},
		})

	case "delete-node-property-metadata":
		res, err = capc.DeleteNodePropertyMetadata(ctx, capture.DeleteNodePropertyMetadata{
			Node: person, PropertyType: "email", MetadataFields: []string{"assurance_level"},
		})

	case "delete-relationships":
		res, err = capc.DeleteRelationships(ctx, owns)

	case "delete-relationship-properties":
		res, err = capc.DeleteRelationshipProperties(ctx, capture.DeleteRelationshipProperties{
			Relationship: owns, PropertyTypes: []string{"status"},
		})

	case "chunked-upsert":
		// Batches above capture.MaxBatchSize are split automatically.
		nodes := make([]capture.UpsertNode, 300)
		for i := range nodes {
			nodes[i] = capture.UpsertNode{Node: capture.Node{
				ExternalID: fmt.Sprintf("%s-%03d", *externalID, i), Type: *nodeType,
			}}
		}
		res, err = capc.UpsertNodesChunked(ctx, nodes, 0)

	default:
		exutil.Usage("capture", subcommands)
	}
	if err != nil {
		exutil.Fatal(err)
	}
	exutil.Print(res)
}
