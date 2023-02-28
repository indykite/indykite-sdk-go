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

package oauth2_test

import (
	"context"
	"fmt"
	"log"

	apicfg "github.com/indykite/jarvis-sdk-go/grpc/config"
	"github.com/indykite/jarvis-sdk-go/oauth2"
)

// This example demonstrates how to generate a token.
func ExampleGetRefreshableTokenSource_default() {
	credentialsLoaders = append(credentialsLoaders, apicfg.DefaultEnvironmentLoader)
	tokenSource, err := oauth2.GetRefreshableTokenSource(context.Background(), credentialsLoaders)
	if err != nil {
		log.Fatalf("failed to generate token source %v", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("failed to get token from tokenSource %v", err)
	}
	if token.Type() != "Bearer" {
		log.Fatalf("unsupported token type, must be 'Bearer' but got %s", token.Type())
	}

	fmt.Println(token.AccessToken)
}

// This example demonstrates how to create an authenticated http client and use it call the Knowledge API.
func ExampleGetHTTPClient_default() {
	credentialsLoaders = append(credentialsLoaders, apicfg.DefaultEnvironmentLoader)
	client, err := oauth2.GetHTTPClient(context.Background(), credentialsLoaders)
	if err != nil {
		log.Fatalf("failed to generate token source %v", err)
	}

	resp, err := client.Get("indykite.com")
	if err != nil {
		log.Fatalf("unable to fetch url")
	}
	fmt.Println(resp.Body)
}
