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

	"github.com/indykite/indykite-sdk-go/config"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

var consentConfigCmd = &cobra.Command{
	Use:   "consent",
	Short: "Consent config",
}

var createConsentConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Consent configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.ConsentConfiguration{
			Purpose: "Taking control again",
			DataPoints: []string{
				"{ \"query\": \"\", \"returns\": [ { \"variable\": \"\"," +
					"\"properties\": [\"name\", \"email\", \"location\"] } ] }",
			},
			ApplicationId:  "gid:AAAABMoo7PXYfkwepSVjj4GTtfc",
			ValidityPeriod: 86400,
			RevokeAfterUse: true,
			TokenStatus:    3,
		}
		createReq, _ := config.NewCreate("like-real-config-node-name")
		createReq.ForLocation("gid:AAAAAvFyVpD_1kd8k2kpNY9rjFM")
		createReq.WithDisplayName("Like real ConfigNode Name")
		createReq.WithConsentConfig(configuration)

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

var updateConsentConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Consent configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.ConsentConfiguration{
			Purpose: "Taking control upd",
			DataPoints: []string{"\"query\": \"\", \"returns\": [ { \"variable\": \"\", " +
				"\"properties\": [\"name\", \"email\", \"location\"]}]"},
			ApplicationId:  "gid:like-real-application-id",
			ValidityPeriod: 86400,
			RevokeAfterUse: true,
			TokenStatus:    3,
		}
		updateReq, _ := config.NewUpdate("gid:id-of-existing-config")
		updateReq.WithConsentConfig(configuration)

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

var deleteConsentConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Consent configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:id-of-existing-config")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(consentConfigCmd)
	consentConfigCmd.AddCommand(createConsentConfigCmd)
	consentConfigCmd.AddCommand(updateConsentConfigCmd)
	consentConfigCmd.AddCommand(deleteConsentConfigCmd)
}
