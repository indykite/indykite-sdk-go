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

		accessToken = "eyJhbGciOiJFUzI1NiIsImN0eSI6IkpXVCIsImtpZCI6ImtRWnIyYUk1TUUwQ0o1ejR3U1AwQk9oNkRNOTI2QTVla2tfLUYtYmJBVnciLCJub25jZSI6InRjYmRXdHBBRGdGSjVOQ0U4UHhta0EiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiM2M4ZWZmN2EtYzljYS00NGZjLThhYzgtMWY0ZjYwMTVkYWIwIl0sImV4cCI6MTY2OTYzMjU4OCwiaWF0IjoxNjY5NjI4OTg4LCJpc3MiOiJodHRwczovL2phcnZpcy1kZXYuaW5keWtpdGUuY29tL29hdXRoMi9jZmRiY2Y1Mi00YjEyLTRjYTMtYWMwNi0wZGQwZGE2OThjMTIiLCJqdGkiOiI4ZWJjOGU4My1lOWE1LTQ2OGYtYjcxZC0yYWY4N2EwYTAwMDgjMCIsIm5iZiI6MTY2OTYyODk4OCwic3ViIjoiMWM2MDNhNjEtY2ZjYy00YzRhLWFiNjktMmQ4NmM5NmZlZjk5In0.Xi22X4ncGjY_W52PLJPkg6PZ_OHuJl3g7GEC_fQOdCuMSfTsFMq4s0NLy_YaY2jb9OEFYEMpS5Akvdsyhe1q1A"

		actions := []string{"ACTION"}
		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
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
		resources := []*authorizationpb.IsAuthorizedRequest_Resource{
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
