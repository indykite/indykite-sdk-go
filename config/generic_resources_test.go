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
	"testing"

	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/transport"
)

// crudCase describes one generic resource sub-API in terms of closures so a
// single table-driven test can walk List/Create/Read/Update/Delete for all of
// them.
type crudCase struct {
	list     func(ctx context.Context, a *config.AdminClient) error
	create   func(ctx context.Context, a *config.AdminClient) (*config.WriteResult, error)
	read     func(ctx context.Context, a *config.AdminClient, id string) (string, error)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				p, err := a.Projects().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				app, err := a.Applications().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				sa, err := a.ServiceAccounts().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				ti, err := a.TokenIntrospects().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				r, err := a.ExternalDataResolvers().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				p, err := a.TrustScoreProfiles().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				p, err := a.EntityMatchingPipelines().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				s, err := a.MCPServers().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				ag, err := a.AppAgents().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				kq, err := a.KnowledgeQueries().Read(ctx, id)
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
			read: func(ctx context.Context, a *config.AdminClient, id string) (string, error) {
				pol, err := a.AuthorizationPolicies().Read(ctx, id)
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

// TestResourceCRUD walks List/Create/Read/Update/Delete for every generic
// resource sub-API and asserts method, path, query, If-Match and ETag capture.
func TestResourceCRUD(t *testing.T) {
	for _, tc := range crudCases() {
		t.Run(tc.name, func(t *testing.T) {
			var rec recorder
			admin := adminTo(t, crudHandler(&rec, tc.basePath, genericCRUDResponses("gid:res")))
			ctx := context.Background()
			idPath := tc.basePath + "/gid:res"

			if err := tc.list(ctx, admin); err != nil {
				t.Fatalf("List: %v", err)
			}
			rec.wantReq(t, http.MethodGet, tc.basePath)
			if rec.query.Get(tc.scopeKey) != "gid:scope" {
				t.Errorf("list query[%s] = %q, want gid:scope", tc.scopeKey, rec.query.Get(tc.scopeKey))
			}

			created, err := tc.create(ctx, admin)
			if err != nil {
				t.Fatalf("Create: %v", err)
			}
			rec.wantReq(t, http.MethodPost, tc.basePath)
			wantWrite(t, created, `"v1"`)

			etag, err := tc.read(ctx, admin, "gid:res")
			if err != nil {
				t.Fatalf("Read: %v", err)
			}
			rec.wantReq(t, http.MethodGet, idPath)
			if etag != `"v2"` {
				t.Errorf("read ETag = %q, want \"v2\"", etag)
			}

			updated, err := tc.update(ctx, admin, "gid:res", etag)
			if err != nil {
				t.Fatalf("Update: %v", err)
			}
			rec.wantReq(t, http.MethodPut, idPath)
			wantWrite(t, updated, `"v3"`)
			if rec.ifMatch != `"v2"` {
				t.Errorf("update If-Match = %q, want \"v2\"", rec.ifMatch)
			}

			if err = tc.remove(ctx, admin, "gid:res", updated.ETag); err != nil {
				t.Fatalf("Delete: %v", err)
			}
			rec.wantReq(t, http.MethodDelete, idPath)
			if rec.ifMatch != `"v3"` {
				t.Errorf("delete If-Match = %q, want \"v3\"", rec.ifMatch)
			}
		})
	}
}

// TestReadAndListOptions asserts the query parameters produced by the
// Read/List option helpers.
func TestReadAndListOptions(t *testing.T) {
	var query url.Values
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query()
		if r.URL.Path == "/configs/v1/projects" {
			_, _ = io.WriteString(w, `{"data":[]}`)
		} else {
			_, _ = io.WriteString(w, `{"id":"gid:proj","name":"proj"}`)
		}
	})
	ctx := context.Background()

	if _, err := admin.Projects().Read(ctx, "proj",
		config.WithLocation("gid:org"), config.WithVersion(3)); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if query.Get("location") != "gid:org" || query.Get("version") != "3" {
		t.Errorf("read query = %v, want location=gid:org version=3", query)
	}

	if _, err := admin.Projects().List(ctx, "gid:org",
		config.WithSearch("prod"), config.WithFullFetch()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if query.Get("organization_id") != "gid:org" || query.Get("search") != "prod" ||
		query.Get("full_fetch") != "true" {
		t.Errorf("list query = %v, want organization_id=gid:org search=prod full_fetch=true", query)
	}
}

// TestAPIErrorPropagation asserts that non-2xx responses surface as
// *transport.APIError on every CRUD verb.
func TestAPIErrorPropagation(t *testing.T) {
	admin := adminTo(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"message":"resource not found","code":"NOT_FOUND"}`)
	})
	ctx := context.Background()

	checkErr := func(t *testing.T, op string, err error) {
		t.Helper()
		if err == nil {
			t.Fatalf("%s: expected error, got nil", op)
		}
		apiErr, ok := transport.AsAPIError(err)
		if !ok {
			t.Fatalf("%s: AsAPIError = false for %v", op, err)
		}
		if apiErr.StatusCode != http.StatusNotFound || !apiErr.IsNotFound() {
			t.Errorf("%s: StatusCode = %d, want 404", op, apiErr.StatusCode)
		}
		if apiErr.Message != "resource not found" {
			t.Errorf("%s: Message = %q", op, apiErr.Message)
		}
	}

	_, err := admin.MCPServers().List(ctx, "gid:proj")
	checkErr(t, "List", err)

	_, err = admin.MCPServers().Read(ctx, "gid:missing")
	checkErr(t, "Read", err)

	_, err = admin.MCPServers().Create(ctx, &config.CreateMCPServer{ProjectID: "gid:proj", Name: "mcp"})
	checkErr(t, "Create", err)

	checkErr(t, "Delete", admin.MCPServers().Delete(ctx, "gid:missing", `"v1"`))

	_, err = admin.AppAgentCredentials().Read(ctx, "gid:missing")
	checkErr(t, "CredentialRead", err)

	_, err = admin.DataSchema().Rebuild(ctx, "gid:proj")
	checkErr(t, "Rebuild", err)
}

// TestAppAgentCredentialLifecycle covers Create/Read/Delete of application
// agent credentials, which have no ETag and return the secret only on Create.
func TestAppAgentCredentialLifecycle(t *testing.T) {
	const base = "/configs/v1/application-agent-credentials"
	var rec recorder
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		if r.Method == http.MethodPost {
			_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1",
				"application_agent_id":"gid:agent","application_agent_config":{"appAgentId":"gid:agent"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1","application_agent_id":"gid:agent"}`)
	})
	ctx := context.Background()
	creds := admin.AppAgentCredentials()

	res, err := creds.Create(ctx, &config.CreateAppAgentCredential{
		ApplicationAgentID: "gid:agent", DisplayName: "ci",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	rec.wantReq(t, http.MethodPost, base)
	if rec.body["application_agent_id"] != "gid:agent" || rec.body["display_name"] != "ci" {
		t.Errorf("create body = %v", rec.body)
	}
	if res.Kid != "k1" || len(res.ApplicationAgentConfig) == 0 {
		t.Errorf("create result = %+v", res)
	}

	got, err := creds.Read(ctx, "gid:cred")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	rec.wantReq(t, http.MethodGet, base+"/gid:cred")
	if got.ID != "gid:cred" || got.ApplicationAgentID != "gid:agent" {
		t.Errorf("read result = %+v", got)
	}

	if err = creds.Delete(ctx, "gid:cred"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rec.wantReq(t, http.MethodDelete, base+"/gid:cred")
}

// TestServiceAccountCredentialReadDelete covers the metadata Read and revoking
// Delete of service account credentials (Create/List are covered elsewhere).
func TestServiceAccountCredentialReadDelete(t *testing.T) {
	var method, path string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		if r.Method == http.MethodGet {
			_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1","service_account_id":"gid:sa"}`)
			return
		}
		_, _ = io.WriteString(w, `{}`)
	})
	ctx := context.Background()
	creds := admin.ServiceAccountCredentials()

	got, err := creds.Read(ctx, "gid:cred")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if method != http.MethodGet || path != "/configs/v1/service-account-credentials/gid:cred" {
		t.Errorf("read method/path = %s %s", method, path)
	}
	if got.Kid != "k1" || got.ServiceAccountID != "gid:sa" {
		t.Errorf("read result = %+v", got)
	}

	if err = creds.Delete(ctx, "gid:cred"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if method != http.MethodDelete || path != "/configs/v1/service-account-credentials/gid:cred" {
		t.Errorf("delete method/path = %s %s", method, path)
	}
}

// TestAdminClientFromCredentialsInvalidJSON asserts the credential-parsing
// error path of NewAdminClientFromCredentials.
func TestAdminClientFromCredentialsInvalidJSON(t *testing.T) {
	if _, err := config.NewAdminClientFromCredentials(context.Background(), []byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid credentials JSON")
	}
}
