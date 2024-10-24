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

package test

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta2"
)

var (
	/*
		policy
			{
			"meta": {
				"policyVersion": "1.0-indykite"
			},
			"subject": {
				"type": "Person"
			},
			"actions": [
				"SUBSCRIBES_TO"
			],
			"resource": {
				"type": "Asset"
			},
			"condition": {
				"cypher": "MATCH (subject:Person)-[:BELONGS_TO]->(:Organization)-[:IS_ON]->
				(s:Subscription) MATCH (s)-[:OFFERS]->(:Service) MATCH (s)-[:COVERS]->
				(resource:Asset) "
			}
		}
	*/

	Resource1 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "pFlpMtkWqCPXVue",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource2 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "pFlpMtkWqCPXVue",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource3 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "pFlpMtkWqCPXVue",
			Type:       "Asset",
			Actions:    []string{"DEMANDS"},
		},
	}

	Resource4 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "QovektcrVBbNmFj",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource5 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "HQKzkgPnGJDiaGo",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource6 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "XcbZruEzGNYHLic",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource7 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "zIDegSbXcRlBeFZ",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource8 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "paLtQSpEcTvzeuC",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource9 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "pFlpMtkWqCPXVue",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
		{
			ExternalId: "HQKzkgPnGJDiaGo",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	ResourceType1 = []*authorizationpb.WhatAuthorizedRequest_ResourceType{
		{Type: "Asset", Actions: []string{"SUBSCRIBES_TO"}},
	}

	ResourceType2 = []*authorizationpb.WhatAuthorizedRequest_ResourceType{
		{Actions: []string{"SUBSCRIBES_TO", "OWNS"}},
	}

	ResourceType3 = []*authorizationpb.WhatAuthorizedRequest_ResourceType{
		{Type: "Asset", Actions: []string{"DEMANDS"}},
	}

	ResourceWho1 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{ExternalId: "pFlpMtkWqCPXVue", Type: "Asset", Actions: []string{"SUBSCRIBES_TO", "OWNS"}},
	}

	ResourceWho2 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{Type: "Asset", Actions: []string{"SUBSCRIBES_TO", "OWNS"}},
	}

	ResourceWho3 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{ExternalId: "pFlpMtkWqCPXVue", Type: "Asset", Actions: []string{"DEMANDS", "REPELS"}},
	}

	NodeBad     = "id"
	NodeNotInDB = "gid:AAAAGaiIPzg6L0DKkbIh22crsFg"
	Node1       = "gid:AAAAHJsPjaxKk0WchnF4wH3Hg10"
	Node2       = "gid:AAAAHO-ocNjhAU1dlkF_1QG22Uo"
	Node3       = "gid:AAAAHNdVLTx1-ExZnjv7nVyRiQc"
	Node4       = "gid:AAAAHH50iSkNRkZni9C12Ed-7fk"
	Node5       = "gid:AAAAHIJg29h5dErYihm_ZRLDB_Y"

	EmailBad   = "test@example.com"
	EmailGood  = "biche@yahoo.co.uk"
	EmailGood2 = "darna@yahoo.co.uk"

	ExternalIDGood = "TrSFiLuoSLGiCIo"

	Asset1  = "HQKzkgPnGJDiaGo"
	Asset2  = "pFlpMtkWqCPXVue"
	Asset3  = "zojWwtKbBLmAXCO"
	Asset4  = "paLtQSpEcTvzeuC"
	Asset5  = "dLZVTSllFCdZfXC"
	Car1    = "gid:AAAAHMjQZuF-L0F6lliyCVGfIRc"
	Car1Ext = "9658744"
	Car2    = "gid:AAAAHM1Lc0CS5EJxpM5QuRUAnrc"
	Car2Ext = "963258"

	Subject1   = "dilZWYdFcmXiojC"
	Subject2   = "fVcaUxJqmOkyOTX"
	Subject3   = "lSPmCXIPRXppszf"
	Subject4   = "NACTFFKUCcceDIz"
	SubjectDT4 = "852147963"
	SubjectDT5 = "741258"

	ConsentConfigID = "gid:AAAAHd5zxpEOlkZtkZQBQ4LEpMg"
	ConsentConfig2  = "gid:AAAAHW5LGqLwe02RgZH6F-p7Tqk"
	ConsentConfig3  = "gid:AAAAHf76WounSEFgg_PZXdvz-eI"
	Application     = "gid:AAAABAm7PCSpUkkej0_iLS8pWrM"
	ConsentInvalid  = "gid:AAAAHXG0xeENcEY6n2qHrf5v7bU"
	ConsentEnforce  = "gid:AAAAHedNk86qM0KHn5KkOWZPSn0"
	ConsentAllow    = "gid:AAAAHfIbn_01-0a9igaJVUOjjg8"

	Resolver = "gid:AAAAIYZfGPsbJEiFnV92GnXklOc"
	URL      = "https://example.com/source2"
	URLUpd   = "https://example.com/sourceupd"
	Method1  = "GET"
	Method2  = "POST"
	Method3  = "ACTION"
	Headers  = map[string]*configpb.ExternalDataResolverConfig_Header{
		"Authorization": {Values: []string{"Bearer edolkUTY"}},
		"Content-Type":  {Values: []string{"application/json"}},
	}
	HeadersUpd = map[string]*configpb.ExternalDataResolverConfig_Header{
		"Authorization": {Values: []string{"Bearer pdnYhjui"}},
		"Content-Type":  {Values: []string{"application/json"}},
	}
	RequestType      = configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON
	RequestPayload   = []byte(`{"key": "value"}`)
	ResponseType     = configpb.ExternalDataResolverConfig_CONTENT_TYPE_JSON
	ResponseSelector = "."

	CustomerID    = os.Getenv("CUSTOMER_ID")
	WrongAppSpace = "gid:AAAAAgDRZxyY6Ecrjhj2GMCtgVI"
	Tags          = []string{"Sitea", "Siteb"}

	NodeFilter1 = &configpb.EntityMatchingPipelineConfig_NodeFilter{
		SourceNodeTypes: []string{"Customer"},
		TargetNodeTypes: []string{"Client"},
	}
	NodeFilter2 = &configpb.EntityMatchingPipelineConfig_NodeFilter{
		SourceNodeTypes: []string{"Employee"},
		TargetNodeTypes: []string{"User"},
	}
	NodeFilter3 = &configpb.EntityMatchingPipelineConfig_NodeFilter{
		SourceNodeTypes: []string{"Employee"},
	}
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func CreateRecordNodeIndividual( //nolint:gocritic // nonamedreturns against unnamedResult
	role string) (*ingestpb.Record, string) {
	externalID := GenerateRandomString(10)
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &knowledgeobjects.Node{
						ExternalId: externalID,
						Type:       "Individual",
						IsIdentity: true,
						Properties: []*knowledgeobjects.Property{
							{
								Type: "email",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: GenerateRandomString(6) + "@yahoo.uk",
									},
								},
								Metadata: &knowledgeobjects.Metadata{
									AssuranceLevel:   1,
									VerificationTime: timestamppb.Now(),
									Source:           "Myself",
									CustomMetadata: map[string]*objects.Value{
										"emaildata": {
											Type: &objects.Value_StringValue{StringValue: "Emaildata"},
										},
									},
								},
							},
							{
								Type: "first_name",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: GenerateRandomString(6),
									},
								},
							},
							{
								Type: "last_name",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: GenerateRandomString(6),
									},
								},
							},
							{
								Type: "role",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: role,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return record, externalID
}

func CreateRecordNoProperty(externalID, nodeType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &knowledgeobjects.Node{
						ExternalId: externalID,
						Type:       nodeType,
						IsIdentity: true,
					},
				},
			},
		},
	}
	return record
}

func DeleteRecord(externalID, nodeType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Delete{
			Delete: &ingestpb.DeleteData{
				Data: &ingestpb.DeleteData_Node{
					Node: &ingestpb.NodeMatch{
						ExternalId: externalID,
						Type:       nodeType,
					},
				},
			},
		},
	}
	return record
}

func DeleteRecordWithProperty(externalID, nodeType, property string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Delete{
			Delete: &ingestpb.DeleteData{
				Data: &ingestpb.DeleteData_NodeProperty{
					NodeProperty: &ingestpb.DeleteData_NodePropertyMatch{
						Match: &ingestpb.NodeMatch{
							ExternalId: externalID,
							Type:       nodeType,
						},
						PropertyType: property,
					},
				},
			},
		},
	}
	return record
}

func UpsertRecordNodeAsset() (*ingestpb.Record, string) { //nolint:gocritic // nonamedreturns against unnamedResult
	externalID := GenerateRandomString(10)
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &knowledgeobjects.Node{
						ExternalId: externalID,
						Type:       "Asset",
						IsIdentity: false,
						Properties: []*knowledgeobjects.Property{
							{
								Type: "maker",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: GenerateRandomString(6),
									},
								},
							},
							{
								Type: "vin",
								Value: &objects.Value{
									Type: &objects.Value_IntegerValue{
										IntegerValue: 123456789,
									},
								},
							},
							{
								Type: "colour",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: "Blue",
									},
								},
							},
							{
								Type: "asset",
								ExternalValue: &knowledgeobjects.ExternalValue{
									Resolver: &knowledgeobjects.ExternalValue_Id{
										Id: Resolver,
									},
								},
							},
							{
								Type: "status",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: "Active",
									},
								},
								Metadata: &knowledgeobjects.Metadata{
									AssuranceLevel:   1,
									VerificationTime: timestamppb.Now(),
									Source:           "Myself",
									CustomMetadata: map[string]*objects.Value{
										"statusdata": {
											Type: &objects.Value_StringValue{StringValue: "StatusData"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return record, externalID
}

func CreateRecordResourceNoProperty(externalID, nodeType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &knowledgeobjects.Node{
						ExternalId: externalID,
						Type:       nodeType,
						IsIdentity: false,
					},
				},
			},
		},
	}
	return record
}

func CreateRecordRelationship(
	sourceExternalID string,
	sourceType string,
	targetExternalID string,
	targetType string,
	relationType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Relationship{
					Relationship: &ingestpb.Relationship{
						Source: &ingestpb.NodeMatch{
							ExternalId: sourceExternalID,
							Type:       sourceType,
						},
						Target: &ingestpb.NodeMatch{
							ExternalId: targetExternalID,
							Type:       targetType,
						},
						Type: relationType,
						Properties: []*knowledgeobjects.Property{
							{
								Type: "property1",
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: "value1",
									},
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
			},
		},
	}
	return record
}

func GetRelationship(
	sourceExternalID string,
	sourceType string,
	targetExternalID string,
	targetType string,
	relationType string) *ingestpb.Relationship {
	relationship := &ingestpb.Relationship{
		Source: &ingestpb.NodeMatch{
			ExternalId: sourceExternalID,
			Type:       sourceType,
		},
		Target: &ingestpb.NodeMatch{
			ExternalId: targetExternalID,
			Type:       targetType,
		},
		Type: relationType,
	}
	return relationship
}

func DeleteRecordRelationship(relationship *ingestpb.Relationship) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Delete{
			Delete: &ingestpb.DeleteData{
				Data: &ingestpb.DeleteData_Relationship{
					Relationship: relationship,
				},
			},
		},
	}
	return record
}

func DeleteRecordRelationshipProperty(
	sourceExternalID string,
	sourceType string,
	targetExternalID string,
	targetType string,
	typeRelation string,
	propertyType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Delete{
			Delete: &ingestpb.DeleteData{
				Data: &ingestpb.DeleteData_RelationshipProperty{
					RelationshipProperty: &ingestpb.DeleteData_RelationshipPropertyMatch{
						Source: &ingestpb.NodeMatch{
							ExternalId: sourceExternalID,
							Type:       sourceType,
						},
						Target: &ingestpb.NodeMatch{
							ExternalId: targetExternalID,
							Type:       targetType,
						},
						Type:         typeRelation,
						PropertyType: propertyType,
					},
				},
			},
		},
	}
	return record
}

func CreateBatchNodes( //nolint:gocritic // nonamedreturns against unnamedResult
	typeNode string) (*knowledgeobjects.Node, string) {
	externalID := GenerateRandomString(10)
	node := &knowledgeobjects.Node{
		ExternalId: externalID,
		Type:       typeNode,
		IsIdentity: true,
		Properties: []*knowledgeobjects.Property{
			{
				Type: "email",
				Value: &objects.Value{
					Type: &objects.Value_StringValue{
						StringValue: GenerateRandomString(6) + "@yahoo.uk",
					},
				},
				Metadata: &knowledgeobjects.Metadata{
					AssuranceLevel:   1,
					VerificationTime: timestamppb.Now(),
					Source:           "Myself",
					CustomMetadata: map[string]*objects.Value{
						"emaildata": {
							Type: &objects.Value_StringValue{StringValue: "Emaildata"},
						},
					},
				},
			},
			{
				Type: "first_name",
				Value: &objects.Value{
					Type: &objects.Value_StringValue{
						StringValue: GenerateRandomString(6),
					},
				},
			},
			{
				Type: "last_name",
				Value: &objects.Value{
					Type: &objects.Value_StringValue{
						StringValue: GenerateRandomString(6),
					},
				},
			},
			{
				Type: "asset",
				ExternalValue: &knowledgeobjects.ExternalValue{
					Resolver: &knowledgeobjects.ExternalValue_Id{
						Id: Resolver,
					},
				},
			},
		},
		Tags: []string{"Sitea", "Siteb"},
	}
	return node, externalID
}

func CreateBatchNodesError(
	typeNode string) *knowledgeobjects.Node {
	externalID := GenerateRandomString(10)
	node := &knowledgeobjects.Node{
		ExternalId: externalID,
		Type:       typeNode,
		IsIdentity: true,
		Properties: []*knowledgeobjects.Property{
			{
				Type: "email",
				Value: &objects.Value{
					Type: &objects.Value_StringValue{
						StringValue: GenerateRandomString(6) + "@yahoo.uk",
					},
				},
			},
			{
				Type: "last_name",
				Value: &objects.Value{
					Type: &objects.Value_StringValue{
						StringValue: GenerateRandomString(6),
					},
				},
				ExternalValue: &knowledgeobjects.ExternalValue{
					Resolver: &knowledgeobjects.ExternalValue_Id{
						Id: Resolver,
					},
				},
			},
		},
	}
	return node
}

func CreateBatchNodesNoResolver(
	typeNode string) *knowledgeobjects.Node {
	externalID := GenerateRandomString(10)
	node := &knowledgeobjects.Node{
		ExternalId: externalID,
		Type:       typeNode,
		IsIdentity: true,
		Properties: []*knowledgeobjects.Property{
			{
				Type: "last_name",
				ExternalValue: &knowledgeobjects.ExternalValue{
					Resolver: &knowledgeobjects.ExternalValue_Id{
						Id: NodeNotInDB,
					},
				},
			},
		},
	}
	return node
}

func BatchNodesType( //nolint:gocritic // nonamedreturns against unnamedResult
	typeNode string) ([]*knowledgeobjects.Node, string, string) {
	node1, externalID1 := CreateBatchNodes(typeNode)
	node2, externalID2 := CreateBatchNodes(typeNode)
	nodes := []*knowledgeobjects.Node{
		node1, node2,
	}
	return nodes, externalID1, externalID2
}

func BatchRelationships(
	relationship *ingestpb.Relationship, relationship2 *ingestpb.Relationship) []*ingestpb.Relationship {
	relationships := []*ingestpb.Relationship{
		relationship,
		relationship2,
	}
	return relationships
}

func CreateBatchNodeMatch(
	externalID string, typeNode string) *ingestpb.NodeMatch {
	nodeMatch := &ingestpb.NodeMatch{
		ExternalId: externalID,
		Type:       typeNode,
	}
	return nodeMatch
}

func BatchNodesMatch(
	nodeMatch *ingestpb.NodeMatch, nodeMatch2 *ingestpb.NodeMatch) []*ingestpb.NodeMatch {
	nodes := []*ingestpb.NodeMatch{
		nodeMatch,
		nodeMatch2,
	}
	return nodes
}
