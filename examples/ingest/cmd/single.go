// Copyright (c) 2022 IndyKite
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

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
)

// single represents the command for ingesting a single record
var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "Ingest a single record using the IndyKite Ingest API",
	Long:  `Ingest a single record using the IndyKite Ingest API.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		record1 := &ingestpb.Record{
			Id: "2",
			Operation: &ingestpb.Record_Upsert{
				Upsert: &ingestpb.UpsertData{
					Data: &ingestpb.UpsertData_Node{
						Node: &ingestpb.Node{
							Type: &ingestpb.Node_Resource{
								Resource: &ingestpb.Resource{
									ExternalId: "0001",
									Type:       "Truck",
									Properties: []*ingestpb.Property{
										{
											Key: "vin",
											Value: &objects.Value{
												Value: &objects.Value_IntegerValue{
													IntegerValue: 1234,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		resp, err := client.IngestRecord(context.Background(), record1)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)
}
