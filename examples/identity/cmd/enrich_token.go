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

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
)

// enrichTokenCmd represents the enrichToken command
var enrichTokenCmd = &cobra.Command{
	Use:   "enrich_token",
	Short: "Enrich token operation",
	Long:  "Enrich session and token with extra properties",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken string
		fmt.Scanln(&accessToken)

		_, err := client.EnrichToken(
			context.Background(),
			&identitypb.EnrichTokenRequest{
				AccessToken: accessToken,
				TokenClaims: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"string_claim": structpb.NewStringValue("string_value"),
						"number_claim": structpb.NewNumberValue(42),
						"bool_claim":   structpb.NewBoolValue(true),
						"null_claim":   structpb.NewNullValue(),
						"map_claim": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
							"key": structpb.NewStringValue("string_value"),
						}}),
						"array_claim": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
							structpb.NewStringValue("string_value"),
						}}),
					},
				},
				SessionClaims: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"string_claim": structpb.NewStringValue("string_value"),
						"number_claim": structpb.NewNumberValue(42),
						"bool_claim":   structpb.NewBoolValue(true),
						"null_claim":   structpb.NewNullValue(),
						"map_claim": structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
							"key": structpb.NewStringValue("string_value"),
						}}),
						"array_claim": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
							structpb.NewStringValue("string_value"),
						}}),
					},
				},
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}

		fmt.Println("successfully enriched token and session")

	},
}

func init() {
	rootCmd.AddCommand(enrichTokenCmd)
}
