// Copyright (c) 2024 IndyKite
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

package entitymatching

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/errors"
	entitymatchingpb "github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1"
)

func (c *Client) ReadEntityMatchingReport(
	ctx context.Context,
	req *entitymatchingpb.ReadEntityMatchingReportRequest,
	opts ...grpc.CallOption,
) (*entitymatchingpb.ReadEntityMatchingReportResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call ReadEntityMatchingReport")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadEntityMatchingReport(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) ReadSuggestedPropertyMapping(
	ctx context.Context,
	req *entitymatchingpb.ReadSuggestedPropertyMappingRequest,
	opts ...grpc.CallOption,
) (*entitymatchingpb.ReadSuggestedPropertyMappingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call ReadSuggestedPropertyMapping")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.ReadSuggestedPropertyMapping(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}

func (c *Client) RunEntityMatchingPipeline(
	ctx context.Context,
	req *entitymatchingpb.RunEntityMatchingPipelineRequest,
	opts ...grpc.CallOption,
) (*entitymatchingpb.RunEntityMatchingPipelineResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call RunEntityMatchingPipeline")
	}

	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.RunEntityMatchingPipeline(ctx, req, opts...)

	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
