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

func (c *Client) CreateConfigNode(ctx context.Context, request *NodeRequest, opts ...grpc.CallOption) (
	*configpb.CreateConfigNodeResponse, error) {
	if request == nil || request.create == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil or not create request")
	}

	if err := request.validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.CreateConfigNode(ctx, request.create, opts...)
	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ReadConfigNode(ctx context.Context, req *NodeRequest, opts ...grpc.CallOption) (
	*configpb.ReadConfigNodeResponse, error) {
	if req == nil || req.read == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil or not read request")
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadConfigNode(ctx, req.read, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) UpdateConfigNode(ctx context.Context, req *NodeRequest, opts ...grpc.CallOption) (
	*configpb.UpdateConfigNodeResponse, error) {
	if req == nil || req.update == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil or not update request")
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.UpdateConfigNode(ctx, req.update, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) DeleteConfigNode(ctx context.Context, req *NodeRequest, opts ...grpc.CallOption) (
	*configpb.DeleteConfigNodeResponse, error) {
	if req == nil || req.delete == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil or not delete request")
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.DeleteConfigNode(ctx, req.delete, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ListConfigNodeVersions(ctx context.Context, req *NodeRequest, opts ...grpc.CallOption) (
	*configpb.ListConfigNodeVersionsResponse, error) {
	if req == nil || req.listVersions == nil {
		return nil, errors.NewInvalidArgumentError("invalid nil or not read request")
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ListConfigNodeVersions(ctx, req.listVersions, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
