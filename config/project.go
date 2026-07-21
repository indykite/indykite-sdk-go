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

const pathProjects = "/configs/v1/projects"

// Project is an organization-scoped project (formerly application space).
// DBConnection is passed through as raw JSON; see the platform API reference.
type Project struct {
	Metadata
	Region        string `json:"region,omitempty"`
	IkgSize       string `json:"ikg_size,omitempty"`
	IkgStatus     string `json:"ikg_status,omitempty"`
	ReplicaRegion string `json:"replica_region,omitempty"`
	Versioned
	DBConnection json.RawMessage `json:"db_connection,omitempty"`
}

// CreateProject is the body to create a project.
type CreateProject struct {
	OrganizationID string          `json:"organization_id"`
	Name           string          `json:"name"`
	DisplayName    string          `json:"display_name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Region         string          `json:"region"`
	IkgSize        string          `json:"ikg_size,omitempty"`
	ReplicaRegion  string          `json:"replica_region,omitempty"`
	DBConnection   json.RawMessage `json:"db_connection,omitempty"`
}

// UpdateProject is the body to update a project.
type UpdateProject struct {
	DisplayName  *string         `json:"display_name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	DBConnection json.RawMessage `json:"db_connection,omitempty"`
}

// ProjectAPI is the /configs/v1/projects sub-API.
type ProjectAPI struct {
	t *transport.Client
}

// List returns the projects in an organization.
func (a *ProjectAPI) List(ctx context.Context, orgID string, opts ...ListOption) ([]Project, error) {
	return listResource[Project](ctx, a.t, pathProjects, orgListQuery(orgID, opts))
}

// Create creates a project.
func (a *ProjectAPI) Create(ctx context.Context, req *CreateProject) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathProjects, req)
}

// Read fetches one project by gid (or by name with WithLocation).
func (a *ProjectAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*Project, error) {
	return readResource[Project](ctx, a.t, pathProjects, id, readOptsQuery(opts))
}

// Update updates a project, optionally guarded by an ETag.
func (a *ProjectAPI) Update(ctx context.Context, id, etag string, req *UpdateProject) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathProjects, id), req, ifMatch(etag)...)
}

// Delete deletes a project, optionally guarded by an ETag.
func (a *ProjectAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathProjects, id, etag)
}
