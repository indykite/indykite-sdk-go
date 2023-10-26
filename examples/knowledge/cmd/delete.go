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
		records, err := client.NodesRecordsWithTypeNode(context.Background(), "DigitalTwin")
		if err != nil {
			fmt.Println(err.Error())
		}
		records2, err := client.NodesRecordsWithTypeNode(context.Background(), "Resource")
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(records2) > 0 {
			records = append(records, records2...)
		}
		fmt.Println(records)
		resp, err := clientIngest.StreamRecords(records)
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, response := range resp {
			fmt.Println(jsonp.Format(response))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
