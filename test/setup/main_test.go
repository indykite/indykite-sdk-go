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
	"testing"
)

// TestLoadFixtures pins the shipped fixture files: they must parse into the
// tool's types and carry the pieces the integration tests depend on.
func TestLoadFixtures(t *testing.T) {
	m, err := load("../fixtures/config.json", "../fixtures/dataset.json")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	assertRequiredFields(t, m)
	assertKBACPolicyValid(t, m)
	assertDatasetSupportsAuthZEN(t, m)
}

// assertRequiredFields checks the config resource names and env keys the
// integration tests depend on are all present.
func assertRequiredFields(t *testing.T, m *manifest) {
	t.Helper()
	for name, s := range map[string]string{
		"kbac_policy.name":              m.KBACPolicy.Name,
		"kbac_policy.status":            m.KBACPolicy.Status,
		"ciq_policy.name":               m.CIQPolicy.Name,
		"knowledge_query.name":          m.Query.Name,
		"entity_matching_pipeline.name": m.Pipeline.Name,
	} {
		if s == "" {
			t.Errorf("%s is empty", name)
		}
	}
	for _, key := range []string{
		"AUTHZEN_SUBJECT_TYPE", "AUTHZEN_SUBJECT_ID", "AUTHZEN_ACTION",
		"AUTHZEN_RESOURCE_TYPE", "AUTHZEN_RESOURCE_ID", "CIQ_INPUT_PARAMS",
	} {
		if m.Env[key] == "" {
			t.Errorf("env %s missing from config.json", key)
		}
	}
}

// assertKBACPolicyValid checks the KBAC policy document is valid JSON of the
// expected version.
func assertKBACPolicyValid(t *testing.T, m *manifest) {
	t.Helper()
	var kbac struct {
		Meta struct {
			PolicyVersion string `json:"policy_version"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(m.KBACPolicy.Policy, &kbac); err != nil {
		t.Fatalf("kbac policy JSON: %v", err)
	}
	if kbac.Meta.PolicyVersion != "2.0-kbac" {
		t.Errorf("kbac policy_version = %q, want 2.0-kbac", kbac.Meta.PolicyVersion)
	}
}

// assertDatasetSupportsAuthZEN checks the dataset carries the AUTHZEN fixture
// subject and resource and the OWNS edge that makes the decision true.
func assertDatasetSupportsAuthZEN(t *testing.T, m *manifest) {
	t.Helper()
	if len(m.loadedDataset.Nodes) == 0 || len(m.loadedDataset.Relationships) == 0 {
		t.Fatalf("dataset must carry nodes and relationships, got %d/%d",
			len(m.loadedDataset.Nodes), len(m.loadedDataset.Relationships))
	}

	subject, resource := m.Env["AUTHZEN_SUBJECT_ID"], m.Env["AUTHZEN_RESOURCE_ID"]
	byID := map[string]bool{}
	for _, n := range m.loadedDataset.Nodes {
		byID[n.ExternalID] = true
	}
	if !byID[subject] || !byID[resource] {
		t.Errorf("AUTHZEN fixture ids %q/%q not present in dataset", subject, resource)
	}

	for _, r := range m.loadedDataset.Relationships {
		if r.Type == "OWNS" && r.Source.ExternalID == subject && r.Target.ExternalID == resource {
			return
		}
	}
	t.Errorf("dataset has no %s -[:OWNS]-> %s edge; the AUTHZEN fixture decision would be false",
		subject, resource)
}
