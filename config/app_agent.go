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

// The platform serves application agents at both paths; we use the canonical one.
const pathAppAgents = "/configs/v1/application-agents"

// AppAgentAPI is the /configs/v1/application-agents sub-API.
type AppAgentAPI struct {
	t *transport.Client
}

// List returns the application agents in a project.
func (a *AppAgentAPI) List(ctx context.Context, projectID string, opts ...ListOption) ([]AppAgent, error) {
	return listResource[AppAgent](ctx, a.t, pathAppAgents, projectListQuery(projectID, opts))
}

// Create creates an application agent.
func (a *AppAgentAPI) Create(ctx context.Context, req *CreateAppAgent) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathAppAgents, req)
}

// Read fetches one application agent by gid (or by name with WithLocation).
func (a *AppAgentAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*AppAgent, error) {
	return readResource[AppAgent](ctx, a.t, pathAppAgents, id, readOptsQuery(opts))
}

// Update updates an application agent, optionally guarded by an ETag.
func (a *AppAgentAPI) Update(ctx context.Context, id, etag string, req *UpdateAppAgent) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathAppAgents, id), req, ifMatch(etag)...)
}

// Delete deletes an application agent, optionally guarded by an ETag.
func (a *AppAgentAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathAppAgents, id, etag)
}
