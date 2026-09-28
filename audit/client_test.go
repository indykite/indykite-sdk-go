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

package audit_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/indykite/indykite-sdk-go/audit"
	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const projectID = "gid:AAAAAmluZHlraURlgAAAAAAAAA8"

type stubProvider struct{}

func (stubProvider) Token(context.Context) (string, error) { return "tok", nil }

// request is what the fake server saw on one call.
type request struct {
	query  url.Values
	header http.Header
	method string
	path   string
}

// newClient wires an audit.Client to an httptest server running h, recording
// every request it sees; the server is torn down when the spec ends.
func newClient(h http.HandlerFunc) (*audit.Client, *[]request) {
	GinkgoHelper()
	var seen []request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, request{
			method: r.Method, path: r.URL.Path, query: r.URL.Query(), header: r.Header.Clone(),
		})
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}))
	DeferCleanup(srv.Close)

	a := auth.NewWithProvider(auth.PlaneRuntime, stubProvider{})
	tc, err := transport.NewClient(a, transport.WithBaseURL(srv.URL))
	Expect(err).NotTo(HaveOccurred())
	return audit.NewClient(tc), &seen
}

func reply(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, body) }
}

var _ = Describe("ListLogs", func() {
	It("issues an authenticated GET with the project, cursor and page size", func(ctx SpecContext) {
		c, seen := newClient(reply(`{
			"next_cursor":"Mg","has_more":true,
			"items":[{"batch_id":"b1","project_id":"gid:p","sequence":1,"hash":"h1","signature":"s1",
				"kid":"k1","alg":"ecdsa-p256-sha256-der","chain_hash":"c1","manifest_id":"m1",
				"data":[{"event":"x"}]}]}`))

		page, err := c.ListLogs(ctx, audit.ListRequest{ProjectID: projectID, Cursor: "MQ", PageSize: 10})
		Expect(err).NotTo(HaveOccurred())

		Expect(*seen).To(HaveLen(1))
		r := (*seen)[0]
		Expect(r.method).To(Equal(http.MethodGet))
		Expect(r.path).To(Equal("/audit/v1/logs"))
		Expect(r.query).To(Equal(url.Values{
			"project_id": {projectID}, "cursor": {"MQ"}, "pagesize": {"10"},
		}))
		Expect(r.header.Get("X-IK-ClientKey")).To(Equal("tok"))

		Expect(page.HasMore).To(BeTrue())
		Expect(page.NextCursor).To(Equal("Mg"))
		Expect(page.Items).To(HaveLen(1))
		Expect(page.Items[0].BatchID).To(Equal("b1"))
		Expect(page.Items[0].Sequence).To(Equal(uint64(1)))
		Expect(page.Items[0].ChainHash).To(Equal("c1"))
		Expect(page.Items[0].Data).To(MatchJSON(`[{"event":"x"}]`))
	})

	It("leaves cursor and page size to the server defaults when unset", func(ctx SpecContext) {
		c, seen := newClient(reply(`{"items":[],"has_more":false,"next_cursor":""}`))

		page, err := c.ListLogs(ctx, audit.ListRequest{ProjectID: projectID})
		Expect(err).NotTo(HaveOccurred())
		Expect(page.Items).To(BeEmpty())
		Expect((*seen)[0].query).To(Equal(url.Values{"project_id": {projectID}}))
	})

	It("fails locally without a project ID, never hitting the server", func(ctx SpecContext) {
		c, seen := newClient(reply(`{}`))

		_, err := c.ListLogs(ctx, audit.ListRequest{})
		Expect(err).To(MatchError(audit.ErrProjectIDRequired))
		_, err = c.AllManifests(ctx, audit.ListRequest{})
		Expect(err).To(MatchError(audit.ErrProjectIDRequired))
		_, err = c.JWKS(ctx, "")
		Expect(err).To(MatchError(audit.ErrProjectIDRequired))
		Expect(*seen).To(BeEmpty())
	})

	It("surfaces a project mismatch as a 403 APIError", func(ctx SpecContext) {
		c, _ := newClient(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w,
				`{"message":"the authenticated credential does not have access to this project"}`)
		})

		_, err := c.ListLogs(ctx, audit.ListRequest{ProjectID: projectID})
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue(), "err = %T (%v), want *transport.APIError", err, err)
		Expect(apiErr.StatusCode).To(Equal(http.StatusForbidden))
		Expect(apiErr.Message).To(ContainSubstring("does not have access"))
	})
})

var _ = Describe("Iterate and All", func() {
	It("follows next_cursor while has_more is true", func(ctx SpecContext) {
		c, seen := newClient(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Query().Get("cursor") {
			case "":
				_, _ = io.WriteString(w,
					`{"items":[{"manifest_id":"m1","sequence":1}],"has_more":true,"next_cursor":"Mg"}`)
			case "Mg":
				_, _ = io.WriteString(w,
					`{"items":[{"manifest_id":"m2","sequence":2}],"has_more":false,"next_cursor":""}`)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		})

		all, err := c.AllManifests(ctx, audit.ListRequest{ProjectID: projectID, PageSize: 1})
		Expect(err).NotTo(HaveOccurred())
		Expect(all).To(HaveLen(2))
		Expect(all[0].ManifestID).To(Equal("m1"))
		Expect(all[1].ManifestID).To(Equal("m2"))

		Expect(*seen).To(HaveLen(2))
		for _, r := range *seen {
			Expect(r.path).To(Equal("/audit/v1/manifests"))
			Expect(r.query.Get("pagesize")).To(Equal("1"))
		}
	})

	It("walks checkpoints and logs through their own endpoints", func(ctx SpecContext) {
		c, seen := newClient(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/audit/v1/checkpoints" {
				_, _ = io.WriteString(w, `{"items":[{"checkpoint_id":"cp2","sequence":2,"head_hash":"h2"},
					{"checkpoint_id":"cp1","sequence":1}],"has_more":false}`)
				return
			}
			_, _ = io.WriteString(w, `{"items":[{"batch_id":"b1"}],"has_more":false}`)
		})

		cps, err := c.AllCheckpoints(ctx, audit.ListRequest{ProjectID: projectID})
		Expect(err).NotTo(HaveOccurred())
		Expect(cps).To(HaveLen(2))
		Expect(cps[0].CheckpointID).To(Equal("cp2"))
		Expect(cps[0].HeadHash).To(Equal("h2"))

		it := c.IterateLogs(audit.ListRequest{ProjectID: projectID})
		Expect(it.Next(ctx)).To(BeTrue())
		Expect(it.Item().BatchID).To(Equal("b1"))
		Expect(it.Next(ctx)).To(BeFalse())
		Expect(it.Err()).NotTo(HaveOccurred())

		Expect((*seen)[0].path).To(Equal("/audit/v1/checkpoints"))
		Expect((*seen)[1].path).To(Equal("/audit/v1/logs"))
	})

	It("stops on the first page error", func(ctx SpecContext) {
		c, _ := newClient(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"message":"invalid cursor"}`)
		})

		_, err := c.AllCheckpoints(ctx, audit.ListRequest{ProjectID: projectID, Cursor: "bogus"})
		apiErr, ok := transport.AsAPIError(err)
		Expect(ok).To(BeTrue())
		Expect(apiErr.StatusCode).To(Equal(http.StatusBadRequest))
	})
})

var _ = Describe("JWKS", func() {
	// jwk renders pub the way the platform does: an EC P-256 JWK, no alg.
	jwk := func(pub *ecdsa.PublicKey, kid string) string {
		point, err := pub.Bytes()
		Expect(err).NotTo(HaveOccurred())
		enc := base64.RawURLEncoding
		return fmt.Sprintf(`{"kty":"EC","crv":"P-256","x":%q,"y":%q,"kid":%q,"use":"sig"}`,
			enc.EncodeToString(point[1:33]), enc.EncodeToString(point[33:]), kid)
	}

	It("fetches the key set and yields a key that verifies a DER signature", func(ctx SpecContext) {
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		Expect(err).NotTo(HaveOccurred())
		c, seen := newClient(reply(`{"keys":[` + jwk(&priv.PublicKey, "platform-1") + `]}`))

		set, err := c.JWKS(ctx, projectID)
		Expect(err).NotTo(HaveOccurred())
		Expect((*seen)[0].method).To(Equal(http.MethodGet))
		Expect((*seen)[0].path).To(Equal("/audit/.well-known/jwks.json"))
		Expect((*seen)[0].query).To(Equal(url.Values{"project_id": {projectID}}))

		key, ok := set.Find("platform-1")
		Expect(ok).To(BeTrue())
		Expect(key.Use).To(Equal("sig"))
		_, ok = set.Find("other")
		Expect(ok).To(BeFalse())

		pub, err := key.ECDSAPublicKey()
		Expect(err).NotTo(HaveOccurred())
		digest := sha256.Sum256([]byte("manifest"))
		sig, err := ecdsa.SignASN1(rand.Reader, priv, digest[:])
		Expect(err).NotTo(HaveOccurred())
		Expect(ecdsa.VerifyASN1(pub, digest[:], sig)).To(BeTrue())
	})

	DescribeTable("ECDSAPublicKey rejects keys it cannot use",
		func(key audit.Key, msg string) {
			_, err := key.ECDSAPublicKey()
			Expect(err).To(MatchError(ContainSubstring(msg)))
		},
		Entry("an RSA key", audit.Key{Kty: "RSA"}, "unsupported key type"),
		Entry("another curve", audit.Key{Kty: "EC", Crv: "P-384"}, "unsupported key type"),
		Entry("a non-base64url coordinate", audit.Key{Kty: "EC", Crv: "P-256", X: "!!"}, "invalid key x"),
		Entry("a short coordinate",
			audit.Key{Kty: "EC", Crv: "P-256", X: base64.RawURLEncoding.EncodeToString(make([]byte, 32)), Y: "AA"},
			"invalid key y"),
		Entry("a point off the curve", audit.Key{
			Kty: "EC", Crv: "P-256",
			X: base64.RawURLEncoding.EncodeToString(make([]byte, 32)),
			Y: base64.RawURLEncoding.EncodeToString(make([]byte, 32)),
		}, "invalid key"),
	)
})
