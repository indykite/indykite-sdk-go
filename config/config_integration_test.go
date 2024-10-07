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

//go:build integration

package config_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/indykite/indykite-sdk-go/config"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Configuration", func() {
	Describe("ExternalDataResolver", func() {
		It("CreateExternalDataResolver", func() {
			var (
				err     error
				timeNow = fmt.Sprintf("%v", time.Now().UnixNano())
			)

			configClient, err := integration.InitConfigConfig()
			Expect(err).To(Succeed())

			displayNamePb := &wrapperspb.StringValue{Value: "AppSpace " + timeNow}
			createAppSpaceReq := &configpb.CreateApplicationSpaceRequest{
				CustomerId:  integration.CustomerID,
				Name:        "appspace-" + timeNow,
				DisplayName: displayNamePb,
				Region:      "europe-west1",
			}
			respAppSpace, err := configClient.CreateApplicationSpace(context.Background(), createAppSpaceReq)
			Expect(err).To(Succeed())
			Expect(respAppSpace).NotTo(BeNil())
			appSpaceID := respAppSpace.Id
			appSpaceEtag := respAppSpace.Etag

			configuration := &configpb.ExternalDataResolverConfig{
				Url:              integration.URL,
				Method:           integration.Method1,
				Headers:          integration.Headers,
				RequestType:      integration.RequestType,
				RequestPayload:   integration.RequestPayload,
				ResponseType:     integration.ResponseType,
				ResponseSelector: integration.ResponseSelector,
			}
			createReq, _ := config.NewCreate("resolver-" + timeNow)
			createReq.ForLocation(appSpaceID)
			createReq.WithDisplayName("Resolver" + timeNow)
			createReq.WithExternalDataResolverConfig(configuration)

			resp, err := configClient.CreateConfigNode(context.Background(), createReq)
			if err != nil {
				log.Fatalf("failed to invoke operation on IndyKite creation config node %v", err)
			}
			Expect(resp).NotTo(BeNil())
			configID := resp.Id
			configEtag := resp.Etag
			Expect(resp.LocationId).To(Equal(appSpaceID))

			readReq, _ := config.NewRead(configID)
			respRead, err := configClient.ReadConfigNode(context.Background(), readReq)
			Expect(err).To(Succeed())
			Expect(respRead).NotTo(BeNil())
			configNode := respRead.ConfigNode
			Expect(configNode).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":   Equal(configID),
				"Name": Equal("resolver-" + timeNow),
				"Config": PointTo(MatchFields(IgnoreExtras, Fields{
					"ExternalDataResolverConfig": integration.EqualProto(configuration),
				})),
			})))

			configurationUpd := &configpb.ExternalDataResolverConfig{
				Url:              integration.URLUpd,
				Method:           integration.Method1,
				Headers:          integration.HeadersUpd,
				RequestType:      integration.RequestType,
				RequestPayload:   integration.RequestPayload,
				ResponseType:     integration.ResponseType,
				ResponseSelector: integration.ResponseSelector,
			}
			updateReq, _ := config.NewUpdate(configID)
			updateReq.WithDisplayName("Resolver2" + timeNow)
			updateReq.WithExternalDataResolverConfig(configurationUpd)
			respUpd, err := configClient.UpdateConfigNode(context.Background(), updateReq)
			if err != nil {
				log.Fatalf("failed to invoke operation on IndyKite update config node Client %v", err)
			}
			Expect(respUpd).NotTo(BeNil())
			configUpdEtag := respUpd.Etag
			Expect(respUpd.Id).To(Equal(configID))
			Expect(respUpd.LocationId).To(Equal(appSpaceID))
			Expect(configUpdEtag).NotTo(Equal(configEtag))

			deleteReq, _ := config.NewDelete(configID)
			respDel, err := configClient.DeleteConfigNode(context.Background(), deleteReq)
			Expect(err).To(Succeed())
			Expect(respDel).NotTo(BeNil())

			time.Sleep(5 * time.Second)
			etagPb := &wrapperspb.StringValue{Value: appSpaceEtag}
			reqDelAS := &configpb.DeleteApplicationSpaceRequest{
				Id:        appSpaceID,
				Etag:      etagPb,
				Bookmarks: []string{},
			}
			respDelAS, err := configClient.DeleteApplicationSpace(context.Background(), reqDelAS)
			Expect(err).To(Succeed())
			Expect(respDelAS).NotTo(BeNil())
		})

		It("CreateExternalDataResolverErrorLocation", func() {
			var (
				err     error
				timeNow = fmt.Sprintf("%v", time.Now().UnixNano())
			)

			configClient, err := integration.InitConfigConfig()
			Expect(err).To(Succeed())

			configuration := &configpb.ExternalDataResolverConfig{
				Url:              integration.URL,
				Method:           integration.Method1,
				Headers:          integration.Headers,
				RequestType:      integration.RequestType,
				RequestPayload:   integration.RequestPayload,
				ResponseType:     integration.ResponseType,
				ResponseSelector: integration.ResponseSelector,
			}
			createReq, _ := config.NewCreate("resolver-" + timeNow)
			createReq.ForLocation(integration.WrongAppSpace)
			createReq.WithDisplayName("Resolver" + timeNow)
			createReq.WithExternalDataResolverConfig(configuration)

			resp, err := configClient.CreateConfigNode(context.Background(), createReq)
			Expect(err).To(MatchError(ContainSubstring(
				"insufficient permission to perform requested action")))
			Expect(resp).To(BeNil())
		})

		It("CreateExternalDataResolverWrongMethod", func() {
			var (
				err     error
				timeNow = fmt.Sprintf("%v", time.Now().UnixNano())
			)

			configClient, err := integration.InitConfigConfig()
			Expect(err).To(Succeed())

			displayNamePb := &wrapperspb.StringValue{Value: "AppSpace " + timeNow}
			createAppSpaceReq := &configpb.CreateApplicationSpaceRequest{
				CustomerId:  integration.CustomerID,
				Name:        "appspace-" + timeNow,
				DisplayName: displayNamePb,
				Region:      "europe-west1",
			}
			respAppSpace, err := configClient.CreateApplicationSpace(context.Background(), createAppSpaceReq)
			Expect(err).To(Succeed())
			Expect(respAppSpace).NotTo(BeNil())
			appSpaceID := respAppSpace.Id
			appSpaceEtag := respAppSpace.Etag

			configuration := &configpb.ExternalDataResolverConfig{
				Url:              integration.URL,
				Method:           integration.Method3,
				RequestType:      integration.RequestType,
				RequestPayload:   integration.RequestPayload,
				ResponseType:     integration.ResponseType,
				ResponseSelector: integration.ResponseSelector,
			}
			createReq, _ := config.NewCreate("resolver-" + timeNow)
			createReq.ForLocation(appSpaceID)
			createReq.WithDisplayName("Resolver" + timeNow)
			createReq.WithExternalDataResolverConfig(configuration)

			resp, err := configClient.CreateConfigNode(context.Background(), createReq)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid ExternalDataResolverConfig.Method: value must be in list [GET POST PUT PATCH]")))
			Expect(resp).To(BeNil())

			time.Sleep(5 * time.Second)
			etagPb := &wrapperspb.StringValue{Value: appSpaceEtag}
			reqDelAS := &configpb.DeleteApplicationSpaceRequest{
				Id:        appSpaceID,
				Etag:      etagPb,
				Bookmarks: []string{},
			}
			respDelAS, err := configClient.DeleteApplicationSpace(context.Background(), reqDelAS)
			Expect(err).To(Succeed())
			Expect(respDelAS).NotTo(BeNil())
		})
	})
})
