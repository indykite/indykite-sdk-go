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

//go:build integration

package entitymatching_test

import (
	"context"
	"errors"
	"os"
	"testing"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/entitymatching"
)

func pipelineID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("EM_PIPELINE_ID")
	if id == "" {
		t.Skip("EM_PIPELINE_ID not set")
	}
	return id
}

func emClient(t *testing.T) *entitymatching.Client {
	t.Helper()
	if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
		t.Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
	}
	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	cli, err := indykite.NewClientFromEnv(context.Background(), opts...)
	if err != nil {
		t.Fatalf("NewClientFromEnv: %v", err)
	}
	return cli.EntityMatching()
}

func TestIntegrationEntityMatchingStatus(t *testing.T) {
	em, id := emClient(t), pipelineID(t)

	status, err := em.Status(context.Background(), id)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if status.ID == "" {
		t.Errorf("status has no pipeline id: %+v", status)
	}
	t.Logf("property_mapping=%s entity_matching=%s",
		status.PropertyMappingStatus, status.EntityMatchingStatus)
}

func TestIntegrationEntityMatchingSuggestedMappings(t *testing.T) {
	em, id := emClient(t), pipelineID(t)

	mappings, err := em.SuggestedPropertyMappings(context.Background(), id)
	if errors.Is(err, entitymatching.ErrMappingNotReady) {
		t.Skip("suggested property mappings not computed yet")
	}
	if err != nil {
		t.Fatalf("SuggestedPropertyMappings: %v", err)
	}
	t.Logf("suggested mappings: %d", len(mappings.SuggestedPropertyMappings))
}

// TestIntegrationEntityMatchingRun triggers a pipeline run only when the
// environment opts in (a run mutates pipeline state and can take a while).
func TestIntegrationEntityMatchingRun(t *testing.T) {
	em, id := emClient(t), pipelineID(t)
	if os.Getenv("EM_RUN") == "" {
		t.Skip("EM_RUN not set; skipping mutating pipeline run")
	}

	res, err := em.Run(context.Background(), id, entitymatching.RunRequest{
		SimilarityScoreCutoff: 0.9,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.ID == "" {
		t.Errorf("run result has no id: %+v", res)
	}
}
