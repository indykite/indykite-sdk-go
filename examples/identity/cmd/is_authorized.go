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

var isAuthorizedCmd = &cobra.Command{
	Use:   "is_authorized",
	Short: "Is Authorized operations",
	Long:  "General commands for Is Authorized",
}

var withTokenCmd = &cobra.Command{
	Use:   "with_token",
	Short: "Is Authorized by token",
	Long:  "Check if a digital twin is authorized to perform action based on token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken string
		fmt.Scanln(&accessToken)

		actions := []string{"ACTION"}
		resources := []*identitypb.IsAuthorizedRequest_Resource{
			{
				Id:    "resourceID",
				Label: "Label",
			},
		}

		resp, err := client.IsAuthorizedByToken(
			context.Background(),
			accessToken,
			actions,
			resources,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

	},
}

var withDigitalTwinCmd = &cobra.Command{
	Use:   "with_dt",
	Short: "Is Authorized by digital twin",
	Long:  "Check if a digital twin is authorized to perform action based on digital twin id",
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID, tenantID string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)
		fmt.Print("Enter tenant_id: ")
		fmt.Scanln(&tenantID)

		digitalTwin := &identitypb.DigitalTwin{
			Id:       digitalTwinID,
			TenantId: tenantID,
		}

		actions := []string{"ACTION"}
		resources := []*identitypb.IsAuthorizedRequest_Resource{
			{
				Id:    "resourceID",
				Label: "Label",
			},
		}

		resp, err := client.IsAuthorized(
			context.Background(),
			digitalTwin,
			actions,
			resources,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

	},
}

func init() {
	rootCmd.AddCommand(isAuthorizedCmd)
	isAuthorizedCmd.AddCommand(withTokenCmd)
	isAuthorizedCmd.AddCommand(withDigitalTwinCmd)
}
