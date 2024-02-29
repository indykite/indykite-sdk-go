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

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
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
				"type": "Individual"
			},
			"actions": [
				""
			],
			"resource": {
				"type": "Asset"
			},
			"condition": {
				"cypher": "MATCH (subject:Individual)-[:BELONGS_TO]->(:Organization)-[:IS_ON]->
				(s:Subscription) MATCH (s)-[:OFFERS]->(:Service) MATCH (s)-[:COVERS]->
				(resource:Asset) WITH subject, resource"
			}
		}
	*/

	Resource1 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "LPcearawBJWDQLR",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource2 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "LPcearawBJWDQLR",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource3 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "LPcearawBJWDQLR",
			Type:       "Asset",
			Actions:    []string{"DEMANDS"},
		},
	}

	Resource4 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "mfYbpowiNPJQCBY",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource5 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "CCbJwkQtLOmCdLq",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource6 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "aXQMRIcTzyIyeKC",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource7 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "EvfDHrEObtYVleh",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource8 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "BLOXgHAvWFMHDsS",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
	}

	Resource9 = []*authorizationpb.IsAuthorizedRequest_Resource{
		{
			ExternalId: "LPcearawBJWDQLR",
			Type:       "Asset",
			Actions:    []string{"SUBSCRIBES_TO"},
		},
		{
			ExternalId: "CCbJwkQtLOmCdLq",
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
		{ExternalId: "LPcearawBJWDQLR", Type: "Asset", Actions: []string{"SUBSCRIBES_TO", "OWNS"}},
	}

	ResourceWho2 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{Type: "Asset", Actions: []string{"SUBSCRIBES_TO", "OWNS"}},
	}

	ResourceWho3 = []*authorizationpb.WhoAuthorizedRequest_Resource{
		{ExternalId: "LPcearawBJWDQLR", Type: "Asset", Actions: []string{"DEMANDS", "REPELS"}},
	}

	DigitalTwinBad     = "id"
	DigitalTwinNotInDB = "gid:AAAAGaiIPzg6L0DKkbIh22crsFg"
	DigitalTwin1       = "gid:AAAAFR3royp640c-gXRGdusXM4Y"
	DigitalTwin2       = "gid:AAAAFf6Y9ZMWhEdsr3INueqfRLU"

	EmailBad  = "test@example.com"
	EmailGood = "paulo@yahoo.uk"

	ExternalIDGood = "TrSFiLuoSLGiCIo"

	Asset1 = "CCbJwkQtLOmCdLq"
	Asset2 = "LPcearawBJWDQLR"
	Asset3 = "zBiBMaYOaDmdCyX"
	Asset4 = "BLOXgHAvWFMHDsS"

	Subject1 = "HLEgiljrtoNEiyX"
	Subject2 = "zvPYDXxXyVgeZHw"
	Subject3 = "TrSFiLuoSLGiCIo"

	// TokenGoodFormat is a valid format for jwt.
	TokenGoodFormat = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9." +
		"dyt0CoTl4WoVjAHI9Q_CwSKhl6d_9rhM3NrXuJttkao" // #nosec G101
	TokenBad = "token_invalid_format"
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
								Value: &objects.Value{
									Type: &objects.Value_StringValue{
										StringValue: "T",
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
