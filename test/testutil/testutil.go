// Copyright (c) 2026 IndyKite
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

// Package testutil holds Ginkgo helpers shared by the SDK's specs. It is test
// support only and not part of the SDK's API.
package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// SetEnv sets an environment variable for the rest of the spec and restores its
// prior value (or absence) when the spec ends, like testing.T.Setenv. It
// mutates process state, so containers whose specs call it are marked Serial.
func SetEnv(key, value string) {
	ginkgo.GinkgoHelper()
	prev, had := os.LookupEnv(key)
	gomega.Expect(os.Setenv(key, value)).To(gomega.Succeed())
	ginkgo.DeferCleanup(func() {
		if had {
			gomega.Expect(os.Setenv(key, prev)).To(gomega.Succeed())
		} else {
			gomega.Expect(os.Unsetenv(key)).To(gomega.Succeed())
		}
	})
}

// NewServer starts an httptest server running h, torn down when the spec ends.
func NewServer(h http.HandlerFunc) *httptest.Server {
	ginkgo.GinkgoHelper()
	srv := httptest.NewServer(h)
	ginkgo.DeferCleanup(srv.Close)
	return srv
}
