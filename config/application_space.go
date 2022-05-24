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

func (c *Client) CreateApplicationSpace(
	ctx context.Context,
	request *configpb.CreateApplicationSpaceRequest,
	opts ...grpc.CallOption) (*configpb.CreateApplicationSpaceResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateApplicationSpace(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ReadApplicationSpace(
	ctx context.Context,
	request *configpb.ReadApplicationSpaceRequest,
	opts ...grpc.CallOption) (*configpb.ReadApplicationSpaceResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadApplicationSpace(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ListApplicationSpaces(
	ctx context.Context,
	request *configpb.ListApplicationSpacesRequest,
	opts ...grpc.CallOption,
) (configpb.ConfigManagementAPI_ListApplicationSpacesClient, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ListApplicationSpaces(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) UpdateApplicationSpace(
	ctx context.Context,
	request *configpb.UpdateApplicationSpaceRequest,
	opts ...grpc.CallOption) (*configpb.UpdateApplicationSpaceResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.UpdateApplicationSpace(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) DeleteApplicationSpace(
	ctx context.Context,
	request *configpb.DeleteApplicationSpaceRequest,
	opts ...grpc.CallOption) (*configpb.DeleteApplicationSpaceResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteApplicationSpace(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
