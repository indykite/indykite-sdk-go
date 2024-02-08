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

	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
)

// readCmd represents the command for making a read query with the Identity Knowledge API
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Make a read query to the IndyKite Identity Knowledge API",
	Long:  `Make a read query to the IndyKite Identity Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		query := "MATCH (n:Person)-[r:CAN_SEE]->(a:Asset) WHERE n.external_id=$external_id AND n.type=$type"
		params := map[string]*objects.Value{
			"external_id": {
				Type: &objects.Value_StringValue{StringValue: "1234"},
			},
			"type": {
				Type: &objects.Value_StringValue{StringValue: "store"},
			},
		}
		returns := []*knowledgepb.Return{
			{
				Variable: "n",
			},
		}

		resp, err := client.IdentityKnowledgeRead(context.Background(), query, params, returns)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
