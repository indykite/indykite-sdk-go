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
)

// introspectCmd represents the plan command
var introspectCmd = &cobra.Command{
	Use:   "introspect",
	Short: "DigitalTwin operation",
	Long: `General commands for DigitalTwin

  This is a sample only.`,
	Run: func(cmd *cobra.Command, args []string) {
		var accessToken string
		for {
			fmt.Print("Enter access_token: ")
			fmt.Scanln(&accessToken)
			tenant, err := client.IntrospectToken(context.Background(), accessToken, retry.WithMax(2))
			if err != nil {
				log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
			}
			fmt.Println(jsonp.Format(tenant))
		}
	},
}

func init() {
	rootCmd.AddCommand(introspectCmd)
}
