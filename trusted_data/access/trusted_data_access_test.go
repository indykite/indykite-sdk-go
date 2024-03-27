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

package access_test

import (
	"context"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accesspb "github.com/indykite/indykite-sdk-go/gen/indykite/trusted_data/access/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	accessmock "github.com/indykite/indykite-sdk-go/test/trusted_data/access/v1beta1"
	"github.com/indykite/indykite-sdk-go/trusted_data/access"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TDA", func() {
	mockErrorCode := codes.InvalidArgument
	mockErrorMessage := "mockError"

	Describe("Grant consent", func() {
		mockGrantConsentRequest := accesspb.GrantConsentRequest{
			User:           &accesspb.GrantConsentRequest_UserId{UserId: "gid:886hfic8fswlz3zjrc2e3nun9xs"},
			ConsentId:      "gid:aof7m7kl966tmzct2kc75kxeisl",
			RevokeAfterUse: true,
		}

		It("grants consent and returns response", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := accessmock.NewMockTrustedDataControlAPIClient(mockCtrl)

			client, err := access.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().GrantConsent(
				gomock.Any(),
				gomock.Eq(&mockGrantConsentRequest),
				gomock.Any(),
			).Return(&accesspb.GrantConsentResponse{}, nil)

			resp, err := client.GrantConsent(context.Background(), &mockGrantConsentRequest)
			Ω(resp).ToNot(BeNil())
			Ω(err).To(Succeed())
		})

		It("handles and returns error", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := accessmock.NewMockTrustedDataControlAPIClient(mockCtrl)

			client, err := access.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().GrantConsent(
				gomock.Any(),
				gomock.Eq(&mockGrantConsentRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := client.GrantConsent(context.Background(), &mockGrantConsentRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("Access Consented Data", func() {
		mockAccessConsentedDataRequest := accesspb.AccessConsentedDataRequest{
			ConsentId: "gid:aof7m7kl966tmzct2kc75kxeisl",
			UserId:    "gid:886hfic8fswlz3zjrc2e3nun9xs",
		}

		It("accesses consented data and returns response", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := accessmock.NewMockTrustedDataControlAPIClient(mockCtrl)

			client, err := access.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().AccessConsentedData(
				gomock.Any(),
				gomock.Eq(&mockAccessConsentedDataRequest),
				gomock.Any(),
			).Return(&accesspb.AccessConsentedDataResponse{}, nil)

			resp, err := client.AccessConsentedData(context.Background(), &mockAccessConsentedDataRequest)
			Ω(resp).ToNot(BeNil())
			Ω(err).To(Succeed())
		})

		It("handles and returns error", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := accessmock.NewMockTrustedDataControlAPIClient(mockCtrl)

			client, err := access.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().AccessConsentedData(
				gomock.Any(),
				gomock.Eq(&mockAccessConsentedDataRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := client.AccessConsentedData(context.Background(), &mockAccessConsentedDataRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("Revoke consent", func() {
		mockRevokeConsentRequest := accesspb.RevokeConsentRequest{
			User:      &accesspb.RevokeConsentRequest_UserId{UserId: "gid:886hfic8fswlz3zjrc2e3nun9xs"},
			ConsentId: "gid:aof7m7kl966tmzct2kc75kxeisl",
		}

		It("revokes consent", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := accessmock.NewMockTrustedDataControlAPIClient(mockCtrl)

			client, err := access.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().RevokeConsent(
				gomock.Any(),
				gomock.Eq(&mockRevokeConsentRequest),
				gomock.Any(),
			).Return(&accesspb.RevokeConsentResponse{}, nil)

			resp, err := client.RevokeConsent(context.Background(), &mockRevokeConsentRequest)
			Ω(resp).ToNot(BeNil())
			Ω(err).To(Succeed())
		})

		It("handles and returns error", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := accessmock.NewMockTrustedDataControlAPIClient(mockCtrl)

			client, err := access.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().RevokeConsent(
				gomock.Any(),
				gomock.Eq(&mockRevokeConsentRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := client.RevokeConsent(context.Background(), &mockRevokeConsentRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})
})
