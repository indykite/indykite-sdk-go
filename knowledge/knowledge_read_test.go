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

	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
	"github.com/indykite/indykite-sdk-go/knowledge"
	knowledgem "github.com/indykite/indykite-sdk-go/test/knowledge/v1beta1"

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

		p := "(:DigitalTwin)-[:SERVICES]->(n:Truck)"
		c := "WHERE n.external_id = $external_id"
		params := map[string]*knowledgepb.InputParam{
			"external_id": {
				Value: &knowledgepb.InputParam_StringValue{StringValue: "1234"},
			},
		}
		mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
			Operation:   knowledgepb.Operation_OPERATION_READ,
			Path:        p,
			Conditions:  c,
			InputParams: params,
		}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
			Paths: []*knowledgepb.Path{
				{
					Nodes: []*knowledgepb.Node{
						{
							Id:         "gid:abc",
							ExternalId: "1010",
							Type:       "Person",
							Properties: []*knowledgepb.Property{
								{
									Key: "abc",
									Value: &objects.Value{
										Value: &objects.Value_StringValue{StringValue: "something"},
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
					Relationships: []*knowledgepb.Relationship{
						{
							Id:   "gid:xxx",
							Type: "SERVICES",
						},
					},
				},
			},
		}, nil)

		resp, err := client.Read(context.Background(), p, c, params)
		Expect(err).To(Succeed())
		Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Paths": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
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
			}))),
		})))
	})

	Context("GetDigitalTwinByID", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:DigitalTwin)"
			c := "WHERE n.id = $id"
			params := map[string]*knowledgepb.InputParam{
				"id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "gid:abc"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
								Properties: []*knowledgepb.Property{
									{
										Key: "abc",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{StringValue: "something"},
										},
									},
								},
							},
						},
					},
				},
			}, nil)

			node, err := client.GetDigitalTwinByID(context.Background(), "gid:abc")
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

			p := "(n:DigitalTwin)"
			c := "WHERE n.id = $id"
			params := map[string]*knowledgepb.InputParam{
				"id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "gid:abc"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{},
			}, nil)

			node, err := client.GetDigitalTwinByID(context.Background(), "gid:abc")
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:DigitalTwin)"
			c := "WHERE n.id = $id"
			params := map[string]*knowledgepb.InputParam{
				"id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "gid:abc"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
								Properties: []*knowledgepb.Property{
									{
										Key: "abc",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{StringValue: "something"},
										},
									},
								},
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:cba",
								ExternalId: "0101",
								Type:       "Person",
							},
						},
					},
				},
			}, nil)

			node, err := client.GetDigitalTwinByID(context.Background(), "gid:abc")
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("GetDigitalTwinByIdentifier", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:DigitalTwin)"
			c := "WHERE n.external_id = $external_id AND n.type = $type"
			params := map[string]*knowledgepb.InputParam{
				"external_id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "1010"},
				},
				"type": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "person"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
								Properties: []*knowledgepb.Property{
									{
										Key: "abc",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{StringValue: "something"},
										},
									},
								},
							},
						},
					},
				},
			}, nil)

			node, err := client.GetDigitalTwinByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "1010",
				Type:       "Person",
			})
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

			p := "(n:DigitalTwin)"
			c := "WHERE n.external_id = $external_id AND n.type = $type"
			params := map[string]*knowledgepb.InputParam{
				"external_id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "1010"},
				},
				"type": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "person"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{},
			}, nil)

			node, err := client.GetDigitalTwinByIdentifier(context.Background(), &knowledge.Identifier{
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

			p := "(n:DigitalTwin)"
			c := "WHERE n.external_id = $external_id AND n.type = $type"
			params := map[string]*knowledgepb.InputParam{
				"external_id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "1010"},
				},
				"type": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "person"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
								Properties: []*knowledgepb.Property{
									{
										Key: "abc",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{StringValue: "something"},
										},
									},
								},
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:cba",
								ExternalId: "0101",
								Type:       "Person",
							},
						},
					},
				},
			}, nil)

			node, err := client.GetDigitalTwinByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "1010",
				Type:       "Person",
			})
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("GetResourceByID", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:Resource)"
			c := "WHERE n.id = $id"
			params := map[string]*knowledgepb.InputParam{
				"id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "gid:xyz"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xyz",
								ExternalId: "0000",
								Type:       "Store",
							},
						},
					},
				},
			}, nil)

			node, err := client.GetResourceByID(context.Background(), "gid:xyz")
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

			p := "(n:Resource)"
			c := "WHERE n.id = $id"
			params := map[string]*knowledgepb.InputParam{
				"id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "gid:xyz"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{},
			}, nil)

			node, err := client.GetResourceByID(context.Background(), "gid:xyz")
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:Resource)"
			c := "WHERE n.id = $id"
			params := map[string]*knowledgepb.InputParam{
				"id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "gid:xyz"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xyz",
								ExternalId: "0000",
								Type:       "Store",
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:cba",
								ExternalId: "0101",
								Type:       "Person",
							},
						},
					},
				},
			}, nil)

			node, err := client.GetResourceByID(context.Background(), "gid:xyz")
			Expect(err).To(MatchError(ContainSubstring("unable to complete request")))
			Expect(node).To(BeNil())
		})
	})

	Context("GetResourceByIdentifier", func() {
		It("returns 1 node - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:Resource)"
			c := "WHERE n.external_id = $external_id AND n.type = $type"
			params := map[string]*knowledgepb.InputParam{
				"external_id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "0000"},
				},
				"type": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "store"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xyz",
								ExternalId: "0000",
								Type:       "Store",
							},
						},
					},
				},
			}, nil)

			node, err := client.GetResourceByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "0000",
				Type:       "Store",
			})
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

			p := "(n:Resource)"
			c := "WHERE n.external_id = $external_id AND n.type = $type"
			params := map[string]*knowledgepb.InputParam{
				"external_id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "0000"},
				},
				"type": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "store"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{},
			}, nil)

			node, err := client.GetResourceByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "0000",
				Type:       "Store",
			})
			Expect(err).To(MatchError(ContainSubstring("node not found")))
			Expect(node).To(BeNil())
		})

		It("returns many nodes - fail", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:Resource)"
			c := "WHERE n.external_id = $external_id AND n.type = $type"
			params := map[string]*knowledgepb.InputParam{
				"external_id": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "0000"},
				},
				"type": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "store"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xyz",
								ExternalId: "0000",
								Type:       "Store",
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:cba",
								ExternalId: "0101",
								Type:       "Person",
							},
						},
					},
				},
			}, nil)

			node, err := client.GetResourceByIdentifier(context.Background(), &knowledge.Identifier{
				ExternalID: "0000",
				Type:       "Store",
			})
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

			p := "(:Resource)"
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation: knowledgepb.Operation_OPERATION_READ,
				Path:      p,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xyz",
								ExternalId: "0000",
								Type:       "Store",
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xxx",
								ExternalId: "0001",
								Type:       "Store",
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:yyyy",
								ExternalId: "0002",
								Type:       "Product",
							},
						},
					},
				},
			}, nil)

			nodes, err := client.ListResources(context.Background())
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(3))
		})

		It("ListDigitalTwins - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(:DigitalTwin)"
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation: knowledgepb.Operation_OPERATION_READ,
				Path:      p,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:aaa",
								ExternalId: "0101",
								Type:       "Vehicle",
							},
						},
					},
				},
			}, nil)

			nodes, err := client.ListDigitalTwins(context.Background())
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(2))
		})

		It("ListNodes - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(:Person)"
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation: knowledgepb.Operation_OPERATION_READ,
				Path:      p,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
							},
						},
					},
				},
			}, nil)

			nodes, err := client.ListNodes(context.Background(), "Person")
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(1))
		})
	})

	Context("List by property", func() {
		It("ListResourcesByProperty - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:Resource)"
			c := "WHERE n.location = $location"
			params := map[string]*knowledgepb.InputParam{
				"location": {
					Value: &knowledgepb.InputParam_StringValue{StringValue: "Seattle"},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:xyz",
								ExternalId: "0000",
								Type:       "Store",
								Properties: []*knowledgepb.Property{
									{
										Key: "location",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: "Seattle",
											},
										},
									},
								},
							},
						},
					},
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:koko",
								ExternalId: "1090",
								Type:       "Store",
								Properties: []*knowledgepb.Property{
									{
										Key: "location",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: "Seattle",
											},
										},
									},
								},
							},
						},
					},
				},
			}, nil)

			nodes, err := client.ListResourcesByProperty(context.Background(), &knowledgepb.Property{
				Key: "location",
				Value: &objects.Value{
					Value: &objects.Value_StringValue{
						StringValue: "Seattle",
					},
				},
			})
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(2))
		})

		It("ListDigitalTwinsByProperty - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:DigitalTwin)"
			c := "WHERE n.ssn = $ssn"
			params := map[string]*knowledgepb.InputParam{
				"ssn": {
					Value: &knowledgepb.InputParam_IntegerValue{IntegerValue: 12345},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
								Properties: []*knowledgepb.Property{
									{
										Key: "ssn",
										Value: &objects.Value{
											Value: &objects.Value_IntegerValue{
												IntegerValue: 12345,
											},
										},
									},
								},
							},
						},
					},
				},
			}, nil)

			nodes, err := client.ListDigitalTwinsByProperty(context.Background(), &knowledgepb.Property{
				Key: "ssn",
				Value: &objects.Value{
					Value: &objects.Value_IntegerValue{
						IntegerValue: 12345,
					},
				},
			})
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(1))
		})

		It("ListNodesByProperty - success", func() {
			mockCtrl := gomock.NewController(GinkgoT())
			mockClient := knowledgem.NewMockIdentityKnowledgeAPIClient(mockCtrl)

			client, err := knowledge.NewTestClient(mockClient)
			Ω(err).To(Succeed())

			p := "(n:Person)"
			c := "WHERE n.ssn = $ssn"
			params := map[string]*knowledgepb.InputParam{
				"ssn": {
					Value: &knowledgepb.InputParam_IntegerValue{IntegerValue: 12345},
				},
			}
			mockClient.EXPECT().IdentityKnowledge(gomock.Any(), gomock.Eq(&knowledgepb.IdentityKnowledgeRequest{
				Operation:   knowledgepb.Operation_OPERATION_READ,
				Path:        p,
				Conditions:  c,
				InputParams: params,
			}), gomock.Any()).Return(&knowledgepb.IdentityKnowledgeResponse{
				Paths: []*knowledgepb.Path{
					{
						Nodes: []*knowledgepb.Node{
							{
								Id:         "gid:abc",
								ExternalId: "1010",
								Type:       "Person",
								Properties: []*knowledgepb.Property{
									{
										Key: "ssn",
										Value: &objects.Value{
											Value: &objects.Value_IntegerValue{
												IntegerValue: 12345,
											},
										},
									},
								},
							},
						},
					},
				},
			}, nil)

			nodes, err := client.ListNodesByProperty(context.Background(), "Person", &knowledgepb.Property{
				Key: "ssn",
				Value: &objects.Value{
					Value: &objects.Value_IntegerValue{
						IntegerValue: 12345,
					},
				},
			})
			Expect(err).To(Succeed())
			Expect(nodes).To(HaveLen(1))
		})
	})
})
