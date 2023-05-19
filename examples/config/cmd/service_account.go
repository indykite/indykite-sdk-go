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

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"

	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

var serviceAccountCmd = &cobra.Command{
	Use:   "service_account",
	Short: "Service Account operations",
	Long:  `General commands for managing Service Accounts`,
}

var createServiceAccountCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Service Account",
	Run: func(cmd *cobra.Command, args []string) {
		var name, displayName, customerID string
		fmt.Print("Enter CustomerID in gid format:")
		fmt.Scanln(&customerID)

		fmt.Print("Enter name (slug): ")
		fmt.Scanln(&name)
		fmt.Print("Enter Display name: ")
		fmt.Scanln(&displayName)
		var displayNamePb *wrapperspb.StringValue
		if len(displayName) > 0 {
			displayNamePb = wrapperspb.String(displayName)
		}

		resp, err := client.CreateServiceAccount(context.Background(), &configpb.CreateServiceAccountRequest{
			Location:    customerID,
			Name:        name,
			DisplayName: displayNamePb,
			Role:        "all_editor",
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(serviceAccountCmd)
	serviceAccountCmd.AddCommand(createServiceAccountCmd)
}
