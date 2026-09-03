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

const pathAuditSignings = "/configs/v1/audit-signings"

// Audit signing provider values: who manages the key used to sign audit records.
const (
	// AuditSigningProviderPlatformManaged lets the platform own the signing key;
	// KeyResource, Kid and AuthParams are not required.
	AuditSigningProviderPlatformManaged = "PLATFORM_MANAGED"
	// AuditSigningProviderCustomerGCPKMS signs with a customer key in Google Cloud KMS.
	AuditSigningProviderCustomerGCPKMS = "CUSTOMER_GCP_KMS"
	// AuditSigningProviderCustomerAWSKMS signs with a customer key in AWS KMS.
	AuditSigningProviderCustomerAWSKMS = "CUSTOMER_AWS_KMS"
	// AuditSigningProviderCustomerAzureKeyVault signs with a customer key in Azure Key Vault.
	AuditSigningProviderCustomerAzureKeyVault = "CUSTOMER_AZURE_KEY_VAULT"
)

// AuditSigning is an audit signing configuration (project scoped). It tells the
// platform which key signs the project's audit records.
//
// On reads the platform never returns AuthParams secret material: the map keeps
// the keys that are set but every value is blank.
type AuditSigning struct {
	Metadata
	AuthParams map[string]string `json:"auth_params,omitempty"`
	Provider   string            `json:"provider,omitempty"`
	// KeyResource is the provider-specific key reference (e.g. a KMS key resource name).
	KeyResource string `json:"key_resource,omitempty"`
	// Kid is the key identifier published alongside signatures.
	Kid string `json:"kid,omitempty"`
	Versioned
}

// AuditSigningConfig is the signing-specific config shared by create/update.
//
// Provider is always required. KeyResource and Kid are required for every
// customer-managed provider. AuthParams holds provider credentials (at most 32
// pairs, non-empty keys and values) and is write-only.
type AuditSigningConfig struct {
	AuthParams  map[string]string `json:"auth_params,omitempty"`
	Provider    string            `json:"provider"`
	KeyResource string            `json:"key_resource,omitempty"`
	Kid         string            `json:"kid,omitempty"`
}

// CreateAuditSigning is the body to create an audit signing config.
type CreateAuditSigning struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	AuditSigningConfig
}

// UpdateAuditSigning is the body to update an audit signing config.
type UpdateAuditSigning struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	AuditSigningConfig
}

// AuditSigningAPI is the /configs/v1/audit-signings sub-API.
type AuditSigningAPI struct {
	t *transport.Client
}

// List returns the audit signing configs in a project.
func (a *AuditSigningAPI) List(ctx context.Context, projectID string, opts ...ListOption) ([]AuditSigning, error) {
	return listResource[AuditSigning](ctx, a.t, pathAuditSignings, projectListQuery(projectID, opts))
}

// Create creates an audit signing config.
func (a *AuditSigningAPI) Create(ctx context.Context, req *CreateAuditSigning) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathAuditSignings, req)
}

// Read fetches one audit signing config by gid (or by name with WithLocation).
func (a *AuditSigningAPI) Read(ctx context.Context, id string, opts ...ReadOption) (*AuditSigning, error) {
	return readResource[AuditSigning](ctx, a.t, pathAuditSignings, id, readOptsQuery(opts))
}

// Update updates an audit signing config, optionally guarded by an ETag.
func (a *AuditSigningAPI) Update(ctx context.Context, id, etag string, req *UpdateAuditSigning) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathAuditSignings, id), req, ifMatch(etag)...)
}

// Delete deletes an audit signing config, optionally guarded by an ETag.
func (a *AuditSigningAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathAuditSignings, id, etag)
}
