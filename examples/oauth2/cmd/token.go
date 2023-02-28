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

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"

	apicfg "github.com/indykite/jarvis-sdk-go/grpc/config"
	"github.com/indykite/jarvis-sdk-go/oauth2"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate bearer token",
	Long:  `Generate bearer token from application credentials`,
}

// generateBearerToken generates a bearer token from the provided credentials
var generateBearerTokenCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate bearer token",
	Long:  `Generate bearer token from application credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		var credentialsLoaders []apicfg.CredentialsLoader

		credentialsLoaders = append(credentialsLoaders, apicfg.DefaultEnvironmentLoader)
		tokenSource, err := oauth2.GetRefreshableTokenSource(context.Background(), credentialsLoaders)
		if err != nil {
			log.Fatalf("failed to generate token source %v", err)
		}
		token, err := tokenSource.Token()
		if err != nil {
			er(err)
		}
		if token.Type() != "Bearer" {
			er(fmt.Errorf("unsupported token type, must be 'Bearer' but got %s", token.Type()))
		}

		fmt.Println(token.AccessToken)
	},
}

// queryKnowledgeAPI generates an authenticated http client and make a Knowledge API query
var queryKnowledgeAPI = &cobra.Command{
	Use:   "query",
	Short: "Query Knowledge API",
	Long:  `Query the Knowledge API, using a http client authenticated with the application credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		var credentialsLoaders []apicfg.CredentialsLoader

		credentialsLoaders = append(credentialsLoaders, apicfg.DefaultEnvironmentLoader)
		client, err := oauth2.GetHTTPClient(context.Background(), credentialsLoaders)
		if err != nil {
			log.Fatalf("failed to generate http client %v", err)
		}

		cfg, err := apicfg.DefaultEnvironmentLoader(context.Background())
		if err != nil {
			log.Fatalf("unable to load credentials")
		}

		var body bytes.Buffer
		json.NewEncoder(&body).Encode(map[string]interface{}{
			"operationName": "Tests",
			"query":         "query Tests {  tests {    testField  }}",
			"variables":     map[string]interface{}{},
		})

		knowledgeAPIUrl := fmt.Sprintf("%s/knowledge/%s", cfg.BaseURL, cfg.AppSpaceID)
		resp, err := client.Post(knowledgeAPIUrl, "application/json", &body)
		if err != nil {
			log.Fatalf("unable to POST message %s", err)
		}
		b, err := io.ReadAll(resp.Body)
		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(generateBearerTokenCmd)
	tokenCmd.AddCommand(queryKnowledgeAPI)
}
