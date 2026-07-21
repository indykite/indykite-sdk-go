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
	"encoding/json"
	"net/http"

	"github.com/indykite/indykite-sdk-go/transport"
)

const pathEntityMatchingPipelines = "/configs/v1/entity-matching-pipelines"

// EntityMatchingPipeline configures a Metis entity-matching pipeline (project
// scoped). NodeFilter is passed through as raw JSON; see the platform API ref.
type EntityMatchingPipeline struct {
	Metadata
	PropertyMappingStatus  string `json:"property_mapping_status,omitempty"`
	PropertyMappingMessage string `json:"property_mapping_message,omitempty"`
	EntityMatchingStatus   string `json:"entity_matching_status,omitempty"`
	EntityMatchingMessage  string `json:"entity_matching_message,omitempty"`
	RerunInterval          string `json:"rerun_interval,omitempty"`
	LastRunTime            string `json:"last_run_time,omitempty"`
	ReportURL              string `json:"report_url,omitempty"`
	ReportType             string `json:"report_type,omitempty"`
	Versioned
	NodeFilter json.RawMessage `json:"node_filter,omitempty"`
	// PropertyMappings are the property pairings the pipeline currently uses.
	PropertyMappings      json.RawMessage `json:"property_mappings,omitempty"`
	MatchedEntities       int64           `json:"matched_entities,omitempty"`
	SimilarityScoreCutoff float32         `json:"similarity_score_cutoff"`
}

// CreateEntityMatchingPipeline is the body to create an entity-matching pipeline.
type CreateEntityMatchingPipeline struct {
	ProjectID             string          `json:"project_id"`
	Name                  string          `json:"name"`
	DisplayName           string          `json:"display_name,omitempty"`
	Description           string          `json:"description,omitempty"`
	RerunInterval         string          `json:"rerun_interval,omitempty"`
	NodeFilter            json.RawMessage `json:"node_filter"`
	SimilarityScoreCutoff float32         `json:"similarity_score_cutoff"`
}

// UpdateEntityMatchingPipeline is the body to update an entity-matching pipeline.
type UpdateEntityMatchingPipeline struct {
	DisplayName           *string `json:"display_name,omitempty"`
	Description           *string `json:"description,omitempty"`
	RerunInterval         string  `json:"rerun_interval,omitempty"`
	SimilarityScoreCutoff float32 `json:"similarity_score_cutoff"`
}

// EntityMatchingPipelineAPI is the /configs/v1/entity-matching-pipelines sub-API.
type EntityMatchingPipelineAPI struct {
	t *transport.Client
}

// List returns the entity-matching pipelines in a project.
func (a *EntityMatchingPipelineAPI) List(
	ctx context.Context,
	projectID string,
	opts ...ListOption,
) ([]EntityMatchingPipeline, error) {
	q := projectListQuery(projectID, opts)
	return listResource[EntityMatchingPipeline](ctx, a.t, pathEntityMatchingPipelines, q)
}

// Create creates an entity-matching pipeline.
func (a *EntityMatchingPipelineAPI) Create(
	ctx context.Context,
	req *CreateEntityMatchingPipeline,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathEntityMatchingPipelines, req)
}

// Read fetches one entity-matching pipeline by gid (or by name with WithLocation).
func (a *EntityMatchingPipelineAPI) Read(
	ctx context.Context,
	id string,
	opts ...ReadOption,
) (*EntityMatchingPipeline, error) {
	return readResource[EntityMatchingPipeline](ctx, a.t, pathEntityMatchingPipelines, id, readOptsQuery(opts))
}

// Update updates an entity-matching pipeline, optionally guarded by an ETag.
func (a *EntityMatchingPipelineAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateEntityMatchingPipeline,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathEntityMatchingPipelines, id), req, ifMatch(etag)...)
}

// Delete deletes an entity-matching pipeline, optionally guarded by an ETag.
func (a *EntityMatchingPipelineAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathEntityMatchingPipelines, id, etag)
}
