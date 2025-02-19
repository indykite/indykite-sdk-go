// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/auditsink/v1beta1/ai.proto

package auditsinkv1beta1

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

// Validate checks the field values on IsChangePoint with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *IsChangePoint) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IsChangePoint with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in IsChangePointMultiError, or
// nil if none found.
func (m *IsChangePoint) ValidateAll() error {
	return m.validate(true)
}

func (m *IsChangePoint) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetResolvedDigitalTwin()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "ResolvedDigitalTwin",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "ResolvedDigitalTwin",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetResolvedDigitalTwin()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IsChangePointValidationError{
				field:  "ResolvedDigitalTwin",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetRequest()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "Request",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "Request",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRequest()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IsChangePointValidationError{
				field:  "Request",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetChangePointDetected()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "ChangePointDetected",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "ChangePointDetected",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetChangePointDetected()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IsChangePointValidationError{
				field:  "ChangePointDetected",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetResponse()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "Response",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "Response",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetResponse()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IsChangePointValidationError{
				field:  "Response",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetEvaluationTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "EvaluationTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IsChangePointValidationError{
					field:  "EvaluationTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEvaluationTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IsChangePointValidationError{
				field:  "EvaluationTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return IsChangePointMultiError(errors)
	}

	return nil
}

// IsChangePointMultiError is an error wrapping multiple validation errors
// returned by IsChangePoint.ValidateAll() if the designated constraints
// aren't met.
type IsChangePointMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IsChangePointMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IsChangePointMultiError) AllErrors() []error { return m }

// IsChangePointValidationError is the validation error returned by
// IsChangePoint.Validate if the designated constraints aren't met.
type IsChangePointValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IsChangePointValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IsChangePointValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IsChangePointValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IsChangePointValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IsChangePointValidationError) ErrorName() string { return "IsChangePointValidationError" }

// Error satisfies the builtin error interface
func (e IsChangePointValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIsChangePoint.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IsChangePointValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IsChangePointValidationError{}

// Validate checks the field values on IsChangePoint_Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IsChangePoint_Request) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IsChangePoint_Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IsChangePoint_RequestMultiError, or nil if none found.
func (m *IsChangePoint_Request) ValidateAll() error {
	return m.validate(true)
}

func (m *IsChangePoint_Request) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Resource

	// no validation rules for Action

	if len(errors) > 0 {
		return IsChangePoint_RequestMultiError(errors)
	}

	return nil
}

// IsChangePoint_RequestMultiError is an error wrapping multiple validation
// errors returned by IsChangePoint_Request.ValidateAll() if the designated
// constraints aren't met.
type IsChangePoint_RequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IsChangePoint_RequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IsChangePoint_RequestMultiError) AllErrors() []error { return m }

// IsChangePoint_RequestValidationError is the validation error returned by
// IsChangePoint_Request.Validate if the designated constraints aren't met.
type IsChangePoint_RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IsChangePoint_RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IsChangePoint_RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IsChangePoint_RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IsChangePoint_RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IsChangePoint_RequestValidationError) ErrorName() string {
	return "IsChangePoint_RequestValidationError"
}

// Error satisfies the builtin error interface
func (e IsChangePoint_RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIsChangePoint_Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IsChangePoint_RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IsChangePoint_RequestValidationError{}

// Validate checks the field values on IsChangePoint_ChangePointDetection with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *IsChangePoint_ChangePointDetection) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IsChangePoint_ChangePointDetection
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// IsChangePoint_ChangePointDetectionMultiError, or nil if none found.
func (m *IsChangePoint_ChangePointDetection) ValidateAll() error {
	return m.validate(true)
}

func (m *IsChangePoint_ChangePointDetection) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for IsChange

	// no validation rules for Explanation

	// no validation rules for ChangeScore

	if len(errors) > 0 {
		return IsChangePoint_ChangePointDetectionMultiError(errors)
	}

	return nil
}

// IsChangePoint_ChangePointDetectionMultiError is an error wrapping multiple
// validation errors returned by
// IsChangePoint_ChangePointDetection.ValidateAll() if the designated
// constraints aren't met.
type IsChangePoint_ChangePointDetectionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IsChangePoint_ChangePointDetectionMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IsChangePoint_ChangePointDetectionMultiError) AllErrors() []error { return m }

// IsChangePoint_ChangePointDetectionValidationError is the validation error
// returned by IsChangePoint_ChangePointDetection.Validate if the designated
// constraints aren't met.
type IsChangePoint_ChangePointDetectionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IsChangePoint_ChangePointDetectionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IsChangePoint_ChangePointDetectionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IsChangePoint_ChangePointDetectionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IsChangePoint_ChangePointDetectionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IsChangePoint_ChangePointDetectionValidationError) ErrorName() string {
	return "IsChangePoint_ChangePointDetectionValidationError"
}

// Error satisfies the builtin error interface
func (e IsChangePoint_ChangePointDetectionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIsChangePoint_ChangePointDetection.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IsChangePoint_ChangePointDetectionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IsChangePoint_ChangePointDetectionValidationError{}

// Validate checks the field values on IsChangePoint_Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IsChangePoint_Response) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IsChangePoint_Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IsChangePoint_ResponseMultiError, or nil if none found.
func (m *IsChangePoint_Response) ValidateAll() error {
	return m.validate(true)
}

func (m *IsChangePoint_Response) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for DecisionTime

	if len(errors) > 0 {
		return IsChangePoint_ResponseMultiError(errors)
	}

	return nil
}

// IsChangePoint_ResponseMultiError is an error wrapping multiple validation
// errors returned by IsChangePoint_Response.ValidateAll() if the designated
// constraints aren't met.
type IsChangePoint_ResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IsChangePoint_ResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IsChangePoint_ResponseMultiError) AllErrors() []error { return m }

// IsChangePoint_ResponseValidationError is the validation error returned by
// IsChangePoint_Response.Validate if the designated constraints aren't met.
type IsChangePoint_ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IsChangePoint_ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IsChangePoint_ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IsChangePoint_ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IsChangePoint_ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IsChangePoint_ResponseValidationError) ErrorName() string {
	return "IsChangePoint_ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e IsChangePoint_ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIsChangePoint_Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IsChangePoint_ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IsChangePoint_ResponseValidationError{}
