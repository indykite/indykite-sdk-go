// Copyright (c) 2022 IndyKite
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

package objectsv1beta1_test

import (
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Objects", func() {
	DescribeTable("SendVerificationEmailActivity - Verify attributes",
		func(fields map[string]*objects.Value, expect map[string]interface{}) {
			result, err := objects.ToMap(fields)

			Expect(err).To(Succeed())
			Expect(result).To(Equal(expect))
		},
		Entry("Nil",
			nil,
			nil,
		),
		Entry("String",
			map[string]*objects.Value{
				"Key": objects.String("String"),
			},
			map[string]interface{}{
				"Key": "String",
			},
		),
		Entry("Int64",
			map[string]*objects.Value{
				"Key": objects.Int64(-64),
			},
			map[string]interface{}{
				"Key": int64(-64),
			},
		),
		Entry("UInt64",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_UnsignedIntegerValue{UnsignedIntegerValue: 64}},
			},
			map[string]interface{}{
				"Key": uint64(64),
			},
		),
		Entry("Bool",
			map[string]*objects.Value{
				"Key": objects.Bool(true),
			},
			map[string]interface{}{
				"Key": true,
			},
		),
		Entry("Float64",
			map[string]*objects.Value{
				"Key": objects.Float64(6.4),
			},
			map[string]interface{}{
				"Key": 6.4,
			},
		),
		Entry("Time",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_ValueTime{
					ValueTime: timestamppb.New(time.Date(2020, 8, 8, 8, 8, 8, 8, time.UTC))}},
			},
			map[string]interface{}{
				"Key": timestamppb.New(
					time.Date(2020, 8, 8, 8, 8, 8, 8, time.UTC)).AsTime().UTC().Format(time.RFC3339),
			},
		),
		Entry("Duration",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_DurationValue{
					DurationValue: durationpb.New(time.Duration(10) * time.Second)}},
			},
			map[string]interface{}{
				"Key": durationpb.New(time.Duration(10) * time.Second).AsDuration().String(),
			},
		),
		Entry("Bytes",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_BytesValue{BytesValue: []byte("somefunnyjokeaboutdinosaurs")}},
			},
			map[string]interface{}{
				"Key": []byte("somefunnyjokeaboutdinosaurs"),
			},
		),
		Entry("GeoPoint",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_GeoPointValue{
					GeoPointValue: &latlng.LatLng{Latitude: 64, Longitude: 64.03}}},
			},
			map[string]interface{}{
				"Key": fmt.Sprintf("POINT (%v %v)", 64, 64.03),
			},
		),
		Entry("Null",
			map[string]*objects.Value{
				"Key": objects.Null(),
			},
			map[string]interface{}{
				"Key": nil,
			},
		),
		Entry("ObjectValueNil",
			map[string]*objects.Value{
				"Key": nil,
			},
			map[string]interface{}{
				"Key": nil,
			},
		),
		Entry("Array",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_ArrayValue{ArrayValue: &objects.ArrayValue{Values: []*objects.Value{
					objects.String("item1"),
					objects.Int64(2),
					objects.Null(),
					nil,
				}}}},
			},
			map[string]interface{}{
				"Key": []interface{}{
					"item1",
					int64(2),
					nil,
					nil,
				},
			},
		),

		Entry("Array - nil",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_ArrayValue{ArrayValue: nil}},
			},
			map[string]interface{}{
				"Key": nil,
			},
		),
		Entry("Array - nil values",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_ArrayValue{ArrayValue: &objects.ArrayValue{Values: nil}}},
			},
			map[string]interface{}{
				"Key": []interface{}{},
			},
		),
		Entry("Map",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_MapValue{MapValue: &objects.MapValue{Fields: map[string]*objects.Value{
					"item1": objects.String("item1"),
					"item2": objects.Int64(2),
					"item3": objects.Null(),
					"item4": nil,
				}}}},
			},
			map[string]interface{}{
				"Key": map[string]interface{}{
					"item1": "item1",
					"item2": int64(2),
					"item3": nil,
					"item4": nil,
				},
			},
		),
		Entry("Map - nil",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_MapValue{MapValue: nil}},
			},
			map[string]interface{}{
				"Key": nil,
			},
		),
		Entry("Map - nil values",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_MapValue{MapValue: &objects.MapValue{Fields: nil}}},
			},
			map[string]interface{}{
				"Key": map[string]interface{}{},
			},
		),
		Entry("Any",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_AnyValue{AnyValue: &anypb.Any{
					TypeUrl: "TypeUrl",
					Value:   []byte("somefunnyjokeaboutdinosaurs"),
				}}},
			},
			map[string]interface{}{
				"Key": map[string]interface{}{
					"typeUrl": "TypeUrl",
					"value":   []byte("somefunnyjokeaboutdinosaurs"),
				},
			},
		),
		Entry("Any - nil",
			map[string]*objects.Value{
				"Key": {Value: nil},
			},
			map[string]interface{}{
				"Key": nil,
			},
		),
		Entry("Any - nil value",
			map[string]*objects.Value{
				"Key": {Value: &objects.Value_AnyValue{AnyValue: nil}},
			},
			map[string]interface{}{
				"Key": nil,
			},
		),
	)
})
