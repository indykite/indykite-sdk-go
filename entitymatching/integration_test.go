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
	"errors"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/entitymatching"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Entity matching (live platform)", Label("integration"), func() {
	var (
		em *entitymatching.Client
		id string
	)

	BeforeEach(func(ctx SpecContext) {
		if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
		}
		id = os.Getenv("EM_PIPELINE_ID")
		if id == "" {
			Skip("EM_PIPELINE_ID not set")
		}
		var opts []indykite.Option
		if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
			opts = append(opts, indykite.WithBaseURL(base))
		}
		cli, err := indykite.NewClientFromEnv(ctx, opts...)
		Expect(err).NotTo(HaveOccurred())
		em = cli.EntityMatching()
	})

	It("reports the pipeline status", func(ctx SpecContext) {
		status, err := em.Status(ctx, id)
		Expect(err).NotTo(HaveOccurred())
		Expect(status.ID).NotTo(BeEmpty(), "status has no pipeline id: %+v", status)
		GinkgoWriter.Printf("property_mapping=%s entity_matching=%s\n",
			status.PropertyMappingStatus, status.EntityMatchingStatus)
	})

	It("returns the suggested property mappings", func(ctx SpecContext) {
		mappings, err := em.SuggestedPropertyMappings(ctx, id)
		if errors.Is(err, entitymatching.ErrMappingNotReady) {
			Skip("suggested property mappings not computed yet")
		}
		Expect(err).NotTo(HaveOccurred())
		GinkgoWriter.Printf("suggested mappings: %d\n", len(mappings.SuggestedPropertyMappings))
	})

	// Triggers a pipeline run only when the environment opts in (a run mutates
	// pipeline state and can take a while).
	It("triggers a pipeline run", func(ctx SpecContext) {
		if os.Getenv("EM_RUN") == "" {
			Skip("EM_RUN not set; skipping mutating pipeline run")
		}

		res, err := em.Run(ctx, id, entitymatching.RunRequest{
			SimilarityScoreCutoff: 0.9,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.ID).NotTo(BeEmpty(), "run result has no id: %+v", res)
	})
})
