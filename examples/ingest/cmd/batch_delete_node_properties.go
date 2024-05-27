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

// batch delete node property represents the command for deleting up to 250 nodes
var batchDeleteNodePropertiesCmd = &cobra.Command{
	Use:   "batch_delete_node_properties",
	Short: "Delete bunch of records using the IndyKite Ingest API",
	Long:  `Delete bunch of records using the IndyKite Ingest API.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		nodeMatch1 := &ingestpb.NodeMatch{
			ExternalId: "0000",
			Type:       "Employee",
		}

		nodeMatch2 := &ingestpb.NodeMatch{
			ExternalId: "0001",
			Type:       "Truck",
		}

		nodeProperties := []*ingestpb.DeleteData_NodePropertyMatch{
			{
				Match:        nodeMatch1,
				PropertyType: "PropertyType1",
			},
			{
				Match:        nodeMatch2,
				PropertyType: "PropertyType2",
			},
		}
		resp, err := client.BatchDeleteNodeProperties(context.Background(), nodeProperties)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(batchDeleteNodePropertiesCmd)
}
