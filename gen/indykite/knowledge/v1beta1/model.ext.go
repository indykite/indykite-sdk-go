// Copyright (c) 2023 IndyKite
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

package knowledgev1beta1

import (
	"strings"

	objects "github.com/indykite/indykite-sdk-go/gen/indykite/objects/v1beta1"
)

// GetProperty is a helper function to get the property with a certain key from the node's list of properties, similar
// to a map lookup. It returns the value and true if found, nil and false if not found.
func (x *Node) GetProperty(key string) (*objects.Value, bool) {
	for _, p := range x.Properties {
		if strings.EqualFold(p.GetKey(), key) {
			return p.GetValue(), true
		}
	}
	return nil, false
}
