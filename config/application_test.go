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

var _ = Describe("Application", func() {
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

	Describe("Application", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadApplication(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadApplicationRequest{
				Identifier: &configpb.ReadApplicationRequest_Id{Id: "gid:like"},
			}
			resp, err := configClient.ReadApplication(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		DescribeTable("ReadSuccess",
			func(req *configpb.ReadApplicationRequest, beResp *configpb.ReadApplicationResponse) {
				mockClient.EXPECT().
					ReadApplication(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, nil)

				resp, err := configClient.ReadApplication(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"ReadId",
				&configpb.ReadApplicationRequest{
					Identifier: &configpb.ReadApplicationRequest_Id{Id: "gid:like-real-application-id"},
				},
				&configpb.ReadApplicationResponse{
					Application: &configpb.Application{
						Id:          "gid:like-real-application-id",
						Name:        "like-real-application-name",
						DisplayName: "Like Real Application Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						AppSpaceId:  "gid:like-real-app-space-id",
					},
				},
			),
			Entry(
				"ReadName",
				&configpb.ReadApplicationRequest{
					Identifier: &configpb.ReadApplicationRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-app-space-id",
						Name:     "like-real-application-name",
					},
					},
				},
				&configpb.ReadApplicationResponse{
					Application: &configpb.Application{
						Id:          "gid:like-real-application-id",
						Name:        "like-real-application-name",
						DisplayName: "Like Real Application Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
						AppSpaceId:  "gid:like-real-app-space-id",
					},
				},
			),
		)

		DescribeTable("ReadError",
			func(req *configpb.ReadApplicationRequest, beResp *configpb.ReadApplicationResponse) {
				mockClient.EXPECT().
					ReadApplication(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

				resp, err := configClient.ReadApplication(ctx, req)
				Expect(err).ToNot(Succeed())
				Expect(resp).To(BeNil())
			},
			Entry(
				"ReadId",
				&configpb.ReadApplicationRequest{
					Identifier: &configpb.ReadApplicationRequest_Id{Id: "gid:like-real-application-id"},
				},
				&configpb.ReadApplicationResponse{},
			),
			Entry(
				"ReadName",
				&configpb.ReadApplicationRequest{
					Identifier: &configpb.ReadApplicationRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-app-space-id",
						Name:     "like-real-application-name",
					},
					},
				},
				&configpb.ReadApplicationResponse{},
			),
		)
	})

	Describe("ApplicationCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateApplication(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Name"}
			req := &configpb.CreateApplicationRequest{
				AppSpaceId:  "gid:like-real-app-space-id",
				Name:        "like-real-application-name",
				DisplayName: displayNamePb,
			}
			beResp := &configpb.CreateApplicationResponse{
				Id:         "gid:like-real-application-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateApplication(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateApplication(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Name"}
			req := &configpb.CreateApplicationRequest{
				AppSpaceId:  "error-app-space-id",
				Name:        "like-real-application-name",
				DisplayName: displayNamePb,
			}

			resp, err := configClient.CreateApplication(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Name"}
			req := &configpb.CreateApplicationRequest{
				AppSpaceId:  "gid:like-real-app-space-id",
				Name:        "like-real-application-name",
				DisplayName: displayNamePb,
			}
			beResp := &configpb.CreateApplicationResponse{}

			mockClient.EXPECT().CreateApplication(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateApplication(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ApplicationUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateApplication(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationRequest{
				Id:          "gid:like-real-application-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateApplicationResponse{
				Id:         "gid:like-real-application-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateApplication(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateApplication(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationRequest{
				Id:          "wrong-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			resp, err := configClient.UpdateApplication(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationRequest{
				Id:          "gid:like-real-application-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateApplicationResponse{}

			mockClient.EXPECT().UpdateApplication(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateApplication(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ApplicationList", func() {
		It("Nil request", func() {
			resp, err := configClient.ListApplications(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("MatchName", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListApplicationsClient(mockCtrl)
			mockListApplicationsRequest := configpb.ListApplicationsRequest{
				AppSpaceId: "gid:like-real-app-space-id",
				Match:      []string{"like-real-application-name"},
			}

			mockClient.EXPECT().ListApplications(
				gomock.Any(),
				gomock.Eq(&mockListApplicationsRequest),
				gomock.Any(),
			).Return(mockResponseClient, nil)

			stream, err := configClient.ListApplications(context.Background(), &mockListApplicationsRequest)
			Ω(stream).ToNot(BeNil())
			Ω(err).To(Succeed())

			mockResp := &configpb.ListApplicationsResponse{
				Application: &configpb.Application{
					Id:          "gid:like-real-application-id",
					Name:        "like-real-application-name",
					DisplayName: "Like Real Application Name",
					CreatedBy:   "creator",
					CreateTime:  timestamppb.Now(),
					Etag:        "123qwert",
					AppSpaceId:  "gid:like-real-app-space-id",
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
			func(mockListApplicationsRequest *configpb.ListApplicationsRequest, message string) {
				resp, err := configClient.ListApplications(context.Background(), mockListApplicationsRequest)
				Ω(resp).To(BeNil())
				Ω(err).To(MatchError(ContainSubstring(message)))
			},
			Entry(
				"AppSpace error",
				&configpb.ListApplicationsRequest{
					AppSpaceId: "app-space-id",
					Match:      []string{"like-real-application-name"},
				},
				"Id: value length must be between 22",
			),
			Entry(
				"Match error",
				&configpb.ListApplicationsRequest{
					AppSpaceId: "gid:like-real-app-space-id",
					Match:      []string{},
				},
				"Match: value must contain at least 1 item",
			),
		)

		It("MatchNameError", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListApplicationsClient(mockCtrl)
			mockListApplicationsRequest := configpb.ListApplicationsRequest{
				AppSpaceId: "gid:like-real-app-space-id",
				Match:      []string{"like-real-application-name"},
			}

			mockClient.EXPECT().ListApplications(
				gomock.Any(),
				gomock.Eq(&mockListApplicationsRequest),
				gomock.Any(),
			).Return(mockResponseClient, status.Error(codes.InvalidArgument, "status error"))

			stream, err := configClient.ListApplications(context.Background(), &mockListApplicationsRequest)
			Expect(err).ToNot(Succeed())
			Expect(stream).To(BeNil())
		})
	})

	Describe("ApplicationDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteApplication(ctx, nil)
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
			req := &configpb.DeleteApplicationRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteApplication(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationRequest{
				Id:   "gid:like-real-application-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationResponse{
				Bookmark: "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().DeleteApplication(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteApplication(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationRequest{
				Id:   "gid:like-real-application-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationResponse{}

			mockClient.EXPECT().DeleteApplication(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteApplication(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
