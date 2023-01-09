// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/identity/v1beta2/document.proto

package identityv1beta2

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

// Validate checks the field values on Document with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Document) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Document with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DocumentMultiError, or nil
// if none found.
func (m *Document) ValidateAll() error {
	return m.validate(true)
}

func (m *Document) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

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
						errors = append(errors, DocumentValidationError{
							field:  fmt.Sprintf("Fields[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, DocumentValidationError{
							field:  fmt.Sprintf("Fields[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return DocumentValidationError{
						field:  fmt.Sprintf("Fields[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if all {
		switch v := interface{}(m.GetCreateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DocumentValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DocumentValidationError{
					field:  "CreateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DocumentValidationError{
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
				errors = append(errors, DocumentValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DocumentValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DocumentValidationError{
				field:  "UpdateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return DocumentMultiError(errors)
	}

	return nil
}

// DocumentMultiError is an error wrapping multiple validation errors returned
// by Document.ValidateAll() if the designated constraints aren't met.
type DocumentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocumentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocumentMultiError) AllErrors() []error { return m }

// DocumentValidationError is the validation error returned by
// Document.Validate if the designated constraints aren't met.
type DocumentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocumentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocumentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocumentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocumentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocumentValidationError) ErrorName() string { return "DocumentValidationError" }

// Error satisfies the builtin error interface
func (e DocumentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocument.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocumentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocumentValidationError{}

// Validate checks the field values on DocumentMask with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DocumentMask) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocumentMask with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DocumentMaskMultiError, or
// nil if none found.
func (m *DocumentMask) ValidateAll() error {
	return m.validate(true)
}

func (m *DocumentMask) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DocumentMaskMultiError(errors)
	}

	return nil
}

// DocumentMaskMultiError is an error wrapping multiple validation errors
// returned by DocumentMask.ValidateAll() if the designated constraints aren't met.
type DocumentMaskMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocumentMaskMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocumentMaskMultiError) AllErrors() []error { return m }

// DocumentMaskValidationError is the validation error returned by
// DocumentMask.Validate if the designated constraints aren't met.
type DocumentMaskValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocumentMaskValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocumentMaskValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocumentMaskValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocumentMaskValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocumentMaskValidationError) ErrorName() string { return "DocumentMaskValidationError" }

// Error satisfies the builtin error interface
func (e DocumentMaskValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocumentMask.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocumentMaskValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocumentMaskValidationError{}

// Validate checks the field values on Precondition with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Precondition) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Precondition with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PreconditionMultiError, or
// nil if none found.
func (m *Precondition) ValidateAll() error {
	return m.validate(true)
}

func (m *Precondition) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	switch m.ConditionType.(type) {

	case *Precondition_Exists:
		// no validation rules for Exists

	case *Precondition_UpdateTime:

		if all {
			switch v := interface{}(m.GetUpdateTime()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PreconditionValidationError{
						field:  "UpdateTime",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PreconditionValidationError{
						field:  "UpdateTime",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdateTime()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PreconditionValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PreconditionMultiError(errors)
	}

	return nil
}

// PreconditionMultiError is an error wrapping multiple validation errors
// returned by Precondition.ValidateAll() if the designated constraints aren't met.
type PreconditionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PreconditionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PreconditionMultiError) AllErrors() []error { return m }

// PreconditionValidationError is the validation error returned by
// Precondition.Validate if the designated constraints aren't met.
type PreconditionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PreconditionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PreconditionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PreconditionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PreconditionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PreconditionValidationError) ErrorName() string { return "PreconditionValidationError" }

// Error satisfies the builtin error interface
func (e PreconditionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPrecondition.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PreconditionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PreconditionValidationError{}

// Validate checks the field values on Write with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Write) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Write with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in WriteMultiError, or nil if none found.
func (m *Write) ValidateAll() error {
	return m.validate(true)
}

func (m *Write) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetUpdateMask()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WriteValidationError{
					field:  "UpdateMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WriteValidationError{
					field:  "UpdateMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdateMask()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WriteValidationError{
				field:  "UpdateMask",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetUpdateTransforms() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, WriteValidationError{
						field:  fmt.Sprintf("UpdateTransforms[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, WriteValidationError{
						field:  fmt.Sprintf("UpdateTransforms[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return WriteValidationError{
					field:  fmt.Sprintf("UpdateTransforms[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetCurrentDocument()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WriteValidationError{
					field:  "CurrentDocument",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WriteValidationError{
					field:  "CurrentDocument",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCurrentDocument()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WriteValidationError{
				field:  "CurrentDocument",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.Operation.(type) {

	case *Write_Update:

		if all {
			switch v := interface{}(m.GetUpdate()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, WriteValidationError{
						field:  "Update",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, WriteValidationError{
						field:  "Update",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdate()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return WriteValidationError{
					field:  "Update",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Write_Delete:
		// no validation rules for Delete

	case *Write_Transform:

		if all {
			switch v := interface{}(m.GetTransform()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, WriteValidationError{
						field:  "Transform",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, WriteValidationError{
						field:  "Transform",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetTransform()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return WriteValidationError{
					field:  "Transform",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return WriteMultiError(errors)
	}

	return nil
}

// WriteMultiError is an error wrapping multiple validation errors returned by
// Write.ValidateAll() if the designated constraints aren't met.
type WriteMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WriteMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WriteMultiError) AllErrors() []error { return m }

// WriteValidationError is the validation error returned by Write.Validate if
// the designated constraints aren't met.
type WriteValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WriteValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WriteValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WriteValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WriteValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WriteValidationError) ErrorName() string { return "WriteValidationError" }

// Error satisfies the builtin error interface
func (e WriteValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWrite.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WriteValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WriteValidationError{}

// Validate checks the field values on WriteResult with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WriteResult) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WriteResult with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WriteResultMultiError, or
// nil if none found.
func (m *WriteResult) ValidateAll() error {
	return m.validate(true)
}

func (m *WriteResult) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetUpdateTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, WriteResultValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, WriteResultValidationError{
					field:  "UpdateTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdateTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return WriteResultValidationError{
				field:  "UpdateTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return WriteResultMultiError(errors)
	}

	return nil
}

// WriteResultMultiError is an error wrapping multiple validation errors
// returned by WriteResult.ValidateAll() if the designated constraints aren't met.
type WriteResultMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WriteResultMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WriteResultMultiError) AllErrors() []error { return m }

// WriteResultValidationError is the validation error returned by
// WriteResult.Validate if the designated constraints aren't met.
type WriteResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WriteResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WriteResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WriteResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WriteResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WriteResultValidationError) ErrorName() string { return "WriteResultValidationError" }

// Error satisfies the builtin error interface
func (e WriteResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWriteResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WriteResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WriteResultValidationError{}

// Validate checks the field values on DocumentTransform with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DocumentTransform) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocumentTransform with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DocumentTransformMultiError, or nil if none found.
func (m *DocumentTransform) ValidateAll() error {
	return m.validate(true)
}

func (m *DocumentTransform) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Document

	for idx, item := range m.GetFieldTransforms() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DocumentTransformValidationError{
						field:  fmt.Sprintf("FieldTransforms[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DocumentTransformValidationError{
						field:  fmt.Sprintf("FieldTransforms[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DocumentTransformValidationError{
					field:  fmt.Sprintf("FieldTransforms[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return DocumentTransformMultiError(errors)
	}

	return nil
}

// DocumentTransformMultiError is an error wrapping multiple validation errors
// returned by DocumentTransform.ValidateAll() if the designated constraints
// aren't met.
type DocumentTransformMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocumentTransformMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocumentTransformMultiError) AllErrors() []error { return m }

// DocumentTransformValidationError is the validation error returned by
// DocumentTransform.Validate if the designated constraints aren't met.
type DocumentTransformValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocumentTransformValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocumentTransformValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocumentTransformValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocumentTransformValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocumentTransformValidationError) ErrorName() string {
	return "DocumentTransformValidationError"
}

// Error satisfies the builtin error interface
func (e DocumentTransformValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocumentTransform.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocumentTransformValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocumentTransformValidationError{}

// Validate checks the field values on DocumentTransform_FieldTransform with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *DocumentTransform_FieldTransform) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DocumentTransform_FieldTransform with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// DocumentTransform_FieldTransformMultiError, or nil if none found.
func (m *DocumentTransform_FieldTransform) ValidateAll() error {
	return m.validate(true)
}

func (m *DocumentTransform_FieldTransform) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for FieldPath

	switch m.TransformType.(type) {

	case *DocumentTransform_FieldTransform_AppendMissingElements:

		if all {
			switch v := interface{}(m.GetAppendMissingElements()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DocumentTransform_FieldTransformValidationError{
						field:  "AppendMissingElements",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DocumentTransform_FieldTransformValidationError{
						field:  "AppendMissingElements",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAppendMissingElements()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DocumentTransform_FieldTransformValidationError{
					field:  "AppendMissingElements",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *DocumentTransform_FieldTransform_RemoveAllFromArray:

		if all {
			switch v := interface{}(m.GetRemoveAllFromArray()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DocumentTransform_FieldTransformValidationError{
						field:  "RemoveAllFromArray",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DocumentTransform_FieldTransformValidationError{
						field:  "RemoveAllFromArray",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetRemoveAllFromArray()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DocumentTransform_FieldTransformValidationError{
					field:  "RemoveAllFromArray",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return DocumentTransform_FieldTransformMultiError(errors)
	}

	return nil
}

// DocumentTransform_FieldTransformMultiError is an error wrapping multiple
// validation errors returned by
// DocumentTransform_FieldTransform.ValidateAll() if the designated
// constraints aren't met.
type DocumentTransform_FieldTransformMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DocumentTransform_FieldTransformMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DocumentTransform_FieldTransformMultiError) AllErrors() []error { return m }

// DocumentTransform_FieldTransformValidationError is the validation error
// returned by DocumentTransform_FieldTransform.Validate if the designated
// constraints aren't met.
type DocumentTransform_FieldTransformValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DocumentTransform_FieldTransformValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DocumentTransform_FieldTransformValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DocumentTransform_FieldTransformValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DocumentTransform_FieldTransformValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DocumentTransform_FieldTransformValidationError) ErrorName() string {
	return "DocumentTransform_FieldTransformValidationError"
}

// Error satisfies the builtin error interface
func (e DocumentTransform_FieldTransformValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDocumentTransform_FieldTransform.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DocumentTransform_FieldTransformValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DocumentTransform_FieldTransformValidationError{}