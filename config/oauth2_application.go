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

func (c *Client) ReadOAuth2Application(
	ctx context.Context,
	request *configpb.ReadOAuth2ApplicationRequest,
	opts ...grpc.CallOption) (*configpb.ReadOAuth2ApplicationResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadOAuth2Application(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) CreateOAuth2Application(
	ctx context.Context,
	request *configpb.CreateOAuth2ApplicationRequest,
	opts ...grpc.CallOption) (*configpb.CreateOAuth2ApplicationResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateOAuth2Application(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) UpdateOAuth2Application(
	ctx context.Context,
	request *configpb.UpdateOAuth2ApplicationRequest,
	opts ...grpc.CallOption) (*configpb.UpdateOAuth2ApplicationResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.UpdateOAuth2Application(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) DeleteOAuth2Application(
	ctx context.Context,
	request *configpb.DeleteOAuth2ApplicationRequest,
	opts ...grpc.CallOption) (*configpb.DeleteOAuth2ApplicationResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteOAuth2Application(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
