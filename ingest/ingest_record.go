// Copyright (c) 2023 IndyKite
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

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
)

// IngestRecord is deprecated and will be removed: use Batch functions.
func (c *Client) IngestRecord(
	ctx context.Context,
	record *ingestpb.Record,
	opts ...grpc.CallOption,
) (*ingestpb.IngestRecordResponse, error) {
	ctx = insertMetadata(ctx, c.xMetadata)
	resp, err := c.client.IngestRecord(
		ctx, &ingestpb.IngestRecordRequest{
			Record: record,
		}, opts...)
	return resp, err
}
