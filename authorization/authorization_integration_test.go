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

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	authorizationpb "github.com/indykite/indykite-sdk-go/gen/indykite/authorization/v1beta1"
	objectpb "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Authorized", func() {
	Describe("IsAuthorized", func() {
		It("IsAuthorizedDT", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin2,
			}

			resources := integration.Resource9
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource1
			inputParams := map[string]*authorizationpb.InputParam{}
			policyTags := []string{"Tag1"}

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource1
			inputParams := map[string]*authorizationpb.InputParam{}
			policyTags := []string{"TagBad"}

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource2
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwinBad,
			}

			resources := integration.Resource2
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
			)
			Expect(err).To(MatchError(ContainSubstring("invalid IsAuthorizedRequest.Subject")))
			Expect(resp).To(BeNil())
		})

		It("IsAuthorizedDTSubjectNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.DigitalTwinNotInDB,
			}

			resources := integration.Resource1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource3
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource4
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource7
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource6
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resources := integration.Resource4
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorized(
				context.Background(),
				digitalTwin,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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

		It("IsAuthorizedTokenNonValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			token := integration.TokenBad
			resources := integration.Resource1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByToken(
				context.Background(),
				token,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
			)
			Expect(err).To(MatchError(ContainSubstring("invalid JWT")))
			Expect(resp).To(BeNil())
		})

		It("IsAuthorizedTokenNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			token := integration.TokenGoodFormat
			resources := integration.Resource1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.IsAuthorizedByToken(
				context.Background(),
				token,
				resources,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
		It("WhatAuthorizedDT", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			digitalTwin := &authorizationpb.DigitalTwin{
				Id: integration.DigitalTwin1,
			}

			resourcesTypes := integration.ResourceType1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
										"ExternalId": Equal(integration.Asset1),
									})),
									"1": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset2),
									})),
									"2": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset3),
									})),
									"3": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset4),
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
				Id: integration.DigitalTwin1,
			}

			resourcesTypes := integration.ResourceType2
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwinBad,
			}

			resourcesTypes := integration.ResourceType1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwinNotInDB,
			}

			resourcesTypes := integration.ResourceType1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Id: integration.DigitalTwin1,
			}

			resourcesTypes := integration.ResourceType3
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorized(
				context.Background(),
				digitalTwin,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
				Value: objectpb.String(integration.EmailGood),
			}

			resourcesTypes := integration.ResourceType1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
										"ExternalId": Equal(integration.Asset1),
									})),
									"1": PointTo(MatchFields(IgnoreExtras, Fields{
										"ExternalId": Equal(integration.Asset2),
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
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByProperty(
				context.Background(),
				digitalTwinProperty,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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

		It("WhatAuthorizedTokenNonValid", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			resourcesTypes := integration.ResourceType1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByToken(
				context.Background(),
				integration.TokenBad,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
			)

			Expect(err).To(MatchError(ContainSubstring("invalid JWT")))
			Expect(resp).To(BeNil())
		})

		It("WhatAuthorizedPropertyNotInDB", func() {
			var err error
			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			resourcesTypes := integration.ResourceType1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			resp, err := authorizationClient.WhatAuthorizedByToken(
				context.Background(),
				integration.TokenGoodFormat,
				resourcesTypes,
				inputParams,
				policyTags,
				retry.WithMax(2),
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
		It("WhoAuthorized", func() {
			var err error

			authorizationClient, err := integration.InitConfigAuthorization()
			Expect(err).To(Succeed())

			resources := integration.ResourceWho1
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			req := &authorizationpb.WhoAuthorizedRequest{
				Resources:   resources,
				InputParams: inputParams,
				PolicyTags:  policyTags,
			}

			resp, err := authorizationClient.WhoAuthorized(
				context.Background(),
				req,
				retry.WithMax(2),
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
												"ExternalId": Equal(integration.Subject2),
											})),
											"2": PointTo(MatchFields(IgnoreExtras, Fields{
												"ExternalId": Equal(integration.Subject3),
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
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			req := &authorizationpb.WhoAuthorizedRequest{
				Resources:   resources,
				InputParams: inputParams,
				PolicyTags:  policyTags,
			}

			resp, err := authorizationClient.WhoAuthorized(
				context.Background(),
				req,
				retry.WithMax(2),
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
			inputParams := map[string]*authorizationpb.InputParam{}
			var policyTags []string

			req := &authorizationpb.WhoAuthorizedRequest{
				Resources:   resources,
				InputParams: inputParams,
				PolicyTags:  policyTags,
			}

			resp, err := authorizationClient.WhoAuthorized(
				context.Background(),
				req,
				retry.WithMax(2),
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
