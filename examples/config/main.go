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

// Command config demonstrates the control-plane (config management) API with a
// Service Account credential. Every mutation is ETag-guarded: Read returns the
// ETag, Update/Delete echo it as If-Match.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

const subcommands = "organization projects list-all app-lifecycle policy-lifecycle event-sink-lifecycle"

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("config", subcommands)
	}

	ctx := context.Background()
	admin, err := indykite.NewAdminFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	projectID := fs.String("project-id", os.Getenv("PROJECT_ID"), "project gid")
	name := fs.String("name", "sdk-example", "name for created resources")
	_ = fs.Parse(os.Args[2:])

	switch os.Args[1] {
	case "organization":
		org, err := admin.Organizations().ReadCurrent(ctx)
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(org)

	case "projects":
		org, err := admin.Organizations().ReadCurrent(ctx)
		if err != nil {
			exutil.Fatal(err)
		}
		projects, err := admin.Projects().List(ctx, org.ID)
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(projects)

	case "list-all":
		// One list per project-scoped resource type.
		listAll(ctx, admin, *projectID)

	case "app-lifecycle":
		// Application -> agent -> credential, then tear down in reverse.
		appLifecycle(ctx, admin, *projectID, *name)

	case "policy-lifecycle":
		policyLifecycle(ctx, admin, *projectID, *name)

	case "event-sink-lifecycle":
		eventSinkLifecycle(ctx, admin, *projectID, *name)

	default:
		exutil.Usage("config", subcommands)
	}
}

func listAll(ctx context.Context, admin *config.AdminClient, projectID string) {
	for label, list := range map[string]func() (any, error){
		"applications": func() (any, error) { return admin.Applications().List(ctx, projectID) },
		"app-agents":   func() (any, error) { return admin.AppAgents().List(ctx, projectID) },
		"app-agent-credentials": func() (any, error) {
			return admin.AppAgentCredentials().List(ctx, projectID)
		},
		"authorization-policies": func() (any, error) {
			return admin.AuthorizationPolicies().List(ctx, projectID, "")
		},
		"knowledge-queries": func() (any, error) { return admin.KnowledgeQueries().List(ctx, projectID) },
		"event-sinks":       func() (any, error) { return admin.EventSinks().List(ctx, projectID) },
		"token-introspects": func() (any, error) { return admin.TokenIntrospects().List(ctx, projectID) },
		"external-data-resolvers": func() (any, error) {
			return admin.ExternalDataResolvers().List(ctx, projectID)
		},
		"trust-score-profiles": func() (any, error) {
			return admin.TrustScoreProfiles().List(ctx, projectID)
		},
		"entity-matching-pipelines": func() (any, error) {
			return admin.EntityMatchingPipelines().List(ctx, projectID)
		},
		"mcp-servers": func() (any, error) { return admin.MCPServers().List(ctx, projectID) },
	} {
		items, err := list()
		if err != nil {
			exutil.Fatal(fmt.Errorf("list %s: %w", label, err))
		}
		fmt.Printf("== %s ==\n", label)
		exutil.Print(items)
	}
}

func appLifecycle(ctx context.Context, admin *config.AdminClient, projectID, name string) {
	app, err := admin.Applications().Create(ctx, &config.CreateApplication{
		ProjectID: projectID, Name: name + "-app",
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("created application:", app.ID)

	agent, err := admin.AppAgents().Create(ctx, &config.CreateAppAgent{
		ApplicationID:  app.ID,
		Name:           name + "-agent",
		APIPermissions: []string{config.PermissionAuthorization, config.PermissionCapture},
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("created app agent:", agent.ID)

	// The signed credential is returned exactly once — persist it now.
	cred, err := admin.AppAgentCredentials().Create(ctx, &config.CreateAppAgentCredential{
		ApplicationAgentID: agent.ID, DisplayName: name + "-cred",
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("minted credential:", cred.ID, "kid:", cred.Kid)

	displayName := name + " agent (renamed)"
	if _, err = admin.AppAgents().Update(ctx, agent.ID, "", &config.UpdateAppAgent{
		DisplayName: &displayName,
	}); err != nil {
		exutil.Fatal(err)
	}

	// Tear down in reverse order.
	for _, del := range []func() error{
		func() error { return admin.AppAgentCredentials().Delete(ctx, cred.ID) },
		func() error { return admin.AppAgents().Delete(ctx, agent.ID, "") },
		func() error { return admin.Applications().Delete(ctx, app.ID, "") },
	} {
		if err = del(); err != nil {
			exutil.Fatal(err)
		}
	}
	fmt.Println("cleaned up")
}

func policyLifecycle(ctx context.Context, admin *config.AdminClient, projectID, name string) {
	policyJSON := `{
	  "meta": {"policy_version": "2.0-kbac"},
	  "subject": {"type": "Person"},
	  "actions": ["CAN_READ"],
	  "resource": {"type": "Asset"},
	  "condition": {"cypher": "MATCH (subject:Person)-[:OWNS]->(resource:Asset)"}
	}`
	created, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: projectID, Name: name + "-policy", Policy: policyJSON, Status: config.StatusDraft,
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("created policy:", created.ID)

	pol, err := admin.AuthorizationPolicies().Read(ctx, created.ID)
	if err != nil {
		exutil.Fatal(err)
	}

	// Activate it, guarded by the ETag captured on Read.
	if _, err = admin.AuthorizationPolicies().Update(ctx, pol.ID, pol.ETag, &config.UpdateAuthorizationPolicy{
		Policy: pol.Policy, Status: config.StatusActive,
	}); err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("activated policy")

	if err = admin.AuthorizationPolicies().Delete(ctx, pol.ID, ""); err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("cleaned up")
}

func eventSinkLifecycle(ctx context.Context, admin *config.AdminClient, projectID, name string) {
	created, err := admin.EventSinks().Create(ctx, &config.CreateEventSink{
		ProjectID: projectID,
		Name:      name + "-sink",
		Providers: map[string]config.EventSinkProvider{
			"kafka-main": {Kafka: &config.KafkaSinkConfig{
				Brokers: []string{"kafka.example.com:9092"},
				Topic:   "indykite-events",
			}},
		},
		Routes: []config.EventSinkRoute{{
			ProviderID: "kafka-main",
			Filter:     config.EventSinkFilter{EventType: "indykite.audit.config.create"},
		}},
	})
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("created event sink:", created.ID)

	sink, err := admin.EventSinks().Read(ctx, created.ID)
	if err != nil {
		exutil.Fatal(err)
	}
	exutil.Print(sink)

	if err = admin.EventSinks().Delete(ctx, sink.ID, sink.ETag); err != nil {
		exutil.Fatal(err)
	}
	fmt.Println("cleaned up")
}
