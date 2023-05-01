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

package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"

	ingestpb "github.com/indykite/jarvis-sdk-go/gen/indykite/ingest/v1beta2"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

// streamCmd represents the upload command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream multiple records to the IndyKite Ingest API",
	Long:  `Stream multiple records to the IndyKite Ingest API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		record1 := &ingestpb.Record{
			Id: "1",
			Operation: &ingestpb.Record_Upsert{
				Upsert: &ingestpb.UpsertData{
					Data: &ingestpb.UpsertData_Relation{
						Relation: &ingestpb.Relation{
							Match: &ingestpb.RelationMatch{
								SourceMatch: &ingestpb.NodeMatch{
									ExternalId: "0000",
									Type:       "Employee",
								},
								TargetMatch: &ingestpb.NodeMatch{
									ExternalId: "0001",
									Type:       "Truck",
								},
								Type: "SERVICES",
							},
							Properties: []*ingestpb.Property{
								{
									Key: "since",
									Value: &objects.Value{
										Value: &objects.Value_StringValue{
											StringValue: "production",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		record2 := &ingestpb.Record{
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

		records := []*ingestpb.Record{
			record1, record2,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		err := client.OpenStreamClient(ctx)
		if err != nil {
			log.Fatalf("failed to open ingest stream %v", err)
		}
		for _, record := range records {
			err := client.SendRecord(record)
			if err != nil {
				log.Fatalf("failed to send record %v", err)
			}
			resp, err := client.ReceiveResponse()
			if err != nil {
				log.Fatalf("failed to receive response %v", err)
			}
			fmt.Println(jsonp.Format(resp))
		}

		err = client.Close()
		if err != nil {
			log.Fatalf("failed to close ingest stream %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
}
