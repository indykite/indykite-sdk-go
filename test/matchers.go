// Copyright (c) 2025 IndyKite
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

package test

import (
	"encoding/json"
	"fmt"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/indykite/indykite-sdk-go/errors"
)

// BeBase64 match string if it is Base64 Standard encoded.
func BeBase64() types.GomegaMatcher {
	return gomega.MatchRegexp(`^(?:[A-Za-z\d+/]{4})*(?:[A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?$]`)
}

// MatchSDKStatusErrorMatcher expects code.
type MatchSDKStatusErrorMatcher struct {
	message  types.GomegaMatcher
	Expected codes.Code
}

// MatchErrorCode succeeds if actual is a non-nil error that has the code the passed in error.
// Error must be google.golang.org/grpc/status.Error
// It is an error for err to be nil or an object that does not implement the Error interface.
func MatchErrorCode(errorCode codes.Code) types.GomegaMatcher {
	return &MatchSDKStatusErrorMatcher{
		Expected: errorCode,
	}
}

// MatchStatusError succeeds if actual is a non-nil error that has the code the passed in error.
//
// Error must be google.golang.org/grpc/status.Error
// It is an error for err to be nil or an object that does not implement the Error interface.
func MatchStatusError(errorCode codes.Code, msg any) types.GomegaMatcher {
	switch m := msg.(type) {
	case string:
		return &MatchSDKStatusErrorMatcher{
			message: gstruct.MatchAllKeys(gstruct.Keys{
				"Code":    gomega.Equal(errorCode),
				"Message": gomega.Equal(m),
				"Details": gomega.BeEmpty(),
			}),
		}
	case types.GomegaMatcher:
		return &MatchSDKStatusErrorMatcher{
			message: gstruct.MatchAllKeys(gstruct.Keys{
				"Code":    gomega.Equal(errorCode),
				"Message": m,
				"Details": gomega.BeEmpty(),
			}),
		}
	case nil:
		return &MatchSDKStatusErrorMatcher{
			Expected: errorCode,
		}
	default:
		panic(fmt.Errorf("unsupported message matcher type %T", msg))
	}
}

// Match returns matcher.
func (matcher *MatchSDKStatusErrorMatcher) Match(actual any) (bool, error) {
	sdkError, ok := actual.(*errors.StatusError)
	if !ok {
		return false, fmt.Errorf("expected a SDK *errors.StatusError, got %T", actual)
	}
	s := sdkError.Status()
	if matcher.message == nil {
		return s.Code() == matcher.Expected, nil
	}
	return matcher.message.Match(map[string]any{
		"Code":    s.Code(),
		"Message": s.Message(),
		"Details": s.Details(),
	})
}

// FailureMessage returns matcher.
func (matcher *MatchSDKStatusErrorMatcher) FailureMessage(actual any) string {
	if matcher.message == nil {
		return format.Message(status.Code(actual.(error)), "to match error", matcher.Expected)
	}
	return matcher.message.FailureMessage(actual)
}

// NegatedFailureMessage returns matcher.
func (matcher *MatchSDKStatusErrorMatcher) NegatedFailureMessage(actual any) string {
	if matcher.message == nil {
		return format.Message(status.Code(actual.(error)), "not to match error", matcher.Expected)
	}
	return matcher.message.NegatedFailureMessage(actual)
}

// MatchProtoMessageMatcher expects proto.
type MatchProtoMessageMatcher struct {
	Expected proto.Message
}

// MatchProtoMessage succeeds if.
func MatchProtoMessage(message proto.Message) types.GomegaMatcher {
	return &MatchProtoMessageMatcher{
		Expected: message,
	}
}

// Match returns proto.
func (matcher *MatchProtoMessageMatcher) Match(actual any) (bool, error) {
	eq := &matchers.EqualMatcher{
		Expected: matcher.Expected,
	}
	if s, e := eq.Match(actual); e != nil {
		return s, e
	}

	e, ok := actual.(proto.Message)
	if !ok {
		return false, fmt.Errorf("expected a proto message, got %T", actual)
	}

	return proto.Equal(e, matcher.Expected), nil
}

// FailureMessage returns format.
func (matcher *MatchProtoMessageMatcher) FailureMessage(actual any) string {
	return format.Message(actual, "to equal", matcher.Expected)
}

// NegatedFailureMessage returns format.
func (matcher *MatchProtoMessageMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to equal", matcher.Expected)
}

// EqualProto uses proto.Equal to compare actual with expected.  Equal is strict about
// types when performing comparisons.
// It is an error for both actual and expected to be nil.  Use BeNil() instead.
func EqualProto(expected protoreflect.ProtoMessage) types.GomegaMatcher {
	return &EqualProtoMatcher{
		Expected: expected,
	}
}

// EqualProtoMatcher expects proto.
type EqualProtoMatcher struct {
	Expected proto.Message
}

// GomegaString returns string.
func (matcher *EqualProtoMatcher) GomegaString() string {
	op := protojson.MarshalOptions{AllowPartial: true, Indent: "  "}
	ex, _ := op.Marshal(matcher.Expected)
	return string(ex)
}

// Match returns proto.
func (matcher *EqualProtoMatcher) Match(actual any) (bool, error) {
	if actual == nil && matcher.Expected == nil {
		//nolint:lll
		return false, fmt.Errorf("refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized")
	}
	var err error
	if a, ok := actual.(*anypb.Any); ok {
		var pa proto.Message
		pa, err = a.UnmarshalNew()
		if err != nil {
			return false, err
		}
		return proto.Equal(pa, matcher.Expected), nil
	}

	pa, ok := actual.(proto.Message)
	if !ok {
		return false, fmt.Errorf("expected a proto.Message.  Got:\n%s", format.Object(actual, 1))
	}
	return proto.Equal(pa, matcher.Expected), nil
}

// FailureMessage returns format.
func (matcher *EqualProtoMatcher) FailureMessage(actual any) string {
	actualMessage, actualOK := actual.(proto.Message)
	if actualOK {
		op := protojson.MarshalOptions{AllowPartial: true}
		ac, _ := op.Marshal(actualMessage)
		ex, _ := op.Marshal(matcher.Expected)
		acFullName := actualMessage.ProtoReflect().Descriptor().FullName()
		extFullName := matcher.Expected.ProtoReflect().Descriptor().FullName()
		return format.MessageWithDiff(
			string(acFullName)+": "+string(ac),
			"to equal",
			string(extFullName)+": "+string(ex),
		)
	}

	return format.Message(actual, "to equal", matcher.Expected)
}

// NegatedFailureMessage returns format.
func (matcher *EqualProtoMatcher) NegatedFailureMessage(actual any) string {
	actualMessage, actualOK := actual.(proto.Message)
	if actualOK {
		op := protojson.MarshalOptions{AllowPartial: true}
		ac, _ := op.Marshal(actualMessage)
		ex, _ := op.Marshal(matcher.Expected)
		acFullName := actualMessage.ProtoReflect().Descriptor().FullName()
		extFullName := matcher.Expected.ProtoReflect().Descriptor().FullName()
		return format.MessageWithDiff(
			string(acFullName)+": "+string(ac),
			"not to equal",
			string(extFullName)+": "+string(ex),
		)
	}
	return format.Message(actual, "not to equal", matcher.Expected)
}

// EqualAnyProto uses proto.Equal to compare actual with expected.  Equal is strict about
// types when performing comparisons.
// It is an error for both actual and expected to be nil.  Use BeNil() instead.
func EqualAnyProto(expected protoreflect.ProtoMessage) types.GomegaMatcher {
	return &EqualAnyProtoMatcher{
		Expected: expected,
	}
}

// EqualAnyProtoMatcher expects proto.
type EqualAnyProtoMatcher struct {
	Expected proto.Message
	actual   proto.Message
}

// Match returns proto.
func (matcher *EqualAnyProtoMatcher) Match(actual any) (bool, error) {
	if actual == nil && matcher.Expected == nil {
		//nolint:lll
		return false, fmt.Errorf("refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized")
	}

	a, ok := actual.(*anypb.Any)
	if !ok {
		return false, fmt.Errorf("expected a proto.Message.  Got:\n%s", format.Object(actual, 1))
	}
	var err error
	matcher.actual, err = a.UnmarshalNew()
	if err != nil {
		return false, err
	}
	return proto.Equal(matcher.actual, matcher.Expected), nil
}

// FailureMessage returns format.
func (matcher *EqualAnyProtoMatcher) FailureMessage(any) string {
	// This function does not display if the type is different!!
	op := protojson.MarshalOptions{AllowPartial: true}
	ac, _ := op.Marshal(matcher.actual)
	ex, _ := op.Marshal(matcher.Expected)
	return format.MessageWithDiff(string(ac), "to equal", string(ex))
}

// NegatedFailureMessage returns format.
func (matcher *EqualAnyProtoMatcher) NegatedFailureMessage(any) string {
	// This function does not display if the type is different!!
	op := protojson.MarshalOptions{AllowPartial: true}
	ac, _ := op.Marshal(matcher.actual)
	ex, _ := op.Marshal(matcher.Expected)
	return format.MessageWithDiff(string(ac), "not to equal", string(ex))
}

// MatchJSONMatcher expects types.
type MatchJSONMatcher struct {
	types.GomegaMatcher
}

// Match returns matcher.
func (matcher *MatchJSONMatcher) Match(actual any) (bool, error) {
	var aval any
	var err error
	switch v := actual.(type) {
	case string:
		err = json.Unmarshal([]byte(v), &aval)
	case fmt.Stringer:
		err = json.Unmarshal([]byte(v.String()), &aval)
	case []byte:
		err = json.Unmarshal(v, &aval)
	case json.RawMessage:
		err = json.Unmarshal(v, &aval)
	default:
		aval = actual
	}
	if err != nil {
		return false, err
	}
	return matcher.GomegaMatcher.Match(aval)
}

//nolint:unused //will cause error
func valueMatcher(v any) types.GomegaMatcher {
	switch jv := v.(type) {
	case nil:
		return gomega.BeNil()
	case bool:
		if jv {
			return gomega.BeTrue()
		}
		return gomega.BeFalse()
	case map[string]any:
		keys := map[any]types.GomegaMatcher{}
		for k, mv := range jv {
			keys[k] = valueMatcher(mv)
		}
		return gstruct.MatchAllKeys(keys)
	case []any:
		m := &matchers.ConsistOfMatcher{}
		for _, sv := range jv {
			m.Elements = append(m.Elements, valueMatcher(sv))
		}
		return m
	case types.GomegaMatcher:
		return jv
	default:
		return gomega.BeEquivalentTo(jv)
	}
}

// matcherWrapper initialization.
type matcherWrapper struct {
	matcher types.GomegaMatcher
	// This is used to save variable between calls to Matches and String in case of error
	// to be able to print better messages on failure
	actual any
}

// WrapMatcher returns matcher.
func WrapMatcher(matcher types.GomegaMatcher) gomock.Matcher {
	return &matcherWrapper{matcher: matcher}
}

// Matches returns ok.
func (m *matcherWrapper) Matches(x any) bool {
	m.actual = x
	ok, err := m.matcher.Match(x)
	if err != nil {
		ok = false
	}
	return ok
}

// matcherWrapper returns format.
func (m *matcherWrapper) String() string {
	return "Wrapped Gomega fail message: " + m.matcher.FailureMessage(m.actual)
}
