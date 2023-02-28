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

package oauth2

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"

	"github.com/indykite/jarvis-sdk-go/errors"
	"github.com/indykite/jarvis-sdk-go/grpc/config"
	"github.com/indykite/jarvis-sdk-go/grpc/jwt"
)

var (
	credentials *config.CredentialsConfig
	tokenSource oauth2.TokenSource
)

// GetRefreshableTokenSource Token returns a token or an error.
// Token must be safe for concurrent use by multiple goroutines.
// The returned Token must not be modified.
func GetRefreshableTokenSource(
	ctx context.Context,
	credentialsLoaders []config.CredentialsLoader,
) (oauth2.TokenSource, error) {
	var err error
	for _, v := range credentialsLoaders {
		credentials, err = v(ctx)
		if err != nil {
			return nil, errors.NewInvalidArgumentErrorWithCause(err, "unable to load credentials")
		}
		if credentials != nil {
			break
		}
	}

	if credentials == nil {
		return nil, errors.NewInvalidArgumentError("credentials not found")
	}

	if tokenSource == nil {
		tokenSource, err = jwt.CreateTokenSource(credentials)
		if err != nil {
			return nil, err
		}
	}

	// first token
	token, err := tokenSource.Token()
	if err != nil {
		return nil, errors.New(codes.Internal, "unable to fetch token %v", err)
	}

	// refreshable token source, it can refresh every time we need a new one transparently.
	return oauth2.ReuseTokenSource(token, tokenSource), nil
}

// GetHTTPClient returns an authenticated HTTP client
// that always injects a valid token.
func GetHTTPClient(ctx context.Context, credentialsLoaders []config.CredentialsLoader) (*http.Client, error) {
	reusableTokenSource, err := GetRefreshableTokenSource(ctx, credentialsLoaders)
	if err != nil {
		return nil, errors.New(codes.Internal, "unable to fetch token %v", err)
	}

	// Create an http.Client that always injects the valid token
	return oauth2.NewClient(ctx, reusableTokenSource), nil
}
