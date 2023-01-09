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
)

// planCmd represents the plan command
var getDtTokenCmd = &cobra.Command{
	Use:   "get_with_token",
	Short: "DigitalTwin operation",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken string
		fmt.Scanln(&accessToken)

		resp, err := client.GetDigitalTwinByToken(
			context.Background(),
			accessToken,
			[]*identitypb.PropertyMask{
				{Definition: &identitypb.PropertyDefinition{Property: "email"}},
				{Definition: &identitypb.PropertyDefinition{Property: "_test0010"}},
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// planCmd represents the plan command
var getDtCmd = &cobra.Command{
	Use:   "get",
	Short: "DigitalTwin operation",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID, tenantID string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)
		fmt.Print("Enter tenant_id: ")
		fmt.Scanln(&tenantID)

		resp, err := client.GetDigitalTwin(
			context.Background(),
			&identitypb.DigitalTwin{Id: digitalTwinID, TenantId: tenantID},
			[]*identitypb.PropertyMask{
				{Definition: &identitypb.PropertyDefinition{Property: "email"}},
				{Definition: &identitypb.PropertyDefinition{Property: "mobile"}},
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	dtCmd.AddCommand(getDtTokenCmd)
	dtCmd.AddCommand(getDtCmd)
}
