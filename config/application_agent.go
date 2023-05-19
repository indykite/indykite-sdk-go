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

	"github.com/indykite/indykite-sdk-go/errors"
	configpb "github.com/indykite/indykite-sdk-go/gen/indykite/config/v1beta1"
)

func (c *Client) ReadApplicationAgent(
	ctx context.Context,
	request *configpb.ReadApplicationAgentRequest,
	opts ...grpc.CallOption) (*configpb.ReadApplicationAgentResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadApplicationAgent(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ListApplicationAgents(
	ctx context.Context,
	request *configpb.ListApplicationAgentsRequest,
	opts ...grpc.CallOption,
) (configpb.ConfigManagementAPI_ListApplicationAgentsClient, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ListApplicationAgents(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) CreateApplicationAgent(
	ctx context.Context,
	request *configpb.CreateApplicationAgentRequest,
	opts ...grpc.CallOption) (*configpb.CreateApplicationAgentResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateApplicationAgent(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) UpdateApplicationAgent(
	ctx context.Context,
	request *configpb.UpdateApplicationAgentRequest,
	opts ...grpc.CallOption) (*configpb.UpdateApplicationAgentResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.UpdateApplicationAgent(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) DeleteApplicationAgent(
	ctx context.Context,
	request *configpb.DeleteApplicationAgentRequest,
	opts ...grpc.CallOption) (*configpb.DeleteApplicationAgentResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteApplicationAgent(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
