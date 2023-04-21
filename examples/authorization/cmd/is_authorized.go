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

	authorizationpb "github.com/indykite/jarvis-sdk-go/gen/indykite/authorization/v1beta1"
	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
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

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				Id:      "resourceID",
				Type:    "Type",
				Actions: []string{"ACTION"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}

		resp, err := client.IsAuthorizedByToken(
			context.Background(),
			accessToken,
			resources,
			inputParams,
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

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				Id:      "resourceID",
				Type:    "Type",
				Actions: []string{"ACTION"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}

		resp, err := client.IsAuthorized(
			context.Background(),
			digitalTwin,
			resources,
			inputParams,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

	},
}

var withPropertyCmd = &cobra.Command{
	Use:   "with_prop",
	Short: "Is Authorized by digital twin property",
	Long:  "Check if a digital twin is authorized to perform action based on digital twin property",
	Run: func(cmd *cobra.Command, args []string) {
		var propertyType, propertyValue string
		fmt.Print("Enter property type: ")
		fmt.Scanln(&propertyType)
		fmt.Print("Enter property value: ")
		fmt.Scanln(&propertyValue)

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				Id:      "resourceID",
				Type:    "Type",
				Actions: []string{"ACTION"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}

		propertyFilter := &identitypb.PropertyFilter{
			Type:  propertyType,
			Value: objects.String(propertyValue),
		}

		resp, err := client.IsAuthorizedByProperty(
			context.Background(),
			propertyFilter,
			resources,
			inputParams,
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
	isAuthorizedCmd.AddCommand(withPropertyCmd)
}
