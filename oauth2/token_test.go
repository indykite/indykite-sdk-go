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

package oauth2_test

import (
	"context"
	"encoding/json"
	"os"

	apicfg "github.com/indykite/indykite-sdk-go/grpc/config"
	"github.com/indykite/indykite-sdk-go/oauth2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Bearer token", func() {

	BeforeEach(func() {
		// Prepare test data
		// #nosec: G304: Potential file inclusion via variable
		data, err := os.ReadFile("../grpc/testdata/valid_jwt_as_string.json")
		Expect(err).To(Succeed())
		// Quickly parse JSON to avoid hardcoded values
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(data, &jsonData)
		Expect(err).To(Succeed())

		credentialsLoaders = append(credentialsLoaders, apicfg.StaticCredentialsJSON(data))
	})

	It("Generate token", func() {
		tokenSource, err := oauth2.GetRefreshableTokenSource(context.Background(), credentialsLoaders)
		Expect(err).To(Succeed())
		Expect(tokenSource.Token()).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"AccessToken": Not(BeNil()),
		})))
	})

	It("Create http client", func() {
		httpClient, err := oauth2.GetHTTPClient(context.Background(), credentialsLoaders)
		Expect(err).To(Succeed())
		Expect(httpClient).To(Not(BeNil()))
	})
})
