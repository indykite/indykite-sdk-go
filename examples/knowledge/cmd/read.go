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

		query := "MATCH (n:Person)-[:BELONGS_TO]->(o:Organization)-[:OWNS]->(t:Truck)-[:HAS]->(p:Property:External) WHERE n.external_id=$external_id AND n.type=$type"
		params := map[string]*objects.Value{
			"external_id": {
				Type: &objects.Value_StringValue{StringValue: "fVcaUxJqmOkyOTX"},
			},
			"type": {
				Type: &objects.Value_StringValue{StringValue: "Person"},
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

// readScoreCmd represents the command for making a read query with the Identity Knowledge API with trust score info
var readScoreCmd = &cobra.Command{
	Use:   "read-score",
	Short: "Make a read query with trust score to the IndyKite Identity Knowledge API",
	Long:  `Make a read query with trust score to the IndyKite Identity Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		query := "MATCH (a:Agent)-[:_HAS]->(t:_TrustScore)"
		params := map[string]*objects.Value{}
		returns := []*knowledgepb.Return{
			{
				Variable:   "a",
				Properties: []string{},
			},
			{
				Variable:   "t",
				Properties: []string{},
			},
		}

		resp, err := client.IdentityKnowledgeRead(context.Background(), query, params, returns)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// readScore2Cmd represents all trust score that match the Trust Score profile name
var readScore2Cmd = &cobra.Command{
	Use:   "read-score2",
	Short: "Make a read query with trust score to the IndyKite Identity Knowledge API",
	Long:  `Make a read query with trust score to the IndyKite Identity Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		query := "MATCH(t:_TrustScore {_ingress: $profile})"
		params := map[string]*objects.Value{
			"profile": {
				Type: &objects.Value_StringValue{StringValue: "like-real-config-node-name-ts3"},
			},
		}
		returns := []*knowledgepb.Return{
			{
				Variable: "t",
			},
		}

		resp, err := client.IdentityKnowledgeRead(context.Background(), query, params, returns)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// readScore3Cmd represents specific dimensions from all Trust Scores
var readScore3Cmd = &cobra.Command{
	Use:   "read-score3",
	Short: "Make a read query with trust score to the IndyKite Identity Knowledge API",
	Long:  `Make a read query with trust score to the IndyKite Identity Knowledge API`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		query := "MATCH (t:_TrustScore)"
		params := map[string]*objects.Value{}
		returns := []*knowledgepb.Return{
			{
				Variable:   "t",
				Properties: []string{"_origin", "_verification"},
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
	rootCmd.AddCommand(readScoreCmd)
	rootCmd.AddCommand(readScore2Cmd)
	rootCmd.AddCommand(readScore3Cmd)
}
