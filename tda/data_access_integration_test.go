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

//go:build integration

package tda_test

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	tdapb "github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("TDA", func() {
	Describe("GrantConsent", func() {
		It("GrantConsentByIdSuccess", func() {
			var err error
			tdaClient, err := integration.InitConfigTda()
			Expect(err).To(Succeed())

			resp, err := tdaClient.GrantConsent(
				context.Background(),
				&tdapb.GrantConsentRequest{
					User: &knowledgeobjects.User{
						User: &knowledgeobjects.User_UserId{
							UserId: integration.DigitalTwin4,
						},
					},
					ConsentId:      integration.ConsentConfig2,
					ValidityPeriod: uint64(86400),
				},
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"PropertiesGrantedCount": Equal(uint64(1)),
			})))

			tdaRevokeRequest := tdapb.RevokeConsentRequest{
				User: &knowledgeobjects.User{
					User: &knowledgeobjects.User_UserId{
						UserId: integration.DigitalTwin4,
					},
				},
				ConsentId: integration.ConsentConfigID,
			}
			resp2, err := tdaClient.RevokeConsent(
				context.Background(),
				&tdaRevokeRequest,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).ToNot(BeNil())
		})

		It("GrantConsentByExternalIdSuccess", func() {
			var err error
			tdaClient, err := integration.InitConfigTda()
			Expect(err).To(Succeed())

			tdaRequest := tdapb.GrantConsentRequest{
				User: &knowledgeobjects.User{
					User: &knowledgeobjects.User_ExternalId{
						ExternalId: &knowledgeobjects.User_ExternalID{
							Type:       "Person",
							ExternalId: integration.SubjectDT4,
						},
					},
				},
				ConsentId:      integration.ConsentConfig2,
				ValidityPeriod: uint64(86400),
			}
			resp, err := tdaClient.GrantConsent(
				context.Background(),
				&tdaRequest,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"PropertiesGrantedCount": Equal(uint64(1)),
			})))

			tdaRevokeRequest := tdapb.RevokeConsentRequest{
				User: &knowledgeobjects.User{
					User: &knowledgeobjects.User_ExternalId{
						ExternalId: &knowledgeobjects.User_ExternalID{
							Type:       "Person",
							ExternalId: integration.SubjectDT4,
						},
					},
				},
				ConsentId: integration.ConsentConfig2,
			}
			resp2, err := tdaClient.RevokeConsent(
				context.Background(),
				&tdaRevokeRequest,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).ToNot(BeNil())
		})
	})
})
