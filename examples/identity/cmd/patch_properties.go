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

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

// patchPropDtTokenCmd represents the patch command
var patchPropDtTokenCmd = &cobra.Command{
	Use:   "patch_with_token",
	Short: "DigitalTwin operation",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken string
		fmt.Scanln(&accessToken)

		resp, err := client.PatchDigitalTwinByToken(
			context.Background(),
			accessToken,
			[]*identitypb.PropertyBatchOperation{
				// {Operation: &identitypb.PropertyBatchOperation_Add{Add: &identitypb.Property{
				//	Definition: &identitypb.PropertyDefinition{Property: "email"},
				//	Value:      &identitypb.Property_ObjectValue{ObjectValue: objects.String("some2@example.com")},
				// }}},
				{Operation: &identitypb.PropertyBatchOperation_Replace{Replace: &identitypb.Property{
					Id:    "54c338a", // 88880010 in hex
					Value: &identitypb.Property_ObjectValue{ObjectValue: objects.String("test_abc")},
				}}},
			},
			false,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		for _, r := range resp {
			fmt.Println(jsonp.Format(r))
		}
	},
}

// patchPropDtCmd represents the patch command
var patchPropDtCmd = &cobra.Command{
	Use:   "patch",
	Short: "DigitalTwin operation",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID, tenantID string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)
		fmt.Print("Enter tenant_id: ")
		fmt.Scanln(&tenantID)

		resp, err := client.PatchDigitalTwin(
			context.Background(),
			&identitypb.DigitalTwin{Id: digitalTwinID, TenantId: tenantID},
			[]*identitypb.PropertyBatchOperation{
				// {Operation: &identitypb.PropertyBatchOperation_Add{Add: &identitypb.Property{
				//	Definition: &identitypb.PropertyDefinition{Property: "email"},
				//	Value:      &identitypb.Property_ObjectValue{ObjectValue: objects.String("someone@example.com")},
				// }}},
				{Operation: &identitypb.PropertyBatchOperation_Replace{Replace: &identitypb.Property{
					Id:    "54c338a", // 88880010 in hex
					Value: &identitypb.Property_ObjectValue{ObjectValue: objects.String("test")},
				}}},
			},
			false,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		for _, r := range resp {
			fmt.Println(jsonp.Format(r))
		}
	},
}

func init() {
	dtCmd.AddCommand(patchPropDtTokenCmd)
	dtCmd.AddCommand(patchPropDtCmd)
}
