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
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/indykite/indykite-sdk-go/config"
	sdkerrors "github.com/indykite/indykite-sdk-go/errors"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	configmock "github.com/indykite/indykite-sdk-go/test/config/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OAuth2Provider", func() {
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

	Describe("OAuth2Provider", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadOAuth2Provider(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadOAuth2ProviderRequest{Id: "gid:like"}
			resp, err := configClient.ReadOAuth2Provider(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))

		})

		It("ReadSuccess", func() {
			req := &configpb.ReadOAuth2ProviderRequest{Id: "gid:like-real-oauth2-provider-id"}
			beResp := &configpb.ReadOAuth2ProviderResponse{
				Oauth2Provider: &configpb.OAuth2Provider{
					Id:          "gid:like-real-oauth2-provider-id",
					Name:        "like-real-oauth2-provider-name",
					DisplayName: "Like Real OAuth2-Provider Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					CustomerId:  "gid:like-real-customer-id",
					AppSpaceId:  "gid:like-real-app-space-id",
					Etag:        "123qwe",
					Config: &configpb.OAuth2ProviderConfig{
						GrantTypes: []configpb.GrantType{configpb.
							GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
						ResponseTypes: []configpb.ResponseType{configpb.
							ResponseType_RESPONSE_TYPE_CODE, configpb.ResponseType_RESPONSE_TYPE_TOKEN},
						Scopes: []string{"openid", "profile", "email"},
						TokenEndpointAuthMethod: []configpb.TokenEndpointAuthMethod{configpb.
							TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST},
						TokenEndpointAuthSigningAlg: []string{"ES256", "ES384", "ES512"},
						RequestUris:                 make([]string, 0),
						RequestObjectSigningAlg:     "ES256",
						FrontChannelLoginUri:        map[string]string{"default": "http://localhost:3000/login/oauth2"},
						FrontChannelConsentUri:      map[string]string{"default": "http://localhost:3000/consent"},
					},
				},
			}
			mockClient.EXPECT().
				ReadOAuth2Provider(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadOAuth2Provider(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			req := &configpb.ReadOAuth2ProviderRequest{Id: "gid:like-real-oauth2-provider-id"}
			beResp := &configpb.ReadOAuth2ProviderResponse{
				Oauth2Provider: &configpb.OAuth2Provider{
					Id:     "gid:like-real-oauth2-provider-id",
					Name:   "like-real-oauth2-provider-name",
					Etag:   "123qwe",
					Config: &configpb.OAuth2ProviderConfig{},
				},
			}
			mockClient.EXPECT().
				ReadOAuth2Provider(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadOAuth2Provider(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("OAuth2ProviderCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateOAuth2Provider(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Provider Name"}
			req := &configpb.CreateOAuth2ProviderRequest{
				AppSpaceId:  "gid:like-real-app-space-id",
				Name:        "like-real-oauth2-provider-name",
				DisplayName: displayNamePb,
				Config: &configpb.OAuth2ProviderConfig{
					GrantTypes: []configpb.GrantType{configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.
						ResponseType_RESPONSE_TYPE_CODE, configpb.ResponseType_RESPONSE_TYPE_TOKEN},
					Scopes: []string{"openid", "profile", "email"},
					TokenEndpointAuthMethod: []configpb.TokenEndpointAuthMethod{configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST},
					TokenEndpointAuthSigningAlg: []string{"ES256", "ES384", "ES512"},
					RequestUris:                 make([]string, 0),
					RequestObjectSigningAlg:     "ES256",
					FrontChannelLoginUri:        map[string]string{"default": "http://localhost:3000/login/oauth2"},
					FrontChannelConsentUri:      map[string]string{"default": "http://localhost:3000/consent"},
				},
			}
			beResp := &configpb.CreateOAuth2ProviderResponse{
				Id:         "gid:like-real-oauth2-provider-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateOAuth2Provider(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateOAuth2Provider(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Provider Name"}
			req := &configpb.CreateOAuth2ProviderRequest{
				AppSpaceId:  "error-app-space-id",
				Name:        "like-real-oauth2-provider-name",
				DisplayName: displayNamePb,
			}

			resp, err := configClient.CreateOAuth2Provider(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Provider Name"}
			req := &configpb.CreateOAuth2ProviderRequest{
				AppSpaceId:  "gid:like-real-app-space-id",
				Name:        "like-real-oauth2-provider-name",
				DisplayName: displayNamePb,
				Config: &configpb.OAuth2ProviderConfig{
					GrantTypes: []configpb.GrantType{configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.
						ResponseType_RESPONSE_TYPE_CODE, configpb.ResponseType_RESPONSE_TYPE_TOKEN},
					Scopes: []string{"openid", "profile", "email"},
					TokenEndpointAuthMethod: []configpb.TokenEndpointAuthMethod{configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST},
					TokenEndpointAuthSigningAlg: []string{"ES256", "ES384", "ES512"},
					RequestUris:                 make([]string, 0),
					RequestObjectSigningAlg:     "ES256",
					FrontChannelLoginUri:        map[string]string{"default": "http://localhost:3000/login/oauth2"},
					FrontChannelConsentUri:      map[string]string{"default": "http://localhost:3000/consent"},
				},
			}
			beResp := &configpb.CreateOAuth2ProviderResponse{}
			mockClient.EXPECT().CreateOAuth2Provider(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateOAuth2Provider(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

	})

	Describe("OAuth2ProviderUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateOAuth2Provider(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Provider Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateOAuth2ProviderRequest{
				Id:          "gid:like-real-oauth2-provider-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
				Config: &configpb.OAuth2ProviderConfig{
					GrantTypes: []configpb.GrantType{configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.
						ResponseType_RESPONSE_TYPE_CODE, configpb.ResponseType_RESPONSE_TYPE_TOKEN},
					Scopes: []string{"openid", "profile", "email"},
					TokenEndpointAuthMethod: []configpb.TokenEndpointAuthMethod{configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST},
					TokenEndpointAuthSigningAlg: []string{"ES256", "ES384", "ES512"},
					RequestUris:                 make([]string, 0),
					RequestObjectSigningAlg:     "ES256",
					FrontChannelLoginUri:        map[string]string{"default": "http://localhost:3000/login/oauth2"},
					FrontChannelConsentUri:      map[string]string{"default": "http://localhost:3000/consent"},
				},
			}
			beResp := &configpb.UpdateOAuth2ProviderResponse{
				Id:         "gid:like-real-oauth2-provider-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateOAuth2Provider(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateOAuth2Provider(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Provider Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateOAuth2ProviderRequest{
				Id:          "wrong-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
				Config:      &configpb.OAuth2ProviderConfig{},
			}
			resp, err := configClient.UpdateOAuth2Provider(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Provider Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateOAuth2ProviderRequest{
				Id:          "gid:like-real-oauth2-provider-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
				Config: &configpb.OAuth2ProviderConfig{
					GrantTypes: []configpb.GrantType{configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.
						ResponseType_RESPONSE_TYPE_CODE, configpb.ResponseType_RESPONSE_TYPE_TOKEN},
					Scopes: []string{"openid", "profile", "email"},
					TokenEndpointAuthMethod: []configpb.TokenEndpointAuthMethod{configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST},
					TokenEndpointAuthSigningAlg: []string{"ES256", "ES384", "ES512"},
					RequestUris:                 make([]string, 0),
					RequestObjectSigningAlg:     "ES256",
					FrontChannelLoginUri:        map[string]string{"default": "http://localhost:3000/login/oauth2"},
					FrontChannelConsentUri:      map[string]string{"default": "http://localhost:3000/consent"},
				},
			}
			beResp := &configpb.UpdateOAuth2ProviderResponse{}

			mockClient.EXPECT().UpdateOAuth2Provider(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateOAuth2Provider(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

	})

	Describe("OAuth2ProviderDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteOAuth2Provider(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("should return an length error in the response", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteOAuth2ProviderRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteOAuth2Provider(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))

		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteOAuth2ProviderRequest{
				Id:   "gid:like-real-oauth2-provider-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteOAuth2ProviderResponse{
				Bookmark: "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().DeleteOAuth2Provider(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteOAuth2Provider(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteOAuth2ProviderRequest{
				Id:   "gid:like-real-oauth2-provider-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteOAuth2ProviderResponse{}

			mockClient.EXPECT().DeleteOAuth2Provider(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteOAuth2Provider(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
