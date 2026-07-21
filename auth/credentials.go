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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Environment variables the default loaders read.
//
//nolint:gosec // G101: these are environment-variable names, not credential values.
const (
	EnvAppAgentCredentials       = "INDYKITE_APPLICATION_CREDENTIALS"
	EnvServiceAccountCredentials = "INDYKITE_SERVICE_ACCOUNT_CREDENTIALS"
)

// ErrNoCredentials is returned when no loader yields a credential.
var ErrNoCredentials = errors.New("indykite: no credentials found")

// Credentials is the Service Account credential artifact: the JSON document
// carrying baseUrl/endpoint, serviceAccountId, a pre-issued "token", and the
// private key material. When the pre-issued Token is still valid it is sent
// as-is; a fresh JWT is minted from the private key otherwise.
//
// (The runtime plane's App Agent credential is not JSON at all — it is the
// credential token itself; see AppAgentFromToken.)
type Credentials struct {
	BaseURL               string `json:"baseUrl,omitempty"`
	Endpoint              string `json:"endpoint,omitempty"`
	ApplicationID         string `json:"applicationId,omitempty"`
	AppSpaceID            string `json:"appSpaceId,omitempty"`
	AppAgentID            string `json:"appAgentId,omitempty"`
	ServiceAccountID      string `json:"serviceAccountId,omitempty"`
	PrivateKeyPKCS8Base64 string `json:"privateKeyPKCS8Base64,omitempty"`
	PrivateKeyPKCS8       string `json:"privateKeyPKCS8,omitempty"`
	TokenLifetime         string `json:"tokenLifetime,omitempty"`
	// Token is a pre-issued credential token, sent verbatim instead of minting
	// a JWT from the private key.
	Token         string          `json:"token,omitempty"`
	PrivateKeyJWK json.RawMessage `json:"privateKeyJWK,omitempty"`
}

// hasPrivateKey reports whether any private-key form is present.
func (c *Credentials) hasPrivateKey() bool {
	return len(c.PrivateKeyJWK) > 0 || c.PrivateKeyPKCS8 != "" || c.PrivateKeyPKCS8Base64 != ""
}

// ClientID returns the JWT subject/issuer for the credential.
func (c *Credentials) ClientID() (string, error) {
	switch {
	case c.AppAgentID != "":
		return c.AppAgentID, nil
	case c.ServiceAccountID != "":
		return c.ServiceAccountID, nil
	default:
		return "", errors.New("indykite: missing client ID (appAgentId or serviceAccountId)")
	}
}

// Loader resolves credentials. It returns (nil, nil) when this source has
// nothing to offer, allowing the next loader in a chain to run.
type Loader func(ctx context.Context) (*Credentials, error)

// StaticJSON builds a Loader from a raw credentials JSON document.
func StaticJSON(credentialsJSON []byte) Loader {
	return func(_ context.Context) (*Credentials, error) {
		return unmarshalCredentials(credentialsJSON)
	}
}

// StaticCredentials builds a Loader from an in-memory credentials value.
func StaticCredentials(c *Credentials) Loader {
	return func(_ context.Context) (*Credentials, error) { return c, nil }
}

// ServiceAccountEnvLoader reads INDYKITE_SERVICE_ACCOUNT_CREDENTIALS[_FILE].
func ServiceAccountEnvLoader(_ context.Context) (*Credentials, error) {
	return loadFromEnv(EnvServiceAccountCredentials)
}

func loadFromEnv(env string) (*Credentials, error) {
	data, err := lookupEnvCredential(env)
	if err != nil {
		return nil, err
	}
	return unmarshalCredentials(data)
}

func unmarshalCredentials(credentialJSON []byte) (*Credentials, error) {
	trimmed := bytes.TrimSpace(credentialJSON)
	if len(trimmed) == 0 {
		return nil, nil
	}
	c := &Credentials{}
	if err := json.Unmarshal(trimmed, c); err != nil {
		return nil, fmt.Errorf("indykite: invalid credentials JSON: %w", err)
	}
	return c, nil
}

func lookupEnvCredential(env string) ([]byte, error) {
	if v, ok := os.LookupEnv(env); ok && v != "" {
		return []byte(v), nil
	}
	if v, ok := os.LookupEnv(env + "_FILE"); ok && v != "" {
		return readFileWithLimit(v, 10)
	}
	return nil, nil
}

func readFileWithLimit(path string, limitKB int64) ([]byte, error) {
	clean := filepath.Clean(path)
	info, err := os.Stat(clean)
	if err != nil {
		return nil, err
	}
	if info.Size() > limitKB*1024 {
		return nil, fmt.Errorf("indykite: credential file %q exceeds %dKB", clean, limitKB)
	}
	return os.ReadFile(clean)
}
