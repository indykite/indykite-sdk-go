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

package objectsv1beta1

import (
	"github.com/leodido/go-urn"
	"github.com/multiformats/go-base36"
	"github.com/pborman/uuid"
	"google.golang.org/protobuf/proto"
)

// NewIdentifier generates a new random (Version 4) UUID identifier.
func NewIdentifier() *Identifier {
	return &Identifier{
		Id: &Identifier_IdBytes{
			IdBytes: uuid.NewRandom(),
		},
	}
}

// FromUUID converts UUID to *Identifier.
//
// Converts from:
// * String
// * UUID
// * *Identifier
//
// Returns nil if failed to parse.
//
func FromUUID(from interface{}) *Identifier {
	var id uuid.UUID
	switch v := from.(type) {
	case string:
		id = uuid.Parse(v)
	case []byte:
		id = v
		if id.Variant() != uuid.RFC4122 {
			return nil
		}
	case *Identifier:
		return v
	default:
		return nil
	}
	return &Identifier{
		Id: &Identifier_IdBytes{
			IdBytes: id,
		},
	}
}

func (x *Identifier) AsString() string {
	switch v := x.GetId().(type) {
	case *Identifier_IdString:
		return v.IdString
	case *Identifier_IdBytes:
		id := uuid.UUID(v.IdBytes)
		if id.Variant() == uuid.RFC4122 {
			return id.String()
		}
	}
	return ""
}

func (x *Identifier) AsUUID() uuid.UUID {
	switch v := x.GetId().(type) {
	case *Identifier_IdString:
		return uuid.Parse(v.IdString)
	case *Identifier_IdBytes:
		id := uuid.UUID(v.IdBytes)
		if id.Variant() == uuid.RFC4122 {
			return id
		}
	}
	return nil
}

func (x *ObjectReference) AsString() string {
	data, err := proto.Marshal(x)
	if err != nil {
		return ""
	}
	return (&urn.URN{
		ID: "indy",
		SS: base36.EncodeToStringLc(data),
	}).String()
}

func ParseReferenceBytes(raw []byte) *ObjectReference {
	if u, ok := urn.Parse(raw); ok {
		if u.ID == "indy" {
			var ref ObjectReference
			data, err := base36.DecodeString(u.SS)
			if err == nil {
				if err = proto.Unmarshal(data, &ref); err == nil {
					return &ref
				}
			}
		}
	}
	return nil
}

func ParseReference(raw string) *ObjectReference {
	return ParseReferenceBytes([]byte(raw))
}
