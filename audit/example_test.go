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
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/audit"
)

// Walk a project's signed chain manifests page by page.
func ExampleClient_IterateManifests() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}

	it := cli.Audit().IterateManifests(audit.ListRequest{ProjectID: os.Getenv("PROJECT_ID")})
	for it.Next(ctx) {
		m := it.Item()
		_ = m.PrevHash // links to the previous manifest's HeadHash
	}
	if err := it.Err(); err != nil {
		return
	}
}

// Verify the newest checkpoint against the published key. A checkpoint is
// signed over sha256("project-checkpoint|<project_id>|<sequence>|<head_hash>|<created_at>");
// the signature is standard base64 of an ASN.1 DER ECDSA P-256 signature.
func ExampleKey_ECDSAPublicKey() {
	ctx := context.Background()
	cli, err := indykite.NewClientFromEnv(ctx)
	if err != nil {
		return
	}
	projectID := os.Getenv("PROJECT_ID")

	keys, err := cli.Audit().JWKS(ctx, projectID)
	if err != nil {
		return
	}
	page, err := cli.Audit().ListCheckpoints(ctx, audit.ListRequest{ProjectID: projectID, PageSize: 1})
	if err != nil || len(page.Items) == 0 {
		return
	}
	cp := page.Items[0] // newest first

	key, ok := keys.Find(cp.KID)
	if !ok {
		return
	}
	pub, err := key.ECDSAPublicKey()
	if err != nil {
		return
	}
	digest := sha256.Sum256(fmt.Appendf(nil, "project-checkpoint|%s|%d|%s|%s",
		cp.ProjectID, cp.Sequence, cp.HeadHash, cp.CreatedAt))
	sig, err := base64.StdEncoding.DecodeString(cp.Signature)
	if err != nil {
		return
	}
	fmt.Println(ecdsa.VerifyASN1(pub, digest[:], sig))
}
