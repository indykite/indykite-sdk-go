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

package ciq

// ExecuteRequest is the body for POST /contx-iq/v1/execute.
type ExecuteRequest struct {
	// InputParams supplies values for input parameters referenced by the query.
	InputParams map[string]any `json:"input_params,omitempty"`
	// PreprocessParams supplies CIQ v2.0 preprocess parameters.
	PreprocessParams map[string]string `json:"preprocess_params,omitempty"`
	// ID is the ContX IQ query ID (gid:...) or query name.
	ID string `json:"id"`
	// PageToken selects the result page; any value < 1 returns the first page.
	PageToken int `json:"page_token"`
	// PageSize is the result-set size per page (server default is 100).
	PageSize int `json:"page_size"`
}

// Record is one result row: nodes, relationships and aggregate values, each a
// map keyed by the alias used in the query.
type Record struct {
	Nodes           map[string]any `json:"nodes,omitempty"`
	Relationships   map[string]any `json:"relationships,omitempty"`
	AggregateValues map[string]any `json:"aggregate_values,omitempty"`
}

// ExecuteResponse is the result of one Execute call (a single page).
type ExecuteResponse struct {
	Data []Record `json:"data"`
}
