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

const pathAuthorizationPolicies = "/configs/v1/authorization-policies"

// AuthorizationPolicyAPI is the /configs/v1/authorization-policies sub-API.
type AuthorizationPolicyAPI struct {
	t *transport.Client
}

// List returns the authorization policies in a project. Pass a non-empty
// policyType ("kbac" or "ciq") to filter, or "" for all.
func (a *AuthorizationPolicyAPI) List(
	ctx context.Context,
	projectID, policyType string,
	opts ...ListOption,
) ([]AuthorizationPolicy, error) {
	q := projectListQuery(projectID, opts)
	if policyType != "" {
		q.Set("type", policyType)
	}
	return listResource[AuthorizationPolicy](ctx, a.t, pathAuthorizationPolicies, q)
}

// Create creates an authorization policy.
func (a *AuthorizationPolicyAPI) Create(
	ctx context.Context,
	req *CreateAuthorizationPolicy,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathAuthorizationPolicies, req)
}

// Read fetches one policy by gid (or by name with WithLocation). The returned
// ETag is needed for a subsequent Update/Delete.
func (a *AuthorizationPolicyAPI) Read(
	ctx context.Context,
	id string,
	opts ...ReadOption,
) (*AuthorizationPolicy, error) {
	return readResource[AuthorizationPolicy](ctx, a.t, pathAuthorizationPolicies, id, readOptsQuery(opts))
}

// Update updates a policy. Pass the ETag from a prior Read for optimistic
// concurrency ("" to skip the If-Match check).
func (a *AuthorizationPolicyAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateAuthorizationPolicy,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathAuthorizationPolicies, id), req, ifMatch(etag)...)
}

// Delete deletes a policy, optionally guarded by an ETag.
func (a *AuthorizationPolicyAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathAuthorizationPolicies, id, etag)
}
