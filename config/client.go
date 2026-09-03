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

// Package config is the client for the IndyKite Config Management API
// (/configs/v1/*). Unlike the runtime packages, config is the CONTROL PLANE: it
// authenticates with a Service Account token (Bearer) and is exposed through an
// AdminClient. Operations are ETag-guarded for optimistic concurrency.
//
//	admin, _ := config.NewAdminClientFromCredentials(ctx, serviceAccountJSON)
//	pol, _ := admin.AuthorizationPolicies().Create(ctx, config.CreateAuthorizationPolicy{
//	    ProjectID: "gid:...", Name: "can-provision", Policy: policyJSON, Status: config.StatusActive,
//	})
package config

import (
	"context"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/transport"
)

const headerETag = "ETag"

// AdminClient is the entry point for control-plane (config management)
// operations. It must be backed by a Service Account (control-plane)
// authenticator.
type AdminClient struct {
	t *transport.Client
}

// NewAdminClient wraps an existing control-plane transport. The caller is
// responsible for having built it from a Service Account authenticator
// (auth.PlaneControl).
func NewAdminClient(t *transport.Client) *AdminClient {
	return &AdminClient{t: t}
}

// NewAdminClientFromCredentials builds a fully-wired AdminClient from a Service
// Account credentials JSON document. It selects the control plane (Bearer auth)
// automatically.
func NewAdminClientFromCredentials(
	ctx context.Context,
	serviceAccountJSON []byte,
	opts ...transport.Option,
) (*AdminClient, error) {
	a, err := auth.ServiceAccountFromJSON(ctx, serviceAccountJSON)
	if err != nil {
		return nil, err
	}
	t, err := transport.NewClient(a, opts...)
	if err != nil {
		return nil, err
	}
	return &AdminClient{t: t}, nil
}

// Organizations returns the organizations sub-API.
func (a *AdminClient) Organizations() *OrganizationAPI { return &OrganizationAPI{t: a.t} }

// Projects returns the projects sub-API.
func (a *AdminClient) Projects() *ProjectAPI { return &ProjectAPI{t: a.t} }

// Applications returns the applications sub-API.
func (a *AdminClient) Applications() *ApplicationAPI { return &ApplicationAPI{t: a.t} }

// ServiceAccounts returns the service-accounts sub-API.
func (a *AdminClient) ServiceAccounts() *ServiceAccountAPI { return &ServiceAccountAPI{t: a.t} }

// ServiceAccountCredentials returns the service-account-credentials sub-API.
func (a *AdminClient) ServiceAccountCredentials() *ServiceAccountCredentialAPI {
	return &ServiceAccountCredentialAPI{t: a.t}
}

// AppAgents returns the application-agents sub-API.
func (a *AdminClient) AppAgents() *AppAgentAPI {
	return &AppAgentAPI{t: a.t}
}

// AppAgentCredentials returns the application-agent-credentials sub-API.
func (a *AdminClient) AppAgentCredentials() *AppAgentCredentialAPI {
	return &AppAgentCredentialAPI{t: a.t}
}

// AuthorizationPolicies returns the authorization-policies sub-API.
func (a *AdminClient) AuthorizationPolicies() *AuthorizationPolicyAPI {
	return &AuthorizationPolicyAPI{t: a.t}
}

// KnowledgeQueries returns the knowledge-queries sub-API.
func (a *AdminClient) KnowledgeQueries() *KnowledgeQueryAPI { return &KnowledgeQueryAPI{t: a.t} }

// TokenIntrospects returns the token-introspects sub-API.
func (a *AdminClient) TokenIntrospects() *TokenIntrospectAPI { return &TokenIntrospectAPI{t: a.t} }

// ExternalDataResolvers returns the external-data-resolvers sub-API.
func (a *AdminClient) ExternalDataResolvers() *ExternalDataResolverAPI {
	return &ExternalDataResolverAPI{t: a.t}
}

// TrustScoreProfiles returns the trust-score-profiles sub-API.
func (a *AdminClient) TrustScoreProfiles() *TrustScoreProfileAPI {
	return &TrustScoreProfileAPI{t: a.t}
}

// EntityMatchingPipelines returns the entity-matching-pipelines sub-API.
func (a *AdminClient) EntityMatchingPipelines() *EntityMatchingPipelineAPI {
	return &EntityMatchingPipelineAPI{t: a.t}
}

// EventSinks returns the event-sinks sub-API.
func (a *AdminClient) EventSinks() *EventSinkAPI { return &EventSinkAPI{t: a.t} }

// MCPServers returns the mcp-servers sub-API.
func (a *AdminClient) MCPServers() *MCPServerAPI { return &MCPServerAPI{t: a.t} }

// AuditSignings returns the audit-signings sub-API.
func (a *AdminClient) AuditSignings() *AuditSigningAPI { return &AuditSigningAPI{t: a.t} }

// DataSchema returns the data-schema sub-API.
func (a *AdminClient) DataSchema() *DataSchemaAPI { return &DataSchemaAPI{t: a.t} }

// ifMatch returns a CallOption that sets the If-Match header when etag is set.
func ifMatch(etag string) []transport.CallOption {
	if etag == "" {
		return nil
	}
	return []transport.CallOption{transport.WithHeader("If-Match", etag)}
}
