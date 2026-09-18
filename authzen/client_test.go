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
	"net/url"

	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

// nodeJSON matches the JSON encoding of a Node: exactly type and id.
func nodeJSON(nodeType, id string) types.GomegaMatcher {
	return gstruct.MatchAllKeys(gstruct.Keys{"type": Equal(nodeType), "id": Equal(id)})
}

// nodeTypeJSON matches the JSON encoding of a NodeType: type only, no id.
func nodeTypeJSON(nodeType string) types.GomegaMatcher {
	return gstruct.MatchAllKeys(gstruct.Keys{"type": Equal(nodeType)})
}

// request is what the fake server saw on the last call.
type request struct {
	body   map[string]any
	query  url.Values
	method string
	path   string
}

// newClient wires an authzen.Client to an httptest server that records the
// last request into rec and answers every call with reply; the server is torn
// down when the spec ends.
func newClient(reply string) (*authzen.Client, *request) {
	GinkgoHelper()
	rec := &request{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.path = r.URL.Path
		rec.query = r.URL.Query()
		rec.body = nil
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &rec.body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, reply)
	}))
	DeferCleanup(srv.Close)

	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	Expect(err).NotTo(HaveOccurred())
	return authzen.NewClient(tc), rec
}

// newErrorClient wires an authzen.Client to a server that answers every call
// with a 403 JSON error body.
func newErrorClient() *authzen.Client {
	GinkgoHelper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"message":"access denied","code":"PERMISSION_DENIED"}`)
	}))
	DeferCleanup(srv.Close)

	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	Expect(err).NotTo(HaveOccurred())
	return authzen.NewClient(tc)
}

var _ = Describe("Evaluate and Allowed", func() {
	It("posts the (subject, action, resource) triple with the optional context", func(ctx SpecContext) {
		c, rec := newClient(`{"decision":true}`)

		ok, err := c.Allowed(ctx,
			authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"),
			authzen.WithInputParams(map[string]any{"budget": 100}),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		Expect(rec.method).To(Equal(http.MethodPost))
		Expect(rec.path).To(Equal("/access/v1/evaluation"))
		Expect(rec.body).To(HaveKeyWithValue("subject", nodeJSON("Person", "ada")))
		Expect(rec.body).To(HaveKeyWithValue("resource", nodeJSON("Server", "gpu-7")))
		Expect(rec.body).To(HaveKeyWithValue("action", HaveKeyWithValue("name", "PROVISION")))
		Expect(rec.body).To(HaveKeyWithValue("context",
			HaveKeyWithValue("input_params", HaveKeyWithValue("budget", BeNumerically("==", 100)))))
	})

	It("omits the context when no options are given", func(ctx SpecContext) {
		c, rec := newClient(`{"decision":false}`)

		ok, err := c.Allowed(ctx, authzen.NewNode("Person", "x"), "READ", authzen.NewNode("Doc", "d1"))
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeFalse())
		Expect(rec.body).NotTo(HaveKey("context"))
	})

	It("sends policy tags and input params together", func(ctx SpecContext) {
		c, rec := newClient(`{"decision":true}`)

		_, err := c.Allowed(ctx,
			authzen.NewNode("Person", "ada"), "READ", authzen.NewNode("Doc", "d1"),
			authzen.WithPolicyTags("prod", "eu"),
			authzen.WithInputParams(map[string]any{"tier": "gold"}),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.body).To(HaveKeyWithValue("context", gstruct.MatchAllKeys(gstruct.Keys{
			"policy_tags":  ConsistOf("prod", "eu"),
			"input_params": HaveKeyWithValue("tier", "gold"),
		})))
	})

	It("returns the full response, including the reasoning context", func(ctx SpecContext) {
		c, _ := newClient(`{"decision":false,"context":{"reason":"no OWNS edge"}}`)

		subject := authzen.NewNode("Person", "ada")
		resource := authzen.NewNode("Doc", "d1")
		resp, err := c.Evaluate(ctx, authzen.EvaluationRequest{
			Subject: &subject, Resource: &resource, Action: &authzen.Action{Name: "READ"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Decision).To(BeFalse())
		Expect(resp.Context).To(HaveValue(Equal(authzen.ResponseContext{Reason: "no OWNS edge"})))
	})
})

var _ = Describe("EvaluateBatch", func() {
	It("returns one decision per entry in request order", func(ctx SpecContext) {
		c, rec := newClient(`{"evaluations":[{"decision":true},{"decision":false}]}`)

		resp, err := c.EvaluateBatch(ctx, authzen.EvaluationsRequest{
			Subject: &authzen.Node{Type: "Person", ID: "ada"}, // default for all entries
			Evaluations: []authzen.EvaluationItem{
				{Action: &authzen.Action{Name: "READ"}, Resource: &authzen.Node{Type: "Doc", ID: "d1"}},
				{Action: &authzen.Action{Name: "DELETE"}, Resource: &authzen.Node{Type: "Doc", ID: "d2"}},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/access/v1/evaluations"))
		Expect(resp.Evaluations).To(HaveExactElements(
			authzen.EvaluationResponse{Decision: true},
			authzen.EvaluationResponse{Decision: false},
		))
	})

	It("sends batch-level defaults and leaves them out of the entries", func(ctx SpecContext) {
		c, rec := newClient(`{"evaluations":[{"decision":true}]}`)

		_, err := c.EvaluateBatch(ctx, authzen.EvaluationsRequest{
			Subject: &authzen.Node{Type: "Person", ID: "ada"},
			Context: &authzen.Context{PolicyTags: []string{"prod"}},
			Evaluations: []authzen.EvaluationItem{
				{Action: &authzen.Action{Name: "READ"}, Resource: &authzen.Node{Type: "Doc", ID: "d1"}},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(rec.body).To(HaveKeyWithValue("subject", nodeJSON("Person", "ada")))
		Expect(rec.body).To(HaveKeyWithValue("context", HaveKeyWithValue("policy_tags", ConsistOf("prod"))))
		Expect(rec.body).To(HaveKeyWithValue("evaluations", HaveExactElements(
			gstruct.MatchAllKeys(gstruct.Keys{
				"action":   HaveKeyWithValue("name", "READ"),
				"resource": nodeJSON("Doc", "d1"),
			}),
		)), "entry subject should be omitted so the batch default applies")
	})
})

var _ = Describe("Search", func() {
	It("SearchAction lists the actions a subject may perform on a resource", func(ctx SpecContext) {
		c, rec := newClient(`{"results":[{"name":"READ"},{"name":"PROVISION"}]}`)

		actions, err := c.SearchAction(ctx, authzen.SearchActionRequest{
			Subject:  &authzen.Node{Type: "Person", ID: "ada"},
			Resource: &authzen.Node{Type: "Server", ID: "gpu-7"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/access/v1/search/action"))
		Expect(actions).To(HaveExactElements(authzen.Action{Name: "READ"}, authzen.Action{Name: "PROVISION"}))
	})

	It("SearchResource lists the resources of a type the subject may act on", func(ctx SpecContext) {
		c, rec := newClient(`{"results":[{"type":"Server","id":"gpu-7"},{"type":"Server","id":"gpu-8"}]}`)

		nodes, err := c.SearchResource(ctx, authzen.SearchResourceRequest{
			Subject:  &authzen.Node{Type: "Person", ID: "ada"},
			Action:   &authzen.Action{Name: "PROVISION"},
			Resource: &authzen.NodeType{Type: "Server"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/access/v1/search/resource"))
		Expect(rec.body).To(HaveKeyWithValue("resource", nodeTypeJSON("Server")),
			"resource id should be absent for a type-only search resource")
		Expect(nodes).To(HaveExactElements(
			authzen.Node{Type: "Server", ID: "gpu-7"},
			authzen.Node{Type: "Server", ID: "gpu-8"},
		))
	})

	It("SearchSubject lists the subjects of a type allowed to act on a resource", func(ctx SpecContext) {
		c, rec := newClient(`{"results":[{"type":"Person","id":"ada"},{"type":"Person","id":"linus"}]}`)

		nodes, err := c.SearchSubject(ctx, authzen.SearchSubjectRequest{
			Subject:  &authzen.NodeType{Type: "Person"},
			Action:   &authzen.Action{Name: "PROVISION"},
			Resource: &authzen.Node{Type: "Server", ID: "gpu-7"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/access/v1/search/subject"))
		Expect(rec.body).To(HaveKeyWithValue("subject", nodeTypeJSON("Person")),
			"subject id should be absent for a type-only search subject")
		Expect(rec.body).To(HaveKeyWithValue("resource", nodeJSON("Server", "gpu-7")))
		Expect(nodes).To(HaveExactElements(
			authzen.Node{Type: "Person", ID: "ada"},
			authzen.Node{Type: "Person", ID: "linus"},
		))
	})
})

var _ = Describe("ListPolicies", func() {
	It("issues a body-less GET and returns the stored policies verbatim", func(ctx SpecContext) {
		c, rec := newClient(`{"results":[
			{"policy":{"meta":{"policy_version":"2.0-kbac"},"subject":{"type":"Person"}},"tags":["prod","eu"]},
			{"policy":{"meta":{"policy_version":"1.0-ciq"}},"tags":[]}
		]}`)

		policies, err := c.ListPolicies(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(rec.method).To(Equal(http.MethodGet))
		Expect(rec.path).To(Equal("/access/v1/policies"))
		Expect(rec.query).To(BeEmpty())
		Expect(rec.body).To(BeNil(), "GET must have no body")

		Expect(policies).To(HaveLen(2))
		Expect(policies[0].Policy).To(MatchJSON(`{"meta":{"policy_version":"2.0-kbac"},"subject":{"type":"Person"}}`))
		Expect(policies[0].Tags).To(Equal([]string{"prod", "eu"}))
		Expect(policies[1].Policy).To(MatchJSON(`{"meta":{"policy_version":"1.0-ciq"}}`))
		Expect(policies[1].Tags).To(And(Not(BeNil()), BeEmpty()), "a policy without tags carries an empty slice")
	})

	It("normalizes omitted or null tags to an empty slice", func(ctx SpecContext) {
		c, _ := newClient(`{"results":[
			{"policy":{"meta":{"policy_version":"2.0-kbac"}}},
			{"policy":{"meta":{"policy_version":"1.0-ciq"}},"tags":null}
		]}`)

		policies, err := c.ListPolicies(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(policies).To(HaveLen(2))
		for i, p := range policies {
			Expect(p.Tags).To(And(Not(BeNil()), BeEmpty()), "policy[%d] tags", i)
		}
	})

	It("narrows by subject type through the subject_type query parameter", func(ctx SpecContext) {
		c, rec := newClient(`{"results":[]}`)

		policies, err := c.ListPolicies(ctx, authzen.WithSubjectType("Person"))
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.path).To(Equal("/access/v1/policies"))
		Expect(rec.query).To(HaveKeyWithValue("subject_type", ConsistOf("Person")))
		Expect(policies).To(BeEmpty())
	})
})

// A non-2xx response surfaces as *transport.APIError from every method.
var _ = Describe("Error propagation", func() {
	subject := authzen.Node{Type: "Person", ID: "ada"}
	resource := authzen.Node{Type: "Server", ID: "gpu-7"}
	action := authzen.Action{Name: "PROVISION"}

	DescribeTable("returns the platform error as an APIError",
		func(ctx SpecContext, call func(context.Context, *authzen.Client) error) {
			err := call(ctx, newErrorClient())
			apiErr, ok := transport.AsAPIError(err)
			Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
			Expect(apiErr.StatusCode).To(Equal(http.StatusForbidden))
			Expect(apiErr.IsUnauthorized()).To(BeTrue())
			Expect(apiErr.Message).To(Equal("access denied"))
			Expect(apiErr.Code).To(Equal("PERMISSION_DENIED"))
		},
		Entry("Evaluate", func(ctx context.Context, c *authzen.Client) error {
			_, err := c.Evaluate(ctx, authzen.EvaluationRequest{
				Subject: &subject, Resource: &resource, Action: &action,
			})
			return err
		}),
		Entry("Allowed", func(ctx context.Context, c *authzen.Client) error {
			ok, err := c.Allowed(ctx, subject, action.Name, resource)
			Expect(ok).To(BeFalse(), "Allowed should be false on error")
			return err
		}),
		Entry("EvaluateBatch", func(ctx context.Context, c *authzen.Client) error {
			_, err := c.EvaluateBatch(ctx, authzen.EvaluationsRequest{
				Evaluations: []authzen.EvaluationItem{{Subject: &subject, Resource: &resource, Action: &action}},
			})
			return err
		}),
		Entry("SearchAction", func(ctx context.Context, c *authzen.Client) error {
			_, err := c.SearchAction(ctx, authzen.SearchActionRequest{Subject: &subject, Resource: &resource})
			return err
		}),
		Entry("SearchResource", func(ctx context.Context, c *authzen.Client) error {
			_, err := c.SearchResource(ctx, authzen.SearchResourceRequest{
				Subject: &subject, Action: &action, Resource: &authzen.NodeType{Type: "Server"},
			})
			return err
		}),
		Entry("SearchSubject", func(ctx context.Context, c *authzen.Client) error {
			_, err := c.SearchSubject(ctx, authzen.SearchSubjectRequest{
				Subject: &authzen.NodeType{Type: "Person"}, Action: &action, Resource: &resource,
			})
			return err
		}),
		Entry("ListPolicies", func(ctx context.Context, c *authzen.Client) error {
			_, err := c.ListPolicies(ctx)
			return err
		}),
	)
})
