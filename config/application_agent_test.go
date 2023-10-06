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

var _ = Describe("ApplicationAgent", func() {
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

	Describe("ApplicationAgent", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadApplicationAgent(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadApplicationAgentRequest{
				Identifier: &configpb.ReadApplicationAgentRequest_Id{Id: "gid:like"},
			}
			resp, err := configClient.ReadApplicationAgent(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		DescribeTable("ReadSuccess",
			func(req *configpb.ReadApplicationAgentRequest, beResp *configpb.ReadApplicationAgentResponse) {
				mockClient.EXPECT().
					ReadApplicationAgent(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, nil)

				resp, err := configClient.ReadApplicationAgent(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"ReadId",
				&configpb.ReadApplicationAgentRequest{
					Identifier: &configpb.ReadApplicationAgentRequest_Id{Id: "gid:like-real-application-agent-id"},
				},
				&configpb.ReadApplicationAgentResponse{
					ApplicationAgent: &configpb.ApplicationAgent{
						Id:            "gid:like-real-application-agent-id",
						Name:          "like-real-application-agent-name",
						DisplayName:   "Like Real Application Agent Name",
						CreatedBy:     "creator",
						CreateTime:    timestamppb.Now(),
						ApplicationId: "gid:like-real-application-id",
					},
				},
			),
			Entry(
				"ReadName",
				&configpb.ReadApplicationAgentRequest{
					Identifier: &configpb.ReadApplicationAgentRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-application-id",
						Name:     "like-real-application-agent-name",
					},
					},
				},
				&configpb.ReadApplicationAgentResponse{
					ApplicationAgent: &configpb.ApplicationAgent{
						Id:            "gid:like-real-application-agent-id",
						Name:          "like-real-application-agent-name",
						DisplayName:   "Like Real Application Agent Name",
						CreatedBy:     "creator",
						CreateTime:    timestamppb.Now(),
						ApplicationId: "gid:like-real-application-id",
					},
				},
			),
		)

		DescribeTable("ReadError",
			func(req *configpb.ReadApplicationAgentRequest, beResp *configpb.ReadApplicationAgentResponse) {
				mockClient.EXPECT().
					ReadApplicationAgent(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

				resp, err := configClient.ReadApplicationAgent(ctx, req)
				Expect(err).ToNot(Succeed())
				Expect(resp).To(BeNil())
			},
			Entry(
				"ReadId",
				&configpb.ReadApplicationAgentRequest{
					Identifier: &configpb.ReadApplicationAgentRequest_Id{Id: "gid:like-real-application-agent-id"},
				},
				&configpb.ReadApplicationAgentResponse{},
			),
			Entry(
				"ReadName",
				&configpb.ReadApplicationAgentRequest{
					Identifier: &configpb.ReadApplicationAgentRequest_Name{Name: &configpb.UniqueNameIdentifier{
						Location: "gid:like-real-application-id",
						Name:     "like-real-application-agent-name",
					},
					},
				},
				&configpb.ReadApplicationAgentResponse{},
			),
		)
	})

	Describe("ApplicationAgentCreate", func() {
		It("Nil request", func() {
			resp, err := configClient.CreateApplicationAgent(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Create", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Agent Name"}
			req := &configpb.CreateApplicationAgentRequest{
				ApplicationId: "gid:like-real-application-id",
				Name:          "like-real-application-agent-name",
				DisplayName:   displayNamePb,
			}
			beResp := &configpb.CreateApplicationAgentResponse{
				Id:         "gid:like-real-application-agent-id",
				Etag:       "123qwe",
				CreatedBy:  "creator",
				CreateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().CreateApplicationAgent(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.CreateApplicationAgent(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("CreateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Agent Name"}
			req := &configpb.CreateApplicationAgentRequest{
				ApplicationId: "gid:like-real-application-id",
				Name:          "like-real-application-agent-name",
				DisplayName:   displayNamePb,
			}
			beResp := &configpb.CreateApplicationAgentResponse{}

			mockClient.EXPECT().CreateApplicationAgent(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.CreateApplicationAgent(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("CreateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real ApplicationAgent Name"}
			req := &configpb.CreateApplicationAgentRequest{
				ApplicationId: "error-app-id",
				Name:          "like-real-application-agent-name",
				DisplayName:   displayNamePb,
			}

			resp, err := configClient.CreateApplicationAgent(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})
	})

	Describe("ApplicationAgentUpdate", func() {
		It("Nil request", func() {
			resp, err := configClient.UpdateApplicationAgent(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Update", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real ApplicationAgent Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationAgentRequest{
				Id:          "gid:like-real-application-agent-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateApplicationAgentResponse{
				Id:         "gid:like-real-application-agent-id",
				Etag:       "123qwert",
				UpdatedBy:  "creator",
				UpdateTime: timestamppb.Now(),
				Bookmark:   "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().UpdateApplicationAgent(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.UpdateApplicationAgent(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("UpdateError", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real ApplicationAgent Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationAgentRequest{
				Id:          "gid:like-real-application-agent-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			beResp := &configpb.UpdateApplicationAgentResponse{}

			mockClient.EXPECT().UpdateApplicationAgent(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.UpdateApplicationAgent(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("UpdateNonValid", func() {
			displayNamePb := &wrapperspb.StringValue{Value: "Like real Application Agent Name"}
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.UpdateApplicationAgentRequest{
				Id:          "wrong-id",
				Etag:        etagPb,
				DisplayName: displayNamePb,
			}
			resp, err := configClient.UpdateApplicationAgent(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})
	})

	Describe("ApplicationAgentList", func() {
		It("Nil request", func() {
			resp, err := configClient.ListApplicationAgents(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("MatchName", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListApplicationAgentsClient(mockCtrl)
			mockListApplicationAgentsRequest := configpb.ListApplicationAgentsRequest{
				AppSpaceId: "gid:like-real-app-space-id",
				Match:      []string{"like-real-application-agent-name"},
			}

			mockClient.EXPECT().ListApplicationAgents(
				gomock.Any(),
				gomock.Eq(&mockListApplicationAgentsRequest),
				gomock.Any(),
			).Return(mockResponseClient, nil)

			stream, err := configClient.ListApplicationAgents(context.Background(), &mockListApplicationAgentsRequest)
			Ω(stream).ToNot(BeNil())
			Ω(err).To(Succeed())

			mockResp := &configpb.ListApplicationAgentsResponse{
				ApplicationAgent: &configpb.ApplicationAgent{
					Id:            "gid:like-real-applicationagent-id",
					Name:          "like-real-applicationagent-name",
					DisplayName:   "Like Real ApplicationAgent Name",
					CreatedBy:     "creator",
					CreateTime:    timestamppb.Now(),
					Etag:          "123qwert",
					ApplicationId: "gid:like-real-application-id",
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

		It("MatchNameError", func() {
			mockResponseClient := configmock.NewMockConfigManagementAPI_ListApplicationAgentsClient(mockCtrl)
			mockListApplicationAgentsRequest := configpb.ListApplicationAgentsRequest{
				AppSpaceId: "gid:like-real-app-space-id",
				Match:      []string{"like-real-application-agent-name"},
			}

			mockClient.EXPECT().ListApplicationAgents(
				gomock.Any(),
				gomock.Eq(&mockListApplicationAgentsRequest),
				gomock.Any(),
			).Return(mockResponseClient, status.Error(codes.InvalidArgument, "status error"))

			stream, err := configClient.ListApplicationAgents(context.Background(), &mockListApplicationAgentsRequest)
			Expect(err).ToNot(Succeed())
			Expect(stream).To(BeNil())
		})

		DescribeTable("ListError",
			func(mockListApplicationAgentsRequest *configpb.ListApplicationAgentsRequest, message string) {
				resp, err := configClient.ListApplicationAgents(context.Background(), mockListApplicationAgentsRequest)
				Ω(resp).To(BeNil())
				Ω(err).To(MatchError(ContainSubstring(message)))
			},
			Entry(
				"AppSpace error",
				&configpb.ListApplicationAgentsRequest{
					AppSpaceId: "app-space-id",
					Match:      []string{"like-real-application-agent-name"},
				},
				"Id: value length must be between 22",
			),
			Entry(
				"Match error",
				&configpb.ListApplicationAgentsRequest{
					AppSpaceId: "gid:like-real-app-space-id",
					Match:      []string{},
				},
				"Match: value must contain at least 1 item",
			),
		)
	})

	Describe("ApplicationAgentDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteApplicationAgent(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationAgentRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteApplicationAgent(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationAgentRequest{
				Id:   "gid:like-real-application-agent-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationAgentResponse{
				Bookmark: "something-like-bookmark-which-is-long-enough",
			}

			mockClient.EXPECT().DeleteApplicationAgent(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteApplicationAgent(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationAgentRequest{
				Id:   "gid:like-real-application-agent-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationAgentResponse{}

			mockClient.EXPECT().DeleteApplicationAgent(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteApplicationAgent(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
