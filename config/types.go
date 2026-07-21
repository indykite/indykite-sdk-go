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

// Metadata holds the fields common to every configuration object.
type Metadata struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	CreateTime     string `json:"create_time"`
	CreatedBy      string `json:"created_by"`
	UpdateTime     string `json:"update_time"`
	UpdatedBy      string `json:"updated_by"`
	OrganizationID string `json:"organization_id,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
}

// WriteResult is the response of a create/update/delete operation. ETag is read
// from the response header (not the body) and is used for optimistic
// concurrency on the next update/delete.
type WriteResult struct {
	ID         string `json:"id"`
	CreateTime string `json:"create_time"`
	CreatedBy  string `json:"created_by"`
	UpdateTime string `json:"update_time"`
	UpdatedBy  string `json:"updated_by"`

	ETag string `json:"-"`
}

// Versioned is embedded by every readable config resource. Its ETag is read from
// the response header (not the body) for optimistic concurrency.
type Versioned struct {
	ETag string `json:"-"`
}

func (v *Versioned) setETag(etag string) { v.ETag = etag }

// etagSettable is satisfied by *T for any resource embedding Versioned; it lets
// the generic read helper populate the ETag from the response header.
type etagSettable interface {
	setETag(etag string)
}

// AuthorizationPolicy status values.
const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusDraft    = "DRAFT"
)

// AuthorizationPolicy is a KBAC/CIQ authorization policy configuration.
type AuthorizationPolicy struct {
	Metadata
	// Policy is the policy document as a JSON string.
	Policy string `json:"policy"`
	// Status is one of ACTIVE, INACTIVE, DRAFT.
	Status string `json:"status"`
	Versioned
	// Tags optionally label the policy.
	Tags []string `json:"tags,omitempty"`
}

// CreateAuthorizationPolicy is the body to create an authorization policy.
type CreateAuthorizationPolicy struct {
	ProjectID   string   `json:"project_id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name,omitempty"`
	Description string   `json:"description,omitempty"`
	Policy      string   `json:"policy"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags,omitempty"`
}

// UpdateAuthorizationPolicy is the body to update an authorization policy.
type UpdateAuthorizationPolicy struct {
	DisplayName *string  `json:"display_name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Policy      string   `json:"policy"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags,omitempty"`
}

// AppAgent API permission values.
const (
	PermissionAuthorization  = "Authorization"
	PermissionCapture        = "Capture"
	PermissionContXIQ        = "ContXIQ"
	PermissionEntityMatching = "EntityMatching"
	PermissionIKGRead        = "IKGRead"
	PermissionReadDataSchema = "ReadDataSchema"
)

// AppAgent is an application agent configuration.
type AppAgent struct {
	Metadata
	ApplicationID string `json:"application_id"`
	Versioned
	APIPermissions []string `json:"api_permissions"`
}

// CreateAppAgent is the body to create an application agent.
type CreateAppAgent struct {
	ApplicationID  string   `json:"application_id"`
	Name           string   `json:"name"`
	DisplayName    string   `json:"display_name,omitempty"`
	Description    string   `json:"description,omitempty"`
	APIPermissions []string `json:"api_permissions"`
}

// UpdateAppAgent is the body to update an application agent.
type UpdateAppAgent struct {
	DisplayName    *string  `json:"display_name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	APIPermissions []string `json:"api_permissions,omitempty"`
}

// listResponse is the generic {"data":[...]} list envelope.
type listResponse[T any] struct {
	Data []T `json:"data"`
}
