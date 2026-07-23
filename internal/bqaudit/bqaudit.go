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

// Package bqaudit verifies, from integration tests, that platform audit events
// landed in the BigQuery audit-log table.
//
// It is enabled by setting SDK_AUDIT_TABLE_NAME (the dataset/table name, e.g.
// "audit_log"). BQ_PROJECT_ID overrides the Google Cloud project that hosts the
// table. Queries authenticate with Application Default Credentials.
package bqaudit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// Audit events are delivered asynchronously; poll before declaring a miss.
const (
	pollAttempts = 6
	pollInterval = 5 * time.Second
	rowLimit     = 50
)

// Enabled reports whether audit-log verification is configured; tests skip the
// check when it is not.
func Enabled() bool { return os.Getenv("SDK_AUDIT_TABLE_NAME") != "" }

// Checker queries the audit-log table for events of the caller's app space.
type Checker struct {
	client     *bigquery.Client
	table      string
	appSpaceID string
}

func gcpProject() string {
	if p := os.Getenv("BQ_PROJECT_ID"); p != "" {
		return p
	}
	return "indykite-rc"
}

// appSpaceFromCredentials reads the app-space gid out of the App Agent
// credential JSON when it carries one, falling back to PROJECT_ID (the project
// gid IS the app-space gid) — token-only credential files carry no appSpaceId.
func appSpaceFromCredentials() (string, error) {
	raw := os.Getenv("INDYKITE_APPLICATION_CREDENTIALS")
	if raw == "" {
		if path := os.Getenv("INDYKITE_APPLICATION_CREDENTIALS_FILE"); path != "" {
			b, err := os.ReadFile(path) //nolint:gosec // test helper reading the caller-supplied credentials file
			if err != nil {
				return "", err
			}
			raw = string(b)
		}
	}
	var cred struct {
		AppSpaceID string `json:"appSpaceId"`
	}
	if raw != "" && strings.HasPrefix(strings.TrimSpace(raw), "{") {
		if err := json.Unmarshal([]byte(raw), &cred); err != nil {
			return "", err
		}
	}
	if cred.AppSpaceID == "" {
		cred.AppSpaceID = os.Getenv("PROJECT_ID")
	}
	if cred.AppSpaceID == "" {
		return "", errors.New(
			"bqaudit: no appSpaceId in INDYKITE_APPLICATION_CREDENTIALS[_FILE] and PROJECT_ID not set")
	}
	return cred.AppSpaceID, nil
}

// New builds a Checker for the configured audit table.
func New(ctx context.Context) (*Checker, error) {
	table := os.Getenv("SDK_AUDIT_TABLE_NAME")
	if table == "" {
		return nil, errors.New("bqaudit: SDK_AUDIT_TABLE_NAME not set")
	}
	appSpaceID, err := appSpaceFromCredentials()
	if err != nil {
		return nil, err
	}
	project := gcpProject()
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	return &Checker{
		client:     client,
		table:      fmt.Sprintf("%s.%s.%s", project, table, table),
		appSpaceID: appSpaceID,
	}, nil
}

// Close releases the underlying BigQuery client.
func (c *Checker) Close() error { return c.client.Close() }

type row struct {
	PublishTime time.Time `bigquery:"publish_time"`
	Data        string    `bigquery:"data"`
}

// query returns today's newest audit rows of the given event type for the
// checker's app space.
func (c *Checker) query(ctx context.Context, eventType string) ([]row, error) {
	q := c.client.Query(fmt.Sprintf(
		"SELECT publish_time, data FROM `%s`"+
			" WHERE TIMESTAMP_TRUNC(publish_time, DAY) = TIMESTAMP(@day)"+
			" AND JSON_VALUE(attributes[\"ce-appspaceid\"]) = @app_space"+
			" AND JSON_VALUE(data.eventType) = @event_type"+
			" ORDER BY publish_time DESC LIMIT %d", c.table, rowLimit))
	q.Parameters = []bigquery.QueryParameter{
		{Name: "day", Value: time.Now().UTC().Format("2006-01-02")},
		{Name: "app_space", Value: c.appSpaceID},
		{Name: "event_type", Value: eventType},
	}
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}
	var rows []row
	for {
		var r row
		err := it.Next(&r)
		if errors.Is(err, iterator.Done) {
			return rows, nil
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, r)
	}
}

// WaitForEvent polls the audit table until an event of the given type whose
// payload contains marker appears, or the poll budget runs out.
func (c *Checker) WaitForEvent(ctx context.Context, eventType, marker string) error {
	var lastErr error
	for attempt := range pollAttempts {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(pollInterval):
			}
		}
		rows, err := c.query(ctx, eventType)
		if err != nil {
			lastErr = err
			continue
		}
		for _, r := range rows {
			if strings.Contains(r.Data, marker) {
				return nil
			}
		}
		lastErr = fmt.Errorf("bqaudit: no %q event containing %q in the newest %d rows",
			eventType, marker, rowLimit)
	}
	return lastErr
}
