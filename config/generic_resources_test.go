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

package config_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// crudCase describes one generic resource sub-API in terms of closures so a
// single table can walk List/Create/Read/Update/Delete for all of them.
type crudCase struct {
	list     func(ctx context.Context, a *config.AdminClient) error
	create   func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error)
	read     func(ctx context.Context, a *config.AdminClient, id string, opts ...config.ReadOption) (string, error)
	update   func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error)
	remove   func(ctx context.Context, a *config.AdminClient, id, etag string) error
	name     string
	basePath string
	scopeKey string
}

func crudCases() []crudCase {
	return []crudCase{
		{
			name: "Projects", basePath: "/configs/v1/projects", scopeKey: "organization_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.Projects().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.Projects().Create(ctx, &config.CreateProject{
					OrganizationID: "gid:scope", Name: "proj", Region: "europe-west1",
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				p, err := a.Projects().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return p.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				name := "renamed"
				return a.Projects().Update(ctx, id, etag, &config.UpdateProject{DisplayName: &name})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.Projects().Delete(ctx, id, etag)
			},
		},
		{
			name: "Applications", basePath: "/configs/v1/applications", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.Applications().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.Applications().Create(ctx, &config.CreateApplication{
					ProjectID: "gid:scope", Name: "app",
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				app, err := a.Applications().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return app.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				name := "renamed"
				return a.Applications().Update(ctx, id, etag, &config.UpdateApplication{DisplayName: &name})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.Applications().Delete(ctx, id, etag)
			},
		},
		{
			name: "ServiceAccounts", basePath: "/configs/v1/service-accounts", scopeKey: "organization_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.ServiceAccounts().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.ServiceAccounts().Create(ctx, &config.CreateServiceAccount{
					OrganizationID: "gid:scope", Name: "sa", Role: config.RoleAllEditor,
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				sa, err := a.ServiceAccounts().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return sa.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				name := "renamed"
				return a.ServiceAccounts().Update(ctx, id, etag, &config.UpdateServiceAccount{DisplayName: &name})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.ServiceAccounts().Delete(ctx, id, etag)
			},
		},
		{
			name: "TokenIntrospects", basePath: "/configs/v1/token-introspects", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.TokenIntrospects().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.TokenIntrospects().Create(ctx, &config.CreateTokenIntrospect{
					ProjectID: "gid:scope", Name: "ti",
					TokenIntrospectConfig: config.TokenIntrospectConfig{IkgNodeType: "Person"},
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				ti, err := a.TokenIntrospects().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return ti.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.TokenIntrospects().Update(ctx, id, etag, &config.UpdateTokenIntrospect{
					TokenIntrospectConfig: config.TokenIntrospectConfig{IkgNodeType: "Person"},
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.TokenIntrospects().Delete(ctx, id, etag)
			},
		},
		{
			name: "ExternalDataResolvers", basePath: "/configs/v1/external-data-resolvers", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.ExternalDataResolvers().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.ExternalDataResolvers().Create(ctx, &config.CreateExternalDataResolver{
					ProjectID: "gid:scope", Name: "resolver",
					ExternalDataResolverConfig: config.ExternalDataResolverConfig{
						URL: "https://example.com", Method: http.MethodGet,
						RequestContentType: "application/json", ResponseContentType: "application/json",
						ResponseSelector: ".",
					},
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				r, err := a.ExternalDataResolvers().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return r.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.ExternalDataResolvers().Update(ctx, id, etag, &config.UpdateExternalDataResolver{
					ExternalDataResolverConfig: config.ExternalDataResolverConfig{
						URL: "https://example.com", Method: http.MethodPost,
						RequestContentType: "application/json", ResponseContentType: "application/json",
						ResponseSelector: ".",
					},
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.ExternalDataResolvers().Delete(ctx, id, etag)
			},
		},
		{
			name: "TrustScoreProfiles", basePath: "/configs/v1/trust-score-profiles", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.TrustScoreProfiles().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.TrustScoreProfiles().Create(ctx, &config.CreateTrustScoreProfile{
					ProjectID: "gid:scope", Name: "tsp", NodeClassification: "Person",
					Schedule: config.ScheduleDaily, Dimensions: json.RawMessage(`[]`),
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				p, err := a.TrustScoreProfiles().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return p.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.TrustScoreProfiles().Update(ctx, id, etag, &config.UpdateTrustScoreProfile{
					Schedule: config.ScheduleSixHours,
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.TrustScoreProfiles().Delete(ctx, id, etag)
			},
		},
		{
			name: "EntityMatchingPipelines", basePath: "/configs/v1/entity-matching-pipelines", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.EntityMatchingPipelines().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.EntityMatchingPipelines().Create(ctx, &config.CreateEntityMatchingPipeline{
					ProjectID: "gid:scope", Name: "emp",
					NodeFilter: json.RawMessage(`{}`), SimilarityScoreCutoff: 0.9,
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				p, err := a.EntityMatchingPipelines().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return p.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.EntityMatchingPipelines().Update(ctx, id, etag, &config.UpdateEntityMatchingPipeline{
					SimilarityScoreCutoff: 0.8,
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.EntityMatchingPipelines().Delete(ctx, id, etag)
			},
		},
		{
			name: "MCPServers", basePath: "/configs/v1/mcp-servers", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.MCPServers().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				enabled := true
				return a.MCPServers().Create(ctx, &config.CreateMCPServer{
					ProjectID: "gid:scope", Name: "mcp",
					MCPServerConfig: config.MCPServerConfig{
						Enabled: &enabled, AppAgentID: "gid:agent", TokenIntrospectID: "gid:ti",
						ScopesSupported: []string{"read"},
					},
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				s, err := a.MCPServers().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return s.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				enabled := false
				return a.MCPServers().Update(ctx, id, etag, &config.UpdateMCPServer{
					MCPServerConfig: config.MCPServerConfig{
						Enabled: &enabled, AppAgentID: "gid:agent", TokenIntrospectID: "gid:ti",
					},
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.MCPServers().Delete(ctx, id, etag)
			},
		},
		{
			name: "AuditSignings", basePath: "/configs/v1/audit-signings", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.AuditSignings().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.AuditSignings().Create(ctx, &config.CreateAuditSigning{
					ProjectID: "gid:scope", Name: "signing",
					AuditSigningConfig: config.AuditSigningConfig{
						Provider:    config.AuditSigningProviderCustomerGCPKMS,
						KeyResource: "projects/p/locations/l/keyRings/r/cryptoKeys/k",
						Kid:         "k1",
						AuthParams:  map[string]string{"credentials_json": "{}"},
					},
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				s, err := a.AuditSignings().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return s.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.AuditSignings().Update(ctx, id, etag, &config.UpdateAuditSigning{
					AuditSigningConfig: config.AuditSigningConfig{
						Provider: config.AuditSigningProviderPlatformManaged,
					},
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.AuditSignings().Delete(ctx, id, etag)
			},
		},
		{
			name: "AppAgents", basePath: "/configs/v1/application-agents", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.AppAgents().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.AppAgents().Create(ctx, &config.CreateAppAgent{
					ApplicationID: "gid:app", Name: "agent",
					APIPermissions: []string{config.PermissionAuthorization},
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				ag, err := a.AppAgents().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return ag.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.AppAgents().Update(ctx, id, etag, &config.UpdateAppAgent{
					APIPermissions: []string{config.PermissionCapture},
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.AppAgents().Delete(ctx, id, etag)
			},
		},
		{
			name: "KnowledgeQueries", basePath: "/configs/v1/knowledge-queries", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.KnowledgeQueries().List(ctx, "gid:scope")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.KnowledgeQueries().Create(ctx, &config.CreateKnowledgeQuery{
					ProjectID: "gid:scope", Name: "kq", Query: "{}",
					Status: config.StatusActive, PolicyID: "gid:pol",
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				kq, err := a.KnowledgeQueries().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return kq.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.KnowledgeQueries().Update(ctx, id, etag, &config.UpdateKnowledgeQuery{
					Query: "{}", Status: config.StatusInactive, PolicyID: "gid:pol",
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.KnowledgeQueries().Delete(ctx, id, etag)
			},
		},
		{
			name: "AuthorizationPolicies", basePath: "/configs/v1/authorization-policies", scopeKey: "project_id",
			list: func(ctx context.Context, a *config.AdminClient) error {
				_, err := a.AuthorizationPolicies().List(ctx, "gid:scope", "")
				return err
			},
			create: func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error) {
				return a.AuthorizationPolicies().Create(ctx, &config.CreateAuthorizationPolicy{
					ProjectID: "gid:scope", Name: "pol", Policy: "{}", Status: config.StatusActive,
				})
			},
			read: func(ctx context.Context, a *config.AdminClient, id string, o ...config.ReadOption) (string, error) {
				pol, err := a.AuthorizationPolicies().Read(ctx, id, o...)
				if err != nil {
					return "", err
				}
				return pol.ETag, nil
			},
			update: func(ctx context.Context, a *config.AdminClient, id, etag string) (*config.WriteResult, error) {
				return a.AuthorizationPolicies().Update(ctx, id, etag, &config.UpdateAuthorizationPolicy{
					Policy: "{}", Status: config.StatusInactive,
				})
			},
			remove: func(ctx context.Context, a *config.AdminClient, id, etag string) error {
				return a.AuthorizationPolicies().Delete(ctx, id, etag)
			},
		},
	}
}

// Every generic resource sub-API walks List/Create/Read/Update/Delete; the
// entries assert method, path, query, If-Match and ETag capture.
var _ = DescribeTable("generic resource CRUD",
	func(ctx SpecContext, tc crudCase) {
		var rec recorder
		admin := adminTo(crudHandler(&rec, tc.basePath, genericCRUDResponses("gid:res")))
		idPath := tc.basePath + "/gid:res"

		Expect(tc.list(ctx, admin)).To(Succeed())
		rec.expectReq(http.MethodGet, tc.basePath)
		Expect(rec.query.Get(tc.scopeKey)).To(Equal("gid:scope"), "list query[%s]", tc.scopeKey)

		created, err := tc.create(ctx, admin)
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodPost, tc.basePath)
		expectWrite(created, `"v1"`)

		etag, err := tc.read(ctx, admin, "gid:res")
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodGet, idPath)
		Expect(etag).To(Equal(`"v2"`))

		// Project-scoped names need their parent project; organization-scoped
		// names resolve within the caller's organization without a scope.
		var byNameOpts []config.ReadOption
		if tc.scopeKey == "project_id" {
			byNameOpts = append(byNameOpts, config.WithProjectID("gid:scope"))
		}
		_, err = tc.read(ctx, admin, "by-name", byNameOpts...)
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodGet, tc.basePath+"/by-name")
		if tc.scopeKey == "project_id" {
			Expect(rec.query).To(Equal(url.Values{"project_id": {"gid:scope"}}))
		} else {
			Expect(rec.query).To(BeEmpty(), "no scope for an organization-scoped resource")
		}

		updated, err := tc.update(ctx, admin, "gid:res", etag)
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodPut, idPath)
		expectWrite(updated, `"v3"`)
		Expect(rec.ifMatch).To(Equal(`"v2"`))

		Expect(tc.remove(ctx, admin, "gid:res", updated.ETag)).To(Succeed())
		rec.expectReq(http.MethodDelete, idPath)
		Expect(rec.ifMatch).To(Equal(`"v3"`))
	},
	crudEntries(),
)

func crudEntries() []TableEntry {
	cases := crudCases()
	entries := make([]TableEntry, 0, len(cases))
	for _, tc := range cases {
		entries = append(entries, Entry(tc.name, tc))
	}
	return entries
}

var _ = Describe("Read and List options", func() {
	var query url.Values
	var admin *config.AdminClient

	BeforeEach(func() {
		admin = adminTo(func(w http.ResponseWriter, r *http.Request) {
			query = r.URL.Query()
			if r.URL.Path == "/configs/v1/projects" {
				_, _ = io.WriteString(w, `{"data":[]}`)
			} else {
				_, _ = io.WriteString(w, `{"id":"gid:proj","name":"proj"}`)
			}
		})
	})

	It("sends project_id and version on read", func(ctx SpecContext) {
		_, err := admin.KnowledgeQueries().Read(ctx, "kq", config.WithProjectID("gid:proj"), config.WithVersion(3))
		Expect(err).NotTo(HaveOccurred())
		Expect(query).To(Equal(url.Values{"project_id": {"gid:proj"}, "version": {"3"}}))
	})

	// The backend rejects project_id together with location, so the last
	// option wins and the other parameter is dropped.
	DescribeTable("keeps only the last of project_id and location",
		func(ctx SpecContext, opts []config.ReadOption, want url.Values) {
			_, err := admin.KnowledgeQueries().Read(ctx, "kq", opts...)
			Expect(err).NotTo(HaveOccurred())
			Expect(query).To(Equal(want))
		},
		Entry("location last",
			[]config.ReadOption{config.WithProjectID("gid:proj"), config.WithLocation("gid:loc")},
			url.Values{"location": {"gid:loc"}}),
		Entry("project_id last",
			[]config.ReadOption{config.WithLocation("gid:loc"), config.WithProjectID("gid:proj")},
			url.Values{"project_id": {"gid:proj"}}),
	)

	It("sends the scope, search and full_fetch on list", func(ctx SpecContext) {
		_, err := admin.Projects().List(ctx, "gid:org", config.WithSearch("prod"), config.WithFullFetch())
		Expect(err).NotTo(HaveOccurred())
		Expect(query).To(Equal(url.Values{
			"organization_id": {"gid:org"}, "search": {"prod"}, "full_fetch": {"true"},
		}))
	})
})

// Non-2xx responses surface as *transport.APIError on every CRUD verb.
var _ = DescribeTable("API error propagation",
	func(ctx SpecContext, call func(context.Context, *config.AdminClient) error) {
		admin := adminTo(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"message":"resource not found","code":"NOT_FOUND"}`)
		})

		err := call(ctx, admin)
		Expect(err).To(HaveOccurred())
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "AsAPIError = false for %v", err)
		Expect(apiErr.StatusCode).To(Equal(http.StatusNotFound))
		Expect(apiErr.IsNotFound()).To(BeTrue())
		Expect(apiErr.Message).To(Equal("resource not found"))
	},
	Entry("List", func(ctx context.Context, a *config.AdminClient) error {
		_, err := a.MCPServers().List(ctx, "gid:proj")
		return err
	}),
	Entry("Read", func(ctx context.Context, a *config.AdminClient) error {
		_, err := a.MCPServers().Read(ctx, "gid:missing")
		return err
	}),
	Entry("Create", func(ctx context.Context, a *config.AdminClient) error {
		_, err := a.MCPServers().Create(ctx, &config.CreateMCPServer{ProjectID: "gid:proj", Name: "mcp"})
		return err
	}),
	Entry("Delete", func(ctx context.Context, a *config.AdminClient) error {
		return a.MCPServers().Delete(ctx, "gid:missing", `"v1"`)
	}),
	Entry("CredentialRead", func(ctx context.Context, a *config.AdminClient) error {
		_, err := a.AppAgentCredentials().Read(ctx, "gid:missing")
		return err
	}),
	Entry("Rebuild", func(ctx context.Context, a *config.AdminClient) error {
		_, err := a.DataSchema().Rebuild(ctx, "gid:proj")
		return err
	}),
)

// Application agent credentials have no ETag and return the secret only on Create.
var _ = It("runs Create/Read/Delete of application agent credentials", func(ctx SpecContext) {
	const base = "/configs/v1/application-agent-credentials"
	var rec recorder
	admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		if r.Method == http.MethodPost {
			_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1",
				"application_agent_id":"gid:agent","application_agent_config":{"appAgentId":"gid:agent"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1","application_agent_id":"gid:agent"}`)
	})
	creds := admin.AppAgentCredentials()

	res, err := creds.Create(ctx, &config.CreateAppAgentCredential{
		ApplicationAgentID: "gid:agent", DisplayName: "ci",
	})
	Expect(err).NotTo(HaveOccurred())
	rec.expectReq(http.MethodPost, base)
	Expect(rec.body).To(HaveKeyWithValue("application_agent_id", "gid:agent"))
	Expect(rec.body).To(HaveKeyWithValue("display_name", "ci"))
	Expect(res.Kid).To(Equal("k1"))
	Expect(res.ApplicationAgentConfig).NotTo(BeEmpty())

	got, err := creds.Read(ctx, "gid:cred")
	Expect(err).NotTo(HaveOccurred())
	rec.expectReq(http.MethodGet, base+"/gid:cred")
	Expect(got.ID).To(Equal("gid:cred"))
	Expect(got.ApplicationAgentID).To(Equal("gid:agent"))

	Expect(creds.Delete(ctx, "gid:cred")).To(Succeed())
	rec.expectReq(http.MethodDelete, base+"/gid:cred")
})

// Create/List of service account credentials are covered elsewhere.
var _ = It("reads and revokes service account credentials", func(ctx SpecContext) {
	const path = "/configs/v1/service-account-credentials/gid:cred"
	var rec recorder
	admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		if r.Method == http.MethodGet {
			_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1","service_account_id":"gid:sa"}`)
			return
		}
		_, _ = io.WriteString(w, `{}`)
	})
	creds := admin.ServiceAccountCredentials()

	got, err := creds.Read(ctx, "gid:cred")
	Expect(err).NotTo(HaveOccurred())
	rec.expectReq(http.MethodGet, path)
	Expect(got.Kid).To(Equal("k1"))
	Expect(got.ServiceAccountID).To(Equal("gid:sa"))

	Expect(creds.Delete(ctx, "gid:cred")).To(Succeed())
	rec.expectReq(http.MethodDelete, path)
})
