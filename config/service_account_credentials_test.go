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

	"github.com/golang/mock/gomock"

	"github.com/indykite/indykite-sdk-go/config"
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

	It("Delete", func() {
		req := &configpb.DeleteServiceAccountCredentialRequest{
			Id: "gid:like-real-service-account-credential-id",
		}
		beResp := &configpb.DeleteServiceAccountCredentialResponse{
			Bookmark: "bookmark",
		}

		mockClient.EXPECT().DeleteServiceAccountCredential(
			gomock.Any(),
			test.WrapMatcher(test.EqualProto(req)),
			gomock.Any(),
		).Return(beResp, nil)

		resp, err := configClient.DeleteServiceAccountCredential(ctx, req)
		Expect(err).To(Succeed())
		Expect(resp).To(test.EqualProto(beResp))
	})
})
