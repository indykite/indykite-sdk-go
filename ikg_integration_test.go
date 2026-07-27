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

package indykite_test

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
)

// Policies, data and decisions propagate asynchronously; poll before failing.
const (
	e2ePollAttempts = 12
	e2ePollInterval = 5 * time.Second
)

// TestIntegrationIKGEndToEnd exercises an ACTUAL IKG in the project with no
// pre-provisioned fixtures — only credentials and PROJECT_ID are required:
//
//	seed    control plane: KBAC policy + CIQ read policy + knowledge query
//	ingest  runtime plane: Person -[:OWNS]-> Server into the IKG
//	assert  AuthZEN decision over the ingested graph (positive + negative)
//	assert  ContX IQ knowledge query returns the ingested node
//	cleanup graph data and config resources (always, via t.Cleanup)
func TestIntegrationIKGEndToEnd(t *testing.T) {
	if os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS_FILE") == "" {
		t.Skip("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE] not set")
	}
	if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
		t.Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
	}
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		t.Skip("PROJECT_ID not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	admin, err := indykite.NewAdminFromEnv(ctx, opts...)
	if err != nil {
		t.Fatalf("NewAdminFromEnv: %v", err)
	}
	cli, err := indykite.NewClientFromEnv(ctx, opts...)
	if err != nil {
		t.Fatalf("NewClientFromEnv: %v", err)
	}

	unique := strconv.FormatInt(time.Now().UnixNano(), 10)
	action := "SDK_IT_CAN_USE"

	// --- seed: KBAC policy (Person may act on a Server they OWN) ---
	kbacPolicy := `{
	  "meta": {"policy_version": "2.0-kbac"},
	  "subject": {"type": "Person"},
	  "actions": ["` + action + `"],
	  "resource": {"type": "Server"},
	  "condition": {"cypher": "MATCH (subject:Person)-[:OWNS]->(resource:Server)"}
	}`
	kbac, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: projectID,
		Name:      "sdk-it-e2e-kbac-" + unique,
		Policy:    kbacPolicy,
		Status:    config.StatusActive,
	})
	if err != nil {
		t.Fatalf("create KBAC policy: %v", err)
	}
	t.Cleanup(func() { _ = admin.AuthorizationPolicies().Delete(context.Background(), kbac.ID, "") })

	// --- seed: CIQ read policy + knowledge query (find a Server by external id).
	// The _Application subject binds to the caller's own application node, so
	// executing needs only the App Agent key — no user Bearer token.
	ciqPolicy := `{
	  "meta": {"policy_version": "1.0-ciq"},
	  "subject": {"type": "_Application"},
	  "condition": {
	    "cypher": "MATCH (subject:_Application), (server:Server)",
	    "filter": [
	      {"operator": "=", "attribute": "subject.external_id", "value": "$_appId"},
	      {"operator": "=", "attribute": "server.external_id", "value": "$server_external_id"}
	    ]
	  },
	  "allowed_reads": {"nodes": ["server", "server.*"]}
	}`
	ciqPol, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: projectID,
		Name:      "sdk-it-e2e-ciq-" + unique,
		Policy:    ciqPolicy,
		Status:    config.StatusActive,
	})
	if err != nil {
		t.Fatalf("create CIQ policy: %v", err)
	}
	t.Cleanup(func() { _ = admin.AuthorizationPolicies().Delete(context.Background(), ciqPol.ID, "") })

	kq, err := admin.KnowledgeQueries().Create(ctx, &config.CreateKnowledgeQuery{
		ProjectID: projectID,
		Name:      "sdk-it-e2e-kq-" + unique,
		Query:     `{"nodes": ["server"]}`,
		Status:    config.StatusActive,
		PolicyID:  ciqPol.ID,
	})
	if err != nil {
		t.Fatalf("create knowledge query: %v", err)
	}
	t.Cleanup(func() { _ = admin.KnowledgeQueries().Delete(context.Background(), kq.ID, "") })

	// --- ingest: the graph the policy and query evaluate over ---
	person := capture.Node{ExternalID: "sdk-it-e2e-person-" + unique, Type: "Person"}
	stranger := capture.Node{ExternalID: "sdk-it-e2e-stranger-" + unique, Type: "Person"}
	server := capture.Node{ExternalID: "sdk-it-e2e-server-" + unique, Type: "Server"}
	owns := capture.Relationship{Type: "OWNS", Source: &person, Target: &server}

	t.Cleanup(func() {
		cctx := context.Background()
		_, _ = cli.Capture().DeleteRelationships(cctx, owns)
		_, _ = cli.Capture().DeleteNodes(cctx, person, stranger, server)
	})

	if _, err = cli.Capture().UpsertNodes(ctx,
		capture.UpsertNode{Node: person, IsIdentity: true},
		capture.UpsertNode{Node: stranger, IsIdentity: true},
		// No region property: the fixture knowledge query (region=eu) must not
		// see this transient server when packages run in parallel (plain
		// `go test ./...` without the make targets' -p 1).
		capture.UpsertNode{Node: server},
	); err != nil {
		t.Fatalf("ingest nodes: %v", err)
	}
	if _, err = cli.Capture().UpsertRelationships(ctx, owns); err != nil {
		t.Fatalf("ingest relationship: %v", err)
	}

	// --- assert: AuthZEN decision over the real graph ---
	subject := authzen.NewNode("Person", person.ExternalID)
	resource := authzen.NewNode("Server", server.ExternalID)

	allowed := pollUntil(ctx, t, "owner is allowed", func() (bool, error) {
		return cli.AuthZEN().Allowed(ctx, subject, action, resource)
	})
	if !allowed {
		t.Fatalf("owner %s never got %s on %s", person.ExternalID, action, server.ExternalID)
	}

	strangerAllowed, err := cli.AuthZEN().Allowed(ctx,
		authzen.NewNode("Person", stranger.ExternalID), action, resource)
	if err != nil {
		t.Fatalf("Allowed (stranger): %v", err)
	}
	if strangerAllowed {
		t.Errorf("stranger %s must not get %s on %s", stranger.ExternalID, action, server.ExternalID)
	}

	// --- assert: the knowledge query reads the ingested node back ---
	found := pollUntil(ctx, t, "knowledge query returns the server", func() (bool, error) {
		rows, qErr := cli.CIQ().All(ctx, ciq.ExecuteRequest{
			ID:          kq.ID,
			InputParams: map[string]any{"server_external_id": server.ExternalID},
		})
		if qErr != nil {
			return false, qErr
		}
		for _, row := range rows {
			raw, _ := json.Marshal(row)
			if strings.Contains(string(raw), server.ExternalID) {
				return true, nil
			}
		}
		return false, nil
	})
	if !found {
		t.Fatalf("knowledge query never returned server %s", server.ExternalID)
	}
}

// pollUntil retries cond every e2ePollInterval until it is true, the attempts
// are exhausted, or ctx ends. Errors are logged and retried (propagation lag
// can surface as transient failures).
func pollUntil(ctx context.Context, t *testing.T, what string, cond func() (bool, error)) bool {
	t.Helper()
	for attempt := range e2ePollAttempts {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				t.Fatalf("context ended while waiting for %s: %v", what, ctx.Err())
			case <-time.After(e2ePollInterval):
			}
		}
		ok, err := cond()
		if err != nil {
			t.Logf("waiting for %s (attempt %d): %v", what, attempt+1, err)
			continue
		}
		if ok {
			return true
		}
		t.Logf("waiting for %s (attempt %d)", what, attempt+1)
	}
	return false
}
