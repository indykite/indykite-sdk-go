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

package config

import (
	"context"
	"net/http"

	"github.com/indykite/indykite-sdk-go/transport"
)

const pathKnowledgeQueries = "/configs/v1/knowledge-queries"

// KnowledgeQuery is a ContX IQ knowledge query configuration (project scoped).
type KnowledgeQuery struct {
	Metadata
	// Query is the knowledge query document as a JSON string.
	Query string `json:"query"`
	// Status is one of ACTIVE, INACTIVE, DRAFT.
	Status string `json:"status"`
	// PolicyID is the authorization policy gid that governs the query.
	PolicyID string `json:"policy_id"`
	Versioned
}

// CreateKnowledgeQuery is the body to create a knowledge query.
type CreateKnowledgeQuery struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	Query       string `json:"query"`
	Status      string `json:"status"`
	PolicyID    string `json:"policy_id"`
}

// UpdateKnowledgeQuery is the body to update a knowledge query.
type UpdateKnowledgeQuery struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	PolicyID    string  `json:"policy_id"`
}

// KnowledgeQueryAPI is the /configs/v1/knowledge-queries sub-API.
type KnowledgeQueryAPI struct {
	t *transport.Client
}

// List returns the knowledge queries in a project.
func (a *KnowledgeQueryAPI) List(ctx context.Context, projectID string, opts ...ListOption) ([]KnowledgeQuery, error) {
	return listResource[KnowledgeQuery](ctx, a.t, pathKnowledgeQueries, projectListQuery(projectID, opts))
}

// Create creates a knowledge query.
func (a *KnowledgeQueryAPI) Create(ctx context.Context, req *CreateKnowledgeQuery) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathKnowledgeQueries, req)
}

// Read fetches one knowledge query by gid (or by name with WithLocation).
func (a *KnowledgeQueryAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*KnowledgeQuery, error) {
	return readResource[KnowledgeQuery](ctx, a.t, pathKnowledgeQueries, id, readOptsQuery(opts))
}

// Update updates a knowledge query, optionally guarded by an ETag.
func (a *KnowledgeQueryAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateKnowledgeQuery,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathKnowledgeQueries, id), req, ifMatch(etag)...)
}

// Delete deletes a knowledge query, optionally guarded by an ETag.
func (a *KnowledgeQueryAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathKnowledgeQueries, id, etag)
}
