// Copyright (c) 2026 IndyKite
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

// Package indykite is the root facade for the IndyKite REST SDK. It wires the
// shared transport and authentication once and exposes the per-service clients,
// so callers do not assemble auth/transport by hand.
//
// The platform has two planes with two credential types:
//
//   - Runtime / data plane (App Agent credential token, sent in
//     X-IK-ClientKey): AuthZEN, ContX IQ, capture and entity matching. Use
//     [NewClient] / [NewClientFromEnv] and the accessors on [Client].
//   - Control plane (Service Account credential JSON, Bearer token): config
//     management. Use [NewAdmin] / [NewAdminFromEnv].
//
// They are separate because they authenticate with different credentials and
// identities; a single credential cannot serve both.
//
//	cli, _ := indykite.NewClientFromEnv(ctx, indykite.WithRegion("eu"))
//	ok, _ := cli.AuthZEN().Allowed(ctx,
//	    authzen.NewNode("Person", "ada"), "PROVISION", authzen.NewNode("Server", "gpu-7"))
package indykite

import (
	"context"
	"net/http"

	"github.com/indykite/indykite-sdk-go/auth"
	"github.com/indykite/indykite-sdk-go/authzen"
	"github.com/indykite/indykite-sdk-go/capture"
	"github.com/indykite/indykite-sdk-go/ciq"
	"github.com/indykite/indykite-sdk-go/config"
	"github.com/indykite/indykite-sdk-go/entitymatching"
	"github.com/indykite/indykite-sdk-go/transport"
)

// Option configures the underlying transport (region/base URL, retries,
// tracing, HTTP client, user agent).
type Option = transport.Option

// RetryConfig is the transport retry policy.
type RetryConfig = transport.RetryConfig

// WithRegion resolves the base URL as https://<region>.api.indykite.com.
func WithRegion(region string) Option { return transport.WithRegion(region) }

// WithBaseURL sets an explicit base URL (highest precedence).
func WithBaseURL(rawURL string) Option { return transport.WithBaseURL(rawURL) }

// WithTracing enables OpenTelemetry HTTP instrumentation.
func WithTracing() Option { return transport.WithTracing() }

// WithRetry overrides the retry policy.
func WithRetry(r RetryConfig) Option { return transport.WithRetry(r) }

// WithHTTPClient supplies a custom *http.Client.
func WithHTTPClient(c *http.Client) Option { return transport.WithHTTPClient(c) }

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option { return transport.WithUserAgent(ua) }

// Client is the runtime-plane facade. It is built from an App Agent credential
// and exposes the runtime services over a single shared transport.
type Client struct {
	authzen        *authzen.Client
	ciq            *ciq.Client
	capture        *capture.Client
	entityMatching *entitymatching.Client
}

// NewClient builds the runtime facade from the App Agent credential token (the
// value sent verbatim in X-IK-ClientKey).
func NewClient(_ context.Context, appAgentToken string, opts ...Option) (*Client, error) {
	a, err := auth.AppAgentFromToken(appAgentToken)
	if err != nil {
		return nil, err
	}
	return buildClient(a, opts...)
}

// NewClientFromEnv builds the runtime facade from the App Agent credential
// token in the environment (INDYKITE_APPLICATION_CREDENTIALS[_FILE] holds the
// token itself).
func NewClientFromEnv(ctx context.Context, opts ...Option) (*Client, error) {
	a, err := auth.AppAgentFromEnv(ctx)
	if err != nil {
		return nil, err
	}
	return buildClient(a, opts...)
}

func buildClient(a *auth.Authenticator, opts ...Option) (*Client, error) {
	t, err := transport.NewClient(a, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		authzen:        authzen.NewClient(t),
		ciq:            ciq.NewClient(t),
		capture:        capture.NewClient(t),
		entityMatching: entitymatching.NewClient(t),
	}, nil
}

// AuthZEN returns the authorization (AuthZEN) client.
func (c *Client) AuthZEN() *authzen.Client { return c.authzen }

// CIQ returns the ContX IQ query client.
func (c *Client) CIQ() *ciq.Client { return c.ciq }

// Capture returns the capture (IKG ingest) client.
func (c *Client) Capture() *capture.Client { return c.capture }

// EntityMatching returns the entity-matching client.
func (c *Client) EntityMatching() *entitymatching.Client { return c.entityMatching }

// NewAdmin builds the control-plane facade from a Service Account credentials
// JSON. Config operations live here because they use a different credential and
// identity than the runtime plane.
func NewAdmin(ctx context.Context, credentialsJSON []byte, opts ...Option) (*config.AdminClient, error) {
	return config.NewAdminClientFromCredentials(ctx, credentialsJSON, opts...)
}

// NewAdminFromEnv builds the control-plane facade from the Service Account
// credentials in the environment (INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE]).
func NewAdminFromEnv(ctx context.Context, opts ...Option) (*config.AdminClient, error) {
	a, err := auth.ServiceAccountFromEnv(ctx)
	if err != nil {
		return nil, err
	}
	t, err := transport.NewClient(a, opts...)
	if err != nil {
		return nil, err
	}
	return config.NewAdminClient(t), nil
}
