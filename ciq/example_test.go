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

package ciq_test

import (
	"context"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/ciq"
)

// Walk a query's result set page by page with the iterator.
func ExampleClient_Iterate() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	it := cli.CIQ().Iterate(ciq.ExecuteRequest{
		ID:          "get-servers",
		InputParams: map[string]any{"region": "eu"},
		PageSize:    100,
	})
	for it.Next(ctx) {
		record := it.Item()
		_ = record.Nodes // one row: nodes, relationships and aggregates keyed by query alias
	}
	if err := it.Err(); err != nil {
		return
	}
}

// Resolve a third-party end-user token to the IKG node it maps to.
func ExampleClient_WhoAmI() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	// The raw end-user token (e.g. from an incoming Authorization header) that
	// the App Agent's Token Introspect configuration accepts.
	endUserToken := os.Getenv("END_USER_TOKEN")
	me, err := cli.CIQ().WhoAmI(ctx, endUserToken)
	if err != nil {
		return
	}
	_ = me.Type // IKG node type, e.g. "Person"
	_ = me.ID   // token subject == node external_id, e.g. "alice@example.com"
}
