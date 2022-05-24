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

package identity_test

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/protobuf/encoding/protojson"

	api "github.com/indykite/jarvis-sdk-go/grpc"

	"github.com/indykite/jarvis-sdk-go/identity"
)

// This example demonstrates how to create a new Identity Client.
func ExampleNewClient_default() {
	client, err := identity.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
}

// This example demonstrates how to create a new Identity Client.
func ExampleNewClient_options() {
	client, err := identity.NewClient(context.Background(),
		api.WithCredentialsJSON([]byte(`{}`)))
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
}

// This example demonstrates the use of a identity client to introspect the
// access_token from the request.
func ExampleClient_IntrospectToken() {
	client, err := identity.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	/* #nosec */
	token := "JWT TOKEN HERE"
	tenant, err := client.IntrospectToken(context.Background(), token, retry.WithMax(2))
	if err != nil {
		// nolint:gocritic
		log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
	}
	json := protojson.MarshalOptions{
		Multiline: true,
	}
	println(json.Format(tenant))
}

// This example demonstrates the use of an identity client to create a new
// invitation and notify invitee via email.
func ExampleClient_CreateEmailInvitation() {
	client, err := identity.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	err = client.CreateEmailInvitation(context.Background(),
		"test@example.com",
		"696e6479-6b69-4465-8000-030100000002",
		"my-reference",
		time.Now().AddDate(0, 0, 7), time.Now(),
		map[string]interface{}{
			"lang": "en",
		})
	if err != nil {
		log.Printf("failed to invoke operation on IndyKite Client %v", err)
	}
}
