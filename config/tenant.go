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

	"github.com/indykite/jarvis-sdk-go/errors"
	configpb "github.com/indykite/jarvis-sdk-go/gen/indykite/config/v1beta1"
)

func (c *Client) CreateTenant(
	ctx context.Context,
	request *configpb.CreateTenantRequest,
	opts ...grpc.CallOption) (*configpb.CreateTenantResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateTenant(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ReadTenant(
	ctx context.Context,
	request *configpb.ReadTenantRequest,
	opts ...grpc.CallOption) (*configpb.ReadTenantResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadTenant(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) UpdateTenant(
	ctx context.Context,
	request *configpb.UpdateTenantRequest,
	opts ...grpc.CallOption) (*configpb.UpdateTenantResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.UpdateTenant(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ListTenants(
	ctx context.Context,
	request *configpb.ListTenantsRequest,
	opts ...grpc.CallOption,
) (configpb.ConfigManagementAPI_ListTenantsClient, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ListTenants(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) DeleteTenant(
	ctx context.Context,
	request *configpb.DeleteTenantRequest,
	opts ...grpc.CallOption) (*configpb.DeleteTenantResponse, error) {
	if request == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil request")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteTenant(ctx, request, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
