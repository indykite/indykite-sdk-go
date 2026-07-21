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

const pathDataSchemaRebuild = "/configs/v1/data-schema/rebuild"

// DataSchemaAPI is the /configs/v1/data-schema sub-API.
type DataSchemaAPI struct {
	t *transport.Client
}

type rebuildDataSchemaRequest struct {
	ProjectID string `json:"project_id"`
}

// RebuildResult is the (asynchronous) response of a data-schema rebuild.
type RebuildResult struct {
	Status string `json:"status"`
}

// Rebuild triggers an asynchronous rebuild of a project's data schema.
//
// The platform marks this endpoint as a temporary debugging/testing aid, not
// intended for public use, and it is absent from the published API reference.
// It may change or disappear without notice.
func (a *DataSchemaAPI) Rebuild(ctx context.Context, projectID string) (*RebuildResult, error) {
	body := rebuildDataSchemaRequest{ProjectID: projectID}
	var out RebuildResult
	if err := a.t.Do(ctx, http.MethodPost, pathDataSchemaRebuild, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
