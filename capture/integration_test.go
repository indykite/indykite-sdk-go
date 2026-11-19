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
	"strings"
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/capture"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NOTE: there is deliberately no integration test for ReadDataSchema — the
// data-schema API is still behind a feature flag.

// liveOptions returns the client options shared by the live specs.
func liveOptions() []indykite.Option {
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		return []indykite.Option{indykite.WithBaseURL(base)}
	}
	return nil
}

// haveGIDs matches a per-location list of n gids.
func haveGIDs(n int) OmegaMatcher {
	return And(HaveLen(n), HaveEach(HavePrefix("gid:")))
}

var _ = Describe("Capture (live platform)", Label("integration"), func() {
	var capc *capture.Client

	BeforeEach(func(ctx SpecContext) {
		if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
		}
		cli, err := indykite.NewClientFromEnv(ctx, liveOptions()...)
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

// These specs need a project whose IKG is a composite database, so they are
// opt-in:
//
//	CAPTURE_COMPOSITE_LOCATIONS                two location names of that IKG other
//	                                           than "global", comma-separated (required to run)
//	CAPTURE_COMPOSITE_APPLICATION_CREDENTIALS  App Agent token of that project; defaults to
//	                                           INDYKITE_APPLICATION_CREDENTIALS[_FILE]
var _ = Describe("Capture on a composite IKG (live platform)", Label("integration"), func() {
	var (
		capc       *capture.Client
		east, west string
	)

	BeforeEach(func(ctx SpecContext) {
		raw := os.Getenv("CAPTURE_COMPOSITE_LOCATIONS")
		if raw == "" {
			Skip("CAPTURE_COMPOSITE_LOCATIONS not set")
		}
		locations := strings.Split(raw, ",")
		Expect(locations).To(HaveLen(2), "CAPTURE_COMPOSITE_LOCATIONS must name two locations")
		east, west = strings.TrimSpace(locations[0]), strings.TrimSpace(locations[1])
		Expect([]string{east, west}).NotTo(Or(
			ContainElement(BeEmpty()),
			ContainElement(capture.LocationGlobal),
			ContainElement(capture.LocationDefault),
		), "CAPTURE_COMPOSITE_LOCATIONS must name two regional locations")
		Expect(east).NotTo(Equal(west), "CAPTURE_COMPOSITE_LOCATIONS must name two different locations")

		var (
			cli *indykite.Client
			err error
		)
		switch token := os.Getenv("CAPTURE_COMPOSITE_APPLICATION_CREDENTIALS"); {
		case token != "":
			cli, err = indykite.NewClient(ctx, token, liveOptions()...)
		case os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "":
			Skip("no App Agent credential set for the composite project")
		default:
			cli, err = indykite.NewClientFromEnv(ctx, liveOptions()...)
		}
		Expect(err).NotTo(HaveOccurred())
		capc = cli.Capture()
	})

	// Ingests nodes into two regional locations, the global one and the default
	// database, links them, and removes everything. Every response must list
	// the gids under the location the platform routed them to.
	It("routes nodes and relationships by location", func(ctx SpecContext) {
		suffix := fmt.Sprintf("sdk-it-cdb-%d", time.Now().UnixNano())
		node := func(name, nodeType, location string) capture.Node {
			return capture.Node{ExternalID: name + "-" + suffix, Type: nodeType, Location: location}
		}
		eastPerson := node("east-person", "Person", east)
		westPerson := node("west-person", "Person", west)
		westCar := node("west-car", "Car", west)
		defaultCar := node("default-car", "Car", "")
		globalPerson := node("global-person", "Person", capture.LocationGlobal)
		nodes := []capture.Node{eastPerson, westPerson, westCar, defaultCar, globalPerson}

		// Both ends in the same location: stored in that location.
		owns := capture.Relationship{Type: "OWNS", Source: &westPerson, Target: &westCar}
		// Ends in different locations: stored in the global location.
		knows := capture.Relationship{Type: "KNOWS", Source: &eastPerson, Target: &globalPerson}

		// Always try to clean up, even on failure part-way.
		DeferCleanup(func(ctx SpecContext) {
			_, _ = capc.DeleteRelationships(ctx, owns, knows)
			_, _ = capc.DeleteNodes(ctx, nodes...)
		})

		// Only the east person asks for a proxy, so the global location holds the
		// global person plus that one proxy.
		nodesByLocation := And(
			HaveLen(4),
			HaveKeyWithValue(east, haveGIDs(1)),
			HaveKeyWithValue(west, haveGIDs(2)),
			HaveKeyWithValue(capture.LocationDefault, haveGIDs(1)),
			HaveKeyWithValue(capture.LocationGlobal, haveGIDs(2)),
		)
		relationshipsByLocation := And(
			HaveLen(2),
			HaveKeyWithValue(west, haveGIDs(1)),
			HaveKeyWithValue(capture.LocationGlobal, haveGIDs(1)),
		)

		upserts := make([]capture.UpsertNode, len(nodes))
		for i, n := range nodes {
			upserts[i] = capture.UpsertNode{Node: n, CreateProxy: n == eastPerson}
		}
		res, err := capc.UpsertNodes(ctx, upserts...)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Results).To(BeEmpty())
		Expect(res.ByLocation).To(nodesByLocation)

		res, err = capc.UpsertRelationships(ctx, owns, knows)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.ByLocation).To(relationshipsByLocation)

		res, err = capc.DeleteRelationships(ctx, owns, knows)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.ByLocation).To(relationshipsByLocation)

		res, err = capc.DeleteNodes(ctx, nodes...)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.ByLocation).To(nodesByLocation)
	})
})
