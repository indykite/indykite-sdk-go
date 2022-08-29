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

package grpc_test

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"

	"github.com/indykite/jarvis-sdk-go/grpc/config"

	"github.com/indykite/jarvis-sdk-go/grpc/internal"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/indykite/jarvis-sdk-go/grpc"
)

var _ = Describe("Test JSON and YAML to JSON credentials", func() {
	It("Empty input returns nils without errors", func() {
		dialSettings := &internal.DialSettings{}
		for _, v := range []grpc.ClientOption{
			grpc.WithEndpoint("test"),
			grpc.WithCredentialsLoader(config.StaticCredentialsJSON([]byte(""))),
		} {
			v.Apply(dialSettings)
		}
		dialOption, cfg, err := dialSettings.Build(context.Background())
		Expect(dialOption).To(HaveLen(6))
		Expect(cfg).To(BeNil())
		Expect(err).To(Succeed())
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

			dialSettings := &internal.DialSettings{}
			for _, v := range []grpc.ClientOption{
				grpc.WithCredentialsLoader(config.StaticCredentialsJSON(data)),
			} {
				v.Apply(dialSettings)
			}

			dialOption, cfg, err := dialSettings.Build(context.Background())
			Expect(err).To(Succeed())
			Expect(dialOption).To(HaveLen(7))
			Expect(cfg).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"AppAgentID": Equal(jsonData["appAgentId"]),
				"Endpoint":   Equal(jsonData["endpoint"]),
			})))
		},
		Entry("Valid JSON with JWT as a JSON map", "testdata/valid_jwt_as_map.json"),
		Entry("Valid JSON with JWT as a string", "testdata/valid_jwt_as_string.json"),
	)

	It("Valid YAML converted to JSON with JWT as a string", func() {
		// Prepare test data
		data, err := ioutil.ReadFile("testdata/valid_jwt_as_json_string.yaml")
		Expect(err).To(Succeed())
		data, err = yaml.YAMLToJSON(data)
		Expect(err).To(Succeed())
		// Quickly parse JSON to avoid hardcoded values
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(data, &jsonData)
		Expect(err).To(Succeed())

		dialSettings := &internal.DialSettings{}
		for _, v := range []grpc.ClientOption{
			grpc.WithCredentialsLoader(config.StaticCredentialsJSON(data)),
		} {
			v.Apply(dialSettings)
		}
		dialOption, cfg, err := dialSettings.Build(context.Background())

		Expect(dialSettings.Endpoint).To(Equal("dns:///" + jsonData["endpoint"].(string)))
		Expect(dialOption).NotTo(BeNil())
		Expect(cfg).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"AppAgentID": Equal(jsonData["appAgentId"]),
			"Endpoint":   Equal(jsonData["endpoint"]),
		})))
		Expect(err).To(Succeed())
	})

	It("JSON file does not overwrite Dial Settings", func() {
		// Prepare test data
		data, err := ioutil.ReadFile("testdata/valid_jwt_as_map.json")
		Expect(err).To(Succeed())
		// Quickly parse JSON to avoid hardcoded values
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(data, &jsonData)
		Expect(err).To(Succeed())

		endpoint := "grpc_endpoint.example.com"
		dialSettings := &internal.DialSettings{}
		for _, v := range []grpc.ClientOption{
			grpc.WithEndpoint(endpoint),
			grpc.WithCredentialsLoader(config.StaticCredentialsJSON(data)),
		} {
			v.Apply(dialSettings)
		}
		dialOption, cfg, err := dialSettings.Build(context.Background())

		// Endpoint from file does not overwrite endpoint set in DialSettings
		Expect(dialSettings.Endpoint).To(Equal("dns:///" + endpoint))
		Expect(dialOption).NotTo(BeNil())
		Expect(cfg).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"AppAgentID": Equal(jsonData["appAgentId"]),
		})))
		Expect(err).To(Succeed())
	})

	It("JSON file is not Service Account", func() {
		// Prepare test data
		data, err := ioutil.ReadFile("testdata/valid_jwt_as_map.json")
		Ω(err).To(Succeed())

		endpoint := "grpc_endpoint.example.com"
		dialSettings := &internal.DialSettings{}
		for _, v := range []grpc.ClientOption{
			grpc.WithEndpoint(endpoint),
			grpc.WithCredentialsLoader(config.StaticCredentialsJSON(data)),
			grpc.WithServiceAccount(),
		} {
			v.Apply(dialSettings)
		}
		_, _, err = dialSettings.Build(context.Background())
		Ω(err).To(MatchError("empty serviceAccountId"))
	})

	It("JSON file is not AppAgent Account", func() {
		// Prepare test data
		data, err := ioutil.ReadFile("testdata/valid_jwt_service.json")
		Ω(err).To(Succeed())

		endpoint := "grpc_endpoint.example.com"
		dialSettings := &internal.DialSettings{}
		for _, v := range []grpc.ClientOption{
			grpc.WithEndpoint(endpoint),
			grpc.WithCredentialsLoader(config.StaticCredentialsJSON(data)),
		} {
			v.Apply(dialSettings)
		}
		_, _, err = dialSettings.Build(context.Background())
		Ω(err).To(MatchError("empty appAgentId"))
	})

	It("Service Account JSON", func() {
		// Prepare test data
		data, err := ioutil.ReadFile("testdata/valid_jwt_service.json")
		Ω(err).To(Succeed())
		// Quickly parse JSON to avoid hardcoded values
		var jsonData map[string]interface{}
		err = json.Unmarshal(data, &jsonData)
		Ω(err).To(Succeed())

		dialSettings := &internal.DialSettings{}
		for _, v := range []grpc.ClientOption{
			grpc.WithCredentialsLoader(config.StaticCredentialsJSON(data)),
			grpc.WithServiceAccount(),
		} {
			v.Apply(dialSettings)
		}
		dialOption, cfg, err := dialSettings.Build(context.Background())
		Ω(err).To(Succeed())

		// Endpoint from file does not overwrite endpoint set in DialSettings
		Ω(dialSettings.Endpoint).To(Equal("dns:///jarvis.local"))
		Ω(dialOption).To(HaveLen(7))
		Ω(cfg).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ServiceAccountID": Equal("20092d59-e4b1-4828-806b-f084f653797b"),
		})))
	})
})
