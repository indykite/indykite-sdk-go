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

const pathServiceAccounts = "/configs/v1/service-accounts"

// Service account role values.
const (
	RoleAllEditor = "all_editor"
	RoleAllViewer = "all_viewer"
)

// ServiceAccount is an organization-scoped service account. The role is
// write-only: it is set on Create but never returned by read/list.
type ServiceAccount struct {
	Metadata
	Versioned
}

// CreateServiceAccount is the body to create a service account.
type CreateServiceAccount struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name,omitempty"`
	Description    string `json:"description,omitempty"`
	Role           string `json:"role"`
}

// UpdateServiceAccount is the body to update a service account.
type UpdateServiceAccount struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ServiceAccountAPI is the /configs/v1/service-accounts sub-API.
type ServiceAccountAPI struct {
	t *transport.Client
}

// List returns the service accounts in an organization.
func (a *ServiceAccountAPI) List(ctx context.Context, orgID string, opts ...ListOption) ([]ServiceAccount, error) {
	return listResource[ServiceAccount](ctx, a.t, pathServiceAccounts, orgListQuery(orgID, opts))
}

// Create creates a service account.
func (a *ServiceAccountAPI) Create(ctx context.Context, req *CreateServiceAccount) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathServiceAccounts, req)
}

// Read fetches one service account by gid (or by name with WithLocation).
func (a *ServiceAccountAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*ServiceAccount, error) {
	return readResource[ServiceAccount](ctx, a.t, pathServiceAccounts, id, readOptsQuery(opts))
}

// Update updates a service account, optionally guarded by an ETag.
func (a *ServiceAccountAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateServiceAccount,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathServiceAccounts, id), req, ifMatch(etag)...)
}

// Delete deletes a service account, optionally guarded by an ETag.
func (a *ServiceAccountAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathServiceAccounts, id, etag)
}
