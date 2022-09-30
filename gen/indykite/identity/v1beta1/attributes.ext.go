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
	"fmt"

	"github.com/pborman/uuid"
	"google.golang.org/protobuf/proto"

	objects "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

type PropertyIssuerPrefix string

const (
	PropertyRawPrefix         PropertyIssuerPrefix = "raw"
	PropertyApplicationPrefix PropertyIssuerPrefix = "app"
	PropertyAppSpacePrefix    PropertyIssuerPrefix = "asp"
)

type PropertyBatchOperations []*PropertyBatchOperation

func (PropertyBatchOperations) idInArray(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func (x *Property) HasSameDefinition(other *Property) bool {
	if other == nil || other.Definition == nil {
		return false
	}
	return proto.Equal(x.GetDefinition(), other.GetDefinition())
}

func (x *Property) GetStringValue() (string, bool) {
	if val, ok := x.Value.(*Property_ObjectValue); ok {
		if objVal, ok := val.ObjectValue.GetValue().(*objects.Value_StringValue); ok {
			return objVal.StringValue, true
		}
	}
	return "", false
}

func (x *Property) GetIntValue() (int64, bool) {
	if val, ok := x.Value.(*Property_ObjectValue); ok {
		if objVal, ok := val.ObjectValue.GetValue().(*objects.Value_IntegerValue); ok {
			return objVal.IntegerValue, true
		}
	}
	return 0, false
}

func (x *Property) GetAnyValue(msg proto.Message) error {
	if val, ok := x.Value.(*Property_ObjectValue); ok {
		if anyVal, ok := val.ObjectValue.GetValue().(*objects.Value_AnyValue); ok {
			return anyVal.AnyValue.UnmarshalTo(msg)
		} else if _, ok = val.ObjectValue.GetValue().(*objects.Value_NullValue); ok {
			return nil
		}
	}
	return errors.New("invalid AnyPB value")
}

func (x *Property) GetMapValue() (val map[string]interface{}, err error) {
	if val, ok := x.Value.(*Property_ObjectValue); ok {
		if mapVal, ok := val.ObjectValue.GetValue().(*objects.Value_MapValue); ok {
			return objects.ToMap(mapVal.MapValue.Fields)
		} else if _, ok = val.ObjectValue.GetValue().(*objects.Value_NullValue); ok {
			return nil, nil
		}
	}
	return nil, errors.New("invalid map value value")
}

// validateValue checks if an operation value is properly set.
//
// It checks if the property values is not nil.
func (PropertyBatchOperations) validateValue(value isProperty_Value, enableNil bool) error {
	switch propertyValue := value.(type) {
	case nil:
		if !enableNil {
			return errors.New("value cannot be nil")
		}
	case *Property_ObjectValue:
		if propertyValue.ObjectValue == nil || propertyValue.ObjectValue.Value == nil {
			return errors.New("object value must be specified")
		}
	case *Property_ReferenceValue:
		if len(propertyValue.ReferenceValue) == 0 {
			return errors.New("reference value must be specified")
		}
	default:
		return errors.New("invalid value type")
	}
	return nil
}

func (PropertyBatchOperations) validateMeta(meta *PropertyMetadata, isTrusted bool) error {
	if meta == nil {
		return nil
	}

	if meta.AssuranceLevel > AssuranceLevel_ASSURANCE_LEVEL_LOW && meta.Verifier == "" {
		return errors.New("verifier is required when assurance level is Substantial or more")
	}
	if meta.Verifier != "" && meta.AssuranceLevel < AssuranceLevel_ASSURANCE_LEVEL_SUBSTANTIAL {
		return errors.New("assurance level must be Substantial or more when verifier is set")
	}
	if meta.Verifier == "" && meta.VerificationTime != nil {
		return errors.New("verification time must be set together with verifier")
	}

	// If is non-trusted, we don't care about issuer as it will be ignored
	if isTrusted && meta.Issuer != "" {
		switch {
		case len(meta.Issuer) == 26 &&
			(meta.Issuer[:3] == string(PropertyApplicationPrefix) || meta.Issuer[:3] == string(PropertyAppSpacePrefix)):
			if uuid.UUID(meta.Issuer[4:]).Variant() != uuid.RFC4122 {
				return errors.New("issuer must be RFC4122 variant UUID")
			}
		case len(meta.Issuer) >= 7 && meta.Issuer[:3] == string(PropertyRawPrefix):
			// This is valid
		default:
			return errors.New("issuer must starts with 'app:' or 'asp:' followed by trimmed base64URL UUID," +
				" or starts with 'raw:' followed by name with minimum length of 3 characters")
		}
	}

	return nil
}

// Validate checks the PropertyBatchOperations values. All public calls are considered non-trusted.
//
// It checks if the operations represent a valid collection of property updates.
func (x PropertyBatchOperations) Validate(isTrusted bool) error {
	_, _, err := x.PreValidate(isTrusted)
	return err
}

// PreValidate checks the PropertyBatchOperations values. All public calls are considered non-trusted.
//
// It checks if the operations represent a valid collection of property updates.
func (x PropertyBatchOperations) PreValidate(isTrusted bool) ([]string, []string, error) {
	if len(x) == 0 {
		return nil, nil, errors.New("empty batch operation")
	}

	// This logic cannot be really tight to internal property logic, because PreValidate is used also in public part
	var replaced, removed []string
	for i, v := range x {
		err := v.Validate()
		if err != nil {
			return nil, nil, fmt.Errorf("invalid operation at index %d: %v", i, err)
		}
		switch bo := v.Operation.(type) {
		case *PropertyBatchOperation_Add:
			if bo.Add == nil || bo.Add.Definition == nil {
				return nil, nil, fmt.Errorf("invalid add operation at position: %d", i)
			}
			if bo.Add.Id != "" {
				return nil, nil, fmt.Errorf("cannot set custom ID for add operation at position: %d", i)
			}
			if isTrusted && (bo.Add.Meta == nil || bo.Add.Meta.Issuer == "") {
				return nil, nil, fmt.Errorf("issuer is required for trusted operation at position: %d", i)
			}
			if err := x.validateMeta(bo.Add.Meta, isTrusted); err != nil {
				return nil, nil, fmt.Errorf("%s at position: %d", err.Error(), i)
			}
			if err := x.validateValue(bo.Add.Value, false); err != nil {
				return nil, nil, fmt.Errorf("%s at position: %d", err.Error(), i)
			}
		case *PropertyBatchOperation_Replace:
			if bo.Replace == nil || bo.Replace.Id == "" {
				return nil, nil, fmt.Errorf("invalid replace operation at position: %d", i)
			}
			if bo.Replace.Value == nil && bo.Replace.Meta == nil {
				return nil, nil, fmt.Errorf("at least Value or Meta must be set at position: %d", i)
			}
			if err := x.validateValue(bo.Replace.Value, true); err != nil {
				return nil, nil, fmt.Errorf("%s at position: %d", err.Error(), i)
			}
			if err := x.validateMeta(bo.Replace.Meta, isTrusted); err != nil {
				return nil, nil, fmt.Errorf("%s at position: %d", err.Error(), i)
			}
			if x.idInArray(replaced, bo.Replace.Id) {
				return nil, nil, fmt.Errorf("replacing same property more than once is not allowed at position: %d", i)
			}
			if x.idInArray(removed, bo.Replace.Id) {
				return nil, nil, fmt.Errorf("replacing property previously deleted is not allowed at position: %d", i)
			}
			replaced = append(replaced, bo.Replace.Id)
		case *PropertyBatchOperation_Remove:
			if bo.Remove == nil || bo.Remove.Id == "" {
				return nil, nil, fmt.Errorf("invalid remove operation at position: %d", i)
			}
			if x.idInArray(replaced, bo.Remove.Id) {
				return nil, nil, fmt.Errorf("deleting previously replaced property is not allowed at position: %d", i)
			}
			if x.idInArray(removed, bo.Remove.Id) {
				return nil, nil, fmt.Errorf("deleting property multiple times is not allowed at position: %d", i)
			}
			removed = append(removed, bo.Remove.Id)
		default:
			return nil, nil, fmt.Errorf("invalid operation at position: %d", i)
		}
	}
	return replaced, removed, nil
}
