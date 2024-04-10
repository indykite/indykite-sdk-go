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

//go:build !integration

package ingest_test

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
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

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeNoProperties", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateSameType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordNoProperty(externalID, "Individual")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).To(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateDifferentTypes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordNoProperty(externalID, "Cat")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).NotTo(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecord2 := integration.DeleteRecord(externalID, "Cat")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord2,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeResource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.UpsertRecordNodeAsset()
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeResourceNoProp", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProperty(externalID, "Asset")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateResourceSameType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProperty(externalID, "Asset")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordResourceNoProperty(externalID, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).To(Equal(id))

			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateResourceDifferentTypes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProperty(externalID, "Asset")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordResourceNoProperty(externalID, "Cat")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).NotTo(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Cat")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateDifferentNodes", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).NotTo(BeNil())
			id := resp.Info.Changes[0].Id

			recordb := integration.CreateRecordResourceNoProperty(externalID, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).NotTo(BeNil())
			id2 := resp2.Info.Changes[0].Id
			Expect(id2).NotTo(Equal(id))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDuplicateTypeNamedResource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordResourceNoProperty(externalID, "Resource")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
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
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       BeEmpty(),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeDelProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord3 := integration.DeleteRecordWithProperty(externalID, "Individual", "first_name")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertNodeResourceDelProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.UpsertRecordNodeAsset()
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecord3 := integration.DeleteRecordWithProperty(externalID, "Asset", "colour")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			delRecord := integration.DeleteRecord(externalID, "Asset")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("ResourceDelNotExistingProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			delRecord3 := integration.DeleteRecordWithProperty(externalID, "Asset", "colour")
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       BeEmpty(),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})
	})

	Describe("IngestRelationship", func() {
		It("UpsertRelationship", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			id := resp.Info.Changes[0].Id
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelationship(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))
			id3 := resp3.Info.Changes[0].Id

			match := integration.GetRelationship(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationship(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id3),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertRelationshipDeleteWrongSource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelationship(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			match := integration.GetRelationship("whatever", "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationship(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertRelationshipDeleteWrongSourceType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelationship(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			match := integration.GetRelationship(externalID2, "Whatever", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationship(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertRelationshipWrongSourceType", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelationship(externalID, "Whatever", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			match := integration.GetRelationship(externalID, "Whatever", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationship(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).NotTo(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertRelationshipWrongSource", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id
			rel := integration.GenerateRandomString(10)

			record3 := integration.CreateRecordRelationship(rel, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp3).NotTo(BeNil())

			match := integration.GetRelationship(rel, "Individual", externalID2, "Asset", "CAN_SEE")
			delRecord3 := integration.DeleteRecordRelationship(match)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).NotTo(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertRelationshipActionLowercase", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelationship(externalID, "Individual", externalID2, "Asset", "can_see")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)
			Expect(err).To(MatchError(ContainSubstring(
				"invalid Relationship.Type: value does not match regex pattern \"^[A-Z]+(?:_[A-Z]+)*$\"")))
			Expect(resp3).To(BeNil())

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("UpsertRelationshipDelProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			record, externalID := integration.CreateRecordNodeIndividual("Employee")
			resp, err := ingestClient.IngestRecord(
				context.Background(),
				record,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id := resp.Info.Changes[0].Id

			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordResourceNoProperty(externalID2, "Asset")
			resp2, err := ingestClient.IngestRecord(
				context.Background(),
				recordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(resp2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
			id2 := resp2.Info.Changes[0].Id

			record3 := integration.CreateRecordRelationship(externalID, "Individual", externalID2, "Asset", "CAN_SEE")
			resp3, err := ingestClient.IngestRecord(
				context.Background(),
				record3,
				retry.WithMax(5),
			)

			Expect(err).To(Succeed())
			Expect(resp3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			delRecord3 := integration.DeleteRecordRelationshipProperty(
				externalID,
				"Individual",
				externalID2,
				"Asset",
				"CAN_SEE",
				"property1",
			)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
					}))),
				})),
			})))

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Asset")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Equal(id2),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("RelationshipDelNotExistingProperty", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			externalID2 := integration.GenerateRandomString(10)
			delRecord3 := integration.DeleteRecordRelationshipProperty(
				externalID,
				"Individual",
				externalID2,
				"Asset",
				"CAN_SEE",
				"property1",
			)
			del3, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord3,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del3).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_RELATIONSHIP),
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
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordNoProperty(externalID2, "Individual")

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
					"Info": PointTo(MatchFields(IgnoreExtras, Fields{
						"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
							"Id":       Not(BeEmpty()),
							"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
						}))),
					})),
				})))
			}

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Individual")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("StreamRecords", func() {
			var err error
			ingestClient, err := integration.InitConfigIngest()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordNoProperty(externalID2, "Individual")

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
							"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
						}))),
					})),
				})))
			}

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Individual")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})

		It("StreamRecordsRetry", func() {
			var err error
			ingestClient, err := integration.InitConfigIngestRetry()
			Expect(err).To(Succeed())

			externalID := integration.GenerateRandomString(10)
			record := integration.CreateRecordNoProperty(externalID, "Individual")
			externalID2 := integration.GenerateRandomString(10)
			recordb := integration.CreateRecordNoProperty(externalID2, "Individual")

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
							"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
						}))),
					})),
				})))
			}

			delRecord := integration.DeleteRecord(externalID, "Individual")
			del, err := ingestClient.IngestRecord(
				context.Background(),
				delRecord,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))

			delRecordb := integration.DeleteRecord(externalID2, "Individual")
			del2, err := ingestClient.IngestRecord(
				context.Background(),
				delRecordb,
				retry.WithMax(5),
			)
			Expect(err).To(Succeed())
			Expect(del2).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"RecordId": Not(BeNil()),
				"Info": PointTo(MatchFields(IgnoreExtras, Fields{
					"Changes": ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
						"Id":       Not(BeEmpty()),
						"DataType": Equal(ingestpb.DataType_DATA_TYPE_NODE),
					}))),
				})),
			})))
		})
	})
})
