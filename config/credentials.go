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

const (
	//nolint:gosec // G101: URL path, not a credential.
	pathServiceAccountCredentials = "/configs/v1/service-account-credentials"
	//nolint:gosec // G101: URL path, not a credential.
	pathAppAgentCredentials = "/configs/v1/application-agent-credentials"
)

// Credentials are not config documents: they have no ETag and cannot be updated.
// Create returns the agent/account config (the signed credential) ONCE.

// CreateServiceAccountCredential is the body to mint a service account credential.
type CreateServiceAccountCredential struct {
	ServiceAccountID string `json:"service_account_id"`
	DisplayName      string `json:"display_name,omitempty"`
	ExpireTime       string `json:"expire_time,omitempty"`
}

// ServiceAccountCredentialResult is returned by Create. ServiceAccountConfig
// holds the signed credential JSON — capture it, it is not retrievable later.
type ServiceAccountCredentialResult struct {
	ID                   string          `json:"id"`
	ServiceAccountID     string          `json:"service_account_id"`
	Kid                  string          `json:"kid"`
	CreateTime           string          `json:"create_time"`
	CreatedBy            string          `json:"created_by"`
	ExpireTime           string          `json:"expire_time"`
	DisplayName          string          `json:"display_name"`
	ServiceAccountConfig json.RawMessage `json:"service_account_config"`
}

// ServiceAccountCredential is the metadata of a credential (no secret).
type ServiceAccountCredential struct {
	ID               string `json:"id"`
	Kid              string `json:"kid"`
	DisplayName      string `json:"display_name"`
	CreateTime       string `json:"create_time"`
	CreatedBy        string `json:"created_by"`
	OrganizationID   string `json:"organization_id"`
	ServiceAccountID string `json:"service_account_id"`
	ExpireTime       string `json:"expire_time"`
}

// ServiceAccountCredentialAPI is the /configs/v1/service-account-credentials sub-API.
type ServiceAccountCredentialAPI struct {
	t *transport.Client
}

// List returns the service account credentials in an organization (metadata only).
func (a *ServiceAccountCredentialAPI) List(
	ctx context.Context,
	organizationID string,
	opts ...ListOption,
) ([]ServiceAccountCredential, error) {
	return listResource[ServiceAccountCredential](
		ctx, a.t, pathServiceAccountCredentials, orgListQuery(organizationID, opts))
}

// Create mints a new service account credential and returns its signed config.
func (a *ServiceAccountCredentialAPI) Create(
	ctx context.Context,
	req *CreateServiceAccountCredential,
) (*ServiceAccountCredentialResult, error) {
	var out ServiceAccountCredentialResult
	if err := a.t.Do(ctx, http.MethodPost, pathServiceAccountCredentials, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Read returns the metadata of a service account credential.
func (a *ServiceAccountCredentialAPI) Read(ctx context.Context, id string) (*ServiceAccountCredential, error) {
	var out ServiceAccountCredential
	if err := a.t.Do(ctx, http.MethodGet, resourcePath(pathServiceAccountCredentials, id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete revokes a service account credential.
func (a *ServiceAccountCredentialAPI) Delete(ctx context.Context, id string) error {
	return a.t.Do(ctx, http.MethodDelete, resourcePath(pathServiceAccountCredentials, id), nil, nil)
}

// CreateAppAgentCredential is the body to mint an application agent credential.
type CreateAppAgentCredential struct {
	ApplicationAgentID string `json:"application_agent_id"`
	DisplayName        string `json:"display_name,omitempty"`
	ExpireTime         string `json:"expire_time,omitempty"`
}

// AppAgentCredentialResult is returned by Create. ApplicationAgentConfig holds
// the signed credential JSON — capture it, it is not retrievable later.
type AppAgentCredentialResult struct {
	ID                     string          `json:"id"`
	ApplicationAgentID     string          `json:"application_agent_id"`
	Kid                    string          `json:"kid"`
	CreateTime             string          `json:"create_time"`
	CreatedBy              string          `json:"created_by"`
	ExpireTime             string          `json:"expire_time"`
	DisplayName            string          `json:"display_name"`
	ApplicationAgentConfig json.RawMessage `json:"application_agent_config"`
}

// AppAgentCredential is the metadata of a credential (no secret).
type AppAgentCredential struct {
	ID                 string `json:"id"`
	Kid                string `json:"kid"`
	DisplayName        string `json:"display_name"`
	CreateTime         string `json:"create_time"`
	CreatedBy          string `json:"created_by"`
	OrganizationID     string `json:"organization_id"`
	ProjectID          string `json:"project_id"`
	ApplicationID      string `json:"application_id"`
	ApplicationAgentID string `json:"application_agent_id"`
	ExpireTime         string `json:"expire_time"`
}

// AppAgentCredentialAPI is the /configs/v1/application-agent-credentials sub-API.
type AppAgentCredentialAPI struct {
	t *transport.Client
}

// List returns the application agent credentials in a project (metadata only).
func (a *AppAgentCredentialAPI) List(
	ctx context.Context,
	projectID string,
	opts ...ListOption,
) ([]AppAgentCredential, error) {
	return listResource[AppAgentCredential](ctx, a.t, pathAppAgentCredentials, projectListQuery(projectID, opts))
}

// Create mints a new application agent credential and returns its signed config.
func (a *AppAgentCredentialAPI) Create(
	ctx context.Context,
	req *CreateAppAgentCredential,
) (*AppAgentCredentialResult, error) {
	var out AppAgentCredentialResult
	if err := a.t.Do(ctx, http.MethodPost, pathAppAgentCredentials, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Read returns the metadata of an application agent credential.
func (a *AppAgentCredentialAPI) Read(ctx context.Context, id string) (*AppAgentCredential, error) {
	var out AppAgentCredential
	if err := a.t.Do(ctx, http.MethodGet, resourcePath(pathAppAgentCredentials, id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete revokes an application agent credential.
func (a *AppAgentCredentialAPI) Delete(ctx context.Context, id string) error {
	return a.t.Do(ctx, http.MethodDelete, resourcePath(pathAppAgentCredentials, id), nil, nil)
}
