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
	"fmt"
	"regexp"

	"github.com/pborman/uuid"

	"github.com/indykite/indykite-sdk-go/errors"
)

const (
	rfc1035NameTemplate = "[a-z](?:[-a-z0-9]{%d,%d}[a-z0-9])"
)

var nameCheck = regexp.MustCompile(fmt.Sprintf(rfc1035NameTemplate, 0, 252))

// IsUUIDv4 checks if id is valid RFC4122 varian UUID.
func IsUUIDv4(key string, id []byte) error {
	if uuid.UUID(id).Variant() != uuid.RFC4122 {
		return errors.NewInvalidArgumentError("expected UUID RFC4122 variant")
	}
	return nil
}

// IsValidName checks if name is valid RFC1035 string.
//
// Value can have lowercase letters, digits, or hyphens. It must start with a lowercase letter
// and end with a letter or number. The minimum length is 2 and the max is 254.
func IsValidName(name string) error {
	if !nameCheck.MatchString(name) {
		return errors.NewInvalidArgumentError("name value must be valid RFC1035 string with length 2-254")
	}
	return nil
}

// ContainsLabel check if string array contains given node label.
func ContainsLabel(arr []string, searchFor string) bool {
	for _, v := range arr {
		if searchFor == v {
			return true
		}
	}
	return false
}
