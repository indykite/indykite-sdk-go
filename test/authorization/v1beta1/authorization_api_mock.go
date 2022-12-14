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
// Source: github.com/indykite/jarvis-sdk-go/gen/indykite/authorization/v1beta1 (interfaces: AuthorizationAPIClient)

// Package authorization is a generated GoMock package.
package authorization

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"

	authorizationv1beta1 "github.com/indykite/jarvis-sdk-go/gen/indykite/authorization/v1beta1"
)

// MockAuthorizationAPIClient is a mock of AuthorizationAPIClient interface.
type MockAuthorizationAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationAPIClientMockRecorder
}

// MockAuthorizationAPIClientMockRecorder is the mock recorder for MockAuthorizationAPIClient.
type MockAuthorizationAPIClientMockRecorder struct {
	mock *MockAuthorizationAPIClient
}

// NewMockAuthorizationAPIClient creates a new mock instance.
func NewMockAuthorizationAPIClient(ctrl *gomock.Controller) *MockAuthorizationAPIClient {
	mock := &MockAuthorizationAPIClient{ctrl: ctrl}
	mock.recorder = &MockAuthorizationAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationAPIClient) EXPECT() *MockAuthorizationAPIClientMockRecorder {
	return m.recorder
}

// IsAuthorized mocks base method.
func (m *MockAuthorizationAPIClient) IsAuthorized(arg0 context.Context, arg1 *authorizationv1beta1.IsAuthorizedRequest, arg2 ...grpc.CallOption) (*authorizationv1beta1.IsAuthorizedResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsAuthorized", varargs...)
	ret0, _ := ret[0].(*authorizationv1beta1.IsAuthorizedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAuthorized indicates an expected call of IsAuthorized.
func (mr *MockAuthorizationAPIClientMockRecorder) IsAuthorized(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAuthorized", reflect.TypeOf((*MockAuthorizationAPIClient)(nil).IsAuthorized), varargs...)
}