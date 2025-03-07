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

package knowledge_test

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
	"github.com/indykite/indykite-sdk-go/knowledge"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Knowledge", func() {
	Describe("KnowledgeApi", func() {
		DescribeTable("Identity Knowledge READ",
			func(query string,
				params map[string]*objects.Value,
				returns []*knowledgepb.Return,
				expectedError string,
				matcher Fields) {
				knowledgeClient, err := integration.InitConfigKnowledge()
				Expect(err).To(Succeed())

				resp, err := knowledgeClient.IdentityKnowledgeRead(
					context.Background(),
					query,
					params,
					returns,
					retry.WithMax(5),
				)

				if expectedError != "" {
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, matcher)))
				}
			},
			Entry("Read single subject", integration.Query1, integration.Params1,
				integration.Returns1, "", integration.Matcher1),
			Entry("Read single subject with external", integration.Query2, integration.Params1,
				integration.Returns1, "", integration.Matcher1),
			Entry("Read all subjects", integration.Query3, integration.Params3,
				integration.Returns1, "", integration.Matcher1),
			Entry("Read all subjects with external", integration.Query4, integration.Params3,
				integration.Returns1, "", integration.Matcher1),
			Entry("Read with query on property", integration.Query5, integration.Params4,
				integration.Returns1, "", integration.Matcher2),
			Entry("Read with query on external property", integration.Query6, integration.Params5,
				integration.Returns1, "", integration.Matcher2),
			Entry("Read single subject error", integration.Query7, integration.Params3,
				integration.Returns1, "invalid cypher syntax", integration.Matcher4),
			Entry("Read trust score by node label", integration.Query8, integration.Params3,
				integration.Returns2, "", integration.Matcher5),
			Entry("Read match the Trust Score profile name", integration.Query9, integration.Params6,
				integration.Returns3, "", integration.Matcher6),
			Entry("Read specific dimensions from all Trust Scores", integration.Query10, integration.Params3,
				integration.Returns4, "", integration.Matcher7),
		)

		DescribeTable("GetByID",
			func(id string,
				isIdentity bool,
				expectedError string,
				matcher Fields) {
				knowledgeClient, err := integration.InitConfigKnowledge()
				Expect(err).To(Succeed())

				resp, err := knowledgeClient.GetNodeByID(
					context.Background(),
					id,
					isIdentity,
					retry.WithMax(5),
				)

				if expectedError != "" {
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, matcher)))
				}
			},
			Entry("Get by id resource", integration.Truck1Id, false, "", integration.Matcher3),
			Entry("Get by id identity", integration.Node2, true, "", integration.Matcher4),
			Entry("Get by id identity error", integration.NodeNotInDB, true, "node not found", integration.Matcher3),
		)

		DescribeTable("GetByIdentifier",
			func(externalId string,
				typeNode string,
				isIdentity bool,
				expectedError string,
				matcher Fields) {
				knowledgeClient, err := integration.InitConfigKnowledge()
				Expect(err).To(Succeed())

				resp, err := knowledgeClient.GetNodeByIdentifier(
					context.Background(),
					&knowledge.Identifier{
						ExternalID: externalId,
						Type:       typeNode,
					},
					isIdentity,
					retry.WithMax(5),
				)

				if expectedError != "" {
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, matcher)))
				}
			},
			Entry("Get node by identifier", integration.Subject2, "Person", true, "", integration.Matcher4),
			Entry("Get node with external", integration.Truck1, "Truck", false, "", integration.Matcher3),
			Entry("Get node error", "123456", "Person", true, "node not found", integration.Matcher1),
		)

		DescribeTable("GetAll",
			func(typeNode string) {
				knowledgeClient, err := integration.InitConfigKnowledge()
				Expect(err).To(Succeed())

				resp, err := knowledgeClient.ListNodes(
					context.Background(),
					typeNode,
					retry.WithMax(5),
				)
				Expect(err).To(Succeed())
				Expect(len(resp)).Should(BeNumerically(">=", 2))
			},
			Entry("Get all nodes", "Truck"),
			Entry("Get all identities", "Person"),
		)
	})
})
