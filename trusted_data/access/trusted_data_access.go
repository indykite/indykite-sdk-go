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

package access

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/errors"
	accesspb "github.com/indykite/indykite-sdk-go/gen/indykite/trusted_data/access/v1beta1"
)

func (c *Client) AccessConsentedData(
	ctx context.Context,
	req *accesspb.AccessConsentedDataRequest,
	opts ...grpc.CallOption,
) (*accesspb.AccessConsentedDataResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call AccessConsentedData")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.AccessConsentedData(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) GrantConsent(
	ctx context.Context,
	req *accesspb.GrantConsentRequest,
	opts ...grpc.CallOption,
) (*accesspb.GrantConsentResponse, error) {
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

func (c *Client) RevokeConsent(
	ctx context.Context,
	req *accesspb.RevokeConsentRequest,
	opts ...grpc.CallOption,
) (*accesspb.RevokeConsentResponse, error) {
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
