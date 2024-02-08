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

	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
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

// This example demonstrates how to use the Identity Knowledge Client to do an IdentityKnowledgeRead query.
func ExampleClient_IdentityKnowledgeRead() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	query := "MATCH (n:Person)-[:HAS]->(s:Store) WHERE n.external_id = $external_id"
	params := map[string]*objects.Value{
		"external_id": {
			Type: &objects.Value_StringValue{StringValue: "1234"},
		},
	}
	returns := []*knowledgepb.Return{
		{
			Variable: "n",
		},
	}

	response, err := client.IdentityKnowledgeRead(context.Background(), query, params, returns)
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to get an Identity by its id.
func ExampleClient_GetIdentityByID() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.GetIdentityByID(context.Background(), "gid:AAAAAmluZHlraURlgAABDwAAAAA")
	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to get an Identity
// by its externalID + type combination.
func ExampleClient_GetIdentityByIdentifier() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.GetIdentityByIdentifier(context.Background(), &knowledge.Identifier{
		ExternalID: "1234",
		Type:       "Person",
	},
	)

	if err != nil {
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	fmt.Println(json.Format(response))
}

// This example demonstrates how to use the Identity Knowledge Client to list all Identities
// that have a certain property.
func ExampleClient_ListIdentitiesByProperty() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListIdentitiesByProperty(context.Background(), &knowledgeobjects.Property{
		Type: "email",
		Value: &objects.Value{
			Type: &objects.Value_StringValue{
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

// This example demonstrates how to use the Identity Knowledge Client to list all Identities.
func ExampleClient_ListIdentities() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListIdentities(context.Background())
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

// This example demonstrates how to use the Identity Knowledge Client to list all nodes.
func ExampleClient_ListNodes() {
	client, err := knowledge.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	response, err := client.ListNodes(context.Background(), "Resource")
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
		&knowledgeobjects.Property{
			Type: "email",
			Value: &objects.Value{
				Type: &objects.Value_StringValue{
					StringValue: "test@test.com",
				},
			},
		},
		false)
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
