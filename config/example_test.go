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

package config_test

import (
	"context"
	"log"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/indykite/jarvis-sdk-go/config"
	configpb "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"
	api "github.com/indykite/jarvis-sdk-go/grpc"
)

// This example demonstrates how to create a new Config Client.
func ExampleNewClient_default() {
	client, err := config.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
}

// This example demonstrates how to create a new Config Client.
func ExampleNewClient_options() {
	client, err := config.NewClient(context.Background(),
		api.WithCredentialsJSON([]byte(`{}`)))
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
}

// This example demonstrates how to use a Config Client to create an ingest mapping.
func ExampleClient_CreateIngestMapping() {
	client, err := config.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	mappingBuilder := config.NewMappingBuilder()
	var dts []*configpb.IngestMappingConfig_Entity
	var entities []*configpb.IngestMappingConfig_Entity
	dtBuilder := config.NewDigitalTwinBuilder()
	dtBuilder.ExternalID("nin").
		TenantID("gid:AAAAA2luZHlraURlgAADDwAAAAE").
		Properties([]*configpb.IngestMappingConfig_Property{
			{
				SourceName: "fn",
				MappedName: "firstname",
				IsRequired: false,
			},
		}).Relationships([]*configpb.IngestMappingConfig_Relationship{
		{
			ExternalId: "familynumber",
			Type:       "MEMBER_OF",
			Direction:  configpb.IngestMappingConfig_DIRECTION_OUTBOUND,
			MatchLabel: "Family",
		},
	})

	dts = append(dts, dtBuilder.Build())

	familyBuilder := config.NewEntityBuilder()
	familyBuilder.Labels([]string{"Family"}).ExternalID("familynumber")

	entities = append(entities, familyBuilder.Build())

	mappingBuilder.Name("DSF mapping").DisplayName("some cool display name").
		Description("Description").DigitalTwins(dts).
		Entities(entities)
	mapping := mappingBuilder.Mapping

	location := "locationGID"
	resp, err := client.CreateIngestMapping(context.Background(), location, *mapping)
	if err != nil {
		// nolint:gocritic
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	println(json.Format(resp))
}

// This example demonstrates how to use a Config Client to get an ingest mapping.
func ExampleClient_GetIngestMapping() {
	client, err := config.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	id := "ingest mapping id"
	resp, err := client.GetIngestMapping(context.Background(), id)
	if err != nil {
		// nolint:gocritic
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	println(json.Format(resp))
}

// This example demonstrates how to use a Config Client to delete an ingest mapping.
func ExampleClient_DeleteIngestMapping() {
	client, err := config.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	id := "ingest mapping id"
	resp, err := client.DeleteIngestMapping(context.Background(), id)
	if err != nil {
		// nolint:gocritic
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	println(json.Format(resp))
}
