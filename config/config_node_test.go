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

package config_test

import (
	"context"
	"errors"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/indykite/indykite-sdk-go/config"
	sdkerrors "github.com/indykite/indykite-sdk-go/errors"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	configmock "github.com/indykite/indykite-sdk-go/test/config/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ConfigNode", func() {
	var (
		ctx          context.Context
		mockCtrl     *gomock.Controller
		mockClient   *configmock.MockConfigManagementAPIClient
		configClient *config.Client
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = configmock.NewMockConfigManagementAPIClient(mockCtrl)

		var err error
		configClient, err = config.NewTestClient(ctx, mockClient)
		Ω(err).To(Succeed())
	})

	Describe("ConfigNode", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadConfigNode(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil or not read request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			configNodeRequest, err := config.NewRead("gid:like")
			Ω(err).To(Succeed())
			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("ReadSuccessAuthorizationPolicy", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			jsonInput := `{
				"person1": 10,
				"person2": 20,
				"person3": 20
			}`
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_AuthorizationPolicyConfig{
						AuthorizationPolicyConfig: &configpb.AuthorizationPolicyConfig{
							Policy: jsonInput,
							Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
							Tags:   []string{},
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessConsentConfiguration", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_ConsentConfig{
						ConsentConfig: &configpb.ConsentConfiguration{
							Purpose: "Taking control",
							DataPoints: []string{
								"{ \"query\": \"\", \"returns\": [ { \"variable\": \"\"," +
									"\"properties\": [\"name\", \"email\", \"location\"] } ] }",
							},
							ApplicationId:  "gid:like-real-application-id",
							ValidityPeriod: uint64(86400),
							RevokeAfterUse: true,
							TokenStatus:    3,
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessTokenIntrospectConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_TokenIntrospectConfig{
						TokenIntrospectConfig: &configpb.TokenIntrospectConfig{
							TokenMatcher: &configpb.TokenIntrospectConfig_Opaque_{
								Opaque: &configpb.TokenIntrospectConfig_Opaque{
									Hint: "hint",
								}},
							Validation: &configpb.TokenIntrospectConfig_Online_{
								Online: &configpb.TokenIntrospectConfig_Online{
									UserinfoEndpoint: "https://data.example.com/userinfo",
								},
							},
							IkgNodeType: "MyUser",
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessIngestPipelineConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_IngestPipelineConfig{
						IngestPipelineConfig: &configpb.IngestPipelineConfig{
							Sources: []string{"source1", "source2"},
							Operations: []configpb.IngestPipelineOperation{
								configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_NODE,
								configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_RELATIONSHIP,
							},
							AppAgentToken: "", // Empty sensitive data
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessExternalDataResolverConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_ExternalDataResolverConfig{
						ExternalDataResolverConfig: &configpb.ExternalDataResolverConfig{
							Url:    "https://example.com/source",
							Method: "GET",
							Headers: map[string]*configpb.ExternalDataResolverConfig_Header{
								"Authorization": {Values: []string{"Bearer edyUTY"}},
								"Content-Type":  {Values: []string{"application/json"}},
							},
							RequestType:      1,
							RequestPayload:   []byte(`{"key": "value"}`),
							ResponseType:     1,
							ResponseSelector: ".",
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessEntityMatchingPipelineConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_EntityMatchingPipelineConfig{
						EntityMatchingPipelineConfig: &configpb.EntityMatchingPipelineConfig{
							NodeFilter: &configpb.EntityMatchingPipelineConfig_NodeFilter{
								SourceNodeTypes: []string{"employee"},
								TargetNodeTypes: []string{"user"},
							},
							SimilarityScoreCutoff: 0.8,
							PropertyMappingStatus: configpb.EntityMatchingPipelineConfig_STATUS_SUCCESS,
							EntityMatchingStatus:  configpb.EntityMatchingPipelineConfig_STATUS_PENDING,
							PropertyMappings: []*configpb.EntityMatchingPipelineConfig_PropertyMapping{
								{
									SourceNodeType:        "employee",
									SourceNodeProperty:    "email",
									TargetNodeType:        "user",
									TargetNodeProperty:    "address",
									SimilarityScoreCutoff: 0.9,
								},
							},
							RerunInterval: "1 day",
							LastRunTime:   timestamppb.New(time.Now()),
							ReportUrl:     wrapperspb.String("gs://some-path"),
							ReportType:    wrapperspb.String("csv"),
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessTrustScoreProfileConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_TrustScoreProfileConfig{
						TrustScoreProfileConfig: &configpb.TrustScoreProfileConfig{
							NodeClassification: "Agent",
							Dimensions: []*configpb.TrustScoreDimension{
								{
									Name:   configpb.TrustScoreDimension_NAME_VERIFICATION,
									Weight: 0.5,
								},
								{
									Name:   configpb.TrustScoreDimension_NAME_ORIGIN,
									Weight: 0.5,
								},
							},
							Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_THREE_HOURS,
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessKnowledgeQueryConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_KnowledgeQueryConfig{
						KnowledgeQueryConfig: &configpb.KnowledgeQueryConfig{
							Query:    `{"something":["like", "query"]}`,
							Status:   configpb.KnowledgeQueryConfig_STATUS_ACTIVE,
							PolicyId: "1 day",
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessEventSinkConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_EventSinkConfig{
						EventSinkConfig: &configpb.EventSinkConfig{
							Providers: map[string]*configpb.EventSinkConfig_Provider{
								"kafka": {
									Provider: &configpb.EventSinkConfig_Provider_Kafka{
										Kafka: &configpb.KafkaSinkConfig{
											Brokers:     []string{"broker.com"},
											Topic:       "your-topic-name",
											Username:    "your-name",
											Password:    "your-password",
											DisplayName: wrapperspb.String("like-real-provider-name"),
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
									DisplayName: wrapperspb.String("like-real-route-name"),
									Id:          wrapperspb.String("like-real-route-id"),
								},
							},
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessEventSinkConfigGrid", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_EventSinkConfig{
						EventSinkConfig: &configpb.EventSinkConfig{
							Providers: map[string]*configpb.EventSinkConfig_Provider{
								"azureeventgrid": {
									Provider: &configpb.EventSinkConfig_Provider_AzureEventGrid{
										AzureEventGrid: &configpb.AzureEventGridSinkConfig{
											TopicEndpoint: "https://ik-test.eventgrid.azure.net/api/events",
											AccessKey:     "your-access-key",
											DisplayName:   wrapperspb.String("like-real-provider-name"),
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
									DisplayName: wrapperspb.String("like-real-route-name"),
									Id:          wrapperspb.String("like-real-route-id"),
								},
							},
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadSuccessEventSinkConfigBus", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_EventSinkConfig{
						EventSinkConfig: &configpb.EventSinkConfig{
							Providers: map[string]*configpb.EventSinkConfig_Provider{
								"azureservicebus": {
									Provider: &configpb.EventSinkConfig_Provider_AzureServiceBus{
										AzureServiceBus: &configpb.AzureServiceBusSinkConfig{
											ConnectionString: "Endpoint=sb://ik-test.servicebus.windows.net/;SharedAccessKeyName=Root", //nolint:lll // easier to read not cut
											QueueOrTopicName: "your-queue",
											DisplayName:      wrapperspb.String("like-real-provider-name"),
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
									DisplayName: wrapperspb.String("like-real-route-name"),
									Id:          wrapperspb.String("like-real-route-id"),
								},
							},
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadCapturePipelineConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_CapturePipelineConfig{
						CapturePipelineConfig: &configpb.CapturePipelineConfig{
							ApiKeyId:     "something-like-a key",
							ApiKeySecret: "secret-key",
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadCapturePipelineTopicConfig", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithVersion(int64(0))
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Version:     0,
					Config: &configpb.ConfigNode_CapturePipelineTopicConfig{
						CapturePipelineTopicConfig: &configpb.CapturePipelineTopicConfig{
							TopicInputName:     "topic-name",
							TopicInputEndpoint: "https://example.com/topic-endpoint",
							TopicErrorName:     "topic-error",
							TopicSuccessName:   "topic-success",
							Script: &configpb.CapturePipelineTopicScriptConfig{
								Content: "content of the script",
							},
						},
					},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:     "gid:like-real-config-node-id",
					Name:   "like-real-config-node-name",
					Etag:   "123qwe",
					Config: &configpb.ConfigNode_AuthorizationPolicyConfig{},
				},
			}
			mockClient.EXPECT().
				ReadConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadConfigNode(ctx, configNodeRequest)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("ReadString", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			resp := configNodeRequest.String()

			Expect(resp).To(Not(BeNil()))
			Expect(err).To(Succeed())
			Expect(configNodeRequest).To(ContainSubstring("Read gid:like-real-config-node-id configuration"))
		})

		It("ReadWithNameString", func() {
			configNodeRequest, err := config.NewReadWithName("like-real-config-node-name")
			Expect(configNodeRequest).To(Not(BeNil()))
			Expect(err).To(Succeed())
		})

		It("ReadWithNameStringError", func() {
			configNodeRequest, err := config.NewReadWithName("1234")
			Expect(configNodeRequest).To(BeNil())
			Expect(err).ToNot(Succeed())
		})

		It("ReadStringNothing", func() {
			var configNodeRequest config.NodeRequest
			Expect(configNodeRequest.String()).To(Not(BeNil()))
			Expect(configNodeRequest.String()).To(ContainSubstring("Invalid empty request"))
		})
	})

	Describe("ConfigNodeCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateConfigNode(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil or not create request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("CreateAuthorizationPolicy", func() {
			jsonInput := `{
				"person1": 10,
				"person2": 20,
				"person3": 20
			}`
			configuration := &configpb.AuthorizationPolicyConfig{
				Policy: jsonInput,
				Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
				Tags:   []string{},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithAuthorizationPolicyConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"AuthorizationPolicyConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Policy": Equal(jsonInput),
							"Status": Equal(configpb.AuthorizationPolicyConfig_STATUS_ACTIVE),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateConsentConfiguration", func() {
			configuration := &configpb.ConsentConfiguration{
				Purpose: "Taking control",
				DataPoints: []string{
					"{ \"query\": \"\", \"returns\": [ { \"variable\": \"\"," +
						"\"properties\": [\"name\", \"email\", \"location\"] } ] }",
				},
				ApplicationId:  "gid:like-real-application-id",
				ValidityPeriod: uint64(86400),
				RevokeAfterUse: true,
				TokenStatus:    3,
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithConsentConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ConsentConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Purpose":        Equal("Taking control"),
							"ApplicationId":  Equal("gid:like-real-application-id"),
							"ValidityPeriod": Equal(uint64(86400)),
							"RevokeAfterUse": Equal(true),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateTokenIntrospectConfig", func() {
			configuration := &configpb.TokenIntrospectConfig{
				TokenMatcher: &configpb.TokenIntrospectConfig_Jwt{Jwt: &configpb.TokenIntrospectConfig_JWT{
					Issuer:   "https://example.com",
					Audience: "audience-id",
				}},
				Validation: &configpb.TokenIntrospectConfig_Offline_{
					Offline: &configpb.TokenIntrospectConfig_Offline{
						PublicJwks: [][]byte{
							[]byte(`{"kid":"abc","use":"sig","alg":"RS256","n":"--nothing-real-just-random-xyqwerasf--","kty":"RSA"}`), //nolint:lll
							[]byte(`{"kid":"jkl","use":"sig","alg":"RS256","n":"--nothing-real-just-random-435asdf43--","kty":"RSA"}`), //nolint:lll
						},
					},
				},
				IkgNodeType: "MyUser",
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithTokenIntrospectConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"TokenIntrospectConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"IkgNodeType": Equal("MyUser"),
							"TokenMatcher": Equal(&configpb.TokenIntrospectConfig_Jwt{
								Jwt: &configpb.TokenIntrospectConfig_JWT{
									Issuer:   "https://example.com",
									Audience: "audience-id",
								}}),
							"Validation": Equal(&configpb.TokenIntrospectConfig_Offline_{
								Offline: &configpb.TokenIntrospectConfig_Offline{
									PublicJwks: [][]byte{
										[]byte(`{"kid":"abc","use":"sig","alg":"RS256","n":"--nothing-real-just-random-xyqwerasf--","kty":"RSA"}`), //nolint:lll
										[]byte(`{"kid":"jkl","use":"sig","alg":"RS256","n":"--nothing-real-just-random-435asdf43--","kty":"RSA"}`), //nolint:lll
									},
								},
							}),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateIngestPipelineConfig", func() {
			appAgentToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaWQ6QUFBQUJXbHVaSGxyYVVSbGdBQUZEd0FBQUFBIiwic3ViIjoiZ2lkOkFBQUFCV2x1WkhscmFVUmxnQUFGRHdBQUFBQSIsImV4cCI6MjUzNDAyMjYxMTk5LCJpYXQiOjE1MTYyMzkwMjJ9.39Kc7pL8Vjf1S4qA6NHBGMP06TahR5Y9JOGSWKOo5Rw" //nolint:gosec,lll // there are no secrets
			configuration := &configpb.IngestPipelineConfig{
				Sources: []string{"source1", "source2", "source3"},
				Operations: []configpb.IngestPipelineOperation{
					configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_NODE,
					configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_RELATIONSHIP,
					configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_DELETE_NODE,
				},
				AppAgentToken: appAgentToken,
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithIngestPipelineConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"IngestPipelineConfig": test.EqualProto(
							&configpb.IngestPipelineConfig{
								Sources: []string{"source1", "source2", "source3"},
								Operations: []configpb.IngestPipelineOperation{
									configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_NODE,
									configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_RELATIONSHIP,
									configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_DELETE_NODE,
								},
								AppAgentToken: appAgentToken,
							},
						),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateExternalDataResolverConfig", func() {
			configuration := &configpb.ExternalDataResolverConfig{
				Url:    "https://example.com/source",
				Method: "GET",
				Headers: map[string]*configpb.ExternalDataResolverConfig_Header{
					"Authorization": {Values: []string{"Bearer edyUTY"}},
					"Content-Type":  {Values: []string{"application/json"}},
				},
				RequestType:      configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
				RequestPayload:   []byte(`{"key": "value"}`),
				ResponseType:     configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
				ResponseSelector: ".",
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithExternalDataResolverConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ExternalDataResolverConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Url":    Equal("https://example.com/source"),
							"Method": Equal("GET"),
							"Headers": Equal(map[string]*configpb.ExternalDataResolverConfig_Header{
								"Authorization": {Values: []string{"Bearer edyUTY"}},
								"Content-Type":  {Values: []string{"application/json"}},
							}),
							"RequestType":  Equal(configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON),
							"ResponseType": Equal(configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateEntityMatchingPipelineConfig", func() {
			now := timestamppb.New(time.Now())
			configuration := &configpb.EntityMatchingPipelineConfig{
				NodeFilter: &configpb.EntityMatchingPipelineConfig_NodeFilter{
					SourceNodeTypes: []string{"employee"},
					TargetNodeTypes: []string{"user"},
				},
				SimilarityScoreCutoff: 0.8,
				PropertyMappingStatus: configpb.EntityMatchingPipelineConfig_STATUS_SUCCESS,
				EntityMatchingStatus:  configpb.EntityMatchingPipelineConfig_STATUS_PENDING,
				PropertyMappings: []*configpb.EntityMatchingPipelineConfig_PropertyMapping{
					{
						SourceNodeType:        "employee",
						SourceNodeProperty:    "email",
						TargetNodeType:        "user",
						TargetNodeProperty:    "address",
						SimilarityScoreCutoff: 0.9,
					},
				},
				RerunInterval: "1 day",
				LastRunTime:   now,
				ReportUrl:     wrapperspb.String("gs://some-path"),
				ReportType:    wrapperspb.String("csv"),
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithEntityMatchingPipelineConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EntityMatchingPipelineConfig": test.EqualProto(
							&configpb.EntityMatchingPipelineConfig{
								NodeFilter: &configpb.EntityMatchingPipelineConfig_NodeFilter{
									SourceNodeTypes: []string{"employee"},
									TargetNodeTypes: []string{"user"},
								},
								SimilarityScoreCutoff: 0.8,
								PropertyMappingStatus: configpb.EntityMatchingPipelineConfig_STATUS_SUCCESS,
								EntityMatchingStatus:  configpb.EntityMatchingPipelineConfig_STATUS_PENDING,
								PropertyMappings: []*configpb.EntityMatchingPipelineConfig_PropertyMapping{
									{
										SourceNodeType:        "employee",
										SourceNodeProperty:    "email",
										TargetNodeType:        "user",
										TargetNodeProperty:    "address",
										SimilarityScoreCutoff: 0.9,
									},
								},
								RerunInterval: "1 day",
								LastRunTime:   now,
								ReportUrl:     wrapperspb.String("gs://some-path"),
								ReportType:    wrapperspb.String("csv"),
							},
						),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateTrustScoreProfileConfig", func() {
			configuration := &configpb.TrustScoreProfileConfig{
				NodeClassification: "Employee",
				Dimensions: []*configpb.TrustScoreDimension{
					{
						Name:   configpb.TrustScoreDimension_NAME_FRESHNESS,
						Weight: 0.9,
					},
				},
				Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_DAILY,
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithTrustScoreProfileConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"TrustScoreProfileConfig": test.EqualProto(
							&configpb.TrustScoreProfileConfig{
								NodeClassification: "Employee",
								Dimensions: []*configpb.TrustScoreDimension{
									{
										Name:   configpb.TrustScoreDimension_NAME_FRESHNESS,
										Weight: 0.9,
									},
								},
								Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_DAILY,
							},
						),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateKnowledgeQueryConfig", func() {
			configuration := &configpb.KnowledgeQueryConfig{
				Query:    `{"something": ["like", "json"]}`,
				Status:   configpb.KnowledgeQueryConfig_STATUS_ACTIVE,
				PolicyId: "gid:like-real-policy-id",
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithKnowledgeQueryConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"KnowledgeQueryConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateEventSinkConfig", func() {
			configuration := &configpb.EventSinkConfig{
				Providers: map[string]*configpb.EventSinkConfig_Provider{
					"kafka": {
						Provider: &configpb.EventSinkConfig_Provider_Kafka{
							Kafka: &configpb.KafkaSinkConfig{
								Brokers:     []string{"broker.com"},
								Topic:       "your-topic-name",
								Username:    "your-name",
								Password:    "your-password",
								DisplayName: wrapperspb.String("like-real-provider-name"),
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
						DisplayName: wrapperspb.String("like-real-route-name"),
						Id:          wrapperspb.String("like-real-route-id"),
					},
				},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithEventSinkConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EventSinkConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateEventSinkConfigGrid", func() {
			configuration := &configpb.EventSinkConfig{
				Providers: map[string]*configpb.EventSinkConfig_Provider{
					"azureeventgrid": {
						Provider: &configpb.EventSinkConfig_Provider_AzureEventGrid{
							AzureEventGrid: &configpb.AzureEventGridSinkConfig{
								TopicEndpoint: "https://ik-test.eventgrid.azure.net/api/events",
								AccessKey:     "your-access-key",
								DisplayName:   wrapperspb.String("like-real-provider-name"),
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
						DisplayName: wrapperspb.String("like-real-route-name"),
						Id:          wrapperspb.String("like-real-route-id"),
					},
				},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithEventSinkConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EventSinkConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateEventSinkConfigBus", func() {
			configuration := &configpb.EventSinkConfig{
				Providers: map[string]*configpb.EventSinkConfig_Provider{
					"azureservicebus": {
						Provider: &configpb.EventSinkConfig_Provider_AzureServiceBus{
							AzureServiceBus: &configpb.AzureServiceBusSinkConfig{
								ConnectionString: "Endpoint=sb://ik-test.servicebus.windows.net/;SharedAccessKeyName=Root", //nolint:lll // easier to read not cut
								QueueOrTopicName: "your-queue",
								DisplayName:      wrapperspb.String("like-real-provider-name"),
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
						DisplayName: wrapperspb.String("like-real-route-name"),
						Id:          wrapperspb.String("like-real-route-id"),
					},
				},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithEventSinkConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EventSinkConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("RegisterCapturePipelineConfig", func() {
			configuration := &configpb.RegisterCapturePipelineConfig{
				AppAgentToken: "eyJhbGciOiJIIkpXVCJ9.eyJpc3MiOiJnaWQkwMjJ9.39Kc7pL8Vjf1S4Oo5Rw",
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithCapturePipelineConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"CapturePipelineConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("RegisterCapturePipelineTopicConfig", func() {
			configuration := &configpb.RegisterCapturePipelineTopicConfig{
				CapturePipelineId: "gid:AAAAG3FQqyfhzEiUrpVHvab4ct4",
				Script: &configpb.RegisterCapturePipelineTopicConfig_Script{
					Content: "content of the script",
				},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithCapturePipelineTopicConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"CapturePipelineTopicConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			configuration := &configpb.ConsentConfiguration{
				Purpose: "Taking control",
				DataPoints: []string{
					"{ \"query\": \"\", \"returns\": [ { \"variable\": \"\"," +
						"\"properties\": [\"name\", \"email\", \"location\"] } ] }",
				},
				ApplicationId:  "gid:like",
				ValidityPeriod: uint64(86400),
				RevokeAfterUse: true,
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithConsentConfig(configuration)
			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("value length must be between 22 and 254 runes")))
		})

		It("CreateNonValidName", func() {
			configNodeRequest, err := config.NewCreate("1234")
			Expect(err).ToNot(Succeed())
			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("invalid nil or not create request")))
		})

		It("CreateError", func() {
			jsonInput := `{
				"person1": 10,
				"person2": 20,
				"person3": 20
			}`
			configuration := &configpb.AuthorizationPolicyConfig{
				Policy: jsonInput,
				Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
				Tags:   []string{},
			}
			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithAuthorizationPolicyConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"AuthorizationPolicyConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Policy": Equal(jsonInput),
							"Status": Equal(configpb.AuthorizationPolicyConfig_STATUS_ACTIVE),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("CreateString", func() {
			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			resp := configNodeRequest.String()

			Expect(resp).To(Not(BeNil()))
			Expect(err).To(Succeed())
			Expect(configNodeRequest).To(ContainSubstring("Create like-real-config-node-name configuration"))
		})
	})

	Describe("ConfigNodeUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateConfigNode(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil or not update request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("UpdateAuthorizationPolicy", func() {
			jsonInput := `{
				"person1": 10,
				"person2": 20,
				"person3": 20
			}`
			configuration := &configpb.AuthorizationPolicyConfig{
				Policy: jsonInput,
				Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
				Tags:   []string{},
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithAuthorizationPolicyConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"AuthorizationPolicyConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Policy": Equal(jsonInput),
							"Status": Equal(configpb.AuthorizationPolicyConfig_STATUS_ACTIVE),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateConsentConfiguration", func() {
			configuration := &configpb.ConsentConfiguration{
				Purpose: "Taking control",
				DataPoints: []string{
					"{ \"query\": \"\", \"returns\": [ { \"variable\": \"\"," +
						"\"properties\": [\"name\", \"email\", \"location\"] } ] }",
				},
				ApplicationId:  "gid:like-real-application-id",
				ValidityPeriod: uint64(86400),
				RevokeAfterUse: true,
				TokenStatus:    3,
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithConsentConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ConsentConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Purpose":       Equal("Taking control"),
							"ApplicationId": Equal("gid:like-real-application-id"),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateTokenIntrospectConfig", func() {
			configuration := &configpb.TokenIntrospectConfig{
				TokenMatcher: &configpb.TokenIntrospectConfig_Jwt{Jwt: &configpb.TokenIntrospectConfig_JWT{
					Issuer:   "https://example.com",
					Audience: "audience-id",
				}},
				Validation: &configpb.TokenIntrospectConfig_Offline_{
					Offline: &configpb.TokenIntrospectConfig_Offline{
						PublicJwks: [][]byte{
							[]byte(`{"kid":"abc","use":"sig","alg":"RS256","n":"--nothing-real-just-random-xyqwerasf--","kty":"RSA"}`), //nolint:lll
							[]byte(`{"kid":"jkl","use":"sig","alg":"RS256","n":"--nothing-real-just-random-435asdf43--","kty":"RSA"}`), //nolint:lll
						},
					},
				},
				IkgNodeType: "MyUser",
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithTokenIntrospectConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"TokenIntrospectConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"IkgNodeType": Equal("MyUser"),
							"TokenMatcher": Equal(&configpb.TokenIntrospectConfig_Jwt{
								Jwt: &configpb.TokenIntrospectConfig_JWT{
									Issuer:   "https://example.com",
									Audience: "audience-id",
								}}),
							"Validation": Equal(&configpb.TokenIntrospectConfig_Offline_{
								Offline: &configpb.TokenIntrospectConfig_Offline{
									PublicJwks: [][]byte{
										[]byte(`{"kid":"abc","use":"sig","alg":"RS256","n":"--nothing-real-just-random-xyqwerasf--","kty":"RSA"}`), //nolint:lll
										[]byte(`{"kid":"jkl","use":"sig","alg":"RS256","n":"--nothing-real-just-random-435asdf43--","kty":"RSA"}`), //nolint:lll
									},
								},
							}),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateIngestPipelineConfig", func() {
			appAgentToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaWQ6QUFBQUJXbHVaSGxyYVVSbGdBQUZEd0FBQUFBIiwic3ViIjoiZ2lkOkFBQUFCV2x1WkhscmFVUmxnQUFGRHdBQUFBQSIsImV4cCI6MjUzNDAyMjYxMTk5LCJpYXQiOjE1MTYyMzkwMjJ9.39Kc7pL8Vjf1S4qA6NHBGMP06TahR5Y9JOGSWKOo5Rw" //nolint:gosec,lll // there are no secrets
			configuration := &configpb.IngestPipelineConfig{
				Sources: []string{"source1", "source2", "source3"},
				Operations: []configpb.IngestPipelineOperation{
					configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_NODE,
					configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_RELATIONSHIP,
					configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_DELETE_NODE,
				},
				AppAgentToken: appAgentToken,
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithIngestPipelineConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"IngestPipelineConfig": test.EqualProto(
							&configpb.IngestPipelineConfig{
								Sources: []string{"source1", "source2", "source3"},
								Operations: []configpb.IngestPipelineOperation{
									configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_NODE,
									configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_UPSERT_RELATIONSHIP,
									configpb.IngestPipelineOperation_INGEST_PIPELINE_OPERATION_DELETE_NODE,
								},
								AppAgentToken: appAgentToken,
							},
						),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateExternalDataResolverConfig", func() {
			configuration := &configpb.ExternalDataResolverConfig{
				Url:    "https://example.com/source",
				Method: "GET",
				Headers: map[string]*configpb.ExternalDataResolverConfig_Header{
					"Authorization": {Values: []string{"Bearer edyUTY"}},
					"Content-Type":  {Values: []string{"application/json"}},
				},
				RequestType:      configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
				RequestPayload:   []byte(`{"key": "value"}`),
				ResponseType:     configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON,
				ResponseSelector: ".",
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithExternalDataResolverConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ExternalDataResolverConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Url":    Equal("https://example.com/source"),
							"Method": Equal("GET"),
							"Headers": Equal(map[string]*configpb.ExternalDataResolverConfig_Header{
								"Authorization": {Values: []string{"Bearer edyUTY"}},
								"Content-Type":  {Values: []string{"application/json"}},
							}),
							"RequestType":  Equal(configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON),
							"ResponseType": Equal(configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateEntityMatchingPipelineConfig", func() {
			now := timestamppb.New(time.Now())
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
				RerunInterval: "1 hour",
				LastRunTime:   now,
				ReportUrl:     wrapperspb.String("gs://some-otherpath"),
				ReportType:    wrapperspb.String("csv"),
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithEntityMatchingPipelineConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EntityMatchingPipelineConfig": test.EqualProto(
							&configpb.EntityMatchingPipelineConfig{
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
								RerunInterval: "1 hour",
								LastRunTime:   now,
								ReportUrl:     wrapperspb.String("gs://some-otherpath"),
								ReportType:    wrapperspb.String("csv"),
							},
						),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateTrustScoreProfileConfig", func() {
			configuration := &configpb.TrustScoreProfileConfig{
				Dimensions: []*configpb.TrustScoreDimension{
					{
						Name:   configpb.TrustScoreDimension_NAME_COMPLETENESS,
						Weight: 0.92,
					},
				},
				Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_DAILY,
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithTrustScoreProfileConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"TrustScoreProfileConfig": test.EqualProto(
							&configpb.TrustScoreProfileConfig{
								Dimensions: []*configpb.TrustScoreDimension{
									{
										Name:   configpb.TrustScoreDimension_NAME_COMPLETENESS,
										Weight: 0.92,
									},
								},
								Schedule: configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_DAILY,
							},
						),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateKnowledgeQueryConfig", func() {
			configuration := &configpb.KnowledgeQueryConfig{
				Query:    `{"something": ["like", "json"]}`,
				Status:   configpb.KnowledgeQueryConfig_STATUS_ACTIVE,
				PolicyId: "gid:like-real-policy-id",
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithKnowledgeQueryConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"KnowledgeQueryConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateEventSinkConfig", func() {
			configuration := &configpb.EventSinkConfig{
				Providers: map[string]*configpb.EventSinkConfig_Provider{
					"kafka": {
						Provider: &configpb.EventSinkConfig_Provider_Kafka{
							Kafka: &configpb.KafkaSinkConfig{
								Brokers:     []string{"broker.com"},
								Topic:       "your-topic-name",
								Username:    "your-name-update",
								Password:    "your-password-update",
								DisplayName: wrapperspb.String("like-real-provider-name"),
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
						DisplayName: wrapperspb.String("like-real-route-name"),
						Id:          wrapperspb.String("like-real-route-id"),
					},
				},
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithEventSinkConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EventSinkConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateEventSinkConfigGrid", func() {
			configuration := &configpb.EventSinkConfig{
				Providers: map[string]*configpb.EventSinkConfig_Provider{
					"azureeventgrid": {
						Provider: &configpb.EventSinkConfig_Provider_AzureEventGrid{
							AzureEventGrid: &configpb.AzureEventGridSinkConfig{
								TopicEndpoint: "https://ik-test.eventgrid.azure.net/api/events",
								AccessKey:     "your-access-key",
								DisplayName:   wrapperspb.String("like-real-provider-name"),
							},
						},
					},
				},
				Routes: []*configpb.EventSinkConfig_Route{
					{
						ProviderId:     "azureeventgrid",
						StopProcessing: true,
						Filter: &configpb.EventSinkConfig_Route_EventType{
							EventType: "indykite.eventsink.config.update",
						},
						DisplayName: wrapperspb.String("like-real-route-name"),
						Id:          wrapperspb.String("like-real-route-id"),
					},
				},
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithEventSinkConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EventSinkConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateEventSinkConfigBus", func() {
			configuration := &configpb.EventSinkConfig{
				Providers: map[string]*configpb.EventSinkConfig_Provider{
					"azureservicebus": {
						Provider: &configpb.EventSinkConfig_Provider_AzureServiceBus{
							AzureServiceBus: &configpb.AzureServiceBusSinkConfig{
								ConnectionString: "Endpoint=sb://ik-test.servicebus.windows.net/;SharedAccessKeyName=Root", //nolint:lll // easier to read not cut
								QueueOrTopicName: "your-queue",
								DisplayName:      wrapperspb.String("like-real-provider-name"),
							},
						},
					},
				},
				Routes: []*configpb.EventSinkConfig_Route{
					{
						ProviderId:     "azureservicebus",
						StopProcessing: true,
						Filter: &configpb.EventSinkConfig_Route_EventType{
							EventType: "indykite.eventsink.config.update",
						},
						DisplayName: wrapperspb.String("like-real-route-name"),
						Id:          wrapperspb.String("like-real-route-id"),
					},
				},
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName()
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithEventSinkConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EventSinkConfig": test.EqualProto(configuration),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			configNodeRequest, err := config.NewUpdate("12345")
			Ω(err).To(Succeed())
			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			jsonInput := `{
				"person1": 10,
				"person2": 20,
				"person3": 20
			}`
			configuration := &configpb.AuthorizationPolicyConfig{
				Policy: jsonInput,
				Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
				Tags:   []string{},
			}
			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithAuthorizationPolicyConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"AuthorizationPolicyConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Policy": Equal(jsonInput),
							"Status": Equal(configpb.AuthorizationPolicyConfig_STATUS_ACTIVE),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("UpdateString", func() {
			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			resp := configNodeRequest.String()

			Expect(resp).To(Not(BeNil()))
			Expect(err).To(Succeed())
			Expect(configNodeRequest).To(ContainSubstring("Update gid:like-real-config-node-id configuration"))
		})
	})

	Describe("ConfigNodeDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteConfigNode(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil or not delete request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			configNodeRequest, err := config.NewDelete("gid:like")
			Ω(err).To(Succeed())
			resp, err := configClient.DeleteConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			configNodeRequest, err := config.NewDelete("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithPreCondition("qwert1234")
			beResp := &configpb.DeleteConfigNodeResponse{}

			mockClient.EXPECT().
				DeleteConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.DeleteConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			configNodeRequest, err := config.NewDelete("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			beResp := &configpb.DeleteConfigNodeResponse{}

			mockClient.EXPECT().
				DeleteConfigNode(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteConfigNode(ctx, configNodeRequest)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("DeleteString", func() {
			configNodeRequest, err := config.NewDelete("gid:like-real-config-node-id")
			resp := configNodeRequest.String()

			Expect(resp).To(Not(BeNil()))
			Expect(err).To(Succeed())
			Expect(configNodeRequest).To(ContainSubstring("Delete gid:like-real-config-node-id configuration"))
		})
	})

	Describe("ListConfigNodeVersions", func() {
		It("Nil request", func() {
			resp, err := configClient.ListConfigNodeVersions(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil or not read request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			configNodeRequest, err := config.NewListVersions("gid:like")
			Ω(err).To(Succeed())
			resp, err := configClient.ListConfigNodeVersions(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("value length must be between 22 and 254 runes, inclusive")))
		})

		It("ListVersionsSuccess", func() {
			configNodeRequest, err := config.NewListVersions("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			Expect(configNodeRequest).ToNot(BeNil())
			mockResp := &configpb.ListConfigNodeVersionsResponse{
				ConfigNodes: []*configpb.ConfigNode{
					{
						Id:          "gid:like-real-config-node-id",
						Name:        "like-real-config-node-name",
						DisplayName: "Like Real Config-Node Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						CustomerId:  "gid:like-real-customer-id",
						AppSpaceId:  "gid:like-real-app-space-id",
						Etag:        "123qwe",
						Version:     0,
						Config: &configpb.ConfigNode_AuthorizationPolicyConfig{
							AuthorizationPolicyConfig: &configpb.AuthorizationPolicyConfig{
								Policy: `{"person1": 10,"person2": 20,"person3": 20}`,
								Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
								Tags:   []string{},
							},
						},
					},
					{
						Id:          "gid:like-real-config-node-id",
						Name:        "like-real-config-node-name",
						DisplayName: "Like Real Config-Node Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						CustomerId:  "gid:like-real-customer-id",
						AppSpaceId:  "gid:like-real-app-space-id",
						Etag:        "123qwe",
						Version:     1,
						Config: &configpb.ConfigNode_AuthorizationPolicyConfig{
							AuthorizationPolicyConfig: &configpb.AuthorizationPolicyConfig{
								Policy: `{"person1": 11,"person2": 22,"person33": 20}`,
								Status: configpb.AuthorizationPolicyConfig_STATUS_ACTIVE,
								Tags:   []string{},
							},
						},
					},
				},
			}

			mockClient.EXPECT().
				ListConfigNodeVersions(
					gomock.Any(),
					test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id": Equal("gid:like-real-config-node-id"),
					}))),
					gomock.Any(),
				).Return(mockResp, status.Error(codes.InvalidArgument, "status error")).AnyTimes()

			resp, err := configClient.ListConfigNodeVersions(ctx, configNodeRequest)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
