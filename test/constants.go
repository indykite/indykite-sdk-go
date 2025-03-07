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
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"google.golang.org/protobuf/types/known/timestamppb"

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
	knowledgeobjects "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/objects/v1beta1"
	knowledgepb "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta2"
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
		policy trust score
		{
		"meta": {
			"policyVersion": "1.0-indykite"
		},
		"subject": {
			"type": "Agent"
		},
		"actions": [
			"CAN_USE"
		],
		"resource": {
			"type": "Sensor"
		},
		"condition": {
			"cypher": "MATCH (:_TrustScore)<-[:_HAS]-(subject)-[:USE]->(resource)"
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

	Resource10 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "Truck1",
			Type:       "Truck",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource11 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "Truck3",
			Type:       "Truck",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource12 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "Truck4",
			Type:       "Truck",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource13 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "Truck5",
			Type:       "Truck",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource14 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "Truck2",
			Type:       "Truck",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resources15 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "Sensor1",
			Type:       "Sensor",
			Actions:    []string{"CAN_USE"},
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

	ResourceType4 = []*authorizationpb.WhatAuthorizedRequest_ResourceType{
		{Type: "Truck", Actions: []string{"SUBSCRIBES_TO"}},
	}

	ResourceType5 = []*authorizationpb.WhatAuthorizedRequest_ResourceType{
		{Type: "Sensor", Actions: []string{"CAN_USE"}},
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

	ResourceWho4 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{ExternalId: "Truck1", Type: "Truck", Actions: []string{"SUBSCRIBES_TO", "OWNS"}},
	}

	ResourceWho5 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{ExternalId: "Sensor1", Type: "Sensor", Actions: []string{"CAN_USE"}},
	}

	NodeBad     = "id"
	NodeNotInDB = "gid:AAAAGaiIPzg6L0DKkbIh22crsFg"
	Node1       = "gid:AAAAHJHidWK1yER6qC2EnUdqZFE"
	Node2       = "gid:AAAAHPeCE5y5IEpelyyeCXrwyLk"
	Node3       = "gid:AAAAHJwDOghFkk5MgOpl1bUEnIM"
	Node4       = "gid:AAAAHGbvYyieDUb4sv5tSVFLTM8"
	Node5       = "gid:AAAAHBWT8xE1JE3Dt1JUp55dqLk"
	Node6       = "gid:AAAAHA_OG-Ky1EN_vY1wGXx_qg8"
	Node7       = "gid:AAAAHFHWzJ3CnkgYmjjdTLUYhUw"

	EmailBad   = "test@example.com"
	EmailGood  = "biche@yahoo.co.uk"
	EmailGood2 = "darna@yahoo.co.uk" // gid:AAAAHFHWzJ3CnkgYmjjdTLUYhUw
	EmailWhat  = "barnabebe@yahoo.com"

	ExternalIDGood = "TrSFiLuoSLGiCIo"

	Asset1   = "HQKzkgPnGJDiaGo"
	Asset2   = "pFlpMtkWqCPXVue"
	Asset3   = "zojWwtKbBLmAXCO"
	Asset4   = "paLtQSpEcTvzeuC"
	Asset5   = "dLZVTSllFCdZfXC"
	Car1     = "gid:AAAAHEYyKGXN3kopq7t3ct6CJIY"
	Car1Ext  = "9658744"
	Car2     = "gid:AAAAHAfh6gGfwEQ8prtS5mZDHn4"
	Car2Ext  = "963258"
	Truck1   = "Truck1"
	Truck1Id = "gid:AAAAHPUWEuZveU5ggLkQe_ige2w"
	Agent1   = "nemo"
	Agent2   = "barracuda"
	Agent3   = "ozzy"
	Agent1Id = "gid:AAAAHAhop0O2Mk_qux4RXxZFYCY"
	Sensor1  = "Sensor1"

	Subject1   = "dilZWYdFcmXiojC" // gid:AAAAHJwDOghFkk5MgOpl1bUEnIM
	Subject2   = "fVcaUxJqmOkyOTX" // gid:AAAAHPeCE5y5IEpelyyeCXrwyLk
	Subject3   = "lSPmCXIPRXppszf" // gid:AAAAHJHidWK1yER6qC2EnUdqZFE
	Subject4   = "NACTFFKUCcceDIz" // gid:AAAAHFHWzJ3CnkgYmjjdTLUYhUw
	Subject5   = "barnabebe"       // gid:AAAAHA_OG-Ky1EN_vY1wGXx_qg8
	SubjectDT4 = "852147963"       // gid:AAAAHGbvYyieDUb4sv5tSVFLTM8
	SubjectDT5 = "741258"          // gid:AAAAHBWT8xE1JE3Dt1JUp55dqLk

	ConsentConfigID = "gid:AAAAHd5zxpEOlkZtkZQBQ4LEpMg"
	ConsentConfig2  = "gid:AAAAHZVkKsIHW0Twv4dIcaLdEkc"
	ConsentConfig3  = "gid:AAAAHcw9InzWn0FioK4TwqyrZS8"
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
	TestAppSpace  = "gid:AAAAAvLdMmzCWEE-hyrJDXVGgOk"
	Tags          = []string{"Sitea", "Siteb"}

	NodeFilter1 = &configpb.EntityMatchingPipelineConfig_NodeFilter{
		SourceNodeTypes: []string{"Customer"},
		TargetNodeTypes: []string{"Client"},
	}
	NodeFilter3 = &configpb.EntityMatchingPipelineConfig_NodeFilter{
		SourceNodeTypes: []string{"Employee"},
	}
	Dimensions1 = []*configpb.TrustScoreDimension{
		{
			Name:   configpb.TrustScoreDimension_NAME_VERIFICATION,
			Weight: 0.5,
		},
		{
			Name:   configpb.TrustScoreDimension_NAME_ORIGIN,
			Weight: 0.5,
		},
	}
	Schedule1 = configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_THREE_HOURS
	Schedule2 = configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_SIX_HOURS
	Schedule3 = configpb.TrustScoreProfileConfig_UPDATE_FREQUENCY_INVALID
	Query1    = `MATCH (n:Person)-[:BELONGS_TO]->(o:Organization)-[:OWNS]->(t:Truck)
	WHERE n.external_id=$external_id AND n.type=$type`
	Query2 = `MATCH (n:Person)-[:BELONGS_TO]->(o:Organization)-[:OWNS]->(t:Truck)-[:HAS]->(p:Property:External)
	 WHERE n.external_id=$external_id AND n.type=$type`
	Query3  = "MATCH (n:Person)-[:BELONGS_TO]->(o:Organization)-[:OWNS]->(t:Truck)"
	Query4  = "MATCH (n:Person)-[:BELONGS_TO]->(o:Organization)-[:OWNS]->(t:Truck)-[:HAS]->(p:Property:External)"
	Query5  = "MATCH (n:Resource)-[:HAS]->(p:Property) WHERE p.type=$type and p.value=$value"
	Query6  = "MATCH (n:Resource)-[:HAS]->(p:Property:External) WHERE p.type=$type and p.value=$value"
	Query7  = "Person "
	Query8  = "MATCH (a:Agent)-[r:_HAS]->(t:_TrustScore)"
	Query9  = "MATCH(t:_TrustScore {_ingress: $profile})"
	Query10 = "MATCH (t:_TrustScore)"
	Params1 = map[string]*objects.Value{
		"external_id": {
			Type: &objects.Value_StringValue{StringValue: Subject2},
		},
		"type": {
			Type: &objects.Value_StringValue{StringValue: "Person"},
		},
	}
	Params2 = map[string]*objects.Value{
		"type": {
			Type: &objects.Value_StringValue{StringValue: "first_name"},
		},
		"value": {
			Type: &objects.Value_StringValue{StringValue: "darna"},
		},
	}
	Params3 = map[string]*objects.Value{}
	Params4 = map[string]*objects.Value{
		"type": {
			Type: &objects.Value_StringValue{StringValue: "color"},
		},
		"value": {
			Type: &objects.Value_StringValue{StringValue: "blue"},
		},
	}
	Params5 = map[string]*objects.Value{
		"type": {
			Type: &objects.Value_StringValue{StringValue: "echo"},
		},
		"value": {
			Type: &objects.Value_StringValue{StringValue: "2024"},
		},
	}
	Params6 = map[string]*objects.Value{
		"profile": {
			Type: &objects.Value_StringValue{StringValue: "like-real-config-node-name-ts3"},
		},
	}
	Returns1 = []*knowledgepb.Return{
		{
			Variable: "n",
		},
	}
	Returns2 = []*knowledgepb.Return{
		{
			Variable:   "a",
			Properties: []string{},
		},
		{
			Variable:   "t",
			Properties: []string{},
		},
	}
	Returns3 = []*knowledgepb.Return{
		{
			Variable: "t",
		},
	}
	Returns4 = []*knowledgepb.Return{
		{
			Variable:   "t",
			Properties: []string{"_origin", "_verification"},
		},
	}
	Matcher1 = gstruct.Fields{
		"Nodes": gomega.ContainElement(gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Id":         gomega.Equal(Node2),
			"ExternalId": gomega.Equal(Subject2),
			"Type":       gomega.Equal("Person"),
		}))),
		"Relationships": gomega.BeEmpty(),
	}
	Matcher2 = gstruct.Fields{
		"Nodes": gomega.ContainElement(gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Id":         gomega.Equal(Truck1Id),
			"ExternalId": gomega.Equal(Truck1),
			"Type":       gomega.Equal("Truck"),
		}))),
		"Relationships": gomega.BeEmpty(),
	}
	Matcher3 = gstruct.Fields{
		"Id":         gomega.Equal(Truck1Id),
		"ExternalId": gomega.Equal(Truck1),
		"Type":       gomega.Equal("Truck"),
	}
	Matcher4 = gstruct.Fields{
		"Id":         gomega.Equal(Node2),
		"ExternalId": gomega.Equal(Subject2),
		"Type":       gomega.Equal("Person"),
	}
	Matcher5 = gstruct.Fields{
		"Nodes": gomega.ContainElement(gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Id":         gomega.Equal(Agent1Id),
			"ExternalId": gomega.Equal(Agent1),
			"Type":       gomega.Equal("Agent"),
		}))),
		"Relationships": gomega.BeEmpty(),
	}
	Matcher6 = gstruct.Fields{
		"Nodes": gomega.ContainElement(gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Type": gomega.Equal("_TrustScore"),
		}))),
		"Relationships": gomega.BeEmpty(),
	}
	Matcher7 = gstruct.Fields{
		"Nodes": gomega.ContainElement(gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Type": gomega.Equal("_TrustScore"),
		}))),
		"Relationships": gomega.BeEmpty(),
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

//revive:disable
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
