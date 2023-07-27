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

	"github.com/spf13/cobra"

	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

var serviceAccountCredentialCmd = &cobra.Command{
	Use:   "service_account_credential",
	Short: "Service Account Credential operations",
	Long:  `General commands for managing Service Account Credentials`,
}

var createServiceAccountCredentialCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Service Account Credential",
	Run: func(cmd *cobra.Command, args []string) {
		var displayName, serviceAccountID string

		fmt.Print("Enter ServiceAccountID in gid format:")
		fmt.Scanln(&serviceAccountID)

		fmt.Print("Enter Display name: ")
		fmt.Scanln(&displayName)

		resp, err := client.RegisterServiceAccountCredential(context.Background(), &configpb.RegisterServiceAccountCredentialRequest{
			ServiceAccountId: serviceAccountID,
			DisplayName:      displayName,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readServiceAccountCredentialCmd = &cobra.Command{
	Use:   "read",
	Short: "Read Service Account Credential",
	Run: func(cmd *cobra.Command, args []string) {
		var credentialID string

		fmt.Print("Enter credentialID in gid format:")
		fmt.Scanln(&credentialID)

		resp, err := client.ReadServiceAccountCredential(context.Background(), &configpb.ReadServiceAccountCredentialRequest{
			Id: credentialID,
		})

		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var deleteServiceAccountCredentialCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Service Account Credential",
	Run: func(cmd *cobra.Command, args []string) {
		var credentialID string

		fmt.Print("Enter credentialID in gid format:")
		fmt.Scanln(&credentialID)

		resp, err := client.DeleteServiceAccountCredential(context.Background(), &configpb.DeleteServiceAccountCredentialRequest{
			Id: credentialID,
		})

		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(serviceAccountCredentialCmd)
	serviceAccountCredentialCmd.AddCommand(createServiceAccountCredentialCmd)
	serviceAccountCredentialCmd.AddCommand(readServiceAccountCredentialCmd)
	serviceAccountCredentialCmd.AddCommand(deleteServiceAccountCredentialCmd)
}
