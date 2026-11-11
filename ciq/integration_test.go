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

package ciq_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwk"
	"github.com/lestrrat-go/jwx/v4/jwt"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/internal/bqaudit"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// executeAuditEventType is emitted by the platform for POST /contx-iq/v1/execute.
const executeAuditEventType = "indykite.audit.ciq.execute"

// queryFixture returns the pre-configured knowledge query (and optional input
// params) the environment provides, skipping the spec when not configured.
func queryFixture() ciq.ExecuteRequest {
	GinkgoHelper()
	id := os.Getenv("CIQ_QUERY_ID")
	if id == "" {
		Skip("CIQ_QUERY_ID not set")
	}
	req := ciq.ExecuteRequest{ID: id}
	if raw := os.Getenv("CIQ_INPUT_PARAMS"); raw != "" {
		Expect(json.Unmarshal([]byte(raw), &req.InputParams)).To(Succeed(), "CIQ_INPUT_PARAMS is not a JSON object")
	}
	return req
}

// ciqClient builds a live runtime-plane client from the environment, skipping
// the spec when no App Agent credential is configured.
func ciqClient(ctx context.Context) *ciq.Client {
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
	return cli.CIQ()
}

var _ = Describe("Integration", Label("integration"), func() {
	var c *ciq.Client

	BeforeEach(func(ctx SpecContext) {
		c = ciqClient(ctx)
	})

	Describe("Execute", func() {
		It("runs the fixture query and lands in the audit log", func(ctx SpecContext) {
			req := queryFixture()

			// The auditLog input param is echoed into the audit event, which lets
			// the BigQuery check below correlate this exact request.
			auditMarker := fmt.Sprintf("sdk-it-ciq-%d", time.Now().UnixNano())
			if req.InputParams == nil {
				req.InputParams = map[string]any{}
			}
			req.InputParams["auditLog"] = auditMarker

			resp, err := c.Execute(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("first page: %d records\n", len(resp.Data))

			if !bqaudit.Enabled() {
				GinkgoWriter.Println("SDK_AUDIT_TABLE_NAME not set; skipping BigQuery audit-log check")
				return
			}
			checker, err := bqaudit.New(ctx)
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() { _ = checker.Close() })
			Expect(checker.WaitForEvent(ctx, executeAuditEventType, auditMarker)).
				To(Succeed(), "audit log not found in BigQuery")
		})
	})

	Describe("All and Iterate", func() {
		// A small page size makes sure pagination is actually exercised whenever
		// the fixture query returns more than two records.
		It("agree on the full result set", func(ctx SpecContext) {
			req := queryFixture()
			req.PageSize = 2

			records, err := c.All(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("total: %d records\n", len(records))

			it := c.Iterate(req)
			count := 0
			for it.Next(ctx) {
				count++
			}
			Expect(it.Err()).NotTo(HaveOccurred())
			Expect(count).To(Equal(len(records)))
		})
	})

	Describe("WhoAmI", func() {
		// Only needs the App Agent credential: the platform must reject a token
		// it cannot introspect.
		It("rejects a bogus end-user token", func(ctx SpecContext) {
			_, err := c.WhoAmI(ctx, "not-a-valid-end-user-token")
			Expect(err).To(HaveOccurred())
			GinkgoWriter.Printf("rejected as expected: %v\n", err)
		})
	})
})

// WhoAmI end-to-end needs no pre-minted end-user token (those live ~2h): the
// spec generates a signing key, publishes its public JWK in a Token Introspect
// config it creates in the project, ingests the Person the token will name,
// and signs a fresh short-lived JWT itself. Needs both credentials + PROJECT_ID.
var _ = Describe("Integration WhoAmI end-to-end", Label("integration"), func() {
	const (
		propagationTimeout = 90 * time.Second
		propagationPoll    = 5 * time.Second
		tokenLifetime      = 5 * time.Minute
	)

	var (
		admin  *config.AdminClient
		cli    *ciq.Client
		runtim *indykite.Client
	)

	BeforeEach(func(ctx SpecContext) {
		cli = ciqClient(ctx)
		admin, runtim = adminAndRuntimeClients(ctx)
	})

	It("resolves a self-signed end-user token to the ingested Person", func(ctx SpecContext) {
		projectID := os.Getenv("PROJECT_ID")
		unique := strconv.FormatInt(time.Now().UnixNano(), 10)
		// A unique issuer per run keeps the platform's config cache from ever
		// matching a previous run's (deleted) config.
		issuer := "https://sdk-it.indykite.com/whoami/" + unique
		audience := "sdk-it-whoami"
		subject := "sdk-it-whoami-person-" + unique

		signer := newSigningKey(unique)

		// --- seed: Token Introspect config trusting our public key ---
		created, err := admin.TokenIntrospects().Create(ctx, &config.CreateTokenIntrospect{
			ProjectID: projectID,
			Name:      "sdk-it-whoami-" + unique,
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType:       "Person",
				JwtMatcher:        mustJSON(map[string]string{"issuer": issuer, "audience": audience}),
				OfflineValidation: mustJSON(map[string][]string{"public_jwks": {signer.publicJWK}}),
				PerformUpsert:     false,
			},
		})
		Expect(err).NotTo(HaveOccurred(), "create token introspect config")
		DeferCleanup(func(ctx SpecContext) {
			_ = admin.TokenIntrospects().Delete(ctx, created.ID, "")
		})

		// --- ingest: the Person the token's sub must resolve to ---
		person := capture.Node{ExternalID: subject, Type: "Person"}
		_, err = runtim.Capture().UpsertNodes(ctx, capture.UpsertNode{Node: person, IsIdentity: true})
		Expect(err).NotTo(HaveOccurred(), "ingest person")
		DeferCleanup(func(ctx SpecContext) {
			_, _ = runtim.Capture().DeleteNodes(ctx, person)
		})

		// --- act: sign a fresh token and poll until config + node propagate ---
		token := signEndUserToken(signer.private, issuer, audience, subject, tokenLifetime)
		Eventually(func() (*ciq.WhoAmIResponse, error) {
			return cli.WhoAmI(ctx, token)
		}).WithContext(ctx).WithTimeout(propagationTimeout).WithPolling(propagationPoll).
			Should(HaveValue(Equal(ciq.WhoAmIResponse{Type: "Person", ID: subject})))

		// --- negative: same iss/aud but signed by a key the config does not trust ---
		impostor := newSigningKey(unique + "-impostor")
		forged := signEndUserToken(impostor.private, issuer, audience, subject, tokenLifetime)
		_, err = cli.WhoAmI(ctx, forged)
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
		Expect(apiErr.IsUnauthorized()).To(BeTrue(), "forged token must be rejected: %v", apiErr)
	})
})

// A 2.0-ciq write query whose upsert values are "@" references: the platform
// reads them from the graph the policy matched and writes them on the matched
// node, so the response carries values the request never sent. Needs both
// credentials + PROJECT_ID.
var _ = Describe("Integration 2.0-ciq reference values", Label("integration"), func() {
	const (
		pollTimeout  = 60 * time.Second
		pollInterval = 5 * time.Second
	)

	var (
		admin  *config.AdminClient
		cli    *ciq.Client
		runtim *indykite.Client
	)

	BeforeEach(func(ctx SpecContext) {
		cli = ciqClient(ctx)
		admin, runtim = adminAndRuntimeClients(ctx)
	})

	It("upserts values referenced from the matched graph", func(ctx SpecContext) {
		projectID := os.Getenv("PROJECT_ID")
		unique := strconv.FormatInt(time.Now().UnixNano(), 10)
		ownerEmail := "sdk-it-ciq2-" + unique + "@example.com"

		// --- seed: the _Application subject binds to the caller's own
		// application node, so executing needs only the App Agent key.
		policy, err := admin.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
			ProjectID: projectID,
			Name:      "sdk-it-ciq2-ref-" + unique,
			Policy: `{
			  "meta": {"policy_version": "2.0-ciq"},
			  "subject": {"type": "_Application"},
			  "condition": {
			    "cypher": "MATCH (subject:_Application) MATCH (owner:Person)-[:OWNS]->(server:Server)",
			    "filter": [{
			      "operator": "AND",
			      "operands": [
			        {"operator": "=", "attribute": "@subject.external_id", "value": "$_appId"},
			        {"operator": "=", "attribute": "@server.external_id", "value": "$server_external_id"}
			      ]
			    }]
			  },
			  "allowed_reads": {"nodes": ["server.*"]},
			  "allowed_upserts": {"nodes": {"existing_nodes": ["server"]}}
			}`,
			Status: config.StatusActive,
		})
		Expect(err).NotTo(HaveOccurred(), "create 2.0-ciq policy")
		DeferCleanup(func(ctx SpecContext) {
			_ = admin.AuthorizationPolicies().Delete(ctx, policy.ID, "")
		})

		kq, err := admin.KnowledgeQueries().Create(ctx, &config.CreateKnowledgeQuery{
			ProjectID: projectID,
			Name:      "sdk-it-ciq2-ref-kq-" + unique,
			Query: `{
			  "nodes": ["server.property.owner_id", "server.property.owner_email"],
			  "upsert_nodes": [{
			    "name": "server",
			    "properties": [
			      {"type": "owner_id", "value": "@owner.external_id"},
			      {"type": "owner_email", "value": "@owner.property.email"}
			    ]
			  }]
			}`,
			Status:   config.StatusActive,
			PolicyID: policy.ID,
		})
		Expect(err).NotTo(HaveOccurred(), "create knowledge query")
		// Registered after the policy cleanup, so it runs first (LIFO): the
		// query references the policy.
		DeferCleanup(func(ctx SpecContext) {
			_ = admin.KnowledgeQueries().Delete(ctx, kq.ID, "")
		})

		// --- ingest: the owner whose attribute and property are referenced ---
		owner := capture.Node{ExternalID: "sdk-it-ciq2-owner-" + unique, Type: "Person"}
		// No region property: the fixture knowledge query (region=eu) must not
		// see this transient server.
		server := capture.Node{ExternalID: "sdk-it-ciq2-server-" + unique, Type: "Server"}
		owns := capture.Relationship{Type: "OWNS", Source: &owner, Target: &server}
		DeferCleanup(func(ctx SpecContext) {
			_, _ = runtim.Capture().DeleteRelationships(ctx, owns)
			_, _ = runtim.Capture().DeleteNodes(ctx, owner, server)
		})

		_, err = runtim.Capture().UpsertNodes(ctx,
			capture.UpsertNode{
				Node:       owner,
				IsIdentity: true,
				Properties: []capture.Property{{BaseProperty: capture.BaseProperty{Type: "email", Value: ownerEmail}}},
			},
			capture.UpsertNode{Node: server},
		)
		Expect(err).NotTo(HaveOccurred(), "ingest nodes")
		_, err = runtim.Capture().UpsertRelationships(ctx, owns)
		Expect(err).NotTo(HaveOccurred(), "ingest relationship")

		// --- act + assert: the write is idempotent, so poll until policy,
		// query and data propagate and the referenced values come back.
		Eventually(func(g Gomega) map[string]any {
			resp, err := cli.Execute(ctx, ciq.ExecuteRequest{
				ID:          kq.ID,
				InputParams: map[string]any{"server_external_id": server.ExternalID},
			})
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(resp.Data).To(HaveLen(1))
			return resp.Data[0].Nodes
		}).WithContext(ctx).WithTimeout(pollTimeout).WithPolling(pollInterval).Should(And(
			HaveKeyWithValue("server.property.owner_id", owner.ExternalID),
			HaveKeyWithValue("server.property.owner_email", ownerEmail),
		), "server %s never got the referenced owner values", server.ExternalID)
	}, SpecTimeout(2*time.Minute))
})

// adminAndRuntimeClients builds the control-plane and runtime clients the
// end-to-end spec needs, skipping when credentials or PROJECT_ID are missing.
func adminAndRuntimeClients(ctx context.Context) (*config.AdminClient, *indykite.Client) {
	GinkgoHelper()
	if os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS_FILE") == "" {
		Skip("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE] not set")
	}
	if os.Getenv("PROJECT_ID") == "" {
		Skip("PROJECT_ID not set")
	}
	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	admin, err := indykite.NewAdminFromEnv(ctx, opts...)
	Expect(err).NotTo(HaveOccurred())
	runtim, err := indykite.NewClientFromEnv(ctx, opts...)
	Expect(err).NotTo(HaveOccurred())
	return admin, runtim
}

// signingKey is a freshly generated RSA key pair: the private jwk.Key used to
// sign tokens and the public half as a JWK JSON string, the format expected by
// offline_validation.public_jwks.
type signingKey struct {
	private   jwk.Key
	publicJWK string
}

// newSigningKey generates a 2048-bit RSA signingKey labelled with kid.
func newSigningKey(kid string) signingKey {
	GinkgoHelper()
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	Expect(err).NotTo(HaveOccurred())
	key, err := jwk.Import[jwk.Key](raw)
	Expect(err).NotTo(HaveOccurred())
	Expect(key.Set(jwk.KeyIDKey, kid)).To(Succeed())
	Expect(key.Set(jwk.AlgorithmKey, jwa.RS256())).To(Succeed())

	pub, err := jwk.PublicKeyOf(key)
	Expect(err).NotTo(HaveOccurred())
	pubJSON, err := json.Marshal(pub)
	Expect(err).NotTo(HaveOccurred())
	return signingKey{private: key, publicJWK: string(pubJSON)}
}

// signEndUserToken mints a JWT with the standard claims the Token Introspect
// config matches on (iss, aud) and resolves by (sub).
func signEndUserToken(key jwk.Key, issuer, audience, subject string, lifetime time.Duration) string {
	GinkgoHelper()
	now := time.Now()
	t := jwt.New()
	Expect(t.Set(jwt.IssuerKey, issuer)).To(Succeed())
	Expect(t.Set(jwt.AudienceKey, audience)).To(Succeed())
	Expect(t.Set(jwt.SubjectKey, subject)).To(Succeed())
	Expect(t.Set(jwt.IssuedAtKey, now)).To(Succeed())
	Expect(t.Set(jwt.ExpirationKey, now.Add(lifetime))).To(Succeed())
	Expect(t.Set(jwt.JwtIDKey, uuid.New().String())).To(Succeed())

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256(), key))
	Expect(err).NotTo(HaveOccurred())
	return string(signed)
}

func mustJSON(v any) json.RawMessage {
	GinkgoHelper()
	raw, err := json.Marshal(v)
	Expect(err).NotTo(HaveOccurred())
	return raw
}
