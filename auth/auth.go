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

// Package auth turns IndyKite credentials into the right HTTP authentication
// header for the REST APIs.
//
// The platform exposes two planes with two different credential artifacts:
//
//   - Runtime / data plane (AuthZEN, ContX IQ, capture, entity-matching): the
//     credential is the App Agent credential token itself — an opaque string
//     sent verbatim in the "X-IK-ClientKey" header.
//     INDYKITE_APPLICATION_CREDENTIALS[_FILE] holds exactly that token.
//   - Control plane (config management, /configs/v1/*): the credential is a
//     JSON artifact (serviceAccountId, baseUrl, a pre-issued "token", and the
//     private key material); the token goes in "Authorization: Bearer <...>".
//     A valid pre-issued token is sent as-is; when it has expired, a fresh
//     self-signed JWT is minted from the private key instead (cached and
//     refreshed via oauth2.ReuseTokenSource).
package auth

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// Header names used by the two planes.
const (
	HeaderClientKey     = "X-IK-ClientKey"
	HeaderAuthorization = "Authorization"
)

// Plane selects which authentication carrier (and therefore which header) is
// used for outgoing requests.
type Plane int

const (
	// PlaneRuntime is the data/runtime plane (App Agent token, X-IK-ClientKey).
	PlaneRuntime Plane = iota
	// PlaneControl is the control plane (Service Account token, Bearer).
	PlaneControl
)

// String renders the plane name.
func (p Plane) String() string {
	if p == PlaneControl {
		return "control"
	}
	return "runtime"
}

// TokenProvider returns a bearer token string, refreshing it as needed.
// It is the extension point for non-JWT auth (e.g. a token-exchange flow):
// implement this interface and pass it via NewWithProvider.
type TokenProvider interface {
	// Token returns a bearer token, refreshing it when needed.
	Token(ctx context.Context) (string, error)
}

// Authenticator applies the correct authentication header to HTTP requests for
// a given plane.
type Authenticator struct {
	provider TokenProvider

	// BaseURL is the endpoint hint carried by the credential, if any. The
	// transport may use it to default the base URL when none is configured.
	BaseURL string

	plane Plane
}

// New builds an Authenticator for the given plane from the first loader that
// yields a non-nil credential.
func New(ctx context.Context, plane Plane, loaders ...Loader) (*Authenticator, error) {
	creds, err := resolve(ctx, loaders...)
	if err != nil {
		return nil, err
	}

	src, err := tokenSourceFromCredentials(creds)
	if err != nil {
		return nil, err
	}

	return &Authenticator{
		provider: &tokenSourceProvider{ts: oauth2.ReuseTokenSource(nil, src)},
		plane:    plane,
		BaseURL:  firstNonEmpty(creds.BaseURL, creds.Endpoint),
	}, nil
}

// NewWithProvider builds an Authenticator from a custom TokenProvider, for auth
// flows other than the built-in self-signed JWT.
func NewWithProvider(plane Plane, provider TokenProvider) *Authenticator {
	return &Authenticator{provider: provider, plane: plane}
}

// AppAgentFromEnv builds a runtime-plane Authenticator from the App Agent
// credential token in INDYKITE_APPLICATION_CREDENTIALS[_FILE]. The value is
// the token itself, not JSON.
func AppAgentFromEnv(_ context.Context) (*Authenticator, error) {
	data, err := lookupEnvCredential(EnvAppAgentCredentials)
	if err != nil {
		return nil, err
	}
	token := string(bytes.TrimSpace(data))
	if token == "" {
		return nil, ErrNoCredentials
	}
	return AppAgentFromToken(token)
}

// AppAgentFromToken builds a runtime-plane Authenticator from the App Agent
// credential token, which is sent verbatim in X-IK-ClientKey.
func AppAgentFromToken(token string) (*Authenticator, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrNoCredentials
	}
	if strings.HasPrefix(token, "{") {
		return nil, errors.New(
			"indykite: the App Agent credential is the raw credential token, not a JSON document")
	}
	return &Authenticator{
		provider: staticTokenProvider(token),
		plane:    PlaneRuntime,
	}, nil
}

// ServiceAccountFromEnv builds a control-plane Authenticator from the environment.
func ServiceAccountFromEnv(ctx context.Context) (*Authenticator, error) {
	return New(ctx, PlaneControl, ServiceAccountEnvLoader)
}

// ServiceAccountFromJSON builds a control-plane Authenticator from a credentials JSON.
func ServiceAccountFromJSON(ctx context.Context, credentialsJSON []byte) (*Authenticator, error) {
	return New(ctx, PlaneControl, StaticJSON(credentialsJSON))
}

// Apply sets the plane-appropriate authentication header on the request.
func (a *Authenticator) Apply(ctx context.Context, req *http.Request) error {
	tok, err := a.provider.Token(ctx)
	if err != nil {
		return err
	}
	switch a.plane {
	case PlaneControl:
		req.Header.Set(HeaderAuthorization, "Bearer "+tok)
	default: // PlaneRuntime
		// X-IK-ClientKey is the platform's documented header name; HTTP header
		// keys are case-insensitive and Set canonicalizes on the wire.
		req.Header.Set(HeaderClientKey, tok) //nolint:canonicalheader // documented platform header, case-insensitive
	}
	return nil
}

// Plane returns the plane this Authenticator serves.
func (a *Authenticator) Plane() Plane { return a.plane }

// staticTokenProvider serves a pre-issued token verbatim.
type staticTokenProvider string

func (t staticTokenProvider) Token(_ context.Context) (string, error) { return string(t), nil }

type tokenSourceProvider struct {
	ts oauth2.TokenSource
}

func (p *tokenSourceProvider) Token(_ context.Context) (string, error) {
	t, err := p.ts.Token()
	if err != nil {
		return "", err
	}
	return t.AccessToken, nil
}

func resolve(ctx context.Context, loaders ...Loader) (*Credentials, error) {
	for _, l := range loaders {
		if l == nil {
			continue
		}
		creds, err := l(ctx)
		if err != nil {
			return nil, err
		}
		if creds != nil {
			return creds, nil
		}
	}
	return nil, ErrNoCredentials
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
