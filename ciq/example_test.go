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
	"fmt"

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
		fmt.Println(record.Nodes)
	}
	if err := it.Err(); err != nil {
		return
	}
}
