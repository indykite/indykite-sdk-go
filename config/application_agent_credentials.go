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

package config

import (
	"context"

	"google.golang.org/grpc"

	configpb "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"

	"github.com/indykite/jarvis-sdk-go/errors"
)

func (c *Client) RegisterApplicationAgentCredential(
	ctx context.Context,
	request *configpb.RegisterApplicationAgentCredentialRequest,
	opts ...grpc.CallOption,
) (*configpb.RegisterApplicationAgentCredentialResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.RegisterApplicationAgentCredential(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ReadApplicationAgentCredential(
	ctx context.Context,
	request *configpb.ReadApplicationAgentCredentialRequest,
	opts ...grpc.CallOption,
) (*configpb.ReadApplicationAgentCredentialResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadApplicationAgentCredential(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) DeleteApplicationAgentCredential(
	ctx context.Context,
	request *configpb.DeleteApplicationAgentCredentialRequest,
	opts ...grpc.CallOption,
) (*configpb.DeleteApplicationAgentCredentialResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteApplicationAgentCredential(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
