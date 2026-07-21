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

// Command entitymatching demonstrates the Entity Matching runtime API. The
// pipeline itself is configured on the control plane
// (config.AdminClient.EntityMatchingPipelines); here it is run and inspected.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/entitymatching"
	"github.com/indykite/indykite-sdk-go/examples/internal/exutil"
)

func main() {
	if len(os.Args) < 2 {
		exutil.Usage("entitymatching", "run", "status", "mappings")
	}

	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx, exutil.Options()...)
	if err != nil {
		exutil.Fatal(err)
	}
	em := cli.EntityMatching()

	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	pipelineID := fs.String("pipeline-id", "", "entity matching pipeline gid")
	cutoff := fs.Float64("cutoff", 0.85, "similarity score cutoff in [0,1]")
	_ = fs.Parse(os.Args[2:])

	switch os.Args[1] {
	case "run":
		res, err := em.Run(ctx, *pipelineID, entitymatching.RunRequest{
			SimilarityScoreCutoff: float32(*cutoff),
		})
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(res)

	case "status":
		status, err := em.Status(ctx, *pipelineID)
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(status)

	case "mappings":
		mappings, err := em.SuggestedPropertyMappings(ctx, *pipelineID)
		if errors.Is(err, entitymatching.ErrMappingNotReady) {
			fmt.Println("property mappings are still being computed; retry later")
			return
		}
		if err != nil {
			exutil.Fatal(err)
		}
		exutil.Print(mappings)

	default:
		exutil.Usage("entitymatching", "run", "status", "mappings")
	}
}
