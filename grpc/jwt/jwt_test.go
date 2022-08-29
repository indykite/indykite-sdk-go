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

package jwt_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"

	. "github.com/indykite/jarvis-sdk-go/grpc/jwt"
)

var (
	IntTestAppAgentID = uuid.MustParse("696e6479-6b69-4465-8000-050f00000000")
)

var _ = Describe("CreateTokenSourceFrom", func() {
	It("invalid string key", func() {
		appAgentID := IntTestAppAgentID

		invalidPrivateKeyJWK := "invalid_key"
		resp, err := CreateTokenSourceFrom(invalidPrivateKeyJWK, appAgentID.String())

		Expect(resp).To(BeNil())
		Expect(err).To(MatchError(ContainSubstring("failed to unmarshal JSON into key hint")))
	})

	DescribeTable(
		"validate JSON correctly create tokenSource",
		func(privateKeyJWK interface{}) {
			resp, err := CreateTokenSourceFrom(privateKeyJWK, IntTestAppAgentID.String())
			Expect(resp).To(Not(BeNil()))
			Expect(err).To(Succeed())
		},
		//nolint:lll
		Entry("valid string key", "{\"kty\":\"EC\",\"d\":\"2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s\",\"use\":\"sig\",\"crv\":\"P-256\",\"kid\":\"vDUXHBZcRw1KyFPyB0EI2XLBzyP9iGyfvaSX3MNtUlk\",\"x\":\"Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA\",\"y\":\"DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4\",\"alg\":\"ES256\"}"),
		Entry("valid map key", map[string]interface{}{
			"kty": "EC",
			"d":   "2to-_wtohfn2PAgHr3RHQbhDf8g9zy6ndr05ZS-hS8s",
			"use": "sig",
			"crv": "P-256",
			"kid": "vDUXHBZcRw1KyFPyB0EI2XLBzyP9iGyfvaSX3MNtUlk",
			"x":   "Cn2tSCxcQYVKuexBTzqRShvrJG8eQeZUq0ISIp9wXSA",
			"y":   "DVSlYTLzns37LmjdscBA8q5ko1N8CZ-ETwviAJ78vW4",
			"alg": "ES256",
		},
		),
	)
})
