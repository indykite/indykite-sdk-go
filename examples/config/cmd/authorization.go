// Copyright (c) 2024 IndyKite
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

	"github.com/indykite/indykite-sdk-go/config"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

var authorizationPolicyConfigCmd = &cobra.Command{
	Use:   "authorizationPolicy",
	Short: "AuthorizationPolicy config",
}

var createAuthorizationPolicyConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Person"},"actions":["SUBSCRIBES_TO"],"resource":{"type":"Truck"},"condition":{"cypher":"MATCH (subject:Person)-[:BELONGS_TO]->(:Organization)-[:OWNS]->(resource:Truck)-[HAS]->(p:Property:External {type: 'echo', value: '2024'}) "}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		createReq, _ := config.NewCreate("like-real-config-node-name")
		createReq.ForLocation("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		createReq.WithDisplayName("Like real ConfigNode Name")
		createReq.WithAuthorizationPolicyConfig(configuration)

		resp, err := client.CreateConfigNode(context.Background(), createReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.Id)
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

var updateAuthorizationPolicyConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Person"},"actions":["SUBSCRIBES_TO"],"resource":{"type":"Asset"},"condition":{"cypher":"MATCH (subject:Person)-[:BELONGS_TO]->(:Organization)-[:OWNS]->(resource:Truck)-[HAS]->(Truck:Property:External {type: echo, value: '2024'}) "}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{"TagA", "TagB"},
		}
		updateReq, _ := config.NewUpdate("gid:AAAAFo7ukfFQHkBjtiQQZiE2zb8")
		updateReq.WithAuthorizationPolicyConfig(configuration)
		updateReq.WithDescription("Desc1")

		resp, err := client.UpdateConfigNode(context.Background(), updateReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.Id)
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

var deleteAuthorizationPolicyConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete AuthorizationPolicy configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:AAAAFvTeAqwrRUinglaK7B891aI")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// policy for trust score
var createAuthorizationPolicyConfig2Cmd = &cobra.Command{
	Use:   "create2",
	Short: "Create AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Agent"},"actions":["CAN_USE"],"resource":{"type":"Sensor"},"condition":{"cypher":"MATCH (:_TrustScore)<-[:_HAS]-(subject)-[:CAN_USE]->(resource)"}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		createReq, _ := config.NewCreate("like-real-config-node-score")
		createReq.ForLocation("gid:AAAAAguDnEIQ1UIIvJEulLwUnnE")
		createReq.WithDisplayName("Like real ConfigNode Name")
		createReq.WithAuthorizationPolicyConfig(configuration)

		resp, err := client.CreateConfigNode(context.Background(), createReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.Id)
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

var updateAuthorizationPolicyConfig2Cmd = &cobra.Command{
	Use:   "update2",
	Short: "Update AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Agent"},"actions":["CAN_USE"],"resource":{"type":"Sensor"},"condition":{"cypher":"MATCH (:_TrustScore)<-[:_HAS]-(subject)-[:CAN_USE]->(resource)"}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		updateReq, _ := config.NewUpdate("gid:AAAAFlj3-0Ixw0zmo_c83L5R60k")
		updateReq.WithAuthorizationPolicyConfig(configuration)
		updateReq.WithDescription("Desc1")

		resp, err := client.UpdateConfigNode(context.Background(), updateReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.Id)
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

func init() {
	rootCmd.AddCommand(authorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(createAuthorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(updateAuthorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(deleteAuthorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(createAuthorizationPolicyConfig2Cmd)
	authorizationPolicyConfigCmd.AddCommand(updateAuthorizationPolicyConfig2Cmd)
}
