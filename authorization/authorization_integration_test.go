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

//go:build integration

package authorization_test

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	objectpb "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
	googlehelper "github.com/indykite/indykite-sdk-go/helpers"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Authorized", func() {
	Describe("IsAuthorized", func() {
		var auditLogIdentifier string
		var noAuditLogEntry bool
		BeforeEach(func() {
			noAuditLogEntry = false
			auditLogIdentifier = fmt.Sprintf("%v", time.Now().UnixNano())
		})
		AfterEach(func(ctx SpecContext) {
			// Check that the audit log is present in BigQuery
			client, err := googlehelper.BqClient(ctx)
			Expect(err).To(Succeed())
			filter, err := googlehelper.FillFilterFieldsFromEnvironment()
			Expect(err).To(Succeed())
			filter.ChangeType = "type.googleapis.com/indykite.auditsink.v1beta1.IsAuthorized"
			filter.EventSource = "AuthorizationService"
			filter.EventType = "indykite.audit.authorization.isauthorized"
			filter.AuditLogIdentifier = auditLogIdentifier
			res, err := googlehelper.QueryAuditLog(ctx, client, &filter)
			Expect(err).To(Succeed())
			if noAuditLogEntry {
				Expect(res.Data).To(BeEmpty())
			} else {
				Expect(res.Data).To(ContainSubstring(auditLogIdentifier))
			}
		}, NodeTimeout(time.Second*10))

		DescribeTable("DT Authorization tests",
			func(digitalTwinId string,
				resources []*authorizationpb.IsAuthorizedRequest_Resource,
				policyTags []string,
				expectedAllow bool,
				expectedError string,
				numberResources int) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				digitalTwin := &authorizationpb.DigitalTwin{
					Id: digitalTwinId,
				}

				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.IsAuthorized(
					context.Background(),
					digitalTwin,
					resources,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					if numberResources == 0 {
						noAuditLogEntry = true
					}
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resources[0].Type
					resource := resources[0].ExternalId
					action := resources[0].Actions[0]

					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllKeys(Keys{
									resource: PointTo(MatchFields(IgnoreExtras, Fields{
										"Actions": MatchAllKeys(Keys{
											action: PointTo(MatchFields(IgnoreExtras, Fields{
												"Allow": Equal(expectedAllow),
											})),
										}),
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("Authorized single resource", integration.Node1, integration.Resource1,
				[]string{}, true, "", 1),
			Entry("Authorized with tags", integration.Node1, integration.Resource1,
				[]string{"TagOne"}, true, "", 1),
			Entry("Unauthorized with wrong tags", integration.Node1, integration.Resource1,
				[]string{"TagBad"}, false, "", 1),
			Entry("Invalid resource type", integration.Node1, integration.Resource2, []string{}, false,
				"invalid IsAuthorizedRequest_Resource.Type: value length must be between 2 and 50 runes", 0),
			Entry("Invalid digital twin subject", integration.NodeBad, integration.Resource2,
				[]string{}, false,
				"invalid IsAuthorizedRequest.Subject", 0),
			Entry("Digital twin not in DB", integration.NodeNotInDB, integration.Resource1,
				[]string{}, false, "", 1),
			Entry("Resource not in DB", integration.Node1, integration.Resource3,
				[]string{}, false, "", 1),
			Entry("Resource without subscription", integration.Node1, integration.Resource4,
				[]string{}, false, "", 1),
			Entry("Resource without organization", integration.Node1, integration.Resource7,
				[]string{}, false, "", 1),
			Entry("Authorized without service", integration.Node1, integration.Resource6,
				[]string{}, false, "", 1),
			Entry("Resource not linked", integration.Node1, integration.Resource4,
				[]string{}, false, "", 1),
			Entry("Authorized with external property", integration.Node3,
				integration.Resource10, []string{}, true, "", 1),
			Entry("Authorized without external property", integration.Node3,
				integration.Resource12, []string{}, false, "", 1),
			Entry("Authorized with external property against policy", integration.Node3,
				integration.Resource11, []string{}, false, "", 1),
		)

		DescribeTable("DT Authorization multiple tests",
			func(digitalTwinId string,
				resources []*authorizationpb.IsAuthorizedRequest_Resource,
				policyTags []string,
				expectedAllow bool,
				expectedError string,
				numberResources int) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				digitalTwin := &authorizationpb.DigitalTwin{
					Id: digitalTwinId,
				}

				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.IsAuthorized(
					context.Background(),
					digitalTwin,
					resources,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					if numberResources == 0 {
						noAuditLogEntry = true
					}
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resources[0].Type
					resource := resources[0].ExternalId
					action := resources[0].Actions[0]
					resource1 := resources[1].ExternalId
					action1 := resources[1].Actions[0]
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllKeys(Keys{
									resource: PointTo(MatchFields(IgnoreExtras, Fields{
										"Actions": MatchAllKeys(Keys{
											action: PointTo(MatchFields(IgnoreExtras, Fields{
												"Allow": Equal(expectedAllow),
											})),
										}),
									})),
									resource1: PointTo(MatchFields(IgnoreExtras, Fields{
										"Actions": MatchAllKeys(Keys{
											action1: PointTo(MatchFields(IgnoreExtras, Fields{
												"Allow": Equal(expectedAllow),
											})),
										}),
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("Authorized multiple resources", integration.Node2, integration.Resource9,
				[]string{}, true, "", 2),
		)

		DescribeTable("Property Authorization tests",
			func(typeNode string,
				property string,
				resources []*authorizationpb.IsAuthorizedRequest_Resource,
				policyTags []string,
				expectedAllow bool,
				expectedError string) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				digitalTwinProperty := &authorizationpb.Property{
					Type:  typeNode,
					Value: objectpb.String(property),
				}
				// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.IsAuthorizedByProperty(
					context.Background(),
					digitalTwinProperty,
					resources,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					noAuditLogEntry = true
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resources[0].Type
					resource := resources[0].ExternalId
					action := resources[0].Actions[0]
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllKeys(Keys{
									resource: PointTo(MatchFields(IgnoreExtras, Fields{
										"Actions": MatchAllKeys(Keys{
											action: PointTo(MatchFields(IgnoreExtras, Fields{
												"Allow": Equal(expectedAllow),
											})),
										}),
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("Authorized property", "email", integration.EmailGood,
				integration.Resource1, []string{}, true, ""),
			Entry("Authorized property not in DB", "email", integration.EmailBad,
				integration.Resource1, []string{}, false, ""),
			Entry("Authorized property with external property", "email", integration.EmailGood,
				integration.Resource14, []string{}, true, ""),
		)

		DescribeTable("ExternalID Authorization tests",
			func(typeNode string,
				id string,
				resources []*authorizationpb.IsAuthorizedRequest_Resource,
				policyTags []string,
				expectedAllow bool,
				expectedError string) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				externalID := &authorizationpb.ExternalID{
					Type:       typeNode,
					ExternalId: id,
				}
				// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.IsAuthorizedByExternalID(
					context.Background(),
					externalID,
					resources,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					noAuditLogEntry = true
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resources[0].Type
					resource := resources[0].ExternalId
					action := resources[0].Actions[0]
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllKeys(Keys{
									resource: PointTo(MatchFields(IgnoreExtras, Fields{
										"Actions": MatchAllKeys(Keys{
											action: PointTo(MatchFields(IgnoreExtras, Fields{
												"Allow": Equal(expectedAllow),
											})),
										}),
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("Authorized external ID", "Person", integration.Subject2,
				integration.Resource1, []string{}, true, ""),
			Entry("Authorized external ID not in DB", "Person", "anythingwrong",
				integration.Resource1, []string{}, false, ""),
			Entry("External ID with external property", "Person", integration.Subject2,
				integration.Resource14, []string{}, true, ""),
			Entry("External ID with external property against policy", "Person", integration.Subject2,
				integration.Resource11, []string{}, false, ""),
		)
	})

	Describe("WhatAuthorized", func() {
		var auditLogIdentifier string
		var noAuditLogEntry bool
		BeforeEach(func() {
			auditLogIdentifier = fmt.Sprintf("%v", time.Now().UnixNano())
			noAuditLogEntry = false
		})
		AfterEach(func(ctx SpecContext) {
			// Check that the audit log is present in BigQuery
			client, err := googlehelper.BqClient(ctx)
			Expect(err).To(Succeed())
			filter, err := googlehelper.FillFilterFieldsFromEnvironment()
			Expect(err).To(Succeed())
			filter.ChangeType = "type.googleapis.com/indykite.auditsink.v1beta1.WhatAuthorized"
			filter.EventSource = "AuthorizationService"
			filter.EventType = "indykite.audit.authorization.whatauthorized"
			filter.AuditLogIdentifier = auditLogIdentifier
			res, err := googlehelper.QueryAuditLog(ctx, client, &filter)
			Expect(err).To(Succeed())
			if noAuditLogEntry {
				Expect(res.Data).To(BeEmpty())
			} else {
				Expect(res.Data).To(ContainSubstring(auditLogIdentifier))
			}
		}, NodeTimeout(time.Second*10))

		DescribeTable("What Authorization DT",
			func(digitalTwinId string,
				resourcesTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
				results []string,
				policyTags []string,
				expectedError string) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				digitalTwin := &authorizationpb.DigitalTwin{
					Id: digitalTwinId,
				}

				// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.WhatAuthorized(
					context.Background(),
					digitalTwin,
					resourcesTypes,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					noAuditLogEntry = true
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resourcesTypes[0].Type
					action := resourcesTypes[0].Actions[0]
					resourceMatcher := BeEmpty() // Default to empty if no results
					if len(results) > 0 {
						elements := Elements{}
						for i, result := range results {
							elements[fmt.Sprintf("%d", i)] = PointTo(MatchFields(IgnoreExtras, Fields{
								"ExternalId": Equal(result),
							}))
						}
						resourceMatcher = MatchAllElementsWithIndex(IndexIdentity, elements)
					}
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Actions": MatchAllKeys(Keys{
									action: PointTo(MatchFields(IgnoreExtras, Fields{
										"Resources": resourceMatcher,
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("What authorized DT", integration.Node1, integration.ResourceType1,
				[]string{integration.Asset4, integration.Asset1, integration.Asset3,
					integration.Asset2, integration.Asset5}, []string{}, ""),
			Entry("What Authorized DT Resource Non Valid", integration.Node1, integration.ResourceType2,
				[]string{}, []string{},
				"invalid WhatAuthorizedRequest_ResourceType.Type: value length must be between 2 and 50 runes"),
			Entry("What Authorized DT Subject Non Valid", integration.NodeBad, integration.ResourceType1,
				[]string{}, []string{}, "invalid DigitalTwin.Id: value length must be between 27 and 100 runes"),
			Entry("What Authorized DT Subject Not In DB", integration.NodeNotInDB,
				integration.ResourceType1, []string{}, []string{}, ""),
			Entry("What Authorized DT Resource Not In DB", integration.Node1,
				integration.ResourceType3, []string{}, []string{}, ""),
			Entry("What Authorized DT With External Property", integration.Node6, integration.ResourceType4,
				[]string{integration.Truck1}, []string{}, ""),
			Entry("What Authorized DT With External Property Wrong Action", integration.Node7,
				integration.ResourceType4, []string{}, []string{}, ""),
		)

		DescribeTable("What Authorization Property",
			func(propertyType string,
				property string,
				resourcesTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
				results []string,
				policyTags []string,
				expectedError string) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				digitalTwinProperty := &authorizationpb.Property{
					Type:  propertyType,
					Value: objectpb.String(property),
				}

				// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.WhatAuthorizedByProperty(
					context.Background(),
					digitalTwinProperty,
					resourcesTypes,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					noAuditLogEntry = true
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resourcesTypes[0].Type
					action := resourcesTypes[0].Actions[0]
					resourceMatcher := BeEmpty() // Default to empty if no results
					if len(results) > 0 {
						elements := Elements{}
						for i, result := range results {
							elements[fmt.Sprintf("%d", i)] = PointTo(MatchFields(IgnoreExtras, Fields{
								"ExternalId": Equal(result),
							}))
						}
						resourceMatcher = MatchAllElementsWithIndex(IndexIdentity, elements)
					}
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Actions": MatchAllKeys(Keys{
									action: PointTo(MatchFields(IgnoreExtras, Fields{
										"Resources": resourceMatcher,
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("What Authorized Property", "email", integration.EmailGood2, integration.ResourceType1,
				[]string{integration.Asset4, integration.Asset1, integration.Asset3}, []string{}, ""),
			Entry("What Authorized Property Not In DB", "email", integration.EmailBad, integration.ResourceType1,
				[]string{}, []string{}, ""),
			Entry("What Authorized Property With External Property", "email", integration.EmailWhat,
				integration.ResourceType4, []string{integration.Truck1}, []string{}, ""),
		)

		DescribeTable("What Authorization ExternalID",
			func(nodeType string,
				id string,
				resourcesTypes []*authorizationpb.WhatAuthorizedRequest_ResourceType,
				results []string,
				policyTags []string,
				expectedError string) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				externalID := &authorizationpb.ExternalID{
					Type:       nodeType,
					ExternalId: id,
				}

				// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				resp, err := authorizationClient.WhatAuthorizedByExternalID(
					context.Background(),
					externalID,
					resourcesTypes,
					inputParams,
					policyTags,
					retry.WithMax(5),
				)

				if expectedError != "" {
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resourcesTypes[0].Type
					action := resourcesTypes[0].Actions[0]
					resourceMatcher := BeEmpty() // Default to empty if no results
					if len(results) > 0 {
						elements := Elements{}
						for i, result := range results {
							elements[fmt.Sprintf("%d", i)] = PointTo(MatchFields(IgnoreExtras, Fields{
								"ExternalId": Equal(result),
							}))
						}
						resourceMatcher = MatchAllElementsWithIndex(IndexIdentity, elements)
					}
					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Actions": MatchAllKeys(Keys{
									action: PointTo(MatchFields(IgnoreExtras, Fields{
										"Resources": resourceMatcher,
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("What Authorized External ID", "Person", integration.Subject4, integration.ResourceType1,
				[]string{integration.Asset4, integration.Asset1, integration.Asset3}, []string{}, ""),
			Entry("What Authorized External ID Not In DB", "Person", "SomethingWrong", integration.ResourceType1,
				[]string{}, []string{}, ""),
			Entry("What Authorized External ID With External Property", "Person", integration.Subject5,
				integration.ResourceType4, []string{integration.Truck1}, []string{}, ""),
		)
	})

	Describe("WhoAuthorized", func() {
		var auditLogIdentifier string
		var noAuditLogEntry bool
		BeforeEach(func() {
			auditLogIdentifier = fmt.Sprintf("%v", time.Now().UnixNano())
			noAuditLogEntry = false
		})
		AfterEach(func(ctx SpecContext) {
			// Check that the audit log is present in BigQuery
			client, err := googlehelper.BqClient(ctx)
			Expect(err).To(Succeed())
			filter, err := googlehelper.FillFilterFieldsFromEnvironment()
			Expect(err).To(Succeed())
			filter.ChangeType = "type.googleapis.com/indykite.auditsink.v1beta1.WhoAuthorized"
			filter.EventSource = "AuthorizationService"
			filter.EventType = "indykite.audit.authorization.whoauthorized"
			filter.AuditLogIdentifier = auditLogIdentifier
			res, err := googlehelper.QueryAuditLog(ctx, client, &filter)
			Expect(err).To(Succeed())
			if noAuditLogEntry {
				Expect(res.Data).To(BeEmpty())
			} else {
				Expect(res.Data).To(ContainSubstring(auditLogIdentifier))
			}
		}, NodeTimeout(time.Second*10))

		DescribeTable("Who Authorization",
			func(resources []*authorizationpb.WhoAuthorizedRequest_Resource,
				subjects []string,
				policyTags []string,
				expectedError string) {
				authorizationClient, err := integration.InitConfigAuthorization()
				Expect(err).To(Succeed())

				// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
				inputParams := map[string]*authorizationpb.InputParam{
					"auditLog": {
						Value: &authorizationpb.InputParam_StringValue{
							StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
						},
					},
				}

				req := &authorizationpb.WhoAuthorizedRequest{
					Resources:   resources,
					InputParams: inputParams,
					PolicyTags:  policyTags,
				}

				resp, err := authorizationClient.WhoAuthorized(
					context.Background(),
					req,
					retry.WithMax(5),
				)

				if expectedError != "" {
					noAuditLogEntry = true
					Expect(err).To(MatchError(ContainSubstring(expectedError)))
					Expect(resp).To(BeNil())
				} else {
					Expect(err).To(Succeed())
					Expect(resp).NotTo(BeNil())

					decision := resources[0].Type
					resource := resources[0].ExternalId
					actions := resources[0].Actions
					actionMatchers := Keys{}
					for i, action := range actions {
						// First action with specific subject matches
						subjectMatcher := BeEmpty()
						if i == 0 && len(subjects) > 0 {
							elements := Elements{}
							for i, subject := range subjects {
								elements[fmt.Sprintf("%d", i)] = PointTo(MatchFields(IgnoreExtras, Fields{
									"ExternalId": Equal(subject),
								}))
							}
							subjectMatcher = MatchAllElementsWithIndex(IndexIdentity, elements)
						}
						actionMatchers[action] = PointTo(MatchFields(IgnoreExtras, Fields{
							"Subjects": subjectMatcher,
						}))
					}

					Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
						"DecisionTime": Not(BeNil()),
						"Decisions": MatchAllKeys(Keys{
							decision: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllKeys(Keys{
									resource: PointTo(MatchFields(IgnoreExtras, Fields{
										"Actions": MatchAllKeys(actionMatchers),
									})),
								}),
							})),
						}),
					})))
				}
			},
			Entry("Who Authorized", integration.ResourceWho1,
				[]string{integration.Subject3, integration.Subject2, integration.Subject1}, []string{}, ""),
			Entry("Who Authorized Resource Not Valid", integration.ResourceWho2, []string{}, []string{},
				"invalid WhoAuthorizedRequest_Resource.ExternalId: value length must be between 2 and 50 runes"),
			Entry("Who Authorized Resource Not In DB", integration.ResourceWho3, []string{}, []string{}, ""),
			Entry("Who Authorized With External Property", integration.ResourceWho4,
				[]string{integration.Subject5, integration.Subject1}, []string{}, ""),
		)
	})
})
