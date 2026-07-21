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

package entitymatching_test

import (
	"context"
	"errors"
	"fmt"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/entitymatching"
)

// Inspect the suggested mappings, then trigger a run with a custom cutoff.
func ExampleClient_Run() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}
	em := cli.EntityMatching()
	pipelineID := "gid:pipeline"

	mappings, err := em.SuggestedPropertyMappings(ctx, pipelineID)
	if errors.Is(err, entitymatching.ErrMappingNotReady) {
		fmt.Println("mappings still being computed; retry later")
		return
	}
	if err != nil {
		return
	}
	fmt.Println("suggested:", len(mappings.SuggestedPropertyMappings))

	if _, err = em.Run(ctx, pipelineID, entitymatching.RunRequest{
		SimilarityScoreCutoff: 0.9,
	}); err != nil {
		return
	}
}
