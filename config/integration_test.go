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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NOTE: there is deliberately no integration test for the data-schema API
// (config.DataSchemaAPI) — it is still behind a feature flag.

func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// expectAPIStatus asserts err is an *transport.APIError carrying status.
func expectAPIStatus(err error, status int) {
	GinkgoHelper()
	Expect(err).To(HaveOccurred())
	apiErr, ok := transport.AsAPIError(err)
	Expect(ok).To(BeTrue(), "got %v, want an APIError", err)
	Expect(apiErr.StatusCode).To(Equal(status), "got %v", err)
}

var _ = Describe("Config (live platform)", Label("integration"), func() {
	var admin *config.AdminClient

	BeforeEach(func(ctx SpecContext) {
		if os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE] not set")
		}
		a, err := auth.ServiceAccountFromEnv(ctx)
		Expect(err).NotTo(HaveOccurred())
		var opts []transport.Option
		if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
			opts = append(opts, transport.WithBaseURL(base))
		}
		tc, err := transport.NewClient(a, opts...)
		Expect(err).NotTo(HaveOccurred())
		admin = config.NewAdminClient(tc)
	})

	requireProject := func() string {
		GinkgoHelper()
		id := os.Getenv("PROJECT_ID")
		if id == "" {
			Skip("PROJECT_ID not set")
		}
		return id
	}

	It("reads the organization and its projects", func(ctx SpecContext) {
		org, err := admin.Organizations().ReadCurrent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(org.ID).NotTo(BeEmpty(), "organization has no id: %+v", org)

		projects, err := admin.Projects().List(ctx, org.ID)
		Expect(err).NotTo(HaveOccurred())
		GinkgoWriter.Printf("organization %s has %d projects\n", org.ID, len(projects))
	})

	// The platform reports v1 for projects that predate the field; an
	// environment whose platform predates it omits the field altogether.
	It("reports the IKG schema version of the project", func(ctx SpecContext) {
		project := requireProject()

		got, err := admin.Projects().Read(ctx, project)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.IkgSchemaVersion).To(Or(BeEmpty(), Equal(config.IkgSchemaVersionV1)))
	})

	// The platform allows only one audit-signing config per project, and
	// `go run ./test/setup apply` provisions it (AUDIT_SIGNING_ID). The
	// audit-signing specs therefore never create their own config: they work on
	// the fixture and restore it when they are done.
	Describe("audit signing", func() {
		// auditSigningFixture returns the id of the fixture audit-signing config
		// and registers a cleanup that restores it to its provisioned
		// PLATFORM_MANAGED state, so a failed spec does not leave it changed for
		// the next run.
		auditSigningFixture := func(ctx context.Context) string {
			GinkgoHelper()
			id := os.Getenv("AUDIT_SIGNING_ID")
			if id == "" {
				Skip("AUDIT_SIGNING_ID not set")
			}
			api := admin.AuditSignings()
			orig, err := api.Read(ctx, id)
			Expect(err).NotTo(HaveOccurred(), "read fixture")
			DeferCleanup(func(ctx SpecContext) {
				cur, err := api.Read(ctx, id)
				Expect(err).NotTo(HaveOccurred(), "restore fixture: Read")
				_, err = api.Update(ctx, id, cur.ETag, &config.UpdateAuditSigning{
					DisplayName: &orig.DisplayName,
					Description: &orig.Description,
					AuditSigningConfig: config.AuditSigningConfig{
						Provider: config.AuditSigningProviderPlatformManaged,
					},
				})
				Expect(err).NotTo(HaveOccurred(), "restore fixture: Update")
			})
			return id
		}

		// Reads the fixture and checks the platform-managed shape: provider
		// set, no key material, secrets never echoed. A second config in the
		// same project is refused.
		It("serves the platform-managed fixture and refuses a second config", func(ctx SpecContext) {
			id := os.Getenv("AUDIT_SIGNING_ID")
			if id == "" {
				Skip("AUDIT_SIGNING_ID not set")
			}
			project := requireProject()

			signing, err := admin.AuditSignings().Read(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(signing.ID).To(Equal(id))
			Expect(signing.Provider).To(Equal(config.AuditSigningProviderPlatformManaged))
			Expect(signing.ETag).NotTo(BeEmpty())
			for k, v := range signing.AuthParams {
				Expect(v).To(BeEmpty(), "auth_params[%q] echoed a value; the platform must mask secrets on read", k)
			}

			items, err := admin.AuditSignings().List(ctx, project)
			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(ContainElement(HaveField("ID", id)),
				"fixture not present in AuditSignings.List of the project")

			// The project already holds the fixture; a second config must be refused.
			created, err := admin.AuditSignings().Create(ctx, &config.CreateAuditSigning{
				ProjectID: project,
				Name:      uniqueName("sdk-it-audit-signing"),
				AuditSigningConfig: config.AuditSigningConfig{
					Provider: config.AuditSigningProviderPlatformManaged,
				},
			})
			if err == nil {
				DeferCleanup(func(ctx SpecContext) { _ = admin.AuditSignings().Delete(ctx, created.ID, "") })
			}
			expectAPIStatus(err, http.StatusConflict)
		})

		// Runs the ETag-guarded read and update paths on the fixture: read by
		// id, by name (project_id and deprecated location; a bare name is
		// refused) -> update -> stale-ETag update is refused -> customer-managed
		// provider without key material is rejected -> list. Create and delete
		// are not run here: the project already holds its only config, and the
		// fixture must survive the spec.
		//
		// It deliberately stays PLATFORM_MANAGED throughout: a customer-managed
		// config pointing at fake KMS material would be accepted by the config
		// API but could then break audit signing for the shared test project
		// until cleanup. The customer-managed path has its own spec below, gated
		// on real key material.
		It("reads and updates the fixture with ETag guards", func(ctx SpecContext) {
			project := requireProject()
			api := admin.AuditSignings()
			id := auditSigningFixture(ctx)

			signing, err := api.Read(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(signing.ID).To(Equal(id))
			Expect(signing.Provider).To(Equal(config.AuditSigningProviderPlatformManaged))
			Expect(signing.ETag).NotTo(BeEmpty(), "Read returned no ETag")
			name := signing.Name

			byName, err := api.Read(ctx, name, config.WithProjectID(project))
			Expect(err).NotTo(HaveOccurred())
			Expect(byName.ID).To(Equal(id))

			// The deprecated location parameter must keep resolving names.
			byLocation, err := api.Read(ctx, name, config.WithLocation(project))
			Expect(err).NotTo(HaveOccurred())
			Expect(byLocation.ID).To(Equal(id))

			// A name is ambiguous without its parent project; the platform must refuse it.
			_, err = api.Read(ctx, name)
			expectAPIStatus(err, http.StatusUnprocessableEntity)

			displayName := uniqueName("SDK IT audit signing (renamed)")
			updated, err := api.Update(ctx, signing.ID, signing.ETag, &config.UpdateAuditSigning{
				DisplayName: &displayName,
				AuditSigningConfig: config.AuditSigningConfig{
					Provider: config.AuditSigningProviderPlatformManaged,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.ETag).NotTo(BeEmpty())
			Expect(updated.ETag).NotTo(Equal(signing.ETag))

			got, err := api.Read(ctx, signing.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.DisplayName).To(Equal(displayName))
			Expect(got.Provider).To(Equal(config.AuditSigningProviderPlatformManaged))

			// The pre-update ETag is stale now; the platform must refuse to apply it.
			_, err = api.Update(ctx, signing.ID, signing.ETag, &config.UpdateAuditSigning{
				AuditSigningConfig: config.AuditSigningConfig{
					Provider: config.AuditSigningProviderPlatformManaged,
				},
			})
			expectAPIStatus(err, http.StatusPreconditionFailed)

			// Customer-managed providers require key material. This write is
			// rejected by validation before anything is persisted, so it never
			// leaves the config pointing at a non-existent key.
			_, err = api.Update(ctx, signing.ID, got.ETag, &config.UpdateAuditSigning{
				AuditSigningConfig: config.AuditSigningConfig{
					Provider: config.AuditSigningProviderCustomerAWSKMS,
				},
			})
			Expect(err).To(HaveOccurred(), "CUSTOMER_AWS_KMS without key_resource/kid succeeded")
			apiErr, ok := transport.AsAPIError(err)
			Expect(ok).To(BeTrue(), "got %v, want an APIError", err)
			Expect(apiErr.StatusCode).To(And(BeNumerically(">=", 400), BeNumerically("<", 500)))

			items, err := api.List(ctx, project)
			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(ConsistOf(HaveField("ID", id)), "want exactly the fixture")
		})

		// Switches the fixture to customer-managed key material and verifies the
		// key fields round-trip while auth-param secrets are masked on read. The
		// fixture is restored to PLATFORM_MANAGED afterwards. It is opt-in:
		//
		//	AUDIT_SIGNING_KEY_RESOURCE  provider key reference (required to run)
		//	AUDIT_SIGNING_KID           key id (required to run)
		//	AUDIT_SIGNING_PROVIDER      CUSTOMER_GCP_KMS (default) | CUSTOMER_AWS_KMS | CUSTOMER_AZURE_KEY_VAULT
		//	AUDIT_SIGNING_AUTH_PARAMS   optional JSON object of provider credentials
		It("round-trips customer-managed key material", func(ctx SpecContext) {
			keyResource, kid := os.Getenv("AUDIT_SIGNING_KEY_RESOURCE"), os.Getenv("AUDIT_SIGNING_KID")
			if keyResource == "" || kid == "" {
				Skip("AUDIT_SIGNING_KEY_RESOURCE / AUDIT_SIGNING_KID not set")
			}
			provider := os.Getenv("AUDIT_SIGNING_PROVIDER")
			if provider == "" {
				provider = config.AuditSigningProviderCustomerGCPKMS
			}
			var authParams map[string]string
			if raw := os.Getenv("AUDIT_SIGNING_AUTH_PARAMS"); raw != "" {
				Expect(json.Unmarshal([]byte(raw), &authParams)).To(Succeed(),
					"AUDIT_SIGNING_AUTH_PARAMS is not a JSON object of strings")
			}
			api := admin.AuditSignings()
			id := auditSigningFixture(ctx)

			current, err := api.Read(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			_, err = api.Update(ctx, id, current.ETag, &config.UpdateAuditSigning{
				AuditSigningConfig: config.AuditSigningConfig{
					Provider:    provider,
					KeyResource: keyResource,
					Kid:         kid,
					AuthParams:  authParams,
				},
			})
			Expect(err).NotTo(HaveOccurred(), "Update to %s", provider)

			got, err := api.Read(ctx, id)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Provider).To(Equal(provider))
			Expect(got.KeyResource).To(Equal(keyResource))
			Expect(got.Kid).To(Equal(kid))
			for k := range authParams {
				Expect(got.AuthParams).To(HaveKey(k), "keys must be preserved on read")
				Expect(got.AuthParams[k]).To(BeEmpty(),
					"AuthParams[%q] echoed a value; the platform must mask secrets on read", k)
			}
		})
	})

	// The full ETag-guarded lifecycle: create draft -> read -> activate -> delete.
	It("runs the authorization policy lifecycle", func(ctx SpecContext) {
		project := requireProject()
		api := admin.AuthorizationPolicies()

		policyJSON := `{
	  "meta": {"policy_version": "2.0-kbac"},
	  "subject": {"type": "Person"},
	  "actions": ["SDK_IT_CAN_READ"],
	  "resource": {"type": "Asset"},
	  "condition": {"cypher": "MATCH (subject:Person)-[:OWNS]->(resource:Asset)"}
	}`
		name := uniqueName("sdk-it-policy")
		created, err := api.Create(ctx, &config.CreateAuthorizationPolicy{
			ProjectID: project,
			Name:      name,
			Policy:    policyJSON,
			Status:    config.StatusDraft,
		})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func(ctx SpecContext) { _ = api.Delete(ctx, created.ID, "") })

		pol, err := api.Read(ctx, created.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(pol.Status).To(Equal(config.StatusDraft))
		Expect(pol.ETag).NotTo(BeEmpty(), "Read returned no ETag")

		byName, err := api.Read(ctx, name, config.WithProjectID(project))
		Expect(err).NotTo(HaveOccurred())
		Expect(byName.ID).To(Equal(created.ID))

		_, err = api.Update(ctx, pol.ID, pol.ETag, &config.UpdateAuthorizationPolicy{
			Policy: pol.Policy,
			Status: config.StatusActive,
		})
		Expect(err).NotTo(HaveOccurred())

		got, err := api.Read(ctx, pol.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Status).To(Equal(config.StatusActive))

		Expect(api.Delete(ctx, got.ID, got.ETag)).To(Succeed())
		_, err = api.Read(ctx, got.ID)
		expectAPIStatus(err, http.StatusNotFound)
	})

	// Covers the validation modes beyond pinned keys: offline validation that
	// discovers keys through RFC 8414 metadata, then a switch to RFC 7662 token
	// introspection, whose client secret is write-only.
	It("runs the token introspect lifecycle across validation modes", func(ctx SpecContext) {
		project := requireProject()
		api := admin.TokenIntrospects()

		name := uniqueName("sdk-it-ti")
		jwtMatcher := json.RawMessage(fmt.Sprintf(
			`{"issuer":"https://sdk-it.example.com/%s","audience":%q}`, name, name))
		created, err := api.Create(ctx, &config.CreateTokenIntrospect{
			ProjectID: project,
			Name:      name,
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType: "Person",
				JwtMatcher:  jwtMatcher,
				OfflineValidation: json.RawMessage(
					`{"discovery":"` + config.TokenIntrospectDiscoveryRFC8414 + `"}`),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func(ctx SpecContext) { _ = api.Delete(ctx, created.ID, "") })

		got, err := api.Read(ctx, created.ID)
		Expect(err).NotTo(HaveOccurred())
		var offline struct {
			Discovery string `json:"discovery"`
		}
		Expect(json.Unmarshal(got.OfflineValidation, &offline)).To(Succeed())
		Expect(offline.Discovery).To(Equal(config.TokenIntrospectDiscoveryRFC8414))

		oauth2 := `{"introspection_endpoint":"https://sdk-it.example.com/introspect",
			"secret_basic":{"client_id":"sdk-it-client","client_secret":"sdk-it-secret"},"cache_ttl":300}`
		_, err = api.Update(ctx, got.ID, got.ETag, &config.UpdateTokenIntrospect{
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType:            "Person",
				JwtMatcher:             jwtMatcher,
				OnlineOAuth2Validation: json.RawMessage(oauth2),
			},
		})
		Expect(err).NotTo(HaveOccurred())

		got, err = api.Read(ctx, got.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.OfflineValidation).To(Or(BeEmpty(), MatchJSON(`null`)))
		var oauth2Read map[string]any
		Expect(json.Unmarshal(got.OnlineOAuth2Validation, &oauth2Read)).To(Succeed())
		Expect(oauth2Read).To(HaveKeyWithValue("introspection_endpoint", "https://sdk-it.example.com/introspect"))
		Expect(oauth2Read).To(HaveKeyWithValue("secret_basic", And(
			HaveKeyWithValue("client_id", "sdk-it-client"),
			Not(HaveKey("client_secret")),
		)), "the client secret must never be returned")

		// Keeping the endpoint and client id lets the update omit the secret.
		_, err = api.Update(ctx, got.ID, got.ETag, &config.UpdateTokenIntrospect{
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType:            "Person",
				JwtMatcher:             jwtMatcher,
				OnlineOAuth2Validation: got.OnlineOAuth2Validation,
				PerformUpsert:          true,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		got, err = api.Read(ctx, got.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.PerformUpsert).To(BeTrue())

		Expect(api.Delete(ctx, got.ID, got.ETag)).To(Succeed())
		_, err = api.Read(ctx, got.ID)
		expectAPIStatus(err, http.StatusNotFound)
	})

	// Chains application -> agent -> credential (the runtime-plane bootstrap)
	// and tears it down in reverse.
	It("runs the app agent lifecycle", func(ctx SpecContext) {
		project := requireProject()

		app, err := admin.Applications().Create(ctx, &config.CreateApplication{
			ProjectID: project, Name: uniqueName("sdk-it-app"),
		})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func(ctx SpecContext) { _ = admin.Applications().Delete(ctx, app.ID, "") })

		agent, err := admin.AppAgents().Create(ctx, &config.CreateAppAgent{
			ApplicationID:  app.ID,
			Name:           uniqueName("sdk-it-agent"),
			APIPermissions: []string{config.PermissionAuthorization},
		})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func(ctx SpecContext) { _ = admin.AppAgents().Delete(ctx, agent.ID, "") })

		displayName := "sdk-it renamed"
		_, err = admin.AppAgents().Update(ctx, agent.ID, "", &config.UpdateAppAgent{
			DisplayName: &displayName,
		})
		Expect(err).NotTo(HaveOccurred())

		cred, err := admin.AppAgentCredentials().Create(ctx, &config.CreateAppAgentCredential{
			ApplicationAgentID: agent.ID, DisplayName: "sdk-it-cred",
		})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func(ctx SpecContext) { _ = admin.AppAgentCredentials().Delete(ctx, cred.ID) })
		Expect(cred.ApplicationAgentConfig).NotTo(BeEmpty(), "credential Create returned no signed config")

		creds, err := admin.AppAgentCredentials().List(ctx, project)
		Expect(err).NotTo(HaveOccurred())
		Expect(creds).To(ContainElement(HaveField("ID", cred.ID)), "minted credential not in list")
	})

	// Every remaining project-scoped List.
	DescribeTable("lists project-scoped resources",
		func(ctx SpecContext, list func(context.Context, *config.AdminClient, string) (int, error)) {
			project := requireProject()
			n, err := list(ctx, admin, project)
			Expect(err).NotTo(HaveOccurred())
			GinkgoWriter.Printf("%d items\n", n)
		},
		Entry("applications", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.Applications().List(ctx, p)
			return len(items), err
		}),
		Entry("app-agents", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.AppAgents().List(ctx, p)
			return len(items), err
		}),
		Entry("authorization-policies", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.AuthorizationPolicies().List(ctx, p, "")
			return len(items), err
		}),
		Entry("knowledge-queries", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.KnowledgeQueries().List(ctx, p)
			return len(items), err
		}),
		Entry("event-sinks", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.EventSinks().List(ctx, p)
			return len(items), err
		}),
		Entry("token-introspects", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.TokenIntrospects().List(ctx, p)
			return len(items), err
		}),
		Entry("external-data-resolvers", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.ExternalDataResolvers().List(ctx, p)
			return len(items), err
		}),
		Entry("trust-score-profiles", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.TrustScoreProfiles().List(ctx, p)
			return len(items), err
		}),
		Entry("entity-matching-pipelines", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.EntityMatchingPipelines().List(ctx, p)
			return len(items), err
		}),
		Entry("mcp-servers", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.MCPServers().List(ctx, p)
			return len(items), err
		}),
		Entry("audit-signings", func(ctx context.Context, a *config.AdminClient, p string) (int, error) {
			items, err := a.AuditSignings().List(ctx, p)
			return len(items), err
		}),
	)

	It("lists the organization's service account credentials", func(ctx SpecContext) {
		org := os.Getenv("ORGANIZATION_ID")
		if org == "" {
			Skip("ORGANIZATION_ID not set")
		}
		creds, err := admin.ServiceAccountCredentials().List(ctx, org)
		Expect(err).NotTo(HaveOccurred())
		GinkgoWriter.Printf("service-account-credentials: %d\n", len(creds))
	})
})
