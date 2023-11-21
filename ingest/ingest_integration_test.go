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

package ingest_test

import (
	"context"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta2"
	integration "github.com/indykite/indykite-sdk-go/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Ingestion", func() {
	Describe("IngestNode", func() {
		It("UpsertNode", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("UpsertNodeNoProperties", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("UpsertNodeWrongTenant", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(integration.WrongTenant, "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)
			Expect(externalID).NotTo(BeNil())
			Expect(err).To(MatchError(ContainSubstring("server was unable to complete the request")))
			Expect(resp).To(BeNil())
		})

		It("UpsertNodeWrongTenantOtherAppSpace", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(integration.WrongTenantOtherAppSpace, "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)
			Expect(externalID).NotTo(BeNil())
			Expect(err).To(MatchError(ContainSubstring("server was unable to complete the request")))
			Expect(resp).To(BeNil())
		})

		It("UpsertNodeDuplicateSameType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).To(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateDifferentTypes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Cat")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"a Digital Twin node with this id, externalId and type already exists")))
			Expect(resp2).To(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("UpsertNodeResource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.UpsertRecordAsset()
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertNodeResourceNoProp", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProp(externalID, "Asset")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateResourceSameType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProp(externalID, "Asset")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordResourceNoProp(externalID, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).To(Equal(id))

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateResourceDifferentTypes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProp(externalID, "Asset")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordResourceNoProp(externalID, "Cat")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).NotTo(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Cat")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateDifferentNodes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordResourceNoProp(externalID, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).NotTo(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateTypeNamedResource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProp(externalID, "Resource")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(externalID).NotTo(BeNil())
			Expect(err).To(MatchError(ContainSubstring("the type 'Resource' is reserved")))
			Expect(resp).To(BeNil())
		})

		It("DeleteNodeResourceNotExist", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       BeEmpty(),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_INVALID),
					}))),
				})),
			})))
		})

		It("UpsertNodeDelProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord3 := integration.DeleteRecordProperty(externalID, "Individual", "first_name")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("UpsertNodeResourceDelProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.UpsertRecordAsset()
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))

			delRecord3 := integration.DeleteRecordProperty(externalID, "Asset", "colour")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("ResourceDelNotExistingProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			delRecord3 := integration.DeleteRecordProperty(externalID, "Asset", "colour")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       BeEmpty(),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_INVALID),
					}))),
				})),
			})))
		})
	})

	Describe("IngestRelation", func() {
		It("UpsertRelation", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))
			id3 := resp3.Info.Changes[0].Id

			match := integration.GetRelationMatch(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelation(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id3),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertRelationDeleteWrongSource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			match := integration.GetRelationMatch("whatever", "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelation(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertRelationDeleteWrongSourceType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			match := integration.GetRelationMatch(externalID2, "Whatever", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelation(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertRelationWrongSourceType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation(externalID, "Whatever", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			match := integration.GetRelationMatch(externalID, "Whatever", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelation(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).NotTo(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertRelationWrongSource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation("whateveragain", "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp3).NotTo(BeNil())

			match := integration.GetRelationMatch("whateveragain", "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelation(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).NotTo(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertRelationActionLowercase", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation(externalID, "Individual", externalID2, "Asset", "can_see")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid RelationMatch.Type: value does not match regex pattern \"^[A-Z]+(?:_[A-Z]+)*$\"")))
			Expect(resp3).To(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("UpsertRelationDelProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordIndividual(os.Getenv("TENANT_ID"), "Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProp(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelation(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(2),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			match := integration.GetRelationMatch(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationProperty(match, "property1")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RESOURCE),
					}))),
				})),
			})))
		})

		It("RelationDelNotExistingProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			externalID2 := integration.GenerateRandomString(10)
			match := integration.GetRelationMatch(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationProperty(match, "property1")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.Change_DATA_TYPE_RELATION),
					}))),
				})),
			})))
		})
	})

	Describe("Stream", func() {
		It("StreamSendRecord", func() {
			var err, err2 error

			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordNoProp(externalID2, os.Getenv("TENANT_ID"), "Individual")

			records := []*ingestpb.Record{
				record, recordb,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			defer cancel()

			err2 = ingestClient.OpenStreamClient(ctx)
			Expect(err2).To(Succeed())

			for _, record := range records {
				err3 := ingestClient.SendRecord(record)
				Expect(err3).To(Succeed())
				resp, err4 := ingestClient.ReceiveResponse()
				Expect(err4).To(Succeed())
				Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"RecordId": Not(BeNil()),
					"Error":    BeNil(),
					"Info": PointTo(MatchFields(IgnoreExtras, Fields{
						"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
							"Id":       Not(BeEmpty()),
							"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
						}))),
					})),
				})))
			}

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Individual")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("StreamRecords", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordNoProp(externalID2, os.Getenv("TENANT_ID"), "Individual")

			records := []*ingestpb.Record{
				record, recordb,
			}
			responses, err := ingestClient.StreamRecords(records)
			Expect(err).To(Succeed())
			for _, response := range responses {
				Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"RecordId": Not(BeNil()),
					"Error":    BeNil(),
					"Info": PointTo(MatchFields(IgnoreExtras, Fields{
						"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
							"Id":       Not(BeEmpty()),
							"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
						}))),
					})),
				})))
			}

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Individual")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})

		It("StreamRecordsRetry", func() {
			var err error
			ingestClient, err := integration.InitConfigIngestRetry()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProp(externalID, os.Getenv("TENANT_ID"), "Individual")
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordNoProp(externalID2, os.Getenv("TENANT_ID"), "Individual")

			records := []*ingestpb.Record{
				record, recordb,
			}
			responses, err := ingestClient.StreamRecords(records)
			Expect(err).To(Succeed())
			for _, response := range responses {
				Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"RecordId": Not(BeNil()),
					"Error":    BeNil(),
					"Info": PointTo(MatchFields(IgnoreExtras, Fields{
						"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
							"Id":       Not(BeEmpty()),
							"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
						}))),
					})),
				})))
			}

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Individual")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(2),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Error":    BeNil(),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.Change_DATA_TYPE_DIGITAL_TWIN),
					}))),
				})),
			})))
		})
	})
})
