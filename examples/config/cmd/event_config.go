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

var eventSinkConfigCmd = &cobra.Command{
	Use:   "eventSink",
	Short: "EventSink config",
}

var createEventSinkConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create EventSink config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EventSinkConfig{
			Providers: map[string]*configpb.EventSinkConfig_Provider{
				"kafka": {
					Provider: &configpb.EventSinkConfig_Provider_Kafka{
						Kafka: &configpb.KafkaSinkConfig{
							Brokers:  []string{"broker.com"},
							Topic:    "your-topic-name",
							Username: "your-name",
							Password: "your-password",
						},
					},
				},
			},
			Routes: []*configpb.EventSinkConfig_Route{
				{
					ProviderId:     "kafka",
					StopProcessing: true,
					Filter: &configpb.EventSinkConfig_Route_EventType{
						EventType: "indykite.eventsink.config.create",
					},
				},
			},
		}
		createReq, _ := config.NewCreate("like-real-config-node-name-ts3")
		createReq.ForLocation("gid:AAAAAguDnEIQ1UIIvJEulLwUnnE")
		createReq.WithDisplayName("Like real ConfigNode Name TS3")
		createReq.WithEventSinkConfig(configuration)

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

var updateEventSinkConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update EventSink  config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EventSinkConfig{
			Providers: map[string]*configpb.EventSinkConfig_Provider{
				"kafka": {
					Provider: &configpb.EventSinkConfig_Provider_Kafka{
						Kafka: &configpb.KafkaSinkConfig{
							Brokers:  []string{"broker.com"},
							Topic:    "your-topic-name",
							Username: "your-name-update",
							Password: "your-password-update",
						},
					},
				},
			},
			Routes: []*configpb.EventSinkConfig_Route{
				{
					ProviderId:     "kafka",
					StopProcessing: true,
					Filter: &configpb.EventSinkConfig_Route_EventType{
						EventType: "indykite.eventsink.config.update",
					},
				},
			},
		}
		updateReq, _ := config.NewUpdate("gid:AAAAG3FQqyfhzEiUrpVHvab4ct4")
		updateReq.WithEventSinkConfig(configuration)

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

var createEventSinkConfigGridCmd = &cobra.Command{
	Use:   "create2",
	Short: "Create EventSink config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EventSinkConfig{
			Providers: map[string]*configpb.EventSinkConfig_Provider{
				"azureeventgrid": {
					Provider: &configpb.EventSinkConfig_Provider_AzureEventGrid{
						AzureEventGrid: &configpb.AzureEventGridSinkConfig{
							TopicEndpoint: "https://ik-test.eventgrid.azure.net/api/events",
							AccessKey:     "your-access-key",
						},
					},
				},
			},
			Routes: []*configpb.EventSinkConfig_Route{
				{
					ProviderId:     "azureeventgrid",
					StopProcessing: true,
					Filter: &configpb.EventSinkConfig_Route_EventType{
						EventType: "indykite.eventsink.config.create",
					},
				},
			},
		}
		createReq, _ := config.NewCreate("like-real-config-node-name-ts3")
		createReq.ForLocation("gid:AAAAAgZQ3QPgJ0gAkbNQ-IjGvXQ")
		createReq.WithDisplayName("Like real ConfigNode Name TS3")
		createReq.WithEventSinkConfig(configuration)

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

var createEventSinkConfigBusCmd = &cobra.Command{
	Use:   "create3",
	Short: "Create EventSink config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.EventSinkConfig{
			Providers: map[string]*configpb.EventSinkConfig_Provider{
				"azureservicebus": {
					Provider: &configpb.EventSinkConfig_Provider_AzureServiceBus{
						AzureServiceBus: &configpb.AzureServiceBusSinkConfig{
							ConnectionString: "Endpoint=sb://ik-test.servicebus.windows.net/;SharedAccessKeyName=Root",
							QueueOrTopicName: "your-queue",
						},
					},
				},
			},
			Routes: []*configpb.EventSinkConfig_Route{
				{
					ProviderId:     "azureservicebus",
					StopProcessing: true,
					Filter: &configpb.EventSinkConfig_Route_EventType{
						EventType: "indykite.eventsink.config.create",
					},
				},
			},
		}
		createReq, _ := config.NewCreate("like-real-config-node-name-ts3")
		createReq.ForLocation("gid:AAAAAgZQ3QPgJ0gAkbNQ-IjGvXQ")
		createReq.WithDisplayName("Like real ConfigNode Name TS3")
		createReq.WithEventSinkConfig(configuration)

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

var deleteEventSinkConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete EventSink configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:AAAAG7ub-V69fEE0kJZqHcpb1I0")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

var readEventSinkConfigCmd = &cobra.Command{
	Use:   "read",
	Short: "Read EventSink by ID",
	Run: func(cmd *cobra.Command, args []string) {
		var entityID string

		fmt.Print("Enter EventSink ID in gid format: ")
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
	rootCmd.AddCommand(eventSinkConfigCmd)
	eventSinkConfigCmd.AddCommand(createEventSinkConfigCmd)
	eventSinkConfigCmd.AddCommand(updateEventSinkConfigCmd)
	eventSinkConfigCmd.AddCommand(createEventSinkConfigGridCmd)
	eventSinkConfigCmd.AddCommand(createEventSinkConfigBusCmd)
	eventSinkConfigCmd.AddCommand(deleteEventSinkConfigCmd)
	eventSinkConfigCmd.AddCommand(readEventSinkConfigCmd)
}
