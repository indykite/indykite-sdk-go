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

package authzen_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/internal/bqaudit"
	"github.com/indykite/indykite-sdk-go/transport"
)

// evaluationAuditEventType is emitted by the platform for POST /access/v1/evaluation.
const evaluationAuditEventType = "indykite.audit.authorization.evaluation"

// requireAuditEvent asserts the audit event carrying marker reached the
// BigQuery audit-log table. It is a no-op unless SDK_AUDIT_TABLE_NAME is set.
func requireAuditEvent(ctx context.Context, t *testing.T, eventType, marker string) {
	t.Helper()
	if !bqaudit.Enabled() {
		t.Log("SDK_AUDIT_TABLE_NAME not set; skipping BigQuery audit-log check")
		return
	}
	checker, err := bqaudit.New(ctx)
	if err != nil {
		t.Fatalf("bqaudit.New: %v", err)
	}
	defer func() { _ = checker.Close() }()
	if err := checker.WaitForEvent(ctx, eventType, marker); err != nil {
		t.Errorf("audit log not found in BigQuery: %v", err)
	}
}

// fixture is the (subject, action, resource) triple the environment provides
// for decision tests.
type fixture struct {
	subject  authzen.Node
	action   string
	resource authzen.Node
}

func fixtures(t *testing.T) fixture {
	t.Helper()
	f := fixture{
		subject:  authzen.NewNode(os.Getenv("AUTHZEN_SUBJECT_TYPE"), os.Getenv("AUTHZEN_SUBJECT_ID")),
		action:   os.Getenv("AUTHZEN_ACTION"),
		resource: authzen.NewNode(os.Getenv("AUTHZEN_RESOURCE_TYPE"), os.Getenv("AUTHZEN_RESOURCE_ID")),
	}
	if f.subject.Type == "" || f.subject.ID == "" || f.action == "" ||
		f.resource.Type == "" || f.resource.ID == "" {
		t.Skip("AUTHZEN_{SUBJECT_TYPE,SUBJECT_ID,ACTION,RESOURCE_TYPE,RESOURCE_ID} not set")
	}
	return f
}

func runtimeClient(t *testing.T) *indykite.Client {
	t.Helper()
	if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
		t.Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
	}
	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	cli, err := indykite.NewClientFromEnv(context.Background(), opts...)
	if err != nil {
		t.Fatalf("NewClientFromEnv: %v", err)
	}
	return cli
}

func TestIntegrationAuthZENEvaluate(t *testing.T) {
	cli := runtimeClient(t)
	f := fixtures(t)
	ctx := context.Background()

	// The auditLog input param is echoed into the audit event, which lets the
	// BigQuery check below correlate this exact request (gRPC SDK convention).
	auditMarker := fmt.Sprintf("sdk-it-authzen-%d", time.Now().UnixNano())
	resp, err := cli.AuthZEN().Evaluate(ctx, authzen.EvaluationRequest{
		Subject:  &f.subject,
		Resource: &f.resource,
		Action:   &authzen.Action{Name: f.action},
		Context:  &authzen.Context{InputParams: map[string]any{"auditLog": auditMarker}},
	})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	t.Logf("decision=%v", resp.Decision)

	allowed, err := cli.AuthZEN().Allowed(ctx, f.subject, f.action, f.resource)
	if err != nil {
		t.Fatalf("Allowed: %v", err)
	}
	if allowed != resp.Decision {
		t.Errorf("Allowed=%v but Evaluate decision=%v", allowed, resp.Decision)
	}

	requireAuditEvent(ctx, t, evaluationAuditEventType, auditMarker)
}

func TestIntegrationAuthZENEvaluateBatch(t *testing.T) {
	cli := runtimeClient(t)
	f := fixtures(t)

	resp, err := cli.AuthZEN().EvaluateBatch(context.Background(), authzen.EvaluationsRequest{
		Subject: &f.subject,
		Action:  &authzen.Action{Name: f.action},
		Evaluations: []authzen.EvaluationItem{
			{Resource: &f.resource},
			{Resource: &authzen.Node{Type: f.resource.Type, ID: "nonexistent-" + f.resource.ID}},
		},
	})
	if err != nil {
		t.Fatalf("EvaluateBatch: %v", err)
	}
	if len(resp.Evaluations) != 2 {
		t.Fatalf("got %d evaluations, want 2", len(resp.Evaluations))
	}
}

func TestIntegrationAuthZENSearch(t *testing.T) {
	cli := runtimeClient(t)
	f := fixtures(t)
	ctx := context.Background()

	actions, err := cli.AuthZEN().SearchAction(ctx, authzen.SearchActionRequest{
		Subject: &f.subject, Resource: &f.resource,
	})
	if err != nil {
		t.Fatalf("SearchAction: %v", err)
	}
	t.Logf("actions=%v", actions)

	resources, err := cli.AuthZEN().SearchResource(ctx, authzen.SearchResourceRequest{
		Subject:  &f.subject,
		Action:   &authzen.Action{Name: f.action},
		Resource: &authzen.NodeType{Type: f.resource.Type},
	})
	if err != nil {
		t.Fatalf("SearchResource: %v", err)
	}
	t.Logf("resources=%v", resources)

	subjects, err := cli.AuthZEN().SearchSubject(ctx, authzen.SearchSubjectRequest{
		Subject:  &authzen.NodeType{Type: f.subject.Type},
		Action:   &authzen.Action{Name: f.action},
		Resource: &f.resource,
	})
	if err != nil {
		t.Fatalf("SearchSubject: %v", err)
	}
	t.Logf("subjects=%v", subjects)
}

// TestIntegrationAuthZENUnknownResourceType asserts platform errors surface as
// *transport.APIError with useful fields.
func TestIntegrationAuthZENErrorShape(t *testing.T) {
	cli := runtimeClient(t)
	f := fixtures(t)

	_, err := cli.AuthZEN().Evaluate(context.Background(), authzen.EvaluationRequest{
		Subject: &f.subject, // missing action & resource
	})
	if err == nil {
		t.Skip("platform accepted an incomplete request")
	}
	apiErr, ok := transport.AsAPIError(err)
	if !ok {
		t.Fatalf("error is %T, want *transport.APIError: %v", err, err)
	}
	if apiErr.StatusCode < 400 || apiErr.StatusCode > 499 {
		t.Errorf("StatusCode = %d, want 4xx", apiErr.StatusCode)
	}
}
