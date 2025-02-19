// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/auditsink/v1beta1/credentials.proto

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

// Validate checks the field values on TokenIntrospected with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *TokenIntrospected) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TokenIntrospected with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TokenIntrospectedMultiError, or nil if none found.
func (m *TokenIntrospected) ValidateAll() error {
	return m.validate(true)
}

func (m *TokenIntrospected) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ResolvedDigitalTwinId

	// no validation rules for ConfigId

	// no validation rules for Issuer

	// no validation rules for Subject

	// no validation rules for Jti

	if len(errors) > 0 {
		return TokenIntrospectedMultiError(errors)
	}

	return nil
}

// TokenIntrospectedMultiError is an error wrapping multiple validation errors
// returned by TokenIntrospected.ValidateAll() if the designated constraints
// aren't met.
type TokenIntrospectedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenIntrospectedMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenIntrospectedMultiError) AllErrors() []error { return m }

// TokenIntrospectedValidationError is the validation error returned by
// TokenIntrospected.Validate if the designated constraints aren't met.
type TokenIntrospectedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenIntrospectedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenIntrospectedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenIntrospectedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenIntrospectedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenIntrospectedValidationError) ErrorName() string {
	return "TokenIntrospectedValidationError"
}

// Error satisfies the builtin error interface
func (e TokenIntrospectedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTokenIntrospected.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenIntrospectedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenIntrospectedValidationError{}
