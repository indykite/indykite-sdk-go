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

package ingest_test

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Ingestion", func() {
	Describe("IngestBatchUpsertNodes", func() {
		It("UpsertNodes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			nodes, externalID, externalID2 := integration.BatchNodesType("Individual")
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
			id := result.Changes[0].Id

			nodesMatch := integration.BatchNodesMatch(
				integration.CreateBatchNodeMatch(externalID, "Individual"),
				integration.CreateBatchNodeMatch(externalID2, "Individual"))
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Equal(id),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("UpsertNodeDuplicateSameType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			node1, externalID := integration.CreateBatchNodes("Individual")
			nodes := []*knowledgeobjects.Node{
				node1, node1,
			}
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
			id := result.Changes[0].Id

			Expect(err).To(Succeed())
			result2 := resp.Results[1]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
			id2 := result2.Changes[0].Id
			Expect(id2).To(Equal(id))

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Individual"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Equal(id),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("UpsertNodeDuplicateDifferentTypes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			node1, externalID := integration.CreateBatchNodes("Individual")
			node2 := &knowledgeobjects.Node{
				ExternalId: externalID,
				Type:       "Cat",
				IsIdentity: true,
			}

			nodes := []*knowledgeobjects.Node{
				node1, node2,
			}
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
			id := result.Changes[0].Id

			Expect(err).To(Succeed())
			result2 := resp.Results[1]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
			id2 := result2.Changes[0].Id
			Expect(id2).NotTo(Equal(id))

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Individual"),
				integration.CreateBatchNodeMatch(externalID, "Cat"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Equal(id),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("UpsertNodeBothValueAndExternalValue", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			node1 := integration.CreateBatchNodesError("Individual")

			nodes := []*knowledgeobjects.Node{
				node1,
			}
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(MatchError(ContainSubstring(
				"externalValue and value are mutually exclusive")))
			Expect(resp).To(BeNil())
		})

		It("UpsertNodeExternalValueWithoutResolver", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			node1 := integration.CreateBatchNodesNoResolver("Individual")

			nodes := []*knowledgeobjects.Node{
				node1,
			}
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(MatchError(ContainSubstring(
				"invalid external value GID, expected GID of type ExternalDataResolver")))
			Expect(resp).To(BeNil())
		})
	})

	Describe("IngestBatchUpsertRelationships", func() {
		It("BatchUpsertRelationships", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			nodes, externalID, externalID2 := integration.BatchNodesType("Asset")
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			result2 := resp.Results[1]
			Expect(result2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			nodes10, externalID3, externalID4 := integration.BatchNodesType("Asset")
			resp10, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes10,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result10 := resp10.Results[0]
			Expect(result10).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			relationships := integration.BatchRelationships(
				integration.GetRelationship(externalID, "Asset", externalID2, "Asset", "CAN_SEE"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "CAN_SEE"),
			)
			resp3, err := ingestClient.BatchUpsertRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result3 := resp3.Results[0]
			Expect(result3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))
			id3 := result3.Changes[0].Id

			del3, err := ingestClient.BatchDeleteRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			resultDel3 := del3.Results[0]
			Expect(resultDel3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Equal(id3),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Asset"),
				integration.CreateBatchNodeMatch(externalID2, "Asset"),
				integration.CreateBatchNodeMatch(externalID3, "Asset"),
				integration.CreateBatchNodeMatch(externalID4, "Asset"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("BatchUpsertRelationshipWrongSourceType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			nodes, externalID, externalID2 := integration.BatchNodesType("Asset")
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			result2 := resp.Results[1]
			Expect(result2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			nodes10, externalID3, externalID4 := integration.BatchNodesType("Asset")
			resp10, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes10,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result10 := resp10.Results[0]
			Expect(result10).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			relationships := integration.BatchRelationships(
				integration.GetRelationship(externalID, "Service", externalID2, "Service", "CAN_SEE"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "CAN_SEE"),
			)
			resp3, err := ingestClient.BatchUpsertRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result3 := resp3.Results[0]
			Expect(result3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))
			id3 := result3.Changes[0].Id

			del3, err := ingestClient.BatchDeleteRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			resultDel3 := del3.Results[0]
			Expect(resultDel3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Equal(id3),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Asset"),
				integration.CreateBatchNodeMatch(externalID2, "Asset"),
				integration.CreateBatchNodeMatch(externalID3, "Asset"),
				integration.CreateBatchNodeMatch(externalID4, "Asset"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("BatchDeleteRelationshipWrongSourceType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			nodes, externalID, externalID2 := integration.BatchNodesType("Asset")
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			result2 := resp.Results[1]
			Expect(result2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			nodes10, externalID3, externalID4 := integration.BatchNodesType("Asset")
			resp10, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes10,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result10 := resp10.Results[0]
			Expect(result10).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			relationships := integration.BatchRelationships(
				integration.GetRelationship(externalID, "Asset", externalID2, "Asset", "CAN_SEE"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "CAN_SEE"),
			)
			resp3, err := ingestClient.BatchUpsertRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result3 := resp3.Results[0]
			Expect(result3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))

			relationships2 := integration.BatchRelationships(
				integration.GetRelationship(externalID, "Service", externalID2, "Service", "CAN_SEE"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "CAN_SEE"),
			)
			del3, err := ingestClient.BatchDeleteRelationships(
				context.Background(),
				relationships2,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			resultDel3 := del3.Results[0]
			Expect(resultDel3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Asset"),
				integration.CreateBatchNodeMatch(externalID2, "Asset"),
				integration.CreateBatchNodeMatch(externalID3, "Asset"),
				integration.CreateBatchNodeMatch(externalID4, "Asset"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("BatchUpsertRelationshipWrongSource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			nodes, externalID, externalID2 := integration.BatchNodesType("Asset")
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			result2 := resp.Results[1]
			Expect(result2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			nodes10, externalID3, externalID4 := integration.BatchNodesType("Asset")
			resp10, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes10,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result10 := resp10.Results[0]
			Expect(result10).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			rel := integration.GenerateRandomString(10)
			relationships := integration.BatchRelationships(
				integration.GetRelationship(rel, "Asset", externalID2, "Asset", "CAN_SEE"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "CAN_SEE"),
			)
			resp3, err := ingestClient.BatchUpsertRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result3 := resp3.Results[0]
			Expect(result3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))
			id3 := result3.Changes[0].Id

			relationships2 := integration.BatchRelationships(
				integration.GetRelationship(rel, "Asset", externalID2, "Asset", "CAN_SEE"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "CAN_SEE"),
			)
			del3, err := ingestClient.BatchDeleteRelationships(
				context.Background(),
				relationships2,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			resultDel3 := del3.Results[0]
			Expect(resultDel3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Equal(id3),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
				}))),
			})))

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Asset"),
				integration.CreateBatchNodeMatch(externalID2, "Asset"),
				integration.CreateBatchNodeMatch(externalID3, "Asset"),
				integration.CreateBatchNodeMatch(externalID4, "Asset"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})

		It("BatchUpsertRelationshipActionLowercase", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			nodes, externalID, externalID2 := integration.BatchNodesType("Asset")
			resp, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result := resp.Results[0]
			Expect(result).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			result2 := resp.Results[1]
			Expect(result2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			nodes10, externalID3, externalID4 := integration.BatchNodesType("Asset")
			resp10, err := ingestClient.BatchUpsertNodes(
				context.Background(),
				nodes10,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			result10 := resp10.Results[0]
			Expect(result10).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))

			relationships := integration.BatchRelationships(
				integration.GetRelationship(externalID, "Asset", externalID2, "Asset", "can_see"),
				integration.GetRelationship(externalID3, "Asset", externalID4, "Asset", "can_see"),
			)
			resp3, err := ingestClient.BatchUpsertRelationships(
				context.Background(),
				relationships,
				retry.WithMax(5),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid Relationship.Type: value does not match regex pattern \"^[A-Z]+(?:_[A-Z]+)*$\"")))
			Expect(resp3).To(BeNil())

			nodesMatch := []*ingestpb.NodeMatch{
				integration.CreateBatchNodeMatch(externalID, "Asset"),
				integration.CreateBatchNodeMatch(externalID2, "Asset"),
				integration.CreateBatchNodeMatch(externalID3, "Asset"),
				integration.CreateBatchNodeMatch(externalID4, "Asset"),
			}
			del, err := ingestClient.BatchDeleteNodes(
				context.Background(),
				nodesMatch,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			resultDel := del.Results[0]
			Expect(resultDel).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
					"Id":       Not(BeEmpty()),
					"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
				}))),
			})))
		})
	})
})
