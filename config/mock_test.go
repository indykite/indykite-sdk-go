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

package config_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/indykite/indykite-sdk-go/config"
)

// recorder captures the most recent request a mock config server received, so
// tests can assert the method, path, query, If-Match header and decoded body
// without repeating the extraction boilerplate.
type recorder struct {
	query   url.Values
	body    map[string]any
	method  string
	path    string
	ifMatch string
}

// capture records the request line, query, If-Match and (when the body is JSON)
// the decoded body.
func (rec *recorder) capture(r *http.Request) {
	rec.method, rec.path = r.Method, r.URL.Path
	rec.query = r.URL.Query()
	rec.ifMatch = r.Header.Get("If-Match")
	rec.body = nil
	if raw, _ := io.ReadAll(r.Body); len(raw) > 0 {
		rec.body = map[string]any{}
		_ = json.Unmarshal(raw, &rec.body)
	}
}

// wantReq asserts the last request's method and path.
func (rec *recorder) wantReq(t *testing.T, method, path string) {
	t.Helper()
	if rec.method != method || rec.path != path {
		t.Errorf("request = %s %s, want %s %s", rec.method, rec.path, method, path)
	}
}

// wantWrite asserts that a create/update returned a populated result carrying
// the expected ETag (captured from the response header).
func wantWrite(t *testing.T, wr *config.WriteResult, wantETag string) {
	t.Helper()
	if wr.ID == "" {
		t.Errorf("write result has no ID: %+v", wr)
	}
	if wr.ETag != wantETag {
		t.Errorf("write ETag = %q, want %s", wr.ETag, wantETag)
	}
}

// crudResponses are the JSON bodies a crudHandler returns per verb. The handler
// adds the ETag headers ("v1" on create, "v2" on read/list, "v3" on update) so
// tests can assert ETag capture and If-Match propagation uniformly.
type crudResponses struct {
	created string
	list    string
	read    string
	updated string
}

// crudHandler is the shared mock for a standard ETag-guarded resource served at
// listPath: it records each request and returns the matching body, tagging
// create/read/update with distinct ETags. A GET on listPath returns the list
// body; a GET on any other path returns the read body.
func crudHandler(rec *recorder, listPath string, resp crudResponses) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec.capture(r)
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("ETag", `"v1"`)
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, resp.created)
		case http.MethodGet:
			w.Header().Set("ETag", `"v2"`)
			if r.URL.Path == listPath {
				_, _ = io.WriteString(w, resp.list)
			} else {
				_, _ = io.WriteString(w, resp.read)
			}
		case http.MethodPut:
			w.Header().Set("ETag", `"v3"`)
			_, _ = io.WriteString(w, resp.updated)
		default:
			_, _ = io.WriteString(w, `{}`)
		}
	}
}

// genericCRUDResponses are minimal responses for resources whose body fields the
// test does not inspect (it only asserts routing and ETag capture).
func genericCRUDResponses(id string) crudResponses {
	return crudResponses{
		created: fmt.Sprintf(`{"id":%q,"create_time":"t0"}`, id),
		list:    fmt.Sprintf(`{"data":[{"id":%q,"name":"n"}]}`, id),
		read:    fmt.Sprintf(`{"id":%q,"name":"n"}`, id),
		updated: fmt.Sprintf(`{"id":%q,"update_time":"t1"}`, id),
	}
}
