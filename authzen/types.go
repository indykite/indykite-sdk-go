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

package authzen

// Node is a graph node used as a subject or resource (type + external id).
type Node struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// NodeType is a node identified by type only (used by the search endpoints).
type NodeType struct {
	Type string `json:"type"`
}

// Action is the action being asked about.
type Action struct {
	Name string `json:"name"`
}

// Context carries extra decision inputs.
type Context struct {
	// InputParams supplies values for input parameters referenced by policies.
	InputParams map[string]any `json:"input_params,omitempty"`
	// PolicyTags limits evaluation to policies carrying the given tags.
	PolicyTags []string `json:"policy_tags,omitempty"`
}

// ResponseContext is the optional reasoning returned with a decision.
type ResponseContext struct {
	Reason  string              `json:"reason,omitempty"`
	Advices []map[string]string `json:"advice,omitempty"`
}

// EvaluationRequest is the body for POST /access/v1/evaluation.
type EvaluationRequest struct {
	Subject  *Node    `json:"subject"`
	Resource *Node    `json:"resource"`
	Action   *Action  `json:"action"`
	Context  *Context `json:"context,omitempty"`
}

// EvaluationResponse is the result of a single evaluation.
type EvaluationResponse struct {
	Context  *ResponseContext `json:"context,omitempty"`
	Decision bool             `json:"decision"`
}

// EvaluationItem is one entry of a batch; any unset field falls back to the
// batch-level default of the same name.
type EvaluationItem struct {
	Subject  *Node    `json:"subject,omitempty"`
	Resource *Node    `json:"resource,omitempty"`
	Action   *Action  `json:"action,omitempty"`
	Context  *Context `json:"context,omitempty"`
}

// EvaluationsRequest is the body for POST /access/v1/evaluations. The top-level
// Subject/Resource/Action/Context act as defaults for entries that omit them.
type EvaluationsRequest struct {
	Subject     *Node            `json:"subject,omitempty"`
	Resource    *Node            `json:"resource,omitempty"`
	Action      *Action          `json:"action,omitempty"`
	Context     *Context         `json:"context,omitempty"`
	Evaluations []EvaluationItem `json:"evaluations"`
}

// EvaluationsResponse holds one decision per entry, in request order.
type EvaluationsResponse struct {
	Evaluations []EvaluationResponse `json:"evaluations"`
}

// SearchActionRequest is the body for POST /access/v1/search/action.
type SearchActionRequest struct {
	Subject  *Node    `json:"subject"`
	Resource *Node    `json:"resource"`
	Context  *Context `json:"context,omitempty"`
}

type searchActionResponse struct {
	Results []Action `json:"results"`
}

// SearchResourceRequest is the body for POST /access/v1/search/resource.
type SearchResourceRequest struct {
	Subject  *Node     `json:"subject"`
	Action   *Action   `json:"action"`
	Resource *NodeType `json:"resource"`
	Context  *Context  `json:"context,omitempty"`
}

// SearchSubjectRequest is the body for POST /access/v1/search/subject.
type SearchSubjectRequest struct {
	Subject  *NodeType `json:"subject"`
	Action   *Action   `json:"action"`
	Resource *Node     `json:"resource"`
	Context  *Context  `json:"context,omitempty"`
}

type searchNodeResponse struct {
	Results []Node `json:"results"`
}
