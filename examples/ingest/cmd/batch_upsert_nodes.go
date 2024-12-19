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

package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
)

// batch upsert node represents the command for ingesting up to 250 nodes
var batchUpsertNodesCmd = &cobra.Command{
	Use:   "batch_upsert_nodes",
	Short: "Ingest bunch of records using the IndyKite Ingest API",
	Long:  `Ingest bunch of records using the IndyKite Ingest API.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		node1 := &knowledgeobjects.Node{
			ExternalId: "Truck6",
			Type:       "Truck",
			IsIdentity: false,
			Properties: []*knowledgeobjects.Property{
				{
					Type: "color",
					Value: &objects.Value{
						Type: &objects.Value_StringValue{
							StringValue: "blue",
						},
					},
				},
				{
					Type: "vin",
					Value: &objects.Value{
						Type: &objects.Value_StringValue{
							StringValue: "iouiihuhu",
						},
					},
				},
				{
					Type: "echo",
					ExternalValue: &knowledgeobjects.ExternalValue{
						Resolver: &knowledgeobjects.ExternalValue_Id{
							Id: "gid:AAAAIYD2NoHgI0YGrAeFN-dz4K8",
						},
					},
				},
			},
			Tags: []string{},
		}

		nodes := []*knowledgeobjects.Node{
			node1,
		}
		resp, err := client.BatchUpsertNodes(context.Background(), nodes)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(batchUpsertNodesCmd)
}
