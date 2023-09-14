// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/knowledge/v1beta1/identity_knowledge_api.proto

package knowledgev1beta1

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

// Validate checks the field values on IdentityKnowledgeRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IdentityKnowledgeRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IdentityKnowledgeRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IdentityKnowledgeRequestMultiError, or nil if none found.
func (m *IdentityKnowledgeRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *IdentityKnowledgeRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if _, ok := _IdentityKnowledgeRequest_Operation_InLookup[m.GetOperation()]; !ok {
		err := IdentityKnowledgeRequestValidationError{
			field:  "Operation",
			reason: "value must be in list [OPERATION_READ]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := Operation_name[int32(m.GetOperation())]; !ok {
		err := IdentityKnowledgeRequestValidationError{
			field:  "Operation",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetPath()) > 512000 {
		err := IdentityKnowledgeRequestValidationError{
			field:  "Path",
			reason: "value length must be at most 512000 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetConditions()) > 512000 {
		err := IdentityKnowledgeRequestValidationError{
			field:  "Conditions",
			reason: "value length must be at most 512000 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return IdentityKnowledgeRequestMultiError(errors)
	}

	return nil
}

// IdentityKnowledgeRequestMultiError is an error wrapping multiple validation
// errors returned by IdentityKnowledgeRequest.ValidateAll() if the designated
// constraints aren't met.
type IdentityKnowledgeRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IdentityKnowledgeRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IdentityKnowledgeRequestMultiError) AllErrors() []error { return m }

// IdentityKnowledgeRequestValidationError is the validation error returned by
// IdentityKnowledgeRequest.Validate if the designated constraints aren't met.
type IdentityKnowledgeRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IdentityKnowledgeRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IdentityKnowledgeRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IdentityKnowledgeRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IdentityKnowledgeRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IdentityKnowledgeRequestValidationError) ErrorName() string {
	return "IdentityKnowledgeRequestValidationError"
}

// Error satisfies the builtin error interface
func (e IdentityKnowledgeRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIdentityKnowledgeRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IdentityKnowledgeRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IdentityKnowledgeRequestValidationError{}

var _IdentityKnowledgeRequest_Operation_InLookup = map[Operation]struct{}{
	1: {},
}

// Validate checks the field values on IdentityKnowledgeResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IdentityKnowledgeResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IdentityKnowledgeResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IdentityKnowledgeResponseMultiError, or nil if none found.
func (m *IdentityKnowledgeResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *IdentityKnowledgeResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPaths() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, IdentityKnowledgeResponseValidationError{
						field:  fmt.Sprintf("Paths[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, IdentityKnowledgeResponseValidationError{
						field:  fmt.Sprintf("Paths[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return IdentityKnowledgeResponseValidationError{
					field:  fmt.Sprintf("Paths[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return IdentityKnowledgeResponseMultiError(errors)
	}

	return nil
}

// IdentityKnowledgeResponseMultiError is an error wrapping multiple validation
// errors returned by IdentityKnowledgeResponse.ValidateAll() if the
// designated constraints aren't met.
type IdentityKnowledgeResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IdentityKnowledgeResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IdentityKnowledgeResponseMultiError) AllErrors() []error { return m }

// IdentityKnowledgeResponseValidationError is the validation error returned by
// IdentityKnowledgeResponse.Validate if the designated constraints aren't met.
type IdentityKnowledgeResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IdentityKnowledgeResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IdentityKnowledgeResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IdentityKnowledgeResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IdentityKnowledgeResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IdentityKnowledgeResponseValidationError) ErrorName() string {
	return "IdentityKnowledgeResponseValidationError"
}

// Error satisfies the builtin error interface
func (e IdentityKnowledgeResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIdentityKnowledgeResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IdentityKnowledgeResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IdentityKnowledgeResponseValidationError{}