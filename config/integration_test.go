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

package config_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"
)

// NOTE: there is deliberately no integration test for the data-schema API
// (config.DataSchemaAPI) — it is still behind a feature flag.

func adminClient(t *testing.T) *config.AdminClient {
	t.Helper()
	if os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS_FILE") == "" {
		t.Skip("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE] not set")
	}
	a, err := auth.ServiceAccountFromEnv(context.Background())
	if err != nil {
		t.Fatalf("ServiceAccountFromEnv: %v", err)
	}
	var opts []transport.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, transport.WithBaseURL(base))
	}
	tc, err := transport.NewClient(a, opts...)
	if err != nil {
		t.Fatalf("transport.NewClient: %v", err)
	}
	return config.NewAdminClient(tc)
}

func projectID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("PROJECT_ID")
	if id == "" {
		t.Skip("PROJECT_ID not set")
	}
	return id
}

func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func TestIntegrationConfigOrganizationAndProjects(t *testing.T) {
	admin := adminClient(t)
	ctx := context.Background()

	org, err := admin.Organizations().ReadCurrent(ctx)
	if err != nil {
		t.Fatalf("Organizations.ReadCurrent: %v", err)
	}
	if org.ID == "" {
		t.Fatalf("organization has no id: %+v", org)
	}

	projects, err := admin.Projects().List(ctx, org.ID)
	if err != nil {
		t.Fatalf("Projects.List: %v", err)
	}
	t.Logf("organization %s has %d projects", org.ID, len(projects))
}

// TestIntegrationConfigAuthorizationPolicyCRUD runs the full ETag-guarded
// lifecycle: create draft -> read -> activate -> delete.
func TestIntegrationConfigAuthorizationPolicyCRUD(t *testing.T) {
	admin := adminClient(t)
	ctx := context.Background()
	project := projectID(t)
	api := admin.AuthorizationPolicies()

	policyJSON := `{
	  "meta": {"policy_version": "2.0-kbac"},
	  "subject": {"type": "Person"},
	  "actions": ["SDK_IT_CAN_READ"],
	  "resource": {"type": "Asset"},
	  "condition": {"cypher": "MATCH (subject:Person)-[:OWNS]->(resource:Asset)"}
	}`
	created, err := api.Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: project,
		Name:      uniqueName("sdk-it-policy"),
		Policy:    policyJSON,
		Status:    config.StatusDraft,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	t.Cleanup(func() { _ = api.Delete(ctx, created.ID, "") })

	pol, err := api.Read(ctx, created.ID)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if pol.Status != config.StatusDraft {
		t.Errorf("status = %q, want DRAFT", pol.Status)
	}
	if pol.ETag == "" {
		t.Error("Read returned no ETag")
	}

	if _, err = api.Update(ctx, pol.ID, pol.ETag, &config.UpdateAuthorizationPolicy{
		Policy: pol.Policy,
		Status: config.StatusActive,
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := api.Read(ctx, pol.ID)
	if err != nil {
		t.Fatalf("Read after update: %v", err)
	}
	if got.Status != config.StatusActive {
		t.Errorf("status after update = %q, want ACTIVE", got.Status)
	}

	if err = api.Delete(ctx, got.ID, got.ETag); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err = api.Read(ctx, got.ID); err == nil {
		t.Error("Read after delete succeeded, want not-found error")
	} else if apiErr, ok := transport.AsAPIError(err); !ok || !apiErr.IsNotFound() {
		t.Errorf("Read after delete: got %v, want 404 APIError", err)
	}
}

// TestIntegrationConfigAppAgentLifecycle chains application -> agent ->
// credential (the runtime-plane bootstrap) and tears it down in reverse.
func TestIntegrationConfigAppAgentLifecycle(t *testing.T) {
	admin := adminClient(t)
	ctx := context.Background()
	project := projectID(t)

	app, err := admin.Applications().Create(ctx, &config.CreateApplication{
		ProjectID: project, Name: uniqueName("sdk-it-app"),
	})
	if err != nil {
		t.Fatalf("Applications.Create: %v", err)
	}
	t.Cleanup(func() { _ = admin.Applications().Delete(ctx, app.ID, "") })

	agent, err := admin.AppAgents().Create(ctx, &config.CreateAppAgent{
		ApplicationID:  app.ID,
		Name:           uniqueName("sdk-it-agent"),
		APIPermissions: []string{config.PermissionAuthorization},
	})
	if err != nil {
		t.Fatalf("AppAgents.Create: %v", err)
	}
	t.Cleanup(func() { _ = admin.AppAgents().Delete(ctx, agent.ID, "") })

	displayName := "sdk-it renamed"
	if _, err = admin.AppAgents().Update(ctx, agent.ID, "", &config.UpdateAppAgent{
		DisplayName: &displayName,
	}); err != nil {
		t.Fatalf("AppAgents.Update: %v", err)
	}

	cred, err := admin.AppAgentCredentials().Create(ctx, &config.CreateAppAgentCredential{
		ApplicationAgentID: agent.ID, DisplayName: "sdk-it-cred",
	})
	if err != nil {
		t.Fatalf("AppAgentCredentials.Create: %v", err)
	}
	t.Cleanup(func() { _ = admin.AppAgentCredentials().Delete(ctx, cred.ID) })
	if len(cred.ApplicationAgentConfig) == 0 {
		t.Error("credential Create returned no signed config")
	}

	creds, err := admin.AppAgentCredentials().List(ctx, project)
	if err != nil {
		t.Fatalf("AppAgentCredentials.List: %v", err)
	}
	found := false
	for _, c := range creds {
		if c.ID == cred.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("minted credential %s not in list of %d", cred.ID, len(creds))
	}
}

// TestIntegrationConfigLists exercises every remaining project-scoped List.
func TestIntegrationConfigLists(t *testing.T) {
	admin := adminClient(t)
	ctx := context.Background()
	project := projectID(t)

	for label, list := range map[string]func() (int, error){
		"applications": func() (int, error) {
			items, err := admin.Applications().List(ctx, project)
			return len(items), err
		},
		"app-agents": func() (int, error) {
			items, err := admin.AppAgents().List(ctx, project)
			return len(items), err
		},
		"authorization-policies": func() (int, error) {
			items, err := admin.AuthorizationPolicies().List(ctx, project, "")
			return len(items), err
		},
		"knowledge-queries": func() (int, error) {
			items, err := admin.KnowledgeQueries().List(ctx, project)
			return len(items), err
		},
		"event-sinks": func() (int, error) {
			items, err := admin.EventSinks().List(ctx, project)
			return len(items), err
		},
		"token-introspects": func() (int, error) {
			items, err := admin.TokenIntrospects().List(ctx, project)
			return len(items), err
		},
		"external-data-resolvers": func() (int, error) {
			items, err := admin.ExternalDataResolvers().List(ctx, project)
			return len(items), err
		},
		"trust-score-profiles": func() (int, error) {
			items, err := admin.TrustScoreProfiles().List(ctx, project)
			return len(items), err
		},
		"entity-matching-pipelines": func() (int, error) {
			items, err := admin.EntityMatchingPipelines().List(ctx, project)
			return len(items), err
		},
		"mcp-servers": func() (int, error) {
			items, err := admin.MCPServers().List(ctx, project)
			return len(items), err
		},
	} {
		n, err := list()
		if err != nil {
			t.Errorf("list %s: %v", label, err)
			continue
		}
		t.Logf("%s: %d", label, n)
	}

	if org := os.Getenv("ORGANIZATION_ID"); org != "" {
		n, err := admin.ServiceAccountCredentials().List(ctx, org)
		if err != nil {
			t.Errorf("list service-account-credentials: %v", err)
		} else {
			t.Logf("service-account-credentials: %d", len(n))
		}
	}
}
