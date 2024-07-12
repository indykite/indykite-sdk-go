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

package tda_test

import (
	"context"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
	objects2 "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
	tdapb "github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1"
	"github.com/indykite/indykite-sdk-go/tda"
	"github.com/indykite/indykite-sdk-go/test"
	tdamock "github.com/indykite/indykite-sdk-go/test/tda/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TDA", func() {
	mockErrorCode := codes.InvalidArgument
	mockErrorMessage := "mockError"

	var (
		mockCtrl   *gomock.Controller
		mockClient *tdamock.MockTrustedDataAccessAPIClient
		tdaClient  *tda.Client
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = tdamock.NewMockTrustedDataAccessAPIClient(mockCtrl)
		var err error
		tdaClient, err = tda.NewTestClient(mockClient)
		Ω(err).To(Succeed())
	})

	Describe("Grant consent", func() {
		mockGrantConsentRequest := tdapb.GrantConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_UserId{
					UserId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
				},
			},
			ConsentId:      "gid:aof7m7kl966tmzct2kc75kxeisl",
			ValidityPeriod: 86400,
		}

		mockGrantConsentRequest2 := tdapb.GrantConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ExternalId{
					ExternalId: &knowledgeobjects.User_ExternalID{
						Type:       "Person",
						ExternalId: "HLEgiljrtoNEiyX",
					},
				},
			},
			ConsentId:      "gid:aof7m7kl966tmzct2kc75kxeisl",
			ValidityPeriod: 86400,
		}

		mockGrantConsentRequest3 := tdapb.GrantConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_Property_{
					Property: &knowledgeobjects.User_Property{
						Type: "abc",
						Value: &objects.Value{
							Value: &objects.Value_StringValue{StringValue: "something"},
						},
					},
				},
			},
			ConsentId:      "gid:aof7m7kl966tmzct2kc75kxeisl",
			ValidityPeriod: 86400,
		}

		mockGrantConsentRequest4 := tdapb.GrantConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ThirdPartyToken{
					ThirdPartyToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6IlVzTjNV",
				},
			},
			ConsentId:      "gid:aof7m7kl966tmzct2kc75kxeisl",
			ValidityPeriod: 86400,
		}

		DescribeTable("GrantConsentSuccess",
			func(req *tdapb.GrantConsentRequest, beResp *tdapb.GrantConsentResponse) {
				mockClient.EXPECT().GrantConsent(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

				resp, err := tdaClient.GrantConsent(context.Background(), req)
				Ω(resp).ToNot(BeNil())
				Ω(err).To(Succeed())
				Ω(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"grants consent user id and returns response",
				&mockGrantConsentRequest,
				&tdapb.GrantConsentResponse{},
			),
			Entry(
				"grants consent with external id and returns response",
				&mockGrantConsentRequest2,
				&tdapb.GrantConsentResponse{},
			),
			Entry(
				"grants consent with property and returns response",
				&mockGrantConsentRequest3,
				&tdapb.GrantConsentResponse{},
			),
			Entry(
				"grants consent with token and returns response",
				&mockGrantConsentRequest4,
				&tdapb.GrantConsentResponse{},
			),
		)

		It("handles and returns error", func() {
			mockClient.EXPECT().GrantConsent(
				gomock.Any(),
				gomock.Eq(&mockGrantConsentRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := tdaClient.GrantConsent(context.Background(), &mockGrantConsentRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("Trusted Consented Data", func() {
		mockDataAccessRequest := tdapb.DataAccessRequest{
			ConsentId:     "gid:aof7m7kl966tmzct2kc75kxeisl",
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_UserId{
					UserId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
				},
			},
		}

		mockDataAccessRequest2 := tdapb.DataAccessRequest{
			ConsentId:     "gid:aof7m7kl966tmzct2kc75kxeisl",
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ExternalId{
					ExternalId: &knowledgeobjects.User_ExternalID{
						Type:       "Person",
						ExternalId: "HLEgiljrtoNEiyX",
					},
				},
			},
		}

		mockDataAccessRequest3 := tdapb.DataAccessRequest{
			ConsentId:     "gid:aof7m7kl966tmzct2kc75kxeisl",
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_Property_{
					Property: &knowledgeobjects.User_Property{
						Type: "abc",
						Value: &objects.Value{
							Value: &objects.Value_StringValue{StringValue: "something"},
						},
					},
				},
			},
		}

		mockDataAccessRequest4 := tdapb.DataAccessRequest{
			ConsentId:     "gid:aof7m7kl966tmzct2kc75kxeisl",
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
		}

		mockDataAccessRequest5 := tdapb.DataAccessRequest{
			ConsentId:     "gid:aof7m7kl966tmzct2kc75kxeisl",
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ThirdPartyToken{
					ThirdPartyToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6IlVzTjNV",
				},
			},
		}

		mockDataAccessResponse := tdapb.DataAccessResponse{
			Nodes: []*tdapb.TrustedDataNode{
				{
					Id:         "gid:xyz",
					ExternalId: "0000",
					Type:       "Person",
					IsIdentity: true,
					Nodes: []*tdapb.TrustedDataNode{
						{
							Id:         "gid:abc",
							ExternalId: "1111",
							Type:       "Car",
							Properties: []*knowledgeobjects.Property{
								{
									Type: "PlateNumber",
									Value: &objects2.Value{
										Type: &objects2.Value_StringValue{
											StringValue: "N4567",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		DescribeTable("DataAccessSuccess",
			func(req *tdapb.DataAccessRequest, beResp *tdapb.DataAccessResponse) {
				mockClient.EXPECT().DataAccess(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

				resp, err := tdaClient.DataAccess(context.Background(), req)
				Ω(resp).ToNot(BeNil())
				Ω(err).To(Succeed())
				Ω(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"data access with user id and returns response",
				&mockDataAccessRequest,
				&mockDataAccessResponse,
			),
			Entry(
				"data access with external id and returns response",
				&mockDataAccessRequest2,
				&mockDataAccessResponse,
			),
			Entry(
				"data access with property and returns response",
				&mockDataAccessRequest3,
				&mockDataAccessResponse,
			),
			Entry(
				"data access without user and returns response",
				&mockDataAccessRequest4,
				&mockDataAccessResponse,
			),
			Entry(
				"data access with token and returns response",
				&mockDataAccessRequest5,
				&mockDataAccessResponse,
			),
		)

		It("handles and returns error", func() {
			mockClient.EXPECT().DataAccess(
				gomock.Any(),
				gomock.Eq(&mockDataAccessRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := tdaClient.DataAccess(context.Background(), &mockDataAccessRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("Revoke consent", func() {
		mockRevokeConsentRequest := tdapb.RevokeConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_UserId{
					UserId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
				},
			},
			ConsentId: "gid:aof7m7kl966tmzct2kc75kxeisl",
		}

		mockRevokeConsentRequest2 := tdapb.RevokeConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ExternalId{
					ExternalId: &knowledgeobjects.User_ExternalID{
						Type:       "Person",
						ExternalId: "HLEgiljrtoNEiyX",
					},
				},
			},
			ConsentId: "gid:aof7m7kl966tmzct2kc75kxeisl",
		}

		mockRevokeConsentRequest3 := tdapb.RevokeConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_Property_{
					Property: &knowledgeobjects.User_Property{
						Type: "abc",
						Value: &objects.Value{
							Value: &objects.Value_StringValue{StringValue: "something"},
						},
					},
				},
			},
			ConsentId: "gid:aof7m7kl966tmzct2kc75kxeisl",
		}

		mockRevokeConsentRequest4 := tdapb.RevokeConsentRequest{
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ThirdPartyToken{
					ThirdPartyToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6IlVzTjNV",
				},
			},
			ConsentId: "gid:aof7m7kl966tmzct2kc75kxeisl",
		}

		DescribeTable("RevokeConsentSuccess",
			func(req *tdapb.RevokeConsentRequest, beResp *tdapb.RevokeConsentResponse) {
				mockClient.EXPECT().RevokeConsent(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

				resp, err := tdaClient.RevokeConsent(context.Background(), req)
				Ω(resp).ToNot(BeNil())
				Ω(err).To(Succeed())
				Ω(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"revokes consent with user id",
				&mockRevokeConsentRequest,
				&tdapb.RevokeConsentResponse{},
			),
			Entry(
				"revokes consent with external id",
				&mockRevokeConsentRequest2,
				&tdapb.RevokeConsentResponse{},
			),
			Entry(
				"revokes consent with property",
				&mockRevokeConsentRequest3,
				&tdapb.RevokeConsentResponse{},
			),
			Entry(
				"revokes consent with token",
				&mockRevokeConsentRequest4,
				&tdapb.RevokeConsentResponse{},
			),
		)

		It("handles and returns error", func() {
			mockClient.EXPECT().RevokeConsent(
				gomock.Any(),
				gomock.Eq(&mockRevokeConsentRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := tdaClient.RevokeConsent(context.Background(), &mockRevokeConsentRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("ListConsents", func() {
		mockListConsentsRequest := tdapb.ListConsentsRequest{
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_UserId{
					UserId: "gid:886hfic8fswlz3zjrc2e3nun9xs",
				},
			},
		}

		mockListConsentsRequest2 := tdapb.ListConsentsRequest{
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ExternalId{
					ExternalId: &knowledgeobjects.User_ExternalID{
						Type:       "Person",
						ExternalId: "HLEgiljrtoNEiyX",
					},
				},
			},
		}

		mockListConsentsRequest3 := tdapb.ListConsentsRequest{
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_Property_{
					Property: &knowledgeobjects.User_Property{
						Type: "abc",
						Value: &objects.Value{
							Value: &objects.Value_StringValue{StringValue: "something"},
						},
					},
				},
			},
		}

		mockListConsentsRequest4 := tdapb.ListConsentsRequest{
			ApplicationId: "gid:AAAAFbJmG6cY2032lHlm1H0HImY",
			User: &knowledgeobjects.User{
				User: &knowledgeobjects.User_ThirdPartyToken{
					ThirdPartyToken: "eyJhbGciOiJFUzI1NiIsImtpZCI6IlVzTjNV",
				},
			},
		}

		DescribeTable("ListConsentsSuccess",
			func(req *tdapb.ListConsentsRequest, beResp *tdapb.ListConsentsResponse) {
				mockClient.EXPECT().ListConsents(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

				resp, err := tdaClient.ListConsents(context.Background(), req)
				Ω(resp).ToNot(BeNil())
				Ω(err).To(Succeed())
				Ω(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"list of consents with user id",
				&mockListConsentsRequest,
				&tdapb.ListConsentsResponse{
					Consents: []*tdapb.Consent{
						{
							Id:         "gid:xyz",
							Properties: []string{"email", "name"},
						},
					},
				},
			),
			Entry(
				"list of consents  with external id",
				&mockListConsentsRequest2,
				&tdapb.ListConsentsResponse{
					Consents: []*tdapb.Consent{
						{
							Id:         "gid:xyz",
							Properties: []string{"email", "name"},
						},
					},
				},
			),
			Entry(
				"list of consents  with property",
				&mockListConsentsRequest3,
				&tdapb.ListConsentsResponse{
					Consents: []*tdapb.Consent{
						{
							Id:         "gid:xyz",
							Properties: []string{"email", "name"},
						},
					},
				},
			),
			Entry(
				"list of consents  with token",
				&mockListConsentsRequest4,
				&tdapb.ListConsentsResponse{
					Consents: []*tdapb.Consent{
						{
							Id:         "gid:xyz",
							Properties: []string{"email", "name"},
						},
					},
				},
			),
		)

		It("handles and returns error", func() {
			mockClient.EXPECT().ListConsents(
				gomock.Any(),
				gomock.Eq(&mockListConsentsRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := tdaClient.ListConsents(context.Background(), &mockListConsentsRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})
})
