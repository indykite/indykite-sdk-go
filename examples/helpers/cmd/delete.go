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

	"github.com/spf13/cobra"
)

// deleteCmd represents the command for deleting the nodes in a KG
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete nodes in the IndyKite Knowledge API",
	Long:  `Delete nodes in the IndyKite Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		responses, err := client.DeleteNodesWithTypeNode(context.Background(), "DigitalTwin")
		if err != nil {
			fmt.Println("DigitalTwin: " + err.Error())
		}
		responses2, err := client.DeleteNodesWithTypeNode(context.Background(), "Resource")
		if err != nil {
			fmt.Println("Resource: " + err.Error())
		}
		if len(responses2) > 0 {
			responses = append(responses, responses2...)
		}
		for _, response := range responses {
			fmt.Println(jsonp.Format(response))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
