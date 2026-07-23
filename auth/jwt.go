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

package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"golang.org/x/oauth2"
)

// tokenSourceFromCredentials mints a self-signed JWT token source from the
// credential's private key. The resulting token string is what is placed in the
// X-IK-ClientKey or Authorization header.
func tokenSourceFromCredentials(c *Credentials) (oauth2.TokenSource, error) {
	// A pre-issued token is preferred: it is served verbatim until its own
	// expiry. When a private key is also present (Service Account credentials
	// carry both) a minting source is used as the refresh fallback, so a
	// long-running client mints a fresh JWT once the pre-issued token lapses
	// instead of sending a stale token. Without a key there is nothing to fall
	// back to, so the token is served as-is (the platform rejects it if expired).
	if c.Token != "" {
		preTok := &oauth2.Token{AccessToken: c.Token, Expiry: preIssuedExpiry(c.Token)}
		if c.hasPrivateKey() {
			if mint, err := mintingSource(c); err == nil {
				return oauth2.ReuseTokenSource(preTok, mint), nil
			}
		}
		return oauth2.StaticTokenSource(preTok), nil
	}
	return mintingSource(c)
}

// mintingSource builds a token source that self-signs (and refreshes) JWTs from
// the credential's private key.
func mintingSource(c *Credentials) (oauth2.TokenSource, error) {
	clientID, err := c.ClientID()
	if err != nil {
		return nil, err
	}

	var lifetime time.Duration
	if c.TokenLifetime != "" {
		if lifetime, err = time.ParseDuration(c.TokenLifetime); err != nil {
			return nil, fmt.Errorf("indykite: unable to parse tokenLifetime: %w", err)
		}
	}

	switch {
	case c.PrivateKeyJWK != nil:
		return jwtTokenSource(c.PrivateKeyJWK, false, clientID, lifetime)
	case c.PrivateKeyPKCS8Base64 != "":
		raw, decErr := base64.StdEncoding.DecodeString(c.PrivateKeyPKCS8Base64)
		if decErr != nil {
			return nil, decErr
		}
		return jwtTokenSource(raw, true, clientID, lifetime)
	case c.PrivateKeyPKCS8 != "":
		return jwtTokenSource([]byte(c.PrivateKeyPKCS8), true, clientID, lifetime)
	default:
		return nil, errors.New("indykite: no private key in credentials")
	}
}

type jwtAccessTokenSource struct {
	template      jwt.Token
	signer        jwk.Key
	tokenLifetime time.Duration
}

func jwtTokenSource(secretKey []byte, pem bool, clientID string, lifetime time.Duration) (oauth2.TokenSource, error) {
	key, err := parseKey(secretKey, pem)
	if err != nil {
		return nil, err
	}
	// Regenerate the kid the same way the backend does.
	_ = key.Remove(jwk.KeyIDKey)
	if err = jwk.AssignKeyID(key); err != nil {
		return nil, errors.New("indykite: failed to assign kid for key")
	}
	// PEM-parsed (PKCS#8) keys carry no alg; infer it so signing works.
	if _, ok := key.Algorithm(); !ok {
		alg, algErr := inferAlgorithm(key)
		if algErr != nil {
			return nil, algErr
		}
		_ = key.Set(jwk.AlgorithmKey, alg)
	}

	t := jwt.New()
	_ = t.Set(jwt.IssuerKey, clientID)
	_ = t.Set(jwt.SubjectKey, clientID)

	// A token is issued roughly every minute, so keep the lifetime above 2 min.
	if lifetime < 2*time.Minute || lifetime > 24*time.Hour {
		lifetime = time.Hour
	}

	ts := &jwtAccessTokenSource{template: t, signer: key, tokenLifetime: lifetime}
	return oauth2.ReuseTokenSource(nil, ts), nil
}

// Token implements oauth2.TokenSource.
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

	alg, _ := ts.signer.Algorithm()
	signed, err := jwt.Sign(token, jwt.WithKey(alg, ts.signer))
	if err != nil {
		return nil, err
	}
	// Expire a minute early so a refresh happens before the server rejects it.
	return &oauth2.Token{TokenType: "Bearer", Expiry: exp.Add(-time.Minute), AccessToken: string(signed)}, nil
}

// preIssuedExpiry returns the effective expiry of a pre-issued token: a minute
// before the JWT's exp claim (so it is refreshed before the platform would
// reject it), or the zero time when there is no discoverable expiry — an opaque
// (non-JWT) token or a JWT without exp — which oauth2 treats as never-expiring.
// The signature is NOT verified; the platform does that.
func preIssuedExpiry(token string) time.Time {
	parsed, err := jwt.ParseString(token, jwt.WithVerify(false), jwt.WithValidate(false))
	if err != nil {
		return time.Time{}
	}
	exp, ok := parsed.Expiration()
	if !ok {
		return time.Time{}
	}
	return exp.Add(-time.Minute)
}

// inferAlgorithm picks the signature algorithm implied by the key type, for
// keys whose serialization (PKCS#8/PEM) does not carry one.
func inferAlgorithm(key jwk.Key) (jwa.SignatureAlgorithm, error) {
	switch key.KeyType() {
	case jwa.EC():
		var crv jwa.EllipticCurveAlgorithm
		if err := key.Get(jwk.ECDSACrvKey, &crv); err != nil {
			return jwa.SignatureAlgorithm{}, fmt.Errorf("indykite: EC key has no curve: %w", err)
		}
		switch crv {
		case jwa.P256():
			return jwa.ES256(), nil
		case jwa.P384():
			return jwa.ES384(), nil
		case jwa.P521():
			return jwa.ES512(), nil
		default:
			return jwa.SignatureAlgorithm{}, fmt.Errorf("indykite: unsupported EC curve %q", crv)
		}
	case jwa.RSA():
		return jwa.RS256(), nil
	case jwa.OKP():
		return jwa.EdDSA(), nil
	default:
		return jwa.SignatureAlgorithm{}, fmt.Errorf("indykite: unsupported key type %q", key.KeyType())
	}
}

func parseKey(secretKey []byte, pem bool) (jwk.Key, error) {
	if pem {
		return jwk.ParseKey(secretKey, jwk.WithPEM(pem))
	}
	if len(secretKey) > 0 && secretKey[0] != '"' {
		return jwk.ParseKey(secretKey)
	}
	var raw string
	if err := json.Unmarshal(secretKey, &raw); err != nil {
		return nil, err
	}
	return jwk.ParseKey([]byte(raw))
}
