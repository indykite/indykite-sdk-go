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
	objectpb2 "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
	"github.com/indykite/indykite-sdk-go/test"
	authorizationmock "github.com/indykite/indykite-sdk-go/test/authorization/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WhatAuthorized", func() {
	var (
		ctx                 context.Context
		mockCtrl            *gomock.Controller
		mockClient          *authorizationmock.MockAuthorizationAPIClient
		authorizationClient *authorization.Client
		resourceTypeExample = []*authorizationpb.WhatAuthorizedRequest_ResourceType{
			{Type: "Asset", Actions: []string{"ACTION1", "ACTION2"}},
		}
		inputParam = map[string]*authorizationpb.InputParam{
			"Color": {Value: &authorizationpb.InputParam_StringValue{StringValue: "red"}},
		}
		inputParamInt = map[string]*authorizationpb.InputParam{
			"Age": {Value: &authorizationpb.InputParam_IntegerValue{IntegerValue: 21}},
		}
		inputParamTime = map[string]*authorizationpb.InputParam{
			"DateT": {Value: &authorizationpb.InputParam_TimeValue{
				TimeValue: timestamppb.New(
					time.Date(2024, 9, 9, 9, 9, 0, 0, time.UTC),
				)}},
		}
		inputParamArray = map[string]*authorizationpb.InputParam{
			"ArrayParam": {
				Value: &authorizationpb.InputParam_ArrayValue{
					ArrayValue: &objectpb2.Array{
						Values: []*objectpb2.Value{
							{
								Type: &objectpb2.Value_StringValue{StringValue: "ParkingLot"},
							},
							{
								Type: &objectpb2.Value_IntegerValue{IntegerValue: 42},
							},
						},
					},
				},
			},
		}
		inputParamMap = map[string]*authorizationpb.InputParam{
			"MapParam": {
				Value: &authorizationpb.InputParam_MapValue{
					MapValue: &objectpb2.Map{
						Fields: map[string]*objectpb2.Value{
							"ParkingLot": {
								Type: &objectpb2.Value_StringValue{StringValue: "Square street"},
							},
							"Lot": {
								Type: &objectpb2.Value_IntegerValue{IntegerValue: 42},
							},
						},
					},
				},
			},
		}
		policyTags = []string{"sometag"}
		tokenGood  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
			"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9." +
			"dyt0CoTl4WoVjAHI9Q_CwSKhl6d_9rhM3NrXuJttkao" // #nosec G101
		tokenBad = "token_invalid_format"
		beResp   = &authorizationpb.WhatAuthorizedResponse{
			DecisionTime: timestamppb.New(time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)),
			Decisions: map[string]*authorizationpb.WhatAuthorizedResponse_ResourceType{
				"ResourceType": {
					Actions: map[string]*authorizationpb.WhatAuthorizedResponse_Action{
						"ACTION1": {
							Resources: []*authorizationpb.WhatAuthorizedResponse_Resource{
								{ExternalId: "external_id_value1"},
								{ExternalId: "external_id_value2"},
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

	Describe("WhatAuthorized", func() {
		It("Nil params", func() {
			resp, err := authorizationClient.WhatAuthorized(ctx, nil, nil, nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Nil token", func() {
			resp, err := authorizationClient.WhatAuthorizedByToken(ctx, "", nil, nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Message()).To(ContainSubstring("unable to call WhatAuthorized client endpoint"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong DT should return a validation error in the response", func() {
			req := &authorizationpb.WhatAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_DigitalTwinId{
						DigitalTwinId: &authorizationpb.DigitalTwin{
							Id: "gid:like",
						},
					},
				},
			}

			resp, err := authorizationClient.WhatAuthorizedWithRawRequest(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring(
				"invalid DigitalTwin.Id: value length must be between 27 and 100 runes",
			)))
		})

		It("WhatAuthorizedDT", func() {
			req := &authorizationpb.WhatAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_DigitalTwinId{
						DigitalTwinId: &authorizationpb.DigitalTwin{
							Id: "gid:like-real-digital_twin-id-at-least-27-char",
						},
					},
				},
				ResourceTypes: resourceTypeExample,
				InputParams:   inputParam,
				PolicyTags:    policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				WhatAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.WhatAuthorizedWithRawRequest(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("WhatAuthorizedProperty", func() {
			req := &authorizationpb.WhatAuthorizedRequest{
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
				ResourceTypes: resourceTypeExample,
				InputParams:   inputParamInt,
				PolicyTags:    policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				WhatAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.WhatAuthorizedByProperty(
				ctx,
				&authorizationpb.Property{
					Type: "email_for_example",
					Value: &objectpb.Value{
						Value: &objectpb.Value_StringValue{
							StringValue: "test@example.com",
						},
					},
				},
				resourceTypeExample,
				inputParamInt,
				policyTags,
			)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("WhatAuthorizedExternalID", func() {
			externalID := &authorizationpb.ExternalID{
				Type:       "Person",
				ExternalId: "576eb486-28f6-4756-95be-ba6da362b2a7",
			}
			req := &authorizationpb.WhatAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_ExternalId{
						ExternalId: externalID,
					},
				},
				ResourceTypes: resourceTypeExample,
				InputParams:   inputParamTime,
				PolicyTags:    policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				WhatAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.WhatAuthorizedByExternalID(
				ctx,
				externalID,
				resourceTypeExample,
				inputParamTime,
				policyTags,
			)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("WhatAuthorizedToken", func() {
			req := &authorizationpb.WhatAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_AccessToken{
						AccessToken: tokenGood,
					},
				},
				ResourceTypes: resourceTypeExample,
				InputParams:   inputParamArray,
				PolicyTags:    policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				WhatAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.WhatAuthorizedByToken(
				ctx,
				tokenGood,
				resourceTypeExample,
				inputParamArray,
				policyTags,
			)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("WhatAuthorizedTokenWrongFormat", func() {
			req := &authorizationpb.WhatAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_AccessToken{
						AccessToken: tokenBad,
					},
				},
				ResourceTypes: resourceTypeExample,
				InputParams:   inputParamMap,
				PolicyTags:    policyTags,
			}
			resp, err := authorizationClient.WhatAuthorizedWithRawRequest(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("invalid token format cause")))
		})

		It("Invalid status", func() {
			req := &authorizationpb.WhatAuthorizedRequest{
				Subject: &authorizationpb.Subject{
					Subject: &authorizationpb.Subject_AccessToken{
						AccessToken: tokenGood,
					},
				},
				ResourceTypes: resourceTypeExample,
				InputParams:   inputParamMap,
				PolicyTags:    policyTags,
			}
			statusErr := status.New(codes.InvalidArgument, "something wrong").Err()
			mockClient.EXPECT().
				WhatAuthorized(gomock.Any(), req).
				Return(nil, statusErr)

			resp, err := authorizationClient.WhatAuthorizedWithRawRequest(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
})
