// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/entitymatching/v1beta1/model.proto

package entitymatchingv1beta1

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
)

// Validate checks the field values on PropertyMapping with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PropertyMapping) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PropertyMapping with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PropertyMappingMultiError, or nil if none found.
func (m *PropertyMapping) ValidateAll() error {
	return m.validate(true)
}

func (m *PropertyMapping) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for SourceNodeType

	// no validation rules for SourceNodeProperty

	// no validation rules for TargetNodeType

	// no validation rules for TargetNodeProperty

	if m.GetSimilarityScoreCutoff() != 0 {

		if val := m.GetSimilarityScoreCutoff(); val < 0 || val > 1 {
			err := PropertyMappingValidationError{
				field:  "SimilarityScoreCutoff",
				reason: "value must be inside range [0, 1]",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return PropertyMappingMultiError(errors)
	}

	return nil
}

// PropertyMappingMultiError is an error wrapping multiple validation errors
// returned by PropertyMapping.ValidateAll() if the designated constraints
// aren't met.
type PropertyMappingMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PropertyMappingMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PropertyMappingMultiError) AllErrors() []error { return m }

// PropertyMappingValidationError is the validation error returned by
// PropertyMapping.Validate if the designated constraints aren't met.
type PropertyMappingValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PropertyMappingValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PropertyMappingValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PropertyMappingValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PropertyMappingValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PropertyMappingValidationError) ErrorName() string { return "PropertyMappingValidationError" }

// Error satisfies the builtin error interface
func (e PropertyMappingValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPropertyMapping.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PropertyMappingValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PropertyMappingValidationError{}

// Validate checks the field values on CustomPropertyMappings with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CustomPropertyMappings) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CustomPropertyMappings with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CustomPropertyMappingsMultiError, or nil if none found.
func (m *CustomPropertyMappings) ValidateAll() error {
	return m.validate(true)
}

func (m *CustomPropertyMappings) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for SourceNodeProperty

	// no validation rules for TargetNodeProperty

	if len(errors) > 0 {
		return CustomPropertyMappingsMultiError(errors)
	}

	return nil
}

// CustomPropertyMappingsMultiError is an error wrapping multiple validation
// errors returned by CustomPropertyMappings.ValidateAll() if the designated
// constraints aren't met.
type CustomPropertyMappingsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CustomPropertyMappingsMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CustomPropertyMappingsMultiError) AllErrors() []error { return m }

// CustomPropertyMappingsValidationError is the validation error returned by
// CustomPropertyMappings.Validate if the designated constraints aren't met.
type CustomPropertyMappingsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CustomPropertyMappingsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CustomPropertyMappingsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CustomPropertyMappingsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CustomPropertyMappingsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CustomPropertyMappingsValidationError) ErrorName() string {
	return "CustomPropertyMappingsValidationError"
}

// Error satisfies the builtin error interface
func (e CustomPropertyMappingsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCustomPropertyMappings.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CustomPropertyMappingsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CustomPropertyMappingsValidationError{}
