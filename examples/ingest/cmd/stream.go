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
	"fmt"
	"log"

	ingestpb "github.com/indykite/jarvis-sdk-go/gen/indykite/ingest/v1beta1"

	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"

	"github.com/spf13/cobra"
)

// streamCmd represents the upload command
var streamCmd = &cobra.Command{
	Use:   "stream [mappingID]",
	Short: "Stream records to the IndyKite ingest service",
	Long: `Stream records to the IndyKite Ingest Engine. A data mapping must have been set up and referenced
		with the mappingID.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		addressRecord := &ingestpb.Record{
			Id:         "1",
			ExternalId: "addressId",
			Data: map[string]*objects.Value{
				"addressId": objects.Int64(123),
				"address":   objects.String("1 Ferry Building, San Francisco, CA 94105, United States"),
			},
		}
		parentRecord := &ingestpb.Record{
			Id:         "2",
			ExternalId: "parentId",
			Data: map[string]*objects.Value{
				"parentId": objects.Int64(123456789),
				"address":  objects.Int64(123),
			},
		}

		playerRecord1 := &ingestpb.Record{
			Id:         "3",
			ExternalId: "playerId",
			Data: map[string]*objects.Value{
				"playerId":       objects.Int64(987654321),
				"firstname":      objects.String("Kid Doe"),
				"gender":         objects.String("male"),
				"yearOfBirth":    objects.Int64(2006),
				"sizeTop":        objects.String("M"),
				"sizeBottom":     objects.String("M"),
				"sizeShoe":       objects.Int64(36),
				"subscriptionId": objects.Int64(10101010),
				"parentId":       objects.Int64(123456789),
			},
		}

		subscriptionRecord := &ingestpb.Record{
			Id:         "4",
			ExternalId: "subscriptionId",
			Data: map[string]*objects.Value{
				"subscriptionId": objects.Int64(10101010),
				"clubId":         objects.Int64(1337),
			},
		}

		clubRecord := &ingestpb.Record{
			Id:         "5",
			ExternalId: "clubId",
			Data: map[string]*objects.Value{
				"clubId":         objects.Int64(1337),
				"subscriptionId": objects.Int64(10101010),
			},
		}

		orderRecord := &ingestpb.Record{
			Id:         "6",
			ExternalId: "orderId",
			Data: map[string]*objects.Value{
				"orderId":        objects.Int64(11),
				"addressId":      objects.Int64(123),
				"subscriptionId": objects.Int64(10101010),
			},
		}

		playerRecord2 := &ingestpb.Record{
			Id:         "7",
			ExternalId: "playerId",
			Data: map[string]*objects.Value{
				"playerId":       objects.Int64(11223344),
				"firstname":      objects.String("Kiddo Doe"),
				"gender":         objects.String("female"),
				"yearOfBirth":    objects.Int64(2005),
				"sizeTop":        objects.String("S"),
				"sizeBottom":     objects.String("S"),
				"sizeShoe":       objects.Int64(32),
				"subscriptionId": objects.Int64(11111111),
				"parentId":       objects.Int64(123456789),
			},
		}

		records := []*ingestpb.Record{
			addressRecord,
			parentRecord,
			playerRecord1,
			subscriptionRecord,
			clubRecord,
			orderRecord,
			playerRecord2,
		}
		responses, err := client.StreamRecords("", records)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		for _, response := range responses {
			fmt.Println(jsonp.Format(response))
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
}
