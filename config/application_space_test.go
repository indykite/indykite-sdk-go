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

var _ = Describe("AppSpace", func() {
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

	Describe("AppSpace", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadApplicationSpace(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadApplicationSpaceRequest{
				Identifier: &configpb.ReadApplicationSpaceRequest_Id{Id: "gid:like"},
			}
			resp, err := configClient.ReadApplicationSpace(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		DescribeTable("ReadSuccess",
			func(req *configpb.ReadApplicationSpaceRequest, beResp *configpb.ReadApplicationSpaceResponse) {
				mockClient.EXPECT().
					ReadApplicationSpace(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, nil)

				resp, err := configClient.ReadApplicationSpace(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"ReadId",
				&configpb.ReadApplicationSpaceRequest{
					Identifier: &configpb.ReadApplicationSpaceRequest_Id{Id: "gid:like-real-appspace-id"},
				},
				&configpb.ReadApplicationSpaceResponse{
					AppSpace: &configpb.ApplicationSpace{
						Id:          "gid:like-real-appspace-id",
						Name:        "like-real-appspace-name",
						DisplayName: "Like Real Appspace Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						Region:      "europe-west1",
					},
				},
			),
			Entry(
				"ReadName",
				&configpb.ReadApplicationSpaceRequest{
					Identifier: &configpb.ReadApplicationSpaceRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-customer-id",
						Name:     "like-real-appspace-name",
					},
					},
				},
				&configpb.ReadApplicationSpaceResponse{
					AppSpace: &configpb.ApplicationSpace{
						Id:          "gid:like-real-appspace-id",
						Name:        "like-real-appspace-name",
						DisplayName: "Like Real Appspace Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						Region:      "europe-west1",
					},
				},
			),
		)
	})

	DescribeTable("ReadError",
		func(req *configpb.ReadApplicationSpaceRequest, beResp *configpb.ReadApplicationSpaceResponse) {
			mockClient.EXPECT().
				ReadApplicationSpace(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadApplicationSpace(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		},
		Entry(
			"ReadId",
			&configpb.ReadApplicationSpaceRequest{
				Identifier: &configpb.ReadApplicationSpaceRequest_Id{Id: "gid:like-real-appspace-id"},
			},
			&configpb.ReadApplicationSpaceResponse{},
		),
		Entry(
			"ReadName",
			&configpb.ReadApplicationSpaceRequest{
				Identifier: &configpb.ReadApplicationSpaceRequest_Name{Name: &configpb.UniqueNameIdentifier{
					Location: "gid:like-real-customer-id",
					Name:     "like-real-appspace-name",
				},
				},
			},
			&configpb.ReadApplicationSpaceResponse{
				AppSpace: &configpb.ApplicationSpace{},
			},
		),
	)

	Describe("AppSpaceCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateApplicationSpace(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real AppSpace Name"}
			req := &configpb.CreateApplicationSpaceRequest{
				CustomerId:  "gid:like-real-customer-id",
				Name:        "like-real-appspace-name",
				DisplayName: displayNamePb,
				Region:      "europe-west1",
			}
			beResp := &configpb.CreateApplicationSpaceResponse{
				Id:         "gid:like-real-app_space-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().CreateApplicationSpace(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateApplicationSpace(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real ApplicationSpace Name"}
			req := &configpb.CreateApplicationSpaceRequest{
				CustomerId:  "error-customer-id",
				Name:        "like-real-app-space-name",
				DisplayName: displayNamePb,
			}

			resp, err := configClient.CreateApplicationSpace(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real AppSpace Name"}
			req := &configpb.CreateApplicationSpaceRequest{
				CustomerId:  "gid:like-real-customer-id",
				Name:        "like-real-appspace-name",
				DisplayName: displayNamePb,
			}
			beResp := &configpb.CreateApplicationSpaceResponse{}

			mockClient.EXPECT().CreateApplicationSpace(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateApplicationSpace(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("AppSpaceUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateApplicationSpace(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real AppSpace Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationSpaceRequest{
				Id:          "gid:like-real-app-space-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateApplicationSpaceResponse{
				Id:         "gid:like-real-app_space-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
			}

			mockClient.EXPECT().UpdateApplicationSpace(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateApplicationSpace(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real ApplicationSpace Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationSpaceRequest{
				Id:          "wrong-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			resp, err := configClient.UpdateApplicationSpace(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real AppSpace Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationSpaceRequest{
				Id:          "gid:like-real-app-space-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateApplicationSpaceResponse{}

			mockClient.EXPECT().UpdateApplicationSpace(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateApplicationSpace(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("AppSpaceList", func() {
		It("Nil request", func() {
			resp, err := configClient.ListApplicationSpaces(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("MatchName", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListApplicationSpacesClient(mockCtrl)
			mockListApplicationSpacesRequest := configpb.ListApplicationSpacesRequest{
				CustomerId: "gid:like-real-customer-id",
				Match:      []string{"like-real-app-space-name"},
			}

			mockClient.EXPECT().ListApplicationSpaces(
				gomock.Any(),
				gomock.Eq(&mockListApplicationSpacesRequest),
				gomock.Any(),
			).Return(mockResponseClient, nil)

			stream, err := configClient.ListApplicationSpaces(context.Background(), &mockListApplicationSpacesRequest)
			Ω(stream).ToNot(BeNil())
			Ω(err).To(Succeed())

			mockResp := &configpb.ListApplicationSpacesResponse{
				AppSpace: &configpb.ApplicationSpace{
					Id:          "gid:like-real-appspace-id",
					Name:        "like-real-appspace-name",
					DisplayName: "Like Real Appspace Name",
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

		DescribeTable("ListError",
			func(mockListApplicationSpacesRequest *configpb.ListApplicationSpacesRequest, message string) {
				resp, err := configClient.ListApplicationSpaces(context.Background(), mockListApplicationSpacesRequest)
				Ω(resp).To(BeNil())
				Ω(err).To(MatchError(ContainSubstring(message)))
			},
			Entry(
				"Customer error",
				&configpb.ListApplicationSpacesRequest{
					CustomerId: "customer-id",
					Match:      []string{"like-real-appspace-name"},
				},
				"Id: value length must be between 22",
			),
			Entry(
				"Match error",
				&configpb.ListApplicationSpacesRequest{
					CustomerId: "gid:like-real-customer-id",
					Match:      []string{},
				},
				"Match: value must contain at least 1 item",
			),
		)

		It("MatchNameError", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListApplicationSpacesClient(mockCtrl)
			mockListApplicationSpacesRequest := configpb.ListApplicationSpacesRequest{
				CustomerId: "gid:like-real-customer-id",
				Match:      []string{"like-real-app-space-name"},
			}

			mockClient.EXPECT().ListApplicationSpaces(
				gomock.Any(),
				gomock.Eq(&mockListApplicationSpacesRequest),
				gomock.Any(),
			).Return(mockResponseClient, status.Error(codes.InvalidArgument, "status error"))

			stream, err := configClient.ListApplicationSpaces(context.Background(), &mockListApplicationSpacesRequest)
			Expect(err).ToNot(Succeed())
			Expect(stream).To(BeNil())
		})
	})

	Describe("AppSpaceDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteApplicationSpace(ctx, nil)
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
			req := &configpb.DeleteApplicationSpaceRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteApplicationSpace(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationSpaceRequest{
				Id:   "gid:like-real-app-space-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationSpaceResponse{}

			mockClient.EXPECT().DeleteApplicationSpace(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteApplicationSpace(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationSpaceRequest{
				Id:   "gid:like-real-app-space-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationSpaceResponse{}

			mockClient.EXPECT().DeleteApplicationSpace(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteApplicationSpace(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
