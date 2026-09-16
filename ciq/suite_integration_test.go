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

//go:build integration

package ciq_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestIntegrationCIQ is the Ginkgo entry point when built with the integration
// tag: it runs the hermetic specs plus the live-platform ones in
// integration_test.go. The TestIntegration prefix matches the `-run` filter in
// `make integration`. See suite_test.go for the tag-less counterpart.
func TestIntegrationCIQ(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContX IQ Suite (integration)")
}
