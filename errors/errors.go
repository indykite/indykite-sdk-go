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

// Package errors provides a convenient way to deal with SDK errors.
package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ClientError represents any error regarding the used of SDK or validating the operation call.
type ClientError struct {
	cause error
	msg   string
	code  codes.Code
}

func (err *ClientError) Error() string {
	if err.cause == nil {
		return fmt.Sprintf("client error: code = %s desc = %s", err.Code(), err.Message())
	}
	return fmt.Sprintf("client error: code = %s desc = %s cause=%v", err.Code(), err.Message(), err.cause)
}

// Code returns codes.Code.
//
// Unknown is returned also if original error is not from GRPC.
func (err *ClientError) Code() codes.Code {
	return err.code
}

func (err *ClientError) Message() string {
	return err.msg
}

func (err *ClientError) Unwrap() error {
	return err.cause
}

func (err *ClientError) As(out any) bool {
	switch t := out.(type) {
	case *ClientError:
		*t = *err
		return true
	default:
		return false
	}
}

// NewInvalidArgumentError returns a new ClientError represents an InvalidArgument.
func NewInvalidArgumentError(msg string, args ...any) error {
	return New(codes.InvalidArgument, msg, args...)
}

// NewInvalidArgumentErrorWithCause returns a new ClientError represents an InvalidArgument.
func NewInvalidArgumentErrorWithCause(cause error, msg string, args ...any) error {
	return NewWithCause(codes.InvalidArgument, cause, msg, args...)
}

// New returns a new ClientError with code and formatted message.
func New(code codes.Code, msg string, args ...any) error {
	return NewWithCause(code, nil, msg, args...)
}

// NewWithCause returns a new ClientError with code and formatted message.
func NewWithCause(code codes.Code, err error, msg string, args ...any) error {
	if code == codes.OK {
		if err == nil {
			return nil
		}
	} else if code > codes.Unauthenticated {
		return &ClientError{
			code:  codes.Internal,
			msg:   "invalid use of SDK - code value is invalid",
			cause: err,
		}
	}

	if len(args) > 0 {
		return &ClientError{
			code:  codes.InvalidArgument,
			msg:   fmt.Sprintf(msg, args...),
			cause: err,
		}
	}
	return &ClientError{
		code:  code,
		msg:   msg,
		cause: err,
	}
}

// IsServiceError checks if the Error code represents a service call error.
func IsServiceError(err error) bool {
	//nolint:errorlint // TODO: Should be fixed, but also heavily tested
	if se, ok := err.(interface {
		Code() codes.Code
	}); ok {
		//nolint:exhaustive // No need to use default, because last statement will handle all
		switch se.Code() {
		case codes.Unknown,
			codes.DeadlineExceeded,
			codes.ResourceExhausted,
			codes.Aborted,
			codes.Unimplemented,
			codes.Internal,
			codes.Unavailable,
			codes.DataLoss:
			// codes.Unauthenticated:
			return true
		}
	}
	return false
}

// StatusError wraps status returned from gRPC call with optional prefix for easier access.
type StatusError struct {
	origin     error
	grpcStatus *status.Status
}

// FromError converts given error into StatusError if possible, otherwise false is returned.
// Returns (nil, false) also if given error is nil.
func FromError(err error) *StatusError {
	if err == nil {
		return nil
	}
	//nolint:errorlint // TODO: Should be fixed, but also heavily tested
	if s, ok := err.(*StatusError); ok {
		return s
	}
	//nolint:errorlint // TODO: Should be fixed, but also heavily tested
	if s, ok := err.(*ClientError); ok {
		return &StatusError{grpcStatus: status.New(s.code, s.msg), origin: s.cause}
	}
	//nolint:errorlint // TODO: Should be fixed, but also heavily tested
	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return &StatusError{grpcStatus: se.GRPCStatus()}
	}
	return &StatusError{grpcStatus: status.New(codes.Unknown, err.Error())}
}

// NewGRPCError unwrap original GRPC status and wraps into GRPCErrorWrapper for easier access.
func NewGRPCError(entryErr any) *StatusError {
	if entryErr == nil {
		return nil
	}
	grpcErr := &StatusError{}

	var ok bool
	if grpcErr.origin, ok = entryErr.(error); !ok {
		grpcErr.origin = fmt.Errorf("%v", entryErr)
	}

	if s, ok := status.FromError(grpcErr.origin); ok {
		grpcErr.grpcStatus = s
	}

	return grpcErr
}

func (err *StatusError) Error() string {
	if err.origin != nil {
		return fmt.Sprintf(
			"client error: code = %s desc = %s | caused by %s", err.Code(), err.Message(), err.origin.Error())
	}
	return fmt.Sprintf("client error: code = %s desc = %s", err.Code(), err.Message())
}

// Message returns GRPC Status message.
func (err *StatusError) Message() string {
	return err.grpcStatus.Message()
}

// Code returns GRPC Status code.
//
// Unknown is returned also if original error is not from GRPC.
func (err *StatusError) Code() codes.Code {
	if err.grpcStatus == nil {
		return codes.Unknown
	}
	return err.grpcStatus.Code()
}

// Status returns GRPC status.
func (err *StatusError) Status() *status.Status {
	return err.grpcStatus
}

// WithPrefix set prefix which will be printed in when Error is called.
func (err *StatusError) WithPrefix(_ string) *StatusError {
	return err
}

// Origin returns underlying error, if any.
func (err *StatusError) Origin() error {
	return err.origin
}
