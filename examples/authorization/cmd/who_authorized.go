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

	authorizationpb "github.com/indykite/jarvis-sdk-go/gen/indykite/authorization/v1beta1"
)

var whoAuthorizedCmd = &cobra.Command{
	Use:   "who_authorized",
	Short: "Who Authorized operations",
	Long:  "General commands for Who Authorized",
	Run: func(cmd *cobra.Command, args []string) {
		resources := []*authorizationpb.WhoAuthorizedRequest_Resource{
			{
				Id:      "resourceID",
				Type:    "Type",
				Actions: []string{"ACTION"},
			},
		}
		options := map[string]*authorizationpb.Option{}

		resp, err := client.WhoAuthorized(
			context.Background(),
			&authorizationpb.WhoAuthorizedRequest{
				Resources: resources,
				Options:   options,
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
	rootCmd.AddCommand(whoAuthorizedCmd)
}