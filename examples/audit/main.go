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

// Command audit demonstrates the tamper-proof audit log API: paging through a
// project's signed logs, manifests and checkpoints, and fetching the public
// keys their signatures verify against.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/audit"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("audit", "logs", "manifests", "checkpoints", "jwks")
	}

	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	a := cli.Audit()

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	projectID := fs.String("project-id", os.Getenv("PROJECT_ID"),
		"project (app space) gid the App Agent belongs to (default $PROJECT_ID)")
	cursor := fs.String("cursor", "", "next_cursor of a previous page; empty starts at the beginning")
	pageSize := fs.Int("page-size", 0, "items per page (server default and maximum: 50)")
	all := fs.Bool("all", false, "follow next_cursor through every page")
	_ = fs.Parse(os.Args[2:])

	// Listing needs the Audit App Agent permission, and only the App Agent's
	// own project can be read.
	req := audit.ListRequest{ProjectID: *projectID, Cursor: *cursor, PageSize: *pageSize}

	switch os.Args[1] {
	case "logs":
		// The signed chain batches, carrying the audit events, in sequence order.
		if *all {
			printAll(a.AllLogs(ctx, req))
			return
		}
		printPage(a.ListLogs(ctx, req))

	case "manifests":
		// The chain-linking metadata, in lockstep with logs but without payloads.
		if *all {
			printAll(a.AllManifests(ctx, req))
			return
		}
		printPage(a.ListManifests(ctx, req))

	case "checkpoints":
		// The signed trust anchors, newest first.
		if *all {
			printAll(a.AllCheckpoints(ctx, req))
			return
		}
		printPage(a.ListCheckpoints(ctx, req))

	case "jwks":
		// Public: signatures are ECDSA P-256/SHA-256, ASN.1 DER encoded; verify
		// them with ecdsa.VerifyASN1 and Key.ECDSAPublicKey.
		keys, err := a.JWKS(ctx, *projectID)
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(keys)

	default:
		exutil.Usage("audit", "logs", "manifests", "checkpoints", "jwks")
	}
}

func printPage[T any](page *audit.Page[T], err error) {
	if err != nil {
		exutil.Fatal(err)
	}
	exutil.Print(page)
}

func printAll[T any](items []T, err error) {
	if err != nil {
		exutil.Fatal(err)
	}
	fmt.Printf("total: %d\n", len(items))
	exutil.Print(items)
}
