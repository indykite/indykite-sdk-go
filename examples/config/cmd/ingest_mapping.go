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

	config "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"

	config2 "github.com/indykite/jarvis-sdk-go/config"
)

var ingestMappingCmd = &cobra.Command{
	Use:   "ingestmapping",
	Short: "ingest mapping",
	Long:  `ingest mapping`,
}

var (
	addressEntity = &config.IngestMappingConfig_Entity{
		Labels: []string{"Address"},
		ExternalId: &config.IngestMappingConfig_Property{
			SourceName: "addressId",
			MappedName: "externalId",
			IsRequired: true,
			IsPii:      false,
		},
		Properties: []*config.IngestMappingConfig_Property{
			{
				SourceName: "address",
				MappedName: "address",
				IsRequired: true,
				IsPii:      false,
			},
		},
	}

	parentEntity = &config.IngestMappingConfig_Entity{
		Labels: []string{"DigitalTwin"},
		ExternalId: &config.IngestMappingConfig_Property{
			SourceName: "parentId",
			MappedName: "externalId",
			IsRequired: true,
			IsPii:      false,
		},
		Properties: nil,
		Relationships: []*config.IngestMappingConfig_Relationship{{
			ExternalId: "address",
			Type:       "HAS",
			Direction:  config.IngestMappingConfig_DIRECTION_OUTBOUND,
			MatchLabel: "Address",
		},
			{
				ExternalId: "club_contact",
				Type:       "CONTACT_PERSON",
				Direction:  config.IngestMappingConfig_DIRECTION_OUTBOUND,
				MatchLabel: "Club",
			},
		},
	}

	playerEntity = &config.IngestMappingConfig_Entity{
		Labels: []string{"Player"},
		ExternalId: &config.IngestMappingConfig_Property{
			SourceName: "playerId",
			MappedName: "externalId",
			IsRequired: true,
			IsPii:      false,
		},
		Properties: []*config.IngestMappingConfig_Property{
			{
				SourceName: "firstname",
				MappedName: "firstname",
				IsRequired: true,
				IsPii:      false,
			},
			{
				SourceName: "gender",
				MappedName: "gender",
				IsRequired: true,
				IsPii:      false,
			},
			{
				SourceName: "yearOfBirth",
				MappedName: "yearOfBirth",
				IsRequired: false,
				IsPii:      false,
			},
			{
				SourceName: "sizeTop",
				MappedName: "sizeTop",
				IsRequired: false,
				IsPii:      false,
			},
			{
				SourceName: "sizeBottom",
				MappedName: "sizeBottom",
				IsRequired: false,
				IsPii:      false,
			},
			{
				SourceName: "sizeShoe",
				MappedName: "sizeShoe",
				IsRequired: false,
				IsPii:      false,
			},
		},
		Relationships: []*config.IngestMappingConfig_Relationship{{
			ExternalId: "subscriptionId",
			Type:       "HAS",
			Direction:  config.IngestMappingConfig_DIRECTION_OUTBOUND,
			MatchLabel: "Subscription",
		},
			{
				ExternalId: "parentId",
				Type:       "CHILD",
				Direction:  config.IngestMappingConfig_DIRECTION_INBOUND,
				MatchLabel: "DigitalTwin",
			},
		},
	}

	subscriptionEntity = &config.IngestMappingConfig_Entity{
		Labels: []string{"Subscription"},
		ExternalId: &config.IngestMappingConfig_Property{
			SourceName: "subscriptionId",
			MappedName: "externalId",
			IsRequired: true,
			IsPii:      false,
		},
		Properties: nil,
		Relationships: []*config.IngestMappingConfig_Relationship{
			{
				ExternalId: "clubId",
				Type:       "TO",
				Direction:  config.IngestMappingConfig_DIRECTION_OUTBOUND,
				MatchLabel: "Club",
			},
		},
	}

	clubEntity = &config.IngestMappingConfig_Entity{
		Labels: []string{"Club"},
		ExternalId: &config.IngestMappingConfig_Property{
			SourceName: "clubId",
			MappedName: "externalId",
			IsRequired: true,
			IsPii:      false,
		},
		Properties:    nil,
		Relationships: nil,
	}

	orderEntity = &config.IngestMappingConfig_Entity{
		Labels: []string{"Order"},
		ExternalId: &config.IngestMappingConfig_Property{
			SourceName: "orderId",
			MappedName: "externalId",
			IsRequired: true,
			IsPii:      false,
		},
		Properties: nil,
		Relationships: []*config.IngestMappingConfig_Relationship{{
			ExternalId: "addressId",
			Type:       "SHIP_TO",
			Direction:  config.IngestMappingConfig_DIRECTION_OUTBOUND,
			MatchLabel: "Address",
		},
			{
				ExternalId: "subscriptionId",
				Type:       "IN",
				Direction:  config.IngestMappingConfig_DIRECTION_OUTBOUND,
				MatchLabel: "Subscription",
			},
		},
	}

	entities = []*config.IngestMappingConfig_Entity{
		addressEntity, parentEntity, playerEntity, subscriptionEntity, clubEntity, orderEntity,
	}

	mapping = config2.Mapping{Name: "goodsports", Description: "Good Sports mapping example", Entities: entities}
)

var createIngestMappingConfig = &cobra.Command{
	Use:   "create",
	Short: "Create ingest mapping",
	Run: func(cmd *cobra.Command, args []string) {

		location := "locationGID"
		resp, err := client.CreateIngestMapping(context.Background(), location, mapping)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var getIngestMappingConfig = &cobra.Command{
	Use:   "get",
	Short: "Get ingest mapping",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.GetIngestMapping(context.Background(), "")
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var deleteIngestMappingConfig = &cobra.Command{
	Use:   "delete",
	Short: "Delete ingest mapping",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.DeleteIngestMapping(context.Background(), "")
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var updateIngestMappingConfig = &cobra.Command{
	Use:   "update",
	Short: "Update ingest mapping",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.UpdateIngestMapping(context.Background(), "", mapping)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(ingestMappingCmd)
	ingestMappingCmd.AddCommand(createIngestMappingConfig)
	ingestMappingCmd.AddCommand(getIngestMappingConfig)
	ingestMappingCmd.AddCommand(updateIngestMappingConfig)
	ingestMappingCmd.AddCommand(deleteIngestMappingConfig)
}
