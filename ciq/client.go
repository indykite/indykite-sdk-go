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

// Package ciq is the client for the IndyKite ContX IQ query API
// (/contx-iq/v1/execute). It runs on the runtime plane (App Agent token) and is
// a thin facade over a *transport.Client.
//
//	q := ciq.NewClient(client)
//	rows, _ := q.All(ctx, ciq.ExecuteRequest{ID: "get-servers"})  // all pages
//	// or page-by-page:
//	it := q.Iterate(ciq.ExecuteRequest{ID: "get-servers"})
//	for it.Next(ctx) { use(it.Item()) }
package ciq

import (
	"context"
	"net/http"
	"strconv"

	"github.com/indykite/indykite-sdk-go/transport"
)

const (
	pathExecute     = "/contx-iq/v1/execute"
	defaultPageSize = 100
	firstPage       = 1
)

// Client calls the ContX IQ Execute endpoint.
type Client struct {
	t *transport.Client
}

// NewClient builds a CIQ client over the shared transport.
func NewClient(t *transport.Client) *Client {
	return &Client{t: t}
}

// Execute runs a single page of a query. PageToken < 1 returns the first page.
// For multi-page results use Iterate or All.
func (c *Client) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	var out ExecuteResponse
	if err := c.t.Do(ctx, http.MethodPost, pathExecute, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Iterate lazily walks every page of a query. CIQ has no server-returned next
// token: a page is the last one when it returns fewer than PageSize records, so
// a concrete PageSize is required and defaults to 100.
func (c *Client) Iterate(req ExecuteRequest) *transport.Iterator[Record] {
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return transport.NewIterator(func(ctx context.Context, token string) (transport.Page[Record], error) {
		page := firstPage
		if token != "" {
			p, err := strconv.Atoi(token)
			if err != nil {
				return transport.Page[Record]{}, err
			}
			page = p
		}

		pageReq := req
		pageReq.PageToken = page
		pageReq.PageSize = pageSize

		resp, err := c.Execute(ctx, pageReq)
		if err != nil {
			return transport.Page[Record]{}, err
		}

		next := ""
		if len(resp.Data) == pageSize {
			next = strconv.Itoa(page + 1)
		}
		return transport.Page[Record]{Items: resp.Data, NextToken: next}, nil
	})
}

// All collects every record across all pages.
func (c *Client) All(ctx context.Context, req ExecuteRequest) ([]Record, error) {
	return c.Iterate(req).Collect(ctx)
}
