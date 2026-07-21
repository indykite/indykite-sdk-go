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

const pathExternalDataResolvers = "/configs/v1/external-data-resolvers"

// ExternalDataResolver fetches a property value from an external HTTP endpoint.
type ExternalDataResolver struct {
	Metadata
	URL                 string              `json:"url"`
	Method              string              `json:"method"`
	RequestContentType  string              `json:"request_content_type"`
	RequestPayload      string              `json:"request_payload"`
	ResponseContentType string              `json:"response_content_type"`
	ResponseSelector    string              `json:"response_selector"`
	Headers             map[string][]string `json:"headers,omitempty"`
	Versioned
}

// ExternalDataResolverConfig is the resolver-specific config shared by create/update.
type ExternalDataResolverConfig struct {
	Headers             map[string][]string `json:"headers,omitempty"`
	URL                 string              `json:"url"`
	Method              string              `json:"method"`
	RequestContentType  string              `json:"request_content_type"`
	RequestPayload      string              `json:"request_payload,omitempty"`
	ResponseContentType string              `json:"response_content_type"`
	ResponseSelector    string              `json:"response_selector"`
}

// CreateExternalDataResolver is the body to create an external data resolver.
type CreateExternalDataResolver struct {
	ExternalDataResolverConfig
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateExternalDataResolver is the body to update an external data resolver.
type UpdateExternalDataResolver struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	ExternalDataResolverConfig
}

// ExternalDataResolverAPI is the /configs/v1/external-data-resolvers sub-API.
type ExternalDataResolverAPI struct {
	t *transport.Client
}

// List returns the external data resolvers in a project.
func (a *ExternalDataResolverAPI) List(
	ctx context.Context,
	projectID string,
	opts ...ListOption,
) ([]ExternalDataResolver, error) {
	return listResource[ExternalDataResolver](ctx, a.t, pathExternalDataResolvers, projectListQuery(projectID, opts))
}

// Create creates an external data resolver.
func (a *ExternalDataResolverAPI) Create(ctx context.Context, req *CreateExternalDataResolver) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathExternalDataResolvers, req)
}

// Read fetches one external data resolver by gid (or by name with WithLocation).
func (a *ExternalDataResolverAPI) Read(
	ctx context.Context,
	id string,
	opts ...ReadOption,
) (*ExternalDataResolver, error) {
	return readResource[ExternalDataResolver](ctx, a.t, pathExternalDataResolvers, id, readOptsQuery(opts))
}

// Update updates an external data resolver, optionally guarded by an ETag.
func (a *ExternalDataResolverAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateExternalDataResolver,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathExternalDataResolvers, id), req, ifMatch(etag)...)
}

// Delete deletes an external data resolver, optionally guarded by an ETag.
func (a *ExternalDataResolverAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathExternalDataResolvers, id, etag)
}
