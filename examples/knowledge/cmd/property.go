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

	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
)

// propertyCmd represents the command for making a listNodesByProperty query with the Identity Knowledge API.
var propertyCmd = &cobra.Command{
	Use:   "property",
	Short: "Make a list query to the IndyKite Identity Knowledge API",
	Long:  `Make a list query to the IndyKite Identity Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		query := "MATCH (n:Resource)-[:HAS]->(p:Property) WHERE p.type=$type and p.value=$value"
		params := map[string]*objects.Value{
			"type": {
				Type: &objects.Value_StringValue{StringValue: "last_name"},
			},
			"value": {
				Type: &objects.Value_StringValue{StringValue: "mushu"},
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

// listPropertyCmd represents the command for making a listNodesByProperty query with the Identity Knowledge API.
var listPropertyCmd = &cobra.Command{
	Use:   "listproperty",
	Short: "Make a list query to the IndyKite Identity Knowledge API",
	Long:  `Make a list query to the IndyKite Identity Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := client.ListNodesByProperty(
			context.Background(),
			&knowledgeobjects.Property{
				Type: "last_name",
				Value: &objects.Value{
					Type: &objects.Value_StringValue{
						StringValue: "mushu",
					},
				},
			},
			true,
		)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(resp)
	},
}

func init() {
	rootCmd.AddCommand(propertyCmd)
	rootCmd.AddCommand(listPropertyCmd)
}
