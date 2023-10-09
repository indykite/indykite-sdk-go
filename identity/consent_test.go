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

package identity_test

import (
	"context"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	identitypb "github.com/indykite/indykite-sdk-go/gen/indykite/identity/v1beta2"
	"github.com/indykite/indykite-sdk-go/identity"
	"github.com/indykite/indykite-sdk-go/test"
	midentity "github.com/indykite/indykite-sdk-go/test/identity/v1beta2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consent", func() {
	mockErrorCode := codes.InvalidArgument
	mockErrorMessage := "mockError"

	Describe("Create consent", func() {
		mockCreateConsentRequest := identitypb.CreateConsentRequest{
			PiiPrincipalId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
			PiiProcessorId: "gid:aof7m7kl966tmzct2kc75kxeisl",
			Properties:     []string{"email"},
		}

		It("creates consent and returns response", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)

			client, err := identity.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().CreateConsent(
				gomock.Any(),
				gomock.Eq(&mockCreateConsentRequest),
				gomock.Any(),
			).Return(&identitypb.CreateConsentResponse{}, nil)

			resp, err := client.CreateConsent(context.Background(), &mockCreateConsentRequest)
			Ω(resp).ToNot(BeNil())
			Ω(err).To(Succeed())
		})

		It("handles and returns error", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)

			client, err := identity.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().CreateConsent(
				gomock.Any(),
				gomock.Eq(&mockCreateConsentRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := client.CreateConsent(context.Background(), &mockCreateConsentRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("Revoke consent", func() {
		mockRevokeConsentRequest := identitypb.RevokeConsentRequest{
			PiiPrincipalId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
			ConsentIds:     []string{"911cfdbe-edd3-4701-a9f3-c90913ad1d7d"},
		}

		It("revokes consent", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)

			client, err := identity.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().RevokeConsent(
				gomock.Any(),
				gomock.Eq(&mockRevokeConsentRequest),
				gomock.Any(),
			).Return(&identitypb.RevokeConsentResponse{}, nil)

			resp, err := client.RevokeConsent(context.Background(), &mockRevokeConsentRequest)
			Ω(resp).ToNot(BeNil())
			Ω(err).To(Succeed())
		})

		It("handles and returns error", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)

			client, err := identity.NewTestClient(mockClient)
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

	Describe("List consents", func() {
		mockListConsentsRequest := identitypb.ListConsentsRequest{
			PiiPrincipalId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
		}

		It("Lists consents and return response", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)
			mockResponseClient := midentity.NewMockIdentityManagementAPI_ListConsentsClient(mockCtrl)

			client, err := identity.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().ListConsents(
				gomock.Any(),
				gomock.Eq(&mockListConsentsRequest),
				gomock.Any(),
			).Return(mockResponseClient, nil)

			resp, err := client.ListConsents(context.Background(), &mockListConsentsRequest)
			Ω(resp).ToNot(BeNil())
			Ω(err).To(Succeed())
		})

		It("handles and returns error", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)

			client, err := identity.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			mockClient.EXPECT().ListConsents(
				gomock.Any(),
				gomock.Eq(&mockListConsentsRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := client.ListConsents(context.Background(), &mockListConsentsRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})
})
