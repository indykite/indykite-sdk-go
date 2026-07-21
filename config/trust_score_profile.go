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

const pathTrustScoreProfiles = "/configs/v1/trust-score-profiles"

// Trust score schedule values.
const (
	ScheduleThreeHours  = "THREE_HOURS"
	ScheduleSixHours    = "SIX_HOURS"
	ScheduleTwelveHours = "TWELVE_HOURS"
	ScheduleDaily       = "DAILY"
)

// TrustScoreProfile configures periodic trust-score computation (project scoped).
// Dimensions is passed through as raw JSON; see the platform API reference.
type TrustScoreProfile struct {
	Metadata
	NodeClassification string `json:"node_classification"`
	Schedule           string `json:"schedule"`
	LastRunID          string `json:"last_run_id,omitempty"`
	LastRunStartTime   string `json:"last_run_start_time,omitempty"`
	LastRunEndTime     string `json:"last_run_end_time,omitempty"`
	Versioned
	Dimensions json.RawMessage `json:"dimensions,omitempty"`
	// DimensionsExecutionTimes maps each dimension to its last execution time.
	DimensionsExecutionTimes json.RawMessage `json:"dimensions_execution_times,omitempty"`
}

// CreateTrustScoreProfile is the body to create a trust score profile.
type CreateTrustScoreProfile struct {
	ProjectID          string          `json:"project_id"`
	Name               string          `json:"name"`
	DisplayName        string          `json:"display_name,omitempty"`
	Description        string          `json:"description,omitempty"`
	NodeClassification string          `json:"node_classification"`
	Schedule           string          `json:"schedule"`
	Dimensions         json.RawMessage `json:"dimensions"`
}

// UpdateTrustScoreProfile is the body to update a trust score profile.
type UpdateTrustScoreProfile struct {
	DisplayName        *string         `json:"display_name,omitempty"`
	Description        *string         `json:"description,omitempty"`
	NodeClassification string          `json:"node_classification,omitempty"`
	Schedule           string          `json:"schedule,omitempty"`
	Dimensions         json.RawMessage `json:"dimensions,omitempty"`
}

// TrustScoreProfileAPI is the /configs/v1/trust-score-profiles sub-API.
type TrustScoreProfileAPI struct {
	t *transport.Client
}

// List returns the trust score profiles in a project.
func (a *TrustScoreProfileAPI) List(
	ctx context.Context,
	projectID string,
	opts ...ListOption,
) ([]TrustScoreProfile, error) {
	return listResource[TrustScoreProfile](ctx, a.t, pathTrustScoreProfiles, projectListQuery(projectID, opts))
}

// Create creates a trust score profile.
func (a *TrustScoreProfileAPI) Create(ctx context.Context, req *CreateTrustScoreProfile) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPost, pathTrustScoreProfiles, req)
}

// Read fetches one trust score profile by gid (or by name with WithLocation).
func (a *TrustScoreProfileAPI) Read(
	ctx context.Context,
	id string,
	opts ...ReadOption,
) (*TrustScoreProfile, error) {
	return readResource[TrustScoreProfile](ctx, a.t, pathTrustScoreProfiles, id, readOptsQuery(opts))
}

// Update updates a trust score profile, optionally guarded by an ETag.
func (a *TrustScoreProfileAPI) Update(
	ctx context.Context,
	id, etag string,
	req *UpdateTrustScoreProfile,
) (*WriteResult, error) {
	return write(ctx, a.t, http.MethodPut, resourcePath(pathTrustScoreProfiles, id), req, ifMatch(etag)...)
}

// Delete deletes a trust score profile, optionally guarded by an ETag.
func (a *TrustScoreProfileAPI) Delete(ctx context.Context, id, etag string) error {
	return deleteResource(ctx, a.t, pathTrustScoreProfiles, id, etag)
}
