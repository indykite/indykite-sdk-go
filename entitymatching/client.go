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

// Package entitymatching is the client for the IndyKite Entity Matching API
// (/entity-matching/v1/*). It runs on the runtime plane (App Agent token) and is
// a thin facade over a *transport.Client.
//
//	em := entitymatching.NewClient(client)
//	run, _ := em.Run(ctx, pipelineID, entitymatching.RunRequest{SimilarityScoreCutoff: 0.8})
package entitymatching

import (
	"context"
	"errors"
	"net/http"

	"github.com/indykite/indykite-sdk-go/transport"
)

// ErrMappingNotReady is returned by SuggestedPropertyMappings when the pipeline
// is still computing its mappings (the API replies HTTP 202).
var ErrMappingNotReady = errors.New("entitymatching: suggested property mappings not ready yet")

// Client calls the Entity Matching API.
type Client struct {
	t *transport.Client
}

// NewClient builds an Entity Matching client over the shared transport.
func NewClient(t *transport.Client) *Client {
	return &Client{t: t}
}

func pipelinePath(pipelineID, suffix string) string {
	return "/entity-matching/v1/pipelines/" + pipelineID + suffix
}

// Run triggers a pipeline run. The API accepts the run asynchronously (HTTP 202).
func (c *Client) Run(ctx context.Context, pipelineID string, req RunRequest) (*RunResult, error) {
	var out RunResult
	if err := c.t.Do(ctx, http.MethodPost, pipelinePath(pipelineID, "/runs"), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SuggestedPropertyMappings returns the system-suggested mappings for a pipeline.
// It returns ErrMappingNotReady if the mappings are still being computed.
func (c *Client) SuggestedPropertyMappings(ctx context.Context, pipelineID string) (*SuggestedMappings, error) {
	var out SuggestedMappings
	resp, err := c.t.DoResp(ctx, http.MethodGet, pipelinePath(pipelineID, "/property-mappings"), nil, &out)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusAccepted {
		return nil, ErrMappingNotReady
	}
	return &out, nil
}

// Status returns the current run status of a pipeline.
func (c *Client) Status(ctx context.Context, pipelineID string) (*PipelineStatus, error) {
	var out PipelineStatus
	if err := c.t.Do(ctx, http.MethodGet, pipelinePath(pipelineID, "/status"), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
