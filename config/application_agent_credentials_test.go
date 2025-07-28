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
)

var _ = Describe("ApplicationAgentCredentials", func() {
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

	Describe("ReadApplicationAgentCredential", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadApplicationAgentCredential(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadApplicationAgentCredentialRequest{
				Id: "gid:like",
			}
			resp, err := configClient.ReadApplicationAgentCredential(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Read", func() {
			req := &configpb.ReadApplicationAgentCredentialRequest{
				Id: "gid:like-real-application-agent-id",
			}
			now := time.Now()
			future := now.AddDate(2, 0, 0)
			beResp := &configpb.ReadApplicationAgentCredentialResponse{
				ApplicationAgentCredential: &configpb.ApplicationAgentCredential{
					Id:                 "gid:like-real-application-agent-credential-id",
					DisplayName:        "Like Real Application Agent Name",
					CreatedBy:          "creator",
					CreateTime:         timestamppb.Now(),
					ApplicationId:      "gid:like-real-application-id",
					ApplicationAgentId: "gid:like-real-application-agent-id",
					Kid:                "G9uMQzWWeP9lLvf7qKLhmeHabgZI_Mp8fnH7FJGRWHQ",
					ExpireTime:         timestamppb.New(future),
				},
			}

			mockClient.EXPECT().
				ReadApplicationAgentCredential(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := configClient.ReadApplicationAgentCredential(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("ReadError", func() {
			req := &configpb.ReadApplicationAgentCredentialRequest{
				Id: "gid:like-real-application-agent-id",
			}
			beResp := &configpb.ReadApplicationAgentCredentialResponse{}

			mockClient.EXPECT().
				ReadApplicationAgentCredential(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.ReadApplicationAgentCredential(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ApplicationAgentCredentialRegister", func() {
		It("Nil request", func() {
			resp, err := configClient.RegisterApplicationAgentCredential(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Register", func() {
			req := &configpb.RegisterApplicationAgentCredentialRequest{
				ApplicationAgentId: "gid:like-real-application-agent-id",
				DisplayName:        "Like real Application Agent Credential Name",
			}
			now := time.Now()
			future := now.AddDate(2, 0, 0)
			beResp := &configpb.RegisterApplicationAgentCredentialResponse{
				Id:                 "gid:like-real-application-agent-credential-id",
				ApplicationAgentId: "gid:like-real-application-agent-id",
				CreateTime:         timestamppb.Now(),
				AgentConfig:        []byte("falwJAyAQawtoWLpIp9OUj"),
				Kid:                "G9uMQzWWeP9lLvf7qKLhmeHabgZI_Mp8fnH7FJGRWHQ",
				ExpireTime:         timestamppb.New(future),
			}

			mockClient.EXPECT().RegisterApplicationAgentCredential(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.RegisterApplicationAgentCredential(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("RegisterError", func() {
			req := &configpb.RegisterApplicationAgentCredentialRequest{
				ApplicationAgentId: "gid:like-real-application-agent-id",
				DisplayName:        "Like real Application Agent Credential Name",
			}
			beResp := &configpb.RegisterApplicationAgentCredentialResponse{}

			mockClient.EXPECT().RegisterApplicationAgentCredential(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.RegisterApplicationAgentCredential(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("RegisterNonValid", func() {
			req := &configpb.RegisterApplicationAgentCredentialRequest{
				ApplicationAgentId: "error-app-id",
				DisplayName:        "Like real Application Agent Credentials Name",
			}

			resp, err := configClient.RegisterApplicationAgentCredential(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})
	})

	Describe("ApplicationAgentDelete", func() {
		It("Nil request", func() {
			resp, err := configClient.DeleteApplicationAgentCredential(ctx, nil)
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
			req := &configpb.DeleteApplicationAgentCredentialRequest{
				Id:   "like-real",
				Etag: etagPb,
			}
			resp, err := configClient.DeleteApplicationAgentCredential(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		It("Delete", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationAgentCredentialRequest{
				Id:   "gid:like-real-application-agent-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationAgentCredentialResponse{}

			mockClient.EXPECT().DeleteApplicationAgentCredential(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, nil)

			resp, err := configClient.DeleteApplicationAgentCredential(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("DeleteError", func() {
			etagPb := &wrapperspb.StringValue{Value: "123qwert"}
			req := &configpb.DeleteApplicationAgentCredentialRequest{
				Id:   "gid:like-real-application-agent-id",
				Etag: etagPb,
			}
			beResp := &configpb.DeleteApplicationAgentCredentialResponse{}

			mockClient.EXPECT().DeleteApplicationAgentCredential(
				gomock.Any(),
				test.WrapMatcher(test.EqualProto(req)),
				gomock.Any(),
			).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

			resp, err := configClient.DeleteApplicationAgentCredential(ctx, req)
			Expect(err).ToNot(Succeed())
			Expect(resp).To(BeNil())
		})
	})
})
