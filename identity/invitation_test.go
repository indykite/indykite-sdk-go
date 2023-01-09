// Copyright (c) 2022 IndyKite
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
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
	"github.com/indykite/jarvis-sdk-go/identity"
	midentity "github.com/indykite/jarvis-sdk-go/test/identity/v1beta2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Invitation", func() {
	It("Create", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		mockClient := midentity.NewMockIdentityManagementAPIClient(mockCtrl)

		client, err := identity.NewTestClient(mockClient)
		Ω(err).To(Succeed())

		tenantID := "gid:AAAAA2luZHlraURlgAADDwAAAAE"
		Ω(err).To(Succeed())

		now := time.Now().Add(time.Minute)

		mockClient.EXPECT().CreateInvitation(gomock.Any(), gomock.Eq(&identitypb.CreateInvitationRequest{
			TenantId:     tenantID,
			ReferenceId:  "my-reference",
			InviteAtTime: timestamppb.New(now),
			ExpireTime:   timestamppb.New(now.AddDate(0, 0, 7)),
			Invitee:      &identitypb.CreateInvitationRequest_Email{Email: "test@example.com"},
			MessageAttributes: &objects.MapValue{
				Fields: map[string]*objects.Value{"lang": objects.String("en")},
			},
		}), gomock.Any()).Return(&identitypb.CreateInvitationResponse{}, nil)

		err = client.CreateEmailInvitation(context.Background(),
			"test@example.com",
			"gid:AAAAA2luZHlraURlgAADDwAAAAE",
			"my-reference",
			now.AddDate(0, 0, 7), now,
			map[string]interface{}{
				"lang": "en",
			})
		Ω(err).To(Succeed())
	})
})
