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

package audit_test

import (
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/audit"
	"github.com/indykite/indykite-sdk-go/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// deployed caches the one-time probe of the audit API. The API is served by
// triton, which the platform umbrella chart still ships disabled outside the
// tamper-proof pipeline rollout (ENG-8630), so on such an environment every
// /audit path falls through to the ingress's empty 404. The public JWKS route
// tells the two cases apart without needing the Audit permission, and the
// hermetic specs in client_test.go pin the paths, so a wrong path cannot hide
// behind this skip.
var deployed *bool

var _ = Describe("Audit (live platform)", Label("integration"), func() {
	var (
		c       *audit.Client
		project string
	)

	skipUnlessDeployed := func(ctx SpecContext) {
		GinkgoHelper()
		if deployed == nil {
			_, err := c.JWKS(ctx, project)
			apiErr, isAPI := transport.AsAPIError(err)
			ok := !isAPI || !apiErr.IsNotFound()
			deployed = &ok
		}
		if !*deployed {
			Skip("audit API not deployed on this platform (GET /audit/.well-known/jwks.json is 404: " +
				"triton is not enabled there yet)")
		}
	}

	BeforeEach(func(ctx SpecContext) {
		if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
			os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
			Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
		}
		project = os.Getenv("PROJECT_ID")
		if project == "" {
			Skip("PROJECT_ID not set")
		}
		var opts []indykite.Option
		if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
			opts = append(opts, indykite.WithBaseURL(base))
		}
		cli, err := indykite.NewClientFromEnv(ctx, opts...)
		Expect(err).NotTo(HaveOccurred())
		c = cli.Audit()
		skipUnlessDeployed(ctx)
	})

	// skipWithoutPermission skips the spec when the App Agent lacks the Audit
	// API permission: that is a per-agent grant rather than a test input.
	skipWithoutPermission := func(err error) {
		GinkgoHelper()
		if apiErr, ok := transport.AsAPIError(err); ok && apiErr.IsUnauthorized() {
			Skip(fmt.Sprintf("App Agent cannot read the audit log of PROJECT_ID: %v", err))
		}
	}

	It("pages logs and manifests in lockstep", func(ctx SpecContext) {
		logs, err := c.ListLogs(ctx, audit.ListRequest{ProjectID: project, PageSize: 5})
		skipWithoutPermission(err)
		Expect(err).NotTo(HaveOccurred())
		manifests, err := c.ListManifests(ctx, audit.ListRequest{ProjectID: project, PageSize: 5})
		Expect(err).NotTo(HaveOccurred())

		GinkgoWriter.Printf("logs=%d manifests=%d\n", len(logs.Items), len(manifests.Items))
		Expect(logs.Items).To(HaveLen(len(manifests.Items)))
		for i := range logs.Items {
			Expect(logs.Items[i].ManifestID).To(Equal(manifests.Items[i].ManifestID))
			Expect(logs.Items[i].Sequence).To(Equal(manifests.Items[i].Sequence))
		}
	})

	It("lists checkpoints newest first", func(ctx SpecContext) {
		page, err := c.ListCheckpoints(ctx, audit.ListRequest{ProjectID: project})
		skipWithoutPermission(err)
		Expect(err).NotTo(HaveOccurred())
		for i := 1; i < len(page.Items); i++ {
			Expect(page.Items[i].Sequence).To(BeNumerically("<=", page.Items[i-1].Sequence))
		}
	})

	It("publishes an EC P-256 verification key", func(ctx SpecContext) {
		set, err := c.JWKS(ctx, project)
		Expect(err).NotTo(HaveOccurred())
		Expect(set.Keys).NotTo(BeEmpty())
		_, err = set.Keys[0].ECDSAPublicKey()
		Expect(err).NotTo(HaveOccurred())
	})
})
