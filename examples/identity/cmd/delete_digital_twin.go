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

// deleteDtTokenCmd represents the delte DigitalTwin command.
var deleteDtTokenCmd = &cobra.Command{
	Use:   "delete_with_token",
	Short: "Deletes the DigitalTwin",
	Long: `Deletes the DigitalTwin with token

  Removes the DigitalTwin from the system and deletes all
  data collected previously.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter access_token: ")
		var accessToken string
		fmt.Scanln(&accessToken)

		resp, err := client.DeleteDigitalTwinByToken(
			context.Background(),
			accessToken,
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// deleteDtCmd represents the delte DigitalTwin command.
var deleteDtCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes the DigitalTwin",
	Long: `Deletes the DigitalTwin

  Removes the DigitalTwin from the system and deletes all
  data collected previously.`,
	Run: func(cmd *cobra.Command, args []string) {
		var digitalTwinID, tenantID string
		fmt.Print("Enter digital_twin_id: ")
		fmt.Scanln(&digitalTwinID)
		fmt.Print("Enter tenant_id: ")
		fmt.Scanln(&tenantID)

		resp, err := client.DeleteDigitalTwin(
			context.Background(),
			&identitypb.DigitalTwin{Id: digitalTwinID, TenantId: tenantID},
			retry.WithMax(2),
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	dtCmd.AddCommand(deleteDtTokenCmd)
	dtCmd.AddCommand(deleteDtCmd)
}
