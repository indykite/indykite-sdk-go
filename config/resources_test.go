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
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/indykite/indykite-sdk-go/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// A standard generic resource exercised end to end.
var _ = It("runs the knowledge query lifecycle", func(ctx SpecContext) {
	const base = "/configs/v1/knowledge-queries"
	var rec recorder
	admin := adminTo(crudHandler(&rec, base, crudResponses{
		created: `{"id":"gid:kq"}`,
		list:    `{"data":[{"id":"gid:kq","name":"q","status":"ACTIVE"}]}`,
		read:    `{"id":"gid:kq","query":"{}","status":"ACTIVE","policy_id":"gid:p"}`,
		updated: `{"id":"gid:kq"}`,
	}))
	kq := admin.KnowledgeQueries()

	_, err := kq.List(ctx, "gid:proj")
	Expect(err).NotTo(HaveOccurred())
	Expect(rec.query.Get("project_id")).To(Equal("gid:proj"))

	created, err := kq.Create(ctx, &config.CreateKnowledgeQuery{
		ProjectID: "gid:proj", Name: "q", Query: "{}", Status: config.StatusActive, PolicyID: "gid:p",
	})
	Expect(err).NotTo(HaveOccurred())
	expectWrite(created, `"v1"`)

	got, err := kq.Read(ctx, "gid:kq")
	Expect(err).NotTo(HaveOccurred())
	Expect(got.PolicyID).To(Equal("gid:p"))
	Expect(got.ETag).To(Equal(`"v2"`))

	_, err = kq.Update(ctx, "gid:kq", got.ETag, &config.UpdateKnowledgeQuery{
		Query: "{}", Status: config.StatusInactive, PolicyID: "gid:p",
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(rec.ifMatch).To(Equal(`"v2"`))

	Expect(kq.Delete(ctx, "gid:kq", got.ETag)).To(Succeed())
	rec.expectReq(http.MethodDelete, base+"/gid:kq")
})

var _ = It("reads the current organization", func(ctx SpecContext) {
	var path string
	admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("ETag", `"o1"`)
		_, _ = io.WriteString(w, `{"id":"gid:org","name":"acme"}`)
	})

	org, err := admin.Organizations().ReadCurrent(ctx)
	Expect(err).NotTo(HaveOccurred())
	Expect(path).To(Equal("/configs/v1/organizations/current"))
	Expect(org.ID).To(Equal("gid:org"))
	Expect(org.ETag).To(Equal(`"o1"`))
})

var _ = It("returns the signed config of a new service account credential", func(ctx SpecContext) {
	var rec recorder
	admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1","service_account_config":{"appAgentId":"x"}}`)
	})

	res, err := admin.ServiceAccountCredentials().Create(ctx,
		&config.CreateServiceAccountCredential{ServiceAccountID: "gid:sa", DisplayName: "ci"})
	Expect(err).NotTo(HaveOccurred())
	Expect(rec.body).To(HaveKeyWithValue("service_account_id", "gid:sa"))
	Expect(res.Kid).To(Equal("k1"))
	Expect(res.ServiceAccountConfig).NotTo(BeEmpty())
})

var _ = It("runs the event sink lifecycle with typed providers and routes", func(ctx SpecContext) {
	const base = "/configs/v1/event-sinks"
	var rec recorder
	admin := adminTo(crudHandler(&rec, base, crudResponses{
		created: `{"id":"gid:sink"}`,
		list:    `{"data":[{"id":"gid:sink","name":"audit"}]}`,
		read: `{"id":"gid:sink","name":"audit",
			"providers":{"kafka-1":{"kafka":{"brokers":["b:9092"],"topic":"events"}}},
			"routes":[{"provider_id":"kafka-1",
				"event_type_key_values_filter":{"event_type":"indykite.audit.config.create"}}],
			"include_cdc_events":true}`,
		updated: `{"id":"gid:sink"}`,
	}))
	sinks := admin.EventSinks()

	_, err := sinks.List(ctx, "gid:proj")
	Expect(err).NotTo(HaveOccurred())
	Expect(rec.query.Get("project_id")).To(Equal("gid:proj"))

	created, err := sinks.Create(ctx, &config.CreateEventSink{
		ProjectID: "gid:proj",
		Name:      "audit",
		Providers: map[string]config.EventSinkProvider{
			"kafka-1": {Kafka: &config.KafkaSinkConfig{Brokers: []string{"b:9092"}, Topic: "events"}},
		},
		Routes: []config.EventSinkRoute{{
			ProviderID: "kafka-1",
			Filter:     config.EventSinkFilter{EventType: "indykite.audit.config.create"},
		}},
	})
	Expect(err).NotTo(HaveOccurred())
	expectWrite(created, `"v1"`)
	Expect(rec.body).To(HaveKeyWithValue("providers", HaveKey("kafka-1")))

	got, err := sinks.Read(ctx, "gid:sink")
	Expect(err).NotTo(HaveOccurred())
	Expect(got.ETag).To(Equal(`"v2"`))
	Expect(got.IncludeCDCEvents).To(BeTrue())
	Expect(got.Providers).To(HaveKey("kafka-1"))
	Expect(got.Providers["kafka-1"].Kafka.Topic).To(Equal("events"))

	_, err = sinks.Read(ctx, "audit", config.WithProjectID("gid:proj"))
	Expect(err).NotTo(HaveOccurred())
	rec.expectReq(http.MethodGet, base+"/audit")
	Expect(rec.query.Get("project_id")).To(Equal("gid:proj"))

	_, err = sinks.Update(ctx, "gid:sink", got.ETag, &config.UpdateEventSink{
		Providers: got.Providers, Routes: got.Routes,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(rec.ifMatch).To(Equal(`"v2"`))

	Expect(sinks.Delete(ctx, "gid:sink", got.ETag)).To(Succeed())
	rec.expectReq(http.MethodDelete, base+"/gid:sink")
})

var _ = It("updates an app agent with If-Match", func(ctx SpecContext) {
	var rec recorder
	admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		w.Header().Set("ETag", `"a2"`)
		_, _ = io.WriteString(w, `{"id":"gid:agent","update_time":"t1"}`)
	})

	name := "renamed"
	res, err := admin.AppAgents().Update(ctx, "gid:agent", `"a1"`, &config.UpdateAppAgent{
		DisplayName:    &name,
		APIPermissions: []string{config.PermissionAuthorization, config.PermissionCapture},
	})
	Expect(err).NotTo(HaveOccurred())
	rec.expectReq(http.MethodPut, "/configs/v1/application-agents/gid:agent")
	Expect(rec.ifMatch).To(Equal(`"a1"`))
	Expect(rec.body).To(HaveKeyWithValue("display_name", "renamed"))
	Expect(res.ETag).To(Equal(`"a2"`))
})

var _ = Describe("credential lists", func() {
	var rec recorder
	var admin *config.AdminClient

	BeforeEach(func() {
		admin = adminTo(func(w http.ResponseWriter, r *http.Request) {
			rec.capture(r)
			_, _ = io.WriteString(w, `{"data":[{"id":"gid:cred","kid":"k1"}]}`)
		})
	})

	It("lists service account credentials of an organization", func(ctx SpecContext) {
		saCreds, err := admin.ServiceAccountCredentials().List(ctx, "gid:org")
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodGet, "/configs/v1/service-account-credentials")
		Expect(rec.query.Get("organization_id")).To(Equal("gid:org"))
		Expect(saCreds).To(HaveLen(1))
		Expect(saCreds[0].Kid).To(Equal("k1"))
	})

	It("lists application agent credentials of a project", func(ctx SpecContext) {
		agentCreds, err := admin.AppAgentCredentials().List(ctx, "gid:proj", config.WithSearch("ci"))
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodGet, "/configs/v1/application-agent-credentials")
		Expect(rec.query).To(Equal(url.Values{"project_id": {"gid:proj"}, "search": {"ci"}}))
		Expect(agentCreds).To(HaveLen(1))
		Expect(agentCreds[0].ID).To(Equal("gid:cred"))
	})
})

var _ = It("starts a data schema rebuild", func(ctx SpecContext) {
	var rec recorder
	admin := adminTo(func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"status":"Rebuilding..."}`)
	})

	res, err := admin.DataSchema().Rebuild(ctx, "gid:proj")
	Expect(err).NotTo(HaveOccurred())
	rec.expectReq(http.MethodPost, "/configs/v1/data-schema/rebuild")
	Expect(rec.body).To(HaveKeyWithValue("project_id", "gid:proj"))
	Expect(res.Status).To(Equal("Rebuilding..."))
})

var _ = Describe("Token Introspect validation modes", func() {
	const base = "/configs/v1/token-introspects"

	It("sends offline validation with RFC 8414 discovery and reads it back", func(ctx SpecContext) {
		var rec recorder
		admin := adminTo(crudHandler(&rec, base, crudResponses{
			created: `{"id":"gid:ti"}`,
			read: `{"id":"gid:ti","ikg_node_type":"Person",
				"jwt_matcher":{"issuer":"https://issuer.example.com/tenant","audience":"api"},
				"offline_validation":{"public_jwks":[],"discovery":"rfc8414"}}`,
		}))

		_, err := admin.TokenIntrospects().Create(ctx, &config.CreateTokenIntrospect{
			ProjectID: "gid:proj", Name: "rfc8414",
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType: "Person",
				JwtMatcher:  json.RawMessage(`{"issuer":"https://issuer.example.com/tenant","audience":"api"}`),
				OfflineValidation: json.RawMessage(
					`{"discovery":"` + config.TokenIntrospectDiscoveryRFC8414 + `"}`),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(HaveKeyWithValue("offline_validation", HaveKeyWithValue("discovery", "rfc8414")))
		Expect(rec.body).NotTo(HaveKey("online_validation"))
		Expect(rec.body).NotTo(HaveKey("online_oauth2_validation"))

		got, err := admin.TokenIntrospects().Read(ctx, "gid:ti")
		Expect(err).NotTo(HaveOccurred())
		Expect(got.OfflineValidation).To(MatchJSON(`{"public_jwks":[],"discovery":"rfc8414"}`))
		Expect(got.OnlineOAuth2Validation).To(BeEmpty())
	})

	It("sends OAuth 2.0 introspection validation and reads it back without the secret", func(ctx SpecContext) {
		var rec recorder
		admin := adminTo(crudHandler(&rec, base, crudResponses{
			created: `{"id":"gid:ti"}`,
			read: `{"id":"gid:ti","ikg_node_type":"Person","opaque_matcher":{"hint":"my.domain.com"},
				"online_oauth2_validation":{"introspection_endpoint":"https://idp.example.com/introspect",
					"secret_basic":{"client_id":"my-client"},"cache_ttl":300}}`,
			updated: `{"id":"gid:ti"}`,
		}))
		oauth2 := json.RawMessage(`{"introspection_endpoint":"https://idp.example.com/introspect",
			"secret_basic":{"client_id":"my-client","client_secret":"s3cret"},"cache_ttl":300}`)

		_, err := admin.TokenIntrospects().Create(ctx, &config.CreateTokenIntrospect{
			ProjectID: "gid:proj", Name: "oauth2",
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType:            "Person",
				OpaqueMatcher:          json.RawMessage(`{"hint":"my.domain.com"}`),
				OnlineOAuth2Validation: oauth2,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(HaveKeyWithValue("online_oauth2_validation", And(
			HaveKeyWithValue("introspection_endpoint", "https://idp.example.com/introspect"),
			HaveKeyWithValue("secret_basic", HaveKeyWithValue("client_secret", "s3cret")),
			HaveKeyWithValue("cache_ttl", BeNumerically("==", 300)),
		)))
		Expect(rec.body).NotTo(HaveKey("offline_validation"))
		Expect(rec.body).NotTo(HaveKey("online_validation"))

		got, err := admin.TokenIntrospects().Read(ctx, "gid:ti")
		Expect(err).NotTo(HaveOccurred())
		Expect(got.OnlineOAuth2Validation).To(MatchJSON(`{"introspection_endpoint":"https://idp.example.com/introspect",
			"secret_basic":{"client_id":"my-client"},"cache_ttl":300}`))

		// An update may omit client_secret to keep the stored one.
		_, err = admin.TokenIntrospects().Update(ctx, "gid:ti", got.ETag, &config.UpdateTokenIntrospect{
			TokenIntrospectConfig: config.TokenIntrospectConfig{
				IkgNodeType:            "Person",
				OpaqueMatcher:          json.RawMessage(`{"hint":"my.domain.com"}`),
				OnlineOAuth2Validation: got.OnlineOAuth2Validation,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		rec.expectReq(http.MethodPut, base+"/gid:ti")
		Expect(rec.body).To(HaveKeyWithValue("online_oauth2_validation",
			HaveKeyWithValue("secret_basic", Equal(map[string]any{"client_id": "my-client"}))))
	})
})

var _ = Describe("Project IKG schema version", func() {
	It("sends the schema version on create and reads it back", func(ctx SpecContext) {
		var rec recorder
		admin := adminTo(crudHandler(&rec, "/configs/v1/projects", crudResponses{
			created: `{"id":"gid:proj"}`,
			read:    `{"id":"gid:proj","name":"proj","ikg_schema_version":"v1"}`,
		}))

		_, err := admin.Projects().Create(ctx, &config.CreateProject{
			OrganizationID: "gid:org", Name: "proj", Region: "europe-west1",
			IkgSchemaVersion: config.IkgSchemaVersionV1,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(HaveKeyWithValue("ikg_schema_version", "v1"))

		got, err := admin.Projects().Read(ctx, "gid:proj")
		Expect(err).NotTo(HaveOccurred())
		Expect(got.IkgSchemaVersion).To(Equal(config.IkgSchemaVersionV1))
	})

	It("leaves the platform default when the version is unset", func(ctx SpecContext) {
		var rec recorder
		admin := adminTo(crudHandler(&rec, "/configs/v1/projects", genericCRUDResponses("gid:proj")))

		_, err := admin.Projects().Create(ctx, &config.CreateProject{
			OrganizationID: "gid:org", Name: "proj", Region: "europe-west1",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).NotTo(HaveKey("ikg_schema_version"))
	})

	It("never sends the schema version on update, which cannot change it", func(ctx SpecContext) {
		var rec recorder
		admin := adminTo(crudHandler(&rec, "/configs/v1/projects", genericCRUDResponses("gid:proj")))

		name := "renamed"
		_, err := admin.Projects().Update(ctx, "gid:proj", "", &config.UpdateProject{DisplayName: &name})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(Equal(map[string]any{"display_name": "renamed"}))
	})
})
