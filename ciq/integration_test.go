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

package ciq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	indykite "github.com/indykite/indykite-sdk-go"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/internal/bqaudit"
)

// executeAuditEventType is emitted by the platform for POST /contx-iq/v1/execute.
const executeAuditEventType = "indykite.audit.ciq.execute"

// queryFixture returns the pre-configured knowledge query (and optional input
// params) the environment provides, skipping when not configured.
func queryFixture(t *testing.T) ciq.ExecuteRequest {
	t.Helper()
	id := os.Getenv("CIQ_QUERY_ID")
	if id == "" {
		t.Skip("CIQ_QUERY_ID not set")
	}
	req := ciq.ExecuteRequest{ID: id}
	if raw := os.Getenv("CIQ_INPUT_PARAMS"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &req.InputParams); err != nil {
			t.Fatalf("CIQ_INPUT_PARAMS is not a JSON object: %v", err)
		}
	}
	return req
}

func ciqClient(t *testing.T) *ciq.Client {
	t.Helper()
	if os.Getenv("INDYKITE_APPLICATION_CREDENTIALS") == "" &&
		os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE") == "" {
		t.Skip("INDYKITE_APPLICATION_CREDENTIALS[_FILE] not set")
	}
	var opts []indykite.Option
	if base := os.Getenv("INDYKITE_BASE_URL"); base != "" {
		opts = append(opts, indykite.WithBaseURL(base))
	}
	cli, err := indykite.NewClientFromEnv(context.Background(), opts...)
	if err != nil {
		t.Fatalf("NewClientFromEnv: %v", err)
	}
	return cli.CIQ()
}

func TestIntegrationCIQExecute(t *testing.T) {
	c := ciqClient(t)
	req := queryFixture(t)
	ctx := context.Background()

	// The auditLog input param is echoed into the audit event, which lets the
	// BigQuery check below correlate this exact request (gRPC SDK convention).
	auditMarker := fmt.Sprintf("sdk-it-ciq-%d", time.Now().UnixNano())
	if req.InputParams == nil {
		req.InputParams = map[string]any{}
	}
	req.InputParams["auditLog"] = auditMarker

	resp, err := c.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	t.Logf("first page: %d records", len(resp.Data))

	if !bqaudit.Enabled() {
		t.Log("SDK_AUDIT_TABLE_NAME not set; skipping BigQuery audit-log check")
		return
	}
	checker, err := bqaudit.New(ctx)
	if err != nil {
		t.Fatalf("bqaudit.New: %v", err)
	}
	defer func() { _ = checker.Close() }()
	if err := checker.WaitForEvent(ctx, executeAuditEventType, auditMarker); err != nil {
		t.Errorf("audit log not found in BigQuery: %v", err)
	}
}

// TestIntegrationCIQAll walks all pages with a small page size so pagination is
// actually exercised whenever the query returns more than two records.
func TestIntegrationCIQAll(t *testing.T) {
	c := ciqClient(t)
	req := queryFixture(t)
	req.PageSize = 2

	records, err := c.All(context.Background(), req)
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	t.Logf("total: %d records", len(records))

	// The iterator must agree with the collected result.
	it := c.Iterate(req)
	count := 0
	for it.Next(context.Background()) {
		count++
	}
	if err := it.Err(); err != nil {
		t.Fatalf("Iterate: %v", err)
	}
	if count != len(records) {
		t.Errorf("Iterate count = %d, All = %d", count, len(records))
	}
}
