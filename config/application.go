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

const pathApplications = "/configs/v1/applications"

// Application is an application configuration (project/app-space scoped).
type Application struct {
	Metadata
	Versioned
}

// CreateApplication is the body to create an application.
type CreateApplication struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateApplication is the body to update an application.
type UpdateApplication struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ApplicationAPI is the /configs/v1/applications sub-API.
type ApplicationAPI struct {
	t *transport.Client
}

// List returns the applications in a project.
func (a *ApplicationAPI) List(ctx context.Context, projectID string, opts ...ListOption) ([]Application, error) {
	return listResource[Application](ctx, a.t, pathApplications, projectListQuery(projectID, opts))
}

// Create creates an application.
func (a *ApplicationAPI) Create(ctx context.Context, req *CreateApplication) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathApplications, req)
}

// Read fetches one application by gid (or by name with WithLocation).
func (a *ApplicationAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*Application, error) {
	return readResource[Application](ctx, a.t, pathApplications, id, readOptsQuery(opts))
}

// Update updates an application, optionally guarded by an ETag.
func (a *ApplicationAPI) Update(ctx context.Context, id, etag string, req *UpdateApplication) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathApplications, id), req, ifMatch(etag)...)
}

// Delete deletes an application, optionally guarded by an ETag.
func (a *ApplicationAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathApplications, id, etag)
}
