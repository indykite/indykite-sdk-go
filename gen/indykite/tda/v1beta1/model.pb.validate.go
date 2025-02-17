// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/tda/v1beta1/model.proto

package tdav1beta1

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

// Validate checks the field values on Consent with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Consent) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Consent with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ConsentMultiError, or nil if none found.
func (m *Consent) ValidateAll() error {
	return m.validate(true)
}

func (m *Consent) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return ConsentMultiError(errors)
	}

	return nil
}

// ConsentMultiError is an error wrapping multiple validation errors returned
// by Consent.ValidateAll() if the designated constraints aren't met.
type ConsentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ConsentMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ConsentMultiError) AllErrors() []error { return m }

// ConsentValidationError is the validation error returned by Consent.Validate
// if the designated constraints aren't met.
type ConsentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConsentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConsentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConsentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConsentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConsentValidationError) ErrorName() string { return "ConsentValidationError" }

// Error satisfies the builtin error interface
func (e ConsentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConsent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConsentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConsentValidationError{}

// Validate checks the field values on TrustedDataNode with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *TrustedDataNode) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TrustedDataNode with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TrustedDataNodeMultiError, or nil if none found.
func (m *TrustedDataNode) ValidateAll() error {
	return m.validate(true)
}

func (m *TrustedDataNode) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() != "" {

		if l := utf8.RuneCountInString(m.GetId()); l < 22 || l > 256 {
			err := TrustedDataNodeValidationError{
				field:  "Id",
				reason: "value length must be between 22 and 256 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if !strings.HasPrefix(m.GetId(), "gid:") {
			err := TrustedDataNodeValidationError{
				field:  "Id",
				reason: "value does not have prefix \"gid:\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if l := utf8.RuneCountInString(m.GetExternalId()); l < 1 || l > 256 {
		err := TrustedDataNodeValidationError{
			field:  "ExternalId",
			reason: "value length must be between 1 and 256 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetType()) > 64 {
		err := TrustedDataNodeValidationError{
			field:  "Type",
			reason: "value length must be at most 64 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_TrustedDataNode_Type_Pattern.MatchString(m.GetType()) {
		err := TrustedDataNodeValidationError{
			field:  "Type",
			reason: "value does not match regex pattern \"^([A-Z][a-z]+)+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetTags()) > 32 {
		err := TrustedDataNodeValidationError{
			field:  "Tags",
			reason: "value must contain no more than 32 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	_TrustedDataNode_Tags_Unique := make(map[string]struct{}, len(m.GetTags()))

	for idx, item := range m.GetTags() {
		_, _ = idx, item

		if _, exists := _TrustedDataNode_Tags_Unique[item]; exists {
			err := TrustedDataNodeValidationError{
				field:  fmt.Sprintf("Tags[%v]", idx),
				reason: "repeated value must contain unique items",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {
			_TrustedDataNode_Tags_Unique[item] = struct{}{}
		}

		if utf8.RuneCountInString(item) > 64 {
			err := TrustedDataNodeValidationError{
				field:  fmt.Sprintf("Tags[%v]", idx),
				reason: "value length must be at most 64 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if !_TrustedDataNode_Tags_Pattern.MatchString(item) {
			err := TrustedDataNodeValidationError{
				field:  fmt.Sprintf("Tags[%v]", idx),
				reason: "value does not match regex pattern \"^([A-Z][a-z]+)+$\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if all {
		switch v := interface{}(m.GetCreateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TrustedDataNodeValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TrustedDataNodeValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TrustedDataNodeValidationError{
				field:  "CreateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TrustedDataNodeValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TrustedDataNodeValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TrustedDataNodeValidationError{
				field:  "UpdateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(m.GetProperties()) > 50 {
		err := TrustedDataNodeValidationError{
			field:  "Properties",
			reason: "value must contain no more than 50 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetProperties() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TrustedDataNodeValidationError{
						field:  fmt.Sprintf("Properties[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TrustedDataNodeValidationError{
						field:  fmt.Sprintf("Properties[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TrustedDataNodeValidationError{
					field:  fmt.Sprintf("Properties[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for IsIdentity

	for idx, item := range m.GetNodes() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TrustedDataNodeValidationError{
						field:  fmt.Sprintf("Nodes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TrustedDataNodeValidationError{
						field:  fmt.Sprintf("Nodes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TrustedDataNodeValidationError{
					field:  fmt.Sprintf("Nodes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return TrustedDataNodeMultiError(errors)
	}

	return nil
}

// TrustedDataNodeMultiError is an error wrapping multiple validation errors
// returned by TrustedDataNode.ValidateAll() if the designated constraints
// aren't met.
type TrustedDataNodeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TrustedDataNodeMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TrustedDataNodeMultiError) AllErrors() []error { return m }

// TrustedDataNodeValidationError is the validation error returned by
// TrustedDataNode.Validate if the designated constraints aren't met.
type TrustedDataNodeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TrustedDataNodeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TrustedDataNodeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TrustedDataNodeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TrustedDataNodeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TrustedDataNodeValidationError) ErrorName() string { return "TrustedDataNodeValidationError" }

// Error satisfies the builtin error interface
func (e TrustedDataNodeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTrustedDataNode.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TrustedDataNodeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TrustedDataNodeValidationError{}

var _TrustedDataNode_Type_Pattern = regexp.MustCompile("^([A-Z][a-z]+)+$")

var _TrustedDataNode_Tags_Pattern = regexp.MustCompile("^([A-Z][a-z]+)+$")
