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

package identity

import (
	"encoding/base64"
	"time"

	guuid "github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pborman/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/indykite/jarvis-sdk-go/errors"
)

// ParseUUID parse the raw input string and check if it's a RFC4122 variant UUID.
func ParseUUID(raw string) (uuid.UUID, error) {
	var (
		uid guuid.UUID
		err error
	)
	switch len(raw) {
	case 22:
		var b []byte
		b, err = base64.RawURLEncoding.DecodeString(raw)
		if err == nil {
			uid, err = guuid.FromBytes(b)
		}
	default:
		uid, err = guuid.Parse(raw)
	}
	if err != nil {
		return nil, err
	}
	if uid.Variant() != guuid.RFC4122 {
		return nil, errors.NewInvalidArgumentError("invalid UUID RFC4122 variant")
	}
	return uid[:], nil
}

func verifyTokenFormat(bearerToken string) error {
	_, err := jwt.ParseString(bearerToken, jwt.WithValidate(true), jwt.WithAcceptableSkew(time.Second))
	if err != nil {
		return errors.NewWithCause(codes.InvalidArgument, err, "invalid token format")
	}
	return nil
}

func optionalTime(at time.Time) *timestamppb.Timestamp {
	if at.IsZero() {
		return nil
	}
	return timestamppb.New(at)
}
