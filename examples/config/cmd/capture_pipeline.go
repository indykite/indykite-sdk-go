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

var capturePipelineConfigCmd = &cobra.Command{
	Use:   "capture-pipeline",
	Short: "Capture pipeline config",
}

var registerCapturePipelineConfigCmd = &cobra.Command{
	Use:   "register-pipeline",
	Short: "Register Capture Pipeline config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.RegisterCapturePipelineConfig{
			AppAgentToken: "eyJhbGciOiJIIkpXVCJ9.eyJpc3MiOiJnaWQkwMjJ9.39Kc7pL8Vjf1S4Oo5Rw",
		}
		createReq, _ := config.NewCreate("like-real-config-node-name-cp")
		createReq.ForLocation("gid:AAAAAguDnEIQ1UIIvJEulLwUnnE")
		createReq.WithDisplayName("Like real ConfigNode Name CP")
		createReq.WithCapturePipelineConfig(configuration)

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

var registerCapturePipelineTopicConfigCmd = &cobra.Command{
	Use:   "register-topic",
	Short: "Register Capture Pipeline Topic config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.RegisterCapturePipelineTopicConfig{
			CapturePipelineId: "gid:AAAAG3FQqyfhzEiUrpVHvab4ct4",
			Script: &configpb.RegisterCapturePipelineTopicConfig_Script{
				Content: "content of the script",
			},
		}
		createReq, _ := config.NewCreate("like-real-config-node-name-cp")
		createReq.ForLocation("gid:AAAAAguDnEIQ1UIIvJEulLwUnnE")
		createReq.WithDisplayName("Like real ConfigNode Name CP")
		createReq.WithCapturePipelineTopicConfig(configuration)

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

var deleteCapturePipelineConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Capture Pipeline configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:AAAAG3FQqyfhzEiUrpVHvab4ct4")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readCapturePipelineConfigCmd = &cobra.Command{
	Use:   "read-pipeline",
	Short: "Read Capture pipeline by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var entityID string

		fmt.Print("Enter Capture Pipeline ID in gid format: ")
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

var readCapturePipelineTopicConfigCmd = &cobra.Command{
	Use:   "read-topic",
	Short: "Read Capture pipeline topic by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var entityID string

		fmt.Print("Enter Capture Pipeline Topic ID in gid format: ")
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
	rootCmd.AddCommand(capturePipelineConfigCmd)
	capturePipelineConfigCmd.AddCommand(registerCapturePipelineConfigCmd)
	capturePipelineConfigCmd.AddCommand(registerCapturePipelineTopicConfigCmd)
	capturePipelineConfigCmd.AddCommand(updateEventSinkConfigCmd)
	capturePipelineConfigCmd.AddCommand(deleteCapturePipelineConfigCmd)
	capturePipelineConfigCmd.AddCommand(readCapturePipelineConfigCmd)
	capturePipelineConfigCmd.AddCommand(readCapturePipelineTopicConfigCmd)
}
