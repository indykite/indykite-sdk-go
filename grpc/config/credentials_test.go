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

package config_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/indykite/jarvis-sdk-go/grpc/config"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("UnmarshalCredentialConfig", func() {
	It("Empty input returns nils without errors", func() {

		config, err := UnmarshalCredentialConfig([]byte(""))
		Expect(config).To(BeNil())
		Expect(err).To(Succeed())
	})

	It("Invalid json", func() {
		config, err := UnmarshalCredentialConfig([]byte("{invalid}"))
		Expect(config).To(BeNil())
		Expect(err).To(HaveOccurred())
	})

	DescribeTable(
		"Valid JSON correctly set JWT token",
		func(filePath string) {
			// Prepare test data
			// #nosec: G304: Potential file inclusion via variable
			data, err := ioutil.ReadFile(filePath)
			Expect(err).To(Succeed())
			// Quickly parse JSON to avoid hardcoded values
			jsonData := make(map[string]interface{})
			err = json.Unmarshal(data, &jsonData)
			Expect(err).To(Succeed())

			config, err := UnmarshalCredentialConfig(data)

			Expect(err).To(Succeed())
			Expect(config).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"AppAgentID": Equal(jsonData["appAgentId"]),
				"Endpoint":   Equal(jsonData["endpoint"]),
			})))
		},
		Entry("Valid JSON with JWT as a JSON map", "../testdata/valid_jwt_as_map.json"),
		Entry("Valid JSON with JWT as a string", "../testdata/valid_jwt_as_string.json"),
	)
})
