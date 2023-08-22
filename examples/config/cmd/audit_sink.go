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

var auditSinkCmd = &cobra.Command{
	Use:   "audit_sink",
	Short: "AuditSink config operations",
}

var createAuditSinkCmd = &cobra.Command{
	Use:   "create",
	Short: "Create AuditSink configuration",
	Run: func(cmd *cobra.Command, args []string) {
		createReq, _ := config.NewCreate("audit-sink-test")
		createReq.ForLocation("gid:fill-your-gid")
		createReq.WithAuditSinkConfig(&configpb.AuditSinkConfig{
			Provider: &configpb.AuditSinkConfig_Kafka{Kafka: &configpb.KafkaSinkConfig{
				Brokers:  []string{"broker.com"},
				Topic:    "your-topic-name",
				Username: "your-name",
				Password: "your-password",
			}},
		})

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

var updateAuditSinkCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AuditSink configuration",
	Run: func(cmd *cobra.Command, args []string) {
		updateReq, _ := config.NewUpdate("gid:id-of-existing-config")
		updateReq.WithAuditSinkConfig(&configpb.AuditSinkConfig{
			Provider: &configpb.AuditSinkConfig_Kafka{Kafka: &configpb.KafkaSinkConfig{
				Brokers:  []string{"broker.com"},
				Topic:    "your-topic-name",
				Username: "your-name",
				Password: "your-password",
			}},
		})

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

var deleteAuditSinkCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete AuditSink configuration",
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
	rootCmd.AddCommand(auditSinkCmd)
	auditSinkCmd.AddCommand(createAuditSinkCmd)
	auditSinkCmd.AddCommand(updateAuditSinkCmd)
	auditSinkCmd.AddCommand(deleteAuditSinkCmd)
}
