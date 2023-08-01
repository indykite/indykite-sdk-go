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

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type (
	CredentialsConfig struct {
		BaseURL         string `json:"baseUrl,omitempty"         yaml:"base_url,omitempty"`
		ApplicationID   string `json:"applicationId,omitempty"   yaml:"application_id,omitempty"`
		AppSpaceID      string `json:"appSpaceId,omitempty"      yaml:"app_space_id,omitempty"`
		DefaultTenantID string `json:"defaultTenantId,omitempty" yaml:"default_tenant_id,omitempty"`

		AppAgentID            string          `json:"appAgentId,omitempty"            yaml:"app_agent_id,omitempty"`
		ServiceAccountID      string          `json:"serviceAccountId,omitempty"      yaml:"service_account_id,omitempty"` //nolint:lll
		Endpoint              string          `json:"endpoint,omitempty"              yaml:"endpoint,omitempty"`
		PrivateKeyPKCS8Base64 string          `json:"privateKeyPKCS8Base64,omitempty" yaml:"private_key_pkcs8_base64,omitempty"` //nolint:lll
		PrivateKeyPKCS8       string          `json:"privateKeyPKCS8,omitempty"       yaml:"private_key_pkcs8,omitempty"`        //nolint:lll
		TokenLifetime         string          `json:"tokenLifetime,omitempty"         yaml:"token_lifetime,omitempty"`
		PrivateKeyJWK         json.RawMessage `json:"privateKeyJWK,omitempty"         yaml:"private_key_jwk,omitempty"`
	}

	CredentialsLoader func(ctx context.Context) (*CredentialsConfig, error)
)

func DefaultEnvironmentLoader(_ context.Context) (*CredentialsConfig, error) {
	data, err := lookupEnvCredentialVariables("INDYKITE_APPLICATION_CREDENTIALS")
	if err != nil {
		return nil, err
	}
	return UnmarshalCredentialConfig(data)
}

func StaticCredentialsJSON(credentialsJSON []byte) CredentialsLoader {
	return func(ctx context.Context) (*CredentialsConfig, error) {
		return UnmarshalCredentialConfig(credentialsJSON)
	}
}

func UnmarshalCredentialConfig(credentialJSON []byte) (*CredentialsConfig, error) {
	if len(credentialJSON) > 0 {
		var cfg = &CredentialsConfig{}
		err := json.Unmarshal(credentialJSON, cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return nil, nil
}

func StaticCredentialConfig(config *CredentialsConfig) CredentialsLoader {
	return func(ctx context.Context) (*CredentialsConfig, error) {
		return config, nil
	}
}

// lookupEnvCredentialVariables tries to load the client credentials from environment configurations.
func lookupEnvCredentialVariables(env string) ([]byte, error) {
	if v, ok := os.LookupEnv(env); ok && v != "" {
		return []byte(v), nil
	}
	if v, ok := os.LookupEnv(env + "_FILE"); ok && v != "" {
		return readFileWithLimit(v, 10)
	}
	return nil, nil
}

func readFileWithLimit(filePath string, limitSizeInKb int) ([]byte, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	var data []byte
	data, err = readWithLimit(filePath, limitSizeInKb, file)
	if err != nil {
		return nil, err
	}
	return data, file.Close()
}

func readWithLimit(filePath string, limitSizeInKb int, file *os.File) ([]byte, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileInfo.Size()
	if filesize > int64(limitSizeInKb)*1024 {
		return nil, fmt.Errorf(
			"credential file '%s' exceeded maximum size of %dKByte, size is %dKByte",
			filePath,
			limitSizeInKb,
			filesize/1024,
		)
	}
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
