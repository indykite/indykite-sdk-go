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
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Person"},
		"actions":["SUBSCRIBES_TO"],"resource":{"type":"Truck"},
		"condition":{"cypher":"MATCH (subject:Person)-[:BELONGS_TO]->(:Organization)-[:OWNS]
		->(resource:Truck)-[HAS]
		->(p:Property:External {type: 'echo', value: '2024'}) "}}`
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

		readReq, _ := config.NewRead(resp.GetId())
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
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Person"},
		"actions":["SUBSCRIBES_TO"],"resource":{"type":"Asset"},
		"condition":{"cypher":"MATCH (subject:Person)-[:BELONGS_TO]->(:Organization)-[:OWNS]
		->(resource:Truck)-[HAS]->(Truck:Property:External {type: echo, value: '2024'}) "}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{"TagA", "TagB"},
		}
		updateReq, _ := config.NewUpdate("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		updateReq.WithAuthorizationPolicyConfig(configuration)
		updateReq.WithDescription("Desc1")

		resp, err := client.UpdateConfigNode(context.Background(), updateReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.GetId())
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
		deleteReq, _ := config.NewDelete("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

// policy for trust score.
var createAuthorizationPolicyConfig2Cmd = &cobra.Command{
	Use:   "create2",
	Short: "Create AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Agent"},
		"actions":["CAN_USE"],"resource":{"type":"Sensor"},
		"condition":{"cypher":"MATCH (:_TrustScore)<-[:_HAS]-(subject)-[:CAN_USE]->(resource)"}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		createReq, _ := config.NewCreate("like-real-config-node-score")
		createReq.ForLocation("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		createReq.WithDisplayName("Like real ConfigNode Name")
		createReq.WithAuthorizationPolicyConfig(configuration)

		resp, err := client.CreateConfigNode(context.Background(), createReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.GetId())
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

// policy for knowledge query.
var createAuthorizationPolicyConfig3Cmd = &cobra.Command{
	Use:   "create3",
	Short: "Create AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policy_version":"1.0-ciq"},
		"subject":{"type":"Person"},
		"condition":{
		"cypher":"MATCH (person:Person)-[r1:BELONGS_TO]->(org:Organization)-[r2:OWNS]->(resource:Truck) ",
		"filter" : [{ "app" : "app-sdk", "attribute" : "person.property.last_name", 
		"operator" : "=", "value" : "$lastname" }]
		},
		"upsert_nodes" : [],
	    "upsert_relationships" : [],
	    "allowed_reads" : {
		  "nodes" : ["resource"],
		"relationships" : ["r2"]
	  }
		}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		createReq, _ := config.NewCreate("like-real-config-node-ciq")
		createReq.ForLocation("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		createReq.WithDisplayName("Like real ConfigNode Name")
		createReq.WithAuthorizationPolicyConfig(configuration)

		resp, err := client.CreateConfigNode(context.Background(), createReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.GetId())
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
		jsonInput := `{"meta":{"policyVersion":"1.0-indykite"},"subject":{"type":"Agent"},
		"actions":["CAN_USE"],"resource":{"type":"Sensor"},
		"condition":{"cypher":"MATCH (:_TrustScore)<-[:_HAS]-(subject)-[:CAN_USE]->(resource)"}}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		updateReq, _ := config.NewUpdate("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		updateReq.WithAuthorizationPolicyConfig(configuration)
		updateReq.WithDescription("Desc1")

		resp, err := client.UpdateConfigNode(context.Background(), updateReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.GetId())
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

var updateAuthorizationPolicyConfig3Cmd = &cobra.Command{
	Use:   "update3",
	Short: "Update AuthorizationPolicy config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonInput := `{"meta":{"policy_version":"1.0-ciq"},
		"subject":{"type":"Person"},
		"condition":{
		"cypher":"MATCH (person:Person)-[r1:BELONGS_TO]->(org:Organization)-[r2:OWNS]->(resource:Truck) ",
		"filter" : [{ "app" : "app-sdk", "attribute" : "person.property.last_name", 
		"operator" : "=", "value" : "$lastname" }]
		},
		"upsert_nodes" : [],
	    "upsert_relationships" : [],
	    "allowed_reads" : {
		  "nodes" : ["resource.property.value", "resource.property.transferrable", "resource.external_id"],
		"relationships" : ["r2"]
	  }
		}`
		configuration := &configpb.AuthorizationPolicyConfig{
			Policy: jsonInput,
			Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
			Tags:   []string{},
		}
		updateReq, _ := config.NewUpdate("gid:AAAAAguDnAAAAAAAAAAAAAAA")
		updateReq.WithAuthorizationPolicyConfig(configuration)
		updateReq.WithDescription("Desc1")

		resp, err := client.UpdateConfigNode(context.Background(), updateReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))

		readReq, _ := config.NewRead(resp.GetId())
		readResp, err := client.ReadConfigNode(context.Background(), readReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(readResp))
	},
}

var readAuthorizationPolicyConfigCmd = &cobra.Command{
	Use:   "read",
	Short: "Read Authorization Policy by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var entityID string

		fmt.Print("Enter Authorization Policy ID in gid format: ")
		if _, err := fmt.Scanln(&entityID); err != nil {
			fmt.Println("Error reading entityID:", err)
		}

		configNodeRequest, err := config.NewRead(entityID)
		if err != nil {
			log.Fatalf("failed to create request %v", err)
		}
		resp, err := client.ReadConfigNode(context.Background(), configNodeRequest)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(authorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(createAuthorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(updateAuthorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(deleteAuthorizationPolicyConfigCmd)
	authorizationPolicyConfigCmd.AddCommand(createAuthorizationPolicyConfig2Cmd)
	authorizationPolicyConfigCmd.AddCommand(updateAuthorizationPolicyConfig2Cmd)
	authorizationPolicyConfigCmd.AddCommand(createAuthorizationPolicyConfig3Cmd)
	authorizationPolicyConfigCmd.AddCommand(updateAuthorizationPolicyConfig3Cmd)
	authorizationPolicyConfigCmd.AddCommand(readAuthorizationPolicyConfigCmd)
}
