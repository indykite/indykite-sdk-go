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

var _ = Describe("ServiceAccount", func() {
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

	Describe("ServiceAccount", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadServiceAccount(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadServiceAccountRequest{
				Identifier: &configpb.ReadServiceAccountRequest_Id{Id: "gid:like"},
			}
			resp, err := configClient.ReadServiceAccount(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		DescribeTable("ReadSuccess",
			func(req *configpb.ReadServiceAccountRequest, beResp *configpb.ReadServiceAccountResponse) {
				mockClient.EXPECT().
					ReadServiceAccount(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, nil)

				resp, err := configClient.ReadServiceAccount(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"ReadId",
				&configpb.ReadServiceAccountRequest{
					Identifier: &configpb.ReadServiceAccountRequest_Id{Id: "gid:like-real-service-account-id"},
				},
				&configpb.ReadServiceAccountResponse{
					ServiceAccount: &configpb.ServiceAccount{
						Id:          "gid:like-real-service-account-id",
						Name:        "like-real-service-account-name",
						DisplayName: "Like Real Service Account Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						CustomerId:  "gid:like-real-customer-id",
						Etag:        "123qwert",
					},
				},
			),
			Entry(
				"ReadName",
				&configpb.ReadServiceAccountRequest{
					Identifier: &configpb.ReadServiceAccountRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-customer-id",
						Name:     "like-real-service-account-name",
					},
					},
				},
				&configpb.ReadServiceAccountResponse{
					ServiceAccount: &configpb.ServiceAccount{
						Id:          "gid:like-real-service-account-id",
						Name:        "like-real-service-account-name",
						DisplayName: "Like Real Service Account Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						CustomerId:  "gid:like-real-customer-id",
						Etag:        "123qwert",
					},
				},
			),
		)

		DescribeTable("ReadError",
			func(req *configpb.ReadServiceAccountRequest, beResp *configpb.ReadServiceAccountResponse) {
				mockClient.EXPECT().
					ReadServiceAccount(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

				resp, err := configClient.ReadServiceAccount(ctx, req)
				Expect(err).ToNot(Succeed())
				Expect(resp).To(BeNil())
			},
			Entry(
				"ReadId",
				&configpb.ReadServiceAccountRequest{
					Identifier: &configpb.ReadServiceAccountRequest_Id{Id: "gid:like-real-service-account-id"},
				},
				&configpb.ReadServiceAccountResponse{},
			),
			Entry(
				"ReadName",
				&configpb.ReadServiceAccountRequest{
					Identifier: &configpb.ReadServiceAccountRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-customer-id",
						Name:     "like-real-service-account-name",
					},
					},
				},
				&configpb.ReadServiceAccountResponse{},
			),
		)
	})

	Describe("ServiceAccountCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateServiceAccount(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Service Account Name"}
			req := &configpb.CreateServiceAccountRequest{
				Location:    "gid:like-real-customer-id",
				Name:        "like-real-service-account-name",
				DisplayName: displayNamePb,
				Role:        "all_editor",
			}
			beResp := &configpb.CreateServiceAccountResponse{
				Id:         "gid:like-real-service-account-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateServiceAccount(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateServiceAccount(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		DescribeTable("CreateNonValid",
			func(mockCreateServiceAccountRequest *configpb.CreateServiceAccountRequest, message string) {
				resp, err := configClient.CreateServiceAccount(context.Background(), mockCreateServiceAccountRequest)
				Ω(resp).To(BeNil())
				Ω(err).To(MatchError(ContainSubstring(message)))
			},
			Entry(
				"CreateNonValidLocation",
				&configpb.CreateServiceAccountRequest{
					Location:    "error-location",
					Name:        "like-real-service-account-name",
					DisplayName: &wrapperspb.StringValue{Value: "Like real Service Account Name"},
				},
				"value length must be between 22 and 254 runes",
			),
			Entry(
				"CreateNonValidRole",
				&configpb.CreateServiceAccountRequest{
					Location:    "gid:like-real-customer-id",
					Name:        "like-real-service-account-name",
					DisplayName: &wrapperspb.StringValue{Value: "Like real Service Account Name"},
					Role:        "none",
				},
				"value must be in list [all_editor all_viewer]",
			),
		)

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Service Account Name"}
			req := &configpb.CreateServiceAccountRequest{
				Location:    "gid:like-real-customer-id",
				Name:        "like-real-service-account-name",
				DisplayName: displayNamePb,
				Role:        "all_editor",
			}
			beResp := &configpb.CreateServiceAccountResponse{}

			mockClient.EXPECT().CreateServiceAccount(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateServiceAccount(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ServiceAccountUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateServiceAccount(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Service Account Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateServiceAccountRequest{
				Id:          "gid:like-real-service-account",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateServiceAccountResponse{
				Id:         "gid:like-real-service-account-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateServiceAccount(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateServiceAccount(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Service Account Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateServiceAccountRequest{
				Id:          "wrong-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			resp, err := configClient.UpdateServiceAccount(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Service Account Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateServiceAccountRequest{
				Id:          "gid:like-real-service-account",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateServiceAccountResponse{}

			mockClient.EXPECT().UpdateServiceAccount(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateServiceAccount(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ServiceAccountDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteServiceAccount(ctx, nil)
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
			req := &configpb.DeleteServiceAccountRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteServiceAccount(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteServiceAccountRequest{
				Id:   "gid:like-real-service-account-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteServiceAccountResponse{}

			mockClient.EXPECT().DeleteServiceAccount(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteServiceAccount(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteServiceAccountRequest{
				Id:   "gid:like-real-service-account-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteServiceAccountResponse{}

			mockClient.EXPECT().DeleteServiceAccount(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteServiceAccount(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
