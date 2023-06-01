// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: indykite/ingest/v1beta2/ingest_api.proto

package ingestv1beta2

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

// Validate checks the field values on StreamRecordsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StreamRecordsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StreamRecordsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StreamRecordsRequestMultiError, or nil if none found.
func (m *StreamRecordsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StreamRecordsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetRecord() == nil {
		err := StreamRecordsRequestValidationError{
			field:  "Record",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetRecord()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, StreamRecordsRequestValidationError{
					field:  "Record",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, StreamRecordsRequestValidationError{
					field:  "Record",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRecord()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return StreamRecordsRequestValidationError{
				field:  "Record",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return StreamRecordsRequestMultiError(errors)
	}

	return nil
}

// StreamRecordsRequestMultiError is an error wrapping multiple validation
// errors returned by StreamRecordsRequest.ValidateAll() if the designated
// constraints aren't met.
type StreamRecordsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StreamRecordsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StreamRecordsRequestMultiError) AllErrors() []error { return m }

// StreamRecordsRequestValidationError is the validation error returned by
// StreamRecordsRequest.Validate if the designated constraints aren't met.
type StreamRecordsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StreamRecordsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StreamRecordsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StreamRecordsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StreamRecordsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StreamRecordsRequestValidationError) ErrorName() string {
	return "StreamRecordsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e StreamRecordsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStreamRecordsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StreamRecordsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StreamRecordsRequestValidationError{}

// Validate checks the field values on StreamRecordsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StreamRecordsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StreamRecordsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StreamRecordsResponseMultiError, or nil if none found.
func (m *StreamRecordsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StreamRecordsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RecordId

	// no validation rules for RecordIndex

	if all {
		switch v := interface{}(m.GetInfo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, StreamRecordsResponseValidationError{
					field:  "Info",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, StreamRecordsResponseValidationError{
					field:  "Info",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return StreamRecordsResponseValidationError{
				field:  "Info",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.Error.(type) {

	case *StreamRecordsResponse_RecordError:

		if all {
			switch v := interface{}(m.GetRecordError()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamRecordsResponseValidationError{
						field:  "RecordError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamRecordsResponseValidationError{
						field:  "RecordError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetRecordError()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamRecordsResponseValidationError{
					field:  "RecordError",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *StreamRecordsResponse_StatusError:

		if all {
			switch v := interface{}(m.GetStatusError()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, StreamRecordsResponseValidationError{
						field:  "StatusError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, StreamRecordsResponseValidationError{
						field:  "StatusError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetStatusError()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return StreamRecordsResponseValidationError{
					field:  "StatusError",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return StreamRecordsResponseMultiError(errors)
	}

	return nil
}

// StreamRecordsResponseMultiError is an error wrapping multiple validation
// errors returned by StreamRecordsResponse.ValidateAll() if the designated
// constraints aren't met.
type StreamRecordsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StreamRecordsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StreamRecordsResponseMultiError) AllErrors() []error { return m }

// StreamRecordsResponseValidationError is the validation error returned by
// StreamRecordsResponse.Validate if the designated constraints aren't met.
type StreamRecordsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StreamRecordsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StreamRecordsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StreamRecordsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StreamRecordsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StreamRecordsResponseValidationError) ErrorName() string {
	return "StreamRecordsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StreamRecordsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStreamRecordsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StreamRecordsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StreamRecordsResponseValidationError{}

// Validate checks the field values on IngestRecordRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IngestRecordRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IngestRecordRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IngestRecordRequestMultiError, or nil if none found.
func (m *IngestRecordRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *IngestRecordRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetRecord() == nil {
		err := IngestRecordRequestValidationError{
			field:  "Record",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetRecord()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IngestRecordRequestValidationError{
					field:  "Record",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IngestRecordRequestValidationError{
					field:  "Record",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRecord()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IngestRecordRequestValidationError{
				field:  "Record",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return IngestRecordRequestMultiError(errors)
	}

	return nil
}

// IngestRecordRequestMultiError is an error wrapping multiple validation
// errors returned by IngestRecordRequest.ValidateAll() if the designated
// constraints aren't met.
type IngestRecordRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IngestRecordRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IngestRecordRequestMultiError) AllErrors() []error { return m }

// IngestRecordRequestValidationError is the validation error returned by
// IngestRecordRequest.Validate if the designated constraints aren't met.
type IngestRecordRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IngestRecordRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IngestRecordRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IngestRecordRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IngestRecordRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IngestRecordRequestValidationError) ErrorName() string {
	return "IngestRecordRequestValidationError"
}

// Error satisfies the builtin error interface
func (e IngestRecordRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIngestRecordRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IngestRecordRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IngestRecordRequestValidationError{}

// Validate checks the field values on IngestRecordResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IngestRecordResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IngestRecordResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IngestRecordResponseMultiError, or nil if none found.
func (m *IngestRecordResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *IngestRecordResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RecordId

	if all {
		switch v := interface{}(m.GetInfo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IngestRecordResponseValidationError{
					field:  "Info",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IngestRecordResponseValidationError{
					field:  "Info",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IngestRecordResponseValidationError{
				field:  "Info",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.Error.(type) {

	case *IngestRecordResponse_RecordError:

		if all {
			switch v := interface{}(m.GetRecordError()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, IngestRecordResponseValidationError{
						field:  "RecordError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, IngestRecordResponseValidationError{
						field:  "RecordError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetRecordError()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return IngestRecordResponseValidationError{
					field:  "RecordError",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *IngestRecordResponse_StatusError:

		if all {
			switch v := interface{}(m.GetStatusError()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, IngestRecordResponseValidationError{
						field:  "StatusError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, IngestRecordResponseValidationError{
						field:  "StatusError",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetStatusError()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return IngestRecordResponseValidationError{
					field:  "StatusError",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return IngestRecordResponseMultiError(errors)
	}

	return nil
}

// IngestRecordResponseMultiError is an error wrapping multiple validation
// errors returned by IngestRecordResponse.ValidateAll() if the designated
// constraints aren't met.
type IngestRecordResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IngestRecordResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IngestRecordResponseMultiError) AllErrors() []error { return m }

// IngestRecordResponseValidationError is the validation error returned by
// IngestRecordResponse.Validate if the designated constraints aren't met.
type IngestRecordResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IngestRecordResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IngestRecordResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IngestRecordResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IngestRecordResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IngestRecordResponseValidationError) ErrorName() string {
	return "IngestRecordResponseValidationError"
}

// Error satisfies the builtin error interface
func (e IngestRecordResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIngestRecordResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IngestRecordResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IngestRecordResponseValidationError{}
