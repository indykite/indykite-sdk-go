// Copyright (c) 2024 IndyKite
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

package entitymatching_test

import (
	"context"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indykite/indykite-sdk-go/entitymatching"
	entitymatchingpb "github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	entitymatchingmock "github.com/indykite/indykite-sdk-go/test/entitymatching/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EntityMatching", func() {
	mockErrorCode := codes.InvalidArgument
	mockErrorMessage := "mockError"
	now := time.Now()

	var (
		mockCtrl             *gomock.Controller
		mockClient           *entitymatchingmock.MockEntityMatchingAPIClient
		entitymatchingClient *entitymatching.Client
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = entitymatchingmock.NewMockEntityMatchingAPIClient(mockCtrl)
		var err error
		entitymatchingClient, err = entitymatching.NewTestClient(mockClient)
		Ω(err).To(Succeed())
	})

	Describe("Read Suggested Property Mapping", func() {
		mockReadSuggestedPropertyMappingRequest := entitymatchingpb.ReadSuggestedPropertyMappingRequest{
			Id: "gid:886hfic8fswlz3zjrc2e3nun9xs",
		}

		DescribeTable("ReadSuggestedPropertyMappingSuccess",
			func(req *entitymatchingpb.ReadSuggestedPropertyMappingRequest,
				beResp *entitymatchingpb.ReadSuggestedPropertyMappingResponse) {
				mockClient.EXPECT().ReadSuggestedPropertyMapping(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

				resp, err := entitymatchingClient.ReadSuggestedPropertyMapping(context.Background(), req)
				Ω(resp).ToNot(BeNil())
				Ω(err).To(Succeed())
				Ω(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"read suggested property mapping id and returns response",
				&mockReadSuggestedPropertyMappingRequest,
				&entitymatchingpb.ReadSuggestedPropertyMappingResponse{
					Id: "gid:886hfic8fswlz3zjrc2e3nun9xs",
					PropertyMappings: []*entitymatchingpb.PropertyMapping{
						{
							SourceNodeType:        "employee",
							SourceNodeProperty:    "email",
							TargetNodeType:        "user",
							TargetNodeProperty:    "address",
							SimilarityScoreCutoff: 0.9,
						},
					},
					PropertyMappingStatus: entitymatchingpb.PipelineStatus_PIPELINE_STATUS_STATUS_SUCCESS,
				},
			),
		)

		It("handles and returns error", func() {
			mockClient.EXPECT().ReadSuggestedPropertyMapping(
				gomock.Any(),
				gomock.Eq(&mockReadSuggestedPropertyMappingRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := entitymatchingClient.ReadSuggestedPropertyMapping(
				context.Background(),
				&mockReadSuggestedPropertyMappingRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})

	Describe("Run Entity Matching Pipeline", func() {
		mockRunEntityMatchingPipelineRequest := entitymatchingpb.RunEntityMatchingPipelineRequest{
			Id: "gid:886hfic8fswlz3zjrc2e3nun9xs",
			CustomPropertyMappings: []*entitymatchingpb.PropertyMapping{
				{
					SourceNodeType:        "employee",
					SourceNodeProperty:    "email",
					TargetNodeType:        "user",
					TargetNodeProperty:    "address",
					SimilarityScoreCutoff: 0.9,
				},
			},
			SimilarityScoreCutoff: 0.95,
		}

		DescribeTable("RunEntityMatchingPipelineSuccess",
			func(req *entitymatchingpb.RunEntityMatchingPipelineRequest,
				beResp *entitymatchingpb.RunEntityMatchingPipelineResponse) {
				mockClient.EXPECT().RunEntityMatchingPipeline(
					gomock.Any(),
					test.WrapMatcher(test.EqualProto(req)),
					gomock.Any(),
				).Return(beResp, nil)

				resp, err := entitymatchingClient.RunEntityMatchingPipeline(context.Background(), req)
				Ω(resp).ToNot(BeNil())
				Ω(err).To(Succeed())
				Ω(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"run entty matching mapping id and returns response",
				&mockRunEntityMatchingPipelineRequest,
				&entitymatchingpb.RunEntityMatchingPipelineResponse{
					Id:          "gid:886hfic8fswlz3zjrc2e3nun9xs",
					LastRunTime: timestamppb.New(now),
					Etag:        "tjNkmiLko",
				},
			),
		)

		It("handles and returns error", func() {
			mockClient.EXPECT().RunEntityMatchingPipeline(
				gomock.Any(),
				gomock.Eq(&mockRunEntityMatchingPipelineRequest),
				gomock.Any(),
			).Return(nil, status.Error(mockErrorCode, mockErrorMessage))

			resp, err := entitymatchingClient.RunEntityMatchingPipeline(
				context.Background(),
				&mockRunEntityMatchingPipelineRequest)
			Ω(resp).To(BeNil())
			Ω(err).To(HaveOccurred())
			Ω(err).To(test.MatchStatusError(mockErrorCode, mockErrorMessage))
		})
	})
})
