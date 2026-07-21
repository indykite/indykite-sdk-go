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

package authzen_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/transport"
)

type capture struct {
	path string
	body map[string]any
}

// server returns an httptest server that records the last request and replies
// with the given JSON, plus an authzen.Client wired to it.
func server(t *testing.T, reply string) (*authzen.Client, *capture) {
	t.Helper()
	rec := &capture{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &rec.body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, reply)
	}))
	t.Cleanup(srv.Close)

	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return authzen.NewClient(tc), rec
}

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

func TestAllowed(t *testing.T) {
	c, rec := server(t, `{"decision":true}`)

	ok, err := c.Allowed(context.Background(),
		authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"),
		authzen.WithInputParams(map[string]any{"budget": 100}),
	)
	if err != nil {
		t.Fatalf("Allowed: %v", err)
	}
	if !ok {
		t.Error("decision = false, want true")
	}
	if rec.path != "/access/v1/evaluation" {
		t.Errorf("path = %q", rec.path)
	}
	// Verify the request body was shaped correctly.
	subj := rec.body["subject"].(map[string]any)
	if subj["type"] != "Person" || subj["id"] != "ada" {
		t.Errorf("subject = %v", subj)
	}
	act := rec.body["action"].(map[string]any)
	if act["name"] != "PROVISION" {
		t.Errorf("action = %v", act)
	}
	ctxObj := rec.body["context"].(map[string]any)
	ip := ctxObj["input_params"].(map[string]any)
	if ip["budget"].(float64) != 100 {
		t.Errorf("input_params = %v", ip)
	}
}

func TestAllowedNoContextOmitsField(t *testing.T) {
	c, rec := server(t, `{"decision":false}`)
	if _, err := c.Allowed(context.Background(),
		authzen.NewNode("Person", "x"), "READ", authzen.NewNode("Doc", "d1")); err != nil {
		t.Fatalf("Allowed: %v", err)
	}
	if _, present := rec.body["context"]; present {
		t.Error("context should be omitted when no options are given")
	}
}

func TestEvaluateBatch(t *testing.T) {
	c, rec := server(t, `{"evaluations":[{"decision":true},{"decision":false}]}`)

	resp, err := c.EvaluateBatch(context.Background(), authzen.EvaluationsRequest{
		Subject: &authzen.Node{Type: "Person", ID: "ada"}, // default for all entries
		Evaluations: []authzen.EvaluationItem{
			{Action: &authzen.Action{Name: "READ"}, Resource: &authzen.Node{Type: "Doc", ID: "d1"}},
			{Action: &authzen.Action{Name: "DELETE"}, Resource: &authzen.Node{Type: "Doc", ID: "d2"}},
		},
	})
	if err != nil {
		t.Fatalf("EvaluateBatch: %v", err)
	}
	if rec.path != "/access/v1/evaluations" {
		t.Errorf("path = %q", rec.path)
	}
	if len(resp.Evaluations) != 2 || !resp.Evaluations[0].Decision || resp.Evaluations[1].Decision {
		t.Errorf("decisions = %+v", resp.Evaluations)
	}
}

func TestSearchResource(t *testing.T) {
	c, rec := server(t, `{"results":[{"type":"Server","id":"gpu-7"},{"type":"Server","id":"gpu-8"}]}`)

	nodes, err := c.SearchResource(context.Background(), authzen.SearchResourceRequest{
		Subject:  &authzen.Node{Type: "Person", ID: "ada"},
		Action:   &authzen.Action{Name: "PROVISION"},
		Resource: &authzen.NodeType{Type: "Server"},
	})
	if err != nil {
		t.Fatalf("SearchResource: %v", err)
	}
	if rec.path != "/access/v1/search/resource" {
		t.Errorf("path = %q", rec.path)
	}
	if len(nodes) != 2 || nodes[0].ID != "gpu-7" || nodes[1].ID != "gpu-8" {
		t.Errorf("results = %+v", nodes)
	}
}

func TestEvaluateBatchBody(t *testing.T) {
	c, rec := server(t, `{"evaluations":[{"decision":true}]}`)

	_, err := c.EvaluateBatch(context.Background(), authzen.EvaluationsRequest{
		Subject: &authzen.Node{Type: "Person", ID: "ada"},
		Context: &authzen.Context{PolicyTags: []string{"prod"}},
		Evaluations: []authzen.EvaluationItem{
			{Action: &authzen.Action{Name: "READ"}, Resource: &authzen.Node{Type: "Doc", ID: "d1"}},
		},
	})
	if err != nil {
		t.Fatalf("EvaluateBatch: %v", err)
	}
	subj := rec.body["subject"].(map[string]any)
	if subj["type"] != "Person" || subj["id"] != "ada" {
		t.Errorf("subject = %v", subj)
	}
	tags := rec.body["context"].(map[string]any)["policy_tags"].([]any)
	if len(tags) != 1 || tags[0] != "prod" {
		t.Errorf("policy_tags = %v", tags)
	}
	evals := rec.body["evaluations"].([]any)
	e0 := evals[0].(map[string]any)
	if e0["action"].(map[string]any)["name"] != "READ" {
		t.Errorf("evaluations[0] = %v", e0)
	}
	if _, present := e0["subject"]; present {
		t.Error("entry subject should be omitted so the batch default applies")
	}
}

func TestSearchSubject(t *testing.T) {
	c, rec := server(t, `{"results":[{"type":"Person","id":"ada"},{"type":"Person","id":"linus"}]}`)

	nodes, err := c.SearchSubject(context.Background(), authzen.SearchSubjectRequest{
		Subject:  &authzen.NodeType{Type: "Person"},
		Action:   &authzen.Action{Name: "PROVISION"},
		Resource: &authzen.Node{Type: "Server", ID: "gpu-7"},
	})
	if err != nil {
		t.Fatalf("SearchSubject: %v", err)
	}
	if rec.path != "/access/v1/search/subject" {
		t.Errorf("path = %q", rec.path)
	}
	subj := rec.body["subject"].(map[string]any)
	if subj["type"] != "Person" {
		t.Errorf("subject = %v", subj)
	}
	if _, present := subj["id"]; present {
		t.Error("subject id should be absent for a type-only search subject")
	}
	res := rec.body["resource"].(map[string]any)
	if res["type"] != "Server" || res["id"] != "gpu-7" {
		t.Errorf("resource = %v", res)
	}
	if len(nodes) != 2 || nodes[0].ID != "ada" || nodes[1].ID != "linus" {
		t.Errorf("results = %+v", nodes)
	}
}

func TestWithPolicyTags(t *testing.T) {
	c, rec := server(t, `{"decision":true}`)

	if _, err := c.Allowed(context.Background(),
		authzen.NewNode("Person", "ada"), "READ", authzen.NewNode("Doc", "d1"),
		authzen.WithPolicyTags("prod", "eu"),
		authzen.WithInputParams(map[string]any{"tier": "gold"}),
	); err != nil {
		t.Fatalf("Allowed: %v", err)
	}
	ctxObj := rec.body["context"].(map[string]any)
	tags := ctxObj["policy_tags"].([]any)
	if len(tags) != 2 || tags[0] != "prod" || tags[1] != "eu" {
		t.Errorf("policy_tags = %v", tags)
	}
	if ctxObj["input_params"].(map[string]any)["tier"] != "gold" {
		t.Errorf("input_params = %v", ctxObj["input_params"])
	}
}

// errorServer replies to every request with a 403 JSON error body.
func errorServer(t *testing.T) *authzen.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"message":"access denied","code":"PERMISSION_DENIED"}`)
	}))
	t.Cleanup(srv.Close)

	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return authzen.NewClient(tc)
}

// A non-2xx response surfaces as *transport.APIError from every method.
func TestErrorPropagation(t *testing.T) {
	c := errorServer(t)
	ctx := context.Background()
	subject := authzen.Node{Type: "Person", ID: "ada"}
	resource := authzen.Node{Type: "Server", ID: "gpu-7"}
	action := authzen.Action{Name: "PROVISION"}

	calls := map[string]func() error{
		"Evaluate": func() error {
			_, err := c.Evaluate(ctx, authzen.EvaluationRequest{
				Subject: &subject, Resource: &resource, Action: &action,
			})
			return err
		},
		"Allowed": func() error {
			ok, err := c.Allowed(ctx, subject, action.Name, resource)
			if ok {
				t.Error("Allowed should be false on error")
			}
			return err
		},
		"EvaluateBatch": func() error {
			_, err := c.EvaluateBatch(ctx, authzen.EvaluationsRequest{
				Evaluations: []authzen.EvaluationItem{{Subject: &subject, Resource: &resource, Action: &action}},
			})
			return err
		},
		"SearchAction": func() error {
			_, err := c.SearchAction(ctx, authzen.SearchActionRequest{Subject: &subject, Resource: &resource})
			return err
		},
		"SearchResource": func() error {
			_, err := c.SearchResource(ctx, authzen.SearchResourceRequest{
				Subject: &subject, Action: &action, Resource: &authzen.NodeType{Type: "Server"},
			})
			return err
		},
		"SearchSubject": func() error {
			_, err := c.SearchSubject(ctx, authzen.SearchSubjectRequest{
				Subject: &authzen.NodeType{Type: "Person"}, Action: &action, Resource: &resource,
			})
			return err
		},
	}
	for name, call := range calls {
		err := call()
		if err == nil {
			t.Errorf("%s: expected error", name)
			continue
		}
		apiErr, ok := transport.AsAPIError(err)
		if !ok {
			t.Errorf("%s: err = %v, want *transport.APIError", name, err)
			continue
		}
		if apiErr.StatusCode != http.StatusForbidden || !apiErr.IsUnauthorized() {
			t.Errorf("%s: status = %d, want 403", name, apiErr.StatusCode)
		}
		if apiErr.Message != "access denied" || apiErr.Code != "PERMISSION_DENIED" {
			t.Errorf("%s: message/code = %q/%q", name, apiErr.Message, apiErr.Code)
		}
	}
}

func TestSearchAction(t *testing.T) {
	c, rec := server(t, `{"results":[{"name":"READ"},{"name":"PROVISION"}]}`)
	actions, err := c.SearchAction(context.Background(), authzen.SearchActionRequest{
		Subject:  &authzen.Node{Type: "Person", ID: "ada"},
		Resource: &authzen.Node{Type: "Server", ID: "gpu-7"},
	})
	if err != nil {
		t.Fatalf("SearchAction: %v", err)
	}
	if rec.path != "/access/v1/search/action" {
		t.Errorf("path = %q", rec.path)
	}
	if len(actions) != 2 || actions[1].Name != "PROVISION" {
		t.Errorf("results = %+v", actions)
	}
}
