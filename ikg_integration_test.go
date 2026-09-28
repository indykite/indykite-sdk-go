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
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Policies, data and decisions propagate asynchronously; poll before failing.
const (
	e2ePollTimeout  = 60 * time.Second
	e2ePollInterval = 5 * time.Second
)

// The IKG end-to-end specs exercise an ACTUAL IKG in the project with no
// pre-provisioned fixtures — only credentials and PROJECT_ID are required:
//
//	seed    control plane: KBAC policy + CIQ read policy + knowledge query
//	ingest  runtime plane: Person -[:OWNS]-> Server into the IKG
//	assert  AuthZEN decision over the ingested graph (positive + negative)
//	assert  AuthZEN policies list shows the seeded policy (if the agent may read it)
//	assert  ContX IQ knowledge query returns the ingested node
//	cleanup graph data and config resources (always, via DeferCleanup)
var _ = Describe("IKG end to end", Ordered, Label("integration"), func() {
	const action = "SDK_IT_CAN_USE"

	var (
		cli              *indykite.Client
		kq               *config.WriteResult
		person, stranger capture.Node
		server           capture.Node
	)

	BeforeAll(func(ctx SpecContext) {
		if os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE] not set")
		}
		if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
		}
		projectID := os.Getenv("PROJECT_ID")
		if projectID == "" {
			Skip("PROJECT_ID not set")
		}

		var opts []indykite.Option
		if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
			opts = append(opts, indykite.WithBaseURL(base))
		}
		admin, err := indykite.NewAdminFromEnv(ctx, opts...)
		Expect(err).NotTo(HaveOccurred(), "NewAdminFromEnv")
		cli, err = indykite.NewClientFromEnv(ctx, opts...)
		Expect(err).NotTo(HaveOccurred(), "NewClientFromEnv")

		unique := strconv.FormatInt(time.Now().UnixNano(), 10)

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
		Expect(err).NotTo(HaveOccurred(), "create KBAC policy")
		DeferCleanup(func() { _ = admin.AuthorizationPolicies().Delete(context.Background(), kbac.ID, "") })

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
		Expect(err).NotTo(HaveOccurred(), "create CIQ policy")
		DeferCleanup(func() { _ = admin.AuthorizationPolicies().Delete(context.Background(), ciqPol.ID, "") })

		kq, err = admin.KnowledgeQueries().Create(ctx, &config.CreateKnowledgeQuery{
			ProjectID: projectID,
			Name:      "sdk-it-e2e-kq-" + unique,
			Query:     `{"nodes": ["server"]}`,
			Status:    config.StatusActive,
			PolicyID:  ciqPol.ID,
		})
		Expect(err).NotTo(HaveOccurred(), "create knowledge query")
		// Registered after the policy cleanup, so it runs first (LIFO): the
		// query references the CIQ policy.
		DeferCleanup(func() { _ = admin.KnowledgeQueries().Delete(context.Background(), kq.ID, "") })

		// --- ingest: the graph the policy and query evaluate over ---
		person = capture.Node{ExternalID: "sdk-it-e2e-person-" + unique, Type: "Person"}
		stranger = capture.Node{ExternalID: "sdk-it-e2e-stranger-" + unique, Type: "Person"}
		server = capture.Node{ExternalID: "sdk-it-e2e-server-" + unique, Type: "Server"}
		owns := capture.Relationship{Type: "OWNS", Source: &person, Target: &server}

		DeferCleanup(func() {
			cctx := context.Background()
			_, _ = cli.Capture().DeleteRelationships(cctx, owns)
			_, _ = cli.Capture().DeleteNodes(cctx, person, stranger, server)
		})

		_, err = cli.Capture().UpsertNodes(ctx,
			capture.UpsertNode{Node: person, IsIdentity: true},
			capture.UpsertNode{Node: stranger, IsIdentity: true},
			// No region property: the fixture knowledge query (region=eu) must not
			// see this transient server when packages run in parallel (plain
			// `go test ./...` without the make targets' -p 1).
			capture.UpsertNode{Node: server},
		)
		Expect(err).NotTo(HaveOccurred(), "ingest nodes")
		_, err = cli.Capture().UpsertRelationships(ctx, owns)
		Expect(err).NotTo(HaveOccurred(), "ingest relationship")
	}, NodeTimeout(2*time.Minute))

	It("allows the owner over the ingested graph", func(ctx SpecContext) {
		subject := authzen.NewNode("Person", person.ExternalID)
		resource := authzen.NewNode("Server", server.ExternalID)

		Eventually(func(g Gomega) {
			ok, err := cli.AuthZEN().Allowed(ctx, subject, action, resource)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(ok).To(BeTrue())
		}).WithContext(ctx).WithTimeout(e2ePollTimeout).WithPolling(e2ePollInterval).Should(Succeed(),
			"owner %s never got %s on %s", person.ExternalID, action, server.ExternalID)
	}, SpecTimeout(2*time.Minute))

	It("denies a stranger once the owner is allowed", func(ctx SpecContext) {
		ok, err := cli.AuthZEN().Allowed(ctx,
			authzen.NewNode("Person", stranger.ExternalID), action, authzen.NewNode("Server", server.ExternalID))
		Expect(err).NotTo(HaveOccurred(), "Allowed (stranger)")
		Expect(ok).To(BeFalse(),
			"stranger %s must not get %s on %s", stranger.ExternalID, action, server.ExternalID)
	})

	It("lists the seeded KBAC policy", func(ctx SpecContext) {
		// Only checked when the App Agent holds ReadAuthZConfigs: that is a
		// per-agent grant rather than a test input.
		_, err := cli.AuthZEN().ListPolicies(ctx, authzen.WithSubjectType("Person"))
		if apiErr, ok := transport.AsAPIError(err); ok && apiErr.IsUnauthorized() {
			Skip("App Agent lacks ReadAuthZConfigs: " + err.Error())
		}

		Eventually(func(g Gomega) {
			policies, lErr := cli.AuthZEN().ListPolicies(ctx, authzen.WithSubjectType("Person"))
			g.Expect(lErr).NotTo(HaveOccurred())
			g.Expect(policies).To(ContainElement(WithTransform(
				func(p authzen.Policy) string { return string(p.Policy) }, ContainSubstring(action))))
		}).WithContext(ctx).WithTimeout(e2ePollTimeout).WithPolling(e2ePollInterval).Should(Succeed(),
			"ListPolicies never returned the seeded policy")
	}, SpecTimeout(2*time.Minute))

	It("reads the ingested server back through the knowledge query", func(ctx SpecContext) {
		Eventually(func(g Gomega) {
			rows, err := cli.CIQ().All(ctx, ciq.ExecuteRequest{
				ID:          kq.ID,
				InputParams: map[string]any{"server_external_id": server.ExternalID},
			})
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(rows).To(ContainElement(WithTransform(func(r ciq.Record) string {
				raw, _ := json.Marshal(r)
				return string(raw)
			}, ContainSubstring(server.ExternalID))))
		}).WithContext(ctx).WithTimeout(e2ePollTimeout).WithPolling(e2ePollInterval).Should(Succeed(),
			"knowledge query never returned server %s", server.ExternalID)
	}, SpecTimeout(2*time.Minute))
})
