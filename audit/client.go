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

// Package audit is the client for the IndyKite tamper-proof audit log API
// (/audit/v1/* and /audit/.well-known/jwks.json). It runs on the runtime plane
// (App Agent token) and is a thin facade over a *transport.Client.
//
// The App Agent needs the Audit API permission (config.PermissionAudit) and
// can only read the project it belongs to.
//
//	a := audit.NewClient(client)
//	logs, _ := a.AllLogs(ctx, audit.ListRequest{ProjectID: projectID})
//	// or page-by-page:
//	it := a.IterateManifests(audit.ListRequest{ProjectID: projectID})
//	for it.Next(ctx) { use(it.Item()) }
//	// the public keys the signatures verify against:
//	keys, _ := a.JWKS(ctx, projectID)
//	// the same keys without any credential, e.g. from a verifier:
//	keys, _ = audit.FetchJWKS(ctx, baseURL, projectID)
//	// or, for repeated lookups, a client that sends no credential:
//	pub, _ := audit.NewPublicClient(baseURL)
//	keys, _ = pub.JWKS(ctx, projectID)
package audit

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/transport"
)

const (
	pathLogs        = "/audit/v1/logs"
	pathManifests   = "/audit/v1/manifests"
	pathCheckpoints = "/audit/v1/checkpoints"
	pathJWKS        = "/audit/.well-known/jwks.json"

	p256CoordinateSize = 32
)

// ErrProjectIDRequired is returned when a call is made without a project ID.
// The platform would reject such a request with 400 anyway; failing locally
// saves the round trip.
var ErrProjectIDRequired = errors.New("audit: project ID is required")

// ErrBaseURLRequired is returned by NewPublicClient and FetchJWKS when no
// platform endpoint is given: with no credential there is nothing to default
// it from.
var ErrBaseURLRequired = errors.New("audit: base URL is required")

// Client calls the audit log API.
type Client struct {
	t *transport.Client
}

// NewClient builds an audit client over the shared transport.
func NewClient(t *transport.Client) *Client {
	return &Client{t: t}
}

// ListLogs returns one page of the project's signed chain batches (the audit
// events), in sequence order. For every page use IterateLogs or AllLogs.
func (c *Client) ListLogs(ctx context.Context, req ListRequest) (*Page[LogEntry], error) {
	return list[LogEntry](ctx, c.t, pathLogs, req)
}

// IterateLogs lazily walks every page of the project's chain batches.
func (c *Client) IterateLogs(req ListRequest) *transport.Iterator[LogEntry] {
	return iterate[LogEntry](c.t, pathLogs, req)
}

// AllLogs collects every chain batch of the project across all pages.
func (c *Client) AllLogs(ctx context.Context, req ListRequest) ([]LogEntry, error) {
	return c.IterateLogs(req).Collect(ctx)
}

// ListManifests returns one page of the project's signed chain manifests, in
// sequence order: the metadata needed to verify chain integrity without
// fetching the batch payloads. Manifests page in lockstep with ListLogs.
func (c *Client) ListManifests(ctx context.Context, req ListRequest) (*Page[Manifest], error) {
	return list[Manifest](ctx, c.t, pathManifests, req)
}

// IterateManifests lazily walks every page of the project's chain manifests.
func (c *Client) IterateManifests(req ListRequest) *transport.Iterator[Manifest] {
	return iterate[Manifest](c.t, pathManifests, req)
}

// AllManifests collects every chain manifest of the project across all pages.
func (c *Client) AllManifests(ctx context.Context, req ListRequest) ([]Manifest, error) {
	return c.IterateManifests(req).Collect(ctx)
}

// ListCheckpoints returns one page of the project's signed checkpoints,
// newest first (unlike logs and manifests).
func (c *Client) ListCheckpoints(ctx context.Context, req ListRequest) (*Page[Checkpoint], error) {
	return list[Checkpoint](ctx, c.t, pathCheckpoints, req)
}

// IterateCheckpoints lazily walks every page of the project's checkpoints,
// newest first.
func (c *Client) IterateCheckpoints(req ListRequest) *transport.Iterator[Checkpoint] {
	return iterate[Checkpoint](c.t, pathCheckpoints, req)
}

// AllCheckpoints collects every checkpoint of the project across all pages.
func (c *Client) AllCheckpoints(ctx context.Context, req ListRequest) ([]Checkpoint, error) {
	return c.IterateCheckpoints(req).Collect(ctx)
}

// JWKS returns the keys the project's batch, manifest and checkpoint
// signatures verify against (GET /audit/.well-known/jwks.json): every key that
// ever signed for the project, current and retired, found by kid with
// KeySet.Find. The endpoint is public and needs no credential.
//
// A project with no signed data yet has no keys, and a project's first key, or
// one added by a key rotation, can take up to two minutes to be listed. Keep
// the set next to an exported trail so each item can be verified later.
//
// The platform answers 503 with a Retry-After header while it loads a large
// key set; the transport retries those on its own, within its retry policy.
func (c *Client) JWKS(ctx context.Context, projectID string) (*KeySet, error) {
	if projectID == "" {
		return nil, ErrProjectIDRequired
	}
	var out KeySet
	path := pathJWKS + "?" + url.Values{"project_id": {projectID}}.Encode()
	if err := c.t.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// NewPublicClient builds an audit client that sends no credential, for
// verifiers that hold an exported trail but no App Agent. Only JWKS works on
// it; every other call answers 401. baseURL is the platform endpoint, e.g.
// https://eu.api.indykite.com, and always wins over a base URL or region in
// opts, which otherwise tune the transport (retry policy, HTTP client).
func NewPublicClient(baseURL string, opts ...transport.Option) (*Client, error) {
	if baseURL == "" {
		return nil, ErrBaseURLRequired
	}
	// Appended last so that it takes precedence over any WithBaseURL/WithRegion
	// in opts; clipped first so the caller's slice is never written to.
	opts = append(slices.Clip(opts), transport.WithBaseURL(baseURL))
	t, err := transport.NewClient(auth.Anonymous(), opts...)
	if err != nil {
		return nil, err
	}
	return NewClient(t), nil
}

// FetchJWKS reads a project's key set in one call without any credential; see
// NewPublicClient for the arguments. Callers that look keys up repeatedly
// should keep a NewPublicClient instead.
func FetchJWKS(ctx context.Context, baseURL, projectID string, opts ...transport.Option) (*KeySet, error) {
	// Checked before the transport is built so a doomed call costs nothing.
	if projectID == "" {
		return nil, ErrProjectIDRequired
	}
	c, err := NewPublicClient(baseURL, opts...)
	if err != nil {
		return nil, err
	}
	return c.JWKS(ctx, projectID)
}

// Find returns the key with the given kid (the KID of a batch, manifest or
// checkpoint), or false when the set has none.
func (s *KeySet) Find(kid string) (Key, bool) {
	for _, k := range s.Keys {
		if k.Kid == kid {
			return k, true
		}
	}
	return Key{}, false
}

// ECDSAPublicKey decodes an EC P-256 key for use with ecdsa.VerifyASN1.
func (k *Key) ECDSAPublicKey() (*ecdsa.PublicKey, error) {
	if k.Kty != "EC" || k.Crv != "P-256" {
		return nil, fmt.Errorf("audit: unsupported key type %q/%q, want EC/P-256", k.Kty, k.Crv)
	}
	x, err := decodeCoordinate(k.X)
	if err != nil {
		return nil, fmt.Errorf("audit: invalid key x: %w", err)
	}
	y, err := decodeCoordinate(k.Y)
	if err != nil {
		return nil, fmt.Errorf("audit: invalid key y: %w", err)
	}
	point := make([]byte, 0, 1+2*p256CoordinateSize)
	point = append(point, 4) // uncompressed point marker
	point = append(point, x...)
	point = append(point, y...)
	pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), point)
	if err != nil {
		return nil, fmt.Errorf("audit: invalid key: %w", err)
	}
	return pub, nil
}

func decodeCoordinate(v string) ([]byte, error) {
	b, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return nil, err
	}
	if len(b) != p256CoordinateSize {
		return nil, fmt.Errorf("got %d bytes, want %d", len(b), p256CoordinateSize)
	}
	return b, nil
}

func list[T any](ctx context.Context, t *transport.Client, path string, req ListRequest) (*Page[T], error) {
	if req.ProjectID == "" {
		return nil, ErrProjectIDRequired
	}
	q := url.Values{"project_id": {req.ProjectID}}
	if req.Cursor != "" {
		q.Set("cursor", req.Cursor)
	}
	if req.PageSize > 0 {
		q.Set("pagesize", strconv.Itoa(req.PageSize))
	}
	var out Page[T]
	if err := t.Do(ctx, http.MethodGet, path+"?"+q.Encode(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func iterate[T any](t *transport.Client, path string, req ListRequest) *transport.Iterator[T] {
	return transport.NewIterator(func(ctx context.Context, token string) (transport.Page[T], error) {
		pageReq := req
		if token != "" {
			pageReq.Cursor = token
		}
		page, err := list[T](ctx, t, path, pageReq)
		if err != nil {
			return transport.Page[T]{}, err
		}
		next := ""
		if page.HasMore {
			next = page.NextCursor
		}
		return transport.Page[T]{Items: page.Items, NextToken: next}, nil
	})
}
