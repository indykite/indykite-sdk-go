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

//nolint:testpackage // We are not able to name that grpc_test for now. Needs bigger refactor.
package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Connection pools", func() {
	It("Get connections in rounds", func() {
		conn1 := &grpc.ClientConn{}
		conn2 := &grpc.ClientConn{}
		conn3 := &grpc.ClientConn{}

		pool := &roundRobinConnPool{
			conns: []*grpc.ClientConn{
				conn1, conn2, conn3,
			},
		}

		Expect(pool.Num()).To(Equal(3))
		// Increment is done first, so Coon() starting from second connection further
		Expect(pool.Conn()).To(BeIdenticalTo(conn2))
		Expect(pool.Conn()).To(BeIdenticalTo(conn3))
		Expect(pool.Conn()).To(BeIdenticalTo(conn1))
		Expect(pool.Conn()).To(BeIdenticalTo(conn2))
	})

	It("Check Close() closes all connections", func() {
		_, listener := mockServer()

		pool := &roundRobinConnPool{}
		for i := 0; i < 4; i++ {
			conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).To(Succeed())
			pool.conns = append(pool.conns, conn)
		}

		if err := pool.Close(); err != nil {
			Fail(fmt.Sprintf("pool.Close: %v", err))
		}
		for _, conn := range pool.conns {
			Expect(conn.GetState()).To(Equal(connectivity.Shutdown))
		}
	})
})

func mockServer() (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", "localhost:0")
	Expect(err).To(Succeed())

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(GinkgoUnaryInterceptor),
		grpc.StreamInterceptor(GinkgoStreamInterceptor),
	)
	go func() { _ = gRPCServer.Serve(listener) }()

	return gRPCServer, listener
}

func GinkgoUnaryInterceptor(ctx context.Context,
	req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer GinkgoRecover()
	return handler(ctx, req)
}

func GinkgoStreamInterceptor(
	srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	defer GinkgoRecover()
	return handler(srv, ss)
}
