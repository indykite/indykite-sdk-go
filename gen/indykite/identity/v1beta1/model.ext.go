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

package identityv1beta1

import (
	"errors"

	"github.com/pborman/uuid"
)

// Verify checks if the *DigitalTwin has RFC4122 TenantId and Id value.
func (x *DigitalTwin) Verify() (tenantID uuid.UUID, digitalTwinID uuid.UUID, _ error) {
	if x == nil {
		return nil, nil, errors.New("DigitalTwin value is required")
	}
	tenantID, digitalTwinID = x.TenantId, x.Id
	if tenantID.Variant() != uuid.RFC4122 {
		return nil, nil, errors.New("tenantId in DigitalTwin must be RFC4122 variant UUID")
	}
	if digitalTwinID.Variant() != uuid.RFC4122 {
		return nil, nil, errors.New("id in DigitalTwin must be RFC4122 variant UUID")
	}
	return tenantID, digitalTwinID, nil
}

func (*DigitalTwin) IsNode() {}
