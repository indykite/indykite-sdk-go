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
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/indykite/indykite-sdk-go/config"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

var entityMatchingPipelineConfigCmd = &cobra.Command{
	Use:   "entityMatchingPipeline",
	Short: "EntityMatchingPipeline config",
}

var createEntityMatchingPipelineConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create EntityMatchingPipeline config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EntityMatchingPipelineConfig{
			NodeFilter: &configpb.EntityMatchingPipelineConfig_NodeFilter{
				SourceNodeTypes: []string{"employee"},
				TargetNodeTypes: []string{"user"},
			},
			SimilarityScoreCutoff: 0.9,
			PropertyMappingStatus: configpb.EntityMatchingPipelineConfig_STATUS_SUCCESS,
			EntityMatchingStatus:  configpb.EntityMatchingPipelineConfig_STATUS_PENDING,
			PropertyMappings: []*configpb.EntityMatchingPipelineConfig_PropertyMapping{
				{
					SourceNodeType:        "employee",
					SourceNodeProperty:    "email",
					TargetNodeType:        "user",
					TargetNodeProperty:    "city",
					SimilarityScoreCutoff: 0.8,
				},
			},
			RerunInterval: "1 day",
			LastRunTime:   timestamppb.New(time.Now()),
			ReportUrl:     wrapperspb.String("gs://some-path"),
			ReportType:    wrapperspb.String("csv"),
		}
		createReq, _ := config.NewCreate("like-real-config-node-name3")
		createReq.ForLocation("gid:AAAAApkaja7LKUQot5UCGh6_Zc4")
		createReq.WithDisplayName("Like real ConfigNode Name3")
		createReq.WithEntityMatchingPipelineConfig(configuration)

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

var updateEntityMatchingPipelineConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update EntityMatchingPipeline config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EntityMatchingPipelineConfig{
			SimilarityScoreCutoff: 0.9,
			PropertyMappingStatus: configpb.EntityMatchingPipelineConfig_STATUS_SUCCESS,
			EntityMatchingStatus:  configpb.EntityMatchingPipelineConfig_STATUS_PENDING,
			PropertyMappings: []*configpb.EntityMatchingPipelineConfig_PropertyMapping{
				{
					SourceNodeType:        "employee",
					SourceNodeProperty:    "email",
					TargetNodeType:        "user",
					TargetNodeProperty:    "city",
					SimilarityScoreCutoff: 0.8,
				},
			},
			RerunInterval: "1 day",
			LastRunTime:   timestamppb.New(time.Now()),
			ReportUrl:     wrapperspb.String("gs://some-path"),
			ReportType:    wrapperspb.String("csv"),
		}
		updateReq, _ := config.NewUpdate("gid:AAAAIOxsvi0BPUb2iGsSM6J8M_Y")
		updateReq.WithEntityMatchingPipelineConfig(configuration)

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

var deleteEntityMatchingPipelineConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete EntityMatchingPipeline configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:AAAAIBkMrARucEOpqeM6vGJm5b0")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readEntityMatchingPipelineConfigCmd = &cobra.Command{
	Use:   "read",
	Short: "Read EntityMatchingPipeline by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var entityID string

		fmt.Print("Enter EntityMatchingPipeline ID in gid format: ")
		fmt.Scanln(&entityID)

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
	rootCmd.AddCommand(entityMatchingPipelineConfigCmd)
	entityMatchingPipelineConfigCmd.AddCommand(createEntityMatchingPipelineConfigCmd)
	entityMatchingPipelineConfigCmd.AddCommand(updateEntityMatchingPipelineConfigCmd)
	entityMatchingPipelineConfigCmd.AddCommand(deleteEntityMatchingPipelineConfigCmd)
	entityMatchingPipelineConfigCmd.AddCommand(readEntityMatchingPipelineConfigCmd)
}
