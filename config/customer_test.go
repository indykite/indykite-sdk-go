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
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indykite/indykite-sdk-go/config"
	sdkerrors "github.com/indykite/indykite-sdk-go/errors"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	"github.com/indykite/indykite-sdk-go/test"
	configmock "github.com/indykite/indykite-sdk-go/test/config/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Customer", func() {
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

	Describe("Customer", func() {
		It("Nil request", func() {
			resp, err := configClient.ReadCustomer(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			var clientErr *sdkerrors.ClientError
			Expect(errors.As(err, &clientErr)).To(BeTrue(), "is client error")
			Expect(clientErr.Unwrap()).To(Succeed())
			Expect(clientErr.Message()).To(Equal("invalid nil request"))
			Expect(clientErr.Code()).To(Equal(codes.InvalidArgument))

		})

		It("Wrong id should return a validation error in the response", func() {
			req := &configpb.ReadCustomerRequest{
				Identifier: &configpb.ReadCustomerRequest_Id{Id: "gid:like"},
			}
			resp, err := configClient.ReadCustomer(ctx, req)
			Expect(resp).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("Id: value length must be between 22")))
		})

		DescribeTable("ReadSuccess",
			func(req *configpb.ReadCustomerRequest, beResp *configpb.ReadCustomerResponse) {
				mockClient.EXPECT().
					ReadCustomer(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, nil)

				resp, err := configClient.ReadCustomer(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp).To(test.EqualProto(beResp))
			},
			Entry(
				"ReadId",
				&configpb.ReadCustomerRequest{
					Identifier: &configpb.ReadCustomerRequest_Id{Id: "gid:like-real-customer-id"},
				},
				&configpb.ReadCustomerResponse{
					Customer: &configpb.Customer{
						Id:          "gid:like-real-customer-id",
						Name:        "like-real-customer-name",
						DisplayName: "Like Real Customer Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
					},
				},
			),
			Entry(
				"ReadName",
				&configpb.ReadCustomerRequest{
					Identifier: &configpb.ReadCustomerRequest_Name{Name: "like-real-customer-name"},
				},
				&configpb.ReadCustomerResponse{
					Customer: &configpb.Customer{
						Id:          "gid:like-real-customer-id",
						Name:        "like-real-customer-name",
						DisplayName: "Like Real Customer Name",
						CreatedBy:   "creator",
						CreateTime:  timestamppb.Now(),
					},
				},
			),
		)

		DescribeTable("ReadError",
			func(req *configpb.ReadCustomerRequest, beResp *configpb.ReadCustomerResponse) {
				mockClient.EXPECT().
					ReadCustomer(
						gomock.Any(),
						test.WrapMatcher(test.EqualProto(req)),
						gomock.Any(),
					).Return(beResp, status.Error(codes.InvalidArgument, "status error"))

				resp, err := configClient.ReadCustomer(ctx, req)
				Expect(err).ToNot(Succeed())
				Expect(resp).To(BeNil())
			},
			Entry(
				"ReadId",
				&configpb.ReadCustomerRequest{
					Identifier: &configpb.ReadCustomerRequest_Id{Id: "gid:like-real-customer-id"},
				},
				&configpb.ReadCustomerResponse{},
			),
			Entry(
				"ReadName",
				&configpb.ReadCustomerRequest{
					Identifier: &configpb.ReadCustomerRequest_Name{Name: "like-real-customer-name"},
				},
				&configpb.ReadCustomerResponse{},
			),
		)
	})
})
