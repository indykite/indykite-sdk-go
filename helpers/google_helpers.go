// Copyright (c) 2024 IndyKite
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

// Package helpers is containing the logic which can be used to connect to BigQuery and run some request to audit log.
package helpers

import (
	"context"
	"encoding/json"
	e "errors"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/indykite/indykite-sdk-go/errors"
)

const (
	appSpaceIDFromAttribute      = "JSON_VALUE(attributes[\"ce-appspaceid\"])"
	appSpaceIDFromContainerPath  = "JSON_VALUE(data[\"data\"][\"containersPath\"][\"applicationSpaceId\"])"
	auditLogIdentifierInputField = "JSON_VALUE(data.data.request.inputParams.auditLog.stringValue)"
	dataType                     = "JSON_VALUE(data.data[\"@type\"])"
	eventSource                  = "JSON_VALUE(data.eventSource)"
	eventType                    = "JSON_VALUE(data.eventType)"
)

// FilterFields is a collection of variable to set up the actual query more easily.
type FilterFields struct {
	AppSpaceID         string
	AppSpaceIDConfig   string
	AuditLogIdentifier string
	ChangeType         string
	EventSource        string
	EventType          string
	RowNumber          string
}

// ReturnValues maps the response from bigquery.
// The getReturnFields specifies what to expect from bigquery and all of those fields should be represented here.
type ReturnValues struct {
	PublishTime time.Time `bigquery:"publish_time"`
	Data        string    `bigquery:"data"`
}

type bqQueryParameters struct {
	ReturnFields string `default:"*"`
	tableName    string
	filters      string
	RowNumber    string `default:"1"`
}

func getProjectID() string {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return "jarvis-dev-268314"
	}

	return projectID
}

func getReturnFields() string {
	return "publish_time, data "
}

func getFilterFields(filter *FilterFields) string {
	today := time.Now()
	var f strings.Builder
	_, _ = f.WriteString("WHERE TIMESTAMP_TRUNC(publish_time, DAY) = TIMESTAMP(\"" + today.Format("2006-01-02") + "\")")
	if filter.AppSpaceID != "" {
		_, _ = f.WriteString(" AND " + appSpaceIDFromAttribute + " LIKE '" + filter.AppSpaceID + "'")
	}
	if filter.AppSpaceIDConfig != "" {
		_, _ = f.WriteString(" AND " + appSpaceIDFromContainerPath + " LIKE '" + filter.AppSpaceIDConfig + "'")
	}
	if filter.AuditLogIdentifier != "" {
		_, _ = f.WriteString(" AND " + auditLogIdentifierInputField + " LIKE '\"" + filter.AuditLogIdentifier + "\"'")
	}
	if filter.ChangeType != "" {
		_, _ = f.WriteString(" AND " + dataType + " LIKE '" + filter.ChangeType + "'")
	}
	if filter.EventSource != "" {
		_, _ = f.WriteString(" AND " + eventSource + " LIKE '" + filter.EventSource + "'")
	}
	if filter.EventType != "" {
		_, _ = f.WriteString(" AND " + eventType + " LIKE '" + filter.EventType + "'")
	}
	return f.String()
}

// getTableName sets the bigquery table name. Each environment can have different table name. The default name is
// `<env>_audit_log` if it was called differently, it can be configured with an environment variable
// `SDK_AUDIT_TABLE_NAME`.
// Later, there is no need to have the `<env>_audit_log` format (when we switch to the new projects).
func getTableName() string {
	var (
		auditTableName = "audit_log"
		projectID      = getProjectID()
	)

	if name := os.Getenv("SDK_AUDIT_TABLE_NAME"); name != "" {
		auditTableName = name
	} else {
		env := "stg"
		runEnv := os.Getenv("RUN_ENV")
		if runEnv == "develop" {
			env = "dev"
		}
		auditTableName = fmt.Sprintf("%s_%s", env, auditTableName)
	}
	return projectID + "." + auditTableName + "." + auditTableName
}

func processBqResults(iter *bigquery.RowIterator) (ReturnValues, error) {
	var row ReturnValues
	for {
		err := iter.Next(&row)
		if e.Is(err, iterator.Done) {
			return row, nil
		}
		if err != nil {
			return row, err
		}
	}
}

// FillFilterFieldsFromEnvironment creates a FilterFields object and prefills the App space information from the
// App credential environment variable.
func FillFilterFieldsFromEnvironment() (FilterFields, error) {
	applicationCred := os.Getenv("INDYKITE_APPLICATION_CREDENTIALS")
	var jsonMap map[string]any
	var f FilterFields
	if applicationCred == "" {
		return f, errors.NewInvalidArgumentError("Missing INDYKITE_APPLICATION_CREDENTIALS environment variable")
	}
	err := json.Unmarshal([]byte(applicationCred), &jsonMap)
	if err != nil {
		return f, err
	}
	appSpaceID := fmt.Sprintf("%v", jsonMap["appSpaceId"])
	f = FilterFields{
		AppSpaceID: appSpaceID,
		RowNumber:  "1",
	}

	return f, nil
}

// BqClient creates a client connection to BigQuery.
func BqClient(ctx context.Context) (*bigquery.Client, error) {
	client, err := bigquery.NewClient(ctx, getProjectID())
	if err != nil {
		return nil, err
	}
	return client, nil
}

// QueryAuditLog sends a query request to the BigQuery with the given parameters and returns with its result
// It returns with error if the RUN_ENV environment variable is not set or not develop/staging.
func QueryAuditLog(ctx context.Context, client *bigquery.Client, fFields *FilterFields) (ReturnValues, error) {
	var params bqQueryParameters
	params.tableName = getTableName()
	if fFields.RowNumber != "" {
		params.RowNumber = fFields.RowNumber
	}
	params.filters = getFilterFields(fFields)
	params.ReturnFields = getReturnFields()
	select {
	case <-ctx.Done():
		return ReturnValues{}, e.New("the context is done")
	case <-time.After(time.Second):
		query := client.Query(
			"SELECT " +
				params.ReturnFields + " FROM `" +
				params.tableName +
				"` " +
				params.filters +
				" ORDER BY publish_time DESC LIMIT " +
				params.RowNumber +
				";")
		iter, err := query.Read(ctx)
		if err != nil {
			return ReturnValues{}, err
		}
		res, err := processBqResults(iter)
		if err != nil {
			return ReturnValues{}, err
		}
		return res, nil
	}
}
