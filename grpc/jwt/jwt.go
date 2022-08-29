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
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

type (
	jwtAccessTokenSource struct {
		template jwt.Token
		signer   jwk.Key
		alg      jwa.SignatureAlgorithm
	}
)

func CreateTokenSourceFrom(privateKeyJWK interface{}, clientID string) (oauth2.TokenSource, error) {
	privateKeyJWKBytes, err := interfaceToBytes(privateKeyJWK)
	if err != nil {
		return nil, err
	}

	return JWTokenSource(privateKeyJWKBytes, false, clientID)
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
		alg:      jwa.SignatureAlgorithm(key.Algorithm()),
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
	h := jws.NewHeaders()
	// This is mandatory !!!
	_ = h.Set(jws.KeyIDKey, ts.signer.KeyID())

	signed, err := jwt.Sign(token, ts.alg, ts.signer, jwt.WithHeaders(h))
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{TokenType: "Bearer", Expiry: exp.Add(-time.Minute), AccessToken: string(signed)}, nil
}
