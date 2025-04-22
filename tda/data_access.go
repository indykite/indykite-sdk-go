// Copyright (c) 2022 IndyKite
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

package tda

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/errors"
	tdapb "github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1"
)

// DataAccess returns resp.
func (c *Client) DataAccess(
	ctx context.Context,
	req *tdapb.DataAccessRequest,
	opts ...grpc.CallOption,
) (*tdapb.DataAccessResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call DataAccess")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DataAccess(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

// GrantConsent returns resp.
func (c *Client) GrantConsent(
	ctx context.Context,
	req *tdapb.GrantConsentRequest,
	opts ...grpc.CallOption,
) (*tdapb.GrantConsentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call GrantConsent")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.GrantConsent(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

// RevokeConsent returns resp.
func (c *Client) RevokeConsent(
	ctx context.Context,
	req *tdapb.RevokeConsentRequest,
	opts ...grpc.CallOption,
) (*tdapb.RevokeConsentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call RevokeConsent")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.RevokeConsent(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

// ListConsents returns resp.
func (c *Client) ListConsents(
	ctx context.Context,
	req *tdapb.ListConsentsRequest,
	opts ...grpc.CallOption,
) (*tdapb.ListConsentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call ListConsents")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ListConsents(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
