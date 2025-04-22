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
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	"github.com/indykite/indykite-sdk-go/grpc/config"
	"github.com/indykite/indykite-sdk-go/grpc/internal"
)

// A ClientOption is an option for an API client.
type ClientOption interface {
	Apply(*internal.DialSettings) //nolint:inamedparam // Not used.
}

// WithEndpoint returns a ClientOption that overrides the default endpoint.
func WithEndpoint(url string) ClientOption {
	return withEndpoint(url)
}

type withEndpoint string

// Apply returns endpoint.
func (w withEndpoint) Apply(o *internal.DialSettings) {
	o.Endpoint = string(w)
}

// WithTokenSource returns a ClientOption that specifies an OAuth2 token source.
// This value is ignored if used WithInsecure.
// Be sure NOT to have set INDYKITE env variables or use WithIgnoredEnvVariables.
func WithTokenSource(s oauth2.TokenSource) ClientOption {
	return withTokenSource{s}
}

type withTokenSource struct{ ts oauth2.TokenSource }

func (w withTokenSource) Apply(o *internal.DialSettings) {
	o.TokenSource = w.ts
}

// WithCredentialsJSON returns a ClientOption that authenticates
// API calls with the given service account or refresh token JSON
// credentials. This value is ignored if used WithInsecure.
// Be sure NOT to have set INDYKITE env variables or use WithIgnoredEnvVariables.
func WithCredentialsJSON(p []byte) ClientOption {
	return WithCredentialsLoader(config.StaticCredentialsJSON(p))
}

// WithGRPCConn returns a ClientOption that specifies the gRPC client to use as
// basis of communications. This option may only be used with services that support gRPC.
func WithGRPCConn(conn *grpc.ClientConn) ClientOption {
	return withGRPCConn{conn}
}

type withGRPCConn struct{ conn *grpc.ClientConn }

func (w withGRPCConn) Apply(o *internal.DialSettings) {
	o.GRPCConn = w.conn
}

// WithGRPCDialOption returns a ClientOption that appends a new grpc.DialOption
// to an underlying gRPC dial. It does not work with WithGRPCConn.
func WithGRPCDialOption(opt grpc.DialOption) ClientOption {
	return withGRPCDialOption{opt}
}

type withGRPCDialOption struct{ opt grpc.DialOption }

func (w withGRPCDialOption) Apply(o *internal.DialSettings) {
	o.GRPCDialOpts = append(o.GRPCDialOpts, w.opt)
}

// WithInsecure returns a ClientOption that sets insecure to true
// This is intended to be used only in tests!
func WithInsecure() ClientOption {
	return withInsecure(true)
}

type withInsecure bool

func (w withInsecure) Apply(o *internal.DialSettings) {
	o.Insecure = bool(w)
}

// WithUserAgent returns a ClientOption that sets the User-Agent.
func WithUserAgent(ua string) ClientOption {
	return withUserAgent(ua)
}

type withUserAgent string

func (w withUserAgent) Apply(o *internal.DialSettings) { o.UserAgent = string(w) }

// WithConnectionPoolSize Returns a ClientOption that sets up Connection Pool Size.
func WithConnectionPoolSize(poolSize int) ClientOption {
	return withConnectionPoolSize(poolSize)
}

type withConnectionPoolSize int

func (w withConnectionPoolSize) Apply(o *internal.DialSettings) {
	o.GRPCConnPoolSize = int(w)
}

// WithTelemetryDisabled returns a ClientOption that disables default telemetry (Open Telemetry)
// settings on gRPC and HTTP clients.
func WithTelemetryDisabled() ClientOption {
	return withTelemetryDisabled{}
}

type withTelemetryDisabled struct{}

func (withTelemetryDisabled) Apply(o *internal.DialSettings) {
	o.TelemetryDisabled = true
}

// WithTelemetryTraceName returns a ClientOption that sets default tracer name
// settings on gRPC and HTTP clients.
// If the name is an empty string then provider uses default name.
func WithTelemetryTraceName(name string) ClientOption {
	return withTelemetryName(name)
}

type withTelemetryName string

func (w withTelemetryName) Apply(o *internal.DialSettings) {
	o.TraceName = string(w)
}

// WithRetryOptions returns a ClientOption that sets default retry options.
func WithRetryOptions(opts ...retry.CallOption) ClientOption {
	return withRetryCallOption(opts)
}

type withRetryCallOption []retry.CallOption

func (w withRetryCallOption) Apply(o *internal.DialSettings) {
	o.RetryOpts = append(o.RetryOpts, w...)
}

// WithCredentialsLoader returns a ClientOption that adds ConfigLoaders.
func WithCredentialsLoader(opts ...config.CredentialsLoader) ClientOption {
	return withConfigLoaderOption(opts)
}

type withConfigLoaderOption []config.CredentialsLoader

func (w withConfigLoaderOption) Apply(o *internal.DialSettings) {
	o.CredentialsLoaders = append(o.CredentialsLoaders, w...)
}

// WithServiceAccount returns a ClientOption that requires service account configuration and fails if
// an application agent credential is loaded.
func WithServiceAccount() ClientOption {
	return withServiceAccountOption(true)
}

type withServiceAccountOption bool

func (w withServiceAccountOption) Apply(o *internal.DialSettings) {
	o.ServiceAccount = bool(w)
}
