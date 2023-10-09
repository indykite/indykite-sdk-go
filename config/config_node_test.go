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

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
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

		It("ReadSuccessReadId", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithBookmarks([]string{"something-like-bookmark-which-is-long-enough"})
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					TenantId:    "gid:like-real-tenant-id",
					Etag:        "123qwe",
					Config: &configpb.ConfigNode_ReadidProviderConfig{
						ReadidProviderConfig: &configpb.ReadIDProviderConfig{
							SubmitterSecret:   "submitter-secret-12345678901234567890",
							ManagerSecret:     "manager-secret-123456789012345678901234",
							SubmitterPassword: "submitter-password",
							HostAddress:       "<https://saas-preprod.readid.com>",
							PropertyMap: map[string]*configpb.ReadIDProviderConfig_Property{
								"Propertymap": {Expression: "property", Enabled: true},
							},
							UniquePropertyName: "uniquepropertyname",
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

		It("ReadSuccessAuthorizationPolicy", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithBookmarks([]string{"something-like-bookmark-which-is-long-enough"})
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
					TenantId:    "gid:like-real-tenant-id",
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

		It("ReadSuccessOAuth2Client", func() {
			configNodeRequest, err := config.NewRead("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithBookmarks([]string{"something-like-bookmark-which-is-long-enough"})
			beResp := &configpb.ReadConfigNodeResponse{
				ConfigNode: &configpb.ConfigNode{
					Id:          "gid:like-real-config-node-id",
					Name:        "like-real-config-node-name",
					DisplayName: "Like Real Config-Node Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					TenantId:    "gid:like-real-tenant-id",
					Etag:        "123qwe",
					Config: &configpb.ConfigNode_Oauth2ClientConfig{
						Oauth2ClientConfig: &configpb.OAuth2ClientConfig{
							ProviderType:  configpb.ProviderType_PROVIDER_TYPE_GOOGLE_COM,
							ClientId:      "manager-secret-123456789012345678901234",
							ClientSecret:  "submitter-password",
							DefaultScopes: []string{"openid", "profile", "email"},
							AllowedScopes: []string{"openid", "profile", "email"},
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
					Config: &configpb.ConfigNode_ReadidProviderConfig{},
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

		It("CreateReadId", func() {
			configuration := &configpb.ReadIDProviderConfig{
				SubmitterSecret:   "submitter-secret-12345678901234567890",
				ManagerSecret:     "manager-secret-123456789012345678901234",
				SubmitterPassword: "submitter-password",
				HostAddress:       "<https://saas-preprod.readid.com>",
				PropertyMap: map[string]*configpb.ReadIDProviderConfig_Property{
					"Propertymap": {Expression: "property", Enabled: true},
				},
				UniquePropertyName: "uniquepropertyname",
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDescription("Description ConfigNode")
			configNodeRequest.WithBookmarks([]string{"something-like-bookmark-which-is-long-enough"})
			configNodeRequest.WithReadidProviderConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ReadidProviderConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"SubmitterSecret": Equal("submitter-secret-12345678901234567890"),
							"ManagerSecret":   Equal("manager-secret-123456789012345678901234"),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateOAuth2Client", func() {
			configuration := &configpb.OAuth2ClientConfig{
				ProviderType:  configpb.ProviderType_PROVIDER_TYPE_GOOGLE_COM,
				ClientId:      "client-id-123456789012345678901234",
				ClientSecret:  "client-secret",
				DefaultScopes: []string{"openid", "profile", "email"},
				AllowedScopes: []string{"openid", "profile", "email"},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithOAuth2ClientConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"Oauth2ClientConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"ClientId":     Equal("client-id-123456789012345678901234"),
							"ClientSecret": Equal("client-secret"),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateWebauthnProvider", func() {
			configuration := &configpb.WebAuthnProviderConfig{
				RelyingParties:          map[string]string{"http://localhost": "localhost"},
				AttestationPreference:   configpb.ConveyancePreference_CONVEYANCE_PREFERENCE_NONE,
				AuthenticatorAttachment: configpb.AuthenticatorAttachment_AUTHENTICATOR_ATTACHMENT_DEFAULT,
				RequireResidentKey:      false,
				UserVerification:        configpb.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED,
				RegistrationTimeout:     &durationpb.Duration{Seconds: 30},
				AuthenticationTimeout:   &durationpb.Duration{Seconds: 60},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithWebauthnProviderConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"WebauthnProviderConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"AttestationPreference": Equal(configpb.ConveyancePreference_CONVEYANCE_PREFERENCE_NONE),
							"AuthenticatorAttachment": Equal(
								configpb.AuthenticatorAttachment_AUTHENTICATOR_ATTACHMENT_DEFAULT,
							),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateEmailNotification", func() {
			configuration := &configpb.EmailServiceConfig{
				DefaultFromAddress: &configpb.Email{Name: "indy", Address: "sos@example.com"},
				Default:            true,
				Provider: &configpb.EmailServiceConfig_Sendgrid{
					Sendgrid: &configpb.SendGridProviderConfig{
						ApiKey:      "12qe2e3r4t5y6u7i8o92t2t3t4t5y6u7i83h5h6",
						SandboxMode: false,
						IpPoolName: &wrapperspb.StringValue{
							Value: "Like real IpPoolName Name",
						},
						Host: &wrapperspb.StringValue{
							Value: "https://api.sendgrid.com",
						},
					},
				},
				InvitationMessage: &configpb.EmailDefinition{
					Email: &configpb.EmailDefinition_Message{
						Message: &configpb.EmailMessage{
							From:    &configpb.Email{Name: "indy", Address: "sos@example.com"},
							Subject: "invitation subject",
						},
					},
				},
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithEmailNotificationConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EmailServiceConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"DefaultFromAddress": PointTo(MatchFields(IgnoreExtras, Fields{
								"Name":    Equal("indy"),
								"Address": Equal("sos@example.com"),
							})),
							"Default": Equal(true),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateAuthFlow", func() {
			configuration := &configpb.AuthFlowConfig{
				SourceFormat: configpb.AuthFlowConfig_FORMAT_BARE_JSON,
				Source:       []byte("wolf"),
				Default:      false,
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithAuthFlowConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"AuthFlowConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"SourceFormat": Equal(configpb.AuthFlowConfig_FORMAT_BARE_JSON),
							"Default":      Equal(false),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
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
				Bookmark:   "something-like-bookmark-which-is-long-enough",
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

		It("CreateKnowledgeGraphSchema", func() {
			schemaInput := `
				type Person implements DigitalTwin {
					externalId: String!
					digitalTwinId: String!
					tenantId: String!
					kind: DigitalTwinKind!
					tags: [String!]!
				}
			`
			configuration := &configpb.KnowledgeGraphSchemaConfig{
				Schema: schemaInput,
			}

			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithKnowledgeGraphSchemaConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"KnowledgeGraphSchemaConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Schema": Equal(schemaInput),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			configuration := &configpb.ReadIDProviderConfig{
				SubmitterSecret:   "submitter-secret",
				ManagerSecret:     "manager-secret-123456789012345678901234",
				SubmitterPassword: "submitter-password",
				HostAddress:       "<https://saas-preprod.readid.com>",
				PropertyMap: map[string]*configpb.ReadIDProviderConfig_Property{
					"Propertymap": {Expression: "property", Enabled: true},
				},
				UniquePropertyName: "uniquepropertyname",
			}
			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithReadidProviderConfig(configuration)
			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("SubmitterSecret: value length must be at least 36 runes")))

		})

		It("CreateNonValidName", func() {

			configNodeRequest, err := config.NewCreate("1234")
			Expect(err).ToNot(Succeed())
			resp, err := configClient.CreateConfigNode(ctx, configNodeRequest)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("invalid nil or not create request")))

		})

		It("CreateError", func() {
			configuration := &configpb.ReadIDProviderConfig{
				SubmitterSecret:   "submitter-secret-12345678901234567890",
				ManagerSecret:     "manager-secret-123456789012345678901234",
				SubmitterPassword: "submitter-password",
				HostAddress:       "<https://saas-preprod.readid.com>",
				PropertyMap: map[string]*configpb.ReadIDProviderConfig_Property{
					"Propertymap": {Expression: "property", Enabled: true},
				},
				UniquePropertyName: "uniquepropertyname",
			}
			configNodeRequest, err := config.NewCreate("like-real-config-node-name")
			Ω(err).To(Succeed())
			configNodeRequest.ForLocation("gid:like-real-customer-id")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithReadidProviderConfig(configuration)

			beResp := &configpb.CreateConfigNodeResponse{}

			mockClient.EXPECT().CreateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":     Equal("like-real-config-node-name"),
					"Location": Equal("gid:like-real-customer-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ReadidProviderConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"SubmitterSecret": Equal("submitter-secret-12345678901234567890"),
							"ManagerSecret":   Equal("manager-secret-123456789012345678901234"),
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

		It("UpdateReadId", func() {
			configuration := &configpb.ReadIDProviderConfig{
				SubmitterSecret:   "submitter-secret-12345678901234567890",
				ManagerSecret:     "manager-secret-123456789012345678901234",
				SubmitterPassword: "submitter-password",
				HostAddress:       "<https://saas-preprod.readid.com>",
				PropertyMap: map[string]*configpb.ReadIDProviderConfig_Property{
					"Propertymap": {Expression: "property", Enabled: true},
				},
				UniquePropertyName: "uniquepropertyname",
			}
			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.EmptyDescription("Description")
			configNodeRequest.WithDescription("Description")
			configNodeRequest.WithPreCondition("qwert1234")
			configNodeRequest.WithBookmarks([]string{"something-like-bookmark-which-is-long-enough"})
			configNodeRequest.WithReadidProviderConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ReadidProviderConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"SubmitterSecret": Equal("submitter-secret-12345678901234567890"),
							"ManagerSecret":   Equal("manager-secret-123456789012345678901234"),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateOAuth2Client", func() {
			configuration := &configpb.OAuth2ClientConfig{
				ProviderType:  configpb.ProviderType_PROVIDER_TYPE_GOOGLE_COM,
				ClientId:      "client-id-123456789012345678901234",
				ClientSecret:  "client-secret",
				DefaultScopes: []string{"openid", "profile", "email"},
				AllowedScopes: []string{"openid", "profile", "email"},
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithOAuth2ClientConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"Oauth2ClientConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"ClientId":     Equal("client-id-123456789012345678901234"),
							"ClientSecret": Equal("client-secret"),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateWebAuthnProvider", func() {
			configuration := &configpb.WebAuthnProviderConfig{
				RelyingParties:          map[string]string{"http://localhost": "localhost"},
				AttestationPreference:   configpb.ConveyancePreference_CONVEYANCE_PREFERENCE_NONE,
				AuthenticatorAttachment: configpb.AuthenticatorAttachment_AUTHENTICATOR_ATTACHMENT_DEFAULT,
				RequireResidentKey:      false,
				UserVerification:        configpb.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED,
				RegistrationTimeout:     &durationpb.Duration{Seconds: 30},
				AuthenticationTimeout:   &durationpb.Duration{Seconds: 60},
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithWebauthnProviderConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"WebauthnProviderConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"AttestationPreference": Equal(configpb.ConveyancePreference_CONVEYANCE_PREFERENCE_NONE),
							"AuthenticatorAttachment": Equal(
								configpb.AuthenticatorAttachment_AUTHENTICATOR_ATTACHMENT_DEFAULT,
							),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateEmailNotification", func() {
			configuration := &configpb.EmailServiceConfig{
				DefaultFromAddress: &configpb.Email{Name: "indy", Address: "sos@example.com"},
				Default:            true,
				Provider: &configpb.EmailServiceConfig_Sendgrid{
					Sendgrid: &configpb.SendGridProviderConfig{
						ApiKey:      "12qe2e3r4t5y6u7i8o92t2t3t4t5y6u7i83h5h6",
						SandboxMode: false,
						IpPoolName: &wrapperspb.StringValue{
							Value: "Like real IpPoolName Name",
						},
						Host: &wrapperspb.StringValue{
							Value: "https://api.sendgrid.com",
						},
					},
				},
				InvitationMessage: &configpb.EmailDefinition{
					Email: &configpb.EmailDefinition_Message{
						Message: &configpb.EmailMessage{
							From:    &configpb.Email{Name: "indy", Address: "sos@example.com"},
							Subject: "invitation subject",
						},
					},
				},
			}
			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithEmailNotificationConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"EmailServiceConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"DefaultFromAddress": PointTo(MatchFields(IgnoreExtras, Fields{
								"Name":    Equal("indy"),
								"Address": Equal("sos@example.com"),
							})),
							"Default": Equal(true),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateAuthFlow", func() {
			configuration := &configpb.AuthFlowConfig{
				SourceFormat: configpb.AuthFlowConfig_FORMAT_BARE_JSON,
				Source:       []byte("wolf"),
				Default:      false,
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithAuthFlowConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"AuthFlowConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"SourceFormat": Equal(configpb.AuthFlowConfig_FORMAT_BARE_JSON),
							"Default":      Equal(false),
						})),
					})),
				}))),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateConfigNode(ctx, configNodeRequest)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
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
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithAuthorizationPolicyConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
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

		It("UpdateKnowledgeGraphSchema", func() {
			schemaInput := `
				type Person implements DigitalTwin {
					externalId: String!
					digitalTwinId: String!
					tenantId: String!
					kind: DigitalTwinKind!
					tags: [String!]!
				}
			`
			configuration := &configpb.KnowledgeGraphSchemaConfig{
				Schema: schemaInput,
			}

			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.EmptyDisplayName("Like real ConfigNode Name")
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithKnowledgeGraphSchemaConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{
				Id:         "gid:like-real-config-node-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"KnowledgeGraphSchemaConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"Schema": Equal(schemaInput),
						})),
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
			configuration := &configpb.ReadIDProviderConfig{
				SubmitterSecret:   "submitter-secret-12345678901234567890",
				ManagerSecret:     "manager-secret-123456789012345678901234",
				SubmitterPassword: "submitter-password",
				HostAddress:       "<https://saas-preprod.readid.com>",
				PropertyMap: map[string]*configpb.ReadIDProviderConfig_Property{
					"Propertymap": {Expression: "property", Enabled: true},
				},
				UniquePropertyName: "uniquepropertyname",
			}
			configNodeRequest, err := config.NewUpdate("gid:like-real-config-node-id")
			Ω(err).To(Succeed())
			configNodeRequest.WithDisplayName("Like real ConfigNode Name Update")
			configNodeRequest.WithReadidProviderConfig(configuration)

			beResp := &configpb.UpdateConfigNodeResponse{}

			mockClient.EXPECT().UpdateConfigNode(
				gomock.Any(),
				test.WrapMatcher(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id": Equal("gid:like-real-config-node-id"),
					"Config": PointTo(MatchFields(IgnoreExtras, Fields{
						"ReadidProviderConfig": PointTo(MatchFields(IgnoreExtras, Fields{
							"SubmitterSecret": Equal("submitter-secret-12345678901234567890"),
							"ManagerSecret":   Equal("manager-secret-123456789012345678901234"),
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
			configNodeRequest.WithBookmarks([]string{"something-like-bookmark-which-is-long-enough"})
			beResp := &configpb.DeleteConfigNodeResponse{
				Bookmark: "something-like-bookmark-which-is-long-enough",
			}

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
						TenantId:    "gid:like-real-tenant-id",
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
						TenantId:    "gid:like-real-tenant-id",
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
