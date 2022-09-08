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

package identity

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/jarvis-sdk-go/errors"
	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta1"
)

func (c *Client) CheckConsentChallenge(
	ctx context.Context,
	req *identitypb.CheckConsentChallengeRequest,
	opts ...grpc.CallOption,
) (*identitypb.CheckConsentChallengeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call CheckConsentChallenge")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CheckConsentChallenge(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) CreateConsentVerifier(
	ctx context.Context,
	req *identitypb.CreateConsentVerifierRequest,
	opts ...grpc.CallOption,
) (*identitypb.CreateConsentVerifierResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call CreateConsentVerifier")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateConsentVerifier(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
