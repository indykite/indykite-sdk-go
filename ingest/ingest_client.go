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
	"fmt"
	"math"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	ingestpb "github.com/indykite/indykite-sdk-go/gen/indykite/ingest/v1beta3"
	api "github.com/indykite/indykite-sdk-go/grpc"
)

// Client initializes Client.
type Client struct {
	// Closing handler of connection pool of gRPC connections to the service.
	closeHandler func() error

	// The gRPC API client.
	client ingestpb.IngestAPIClient
	// The gRPC StreamRecords client
	stream ingestpb.IngestAPI_StreamRecordsClient

	// The metadata to be sent with each request.
	xMetadata metadata.MD
}

// NewClient creates a new Ingest API gRPC Client.
func NewClient(ctx context.Context, opts ...api.ClientOption) (*Client, error) {
	clientOpts := defaultClientOptions()
	connPool, _, err := api.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	c := &Client{
		xMetadata:    metadata.Pairs("x-jarvis-client", fmt.Sprintf("client/%s grpc/%s", versionClient, grpc.Version)),
		client:       ingestpb.NewIngestAPIClient(connPool),
		closeHandler: connPool.Close,
	}
	return c, nil
}

// NewTestClient creates a new Ingest gRPC Client.
func NewTestClient(client ingestpb.IngestAPIClient) (*Client, error) {
	c := &Client{
		xMetadata: metadata.Pairs("x-jarvis-client",
			fmt.Sprintf("client/%s grpc/%s", versionClient, grpc.Version)),
		client:       client,
		closeHandler: func() error { return nil },
	}
	return c, nil
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *Client) Close() error {
	if c.stream != nil {
		_ = c.stream.CloseSend()
	}
	return c.closeHandler()
}

func insertMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	out, _ := metadata.FromOutgoingContext(ctx)
	out = out.Copy()
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return metadata.NewOutgoingContext(ctx, out)
}

func defaultClientOptions() []api.ClientOption {
	return []api.ClientOption{
		api.WithGRPCDialOption(grpc.WithDisableServiceConfig()),
		api.WithGRPCDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}
