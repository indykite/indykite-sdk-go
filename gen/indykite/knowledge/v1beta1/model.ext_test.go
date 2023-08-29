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

package knowledgev1beta1_test

import (
	knowledgev1beta1 "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta1"
	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Node test", func() {

	testNode := &knowledgev1beta1.Node{
		Properties: []*knowledgev1beta1.Property{
			{
				Key: "name",
				Value: &objects.Value{
					Value: &objects.Value_StringValue{StringValue: "John"},
				},
			},
			{
				Key: "email",
				Value: &objects.Value{
					Value: &objects.Value_StringValue{StringValue: "test@test.com"},
				},
			},
		},
	}

	It("GetProperty that exists", func() {
		v, ok := testNode.GetProperty("name")
		Expect(ok).To(BeTrue())
		Expect(v).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Value": PointTo(MatchFields(IgnoreExtras, Fields{
				"StringValue": Equal("John"),
			})),
		})))
	})

	It("GetProperty that doesn't exist", func() {
		v, ok := testNode.GetProperty("something")
		Expect(ok).To(BeFalse())
		Expect(v).To(BeNil())
	})

})
