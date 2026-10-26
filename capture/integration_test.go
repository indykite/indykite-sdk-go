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

package capture_test

import (
	"fmt"
	"os"
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/capture"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NOTE: there is deliberately no integration test for ReadDataSchema — the
// data-schema API is still behind a feature flag.

var _ = Describe("Capture (live platform)", Label("integration"), func() {
	var capc *capture.Client

	BeforeEach(func(ctx SpecContext) {
		if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
		}
		var opts []indykite.Option
		if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
			opts = append(opts, indykite.WithBaseURL(base))
		}
		cli, err := indykite.NewClientFromEnv(ctx, opts...)
		Expect(err).NotTo(HaveOccurred())
		capc = cli.Capture()
	})

	// Ingests two nodes and a relationship, strips a property, then removes
	// everything it created.
	It("runs the node and relationship lifecycle", func(ctx SpecContext) {
		suffix := fmt.Sprintf("sdk-it-%d", time.Now().UnixNano())
		person := capture.Node{ExternalID: "person-" + suffix, Type: "Person"}
		asset := capture.Node{ExternalID: "asset-" + suffix, Type: "Asset"}
		owns := capture.Relationship{Type: "OWNS", Source: &person, Target: &asset}

		// Always try to clean up, even on failure part-way.
		DeferCleanup(func(ctx SpecContext) {
			_, _ = capc.DeleteRelationships(ctx, owns)
			_, _ = capc.DeleteNodes(ctx, person, asset)
		})

		res, err := capc.UpsertNodes(ctx,
			capture.UpsertNode{
				Node:       person,
				IsIdentity: true,
				Properties: []capture.Property{{
					BaseProperty: capture.BaseProperty{Type: "email", Value: suffix + "@example.com"},
					Metadata:     &capture.Metadata{Source: "sdk-integration-test", AssuranceLevel: 1},
				}},
			},
			capture.UpsertNode{Node: asset},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Results).To(HaveLen(2))

		_, err = capc.UpsertRelationships(ctx, owns)
		Expect(err).NotTo(HaveOccurred())

		_, err = capc.DeleteNodePropertyMetadata(ctx, capture.DeleteNodePropertyMetadata{
			Node: person, PropertyType: "email", MetadataFields: []string{"assurance_level"},
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = capc.DeleteNodeProperties(ctx, capture.DeleteNodeProperties{
			Node: person, PropertyTypes: []string{"email"},
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = capc.DeleteRelationships(ctx, owns)
		Expect(err).NotTo(HaveOccurred())

		_, err = capc.DeleteNodes(ctx, person, asset)
		Expect(err).NotTo(HaveOccurred())
	})

	// Verifies the chunked helper splits batches above MaxBatchSize
	// transparently.
	It("splits batches above MaxBatchSize transparently", func(ctx SpecContext) {
		suffix := fmt.Sprintf("sdk-it-chunk-%d", time.Now().UnixNano())
		const n = capture.MaxBatchSize + 10
		nodes := make([]capture.UpsertNode, n)
		deletes := make([]capture.Node, n)
		for i := range nodes {
			node := capture.Node{ExternalID: fmt.Sprintf("%s-%d", suffix, i), Type: "Asset"}
			nodes[i] = capture.UpsertNode{Node: node}
			deletes[i] = node
		}
		DeferCleanup(func(ctx SpecContext) {
			for start := 0; start < n; start += capture.MaxBatchSize {
				end := min(start+capture.MaxBatchSize, n)
				_, _ = capc.DeleteNodes(ctx, deletes[start:end]...)
			}
		})

		res, err := capc.UpsertNodesChunked(ctx, nodes, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Results).To(HaveLen(n))
	})
})
