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

//nolint:gosec // G101: URL path, not a credential.
const pathTokenIntrospects = "/configs/v1/token-introspects"

// TokenIntrospect configures how bearer tokens are validated (project scoped).
// The matcher/validation sub-objects are passed through as raw JSON; see the
// platform API reference for their structure.
type TokenIntrospect struct {
	Metadata
	IkgNodeType string `json:"ikg_node_type"`
	Versioned
	JwtMatcher        json.RawMessage `json:"jwt_matcher,omitempty"`
	OpaqueMatcher     json.RawMessage `json:"opaque_matcher,omitempty"`
	OfflineValidation json.RawMessage `json:"offline_validation,omitempty"`
	OnlineValidation  json.RawMessage `json:"online_validation,omitempty"`
	ClaimsMapping     json.RawMessage `json:"claims_mapping,omitempty"`
	SubClaim          json.RawMessage `json:"sub_claim,omitempty"`
	PerformUpsert     bool            `json:"perform_upsert"`
}

// TokenIntrospectConfig is the introspect-specific config shared by create/update.
type TokenIntrospectConfig struct {
	IkgNodeType       string          `json:"ikg_node_type"`
	JwtMatcher        json.RawMessage `json:"jwt_matcher,omitempty"`
	OpaqueMatcher     json.RawMessage `json:"opaque_matcher,omitempty"`
	OfflineValidation json.RawMessage `json:"offline_validation,omitempty"`
	OnlineValidation  json.RawMessage `json:"online_validation,omitempty"`
	ClaimsMapping     json.RawMessage `json:"claims_mapping,omitempty"`
	SubClaim          json.RawMessage `json:"sub_claim,omitempty"`
	PerformUpsert     bool            `json:"perform_upsert"`
}

// CreateTokenIntrospect is the body to create a token introspect config.
type CreateTokenIntrospect struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	TokenIntrospectConfig
}

// UpdateTokenIntrospect is the body to update a token introspect config.
type UpdateTokenIntrospect struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	TokenIntrospectConfig
}

// TokenIntrospectAPI is the /configs/v1/token-introspects sub-API.
type TokenIntrospectAPI struct {
	t *transport.Client
}

// List returns the token introspect configs in a project.
func (a *TokenIntrospectAPI) List(
	ctx context.Context,
	projectID string,
	opts ...ListOption,
) ([]TokenIntrospect, error) {
	return listResource[TokenIntrospect](ctx, a.t, pathTokenIntrospects, projectListQuery(projectID, opts))
}

// Create creates a token introspect config.
func (a *TokenIntrospectAPI) Create(ctx context.Context, req *CreateTokenIntrospect) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathTokenIntrospects, req)
}

// Read fetches one token introspect config by gid (or by name with WithLocation).
func (a *TokenIntrospectAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*TokenIntrospect, error) {
	return readResource[TokenIntrospect](ctx, a.t, pathTokenIntrospects, id, readOptsQuery(opts))
}

// Update updates a token introspect config, optionally guarded by an ETag.
func (a *TokenIntrospectAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateTokenIntrospect,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathTokenIntrospects, id), req, ifMatch(etag)...)
}

// Delete deletes a token introspect config, optionally guarded by an ETag.
func (a *TokenIntrospectAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathTokenIntrospects, id, etag)
}
