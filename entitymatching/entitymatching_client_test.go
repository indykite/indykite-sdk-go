// Copyright (c) 2024 IndyKite
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

package entitymatching_test

import (
	"github.com/indykite/indykite-sdk-go/entitymatching"
	entitymatchingmock "github.com/indykite/indykite-sdk-go/test/entitymatching/v1beta1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigNode", func() {
	var (
		mockClient           *entitymatchingmock.MockEntityMatchingAPIClient
		entitymatchingClient *entitymatching.Client
	)

	BeforeEach(func() {
		mockClient = &entitymatchingmock.MockEntityMatchingAPIClient{}
	})

	Describe("Client", func() {
		It("New", func() {
			var err error
			entitymatchingClient, err = entitymatching.NewTestClient(mockClient)
			Ω(err).To(Succeed())
			Expect(entitymatchingClient).To(Not(BeNil()))

			if entitymatchingClient != nil {
				Expect(entitymatchingClient.Close()).To(Succeed())
			}
		})
	})
})
