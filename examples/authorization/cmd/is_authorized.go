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

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
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
		if _, err := fmt.Scanln(&accessToken); err != nil {
			fmt.Println("Error reading accessToken:", err)
		}
		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				ExternalId: "resourceID",
				Type:       "Type",
				Actions:    []string{"ACTION"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.IsAuthorizedByToken(
			context.Background(),
			accessToken,
			resources,
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

var withDigitalTwinCmd = &cobra.Command{
	Use:   "with_dt",
	Short: "Is Authorized by digital twin",
	Long:  "Check if a digital twin is authorized to perform action based on digital twin id",
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID string
		fmt.Print("Enter digital_twin_id: ")
		if _, err := fmt.Scanln(&digitalTwinID); err != nil {
			fmt.Println("Error reading digitalTwinID:", err)
		}

		digitalTwin := &authorizationpb.DigitalTwin{
			Id: digitalTwinID,
		}

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				ExternalId: "Truck1",
				Type:       "Truck",
				Actions:    []string{"SUBSCRIBES_TO"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.IsAuthorized(
			context.Background(),
			digitalTwin,
			resources,
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

var withPropertyCmd = &cobra.Command{
	Use:   "with_prop",
	Short: "Is Authorized by digital twin property",
	Long:  "Check if a digital twin is authorized to perform action based on digital twin property",
	Run: func(cmd *cobra.Command, args []string) {
		var propertyType, propertyValue string
		fmt.Print("Enter property type: ")
		if _, err := fmt.Scanln(&propertyType); err != nil {
			fmt.Println("Error reading propertyType:", err)
		}
		fmt.Print("Enter property value: ")
		if _, err := fmt.Scanln(&propertyValue); err != nil {
			fmt.Println("Error reading propertyValue:", err)
		}

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				ExternalId: "Truck2",
				Type:       "Truck",
				Actions:    []string{"SUBSCRIBES_TO"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		property := &authorizationpb.Property{
			Type: propertyType,
			Value: &objects.Value{
				Value: &objects.Value_StringValue{
					StringValue: propertyValue,
				},
			},
		}

		resp, err := client.IsAuthorizedByProperty(
			context.Background(),
			property,
			resources,
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

var withExternalIDCmd = &cobra.Command{
	Use:   "with_external_id",
	Short: "Is Authorized by external id",
	Long:  "Check if a digital twin is authorized to perform action based on external id",
	RunE: func(cmd *cobra.Command, args []string) error {
		var externalID authorizationpb.ExternalID
		fmt.Print("Enter digital twin type: ")
		if _, err := fmt.Scanln(&externalID.Type); err != nil {
			fmt.Println("Error reading externalID:", err)
		}
		fmt.Print("Enter external id value: ")
		if _, err := fmt.Scanln(&externalID.ExternalId); err != nil {
			fmt.Println("Error reading externalID:", err)
		}

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				ExternalId: "Truck4",
				Type:       "Truck",
				Actions:    []string{"SUBSCRIBES_TO"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.IsAuthorizedByExternalID(
			cmd.Context(),
			&externalID,
			resources,
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

var withExternalID2Cmd = &cobra.Command{
	Use:   "with_trust_score",
	Short: "Is Authorized with trust score",
	Long:  "Check if a digital twin is authorized to perform action based on external id",
	RunE: func(cmd *cobra.Command, args []string) error {
		var externalID authorizationpb.ExternalID
		fmt.Print("Enter node type: ")
		if _, err := fmt.Scanln(&externalID.Type); err != nil {
			fmt.Println("Error reading externalID:", err)
		}
		fmt.Print("Enter external id value: ")
		if _, err := fmt.Scanln(&externalID.ExternalId); err != nil {
			fmt.Println("Error reading externalID:", err)
		}

		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
			{
				ExternalId: "Sensor1",
				Type:       "Sensor",
				Actions:    []string{"CAN_USE"},
			},
		}
		inputParams := map[string]*authorizationpb.InputParam{}
		var policyTags []string

		resp, err := client.IsAuthorizedByExternalID(
			cmd.Context(),
			&externalID,
			resources,
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
	rootCmd.AddCommand(isAuthorizedCmd)
	isAuthorizedCmd.AddCommand(withTokenCmd)
	isAuthorizedCmd.AddCommand(withDigitalTwinCmd)
	isAuthorizedCmd.AddCommand(withPropertyCmd)
	isAuthorizedCmd.AddCommand(withExternalIDCmd)
	isAuthorizedCmd.AddCommand(withExternalID2Cmd)
}
