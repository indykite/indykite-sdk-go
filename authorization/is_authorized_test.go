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

package authorization_test

import (
	"context"
	"errors"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indykite/indykite-sdk-go/authorization"
	sdkerrors "github.com/indykite/indykite-sdk-go/errors"
	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	objectpb "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	authorizationmock "github.com/indykite/indykite-sdk-go/test/authorization/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsAuthorized", func() {
	var (
		ctx                 context.Context
		mockCtrl            *gomock.Controller
		mockClient          *authorizationmock.MockAuthorizationAPIClient
		authorizationClient *authorization.Client
		resourceExample     = []*authorizationpb.IsAuthorizedRequest_Resource{
			{ExternalId: "external_id_value", Type: "Asset", Actions: []string{"ACTION1", "ACTION2"}},
		}
		inputParam = map[string]*authorizationpb.InputParam{
			"Color": {Value: &authorizationpb.InputParam_StringValue{StringValue: "red"}},
		}
		policyTags = []string{"sometag"}
		tokenGood  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
			"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9." +
			"dyt0CoTl4WoVjAHI9Q_CwSKhl6d_9rhM3NrXuJttkao" // #nosec G101
		tokenBad = "token_invalid_format"
		beResp   = &authorizationpb.IsAuthorizedResponse{
			DecisionTime: timestamppb.New(time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)),
			Decisions: map[string]*authorizationpb.IsAuthorizedResponse_ResourceType{
				"Resource": {
					Resources: map[string]*authorizationpb.IsAuthorizedResponse_Resource{
						"external_id_value": {
							Actions: map[string]*authorizationpb.IsAuthorizedResponse_Action{
								"ACTION1": {Allow: true},
								"ACTION2": {Allow: false},
							},
						},
					},
				},
			},
		}
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = authorizationmock.NewMockAuthorizationAPIClient(mockCtrl)

		var err error
		authorizationClient, err = authorization.NewTestClient(ctx, mockClient)
		Expect(err).To(Succeed())
	})

	Describe("IsAuthorized", func() {
		It("Nil params", func() {
			resp, err := authorizationClient.IsAuthorized(ctx, nil, nil, nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Nil token", func() {
			resp, err := authorizationClient.IsAuthorizedByToken(ctx, "", nil, nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Message()).To(ContainSubstring("unable to call IsAuthorized client endpoint"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong DT should return a validation error in the response", func() {
			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_DigitalTwinId{
						DigitalTwinId: &authorizationpb.DigitalTwin{
							Id: "gid:like",
						},
					},
				},
			}

			resp, err := authorizationClient.IsAuthorizedWithRawRequest(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring(
				"invalid DigitalTwin.Id: value length must be between 27 and 100 runes",
			)))
		})

		It("IsAuthorizedDT", func() {
			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_DigitalTwinId{
						DigitalTwinId: &authorizationpb.DigitalTwin{
							Id: "gid:like-real-digital_twin-id-at-least-27-char",
						},
					},
				},
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				IsAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.IsAuthorizedWithRawRequest(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("IsAuthorizedProperty", func() {
			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_DigitalTwinProperty{
						DigitalTwinProperty: &authorizationpb.Property{
							Type: "email_for_example",
							Value: &objectpb.Value{
								Value: &objectpb.Value_StringValue{
									StringValue: "test@example.com",
								},
							},
						},
					},
				},
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				IsAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.IsAuthorizedByProperty(
				ctx,
				&authorizationpb.Property{
					Type: "email_for_example",
					Value: &objectpb.Value{
						Value: &objectpb.Value_StringValue{
							StringValue: "test@example.com",
						},
					},
				},
				resourceExample,
				inputParam,
				policyTags,
			)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("IsAuthorizedExternalID", func() {
			externalID := &authorizationpb.ExternalID{
				Type:       "Person",
				ExternalId: "576eb486-28f6-4756-95be-ba6da362b2a7",
			}

			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_ExternalId{
						ExternalId: externalID,
					},
				},
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				IsAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.IsAuthorizedByExternalID(
				ctx,
				externalID,
				resourceExample,
				inputParam,
				policyTags,
			)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("IsAuthorizedToken", func() {
			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_IndykiteAccessToken{
						IndykiteAccessToken: tokenGood,
					},
				},
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				IsAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.IsAuthorizedByToken(
				ctx,
				tokenGood,
				resourceExample,
				inputParam,
				policyTags,
			)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("IsAuthorizedTokenWrongFormat", func() {
			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_IndykiteAccessToken{
						IndykiteAccessToken: tokenBad,
					},
				},
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			resp, err := authorizationClient.IsAuthorizedWithRawRequest(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("invalid token format cause")))
		})

		It("Invalid status", func() {
			req := &authorizationpb.IsAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_IndykiteAccessToken{
						IndykiteAccessToken: tokenGood,
					},
				},
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			statusErr := status.New(codes.AlreadyExists, "something wrong").Err()
			mockClient.EXPECT().
				IsAuthorized(gomock.Any(), req).
				Return(nil, statusErr)

			resp, err := authorizationClient.IsAuthorizedWithRawRequest(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
})
