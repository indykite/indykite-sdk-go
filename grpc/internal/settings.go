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

// Package internal contains Dial settings and helpers
package internal

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/indykite/indykite-sdk-go/grpc/config"
	"github.com/indykite/indykite-sdk-go/grpc/jwt"
)

const (
	tcpUserTimeout = 20 * time.Second
)

// DialSettings holds information needed to establish a connection with a service.
//
//nolint:govet // TODO Align fields
type DialSettings struct {
	GRPCConn           *grpc.ClientConn
	RetryOpts          []retry.CallOption
	GRPCDialOpts       []grpc.DialOption
	TLSConfig          *tls.Config
	TokenSource        oauth2.TokenSource
	CredentialsLoaders []config.CredentialsLoader
	Endpoint           string
	//
	UserAgent        string
	TraceName        string
	GRPCConnPoolSize int
	Insecure         bool

	TelemetryDisabled bool
	// True if the credential is for configuration service.
	ServiceAccount bool
	credentials    *config.CredentialsConfig
}

// Build validates the settings and builds the client configurations.
func (ds *DialSettings) Build(ctx context.Context) ([]grpc.DialOption, *config.CredentialsConfig, error) {
	var dialOptions []grpc.DialOption
	for _, v := range ds.CredentialsLoaders {
		var err error
		ds.credentials, err = v(ctx)
		if err != nil {
			return nil, nil, err
		}
		if ds.credentials != nil {
			break
		}
	}
	//nolint:nestif // Needs bigger refactor
	if ds.credentials != nil {
		if ds.ServiceAccount && ds.credentials.ServiceAccountID == "" {
			return nil, nil, errors.New("empty serviceAccountId")
		}
		if !ds.ServiceAccount && ds.credentials.AppAgentID == "" {
			return nil, nil, errors.New("empty appAgentId")
		}

		var err error
		if ds.TokenSource == nil {
			ds.TokenSource, err = jwt.CreateTokenSource(ds.credentials)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if ds.TokenSource != nil {
		token, err := ds.TokenSource.Token()
		if err != nil {
			return nil, nil, err
		}
		if token.Type() != "Bearer" {
			return nil, nil, fmt.Errorf("unsupported token type, must be 'Bearer' but got %s", token.Type())
		}
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(&oauth.TokenSource{TokenSource: ds.TokenSource}))
	}

	if ds.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		clientTLSConfig, err := config.ClientTLSConfig()
		if err != nil {
			return nil, nil, err
		}
		// TODO fix me
		if false {
			clientTLSConfig.GetCertificate = func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
				// TODO fix #6uk4ct
				return nil, nil
			}
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(clientTLSConfig)))
	}

	var endpoint string
	if ds.credentials != nil {
		endpoint = ds.credentials.Endpoint
	}
	if ds.Endpoint != "" {
		endpoint = ds.Endpoint
	}
	if endpoint == "" {
		return nil, nil, errors.New("missing endpoint")
	}

	if !strings.HasPrefix(endpoint, "dns:///") {
		ds.Endpoint = "dns:///" + endpoint
	}

	dialOptions = append(dialOptions,
		// https://github.com/grpc/grpc/blob/master/doc/service_config.md
		// grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"grpclb":{"childPolicy":[{"pick_first":{}}]}}]}`),
		grpc.WithDisableServiceConfig(),
	)
	if len(ds.UserAgent) > 0 {
		dialOptions = append(dialOptions, grpc.WithUserAgent(ds.UserAgent))
	}

	dialOptions = addInterceptors(dialOptions, ds)
	dialOptions = append(dialOptions, ds.GRPCDialOpts...)

	dialOptions = append(dialOptions,
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{Timeout: tcpUserTimeout}
			return d.DialContext(ctx, "tcp", addr)
		}))
	return dialOptions, ds.credentials, nil
}

func addInterceptors(opts []grpc.DialOption, settings *DialSettings) []grpc.DialOption {
	var (
		retOpts []retry.CallOption
		si      []grpc.StreamClientInterceptor
		ui      []grpc.UnaryClientInterceptor
	)

	if len(settings.RetryOpts) > 0 {
		retOpts = settings.RetryOpts
		si = append(si, retry.StreamClientInterceptor(settings.RetryOpts...))
		ui = append(ui, retry.UnaryClientInterceptor(settings.RetryOpts...))
	} else {
		// disable by default - each call can enable if needed
		retOpts = []retry.CallOption{
			retry.WithBackoff(retry.BackoffLinear(100 * time.Millisecond)),
			// retry.WithBackoff(retry.BackoffExponential(100 * time.Millisecond)),
			retry.WithCodes(codes.ResourceExhausted, codes.Unavailable),
			retry.WithMax(12),
		}
		si = append(si, retry.StreamClientInterceptor())
		ui = append(ui, retry.UnaryClientInterceptor())
	}

	if settings.TelemetryDisabled {
		return append(opts,
			grpc.WithChainStreamInterceptor(retry.StreamClientInterceptor(retOpts...)),
			grpc.WithChainUnaryInterceptor(retry.UnaryClientInterceptor(retOpts...)))
	}

	return append(opts,
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainStreamInterceptor(si...),
		grpc.WithChainUnaryInterceptor(ui...))
}
