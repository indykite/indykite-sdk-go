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

var _ = Describe("OAuth2Application", func() {
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

	Describe("OAuth2Application", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadOAuth2Application(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadOAuth2ApplicationRequest{Id: "gid:like"}
			resp, err := configClient.ReadOAuth2Application(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("ReadSuccess", func() {
			req := &configpb.ReadOAuth2ApplicationRequest{Id: "gid:like-real-oauth2-application-id"}
			beResp := &configpb.ReadOAuth2ApplicationResponse{
				Oauth2Application: &configpb.OAuth2Application{
					Id:               "gid:like-real-oauth2-application-id",
					Name:             "like-real-oauth2-application-name",
					DisplayName:      "Like Real Oauth2-Application Name",
					CreatedBy:        "creator",
					CreateTime:       timestamppb.Now(),
					CustomerId:       "gid:like-real-customer-id",
					AppSpaceId:       "gid:like-real-app-space-id",
					Etag:             "123qwe",
					Oauth2ProviderId: "gid:like-real-oauth2-provider-id",
					Config: &configpb.OAuth2ApplicationConfig{
						ClientId:    "00000000-90af-4ef9-9928-aaaaaaaaaaaa",
						DisplayName: "Some cool public display name",
						SubjectType: configpb.ClientSubjectType_CLIENT_SUBJECT_TYPE_PUBLIC,
						GrantTypes: []configpb.GrantType{
							configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE,
						},
						ResponseTypes: []configpb.ResponseType{configpb.ResponseType_RESPONSE_TYPE_CODE},
						Scopes:        []string{"openid", "profile", "email"},
						Audiences:     []string{"7d2e906e-541a-49da-b5b2-a28840ff8721"},
						TokenEndpointAuthMethod: configpb.
							TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST,
						TokenEndpointAuthSigningAlg: "ES256",
						RedirectUris:                []string{"http://localhost:3000/redirect"},
						Owner:                       "Owner",
						PolicyUri:                   "http://localhost:3000/policy",
						TermsOfServiceUri:           "http://localhost:3000/policy",
						ClientUri:                   "http://localhost:3000/client",
						LogoUri:                     "http://localhost:3000/logo",
						UserSupportEmailAddress:     "test@example.com",
					},
				},
			}
			mockClient.EXPECT().
				ReadOAuth2Application(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadOAuth2Application(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			req := &configpb.ReadOAuth2ApplicationRequest{Id: "gid:like-real-oauth2-application-id"}
			beResp := &configpb.ReadOAuth2ApplicationResponse{
				Oauth2Application: &configpb.OAuth2Application{
					Id:               "gid:like-real-oauth2-application-id",
					Name:             "like-real-oauth2-application-name",
					Etag:             "123qwe",
					Oauth2ProviderId: "gid:like-real-oauth2-provider-id",
					Config:           &configpb.OAuth2ApplicationConfig{},
				},
			}
			mockClient.EXPECT().
				ReadOAuth2Application(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadOAuth2Application(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("OAuth2ApplicationCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateOAuth2Application(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Application Name"}
			req := &configpb.CreateOAuth2ApplicationRequest{
				Oauth2ProviderId: "gid:like-real-oauth2-provider-id",
				Name:             "like-real-oauth2-application-name",
				DisplayName:      displayNamePb,
				Config: &configpb.OAuth2ApplicationConfig{
					ClientId:    "00000000-90af-4ef9-9928-aaaaaaaaaaaa",
					DisplayName: "Some cool public display name",
					SubjectType: configpb.ClientSubjectType_CLIENT_SUBJECT_TYPE_PUBLIC,
					GrantTypes:  []configpb.GrantType{configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.
						ResponseType_RESPONSE_TYPE_CODE},
					Scopes:    []string{"openid", "profile", "email"},
					Audiences: []string{"7d2e906e-541a-49da-b5b2-a28840ff8721"},
					TokenEndpointAuthMethod: configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST,
					TokenEndpointAuthSigningAlg: "ES256",
					RedirectUris:                []string{"http://localhost:3000/redirect"},
					Owner:                       "Owner",
					PolicyUri:                   "http://localhost:3000/policy",
					TermsOfServiceUri:           "http://localhost:3000/policy",
					ClientUri:                   "http://localhost:3000/client",
					LogoUri:                     "http://localhost:3000/logo",
					UserSupportEmailAddress:     "test@example.com",
				},
			}
			beResp := &configpb.CreateOAuth2ApplicationResponse{
				Id:         "gid:like-real-oauth2-application-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateOAuth2Application(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateOAuth2Application(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Application Name"}
			req := &configpb.CreateOAuth2ApplicationRequest{
				Oauth2ProviderId: "error-id",
				Name:             "like-real-oauth2-application-name",
				DisplayName:      displayNamePb,
				Config:           &configpb.OAuth2ApplicationConfig{},
			}

			resp, err := configClient.CreateOAuth2Application(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Application Name"}
			req := &configpb.CreateOAuth2ApplicationRequest{
				Oauth2ProviderId: "gid:like-real-oauth2-provider-id",
				Name:             "like-real-oauth2-application-name",
				DisplayName:      displayNamePb,
				Config: &configpb.OAuth2ApplicationConfig{
					ClientId:    "00000000-90af-4ef9-9928-aaaaaaaaaaaa",
					DisplayName: "Some cool public display name",
					SubjectType: configpb.ClientSubjectType_CLIENT_SUBJECT_TYPE_PUBLIC,
					GrantTypes: []configpb.GrantType{configpb.
						GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.ResponseType_RESPONSE_TYPE_CODE},
					Scopes:        []string{"openid", "profile", "email"},
					Audiences:     []string{"7d2e906e-541a-49da-b5b2-a28840ff8721"},
					TokenEndpointAuthMethod: configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST,
					TokenEndpointAuthSigningAlg: "ES256",
					RedirectUris:                []string{"http://localhost:3000/redirect"},
					Owner:                       "Owner",
					PolicyUri:                   "http://localhost:3000/policy",
					TermsOfServiceUri:           "http://localhost:3000/policy",
					ClientUri:                   "http://localhost:3000/client",
					LogoUri:                     "http://localhost:3000/logo",
					UserSupportEmailAddress:     "test@example.com",
				},
			}
			beResp := &configpb.CreateOAuth2ApplicationResponse{}
			mockClient.EXPECT().CreateOAuth2Application(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateOAuth2Application(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("OAuth2ApplicationUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateOAuth2Application(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Application Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateOAuth2ApplicationRequest{
				Id:          "gid:like-real-oauth2-application-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
				Config: &configpb.OAuth2ApplicationConfig{
					ClientId:    "00000000-90af-4ef9-9928-aaaaaaaaaaaa",
					DisplayName: "Some cool public display name",
					SubjectType: configpb.ClientSubjectType_CLIENT_SUBJECT_TYPE_PUBLIC,
					GrantTypes: []configpb.GrantType{configpb.
						GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.ResponseType_RESPONSE_TYPE_CODE},
					Scopes:        []string{"openid", "profile", "email"},
					Audiences:     []string{"7d2e906e-541a-49da-b5b2-a28840ff8721"},
					TokenEndpointAuthMethod: configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST,
					TokenEndpointAuthSigningAlg: "ES256",
					RedirectUris:                []string{"http://localhost:3000/redirect"},
					Owner:                       "Owner",
					PolicyUri:                   "http://localhost:3000/policy",
					TermsOfServiceUri:           "http://localhost:3000/policy",
					ClientUri:                   "http://localhost:3000/client",
					LogoUri:                     "http://localhost:3000/logo",
					UserSupportEmailAddress:     "test@example.com",
				},
			}
			beResp := &configpb.UpdateOAuth2ApplicationResponse{
				Id:         "gid:like-real-oauth2-application-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateOAuth2Application(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateOAuth2Application(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Application Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateOAuth2ApplicationRequest{
				Id:          "wrong-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
				Config:      &configpb.OAuth2ApplicationConfig{},
			}
			resp, err := configClient.UpdateOAuth2Application(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real OAuth2Application Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateOAuth2ApplicationRequest{
				Id:          "gid:like-real-oauth2-application-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
				Config: &configpb.OAuth2ApplicationConfig{
					ClientId:    "00000000-90af-4ef9-9928-aaaaaaaaaaaa",
					DisplayName: "Some cool public display name",
					SubjectType: configpb.ClientSubjectType_CLIENT_SUBJECT_TYPE_PUBLIC,
					GrantTypes: []configpb.
						GrantType{configpb.GrantType_GRANT_TYPE_AUTHORIZATION_CODE},
					ResponseTypes: []configpb.ResponseType{configpb.ResponseType_RESPONSE_TYPE_CODE},
					Scopes:        []string{"openid", "profile", "email"},
					Audiences:     []string{"7d2e906e-541a-49da-b5b2-a28840ff8721"},
					TokenEndpointAuthMethod: configpb.
						TokenEndpointAuthMethod_TOKEN_ENDPOINT_AUTH_METHOD_CLIENT_SECRET_POST,
					TokenEndpointAuthSigningAlg: "ES256",
					RedirectUris:                []string{"http://localhost:3000/redirect"},
					Owner:                       "Owner",
					PolicyUri:                   "http://localhost:3000/policy",
					TermsOfServiceUri:           "http://localhost:3000/policy",
					ClientUri:                   "http://localhost:3000/client",
					LogoUri:                     "http://localhost:3000/logo",
					UserSupportEmailAddress:     "test@example.com",
				},
			}
			beResp := &configpb.UpdateOAuth2ApplicationResponse{}

			mockClient.EXPECT().UpdateOAuth2Application(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateOAuth2Application(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("OAuth2ApplicationDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteOAuth2Application(ctx, nil)
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
			req := &configpb.DeleteOAuth2ApplicationRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteOAuth2Application(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteOAuth2ApplicationRequest{
				Id:   "gid:like-real-oauth2-application-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteOAuth2ApplicationResponse{
				Bookmark: "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().DeleteOAuth2Application(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteOAuth2Application(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteOAuth2ApplicationRequest{
				Id:   "gid:like-real-oauth2-application-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteOAuth2ApplicationResponse{}

			mockClient.EXPECT().DeleteOAuth2Application(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteOAuth2Application(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
