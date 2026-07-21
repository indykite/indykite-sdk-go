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
	"context"
	"log"
	"os"

	"github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/config"
)

func main() {
	ctx := context.Background()

	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	} else if region := os.Getenv("INDYKITE_REGION"); region != "" {
		opts = append(opts, indykite.WithRegion(region))
	}

	// Control plane (Service Account) for config, runtime (App Agent) for the rest.
	admin, err := indykite.NewAdminFromEnv(ctx, opts...)
	if err != nil {
		log.Fatalf("e2e-driver: build admin client: %v", err)
	}
	ik, err := indykite.NewClientFromEnv(ctx, opts...)
	if err != nil {
		log.Fatalf("e2e-driver: build runtime client: %v", err)
	}

	projectID := os.Getenv("INDYKITE_PROJECT_ID")
	if projectID == "" {
		log.Fatal("e2e-driver: set INDYKITE_PROJECT_ID")
	}

	scenario := Scenario{
		Policy: config.CreateAuthorizationPolicy{
			ProjectID:   projectID,
			Name:        "e2e-can-read-doc",
			Description: "e2e-driver: Person may READ a Document they own.",
			Status:      config.StatusActive,
			Policy: `{"meta":{"policy_version":"2.0-kbac"},"subject":{"type":"Person"},` +
				`"actions":["READ"],"resource":{"type":"Document"},` +
				`"condition":{"cypher":"MATCH (subject:Person)-[:OWNS]->(resource:Document)"}}`,
			Tags: []string{"e2e"},
		},
		Nodes: []capture.UpsertNode{
			{Node: capture.Node{ExternalID: "alice", Type: "Person"}, IsIdentity: true},
			{Node: capture.Node{ExternalID: "doc-1", Type: "Document"}},
		},
		Relationship: &capture.Relationship{
			Type:   "OWNS",
			Source: &capture.Node{ExternalID: "alice", Type: "Person"},
			Target: &capture.Node{ExternalID: "doc-1", Type: "Document"},
		},
		Subject:   authzen.NewNode("Person", "alice"),
		Action:    "READ",
		Resource:  authzen.NewNode("Document", "doc-1"),
		WantAllow: true,
		QueryID:   os.Getenv("CIQ_QUERY_ID"), // optional
	}

	steps, runErr := NewDriver(admin, ik).Run(ctx, &scenario)
	for _, s := range steps {
		status := "PASS"
		if !s.OK {
			status = "FAIL"
		}
		log.Printf("[%s] %-20s %s", status, s.Name, s.Detail)
	}
	if runErr != nil {
		log.Fatalf("e2e-driver: aborted: %v", runErr)
	}
	if !AllOK(steps) {
		os.Exit(1)
	}
	log.Print("e2e-driver: all steps passed")
}
