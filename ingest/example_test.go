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
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	ingestv2pb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta2"
	api "github.com/indykite/indykite-sdk-go/grpc"
	"github.com/indykite/indykite-sdk-go/ingest"
)

// This example demonstrates how to create a new Ingest Client.
func ExampleNewClient_default() {
	client, err := ingest.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	_ = client.Close()
}

// This example demonstrates how to create a new Ingest Client.
func ExampleNewClient_options() {
	client, err := ingest.NewClient(context.Background(),
		api.WithCredentialsJSON([]byte(`{}`)))
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	_ = client.Close()
}

// This example demonstrates how to use the Ingest Client to stream multiple records.
func ExampleClient_StreamRecords() {
	client, err := ingest.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err = client.OpenStreamClient(ctx)
	if err != nil {
		log.Fatalf("failed to open stream %v", err)
	}

	for _, record := range []*ingestv2pb.Record{record1, record2} {
		err = client.SendRecord(record)
		if err != nil {
			log.Fatalf("failed to send record on stream %v", err)
		}
	}

	for {
		resp, err := client.ReceiveResponse()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("failed to receive responses %v", err)
		}
		json := protojson.MarshalOptions{
			Multiline: true,
		}
		fmt.Println(json.Format(resp))
	}
}

// This example demonstrates how to use the Ingest Client to ingest a single record.
func ExampleClient_IngestRecord() {
	client, err := ingest.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.IngestRecord(context.Background(), record1)
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}
