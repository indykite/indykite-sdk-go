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

package grpc_test

import (
	"context"
	"net"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/credentials/insecure"

	google_grpc "google.golang.org/grpc"

	"github.com/indykite/jarvis-sdk-go/grpc"
)

var _ = Describe("Test creating gRPC connection from setting", func() {

	var (
		server   *google_grpc.Server
		listener net.Listener
	)
	BeforeEach(func() {
		l, err := net.Listen("tcp", "localhost:0")
		Expect(err).To(Succeed())

		s := google_grpc.NewServer()
		// nolint
		go s.Serve(l)

		server, listener = s, l
	})
	AfterEach(func() {
		server.Stop()
	})

	It("Dial should return passed connection when is specified", func() {
		originConn, err := google_grpc.Dial(listener.Addr().String(),
			google_grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).To(Succeed())

		conn, _, err := grpc.Dial(context.Background(), grpc.WithGRPCConn(originConn))
		Expect(err).To(Succeed())
		Expect(conn).To(BeIdenticalTo(originConn))
	})

	It("DialPool should return single connection pool with given connection", func() {
		originConn, err := google_grpc.Dial(listener.Addr().String(),
			google_grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).To(Succeed())

		conn, _, err := grpc.DialPool(
			context.Background(),
			grpc.WithGRPCConn(originConn),
			grpc.WithConnectionPoolSize(10), // Is ignored because is passed gRPC Connection
		)
		Expect(err).To(Succeed())
		Expect(conn.Num()).To(Equal(1))
		Expect(conn.Conn()).To(BeIdenticalTo(originConn))
	})

	It("Check DialPool returns pool of connections", func() {
		poolSize := 5

		conn, _, err := grpc.DialPool(
			context.Background(),
			grpc.WithEndpoint(listener.Addr().String()),
			grpc.WithConnectionPoolSize(poolSize),
			grpc.WithInsecure(),
		)
		Expect(err).To(Succeed())
		Expect(conn.Num()).To(Equal(poolSize))

	})
})
