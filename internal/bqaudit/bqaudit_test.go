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

package bqaudit_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/indykite/indykite-sdk-go/internal/bqaudit"
	"github.com/indykite/indykite-sdk-go/test/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// configEnv is every variable the checker reads; each spec starts from all of
// them empty so the developer's or the pipeline's environment cannot leak in.
var configEnv = []string{
	"SDK_AUDIT_TABLE_NAME", "BQ_PROJECT_ID", "PROJECT_ID",
	"INDYKITE_APPLICATION_CREDENTIALS", "INDYKITE_APPLICATION_CREDENTIALS_FILE",
}

// writeFile writes data to name inside a per-spec temp directory and returns
// its path.
func writeFile(name string, data []byte) string {
	GinkgoHelper()
	path := filepath.Join(GinkgoT().TempDir(), name)
	Expect(os.WriteFile(path, data, 0o600)).To(Succeed())
	return path
}

// useFakeGoogleCredentials points Application Default Credentials at a
// throwaway service-account key, so creating the BigQuery client needs neither
// real credentials nor network access.
func useFakeGoogleCredentials() {
	GinkgoHelper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	Expect(err).NotTo(HaveOccurred())
	der, err := x509.MarshalPKCS8PrivateKey(key)
	Expect(err).NotTo(HaveOccurred())
	//nolint:gosec // G101: a throwaway key generated above, not a real credential.
	sa, err := json.Marshal(map[string]string{
		"type":           "service_account",
		"project_id":     "fake-project",
		"private_key_id": "fake-key",
		"private_key":    string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})),
		"client_email":   "fake@fake-project.iam.gserviceaccount.com",
		"client_id":      "1",
		"token_uri":      "https://oauth2.googleapis.com/token",
	})
	Expect(err).NotTo(HaveOccurred())
	testutil.SetEnv("GOOGLE_APPLICATION_CREDENTIALS", writeFile("sa.json", sa))
}

var _ = Describe("BigQuery audit checker", Serial, func() {
	BeforeEach(func() {
		for _, key := range configEnv {
			testutil.SetEnv(key, "")
		}
	})

	Describe("Enabled", func() {
		It("is off without an audit table", func() {
			Expect(bqaudit.Enabled()).To(BeFalse())
		})

		It("is on once the audit table is named", func() {
			testutil.SetEnv("SDK_AUDIT_TABLE_NAME", "audit_log")
			Expect(bqaudit.Enabled()).To(BeTrue())
		})
	})

	Describe("New", func() {
		It("requires the audit table", func(ctx SpecContext) {
			_, err := bqaudit.New(ctx)
			Expect(err).To(MatchError("bqaudit: SDK_AUDIT_TABLE_NAME not set"))
		})

		// The project is never guessed: one that does not hold the table would
		// fail every check, so the misconfiguration fails up front instead.
		It("requires the Google Cloud project of the table", func(ctx SpecContext) {
			testutil.SetEnv("SDK_AUDIT_TABLE_NAME", "audit_log")
			_, err := bqaudit.New(ctx)
			Expect(err).To(MatchError("bqaudit: BQ_PROJECT_ID not set"))
		})

		Context("with the table and its project", func() {
			BeforeEach(func() {
				testutil.SetEnv("SDK_AUDIT_TABLE_NAME", "audit_log")
				testutil.SetEnv("BQ_PROJECT_ID", "bq-project")
			})

			It("requires the app space to filter events by", func(ctx SpecContext) {
				_, err := bqaudit.New(ctx)
				Expect(err).To(MatchError(ContainSubstring("PROJECT_ID not set")))
			})

			DescribeTable("resolves the app space",
				func(ctx SpecContext, setup func(), wantAppSpace string) {
					useFakeGoogleCredentials()
					setup()

					c, err := bqaudit.New(ctx)
					Expect(err).NotTo(HaveOccurred())
					DeferCleanup(func() { Expect(c.Close()).To(Succeed()) })
					Expect(c.Table()).To(Equal("bq-project.audit_log.audit_log"))
					Expect(c.AppSpaceID()).To(Equal(wantAppSpace))
				},
				Entry("from PROJECT_ID with a token-only credential", func() {
					testutil.SetEnv("INDYKITE_APPLICATION_CREDENTIALS", "raw-agent-token")
					testutil.SetEnv("PROJECT_ID", "gid:project")
				}, "gid:project"),
				Entry("from the credential JSON before PROJECT_ID", func() {
					testutil.SetEnv("INDYKITE_APPLICATION_CREDENTIALS", `{"appSpaceId":"gid:from-cred"}`)
					testutil.SetEnv("PROJECT_ID", "gid:project")
				}, "gid:from-cred"),
				Entry("from the credential file", func() {
					testutil.SetEnv("INDYKITE_APPLICATION_CREDENTIALS_FILE",
						writeFile("cred.json", []byte(` {"appSpaceId":"gid:from-file"}`)))
				}, "gid:from-file"),
			)

			It("reports a credential JSON it cannot parse", func(ctx SpecContext) {
				testutil.SetEnv("INDYKITE_APPLICATION_CREDENTIALS", `{broken`)
				_, err := bqaudit.New(ctx)
				Expect(err).To(HaveOccurred())
			})

			It("reports a credential file it cannot read", func(ctx SpecContext) {
				testutil.SetEnv("INDYKITE_APPLICATION_CREDENTIALS_FILE",
					filepath.Join(GinkgoT().TempDir(), "missing.json"))
				_, err := bqaudit.New(ctx)
				Expect(err).To(MatchError(os.ErrNotExist))
			})
		})
	})
})
