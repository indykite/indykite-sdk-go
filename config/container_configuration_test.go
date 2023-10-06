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

var _ = Describe("Container Configuration", func() {
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

	Describe("Customer Configuration", func() {
		DescribeTable("Request error during Read ",
			func(req *configpb.ReadCustomerConfigRequest, causeMatcher, msgMatcher, codeMatcher OmegaMatcher) {
				resp, err := configClient.ReadCustomerConfig(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				var clientErr *sdkerrors.ClientError
				Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
				Expect(clientErr.Unwrap()).To(causeMatcher)
				Expect(clientErr.Message()).To(msgMatcher)
				Expect(clientErr.Code()).To(codeMatcher)
			},
			Entry("nil request", nil, BeNil(), Equal("invalid nil request"), Equal(codes.InvalidArgument)),
			Entry("invalid request", &configpb.ReadCustomerConfigRequest{},
				MatchError(ContainSubstring("Id: value length must be between 22")),
				Equal("invalid request"),
				Equal(codes.InvalidArgument)),
		)

		It("Read", func() {
			req := &configpb.ReadCustomerConfigRequest{
				Id:        "gid:like-real-customer-gid",
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.ReadCustomerConfigResponse{
				Id:   "gid:like-real-customer-gid",
				Etag: "etag-value",
				Config: &configpb.CustomerConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
				},
			}
			mockClient.EXPECT().
				ReadCustomerConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, nil)

			resp, err := configClient.ReadCustomerConfig(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			req := &configpb.ReadCustomerConfigRequest{
				Id:        "gid:like-real-customer-gid",
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.ReadCustomerConfigResponse{}
			mockClient.EXPECT().
				ReadCustomerConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadCustomerConfig(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		DescribeTable("Request error during Update",
			func(req *configpb.UpdateCustomerConfigRequest, causeMatcher, msgMatcher, codeMatcher OmegaMatcher) {
				resp, err := configClient.UpdateCustomerConfig(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				var clientErr *sdkerrors.ClientError
				Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
				Expect(clientErr.Unwrap()).To(causeMatcher)
				Expect(clientErr.Message()).To(msgMatcher)
				Expect(clientErr.Code()).To(codeMatcher)
			},
			Entry("nil request", nil, BeNil(), Equal("invalid nil request"), Equal(codes.InvalidArgument)),
			Entry("invalid request", &configpb.UpdateCustomerConfigRequest{},
				MatchError(ContainSubstring("Id: value length must be between 22")),
				Equal("invalid request"),
				Equal(codes.InvalidArgument)),
		)

		It("Update", func() {
			req := &configpb.UpdateCustomerConfigRequest{
				Id:   "gid:like-real-customer-gid",
				Etag: wrapperspb.String("etag-value"),
				Config: &configpb.CustomerConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow-created-under-customer",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
				},
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.UpdateCustomerConfigResponse{
				Id:         "gid:like-real-customer-gid",
				CreateTime: timestamppb.Now(),
				CreatedBy:  "gid:id-of-user-or-service-account",
				UpdateTime: timestamppb.Now(),
				UpdatedBy:  "gid:id-of-user-or-service-account",
				Bookmark:   "something-like-old-bookmark-which-is-long-enough",
				Etag:       "etag-value",
			}
			mockClient.EXPECT().
				UpdateCustomerConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, nil)

			resp, err := configClient.UpdateCustomerConfig(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateError", func() {
			req := &configpb.UpdateCustomerConfigRequest{
				Id:   "gid:like-real-customer-gid",
				Etag: wrapperspb.String("etag-value"),
				Config: &configpb.CustomerConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow-created-under-customer",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
				},
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.UpdateCustomerConfigResponse{}
			mockClient.EXPECT().
				UpdateCustomerConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateCustomerConfig(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ApplicationSpace Configuration", func() {
		DescribeTable("Request error during Read ",
			func(req *configpb.ReadApplicationSpaceConfigRequest, causeMatcher, msgMatcher, codeMatcher OmegaMatcher) {
				resp, err := configClient.ReadApplicationSpaceConfig(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				var clientErr *sdkerrors.ClientError
				Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
				Expect(clientErr.Unwrap()).To(causeMatcher)
				Expect(clientErr.Message()).To(msgMatcher)
				Expect(clientErr.Code()).To(codeMatcher)
			},
			Entry("nil request", nil, BeNil(), Equal("invalid nil request"), Equal(codes.InvalidArgument)),
			Entry("invalid request", &configpb.ReadApplicationSpaceConfigRequest{},
				MatchError(ContainSubstring("Id: value length must be between 22")),
				Equal("invalid request"),
				Equal(codes.InvalidArgument)),
		)

		It("Read", func() {
			req := &configpb.ReadApplicationSpaceConfigRequest{
				Id:        "gid:like-real-ApplicationSpace-gid",
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.ReadApplicationSpaceConfigResponse{
				Id:   "gid:like-real-ApplicationSpace-gid",
				Etag: "etag-value",
				Config: &configpb.ApplicationSpaceConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
					DefaultTenantId:       "gid:id-of-tenant-under-this-app-space",
					UsernamePolicy: &configpb.UsernamePolicy{
						AllowedUsernameFormats: []string{"email"},
						ValidEmail:             true,
						VerifyEmail:            true,
						AllowedEmailDomains:    []string{"gmail.com"},
					},
				},
			}
			mockClient.EXPECT().
				ReadApplicationSpaceConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, nil)

			resp, err := configClient.ReadApplicationSpaceConfig(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			req := &configpb.ReadApplicationSpaceConfigRequest{
				Id:        "gid:like-real-ApplicationSpace-gid",
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.ReadApplicationSpaceConfigResponse{}
			mockClient.EXPECT().
				ReadApplicationSpaceConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadApplicationSpaceConfig(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		DescribeTable("Request error during Update",
			func(req *configpb.UpdateApplicationSpaceConfigRequest, causeMatch, msgMatcher, codeMatcher OmegaMatcher) {
				resp, err := configClient.UpdateApplicationSpaceConfig(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				var clientErr *sdkerrors.ClientError
				Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
				Expect(clientErr.Unwrap()).To(causeMatch)
				Expect(clientErr.Message()).To(msgMatcher)
				Expect(clientErr.Code()).To(codeMatcher)
			},
			Entry("nil request", nil, BeNil(), Equal("invalid nil request"), Equal(codes.InvalidArgument)),
			Entry("invalid request", &configpb.UpdateApplicationSpaceConfigRequest{},
				MatchError(ContainSubstring("Id: value length must be between 22")),
				Equal("invalid request"),
				Equal(codes.InvalidArgument)),
		)

		It("Update", func() {
			req := &configpb.UpdateApplicationSpaceConfigRequest{
				Id:   "gid:like-real-ApplicationSpace-gid",
				Etag: wrapperspb.String("etag-value"),
				Config: &configpb.ApplicationSpaceConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow-created-under-ApplicationSpace",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
					DefaultTenantId:       "gid:id-of-tenant-under-this-app-space",
					UsernamePolicy: &configpb.UsernamePolicy{
						AllowedUsernameFormats: []string{"email"},
						ValidEmail:             true,
						VerifyEmail:            true,
						AllowedEmailDomains:    []string{"gmail.com"},
					},
				},
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.UpdateApplicationSpaceConfigResponse{
				Id:         "gid:like-real-ApplicationSpace-gid",
				CreateTime: timestamppb.Now(),
				CreatedBy:  "gid:id-of-user-or-service-account",
				UpdateTime: timestamppb.Now(),
				UpdatedBy:  "gid:id-of-user-or-service-account",
				Bookmark:   "something-like-old-bookmark-which-is-long-enough",
				Etag:       "etag-value",
			}
			mockClient.EXPECT().
				UpdateApplicationSpaceConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, nil)

			resp, err := configClient.UpdateApplicationSpaceConfig(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateError", func() {
			req := &configpb.UpdateApplicationSpaceConfigRequest{
				Id:   "gid:like-real-ApplicationSpace-gid",
				Etag: wrapperspb.String("etag-value"),
				Config: &configpb.ApplicationSpaceConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow-created-under-ApplicationSpace",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
					DefaultTenantId:       "gid:id-of-tenant-under-this-app-space",
					UsernamePolicy: &configpb.UsernamePolicy{
						AllowedUsernameFormats: []string{"email"},
						ValidEmail:             true,
						VerifyEmail:            true,
						AllowedEmailDomains:    []string{"gmail.com"},
					},
				},
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.UpdateApplicationSpaceConfigResponse{}
			mockClient.EXPECT().
				UpdateApplicationSpaceConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateApplicationSpaceConfig(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("Tenant Configuration", func() {
		DescribeTable("Request error during Read ",
			func(req *configpb.ReadTenantConfigRequest, causeMatcher, msgMatcher, codeMatcher OmegaMatcher) {
				resp, err := configClient.ReadTenantConfig(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				var clientErr *sdkerrors.ClientError
				Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
				Expect(clientErr.Unwrap()).To(causeMatcher)
				Expect(clientErr.Message()).To(msgMatcher)
				Expect(clientErr.Code()).To(codeMatcher)
			},
			Entry("nil request", nil, BeNil(), Equal("invalid nil request"), Equal(codes.InvalidArgument)),
			Entry("invalid request", &configpb.ReadTenantConfigRequest{},
				MatchError(ContainSubstring("Id: value length must be between 22")),
				Equal("invalid request"),
				Equal(codes.InvalidArgument)),
		)

		It("Read", func() {
			req := &configpb.ReadTenantConfigRequest{
				Id:        "gid:like-real-Tenant-gid",
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.ReadTenantConfigResponse{
				Id:   "gid:like-real-Tenant-gid",
				Etag: "etag-value",
				Config: &configpb.TenantConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
				},
			}
			mockClient.EXPECT().
				ReadTenantConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, nil)

			resp, err := configClient.ReadTenantConfig(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			req := &configpb.ReadTenantConfigRequest{
				Id:        "gid:like-real-Tenant-gid",
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.ReadTenantConfigResponse{}
			mockClient.EXPECT().
				ReadTenantConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadTenantConfig(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		DescribeTable("Request error during Update",
			func(req *configpb.UpdateTenantConfigRequest, causeMatcher, msgMatcher, codeMatcher OmegaMatcher) {
				resp, err := configClient.UpdateTenantConfig(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				var clientErr *sdkerrors.ClientError
				Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
				Expect(clientErr.Unwrap()).To(causeMatcher)
				Expect(clientErr.Message()).To(msgMatcher)
				Expect(clientErr.Code()).To(codeMatcher)
			},
			Entry("nil request", nil, BeNil(), Equal("invalid nil request"), Equal(codes.InvalidArgument)),
			Entry("invalid request", &configpb.UpdateTenantConfigRequest{},
				MatchError(ContainSubstring("Id: value length must be between 22")),
				Equal("invalid request"),
				Equal(codes.InvalidArgument)),
		)

		It("Update", func() {
			req := &configpb.UpdateTenantConfigRequest{
				Id:   "gid:like-real-Tenant-gid",
				Etag: wrapperspb.String("etag-value"),
				Config: &configpb.TenantConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow-created-under-Tenant",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
				},
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.UpdateTenantConfigResponse{
				Id:         "gid:like-real-Tenant-gid",
				CreateTime: timestamppb.Now(),
				CreatedBy:  "gid:id-of-user-or-service-account",
				UpdateTime: timestamppb.Now(),
				UpdatedBy:  "gid:id-of-user-or-service-account",
				Bookmark:   "something-like-old-bookmark-which-is-long-enough",
				Etag:       "etag-value",
			}
			mockClient.EXPECT().
				UpdateTenantConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, nil)

			resp, err := configClient.UpdateTenantConfig(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateError", func() {
			req := &configpb.UpdateTenantConfigRequest{
				Id:   "gid:like-real-Tenant-gid",
				Etag: wrapperspb.String("etag-value"),
				Config: &configpb.TenantConfig{
					DefaultAuthFlowId:     "gid:id-of-authflow-created-under-Tenant",
					DefaultEmailServiceId: "gid:id-of-email-service-provider",
				},
				Bookmarks: []string{"something-like-bookmark-which-is-long-enough"},
			}
			beResp := &configpb.UpdateTenantConfigResponse{}
			mockClient.EXPECT().
				UpdateTenantConfig(gomock.Any(), test.WrapMatcher(test.EqualProto(req)), gomock.Any()).
				Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateTenantConfig(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
