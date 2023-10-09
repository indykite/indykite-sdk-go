// Copyright (c) 2023 IndyKite
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

type MatchSDKStatusErrorMatcher struct {
	message  types.GomegaMatcher
	Expected codes.Code
}

// MatchErrorCode succeeds if actual is a non-nil error that has the code the passed in error.
//
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
func MatchStatusError(errorCode codes.Code, msg interface{}) types.GomegaMatcher {
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

func (matcher *MatchSDKStatusErrorMatcher) Match(actual interface{}) (success bool, err error) {
	sdkError, ok := actual.(*errors.StatusError)
	if !ok {
		return false, fmt.Errorf("expected a SDK *errors.StatusError, got %T", actual)
	}
	s := sdkError.Status()
	if matcher.message == nil {
		return s.Code() == matcher.Expected, nil
	}
	return matcher.message.Match(map[string]interface{}{
		"Code":    s.Code(),
		"Message": s.Message(),
		"Details": s.Details(),
	})
}
func (matcher *MatchSDKStatusErrorMatcher) FailureMessage(actual interface{}) (message string) {
	if matcher.message == nil {
		return format.Message(status.Code(actual.(error)), "to match error", matcher.Expected)
	}
	return matcher.message.FailureMessage(actual)
}
func (matcher *MatchSDKStatusErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	if matcher.message == nil {
		return format.Message(status.Code(actual.(error)), "not to match error", matcher.Expected)
	}
	return matcher.message.NegatedFailureMessage(actual)
}

type MatchProtoMessageMatcher struct {
	Expected proto.Message
}

// MatchProtoMessage succeeds if.
func MatchProtoMessage(message proto.Message) types.GomegaMatcher {
	return &MatchProtoMessageMatcher{
		Expected: message,
	}
}

func (matcher *MatchProtoMessageMatcher) Match(actual interface{}) (success bool, err error) {
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
func (matcher *MatchProtoMessageMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to equal", matcher.Expected)
}
func (matcher *MatchProtoMessageMatcher) NegatedFailureMessage(actual interface{}) (message string) {
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

type EqualProtoMatcher struct {
	Expected proto.Message
}

func (matcher *EqualProtoMatcher) GomegaString() string {
	op := protojson.MarshalOptions{AllowPartial: true, Indent: "  "}
	ex, _ := op.Marshal(matcher.Expected)
	return string(ex)
}

func (matcher *EqualProtoMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		//nolint:lll
		return false, fmt.Errorf("refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized")
	}

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

func (matcher *EqualProtoMatcher) FailureMessage(actual interface{}) (message string) {
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

func (matcher *EqualProtoMatcher) NegatedFailureMessage(actual interface{}) (message string) {
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

type EqualAnyProtoMatcher struct {
	Expected proto.Message
	actual   proto.Message
}

func (matcher *EqualAnyProtoMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		//nolint:lll
		return false, fmt.Errorf("refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized")
	}

	a, ok := actual.(*anypb.Any)
	if !ok {
		return false, fmt.Errorf("expected a proto.Message.  Got:\n%s", format.Object(actual, 1))
	}

	matcher.actual, err = a.UnmarshalNew()
	if err != nil {
		return false, err
	}
	return proto.Equal(matcher.actual, matcher.Expected), nil
}

func (matcher *EqualAnyProtoMatcher) FailureMessage(interface{}) (message string) {
	// This function does not display if the type is different!!
	op := protojson.MarshalOptions{AllowPartial: true}
	ac, _ := op.Marshal(matcher.actual)
	ex, _ := op.Marshal(matcher.Expected)
	return format.MessageWithDiff(string(ac), "to equal", string(ex))
}

func (matcher *EqualAnyProtoMatcher) NegatedFailureMessage(interface{}) (message string) {
	// This function does not display if the type is different!!
	op := protojson.MarshalOptions{AllowPartial: true}
	ac, _ := op.Marshal(matcher.actual)
	ex, _ := op.Marshal(matcher.Expected)
	return format.MessageWithDiff(string(ac), "not to equal", string(ex))
}

type MatchJSONMatcher struct {
	types.GomegaMatcher
}

// MatchJSON uses the expected value to build a value matcher.
func MatchJSON(expected interface{}) types.GomegaMatcher {
	var (
		m   interface{}
		err error
	)
	switch v := expected.(type) {
	case string:
		err = json.Unmarshal([]byte(v), &m)
	case fmt.Stringer:
		err = json.Unmarshal([]byte(v.String()), &m)
	case []byte:
		err = json.Unmarshal(v, &m)
	case json.RawMessage:
		err = json.Unmarshal(v, &m)
	case types.GomegaMatcher:
		return &MatchJSONMatcher{GomegaMatcher: v}
	default:
		err = fmt.Errorf("unsupported exptected type value %T", expected)
	}
	if err != nil {
		panic(err)
	}
	return &MatchJSONMatcher{GomegaMatcher: valueMatcher(m)}
}

func (matcher *MatchJSONMatcher) Match(actual interface{}) (success bool, err error) {
	var aval interface{}
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

func valueMatcher(v interface{}) types.GomegaMatcher {
	switch jv := v.(type) {
	case nil:
		return gomega.BeNil()
	case bool:
		if jv {
			return gomega.BeTrue()
		}
		return gomega.BeFalse()
	case map[string]interface{}:
		keys := map[interface{}]types.GomegaMatcher{}
		for k, mv := range jv {
			keys[k] = valueMatcher(mv)
		}
		return gstruct.MatchAllKeys(keys)
	case []interface{}:
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

type matcherWrapper struct {
	matcher types.GomegaMatcher
	// This is used to save variable between calls to Matches and String in case of error
	// to be able to print better messages on failure
	actual interface{}
}

func WrapMatcher(matcher types.GomegaMatcher) gomock.Matcher {
	return &matcherWrapper{matcher: matcher}
}

func (m *matcherWrapper) Matches(x interface{}) (ok bool) {
	m.actual = x
	var err error
	if ok, err = m.matcher.Match(x); err != nil {
		ok = false
	}
	return
}

func (m *matcherWrapper) String() string {
	return fmt.Sprintf("Wrapped Gomega fail message: %s", m.matcher.FailureMessage(m.actual))
}
