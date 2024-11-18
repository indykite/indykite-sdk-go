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

package knowledge_test

import (
	"context"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
	"github.com/indykite/indykite-sdk-go/knowledge"
	knowledgem "github.com/indykite/indykite-sdk-go/test/knowledge/v1beta2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Identity Knowledge API", func() {
	It("Identity Knowledge READ", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

		client, err := knowledge.NewTestClient(mockClient)
		Ω(err).To(Succeed())

		query := "MATCH (n:Person)-[r:CAN_SEE]->(a:Asset) WHERE n.external_id=$external_id AND n.type=$type"
		params := map[string]*objects.Value{
			"external_id": {
				Type: &objects.Value_StringValue{StringValue: "1234"},
			},
			"type": {
				Type: &objects.Value_StringValue{StringValue: "person"},
			},
		}
		returns := []*knowledgepb.Return{
			{
				Variable: "n",
			},
		}

		mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
			Query:       query,
			InputParams: params,
			Returns:     returns,
		}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{

			Nodes: []*knowledgeobjects.Node{
				{
					Id:         "gid:abc",
					ExternalId: "1010",
					Type:       "Person",
					Properties: []*knowledgeobjects.Property{
						{
							Type: "abc",
							Value: &objects.Value{
								Type: &objects.Value_StringValue{StringValue: "something"},
							},
							Metadata: &knowledgeobjects.Metadata{
								AssuranceLevel:   1,
								VerificationTime: timestamppb.Now(),
								Source:           "Myself",
								CustomMetadata: map[string]*objects.Value{
									"customdata": {
										Type: &objects.Value_StringValue{StringValue: "SomeCustomData"},
									},
								},
							},
						},
					},
				},
				{
					Id:         "gid:cba",
					ExternalId: "0101",
					Type:       "Truck",
				},
			},
			Relationships: []*knowledgeobjects.Relationship{
				{
					Id:   "gid:xxx",
					Type: "SERVICES",
				},
			},
		}, nil)

		resp, err := client.IdentityKnowledgeRead(context.Background(), query, params, returns)
		Expect(err).To(Succeed())
		Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Nodes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":         Equal("gid:abc"),
				"ExternalId": Equal("1010"),
				"Type":       Equal("Person"),
				"Properties": HaveLen(1),
			}))),
			"Relationships": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":         Equal("gid:xxx"),
				"Type":       Equal("SERVICES"),
				"Properties": HaveLen(0),
			}))),
		})))
	})

	Context("GetIdentityByID", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:DigitalTwin) WHERE n.id=$id"
			params := map[string]*objects.Value{
				"id": {
					Type: &objects.Value_StringValue{StringValue: "gid:abc"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{

				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:abc",
						ExternalId: "1010",
						Type:       "Person",
						Properties: []*knowledgeobjects.Property{
							{
								Type: "abc",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{StringValue: "something"},
								},
								Metadata: &knowledgeobjects.Metadata{
									AssuranceLevel:   1,
									VerificationTime: timestamppb.Now(),
									Source:           "Myself",
									CustomMetadata: map[string]*objects.Value{
										"customdata": {
											Type: &objects.Value_StringValue{StringValue: "SomeCustomData"},
										},
									},
								},
							},
						},
					},
				},
			}, nil)

			node, err := client.GetIdentityByID(context.Background(), "gid:abc")
			Expect(err).To(Succeed())
			Expect(node).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":         Equal("gid:abc"),
				"ExternalId": Equal("1010"),
				"Type":       Equal("Person"),
				"Properties": HaveLen(1),
			})))
		})

		It("return 0 nodes - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:DigitalTwin) WHERE n.id=$id"
			params := map[string]*objects.Value{
				"id": {
					Type: &objects.Value_StringValue{StringValue: "gid:abc"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{}, nil)

			node, err := client.GetIdentityByID(context.Background(), "gid:abc")
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:DigitalTwin) WHERE n.id=$id"
			params := map[string]*objects.Value{
				"id": {
					Type: &objects.Value_StringValue{StringValue: "gid:abc"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:abc",
						ExternalId: "1010",
						Type:       "Person",
						Properties: []*knowledgeobjects.Property{
							{
								Type: "abc",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{StringValue: "something"},
								},
							},
						},
					},
					{
						Id:         "gid:cba",
						ExternalId: "0101",
						Type:       "Person",
					},
				},
			}, nil)

			node, err := client.GetIdentityByID(context.Background(), "gid:abc")
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("GetIdentityByIdentifier", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:DigitalTwin) WHERE n.external_id=$external_id AND n.type=$type"
			params := map[string]*objects.Value{
				"external_id": {
					Type: &objects.Value_StringValue{StringValue: "1010"},
				},
				"type": {
					Type: &objects.Value_StringValue{StringValue: "Person"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:abc",
						ExternalId: "1010",
						Type:       "Person",
						Properties: []*knowledgeobjects.Property{
							{
								Type: "abc",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{StringValue: "something"},
								},
							},
						},
					},
				},
			}, nil)

			node, err := client.GetIdentityByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "1010",
				Type:       "Person",
			},
			)
			Expect(err).To(Succeed())
			Expect(node).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":         Equal("gid:abc"),
				"ExternalId": Equal("1010"),
				"Type":       Equal("Person"),
				"Properties": HaveLen(1),
			})))
		})

		It("return 0 nodes - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:DigitalTwin) WHERE n.external_id=$external_id AND n.type=$type"
			params := map[string]*objects.Value{
				"external_id": {
					Type: &objects.Value_StringValue{StringValue: "1010"},
				},
				"type": {
					Type: &objects.Value_StringValue{StringValue: "Person"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{}, nil)

			node, err := client.GetIdentityByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "1010",
				Type:       "Person",
			})
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:DigitalTwin) WHERE n.external_id=$external_id AND n.type=$type"
			params := map[string]*objects.Value{
				"external_id": {
					Type: &objects.Value_StringValue{StringValue: "1010"},
				},
				"type": {
					Type: &objects.Value_StringValue{StringValue: "Person"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:abc",
						ExternalId: "1010",
						Type:       "Person",
						Properties: []*knowledgeobjects.Property{
							{
								Type: "abc",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{StringValue: "something"},
								},
							},
						},
					},
					{
						Id:         "gid:cba",
						ExternalId: "0101",
						Type:       "Person",
					},
				},
			}, nil)

			node, err := client.GetIdentityByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "1010",
				Type:       "Person",
			})
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("GetNodeByID", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource) WHERE n.id=$id"
			params := map[string]*objects.Value{
				"id": {
					Type: &objects.Value_StringValue{StringValue: "gid:xyz"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:xyz",
						ExternalId: "0000",
						Type:       "Store",
					},
				},
			}, nil)

			node, err := client.GetNodeByID(context.Background(), "gid:xyz", false)
			Expect(err).To(Succeed())
			Expect(node).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":         Equal("gid:xyz"),
				"ExternalId": Equal("0000"),
				"Type":       Equal("Store"),
			})))
		})

		It("return 0 nodes - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource) WHERE n.id=$id"
			params := map[string]*objects.Value{
				"id": {
					Type: &objects.Value_StringValue{StringValue: "gid:xyz"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{}, nil)

			node, err := client.GetNodeByID(context.Background(), "gid:xyz", false)
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource) WHERE n.id=$id"
			params := map[string]*objects.Value{
				"id": {
					Type: &objects.Value_StringValue{StringValue: "gid:xyz"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:xyz",
						ExternalId: "0000",
						Type:       "Store",
					},
					{
						Id:         "gid:cba",
						ExternalId: "0101",
						Type:       "Person",
					},
				},
			}, nil)

			node, err := client.GetNodeByID(context.Background(), "gid:xyz", false)
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("GetNodeByIdentifier", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource) WHERE n.external_id=$external_id AND n.type=$type"
			params := map[string]*objects.Value{
				"external_id": {
					Type: &objects.Value_StringValue{StringValue: "0000"},
				},
				"type": {
					Type: &objects.Value_StringValue{StringValue: "Store"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:xyz",
						ExternalId: "0000",
						Type:       "Store",
					},
				},
			}, nil)

			node, err := client.GetNodeByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "0000",
				Type:       "Store",
			},
				false)
			Expect(err).To(Succeed())
			Expect(node).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Id":         Equal("gid:xyz"),
				"ExternalId": Equal("0000"),
				"Type":       Equal("Store"),
			})))
		})

		It("return 0 nodes - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource) WHERE n.external_id=$external_id AND n.type=$type"
			params := map[string]*objects.Value{
				"external_id": {
					Type: &objects.Value_StringValue{StringValue: "0000"},
				},
				"type": {
					Type: &objects.Value_StringValue{StringValue: "Store"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{}, nil)

			node, err := client.GetNodeByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "0000",
				Type:       "Store",
			},
				false)
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())
			query := "MATCH (n:Resource) WHERE n.external_id=$external_id AND n.type=$type"
			params := map[string]*objects.Value{
				"external_id": {
					Type: &objects.Value_StringValue{StringValue: "0000"},
				},
				"type": {
					Type: &objects.Value_StringValue{StringValue: "Store"},
				},
			}
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: params,
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{

				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:xyz",
						ExternalId: "0000",
						Type:       "Store",
					},
					{
						Id:         "gid:cba",
						ExternalId: "0101",
						Type:       "Person",
					},
				},
			}, nil)

			node, err := client.GetNodeByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "0000",
				Type:       "Store",
			},
				false)
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("List all", func() {
		It("ListResources - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource)"
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: map[string]*objects.Value{},
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:xyz",
						ExternalId: "0000",
						Type:       "Store",
					},
					{
						Id:         "gid:xxx",
						ExternalId: "0001",
						Type:       "Store",
					},
					{
						Id:         "gid:yyyy",
						ExternalId: "0002",
						Type:       "Product",
					},
				},
			}, nil)

			nodes, err := client.ListNodes(context.Background(), "Resource")
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(3))
		})

		It("ListIdentities - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}

			query := "MATCH (n:DigitalTwin)"
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: map[string]*objects.Value{},
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:abc",
						ExternalId: "1010",
						Type:       "Person",
					},
					{
						Id:         "gid:aaa",
						ExternalId: "0101",
						Type:       "Vehicle",
					},
				},
			}, nil)

			nodes, err := client.ListIdentities(context.Background())
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(2))
		})

		It("ListNodes - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			query := "MATCH (n:Resource)"
			returns := []*knowledgepb.Return{
				{
					Variable: "n",
				},
			}
			mockClient.EXPECT().IdentityKnowledgeRead(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeReadRequest{
				Query:       query,
				InputParams: map[string]*objects.Value{},
				Returns:     returns,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeReadResponse{
				Nodes: []*knowledgeobjects.Node{
					{
						Id:         "gid:abc",
						ExternalId: "1010",
						Type:       "Person",
					},
				},
			}, nil)

			nodes, err := client.ListNodes(context.Background(), "Resource")
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(1))
		})
	})
})
