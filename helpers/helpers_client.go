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

package helpers

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/entitymatching"
	"github.com/indykite/indykite-sdk-go/ingest"
	"github.com/indykite/indykite-sdk-go/knowledge"
)

type Client struct {
	// Closing handler of connection pool of gRPC connections to the service.
	closeHandler func() error

	// The gRPC API client.
	ClientKnowledge *knowledge.Client
	// The second gRPC API client.
	ClientIngest *ingest.Client

	ClientConfig         *config.Client
	ClientEntitymatching *entitymatching.Client

	// The metadata to be sent with each request.
	xMetadata metadata.MD
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *Client) Close() error {
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
