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

// Command ciq demonstrates the ContX IQ query API: you execute a
// pre-configured knowledge query by id and pass its input parameters.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("ciq", "execute", "all")
	}

	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	queryID := fs.String("query-id", "", "knowledge query gid or name")
	params := fs.String("input-params", "{}", "input parameters as JSON object")
	pageSize := fs.Int("page-size", 100, "records per page")
	_ = fs.Parse(os.Args[2:])

	var inputParams map[string]any
	if err = json.Unmarshal([]byte(*params), &inputParams); err != nil {
		exutil.Fatal(fmt.Errorf("-input-params must be a JSON object: %w", err))
	}
	req := ciq.ExecuteRequest{ID: *queryID, InputParams: inputParams, PageSize: *pageSize}

	switch os.Args[1] {
	case "execute":
		// A single page; PageToken selects which one.
		resp, err := cli.CIQ().Execute(ctx, req)
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(resp.Data)

	case "all":
		// The iterator walks size-based pages until a short page ends the set.
		records, err := cli.CIQ().All(ctx, req)
		if err != nil {
			exutil.Fatal(err)
		}
		fmt.Printf("total records: %d\n", len(records))
		exutil.Print(records)

	default:
		exutil.Usage("ciq", "execute", "all")
	}
}
