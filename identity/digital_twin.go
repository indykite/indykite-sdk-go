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
	"google.golang.org/grpc/codes"

	"github.com/indykite/indykite-sdk-go/errors"
	identitypb "github.com/indykite/indykite-sdk-go/gen/indykite/identity/v1beta2"
)

// IntrospectToken function validates the token and returns information about it.
//
// This is a protected operation and it can be accessed only with valid agent credentials!
func (c *Client) IntrospectToken(ctx context.Context,
	token string, opts ...grpc.CallOption) (*identitypb.TokenIntrospectResponse, error) {
	if err := verifyTokenFormat(token); err != nil {
		//nolint:nilerr // If there is error, we want to ignore error and just return false here
		return &identitypb.TokenIntrospectResponse{Active: false}, nil
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.TokenIntrospect(ctx, &identitypb.TokenIntrospectRequest{Token: token}, opts...)

	switch s := errors.FromError(err); {
	case s == nil:
		return resp, nil
	case s.Code() == codes.InvalidArgument:
		return &identitypb.TokenIntrospectResponse{Active: false}, nil
	default:
		return nil, s
	}
}
