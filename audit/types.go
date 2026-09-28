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

package audit

import "encoding/json"

// ListRequest selects one page of a listing endpoint.
type ListRequest struct {
	// ProjectID is the project (app space) GID. Required: it must be the
	// project the App Agent credential belongs to, or the platform answers 403.
	ProjectID string
	// Cursor is the NextCursor of a previous page; empty starts from the
	// beginning (the newest checkpoint for ListCheckpoints).
	Cursor string
	// PageSize is the maximum number of items to return. Zero or less uses the
	// server default (50); the server caps it at 50.
	PageSize int
}

// Page is one page of a listing endpoint.
type Page[T any] struct {
	// NextCursor resumes after this page; empty when HasMore is false.
	NextCursor string `json:"next_cursor"`
	Items      []T    `json:"items"`
	HasMore    bool   `json:"has_more"`
}

// LogEntry is one signed chain batch (GET /audit/v1/logs).
type LogEntry struct {
	BatchID    string `json:"batch_id"`
	ProjectID  string `json:"project_id"`
	Hash       string `json:"hash"`
	Signature  string `json:"signature"`
	KID        string `json:"kid"`
	Alg        string `json:"alg"`
	ChainHash  string `json:"chain_hash"`
	ManifestID string `json:"manifest_id"`
	// Data is the batch payload (the audit events) exactly as stored.
	Data     json.RawMessage `json:"data"`
	Sequence uint64          `json:"sequence"`
}

// Manifest is one signed chain manifest (GET /audit/v1/manifests): the
// chain-linking metadata of a batch, without the batch payload.
type Manifest struct {
	ManifestID string `json:"manifest_id"`
	BatchID    string `json:"batch_id"`
	BatchURI   string `json:"batch_uri"`
	ProjectID  string `json:"project_id"`
	PrevHash   string `json:"prev_hash"`
	HeadHash   string `json:"head_hash"`
	DataHash   string `json:"data_hash"`
	KID        string `json:"kid"`
	Alg        string `json:"alg"`
	Signature  string `json:"signature"`
	CreatedAt  string `json:"created_at"`
	Sequence   uint64 `json:"sequence"`
}

// Checkpoint is one signed project checkpoint (GET /audit/v1/checkpoints), the
// trust anchor a chain is verified against.
type Checkpoint struct {
	CheckpointID string `json:"checkpoint_id"`
	ProjectID    string `json:"project_id"`
	HeadHash     string `json:"head_hash"`
	CreatedAt    string `json:"created_at"`
	Signature    string `json:"signature"`
	KID          string `json:"kid"`
	Alg          string `json:"alg"`
	Sequence     uint64 `json:"sequence"`
}

// KeySet is the JWKS returned by GET /audit/.well-known/jwks.json.
type KeySet struct {
	Keys []Key `json:"keys"`
}

// Key is one verification key of a KeySet. It carries no "alg": signatures are
// ECDSA P-256/SHA-256 in ASN.1 DER encoding (verify with ecdsa.VerifyASN1), not
// the raw R||S encoding a JOSE "ES256" verifier expects.
type Key struct {
	Kty string `json:"kty"`
	Crv string `json:"crv,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`
	Kid string `json:"kid,omitempty"`
	Use string `json:"use,omitempty"`
}
