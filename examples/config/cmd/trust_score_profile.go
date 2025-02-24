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

var trustScoreProfileConfigCmd = &cobra.Command{
	Use:   "trustScoreProfile",
	Short: "TrustScoreProfile config",
}

var createTrustScoreProfileConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create TrustScoreProfile config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.TrustScoreProfileConfig{
			NodeClassification: "Person",
			Dimensions: []*configpb.TrustScoreDimension{
				{
					Name:   configpb.TrustScoreDimension_NAME_FRESHNESS,
					Weight: 0.9,
				},
			},
			Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_DAILY,
		}
		createReq, _ := config.NewCreate("like-real-config-node-name-ts")
		createReq.ForLocation("gid:AAAAAiOpnZsjRkIVid6bET5RRE4")
		createReq.WithDisplayName("Like real ConfigNode Name TS")
		createReq.WithTrustScoreProfileConfig(configuration)

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

var updateTrustScoreProfileConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update TrustScoreProfile config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.TrustScoreProfileConfig{
			Dimensions: []*configpb.TrustScoreDimension{
				{
					Name:   configpb.TrustScoreDimension_NAME_COMPLETENESS,
					Weight: 0.92,
				},
			},
			Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_DAILY,
		}
		updateReq, _ := config.NewUpdate("gid:AAAAIsOrnjcZNUNrmuXxnWoJS6s")
		updateReq.WithTrustScoreProfileConfig(configuration)

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

var deleteTrustScoreProfileConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete TrustScoreProfile configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:AAAAIsOrnjcZNUNrmuXxnWoJS6s")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readTrustScoreProfileConfigCmd = &cobra.Command{
	Use:   "read",
	Short: "Read TrustScoreProfile by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var entityID string

		fmt.Print("Enter TrustScoreProfile ID in gid format: ")
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
	rootCmd.AddCommand(trustScoreProfileConfigCmd)
	trustScoreProfileConfigCmd.AddCommand(createTrustScoreProfileConfigCmd)
	trustScoreProfileConfigCmd.AddCommand(updateTrustScoreProfileConfigCmd)
	trustScoreProfileConfigCmd.AddCommand(deleteTrustScoreProfileConfigCmd)
	trustScoreProfileConfigCmd.AddCommand(readTrustScoreProfileConfigCmd)
}
