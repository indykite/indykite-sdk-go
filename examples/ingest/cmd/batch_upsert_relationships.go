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

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
)

// batch upsert relationships represents the command for ingesting up to 250 relationships.
var batchUpsertRelationshipsCmd = &cobra.Command{
	Use:   "batch_upsert_relationships",
	Short: "Ingest bunch of records using the IndyKite Ingest API",
	Long:  `Ingest bunch of records using the IndyKite Ingest API.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		relationship1 := &ingestpb.Relationship{
			Source: &ingestpb.NodeMatch{
				ExternalId: "tRVeocDOOzNfTIN",
				Type:       "Organization",
			},
			Target: &ingestpb.NodeMatch{
				ExternalId: "Truck5",
				Type:       "Truck",
			},
			Type: "OWNS",
		}

		relationships := []*ingestpb.Relationship{
			relationship1,
		}
		resp, err := client.BatchUpsertRelationships(context.Background(), relationships)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(batchUpsertRelationshipsCmd)
}
