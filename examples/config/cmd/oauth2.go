// Copyright (c) 2022 IndyKite
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
	"google.golang.org/protobuf/types/known/wrapperspb"

	config "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"
)

var oauth2Cmd = &cobra.Command{
	Use:   "oauth2",
	Short: "...",
	Run: func(cmd *cobra.Command, args []string) {

		providerCfg := &config.OAuth2ProviderConfig{
			GrantTypes:    []config.GrantType{config.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
			ResponseTypes: []config.ResponseType{config.ResponseType_RESPONSE_TYPE_CODE},
			Scopes:        []string{"openid", "profile", "email"},
			TokenEndpointAuthMethod: []config.TokenEndpointAuthMethod{
				config.TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_BASIC,
			},
			TokenEndpointAuthSigningAlg: []string{"ES256"},
			RequestUris:                 make([]string, 0),
			RequestObjectSigningAlg:     "ES256",
			FrontChannelLoginUri: map[string]string{
				"default": "http://www.google.com",
			},
			FrontChannelConsentUri: map[string]string{
				"default": "http://www.google.com",
			},
		}

		// /////////////////
		// PROVIDER part //
		// /////////////////

		// Create Provider
		providerResp, err := client.CreateOAuth2Provider(context.Background(), &config.CreateOAuth2ProviderRequest{
			AppSpaceId: "gid:AAAAAmluZHlraURlgAABDwAAAAA",
			Name:       "test-sdk-3",
			Config:     providerCfg,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Create", jsonp.Format(providerResp))

		// Read Provider
		rr, err := client.ReadOAuth2Provider(context.Background(), &config.ReadOAuth2ProviderRequest{
			Id: providerResp.Id,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Read", jsonp.Format(rr))

		// Update Provider
		ur, err := client.UpdateOAuth2Provider(context.Background(), &config.UpdateOAuth2ProviderRequest{
			Id:          providerResp.Id,
			DisplayName: wrapperspb.String("New display name"),
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Update", jsonp.Format(ur))

		// Read Provider
		rr, err = client.ReadOAuth2Provider(context.Background(), &config.ReadOAuth2ProviderRequest{
			Id: providerResp.Id,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Read", jsonp.Format(rr))

		// Update Provider
		providerCfg.RequestObjectSigningAlg = "ES512"
		ur, err = client.UpdateOAuth2Provider(context.Background(), &config.UpdateOAuth2ProviderRequest{
			Id:     providerResp.Id,
			Config: providerCfg,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Update", jsonp.Format(ur))

		// Read Provider
		rr, err = client.ReadOAuth2Provider(context.Background(), &config.ReadOAuth2ProviderRequest{
			Id: providerResp.Id,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Read", jsonp.Format(rr))

		// ////////////////////
		// APPLICATION part //
		// ////////////////////

		appCfg := &config.OAuth2ApplicationConfig{
			ClientId:                    "00000000-90af-4ef9-9928-aaaaaaaaaaaa",
			DisplayName:                 "Some cool public display name",
			SubjectType:                 config.ClientSubjectType_CLIENT_SUBJECT_TYPE_PAIRWISE,
			GrantTypes:                  []config.GrantType{config.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
			ResponseTypes:               []config.ResponseType{config.ResponseType_RESPONSE_TYPE_CODE},
			Scopes:                      []string{"openid", "profile", "email"},
			Audiences:                   []string{"7d2e906e-541a-49da-b5b2-a28840ff8721"},
			TokenEndpointAuthMethod:     config.TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_BASIC,
			TokenEndpointAuthSigningAlg: "ES256",
		}

		// Create Application
		appResp, err := client.CreateOAuth2Application(context.Background(), &config.CreateOAuth2ApplicationRequest{
			Oauth2ProviderId: providerResp.Id,
			Name:             "test-sdk-3",
			Config:           appCfg,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Create", jsonp.Format(appResp))

		// Read Application
		ra, err := client.ReadOAuth2Application(context.Background(), &config.ReadOAuth2ApplicationRequest{
			Id: appResp.Id,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Read", jsonp.Format(ra))

		// Update Application
		ua, err := client.UpdateOAuth2Application(context.Background(), &config.UpdateOAuth2ApplicationRequest{
			Id:          appResp.Id,
			DisplayName: wrapperspb.String("New display name"),
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Update", jsonp.Format(ua))

		// Read Application
		ra, err = client.ReadOAuth2Application(context.Background(), &config.ReadOAuth2ApplicationRequest{
			Id: appResp.Id,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Read", jsonp.Format(ra))
		// Update Application
		appCfg.LogoUri = "http://localhost/logo.png"
		ua, err = client.UpdateOAuth2Application(context.Background(), &config.UpdateOAuth2ApplicationRequest{
			Id:     appResp.Id,
			Config: appCfg,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Update", jsonp.Format(ua))

		// Read Application
		ra, err = client.ReadOAuth2Application(context.Background(), &config.ReadOAuth2ApplicationRequest{
			Id: appResp.Id,
		})
		if err != nil {
			log.Fatalf("failed to invoke operation on IndyKite Client %v", err)
		}
		fmt.Println("Read", jsonp.Format(ra))
	},
}

func init() {
	rootCmd.AddCommand(oauth2Cmd)
}
