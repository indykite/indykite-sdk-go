// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/objects/v1beta1/struct.proto

package objectsv1beta1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort

	_ = structpb.NullValue(0)
)

// Validate checks the field values on Value with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Value) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Value with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ValueMultiError, or nil if none found.
func (m *Value) ValidateAll() error {
	return m.validate(true)
}

func (m *Value) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	switch m.Value.(type) {

	case *Value_NullValue:
		// no validation rules for NullValue

	case *Value_BoolValue:
		// no validation rules for BoolValue

	case *Value_IntegerValue:
		// no validation rules for IntegerValue

	case *Value_UnsignedIntegerValue:
		// no validation rules for UnsignedIntegerValue

	case *Value_DoubleValue:
		// no validation rules for DoubleValue

	case *Value_AnyValue:

		if all {
			switch v := interface{}(m.GetAnyValue()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "AnyValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "AnyValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAnyValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "AnyValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Value_ValueTime:

		if all {
			switch v := interface{}(m.GetValueTime()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "ValueTime",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "ValueTime",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetValueTime()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "ValueTime",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Value_DurationValue:

		if all {
			switch v := interface{}(m.GetDurationValue()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "DurationValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "DurationValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDurationValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "DurationValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Value_IdentifierValue:

		if all {
			switch v := interface{}(m.GetIdentifierValue()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "IdentifierValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "IdentifierValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetIdentifierValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "IdentifierValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Value_StringValue:
		// no validation rules for StringValue

	case *Value_BytesValue:
		// no validation rules for BytesValue

	case *Value_GeoPointValue:

		if all {
			switch v := interface{}(m.GetGeoPointValue()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "GeoPointValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "GeoPointValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetGeoPointValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "GeoPointValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Value_ArrayValue:

		if all {
			switch v := interface{}(m.GetArrayValue()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "ArrayValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "ArrayValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetArrayValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "ArrayValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Value_MapValue:

		if all {
			switch v := interface{}(m.GetMapValue()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "MapValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ValueValidationError{
						field:  "MapValue",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMapValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ValueValidationError{
					field:  "MapValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ValueMultiError(errors)
	}

	return nil
}

// ValueMultiError is an error wrapping multiple validation errors returned by
// Value.ValidateAll() if the designated constraints aren't met.
type ValueMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValueMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValueMultiError) AllErrors() []error { return m }

// ValueValidationError is the validation error returned by Value.Validate if
// the designated constraints aren't met.
type ValueValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValueValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValueValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValueValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValueValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValueValidationError) ErrorName() string { return "ValueValidationError" }

// Error satisfies the builtin error interface
func (e ValueValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValue.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValueValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValueValidationError{}

// Validate checks the field values on ArrayValue with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ArrayValue) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ArrayValue with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ArrayValueMultiError, or
// nil if none found.
func (m *ArrayValue) ValidateAll() error {
	return m.validate(true)
}

func (m *ArrayValue) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetValues() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ArrayValueValidationError{
						field:  fmt.Sprintf("Values[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ArrayValueValidationError{
						field:  fmt.Sprintf("Values[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ArrayValueValidationError{
					field:  fmt.Sprintf("Values[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ArrayValueMultiError(errors)
	}

	return nil
}

// ArrayValueMultiError is an error wrapping multiple validation errors
// returned by ArrayValue.ValidateAll() if the designated constraints aren't met.
type ArrayValueMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ArrayValueMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ArrayValueMultiError) AllErrors() []error { return m }

// ArrayValueValidationError is the validation error returned by
// ArrayValue.Validate if the designated constraints aren't met.
type ArrayValueValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ArrayValueValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ArrayValueValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ArrayValueValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ArrayValueValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ArrayValueValidationError) ErrorName() string { return "ArrayValueValidationError" }

// Error satisfies the builtin error interface
func (e ArrayValueValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sArrayValue.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ArrayValueValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ArrayValueValidationError{}

// Validate checks the field values on MapValue with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *MapValue) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MapValue with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MapValueMultiError, or nil
// if none found.
func (m *MapValue) ValidateAll() error {
	return m.validate(true)
}

func (m *MapValue) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	{
		sorted_keys := make([]string, len(m.GetFields()))
		i := 0
		for key := range m.GetFields() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetFields()[key]
			_ = val

			// no validation rules for Fields[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, MapValueValidationError{
							field:  fmt.Sprintf("Fields[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, MapValueValidationError{
							field:  fmt.Sprintf("Fields[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return MapValueValidationError{
						field:  fmt.Sprintf("Fields[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if len(errors) > 0 {
		return MapValueMultiError(errors)
	}

	return nil
}

// MapValueMultiError is an error wrapping multiple validation errors returned
// by MapValue.ValidateAll() if the designated constraints aren't met.
type MapValueMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MapValueMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MapValueMultiError) AllErrors() []error { return m }

// MapValueValidationError is the validation error returned by
// MapValue.Validate if the designated constraints aren't met.
type MapValueValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MapValueValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MapValueValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MapValueValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MapValueValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MapValueValidationError) ErrorName() string { return "MapValueValidationError" }

// Error satisfies the builtin error interface
func (e MapValueValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMapValue.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MapValueValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MapValueValidationError{}
