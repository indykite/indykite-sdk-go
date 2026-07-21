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

const pathMCPServers = "/configs/v1/mcp-servers"

// MCPServer configures an MCP server fronting a project (project scoped).
type MCPServer struct {
	Metadata
	AppAgentID        string `json:"app_agent_id"`
	TokenIntrospectID string `json:"token_introspect_id"`
	Versioned
	ScopesSupported []string `json:"scopes_supported"`
	Enabled         bool     `json:"enabled"`
}

// MCPServerConfig is the MCP-specific config shared by create/update.
type MCPServerConfig struct {
	Enabled           *bool    `json:"enabled"`
	AppAgentID        string   `json:"app_agent_id"`
	TokenIntrospectID string   `json:"token_introspect_id"`
	ScopesSupported   []string `json:"scopes_supported"`
}

// CreateMCPServer is the body to create an MCP server config.
type CreateMCPServer struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	MCPServerConfig
}

// UpdateMCPServer is the body to update an MCP server config.
type UpdateMCPServer struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	MCPServerConfig
}

// MCPServerAPI is the /configs/v1/mcp-servers sub-API.
type MCPServerAPI struct {
	t *transport.Client
}

// List returns the MCP server configs in a project.
func (a *MCPServerAPI) List(ctx context.Context, projectID string, opts ...ListOption) ([]MCPServer, error) {
	return listResource[MCPServer](ctx, a.t, pathMCPServers, projectListQuery(projectID, opts))
}

// Create creates an MCP server config.
func (a *MCPServerAPI) Create(ctx context.Context, req *CreateMCPServer) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathMCPServers, req)
}

// Read fetches one MCP server config by gid (or by name with WithLocation).
func (a *MCPServerAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*MCPServer, error) {
	return readResource[MCPServer](ctx, a.t, pathMCPServers, id, readOptsQuery(opts))
}

// Update updates an MCP server config, optionally guarded by an ETag.
func (a *MCPServerAPI) Update(ctx context.Context, id, etag string, req *UpdateMCPServer) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathMCPServers, id), req, ifMatch(etag)...)
}

// Delete deletes an MCP server config, optionally guarded by an ETag.
func (a *MCPServerAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathMCPServers, id, etag)
}
