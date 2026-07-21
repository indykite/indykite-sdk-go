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

// Package authzen is the client for the IndyKite AuthZEN authorization API
// (/access/v1/*). It is a thin, ergonomic facade over a *transport.Client and
// runs on the runtime plane (App Agent token).
//
//	az := authzen.NewClient(client)
//	ok, _ := az.Allowed(ctx,
//	    authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
package authzen

import (
	"context"
	"net/http"

	"github.com/indykite/indykite-sdk-go/transport"
)

const (
	pathEvaluation     = "/access/v1/evaluation"
	pathEvaluations    = "/access/v1/evaluations"
	pathSearchAction   = "/access/v1/search/action"
	pathSearchResource = "/access/v1/search/resource"
	pathSearchSubject  = "/access/v1/search/subject"
)

// Client calls the AuthZEN API.
type Client struct {
	t *transport.Client
}

// NewClient builds an AuthZEN client over the shared transport.
func NewClient(t *transport.Client) *Client {
	return &Client{t: t}
}

// NewNode is a convenience constructor for a subject/resource node.
func NewNode(nodeType, id string) Node { return Node{Type: nodeType, ID: id} }

// Option customizes the optional Context of a convenience call.
type Option func(*Context)

// WithInputParams sets policy input parameters.
func WithInputParams(params map[string]any) Option {
	return func(c *Context) { c.InputParams = params }
}

// WithPolicyTags limits evaluation to policies carrying the given tags.
func WithPolicyTags(tags ...string) Option {
	return func(c *Context) { c.PolicyTags = tags }
}

func buildContext(opts ...Option) *Context {
	if len(opts) == 0 {
		return nil
	}
	c := &Context{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Evaluate makes a single authorization decision.
func (c *Client) Evaluate(ctx context.Context, req EvaluationRequest) (*EvaluationResponse, error) {
	var out EvaluationResponse
	if err := c.t.Do(ctx, http.MethodPost, pathEvaluation, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Allowed is the common case: a boolean decision for one (subject, action,
// resource) triple.
func (c *Client) Allowed(
	ctx context.Context,
	subject Node,
	action string,
	resource Node,
	opts ...Option,
) (bool, error) {
	resp, err := c.Evaluate(ctx, EvaluationRequest{
		Subject:  &subject,
		Resource: &resource,
		Action:   &Action{Name: action},
		Context:  buildContext(opts...),
	})
	if err != nil {
		return false, err
	}
	return resp.Decision, nil
}

// EvaluateBatch makes many decisions in one call, returning one decision per
// entry in request order.
func (c *Client) EvaluateBatch(ctx context.Context, req EvaluationsRequest) (*EvaluationsResponse, error) {
	var out EvaluationsResponse
	if err := c.t.Do(ctx, http.MethodPost, pathEvaluations, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchAction lists the actions a subject may perform on a resource.
func (c *Client) SearchAction(ctx context.Context, req SearchActionRequest) ([]Action, error) {
	var out searchActionResponse
	if err := c.t.Do(ctx, http.MethodPost, pathSearchAction, req, &out); err != nil {
		return nil, err
	}
	return out.Results, nil
}

// SearchResource lists the resources of a type a subject may act on.
func (c *Client) SearchResource(ctx context.Context, req SearchResourceRequest) ([]Node, error) {
	var out searchNodeResponse
	if err := c.t.Do(ctx, http.MethodPost, pathSearchResource, req, &out); err != nil {
		return nil, err
	}
	return out.Results, nil
}

// SearchSubject lists the subjects of a type allowed to perform an action on a
// resource.
func (c *Client) SearchSubject(ctx context.Context, req SearchSubjectRequest) ([]Node, error) {
	var out searchNodeResponse
	if err := c.t.Do(ctx, http.MethodPost, pathSearchSubject, req, &out); err != nil {
		return nil, err
	}
	return out.Results, nil
}
