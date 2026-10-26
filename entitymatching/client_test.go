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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/entitymatching"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

// newClient wires an entitymatching.Client to an httptest server running h;
// the server is torn down when the spec ends.
func newClient(h http.HandlerFunc) *entitymatching.Client {
	GinkgoHelper()
	srv := httptest.NewServer(h)
	DeferCleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	Expect(err).NotTo(HaveOccurred())
	return entitymatching.NewClient(tc)
}

var _ = Describe("Run", func() {
	It("POSTs the run request with runtime-plane auth and decodes the 202 reply", func(ctx SpecContext) {
		var path, key string
		var body map[string]any
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			path = r.URL.Path
			key = r.Header.Get("X-Ik-Clientkey") // canonical form of the X-IK-ClientKey header
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &body)
			w.WriteHeader(http.StatusAccepted) // 202
			_, _ = io.WriteString(w, `{"id":"gid:em1","last_run_time":"t0","etag":"\"v1\""}`)
		})

		run, err := c.Run(ctx, "gid:em1", entitymatching.RunRequest{
			SimilarityScoreCutoff: 0.8,
			CustomPropertyMappings: []entitymatching.CustomPropertyMapping{
				{SourceNodeProperty: "email", TargetNodeProperty: "email"},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(path).To(Equal("/entity-matching/v1/pipelines/gid:em1/runs"))
		Expect(key).To(Equal("tok"), "expected runtime-plane X-IK-ClientKey")
		Expect(body).To(HaveKeyWithValue("similarity_score_cutoff", 0.8))
		cpm := body["custom_property_mappings"].([]any)[0]
		Expect(cpm).To(HaveKeyWithValue("source_node_property", "email"))
		Expect(run.ID).To(Equal("gid:em1"))
		Expect(run.Etag).To(Equal(`"v1"`))
	})
})

var _ = Describe("SuggestedPropertyMappings", func() {
	It("decodes the mappings once they are ready", func(ctx SpecContext) {
		var path string
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			path = r.URL.Path
			_, _ = io.WriteString(w, `{"id":"gid:em1","suggested_property_mappings":[
				{"source_node_type":"Person","source_node_property":"email",
				 "target_node_type":"Lead","target_node_property":"email","similarity_score_cutoff":0.9}]}`)
		})

		m, err := c.SuggestedPropertyMappings(ctx, "gid:em1")
		Expect(err).NotTo(HaveOccurred())
		Expect(path).To(Equal("/entity-matching/v1/pipelines/gid:em1/property-mappings"))
		Expect(m.SuggestedPropertyMappings).To(HaveLen(1))
		sm := m.SuggestedPropertyMappings[0]
		Expect(sm.SourceNodeType).To(Equal("Person"))
		Expect(sm.SimilarityScoreCutoff).To(BeEquivalentTo(float32(0.9)))
	})

	It("returns ErrMappingNotReady while the mappings are still computing", func(ctx SpecContext) {
		c := newClient(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted) // 202 -> still computing
			_, _ = io.WriteString(w, `{"message":"mapping in progress"}`)
		})

		_, err := c.SuggestedPropertyMappings(ctx, "gid:em1")
		Expect(err).To(MatchError(entitymatching.ErrMappingNotReady))
	})
})

var _ = Describe("Status", func() {
	It("GETs the pipeline status", func(ctx SpecContext) {
		var method, path string
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			method, path = r.Method, r.URL.Path
			_, _ = io.WriteString(w,
				`{"id":"gid:em1","property_mapping_status":"DONE","entity_matching_status":"RUNNING"}`)
		})

		st, err := c.Status(ctx, "gid:em1")
		Expect(err).NotTo(HaveOccurred())
		Expect(method).To(Equal(http.MethodGet))
		Expect(path).To(Equal("/entity-matching/v1/pipelines/gid:em1/status"))
		Expect(st.PropertyMappingStatus).To(Equal("DONE"))
		Expect(st.EntityMatchingStatus).To(Equal("RUNNING"))
	})
})
