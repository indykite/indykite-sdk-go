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

package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/grpc/config"
	"github.com/indykite/indykite-sdk-go/grpc/internal"
)

// Dial returns a gRPC connection configured with the given ClientOptions.
// In this call WithConnectionPoolSize() is ignored! Use DialPool instead.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, *config.CredentialsConfig, error) {
	o := new(internal.DialSettings)
	for _, opt := range opts {
		opt.Apply(o)
	}
	if o.GRPCConn != nil {
		return o.GRPCConn, nil, nil
	}
	var (
		dop []grpc.DialOption
		cfg *config.CredentialsConfig
		err error
	)
	if dop, cfg, err = o.Build(ctx); err != nil {
		return nil, nil, err
	}
	c, err := grpc.NewClient(o.Endpoint, dop...)
	if err != nil {
		return nil, nil, err
	}
	return c, cfg, nil
}

// DialPool returns a connection pool configured with ClientOptions.
// If grpcConn is specified in ClientOptions, pool size is reset to 1.
func DialPool(ctx context.Context, opts ...ClientOption) (ConnPool, *config.CredentialsConfig, error) {
	o := new(internal.DialSettings)
	for _, opt := range opts {
		opt.Apply(o)
	}
	if o.GRPCConn != nil {
		return &singleConnPool{o.GRPCConn}, nil, nil
	}
	var (
		dop  []grpc.DialOption
		cfg  *config.CredentialsConfig
		err  error
		conn *grpc.ClientConn
	)
	if dop, cfg, err = o.Build(ctx); err != nil {
		return nil, nil, err
	}
	poolSize := o.GRPCConnPoolSize
	if o.GRPCConn != nil {
		// WithGRPCConn is technically incompatible with WithGRPCConnectionPool.
		// Always assume pool size is 1 when a grpc.ClientConn is explicitly used.
		poolSize = 1
	}

	if poolSize < 2 {
		// Fast path for common case for a connection pool with a single connection.
		conn, err = grpc.NewClient(o.Endpoint, dop...)
		if err != nil {
			return nil, nil, err
		}
		return &singleConnPool{conn}, cfg, nil
	}

	pool := &roundRobinConnPool{}
	for i := 0; i < poolSize; i++ {
		conn, err = grpc.NewClient(o.Endpoint, dop...)
		if err != nil {
			defer func() { _ = pool.Close() }() //nolint:revive,gocritic // If this happen, loop is exited
			return nil, nil, err
		}
		pool.conns = append(pool.conns, conn)
	}
	return pool, cfg, nil
}
