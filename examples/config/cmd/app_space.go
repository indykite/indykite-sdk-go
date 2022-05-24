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

	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"

	config "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"
)

var appSpaceCmd = &cobra.Command{
	Use:   "app_space",
	Short: "Application Spaces operations",
	Long:  `General commands for managing Application Spaces`,
}

var createAppSpaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Application Space",
	Run: func(cmd *cobra.Command, args []string) {
		var name, displayName, customerID string
		fmt.Print("Enter CustomerID in format 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx': ")
		fmt.Scanln(&customerID)

		fmt.Print("Enter name (slug): ")
		fmt.Scanln(&name)
		fmt.Print("Enter Display name: ")
		fmt.Scanln(&displayName)
		var displayNamePb *wrapperspb.StringValue
		if len(displayName) > 0 {
			displayNamePb = wrapperspb.String(displayName)
		}

		resp, err := client.CreateApplicationSpace(context.Background(), &config.CreateApplicationSpaceRequest{
			CustomerId:  customerID,
			Name:        name,
			DisplayName: displayNamePb,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var updateAppSpaceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update given Application Space",
	Run: func(cmd *cobra.Command, args []string) {
		var appSpaceID, customerID string
		fmt.Print("Enter CustomerID in format 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx': ")
		fmt.Scanln(&customerID)
		customerUUID := uuid.Parse(customerID)
		if customerUUID == nil {
			er("failed to parse digitalTwinID, not a valid UUID")
		}

		fmt.Print("Enter AppSpace ID in format 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx': ")
		fmt.Scanln(&appSpaceID)

		fmt.Print("Enter Display name: ")
		displayNamePb := &wrapperspb.StringValue{Value: ""}
		fmt.Scanln(&(displayNamePb.Value))

		resp, err := client.UpdateApplicationSpace(context.Background(), &config.UpdateApplicationSpaceRequest{
			Id:          appSpaceID,
			DisplayName: displayNamePb,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readAppSpaceCmd = &cobra.Command{
	Use:   "read",
	Short: "Read Application Space by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var appSpaceID string

		fmt.Print("Enter AppSpace ID in format 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx': ")
		fmt.Scanln(&appSpaceID)

		resp, err := client.ReadApplicationSpace(context.Background(), &config.ReadApplicationSpaceRequest{
			Identifier: &config.ReadApplicationSpaceRequest_Id{Id: appSpaceID},
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(appSpaceCmd)
	appSpaceCmd.AddCommand(createAppSpaceCmd)
	appSpaceCmd.AddCommand(updateAppSpaceCmd)
	appSpaceCmd.AddCommand(readAppSpaceCmd)
}
