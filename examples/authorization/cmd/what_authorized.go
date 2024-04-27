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
	"context"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/spf13/cobra"

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
)

var whatAuthorizedCmd = &cobra.Command{
	Use:   "what_authorized",
	Short: "What Authorized operations",
	Long:  "General commands for What Authorized",
}

var whatWithTokenCmd = &cobra.Command{
	Use:   "with_token",
	Short: "What Authorized by token",
	Long:  "List resources based on provided type for a subject based on token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken string
		fmt.Scanln(&accessToken)

		resourceTypes := []*authorizationpb.WhatAuthorizedRequest_ResourceType{
			{Type: "TypeA"},
			{Type: "TypeB", Actions: []string{"ACTION"}},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.WhatAuthorizedByToken(
			context.Background(),
			accessToken,
			resourceTypes,
			inputParams,
			policyTags,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var whatWithDigitalTwinCmd = &cobra.Command{
	Use:   "with_dt",
	Short: "What Authorized by digital twin",
	Long:  "List resources based on provided type for a subject based on digital twin id",
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)

		digitalTwin := &authorizationpb.DigitalTwin{
			Id: digitalTwinID,
		}

		resourceTypes := []*authorizationpb.WhatAuthorizedRequest_ResourceType{
			{Type: "TypeA"},
			{Type: "TypeB", Actions: []string{"ACTION"}},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.WhatAuthorized(
			context.Background(),
			digitalTwin,
			resourceTypes,
			inputParams,
			policyTags,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var whatWithPropertyCmd = &cobra.Command{
	Use:   "with_prop",
	Short: "What Authorized by digital twin property",
	Long:  "List resources based on provided type for a subject based on digital twin property",
	Run: func(cmd *cobra.Command, args []string) {
		var propertyType, propertyValue string
		fmt.Print("Enter property type: ")
		fmt.Scanln(&propertyType)
		fmt.Print("Enter property value: ")
		fmt.Scanln(&propertyValue)

		resourceTypes := []*authorizationpb.WhatAuthorizedRequest_ResourceType{
			{Type: "TypeA"},
			{Type: "TypeB", Actions: []string{"ACTION"}},
		}

		property := &authorizationpb.Property{
			Type: propertyType,
			Value: &objects.Value{
				Value: &objects.Value_StringValue{
					StringValue: propertyValue,
				},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.WhatAuthorizedByProperty(
			context.Background(),
			property,
			resourceTypes,
			inputParams,
			policyTags,
			retry.WithMax(2),
		)

		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var whatWithExternalIDCmd = &cobra.Command{
	Use:   "with_external_id",
	Short: "What Authorized by external id",
	Long:  "List resources based on provided type for a subject based on external id",
	RunE: func(cmd *cobra.Command, args []string) error {
		var externalID authorizationpb.ExternalID
		fmt.Print("Enter digital twin type: ")
		fmt.Scanln(&(externalID.Type))
		fmt.Print("Enter external id value: ")
		fmt.Scanln(&(externalID.ExternalId))

		resourceTypes := []*authorizationpb.WhatAuthorizedRequest_ResourceType{
			{Type: "TypeA", Actions: []string{"ACTION1", "ACTION2"}},
			{Type: "TypeB", Actions: []string{"ACTION"}},
		}

		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.WhatAuthorizedByExternalID(
			cmd.Context(),
			&externalID,
			resourceTypes,
			inputParams,
			policyTags,
			retry.WithMax(2),
		)
		if err != nil {
			return err
		}

		fmt.Println(jsonp.Format(resp))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whatAuthorizedCmd)
	whatAuthorizedCmd.AddCommand(whatWithTokenCmd)
	whatAuthorizedCmd.AddCommand(whatWithDigitalTwinCmd)
	whatAuthorizedCmd.AddCommand(whatWithPropertyCmd)
	whatAuthorizedCmd.AddCommand(whatWithExternalIDCmd)
}
