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
	"encoding/json"
	"fmt"
	"os"
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/internal/bqaudit"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// evaluationAuditEventType is emitted by the platform for POST /access/v1/evaluation.
const evaluationAuditEventType = "indykite.audit.authorization.evaluation"

// fixture is the (subject, action, resource) triple the environment provides
// for decision specs.
type fixture struct {
	subject  authzen.Node
	action   string
	resource authzen.Node
}

// fixtures reads the decision triple from the environment, skipping the spec
// when any part of it is missing.
func fixtures() fixture {
	GinkgoHelper()
	f := fixture{
		subject:  authzen.NewNode(os.Getenv("AUTHZEN_SUBJECT_TYPE"), os.Getenv("AUTHZEN_SUBJECT_ID")),
		action:   os.Getenv("AUTHZEN_ACTION"),
		resource: authzen.NewNode(os.Getenv("AUTHZEN_RESOURCE_TYPE"), os.Getenv("AUTHZEN_RESOURCE_ID")),
	}
	if f.subject.Type == "" || f.subject.ID == "" || f.action == "" ||
		f.resource.Type == "" || f.resource.ID == "" {
		Skip("AUTHZEN_{SUBJECT_TYPE,SUBJECT_ID,ACTION,RESOURCE_TYPE,RESOURCE_ID} not set")
	}
	return f
}

// authzenClient builds a live runtime-plane client from the environment,
// skipping the spec when no App Agent credential is configured.
func authzenClient(ctx context.Context) *authzen.Client {
	GinkgoHelper()
	if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
		Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
	}
	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	cli, err := indykite.NewClientFromEnv(ctx, opts...)
	Expect(err).NotTo(HaveOccurred())
	return cli.AuthZEN()
}

// expectAuditEvent asserts the audit event carrying marker reached the
// BigQuery audit-log table. It is a no-op unless SDK_AUDIT_TABLE_NAME is set.
func expectAuditEvent(ctx context.Context, eventType, marker string) {
	GinkgoHelper()
	if !bqaudit.Enabled() {
		GinkgoWriter.Println("SDK_AUDIT_TABLE_NAME not set; skipping BigQuery audit-log check")
		return
	}
	checker, err := bqaudit.New(ctx)
	Expect(err).NotTo(HaveOccurred())
	DeferCleanup(func() { _ = checker.Close() })
	Expect(checker.WaitForEvent(ctx, eventType, marker)).To(Succeed(), "audit log not found in BigQuery")
}

var _ = Describe("Integration", Label("integration"), func() {
	var c *authzen.Client

	BeforeEach(func(ctx SpecContext) {
		c = authzenClient(ctx)
	})

	Describe("Evaluate", func() {
		It("decides over the fixture triple, agrees with Allowed, and lands in the audit log", func(ctx SpecContext) {
			f := fixtures()

			// The auditLog input param is echoed into the audit event, which lets
			// the BigQuery check below correlate this exact request.
			auditMarker := fmt.Sprintf("sdk-it-authzen-%d", time.Now().UnixNano())
			resp, err := c.Evaluate(ctx, authzen.EvaluationRequest{
				Subject:  &f.subject,
				Resource: &f.resource,
				Action:   &authzen.Action{Name: f.action},
				Context:  &authzen.Context{InputParams: map[string]any{"auditLog": auditMarker}},
			})
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("decision=%v\n", resp.Decision)

			allowed, err := c.Allowed(ctx, f.subject, f.action, f.resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(allowed).To(Equal(resp.Decision), "Allowed and Evaluate must agree")

			expectAuditEvent(ctx, evaluationAuditEventType, auditMarker)
		})

		// Platform errors surface as *transport.APIError with useful fields.
		It("surfaces an incomplete request as a 4xx APIError", func(ctx SpecContext) {
			f := fixtures()

			_, err := c.Evaluate(ctx, authzen.EvaluationRequest{
				Subject: &f.subject, // missing action & resource
			})
			if err == nil {
				Skip("platform accepted an incomplete request")
			}
			apiErr, ok := transport.AsAPIError(err)
			Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
			Expect(apiErr.StatusCode).To(And(BeNumerically(">=", 400), BeNumerically("<", 500)))
		})
	})

	Describe("EvaluateBatch", func() {
		It("returns one decision per entry", func(ctx SpecContext) {
			f := fixtures()

			resp, err := c.EvaluateBatch(ctx, authzen.EvaluationsRequest{
				Subject: &f.subject,
				Action:  &authzen.Action{Name: f.action},
				Evaluations: []authzen.EvaluationItem{
					{Resource: &f.resource},
					{Resource: &authzen.Node{Type: f.resource.Type, ID: "nonexistent-" + f.resource.ID}},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Evaluations).To(HaveLen(2))
		})
	})

	Describe("Search", func() {
		It("enumerates actions, resources and subjects around the fixture triple", func(ctx SpecContext) {
			f := fixtures()

			actions, err := c.SearchAction(ctx, authzen.SearchActionRequest{
				Subject: &f.subject, Resource: &f.resource,
			})
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("actions=%v\n", actions)

			resources, err := c.SearchResource(ctx, authzen.SearchResourceRequest{
				Subject:  &f.subject,
				Action:   &authzen.Action{Name: f.action},
				Resource: &authzen.NodeType{Type: f.resource.Type},
			})
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("resources=%v\n", resources)

			subjects, err := c.SearchSubject(ctx, authzen.SearchSubjectRequest{
				Subject:  &authzen.NodeType{Type: f.subject.Type},
				Action:   &authzen.Action{Name: f.action},
				Resource: &f.resource,
			})
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("subjects=%v\n", subjects)
		})
	})

	Describe("ListPolicies", func() {
		// listAll reads every active policy, skipping the spec when the App
		// Agent lacks the ReadAuthZConfigs permission: that is a per-agent
		// grant rather than a test input.
		listAll := func(ctx context.Context) []authzen.Policy {
			GinkgoHelper()
			all, err := c.ListPolicies(ctx)
			if apiErr, ok := transport.AsAPIError(err); ok && apiErr.IsUnauthorized() {
				Skip(fmt.Sprintf("App Agent lacks the ReadAuthZConfigs permission: %v", err))
			}
			Expect(err).NotTo(HaveOccurred())
			return all
		}

		It("returns every policy as a JSON object with non-null tags", func(ctx SpecContext) {
			all := listAll(ctx)
			GinkgoWriter.Printf("policies=%d\n", len(all))
			for i, p := range all {
				Expect(json.Valid(p.Policy)).To(BeTrue(), "policy[%d] is not valid JSON: %s", i, p.Policy)
				Expect(string(p.Policy)).To(HavePrefix("{"), "policy[%d] is not a JSON object", i)
				Expect(p.Tags).NotTo(BeNil(), "policy[%d] tags must be a (possibly empty) slice", i)
			}
		})

		It("narrows to a subset when filtered by the fixture subject type", func(ctx SpecContext) {
			all := listAll(ctx)
			f := fixtures()

			subset, err := c.ListPolicies(ctx, authzen.WithSubjectType(f.subject.Type))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(subset)).To(BeNumerically("<=", len(all)),
				"subject_type=%s must not return more policies than the unfiltered list", f.subject.Type)
		})
	})
})
