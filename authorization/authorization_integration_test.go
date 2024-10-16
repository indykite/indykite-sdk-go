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

		It("IsAuthorizedDT", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(true),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTMultiple", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node2,
			}

			resources := integration.Resource9
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(true),
									})),
								}),
							})),
							resource1: PointTo(MatchFields(IgnoreExtras, Fields{
								"Actions": MatchAllKeys(Keys{
									action1: PointTo(MatchFields(IgnoreExtras, Fields{
										"Allow": Equal(true),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedTags", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			policyTags := []string{"TagOne"}

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(true),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedWrongTags", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			policyTags := []string{"TagBad"}

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTResourceNonValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource2
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			noAuditLogEntry = true
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid IsAuthorizedRequest_Resource.Type: value length must be between 2 and 50 runes")))
			Expect(resp).To(BeNil())
		})

		It("IsAuthorizedDTSubjectNonValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.NodeBad,
			}

			resources := integration.Resource2
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			noAuditLogEntry = true
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)
			Expect(err).To(MatchError(ContainSubstring("invalid IsAuthorizedRequest.Subject")))
			Expect(resp).To(BeNil())
		})

		It("IsAuthorizedDTSubjectNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.NodeNotInDB,
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTResourceNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource3
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTResourceNoSubscription", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource4
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTResourceNoOrganization", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource7
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTNoService", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource6
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedDTResourceNotLinked", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resources := integration.Resource4
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedProperty", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwinProperty := &authorizationpb.Property{
				Type:  "email",
				Value: objectpb.String(integration.EmailGood),
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(true),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedPropertyNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwinProperty := &authorizationpb.Property{
				Type:  "email",
				Value: objectpb.String(integration.EmailBad),
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedExternalID", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			externalID := &authorizationpb.ExternalID{
				Type:       "Person",
				ExternalId: integration.Subject2,
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByExternalID(
				context.Background(),
				externalID,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(true),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("IsAuthorizedExternalIDNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			externalID := &authorizationpb.ExternalID{
				Type:       "Person",
				ExternalId: "anythingwrong",
			}

			resources := integration.Resource1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByExternalID(
				context.Background(),
				externalID,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

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
										"Allow": Equal(false),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})
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

		It("WhatAuthorizedDT", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllElementsWithIndex(IndexIdentity, Elements{
									"0": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset4),
									})),
									"1": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset3),
									})),
									"2": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset1),
									})),
									"3": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset2),
									})),
									"4": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset5),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhatAuthorizedDTResourceNonValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resourcesTypes := integration.ResourceType2
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			noAuditLogEntry = true
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid WhatAuthorizedRequest_ResourceType.Type: value length must be between 2 and 50 runes")))
			Expect(resp).To(BeNil())
		})

		It("WhatAuthorizedDTSubjectNonValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.NodeBad,
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			noAuditLogEntry = true
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid DigitalTwin.Id: value length must be between 27 and 100 runes")))
			Expect(resp).To(BeNil())
		})

		It("WhatAuthorizedDTSubjectNotInDB", func() {
			var err error

			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.NodeNotInDB,
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": BeEmpty(),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhatAuthorizedDTResourceNotInDB", func() {
			var err error

			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.Node1,
			}

			resourcesTypes := integration.ResourceType3
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": BeEmpty(),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhatAuthorizedProperty", func() {
			var err error

			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwinProperty := &authorizationpb.Property{
				Type:  "email",
				Value: objectpb.String(integration.EmailGood2),
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllElementsWithIndex(IndexIdentity, Elements{
									"0": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset4),
									})),
									"1": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset3),
									})),
									"2": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset1),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhatAuthorizedPropertyNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwinProperty := &authorizationpb.Property{
				Type:  "Email",
				Value: objectpb.String(integration.EmailBad),
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": BeEmpty(),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhatAuthorizedExternalID", func() {
			var err error

			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			externalID := &authorizationpb.ExternalID{
				Type:       "Person",
				ExternalId: integration.Subject4,
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByExternalID(
				context.Background(),
				externalID,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": MatchAllElementsWithIndex(IndexIdentity, Elements{
									"0": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset4),
									})),
									"1": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset3),
									})),
									"2": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset1),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhatAuthorizedExternalIDNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			externalID := &authorizationpb.ExternalID{
				Type:       "Person",
				ExternalId: "SomethingWrong",
			}

			resourcesTypes := integration.ResourceType1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByExternalID(
				context.Background(),
				externalID,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resourcesTypes[0].Type
			action := resourcesTypes[0].Actions[0]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Actions": MatchAllKeys(Keys{
							action: PointTo(MatchFields(IgnoreExtras, Fields{
								"Resources": BeEmpty(),
							})),
						}),
					})),
				}),
			})))
		})
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

		It("WhoAuthorized", func() {
			var err error

			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			resources := integration.ResourceWho1
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

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

			decision := resources[0].Type
			resource := resources[0].ExternalId
			action0 := resources[0].Actions[0]
			action1 := resources[0].Actions[1]

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Resources": MatchAllKeys(Keys{
							resource: PointTo(MatchFields(IgnoreExtras, Fields{
								"Actions": MatchAllKeys(Keys{
									action0: PointTo(MatchFields(IgnoreExtras, Fields{
										"Subjects": MatchAllElementsWithIndex(IndexIdentity, Elements{
											"0": PointTo(MatchFields(IgnoreExtras, Fields{
												"ExternalId": Equal(integration.Subject1),
											})),
											"1": PointTo(MatchFields(IgnoreExtras, Fields{
												"ExternalId": Equal(integration.Subject3),
											})),
											"2": PointTo(MatchFields(IgnoreExtras, Fields{
												"ExternalId": Equal(integration.Subject2),
											})),
										}),
									})),
									action1: PointTo(MatchFields(IgnoreExtras, Fields{
										"Subjects": BeEmpty(),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})

		It("WhoAuthorizedResourceNotValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			resources := integration.ResourceWho2
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			noAuditLogEntry = true
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

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

			Expect(err).To(MatchError(ContainSubstring(
				"invalid WhoAuthorizedRequest_Resource.ExternalId: value length must be between 2 and 50 runes")))
			Expect(resp).To(BeNil())
		})

		It("WhoAuthorizedResourceNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			resources := integration.ResourceWho3
			// To make sure that the proper audit log was queried from BigQuery, need to add a unique identifier.
			inputParams := map[string]*authorizationpb.InputParam{
				"auditLog": {
					Value: &authorizationpb.InputParam_StringValue{
						StringValue: fmt.Sprintf("\"%v\"", auditLogIdentifier),
					},
				},
			}
			var policyTags []string

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

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())

			decision := resources[0].Type
			resource := resources[0].ExternalId
			action0 := resources[0].Actions[0]
			action1 := resources[0].Actions[1]

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"DecisionTime": Not(BeNil()),
				"Decisions": MatchAllKeys(Keys{
					decision: PointTo(MatchFields(IgnoreExtras, Fields{
						"Resources": MatchAllKeys(Keys{
							resource: PointTo(MatchFields(IgnoreExtras, Fields{
								"Actions": MatchAllKeys(Keys{
									action0: PointTo(MatchFields(IgnoreExtras, Fields{
										"Subjects": BeEmpty(),
									})),
									action1: PointTo(MatchFields(IgnoreExtras, Fields{
										"Subjects": BeEmpty(),
									})),
								}),
							})),
						}),
					})),
				}),
			})))
		})
	})
})
