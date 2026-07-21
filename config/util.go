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

package config

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/indykite/indykite-sdk-go/transport"
)

// ReadOption customizes a read (by name requires a location; a version pins a
// specific config version).
type ReadOption func(url.Values)

// WithLocation sets the location (project/organization gid) for name lookups.
func WithLocation(location string) ReadOption {
	return func(q url.Values) { q.Set("location", location) }
}

// WithVersion pins a specific config version.
func WithVersion(version int64) ReadOption {
	return func(q url.Values) { q.Set("version", strconv.FormatInt(version, 10)) }
}

// ListOption customizes a list call.
type ListOption func(url.Values)

// WithSearch filters by a search term across name/display name/description.
func WithSearch(term string) ListOption {
	return func(q url.Values) { q.Set("search", term) }
}

// WithFullFetch returns full data instead of metadata only.
func WithFullFetch() ListOption {
	return func(q url.Values) { q.Set("full_fetch", "true") }
}

func appendQuery(path string, q url.Values) string {
	if len(q) == 0 {
		return path
	}
	return path + "?" + q.Encode()
}

// write performs a create/update/delete and captures the ETag header.
func write(
	ctx context.Context,
	t *transport.Client,
	method, path string,
	body any,
	opts ...transport.CallOption,
) (*WriteResult, error) {
	var wr WriteResult
	resp, err := t.DoResp(ctx, method, path, body, &wr, opts...)
	if err != nil {
		return nil, err
	}
	wr.ETag = resp.Header.Get(headerETag)
	return &wr, nil
}

// resourcePath builds "<base>/<id>" with id percent-escaped for use as a single
// path segment, so a gid or config name never alters the request path.
func resourcePath(base, id string) string {
	return base + "/" + url.PathEscape(id)
}

// Generic CRUD helpers shared by all resource sub-APIs.

// listResource lists resources of type T at path with the given query.
func listResource[T any](ctx context.Context, t *transport.Client, path string, q url.Values) ([]T, error) {
	var out listResponse[T]
	if err := t.Do(ctx, http.MethodGet, appendQuery(path, q), nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// readResource reads one resource of type T by id and populates its ETag from
// the response header. PT is *T constrained to expose setETag (via Versioned).
func readResource[T any, PT interface {
	*T
	etagSettable
}](ctx context.Context, t *transport.Client, path, id string, q url.Values) (PT, error) {
	var out T
	p := PT(&out)
	resp, err := t.DoResp(ctx, http.MethodGet, appendQuery(resourcePath(path, id), q), nil, p)
	if err != nil {
		return nil, err
	}
	p.setETag(resp.Header.Get(headerETag))
	return p, nil
}

// deleteResource deletes a resource by id, optionally guarded by an ETag.
func deleteResource(ctx context.Context, t *transport.Client, path, id, etag string) error {
	return t.Do(ctx, http.MethodDelete, resourcePath(path, id), nil, nil, ifMatch(etag)...)
}

// projectListQuery builds the query for a project-scoped list.
func projectListQuery(projectID string, opts []ListOption) url.Values {
	q := url.Values{}
	q.Set("project_id", projectID)
	for _, o := range opts {
		o(q)
	}
	return q
}

// orgListQuery builds the query for an organization-scoped list.
func orgListQuery(organizationID string, opts []ListOption) url.Values {
	q := url.Values{}
	q.Set("organization_id", organizationID)
	for _, o := range opts {
		o(q)
	}
	return q
}

// readOptsQuery builds the query for a read.
func readOptsQuery(opts []ReadOption) url.Values {
	q := url.Values{}
	for _, o := range opts {
		o(q)
	}
	return q
}
