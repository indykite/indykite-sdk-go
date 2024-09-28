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
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1 (interfaces: EntityMatchingAPIClient)
//
// Generated by this command:
//
//	mockgen -copyright_file ./doc/LICENSE -package entitymatching -destination ./test/entitymatching/v1beta1/entitymatching_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1 EntityMatchingAPIClient
//

// Package entitymatching is a generated GoMock package.
package entitymatching

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"

	entitymatchingv1beta1 "github.com/indykite/indykite-sdk-go/gen/indykite/entitymatching/v1beta1"
)

// MockEntityMatchingAPIClient is a mock of EntityMatchingAPIClient interface.
type MockEntityMatchingAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockEntityMatchingAPIClientMockRecorder
}

// MockEntityMatchingAPIClientMockRecorder is the mock recorder for MockEntityMatchingAPIClient.
type MockEntityMatchingAPIClientMockRecorder struct {
	mock *MockEntityMatchingAPIClient
}

// NewMockEntityMatchingAPIClient creates a new mock instance.
func NewMockEntityMatchingAPIClient(ctrl *gomock.Controller) *MockEntityMatchingAPIClient {
	mock := &MockEntityMatchingAPIClient{ctrl: ctrl}
	mock.recorder = &MockEntityMatchingAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntityMatchingAPIClient) EXPECT() *MockEntityMatchingAPIClientMockRecorder {
	return m.recorder
}

// ReadEntityMatchingReport mocks base method.
func (m *MockEntityMatchingAPIClient) ReadEntityMatchingReport(arg0 context.Context, arg1 *entitymatchingv1beta1.ReadEntityMatchingReportRequest, arg2 ...grpc.CallOption) (*entitymatchingv1beta1.ReadEntityMatchingReportResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadEntityMatchingReport", varargs...)
	ret0, _ := ret[0].(*entitymatchingv1beta1.ReadEntityMatchingReportResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadEntityMatchingReport indicates an expected call of ReadEntityMatchingReport.
func (mr *MockEntityMatchingAPIClientMockRecorder) ReadEntityMatchingReport(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEntityMatchingReport", reflect.TypeOf((*MockEntityMatchingAPIClient)(nil).ReadEntityMatchingReport), varargs...)
}

// ReadSuggestedPropertyMapping mocks base method.
func (m *MockEntityMatchingAPIClient) ReadSuggestedPropertyMapping(arg0 context.Context, arg1 *entitymatchingv1beta1.ReadSuggestedPropertyMappingRequest, arg2 ...grpc.CallOption) (*entitymatchingv1beta1.ReadSuggestedPropertyMappingResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReadSuggestedPropertyMapping", varargs...)
	ret0, _ := ret[0].(*entitymatchingv1beta1.ReadSuggestedPropertyMappingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSuggestedPropertyMapping indicates an expected call of ReadSuggestedPropertyMapping.
func (mr *MockEntityMatchingAPIClientMockRecorder) ReadSuggestedPropertyMapping(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSuggestedPropertyMapping", reflect.TypeOf((*MockEntityMatchingAPIClient)(nil).ReadSuggestedPropertyMapping), varargs...)
}

// RunEntityMatchingPipeline mocks base method.
func (m *MockEntityMatchingAPIClient) RunEntityMatchingPipeline(arg0 context.Context, arg1 *entitymatchingv1beta1.RunEntityMatchingPipelineRequest, arg2 ...grpc.CallOption) (*entitymatchingv1beta1.RunEntityMatchingPipelineResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunEntityMatchingPipeline", varargs...)
	ret0, _ := ret[0].(*entitymatchingv1beta1.RunEntityMatchingPipelineResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunEntityMatchingPipeline indicates an expected call of RunEntityMatchingPipeline.
func (mr *MockEntityMatchingAPIClientMockRecorder) RunEntityMatchingPipeline(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunEntityMatchingPipeline", reflect.TypeOf((*MockEntityMatchingAPIClient)(nil).RunEntityMatchingPipeline), varargs...)
}