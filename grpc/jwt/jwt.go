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

package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/oauth2"

	"github.com/indykite/indykite-sdk-go/grpc/config"
)

type (
	jwtAccessTokenSource struct {
		template      jwt.Token
		signer        jwk.Key
		tokenLifetime time.Duration
	}
)

func CreateTokenSourceFromPrivateKey(privateKeyJWK any, clientID string) (oauth2.TokenSource, error) {
	privateKeyJWKBytes, err := interfaceToBytes(privateKeyJWK)
	if err != nil {
		return nil, err
	}

	return JWTokenSource(privateKeyJWKBytes, false, clientID, 0)
}

func CreateTokenSource(credentials *config.CredentialsConfig) (oauth2.TokenSource, error) {
	var clientID string
	switch {
	case credentials.AppAgentID != "":
		clientID = credentials.AppAgentID
	case credentials.ServiceAccountID != "":
		clientID = credentials.ServiceAccountID
	default:
		return nil, errors.New("missing client ID, AppAgentID or ServiceAccountID must be specified")
	}

	var (
		tokenLifetime time.Duration
		err           error
	)
	if credentials.TokenLifetime != "" {
		tokenLifetime, err = time.ParseDuration(credentials.TokenLifetime)
		if err != nil {
			return nil, fmt.Errorf("unable to parse 'tokenLifetime': %w", err)
		}
	}
	switch {
	case credentials.PrivateKeyJWK != nil:
		return JWTokenSource(credentials.PrivateKeyJWK, false, clientID, tokenLifetime)
	case credentials.PrivateKeyPKCS8Base64 != "":
		var raw []byte
		raw, err = base64.StdEncoding.DecodeString(credentials.PrivateKeyPKCS8Base64)
		if err != nil {
			return nil, err
		}
		return JWTokenSource(raw, true, clientID, tokenLifetime)
	case credentials.PrivateKeyPKCS8 != "":
		return JWTokenSource([]byte(credentials.PrivateKeyPKCS8), true, clientID, tokenLifetime)
	default:
		return nil, errors.New("unable to find secret credential")
	}
}

func interfaceToBytes(privateKeyJWK any) ([]byte, error) {
	if stringValue, ok := privateKeyJWK.(string); ok {
		return []byte(stringValue), nil
	}

	privateKeyJWKBytes, err := json.Marshal(privateKeyJWK)
	if err != nil {
		return nil, err
	}

	return privateKeyJWKBytes, nil
}

func parseKey(secretKey []byte, pem bool) (jwk.Key, error) {
	if pem {
		return jwk.ParseKey(secretKey, jwk.WithPEM(pem))
	}

	if secretKey[0] != '"' {
		return jwk.ParseKey(secretKey)
	}

	var raw string
	err := json.Unmarshal(secretKey, &raw)
	if err != nil {
		return nil, err
	}
	return jwk.ParseKey([]byte(raw))
}

func JWTokenSource(secretKey []byte, pem bool, clientID string,
	tokenLifetime time.Duration) (oauth2.TokenSource, error) {
	key, err := parseKey(secretKey, pem)
	if err != nil {
		return nil, err
	}
	// Remove user defined kid and generate new one same way as we do on BE
	// This is micro-optimization on client-side.
	_ = key.Remove(jwk.KeyIDKey)
	err = jwk.AssignKeyID(key)
	if err != nil {
		return nil, errors.New("failed to assign kid for public key")
	}

	t := jwt.New()
	_ = t.Set(jwt.IssuerKey, clientID)
	_ = t.Set(jwt.SubjectKey, clientID)

	// Don't let it be smaller than 2 min because token will be issued in every 1 minute by default.
	if tokenLifetime < 2*time.Minute || tokenLifetime > 24*time.Hour {
		tokenLifetime = time.Hour
	}

	ts := &jwtAccessTokenSource{
		template:      t,
		signer:        key,
		tokenLifetime: tokenLifetime,
	}
	return oauth2.ReuseTokenSource(nil, ts), nil
}

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) {
	iat := time.Now()
	exp := iat.Add(ts.tokenLifetime)

	token, err := ts.template.Clone()
	if err != nil {
		return nil, err
	}

	_ = token.Set(jwt.IssuedAtKey, iat)
	_ = token.Set(jwt.ExpirationKey, exp)
	_ = token.Set(jwt.JwtIDKey, uuid.New().String())

	signed, err := jwt.Sign(token, jwt.WithKey(ts.signer.Algorithm(), ts.signer))
	if err != nil {
		return nil, err
	}
	// Exp shell not be less than 2 min!
	return &oauth2.Token{TokenType: "Bearer", Expiry: exp.Add(-time.Minute), AccessToken: string(signed)}, nil
}
