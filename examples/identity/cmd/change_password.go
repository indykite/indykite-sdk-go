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

// changePwdCmd represents the plan command
var changeDtPwdTokenCmd = &cobra.Command{
	Use:   "change_password_with_token",
	Short: "Change the password of DigitalTwin",
	Long: `Change the Password of DigitalTwin by token.

  Valid token of the user whos password to change is required.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken, newPassword string
		fmt.Scanln(&accessToken)
		fmt.Print("Enter new password: ")
		fmt.Scanln(&newPassword)

		err := client.ChangeMyPassword(context.Background(), accessToken, newPassword, retry.WithMax(2))
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
	},
}

// changePwdIDCmd represents the plan command
var changeDtPwdCmd = &cobra.Command{
	Use:   "change_password",
	Short: "Change the password of DigitalTwin",
	Long:  `Change the Password of DigitalTwin.`,
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID, tenantID, newPassword string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)
		fmt.Print("Enter tenant_id: ")
		fmt.Scanln(&tenantID)
		fmt.Print("Enter new password: ")
		fmt.Scanln(&newPassword)

		resp, err := client.ChangePasswordOfDigitalTwin(
			context.Background(),
			&identitypb.DigitalTwin{
				Id:       digitalTwinID,
				TenantId: tenantID,
			},
			newPassword,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	dtCmd.AddCommand(changeDtPwdTokenCmd)
	dtCmd.AddCommand(changeDtPwdCmd)
}
