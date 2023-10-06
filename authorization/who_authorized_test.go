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
	"github.com/indykite/indykite-sdk-go/test"
	authorizationmock "github.com/indykite/indykite-sdk-go/test/authorization/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WhoAuthorized", func() {
	var (
		ctx                 context.Context
		mockCtrl            *gomock.Controller
		mockClient          *authorizationmock.MockAuthorizationAPIClient
		authorizationClient *authorization.Client
		resourceExample     = []*authorizationpb.WhoAuthorizedRequest_Resource{
			{ExternalId: "external_id_value", Type: "Asset", Actions: []string{"ACTION1", "ACTION2"}},
		}
		inputParam = map[string]*authorizationpb.InputParam{
			"Color": {Value: &authorizationpb.InputParam_StringValue{StringValue: "red"}},
		}
		policyTags = []string{"sometag"}
		beResp     = &authorizationpb.WhoAuthorizedResponse{
			DecisionTime: timestamppb.New(time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)),
			Decisions: map[string]*authorizationpb.WhoAuthorizedResponse_ResourceType{
				"ResourceTypeValue": {
					Resources: map[string]*authorizationpb.WhoAuthorizedResponse_Resource{
						"Resource": {
							Actions: map[string]*authorizationpb.WhoAuthorizedResponse_Action{
								"ACTION1": {
									Subjects: []*authorizationpb.WhoAuthorizedResponse_Subject{
										{ExternalId: "externalId1"},
										{ExternalId: "externalId2"},
									},
								},
								"ACTION2": {
									Subjects: []*authorizationpb.WhoAuthorizedResponse_Subject{
										{ExternalId: "externalId1"},
										{ExternalId: "externalId2"},
									},
								},
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

	Describe("WhoAuthorized", func() {
		It("Nil params", func() {
			req := &authorizationpb.WhoAuthorizedRequest{
				Resources: []*authorizationpb.WhoAuthorizedRequest_Resource{
					{},
				},
			}
			resp, err := authorizationClient.WhoAuthorized(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))
		})

		It("Wrong resource should return a validation error in the response", func() {
			req := &authorizationpb.WhoAuthorizedRequest{
				Resources: []*authorizationpb.WhoAuthorizedRequest_Resource{
					{Type: "Asset", Actions: []string{"ACTION1", "ACTION2"}},
				},
			}

			resp, err := authorizationClient.WhoAuthorized(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("invalid WhoAuthorizedRequest_Resource.ExternalId")))
		})

		It("WhoAuthorized", func() {
			req := &authorizationpb.WhoAuthorizedRequest{
				Resources:   resourceExample,
				InputParams: inputParam,
				PolicyTags:  policyTags,
			}
			beResp := beResp
			mockClient.EXPECT().
				WhoAuthorized(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

			resp, err := authorizationClient.WhoAuthorized(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp).To(test.EqualProto(beResp))
		})

		It("Invalid status", func() {
			statusErr := status.New(codes.InvalidArgument, "something wrong").Err()
			mockClient.EXPECT().
				WhoAuthorized(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, statusErr)

			resp, err := authorizationClient.WhoAuthorized(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})
})
