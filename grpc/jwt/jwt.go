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
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/oauth2"

	"github.com/indykite/jarvis-sdk-go/grpc/config"
)

type (
	jwtAccessTokenSource struct {
		template jwt.Token
		signer   jwk.Key
	}
)

func CreateTokenSourceFromPrivateKey(privateKeyJWK interface{}, clientID string) (oauth2.TokenSource, error) {
	privateKeyJWKBytes, err := interfaceToBytes(privateKeyJWK)
	if err != nil {
		return nil, err
	}

	return JWTokenSource(privateKeyJWKBytes, false, clientID)
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

	var err error
	switch {
	case credentials.PrivateKeyJWK != nil:
		return JWTokenSource(credentials.PrivateKeyJWK, false, clientID)
	case credentials.PrivateKeyPKCS8Base64 != "":
		var raw []byte
		raw, err = base64.StdEncoding.DecodeString(credentials.PrivateKeyPKCS8Base64)
		if err != nil {
			return nil, err
		}
		return JWTokenSource(raw, true, clientID)
	case credentials.PrivateKeyPKCS8 != "":
		return JWTokenSource([]byte(credentials.PrivateKeyPKCS8), true, clientID)
	default:
		return nil, errors.New("unable to find secret credential")
	}
}

func interfaceToBytes(privateKeyJWK interface{}) ([]byte, error) {
	if stringValue, ok := privateKeyJWK.(string); ok {
		return []byte(stringValue), nil
	}

	privateKeyJWKBytes, err := json.Marshal(privateKeyJWK)
	if err != nil {
		return nil, err
	}

	return privateKeyJWKBytes, nil
}

func JWTokenSource(secretKey []byte, pem bool, clientID string) (oauth2.TokenSource, error) {
	var (
		key jwk.Key
		err error
	)
	if pem {
		key, err = jwk.ParseKey(secretKey, jwk.WithPEM(pem))
	} else {
		if secretKey[0] == '"' {
			var raw string
			err = json.Unmarshal(secretKey, &raw)
			if err != nil {
				return nil, err
			}
			key, err = jwk.ParseKey([]byte(raw))
		} else {
			key, err = jwk.ParseKey(secretKey)
		}
	}
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

	ts := &jwtAccessTokenSource{
		template: t,
		signer:   key,
	}
	return oauth2.ReuseTokenSource(nil, ts), nil
}

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) {
	iat := time.Now()
	exp := iat.Add(time.Hour)

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
	return &oauth2.Token{TokenType: "Bearer", Expiry: exp.Add(-time.Minute), AccessToken: string(signed)}, nil
}
