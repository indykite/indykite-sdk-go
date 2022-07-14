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

package ingest_test

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"

	ingestpb "github.com/indykite/jarvis-sdk-go/gen/indykite/ingest/v1beta1"

	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"

	api "github.com/indykite/jarvis-sdk-go/grpc"

	"github.com/indykite/jarvis-sdk-go/ingest"
)

// This example demonstrates how to create a new Ingest Client.
func ExampleNewClient_default() {
	client, err := ingest.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
}

// This example demonstrates how to create a new Ingest Client.
func ExampleNewClient_options() {
	client, err := ingest.NewClient(context.Background(),
		api.WithCredentialsJSON([]byte(`{}`)))
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
}

// This example demonstrates how to use an Ingest Client to stream records.
func ExampleClient_StreamRecords() {
	client, err := ingest.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	record := map[string]*objects.Value{
		"SomeKey":      objects.Int64(12345),
		"SomeOtherKey": objects.String("SomeValue"),
	}
	records := []*ingestpb.Record{
		{
			ExternalId: "SomeKey",
			Data:       record,
		},
	}
	responses, err := client.StreamRecords("gid:AAAAFBtaAlxjDE8GuIWAPEFoSPs", records)
	if err != nil {
		// nolint:gocritic
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}

	for _, response := range responses {
		fmt.Println(json.Format(response))
	}
}
