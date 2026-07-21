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

package entitymatching

// CustomPropertyMapping overrides which properties to compare between source and
// target nodes for a run.
type CustomPropertyMapping struct {
	SourceNodeProperty string `json:"source_node_property"`
	TargetNodeProperty string `json:"target_node_property"`
}

// RunRequest is the body for POST /entity-matching/v1/pipelines/:id/runs.
type RunRequest struct {
	// CustomPropertyMappings optionally overrides the suggested mappings.
	CustomPropertyMappings []CustomPropertyMapping `json:"custom_property_mappings,omitempty"`
	// SimilarityScoreCutoff is the match threshold in [0,1].
	SimilarityScoreCutoff float32 `json:"similarity_score_cutoff"`
}

// RunResult is returned (HTTP 202) when a run is accepted.
type RunResult struct {
	ID          string `json:"id"`
	LastRunTime string `json:"last_run_time"`
	Etag        string `json:"etag"`
}

// SuggestedPropertyMapping is one system-suggested property pairing.
type SuggestedPropertyMapping struct {
	SourceNodeType        string  `json:"source_node_type"`
	SourceNodeProperty    string  `json:"source_node_property"`
	TargetNodeType        string  `json:"target_node_type"`
	TargetNodeProperty    string  `json:"target_node_property"`
	SimilarityScoreCutoff float32 `json:"similarity_score_cutoff"`
}

// SuggestedMappings is the result of reading suggested property mappings.
type SuggestedMappings struct {
	ID                        string                     `json:"id"`
	SuggestedPropertyMappings []SuggestedPropertyMapping `json:"suggested_property_mappings"`
}

// PipelineStatus is the run status of an entity-matching pipeline.
type PipelineStatus struct {
	ID                    string `json:"id"`
	PropertyMappingStatus string `json:"property_mapping_status"`
	EntityMatchingStatus  string `json:"entity_matching_status"`
}
