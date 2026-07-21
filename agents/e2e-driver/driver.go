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

// Command e2e-driver is an SDK-backed end-to-end driver. It is the Go,
// REST-native counterpart to the agent-gateway-e2e Bruno collection's
// platform-side setup and verification:
//
//	seed   -> create a KBAC authorization policy   (control plane, config.AdminClient)
//	ingest -> upsert the subject/resource graph     (runtime plane, capture)
//	assert -> AuthZEN decision + ContX IQ read       (runtime plane, authzen + ciq)
//	cleanup-> delete the policy                      (control plane)
//
// It drives the platform directly over REST; it does not exercise the gateway's
// agent-to-agent (A2A) hops, which are a separate axis the gateway owns.
package main

import (
	"context"
	"fmt"

	"github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
)

// Scenario describes one end-to-end check.
type Scenario struct {
	Relationship *capture.Relationship // optional edge to ingest
	InputParams  map[string]any        // input params for the decision
	QueryParams  map[string]any        // input params for the CIQ read
	Policy       config.CreateAuthorizationPolicy
	Subject      authzen.Node
	Resource     authzen.Node
	Action       string
	QueryID      string               // optional ContX IQ query to read
	Nodes        []capture.UpsertNode // subject/resource graph to ingest
	WantAllow    bool
}

// Step is the outcome of one driver step.
type Step struct {
	Name   string
	Detail string
	OK     bool
}

// Driver runs scenarios using both planes: AdminClient (control) for config and
// the runtime Client for ingest/decisions.
type Driver struct {
	admin *config.AdminClient
	ik    *indykite.Client
}

// NewDriver builds a driver from a control-plane admin client and a runtime client.
func NewDriver(admin *config.AdminClient, ik *indykite.Client) *Driver {
	return &Driver{admin: admin, ik: ik}
}

// Run executes seed -> ingest -> assert and always attempts cleanup. It returns
// one Step per stage; err is non-nil only on a fatal failure that aborts the run.
// policyRef identifies a created policy for cleanup.
type policyRef struct {
	ID   string
	ETag string
}

func (d *Driver) Run(ctx context.Context, sc *Scenario) ([]Step, error) {
	steps := make([]Step, 0, 6)
	ref, err := d.execute(ctx, sc, &steps)

	// Always attempt cleanup of the policy we created, even on a failed assertion.
	if ref.ID != "" {
		if delErr := d.admin.AuthorizationPolicies().Delete(ctx, ref.ID, ref.ETag); delErr != nil {
			steps = append(steps, Step{Name: "cleanup-policy", Detail: delErr.Error()})
		} else {
			steps = append(steps, Step{Name: "cleanup-policy", OK: true})
		}
	}
	return steps, err
}

// execute runs seed -> ingest -> assert, appending one Step per stage to *steps.
// It returns the created policy reference (for cleanup by Run) and a fatal error.
func (d *Driver) execute(ctx context.Context, sc *Scenario, steps *[]Step) (policyRef, error) {
	// 1. Seed: create the authorization policy (control plane).
	pol, err := d.admin.AuthorizationPolicies().Create(ctx, &sc.Policy)
	if err != nil {
		*steps = append(*steps, Step{Name: "seed-policy", Detail: err.Error()})
		return policyRef{}, fmt.Errorf("seed policy: %w", err)
	}
	ref := policyRef{ID: pol.ID, ETag: pol.ETag}
	*steps = append(*steps, Step{Name: "seed-policy", OK: true, Detail: pol.ID})

	// 2. Ingest: the subject/resource graph the policy evaluates over (runtime).
	if len(sc.Nodes) > 0 {
		if _, err = d.ik.Capture().UpsertNodes(ctx, sc.Nodes...); err != nil {
			*steps = append(*steps, Step{Name: "ingest-nodes", Detail: err.Error()})
			return ref, fmt.Errorf("ingest nodes: %w", err)
		}
		*steps = append(*steps, Step{Name: "ingest-nodes", OK: true, Detail: fmt.Sprintf("%d nodes", len(sc.Nodes))})
	}
	if sc.Relationship != nil {
		if _, err = d.ik.Capture().UpsertRelationships(ctx, *sc.Relationship); err != nil {
			*steps = append(*steps, Step{Name: "ingest-relationship", Detail: err.Error()})
			return ref, fmt.Errorf("ingest relationship: %w", err)
		}
		*steps = append(*steps, Step{Name: "ingest-relationship", OK: true})
	}

	// 3a. Assert: the AuthZEN decision matches expectation.
	var opts []authzen.Option
	if len(sc.InputParams) > 0 {
		opts = append(opts, authzen.WithInputParams(sc.InputParams))
	}
	allowed, err := d.ik.AuthZEN().Allowed(ctx, sc.Subject, sc.Action, sc.Resource, opts...)
	if err != nil {
		*steps = append(*steps, Step{Name: "assert-decision", Detail: err.Error()})
		return ref, fmt.Errorf("assert decision: %w", err)
	}
	*steps = append(*steps, Step{
		Name:   "assert-decision",
		OK:     allowed == sc.WantAllow,
		Detail: fmt.Sprintf("allowed=%v want=%v", allowed, sc.WantAllow),
	})

	// 3b. Assert: the optional ContX IQ query returns rows.
	if sc.QueryID != "" {
		rows, qErr := d.ik.CIQ().All(ctx, ciq.ExecuteRequest{ID: sc.QueryID, InputParams: sc.QueryParams})
		if qErr != nil {
			*steps = append(*steps, Step{Name: "assert-query", Detail: qErr.Error()})
			return ref, fmt.Errorf("assert query: %w", qErr)
		}
		*steps = append(*steps, Step{
			Name:   "assert-query",
			OK:     len(rows) > 0,
			Detail: fmt.Sprintf("%d rows", len(rows)),
		})
	}

	return ref, nil
}

// AllOK reports whether every step succeeded.
func AllOK(steps []Step) bool {
	for _, s := range steps {
		if !s.OK {
			return false
		}
	}
	return true
}
