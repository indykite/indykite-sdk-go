// Copyright (c) 2024 IndyKite
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

	tdapb "github.com/indykite/indykite-sdk-go/gen/indykite/trusted_data/access/v1beta1"
)

// planCmd represents the plan command
var consentCmd = &cobra.Command{
	Use:   "consent",
	Short: "Consent operation",
	Long: `General commands for Consent

  This is a sample only.`,
}

var grantConsent = &cobra.Command{
	Use:   "grant",
	Short: "Grant consent by id operation",
	Long:  `Grant consent by id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter User ID (identity node gid): ")
		var userId string
		fmt.Scanln(&userId)

		fmt.Print("Enter Consent ID (Consentconfig node gid): ")
		var consentId string
		fmt.Scanln(&consentId)

		resp, err := client.GrantConsent(
			context.Background(),
			&tdapb.GrantConsentRequest{
				User:           &tdapb.GrantConsentRequest_UserId{UserId: userId},
				ConsentId:      consentId,
				RevokeAfterUse: true,
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var revokeConsent = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke consent by id operation",
	Long:  `Revoke consent by id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter User ID (identity node gid): ")
		var userId string
		fmt.Scanln(&userId)

		fmt.Print("Enter Consent ID (Consentconfig node gid): ")
		var consentId string
		fmt.Scanln(&consentId)

		resp, err := client.RevokeConsent(
			context.Background(),
			&tdapb.RevokeConsentRequest{
				User:      &tdapb.RevokeConsentRequest_UserId{UserId: userId},
				ConsentId: consentId,
			},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var accessConsentedData = &cobra.Command{
	Use:   "access",
	Short: "Access consented data operation",
	Long:  `Access consented data with user id and consent id`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter Consent ID (Consentconfig node gid): ")
		var consentId string
		fmt.Scanln(&consentId)

		fmt.Print("Enter User ID (identity node gid): ")
		var userId string
		fmt.Scanln(&userId)

		resp, err := client.AccessConsentedData(
			context.Background(),
			&tdapb.AccessConsentedDataRequest{
				ConsentId: consentId,
				UserId:    userId,
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
	rootCmd.AddCommand(consentCmd)
	consentCmd.AddCommand(grantConsent)
	consentCmd.AddCommand(revokeConsent)
	consentCmd.AddCommand(accessConsentedData)
}
