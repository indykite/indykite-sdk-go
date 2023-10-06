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

package knowledge_test

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/protojson"

	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
	api "github.com/indykite/indykite-sdk-go/grpc"
	"github.com/indykite/indykite-sdk-go/knowledge"
)

// This example demonstrates how to create a new Identity Knowledge Client.
func ExampleNewClient_default() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	_ = client.Close()
}

// This example demonstrates how to create a new Identity Knowledge Client.
func ExampleNewClient_options() {
	client, err := knowledge.NewClient(context.Background(),
		api.WithCredentialsJSON([]byte(`{}`)))
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	_ = client.Close()
}

// This example demonstrates how to use the Identity Knowledge Client to do a READ query.
func ExampleClient_Read() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	path := "(n:Person)-[:HAS]->(s:Store)"
	conditions := "WHERE n.external_id = $external_id"
	params := map[string]*knowledgepb.InputParam{
		"external_id": {
			Value: &knowledgepb.InputParam_StringValue{StringValue: "1234"},
		},
	}

	response, err := client.Read(context.Background(), path, conditions, params)
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to get a DigitalTwin by its id.
func ExampleClient_GetDigitalTwinByID() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.GetDigitalTwinByID(context.Background(), "gid:AAAAAmluZHlraURlgAABDwAAAAA")
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to get a Resource by its id.
func ExampleClient_GetResourceByID() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.GetResourceByID(context.Background(), "gid:AAAAAmluZHlraURlgAABDwAAAAA")
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to get a DigitalTwin
// by its externalID + type combination.
func ExampleClient_GetDigitalTwinByIdentifier() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.GetDigitalTwinByIdentifier(context.Background(), &knowledge.Identifier{
		ExternalID: "1234",
		Type:       "Person",
	})
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to get a DigitalTwin
// by its externalID + type combination.
func ExampleClient_GetResourceByIdentifier() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.GetResourceByIdentifier(context.Background(), &knowledge.Identifier{
		ExternalID: "1337",
		Type:       "Store",
	})
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to list all DigitalTwins
// that have a certain property.
func ExampleClient_ListDigitalTwinsByProperty() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListDigitalTwinsByProperty(context.Background(), &knowledgepb.Property{
		Key: "email",
		Value: &objects.Value{
			Value: &objects.Value_StringValue{
				StringValue: "test@test.com",
			},
		},
	},
	)
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	for _, n := range response {
		fmt.Println(json.Format(n))
	}
}

// This example demonstrates how to use the Identity Knowledge Client to list all Resources
// that have a certain property.
func ExampleClient_ListResourcesByProperty() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListResourcesByProperty(context.Background(), &knowledgepb.Property{
		Key: "email",
		Value: &objects.Value{
			Value: &objects.Value_StringValue{
				StringValue: "test@test.com",
			},
		},
	},
	)
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	for _, n := range response {
		fmt.Println(json.Format(n))
	}
}

// This example demonstrates how to use the Identity Knowledge Client to list all Resources.
func ExampleClient_ListResources() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListResources(context.Background())
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	for _, n := range response {
		fmt.Println(json.Format(n))
	}
}

// This example demonstrates how to use the Identity Knowledge Client to list all DigitalTwins.
func ExampleClient_ListDigitalTwins() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListDigitalTwins(context.Background())
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	for _, n := range response {
		fmt.Println(json.Format(n))
	}
}

// This example demonstrates how to use the Identity Knowledge Client to list all nodes with a certain type.
func ExampleClient_ListNodes() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListNodes(context.Background(), "Person")
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	for _, n := range response {
		fmt.Println(json.Format(n))
	}
}

// This example demonstrates how to use the Identity Knowledge Client to list all nodes
// with a certain type and property.
func ExampleClient_ListNodesByProperty() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListNodesByProperty(
		context.Background(),
		"Person",
		&knowledgepb.Property{
			Key: "email",
			Value: &objects.Value{
				Value: &objects.Value_StringValue{
					StringValue: "test@test.com"}}})
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	for _, n := range response {
		fmt.Println(json.Format(n))
	}
}
