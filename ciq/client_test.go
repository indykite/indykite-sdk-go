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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

// newClient wires a ciq.Client to an httptest server running h; the server is
// torn down when the spec ends.
func newClient(h http.HandlerFunc) *ciq.Client {
	GinkgoHelper()
	srv := httptest.NewServer(h)
	DeferCleanup(srv.Close)
	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	Expect(err).NotTo(HaveOccurred())
	return ciq.NewClient(tc)
}

// decodeJSONBody reads the request body into a JSON object.
func decodeJSONBody(r *http.Request) map[string]any {
	GinkgoHelper()
	var body map[string]any
	raw, err := io.ReadAll(r.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(json.Unmarshal(raw, &body)).To(Succeed())
	return body
}

// nodeIDs extracts nodes.n.id from every record, in order.
func nodeIDs(rows []ciq.Record) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.Nodes["n"].(map[string]any)["id"].(string))
	}
	return ids
}

var _ = Describe("Execute", func() {
	It("posts the query and decodes a single page", func(ctx SpecContext) {
		var gotPath string
		var gotBody map[string]any
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotBody = decodeJSONBody(r)
			_, _ = io.WriteString(w, `{"data":[{"nodes":{"s":{"id":"gpu-7"}},"aggregate_values":{"count":2}}]}`)
		})

		resp, err := c.Execute(ctx, ciq.ExecuteRequest{
			ID:          "get-servers",
			InputParams: map[string]any{"region": "eu"},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(gotPath).To(Equal("/contx-iq/v1/execute"))
		Expect(gotBody).To(HaveKeyWithValue("id", "get-servers"))
		Expect(gotBody).To(HaveKeyWithValue("input_params", HaveKeyWithValue("region", "eu")))

		Expect(resp.Data).To(HaveLen(1))
		Expect(resp.Data[0].Nodes).To(HaveKeyWithValue("s", HaveKeyWithValue("id", "gpu-7")))
		Expect(resp.Data[0].AggregateValues).To(HaveKeyWithValue("count", BeNumerically("==", 2)))
	})

	It("forwards preprocess_params (CIQ v2.0) when set and omits them otherwise", func(ctx SpecContext) {
		var body map[string]any
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			body = decodeJSONBody(r)
			_, _ = io.WriteString(w, `{"data":[]}`)
		})

		_, err := c.Execute(ctx, ciq.ExecuteRequest{
			ID:               "gid:q",
			PreprocessParams: map[string]string{"mode": "strict"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(HaveKeyWithValue("preprocess_params", HaveKeyWithValue("mode", "strict")))

		body = nil
		_, err = c.Execute(ctx, ciq.ExecuteRequest{ID: "gid:q"})
		Expect(err).NotTo(HaveOccurred())
		Expect(body).NotTo(HaveKey("preprocess_params"), "preprocess_params must be omitted when unset")
	})

	It("surfaces a non-2xx response as an APIError carrying the platform message", func(ctx SpecContext) {
		c := newClient(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"message":"knowledge query not found"}`)
		})

		resp, err := c.Execute(ctx, ciq.ExecuteRequest{ID: "gid:missing"})
		Expect(resp).To(BeNil())
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
		Expect(apiErr.IsNotFound()).To(BeTrue())
		Expect(apiErr.StatusCode).To(Equal(http.StatusNotFound))
		Expect(apiErr.Message).To(Equal("knowledge query not found"))
	})
})

var _ = Describe("Iterate and All", func() {
	// CIQ has no server-side next token: a page shorter than PageSize is the
	// last one.
	It("requests pages until a short page ends the iteration", func(ctx SpecContext) {
		var requestedTokens []int
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			var body ciq.ExecuteRequest
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &body)
			requestedTokens = append(requestedTokens, body.PageToken)

			switch body.PageToken {
			case 1:
				_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"a"}}},{"nodes":{"n":{"id":"b"}}}]}`)
			case 2:
				_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"c"}}}]}`) // short page -> last
			default:
				_, _ = io.WriteString(w, `{"data":[]}`)
			}
		})

		rows, err := c.All(ctx, ciq.ExecuteRequest{ID: "q", PageSize: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodeIDs(rows)).To(Equal([]string{"a", "b", "c"}))
		Expect(requestedTokens).To(Equal([]int{1, 2}))
	})

	It("fetches one extra empty page when the total is an exact multiple of PageSize", func(ctx SpecContext) {
		var requestedTokens []int
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			var body ciq.ExecuteRequest
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &body)
			requestedTokens = append(requestedTokens, body.PageToken)
			if body.PageToken == 1 {
				_, _ = io.WriteString(w, `{"data":[{"nodes":{"n":{"id":"a"}}},{"nodes":{"n":{"id":"b"}}}]}`)
				return
			}
			_, _ = io.WriteString(w, `{"data":[]}`) // page 2 empty -> stop
		})

		rows, err := c.All(ctx, ciq.ExecuteRequest{ID: "q", PageSize: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodeIDs(rows)).To(Equal([]string{"a", "b"}))
		Expect(requestedTokens).To(Equal([]int{1, 2}))
	})
})

var _ = Describe("WhoAmI", func() {
	It("sends the App Agent key plus the end-user Bearer token on a body-less GET", func(ctx SpecContext) {
		var gotMethod, gotPath, gotClientKey, gotAuthz string
		var gotBody []byte
		c := newClient(func(w http.ResponseWriter, r *http.Request) {
			gotMethod, gotPath = r.Method, r.URL.Path
			gotClientKey = r.Header.Get("X-IK-ClientKey")
			gotAuthz = r.Header.Get("Authorization")
			gotBody, _ = io.ReadAll(r.Body)
			_, _ = io.WriteString(w, `{"type":"Person","id":"knightrider"}`)
		})

		me, err := c.WhoAmI(ctx, "eyJ.end.user")
		Expect(err).NotTo(HaveOccurred())

		Expect(gotMethod).To(Equal(http.MethodGet))
		Expect(gotPath).To(Equal("/contx-iq/v1/whoami"))
		Expect(gotClientKey).To(Equal("tok"), "App Agent token must still be sent")
		Expect(gotAuthz).To(Equal("Bearer eyJ.end.user"))
		Expect(gotBody).To(BeEmpty(), "GET must have no body")

		Expect(me).To(HaveValue(Equal(ciq.WhoAmIResponse{Type: "Person", ID: "knightrider"})))
	})

	It("fails locally without an end-user token, never hitting the server", func(ctx SpecContext) {
		called := false
		c := newClient(func(http.ResponseWriter, *http.Request) { called = true })

		me, err := c.WhoAmI(ctx, "")
		Expect(err).To(MatchError(ciq.ErrEndUserTokenRequired))
		Expect(me).To(BeNil())
		Expect(called).To(BeFalse())
	})

	It("surfaces a rejected token as a 401 APIError carrying the platform message", func(ctx SpecContext) {
		c := newClient(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"message":"end-user token is required"}`)
		})

		_, err := c.WhoAmI(ctx, "expired")
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
		Expect(apiErr.IsUnauthorized()).To(BeTrue())
		Expect(apiErr.StatusCode).To(Equal(http.StatusUnauthorized))
		Expect(apiErr.Message).To(Equal("end-user token is required"))
	})
})
