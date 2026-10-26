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

package main

import (
	"encoding/json"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// loadShipped parses the shipped fixture files once; the specs only read the
// manifest, so they can share it.
var loadShipped = sync.OnceValues(func() (*manifest, error) {
	return load("../fixtures/config.json", "../fixtures/dataset.json")
})

// The shipped fixture files are pinned here: they must parse into the tool's
// types and carry the pieces the integration tests depend on.
var _ = Describe("shipped fixtures", func() {
	var m *manifest

	BeforeEach(func() {
		var err error
		m, err = loadShipped()
		Expect(err).NotTo(HaveOccurred())
	})

	// The config resource names and env keys the integration tests depend on
	// are all present.
	DescribeTable("carry the config resource fields the integration tests depend on",
		func(field func(*manifest) string) {
			Expect(field(m)).NotTo(BeEmpty())
		},
		Entry("kbac_policy.name", func(m *manifest) string { return m.KBACPolicy.Name }),
		Entry("kbac_policy.status", func(m *manifest) string { return m.KBACPolicy.Status }),
		Entry("ciq_policy.name", func(m *manifest) string { return m.CIQPolicy.Name }),
		Entry("knowledge_query.name", func(m *manifest) string { return m.Query.Name }),
		Entry("entity_matching_pipeline.name", func(m *manifest) string { return m.Pipeline.Name }),
		Entry("audit_signing.name", func(m *manifest) string { return m.AuditSigning.Name }),
		Entry("audit_signing.provider", func(m *manifest) string { return m.AuditSigning.Provider }),
	)

	DescribeTable("carry the env keys the integration tests depend on",
		func(key string) {
			Expect(m.Env[key]).NotTo(BeEmpty(), "env %s missing from config.json", key)
		},
		Entry(nil, "AUTHZEN_SUBJECT_TYPE"),
		Entry(nil, "AUTHZEN_SUBJECT_ID"),
		Entry(nil, "AUTHZEN_ACTION"),
		Entry(nil, "AUTHZEN_RESOURCE_TYPE"),
		Entry(nil, "AUTHZEN_RESOURCE_ID"),
		Entry(nil, "CIQ_INPUT_PARAMS"),
	)

	// The KBAC policy document is valid JSON of the expected version.
	It("carry a valid 2.0-kbac policy document", func() {
		var kbac struct {
			Meta struct {
				PolicyVersion string `json:"policy_version"`
			} `json:"meta"`
		}
		Expect(json.Unmarshal(m.KBACPolicy.Policy, &kbac)).To(Succeed(), "kbac policy JSON")
		Expect(kbac.Meta.PolicyVersion).To(Equal("2.0-kbac"))
	})

	// The dataset carries the AUTHZEN fixture subject and resource and the
	// OWNS edge that makes the decision true.
	It("carry a dataset that supports the AUTHZEN fixture decision", func() {
		Expect(m.loadedDataset.Nodes).NotTo(BeEmpty(), "dataset must carry nodes")
		Expect(m.loadedDataset.Relationships).NotTo(BeEmpty(), "dataset must carry relationships")

		subject, resource := m.Env["AUTHZEN_SUBJECT_ID"], m.Env["AUTHZEN_RESOURCE_ID"]
		ids := make([]string, 0, len(m.loadedDataset.Nodes))
		for _, n := range m.loadedDataset.Nodes {
			ids = append(ids, n.ExternalID)
		}
		Expect(ids).To(ContainElements(subject, resource), "AUTHZEN fixture ids not present in dataset")

		owns := false
		for _, r := range m.loadedDataset.Relationships {
			if r.Type == "OWNS" && r.Source.ExternalID == subject && r.Target.ExternalID == resource {
				owns = true
				break
			}
		}
		Expect(owns).To(BeTrue(),
			"dataset has no %s -[:OWNS]-> %s edge; the AUTHZEN fixture decision would be false",
			subject, resource)
	})
})
