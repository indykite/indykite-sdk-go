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

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/capture"
)

// Ingest an identity node with a verified property and link it to an asset.
func ExampleClient_UpsertNodes() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}
	capc := cli.Capture()

	person := capture.Node{ExternalID: "ada", Type: "Person"}
	server := capture.Node{ExternalID: "gpu-7", Type: "Server"}

	if _, err = capc.UpsertNodes(ctx,
		capture.UpsertNode{
			Node:       person,
			IsIdentity: true,
			Properties: []capture.Property{{
				BaseProperty: capture.BaseProperty{Type: "email", Value: "ada@example.com"},
				Metadata:     &capture.Metadata{Source: "hr-system", AssuranceLevel: 2},
			}},
		},
		capture.UpsertNode{Node: server},
	); err != nil {
		return
	}

	if _, err = capc.UpsertRelationships(ctx, capture.Relationship{
		Type: "CAN_USE", Source: &person, Target: &server,
	}); err != nil {
		return
	}
}

// Batches above MaxBatchSize are split transparently — the REST replacement
// for the gRPC SDK's streaming ingest.
func ExampleClient_UpsertNodesChunked() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	nodes := make([]capture.UpsertNode, 1000)
	for i := range nodes {
		nodes[i] = capture.UpsertNode{Node: capture.Node{
			ExternalID: "asset-" + string(rune('0'+i%10)), Type: "Asset",
		}}
	}
	if _, err = cli.Capture().UpsertNodesChunked(ctx, nodes, 0); err != nil {
		return
	}
}
