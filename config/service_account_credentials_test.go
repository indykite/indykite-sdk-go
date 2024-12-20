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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/indykite/indykite-sdk-go/config"
	sdkerrors "github.com/indykite/indykite-sdk-go/errors"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	configmock "github.com/indykite/indykite-sdk-go/test/config/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Account Credentials", func() {
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

	It("Register nil request", func() {
		resp, err := configClient.RegisterServiceAccountCredential(ctx, nil)
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())

		var clientErr *sdkerrors.ClientError
		Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
		Expect(clientErr.Unwrap()).To(Succeed())
		Expect(clientErr.Message()).To(Equal("invalid nil request"))
		Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
	})

	It("Register", func() {
		req := &configpb.RegisterServiceAccountCredentialRequest{
			ServiceAccountId: "gid:like-real-service-account-id",
			DisplayName:      "My Credentials",
		}
		beResp := &configpb.RegisterServiceAccountCredentialResponse{
			Kid: "kid",
		}

		mockClient.EXPECT().RegisterServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, nil)

		resp, err := configClient.RegisterServiceAccountCredential(ctx, req)
		Expect(err).To(Succeed())
		Expect(resp).To(test.EqualProto(beResp))
	})

	It("RegisterNonValid", func() {
		req := &configpb.RegisterServiceAccountCredentialRequest{
			ServiceAccountId: "error-sa-id",
			DisplayName:      "My Credentials",
		}

		resp, err := configClient.RegisterServiceAccountCredential(ctx, req)
		Expect(resp).To(BeNil())
		Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
	})

	It("RegisterError", func() {
		req := &configpb.RegisterServiceAccountCredentialRequest{
			ServiceAccountId: "gid:like-real-service-account-id",
			DisplayName:      "My Credentials",
		}
		beResp := &configpb.RegisterServiceAccountCredentialResponse{}

		mockClient.EXPECT().RegisterServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

		resp, err := configClient.RegisterServiceAccountCredential(ctx, req)
		Expect(err).ToNot(Succeed())
		Expect(resp).To(BeNil())
	})

	It("Read nil request", func() {
		resp, err := configClient.ReadServiceAccountCredential(ctx, nil)
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())

		var clientErr *sdkerrors.ClientError
		Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
		Expect(clientErr.Unwrap()).To(Succeed())
		Expect(clientErr.Message()).To(Equal("invalid nil request"))
		Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
	})

	It("Read", func() {
		req := &configpb.ReadServiceAccountCredentialRequest{
			Id: "gid:like-real-service-account-credential-id",
		}

		beResp := &configpb.ReadServiceAccountCredentialResponse{
			ServiceAccountCredential: &configpb.ServiceAccountCredential{
				Id: "gid:like-real-service-account-credential-id",
			},
		}

		mockClient.EXPECT().ReadServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, nil)

		resp, err := configClient.ReadServiceAccountCredential(ctx, req)
		Expect(err).To(Succeed())
		Expect(resp).To(test.EqualProto(beResp))
	})

	It("ReadNonValid", func() {
		req := &configpb.ReadServiceAccountCredentialRequest{
			Id: "gid:like",
		}
		resp, err := configClient.ReadServiceAccountCredential(ctx, req)
		Expect(resp).To(BeNil())
		Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
	})

	It("ReadError", func() {
		req := &configpb.ReadServiceAccountCredentialRequest{
			Id: "gid:like-real-service-account-credential-id",
		}

		beResp := &configpb.ReadServiceAccountCredentialResponse{}

		mockClient.EXPECT().ReadServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

		resp, err := configClient.ReadServiceAccountCredential(ctx, req)
		Expect(err).ToNot(Succeed())
		Expect(resp).To(BeNil())
	})

	It("Delete nil request", func() {
		resp, err := configClient.DeleteServiceAccountCredential(ctx, nil)
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())

		var clientErr *sdkerrors.ClientError
		Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
		Expect(clientErr.Unwrap()).To(Succeed())
		Expect(clientErr.Message()).To(Equal("invalid nil request"))
		Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
	})

	It("Delete", func() {
		req := &configpb.DeleteServiceAccountCredentialRequest{
			Id: "gid:like-real-service-account-credential-id",
		}
		beResp := &configpb.DeleteServiceAccountCredentialResponse{}

		mockClient.EXPECT().DeleteServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, nil)

		resp, err := configClient.DeleteServiceAccountCredential(ctx, req)
		Expect(err).To(Succeed())
		Expect(resp).To(test.EqualProto(beResp))
	})

	It("DeleteNonValid", func() {
		etagPb := &wrapperspb.StringValue{Value: "123qwert"}
		req := &configpb.DeleteServiceAccountCredentialRequest{
			Id:   "like-real",
			Etag: etagPb,
		}
		resp, err := configClient.DeleteServiceAccountCredential(ctx, req)
		Expect(resp).To(BeNil())
		Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
	})

	It("Delete", func() {
		req := &configpb.DeleteServiceAccountCredentialRequest{
			Id: "gid:like-real-service-account-credential-id",
		}
		beResp := &configpb.DeleteServiceAccountCredentialResponse{}

		mockClient.EXPECT().DeleteServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

		resp, err := configClient.DeleteServiceAccountCredential(ctx, req)
		Expect(err).ToNot(Succeed())
		Expect(resp).To(BeNil())
	})
})
