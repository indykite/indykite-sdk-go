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

var externalDataResolverConfigCmd = &cobra.Command{
	Use:   "externalDataResolver",
	Short: "ExternalDataResolver config",
}

var createExternalDataResolverConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ExternalDataResolver config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.ExternalDataResolverConfig{
			Url:              "http://super-octo-waffle.indykite.com/magic?data=2024",
			Method:           "GET",
			Headers:          map[string]*configpb.ExternalDataResolverConfig_Header{},
			RequestType:      configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
			RequestPayload:   []byte(``),
			ResponseType:     configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
			ResponseSelector: ".echo",
		}
		createReq, _ := config.NewCreate("like-real-config-node-name2")
		createReq.ForLocation("gid:AAAAAguAAAAAAAAAAAAAAAAAAAA")
		createReq.WithDisplayName("Like real ConfigNode Name2")
		createReq.WithExternalDataResolverConfig(configuration)

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

var updateExternalDataResolverConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ExternalDataResolver config",
	Run: func(cmd *cobra.Command, args []string) {
		configuration := &configpb.ExternalDataResolverConfig{
			Url:              "http://super-octo-waffle.indykite.com/magic?data=2024",
			Method:           "GET",
			Headers:          map[string]*configpb.ExternalDataResolverConfig_Header{},
			RequestType:      configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
			RequestPayload:   []byte(``),
			ResponseType:     configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
			ResponseSelector: ".echo",
		}
		updateReq, _ := config.NewUpdate("gid:AAAAIZISzhPyS0i-hT-OnuiGkKE")
		updateReq.WithDescription("Desc2")
		updateReq.WithExternalDataResolverConfig(configuration)

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

var deleteExternalDataResolverConfigCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete ExternalDataResolver configuration",
	Run: func(cmd *cobra.Command, args []string) {
		deleteReq, _ := config.NewDelete("gid:AAAAIWWRI5nLaEkjhFBa4v8Gi-4")
		resp, err := client.DeleteConfigNode(context.Background(), deleteReq)
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println(jsonp.Format(resp))
	},
}

func init() {
	rootCmd.AddCommand(externalDataResolverConfigCmd)
	externalDataResolverConfigCmd.AddCommand(createExternalDataResolverConfigCmd)
	externalDataResolverConfigCmd.AddCommand(updateExternalDataResolverConfigCmd)
	externalDataResolverConfigCmd.AddCommand(deleteExternalDataResolverConfigCmd)
}
