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
	"io"

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

var _ = Describe("Tenant", func() {
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

	Describe("Tenant", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadTenant(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadTenantRequest{
				Identifier: &configpb.ReadTenantRequest_Id{Id: "gid:like"},
			}
			resp, err := configClient.ReadTenant(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		DescribeTable("ReadSuccess",
			func(req *configpb.ReadTenantRequest, beResp *configpb.ReadTenantResponse) {
				mockClient.EXPECT().
					ReadTenant(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, nil)

				resp, err := configClient.ReadTenant(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"ReadId",
				&configpb.ReadTenantRequest{
					Identifier: &configpb.ReadTenantRequest_Id{Id: "gid:like-real-tenant-id"},
				},
				&configpb.ReadTenantResponse{
					Tenant: &configpb.Tenant{
						Id:          "gid:like-real-tenant-id",
						Name:        "like-real-tenant-name",
						DisplayName: "Like Real Tenant Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
					},
				},
			),
			Entry(
				"ReadName",
				&configpb.ReadTenantRequest{
					Identifier: &configpb.ReadTenantRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-app-space-id",
						Name:     "like-real-tenant-name",
					},
					},
				},
				&configpb.ReadTenantResponse{
					Tenant: &configpb.Tenant{
						Id:          "gid:like-real-tenant-id",
						Name:        "like-real-tenant-name",
						DisplayName: "Like Real Tenant Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
					},
				},
			),
		)

		DescribeTable("ReadError",
			func(req *configpb.ReadTenantRequest, beResp *configpb.ReadTenantResponse) {
				mockClient.EXPECT().
					ReadTenant(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

				resp, err := configClient.ReadTenant(ctx, req)
				Expect(err).ToNot(Succeed())
				Expect(resp).To(BeNil())
			},
			Entry(
				"ReadId",
				&configpb.ReadTenantRequest{
					Identifier: &configpb.ReadTenantRequest_Id{Id: "gid:like-real-tenant-id"},
				},
				&configpb.ReadTenantResponse{},
			),
			Entry(
				"ReadName",
				&configpb.ReadTenantRequest{
					Identifier: &configpb.ReadTenantRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-app-space-id",
						Name:     "like-real-tenant-name",
					},
					},
				},
				&configpb.ReadTenantResponse{},
			),
		)
	})

	Describe("TenantCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateTenant(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Tenant Name"}
			req := &configpb.CreateTenantRequest{
				IssuerId:    "gid:like-real-issuer-id",
				Name:        "like-real-tenant-name",
				DisplayName: displayNamePb,
			}
			beResp := &configpb.CreateTenantResponse{
				Id:         "gid:like-real-tenant-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateTenant(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateTenant(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Tenant Name"}
			req := &configpb.CreateTenantRequest{
				IssuerId:    "error-issuer-id",
				Name:        "like-real-tenant-name",
				DisplayName: displayNamePb,
			}

			resp, err := configClient.CreateTenant(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Tenant Name"}
			req := &configpb.CreateTenantRequest{
				IssuerId:    "gid:like-real-issuer-id",
				Name:        "like-real-tenant-name",
				DisplayName: displayNamePb,
			}
			beResp := &configpb.CreateTenantResponse{}

			mockClient.EXPECT().CreateTenant(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateTenant(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
	Describe("TenantUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateTenant(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Tenant Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateTenantRequest{
				Id:          "gid:like-real-tenant-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateTenantResponse{
				Id:         "gid:like-real-tenant-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateTenant(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateTenant(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Tenant Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateTenantRequest{
				Id:          "non-valid-tenant-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			resp, err := configClient.UpdateTenant(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Tenant Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateTenantRequest{
				Id:          "gid:like-real-tenant-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateTenantResponse{}

			mockClient.EXPECT().UpdateTenant(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateTenant(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

	})
	Describe("TenantList", func() {
		It("Nil request", func() {
			resp, err := configClient.ListTenants(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("MatchName", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListTenantsClient(mockCtrl)
			mockListTenantsRequest := configpb.ListTenantsRequest{
				AppSpaceId: "gid:like-real-app-space-id",
				Match:      []string{"like-real-tenant-name"},
			}

			mockClient.EXPECT().ListTenants(
				gomock.Any(),
				gomock.Eq(&mockListTenantsRequest),
				gomock.Any(),
			).Return(mockResponseClient, nil)

			stream, err := configClient.ListTenants(context.Background(), &mockListTenantsRequest)
			Ω(stream).ToNot(BeNil())
			Ω(err).To(Succeed())

			mockResp := &configpb.ListTenantsResponse{
				Tenant: &configpb.Tenant{
					Id:          "gid:like-tenant-id",
					Name:        "like-real-tenant-name",
					DisplayName: "Like Real Tenant Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					Etag:        "123qwert",
				},
			}

			mockResponseClient.EXPECT().Recv().Return(mockResp, nil)

			mockResponseClient.EXPECT().Recv().Return(nil, io.EOF)

			resp, err := stream.Recv()
			Ω(err).To(Succeed())
			Ω(resp).To(test.EqualProto(mockResp))

			resp, err = stream.Recv()
			Ω(err).To(Equal(io.EOF))
			Ω(resp).To(BeNil())

		})

		DescribeTable("ListNonValid",
			func(mockListTenantsRequest *configpb.ListTenantsRequest, message string) {
				resp, err := configClient.ListTenants(context.Background(), mockListTenantsRequest)
				Ω(resp).To(BeNil())
				Ω(err).To(MatchError(ContainSubstring(message)))
			},
			Entry(
				"ApplicationSpace error",
				&configpb.ListTenantsRequest{
					AppSpaceId: "app-space-id",
					Match:      []string{"like-real-tenant-name"},
				},
				"Id: value length must be between 22",
			),
			Entry(
				"Match error",
				&configpb.ListTenantsRequest{
					AppSpaceId: "gid:like-real-app-space-id",
					Match:      []string{},
				},
				"Match: value must contain at least 1 item",
			),
		)

		It("MatchNameError", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListTenantsClient(mockCtrl)
			mockListTenantsRequest := configpb.ListTenantsRequest{
				AppSpaceId: "gid:like-real-app-space-id",
				Match:      []string{"like-real-tenant-name"},
			}

			mockClient.EXPECT().ListTenants(
				gomock.Any(),
				gomock.Eq(&mockListTenantsRequest),
				gomock.Any(),
			).Return(mockResponseClient, status.Error(codes.InvalidArgument, "status error"))

			stream, err := configClient.ListTenants(context.Background(), &mockListTenantsRequest)
			Expect(err).ToNot(Succeed())
			Expect(stream).To(BeNil())
		})
	})

	Describe("TenantDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteTenant(ctx, nil)
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
			req := &configpb.DeleteTenantRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteTenant(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))

		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteTenantRequest{
				Id:   "gid:like-real-tenant-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteTenantResponse{
				Bookmark: "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().DeleteTenant(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteTenant(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteTenantRequest{
				Id:   "gid:like-real-tenant-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteTenantResponse{}

			mockClient.EXPECT().DeleteTenant(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteTenant(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
