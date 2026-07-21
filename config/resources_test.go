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
)

// TestKnowledgeQueryCRUD exercises a standard generic resource end to end.
func TestKnowledgeQueryCRUD(t *testing.T) {
	const base = "/configs/v1/knowledge-queries"
	var rec recorder
	admin := adminTo(t, crudHandler(&rec, base, crudResponses{
		created: `{"id":"gid:kq"}`,
		list:    `{"data":[{"id":"gid:kq","name":"q","status":"ACTIVE"}]}`,
		read:    `{"id":"gid:kq","query":"{}","status":"ACTIVE","policy_id":"gid:p"}`,
		updated: `{"id":"gid:kq"}`,
	}))
	kq := admin.KnowledgeQueries()
	ctx := context.Background()

	if _, err := kq.List(ctx, "gid:proj"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.query.Get("project_id") != "gid:proj" {
		t.Errorf("list query = %v", rec.query)
	}

	created, err := kq.Create(ctx, &config.CreateKnowledgeQuery{
		ProjectID: "gid:proj", Name: "q", Query: "{}", Status: config.StatusActive, PolicyID: "gid:p",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	wantWrite(t, created, `"v1"`)

	got, err := kq.Read(ctx, "gid:kq")
	if err != nil || got.PolicyID != "gid:p" || got.ETag != `"v2"` {
		t.Fatalf("Read: %+v err=%v", got, err)
	}

	if _, err = kq.Update(ctx, "gid:kq", got.ETag, &config.UpdateKnowledgeQuery{
		Query: "{}", Status: config.StatusInactive, PolicyID: "gid:p",
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if rec.ifMatch != `"v2"` {
		t.Errorf("update If-Match = %q, want \"v2\"", rec.ifMatch)
	}

	if err = kq.Delete(ctx, "gid:kq", got.ETag); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rec.wantReq(t, http.MethodDelete, base+"/gid:kq")
}

func TestOrganizationReadCurrent(t *testing.T) {
	var path string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("ETag", `"o1"`)
		_, _ = io.WriteString(w, `{"id":"gid:org","name":"acme"}`)
	})
	org, err := admin.Organizations().ReadCurrent(context.Background())
	if err != nil {
		t.Fatalf("ReadCurrent: %v", err)
	}
	if path != "/configs/v1/organizations/current" {
		t.Errorf("path = %q", path)
	}
	if org.ID != "gid:org" || org.ETag != `"o1"` {
		t.Errorf("org = %+v", org)
	}
}

func TestServiceAccountCredentialCreateReturnsConfig(t *testing.T) {
	var body map[string]any
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		_, _ = io.WriteString(w, `{"id":"gid:cred","kid":"k1","service_account_config":{"appAgentId":"x"}}`)
	})
	res, err := admin.ServiceAccountCredentials().Create(context.Background(),
		&config.CreateServiceAccountCredential{ServiceAccountID: "gid:sa", DisplayName: "ci"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if body["service_account_id"] != "gid:sa" {
		t.Errorf("body = %v", body)
	}
	if res.Kid != "k1" || len(res.ServiceAccountConfig) == 0 {
		t.Errorf("result = %+v", res)
	}
}

func TestEventSinkCRUD(t *testing.T) {
	const base = "/configs/v1/event-sinks"
	var rec recorder
	admin := adminTo(t, crudHandler(&rec, base, crudResponses{
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
	ctx := context.Background()

	if _, err := sinks.List(ctx, "gid:proj"); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.query.Get("project_id") != "gid:proj" {
		t.Errorf("list query = %v", rec.query)
	}

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
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	wantWrite(t, created, `"v1"`)
	if _, ok := rec.body["providers"].(map[string]any)["kafka-1"]; !ok {
		t.Errorf("create body providers = %v", rec.body["providers"])
	}

	got, err := sinks.Read(ctx, "gid:sink")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.ETag != `"v2"` || !got.IncludeCDCEvents || got.Providers["kafka-1"].Kafka.Topic != "events" {
		t.Fatalf("Read: %+v", got)
	}

	if _, err = sinks.Update(ctx, "gid:sink", got.ETag, &config.UpdateEventSink{
		Providers: got.Providers, Routes: got.Routes,
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if rec.ifMatch != `"v2"` {
		t.Errorf("update If-Match = %q, want \"v2\"", rec.ifMatch)
	}

	if err = sinks.Delete(ctx, "gid:sink", got.ETag); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rec.wantReq(t, http.MethodDelete, base+"/gid:sink")
}

func TestAppAgentUpdate(t *testing.T) {
	var ifMatch, method, path string
	var body map[string]any
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		ifMatch = r.Header.Get("If-Match")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("ETag", `"a2"`)
		_, _ = io.WriteString(w, `{"id":"gid:agent","update_time":"t1"}`)
	})

	name := "renamed"
	res, err := admin.AppAgents().Update(context.Background(), "gid:agent", `"a1"`, &config.UpdateAppAgent{
		DisplayName:    &name,
		APIPermissions: []string{config.PermissionAuthorization, config.PermissionCapture},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if method != http.MethodPut || path != "/configs/v1/application-agents/gid:agent" {
		t.Errorf("method/path = %s %s", method, path)
	}
	if ifMatch != `"a1"` {
		t.Errorf("If-Match = %q", ifMatch)
	}
	if body["display_name"] != "renamed" {
		t.Errorf("body = %v", body)
	}
	if res.ETag != `"a2"` {
		t.Errorf("result = %+v", res)
	}
}

func TestCredentialLists(t *testing.T) {
	var query url.Values
	var path string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query()
		_, _ = io.WriteString(w, `{"data":[{"id":"gid:cred","kid":"k1"}]}`)
	})

	saCreds, err := admin.ServiceAccountCredentials().List(context.Background(), "gid:org")
	if err != nil {
		t.Fatalf("List service-account-credentials: %v", err)
	}
	if path != "/configs/v1/service-account-credentials" || query.Get("organization_id") != "gid:org" {
		t.Errorf("path=%q query=%v", path, query)
	}
	if len(saCreds) != 1 || saCreds[0].Kid != "k1" {
		t.Errorf("saCreds = %+v", saCreds)
	}

	agentCreds, err := admin.AppAgentCredentials().List(context.Background(), "gid:proj", config.WithSearch("ci"))
	if err != nil {
		t.Fatalf("List application-agent-credentials: %v", err)
	}
	if path != "/configs/v1/application-agent-credentials" ||
		query.Get("project_id") != "gid:proj" || query.Get("search") != "ci" {
		t.Errorf("path=%q query=%v", path, query)
	}
	if len(agentCreds) != 1 || agentCreds[0].ID != "gid:cred" {
		t.Errorf("agentCreds = %+v", agentCreds)
	}
}

func TestDataSchemaRebuild(t *testing.T) {
	var body map[string]any
	var method, path string
	admin := adminTo(t, func(w http.ResponseWriter, r *http.Request) {
		method, path = r.Method, r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"status":"Rebuilding..."}`)
	})
	res, err := admin.DataSchema().Rebuild(context.Background(), "gid:proj")
	if err != nil {
		t.Fatalf("Rebuild: %v", err)
	}
	if method != http.MethodPost || path != "/configs/v1/data-schema/rebuild" {
		t.Errorf("method/path = %s %s", method, path)
	}
	if body["project_id"] != "gid:proj" || res.Status != "Rebuilding..." {
		t.Errorf("body=%v res=%+v", body, res)
	}
}
