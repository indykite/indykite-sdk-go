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
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"

	identitypb "github.com/indykite/jarvis-sdk-go/gen/indykite/identity/v1beta2"
)

type (
	tokenSource struct {
		ctx              context.Context
		client           identitypb.IdentityManagementAPIClient
		applicationToken oauth2.TokenSource
	}
)

var (
	defaultGrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	_                = defaultGrantType
)

func IndyKiteTokenSource(
	ctx context.Context,
	source identitypb.IdentityManagementAPIClient,
) (oauth2.TokenSource, error) {
	return &tokenSource{
		ctx:    ctx,
		client: source,
	}, nil
}

// Token returns a token or an error.
// Token must be safe for concurrent use by multiple goroutines.
// The returned Token must not be modified.
func (t *tokenSource) Token() (*oauth2.Token, error) {
	resp, err := t.client.GetAccessToken(t.ctx, &identitypb.GetAccessTokenRequest{
		// Id: "urn:indy:token",
	}, grpc.PerRPCCredentials(oauth.TokenSource{TokenSource: t.applicationToken}))
	if err != nil {
		return nil, err
	}
	if resp.Token != nil {
		return &oauth2.Token{
			AccessToken: resp.Token.AccessToken,
			TokenType:   resp.Token.TokenType,
			Expiry:      time.Now().Add(time.Duration(resp.Token.ExpiresIn-60) * time.Second),
		}, nil
	}
	return nil, fmt.Errorf("indykite: token is not found")
}
