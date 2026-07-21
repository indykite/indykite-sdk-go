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

package indykite_test

import (
	"context"
	"fmt"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
)

// Build the runtime-plane client from INDYKITE_APPLICATION_CREDENTIALS[_FILE]
// and make one authorization decision.
func ExampleNewClientFromEnv() {
	ctx := context.Background()

	cli, err := indykite.NewClientFromEnv(ctx, indykite.WithRegion("eu"))
	if err != nil {
		return
	}

	allowed, err := cli.AuthZEN().Allowed(ctx,
		authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
	if err != nil {
		return
	}
	fmt.Println("decision:", allowed)
}

// Run a ContX IQ knowledge query and collect every page.
func ExampleClient_CIQ() {
	ctx := context.Background()

	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	records, err := cli.CIQ().All(ctx, ciq.ExecuteRequest{
		ID:          "get-servers",
		InputParams: map[string]any{"region": "eu"},
	})
	if err != nil {
		return
	}
	fmt.Println("records:", len(records))
}

// Build the control-plane client from INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE]
// and walk the ETag-guarded lifecycle of an authorization policy.
func ExampleNewAdminFromEnv() {
	ctx := context.Background()

	admin, err := indykite.NewAdminFromEnv(ctx, indykite.WithRegion("eu"))
	if err != nil {
		return
	}

	created, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
		ProjectID: "gid:project",
		Name:      "can-provision",
		Policy: `{"meta":{"policy_version":"2.0-kbac"},"subject":{"type":"Person"},` +
			`"actions":["PROVISION"],"resource":{"type":"Server"},` +
			`"condition":{"cypher":"MATCH (subject:Person)-[:CAN_USE]->(resource:Server)"}}`,
		Status: config.StatusActive,
	})
	if err != nil {
		return
	}

	// Read captures the ETag; Update echoes it as If-Match.
	pol, err := admin.AuthorizationPolicies().Read(ctx, created.ID)
	if err != nil {
		return
	}
	if _, err = admin.AuthorizationPolicies().Update(ctx, pol.ID, pol.ETag, &config.UpdateAuthorizationPolicy{
		Policy: pol.Policy,
		Status: config.StatusInactive,
	}); err != nil {
		return
	}
}
