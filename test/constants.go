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

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta2"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
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

func CreateRecordIndividual( //nolint:gocritic // nonamedreturns against unnamedResult
	role string) (*ingestpb.Record, string) {
	externalID := GenerateRandomString(10)
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &ingestpb.Node{
						Type: &ingestpb.Node_DigitalTwin{
							DigitalTwin: &ingestpb.DigitalTwin{
								ExternalId: externalID,
								Type:       "Individual",
								Properties: []*ingestpb.Property{
									{
										Key: "email",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: GenerateRandomString(6) + "@yahoo.uk",
											},
										},
									},
									{
										Key: "first_name",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: GenerateRandomString(6),
											},
										},
									},
									{
										Key: "last_name",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: GenerateRandomString(6),
											},
										},
									},
									{
										Key: "role",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: role,
											},
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

func CreateRecordNoProp(externalID, nodeType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &ingestpb.Node{
						Type: &ingestpb.Node_DigitalTwin{
							DigitalTwin: &ingestpb.DigitalTwin{
								ExternalId: externalID,
								Type:       nodeType,
							},
						},
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

func DeleteRecordProperty(externalID, nodeType, property string) *ingestpb.Record {
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
						Key: property,
					},
				},
			},
		},
	}
	return record
}

func UpsertRecordAsset() (*ingestpb.Record, string) { //nolint:gocritic // nonamedreturns against unnamedResult
	externalID := GenerateRandomString(10)
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &ingestpb.Node{
						Type: &ingestpb.Node_Resource{
							Resource: &ingestpb.Resource{
								ExternalId: externalID,
								Type:       "Asset",
								Properties: []*ingestpb.Property{
									{
										Key: "maker",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: GenerateRandomString(6),
											},
										},
									},
									{
										Key: "vin",
										Value: &objects.Value{
											Value: &objects.Value_IntegerValue{
												IntegerValue: 123456789,
											},
										},
									},
									{
										Key: "colour",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: "Blue",
											},
										},
									},
									{
										Key: "asset",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: "T",
											},
										},
									},
									{
										Key: "status",
										Value: &objects.Value{
											Value: &objects.Value_StringValue{
												StringValue: "Active",
											},
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

func CreateRecordResourceNoProp(externalID, nodeType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Node{
					Node: &ingestpb.Node{
						Type: &ingestpb.Node_Resource{
							Resource: &ingestpb.Resource{
								ExternalId: externalID,
								Type:       nodeType,
							},
						},
					},
				},
			},
		},
	}
	return record
}

func CreateRecordRelation(
	sourceExternalID string,
	sourceType string,
	targetExternalID string,
	targetType string,
	relationType string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Upsert{
			Upsert: &ingestpb.UpsertData{
				Data: &ingestpb.UpsertData_Relation{
					Relation: &ingestpb.Relation{
						Match: &ingestpb.RelationMatch{
							SourceMatch: &ingestpb.NodeMatch{
								ExternalId: sourceExternalID,
								Type:       sourceType,
							},
							TargetMatch: &ingestpb.NodeMatch{
								ExternalId: targetExternalID,
								Type:       targetType,
							},
							Type: relationType,
						},
						Properties: []*ingestpb.Property{
							{
								Key: "property1",
								Value: &objects.Value{
									Value: &objects.Value_StringValue{
										StringValue: "value1",
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

func GetRelationMatch(
	sourceExternalID string,
	sourceType string,
	targetExternalID string,
	targetType string,
	relationType string) *ingestpb.RelationMatch {
	match := &ingestpb.RelationMatch{
		SourceMatch: &ingestpb.NodeMatch{
			ExternalId: sourceExternalID,
			Type:       sourceType,
		},
		TargetMatch: &ingestpb.NodeMatch{
			ExternalId: targetExternalID,
			Type:       targetType,
		},
		Type: relationType,
	}
	return match
}

func DeleteRecordRelation(match *ingestpb.RelationMatch) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Delete{
			Delete: &ingestpb.DeleteData{
				Data: &ingestpb.DeleteData_Relation{
					Relation: match,
				},
			},
		},
	}
	return record
}

func DeleteRecordRelationProperty(match *ingestpb.RelationMatch, property string) *ingestpb.Record {
	record := &ingestpb.Record{
		Id: uuid.New().String(),
		Operation: &ingestpb.Record_Delete{
			Delete: &ingestpb.DeleteData{
				Data: &ingestpb.DeleteData_RelationProperty{
					RelationProperty: &ingestpb.DeleteData_RelationPropertyMatch{
						Match: match,
						Key:   property,
					},
				},
			},
		},
	}
	return record
}
