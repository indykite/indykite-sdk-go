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

package errors_test

import (
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/indykite/indykite-sdk-go/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}

var _ = Describe("IsNotFoundError", func() {
	Context("when the error is a NotFound error", func() {
		It("should return true", func() {
			err := errors.New(codes.NotFound, "not found")
			Expect(errors.IsNotFoundError(err)).To(BeTrue())
		})
	})

	Context("when the error is not a NotFound error", func() {
		It("should return false", func() {
			err := errors.New(codes.Internal, "internal")
			Expect(errors.IsNotFoundError(err)).To(BeFalse())
		})
	})

	Context("when the error is nil", func() {
		It("should return false", func() {
			Expect(errors.IsNotFoundError(nil)).To(BeFalse())
		})
	})
})
