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

package ingest

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/errors"
	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
)

// BatchDeleteNodeProperties returns resp.
func (c *Client) BatchDeleteNodeProperties(
	ctx context.Context,
	nodeProperties []*ingestpb.DeleteData_NodePropertyMatch,
	opts ...grpc.CallOption,
) (*ingestpb.BatchDeleteNodePropertiesResponse, error) {
	req := &ingestpb.BatchDeleteNodePropertiesRequest{
		NodeProperties: nodeProperties,
	}
	if err := req.Validate(); err != nil {
		return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to call BatchDeleteNodeProperties")
	}
	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.BatchDeleteNodeProperties(ctx, req, opts...)
	if s := errors.FromError(err); s != nil {
		return nil, s
	}
	return resp, nil
}
