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
// Source: github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1 (interfaces: TrustedDataAccessAPIClient)
//
// Generated by this command:
//
//	mockgen -copyright_file ./doc/LICENSE -package tda -destination ./test/tda/v1beta1/trusted_data_access_api_mock.go github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1 TrustedDataAccessAPIClient
//

// Package tda is a generated GoMock package.
package tda

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"

	tdav1beta1 "github.com/indykite/indykite-sdk-go/gen/indykite/tda/v1beta1"
)

// MockTrustedDataAccessAPIClient is a mock of TrustedDataAccessAPIClient interface.
type MockTrustedDataAccessAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockTrustedDataAccessAPIClientMockRecorder
	isgomock struct{}
}

// MockTrustedDataAccessAPIClientMockRecorder is the mock recorder for MockTrustedDataAccessAPIClient.
type MockTrustedDataAccessAPIClientMockRecorder struct {
	mock *MockTrustedDataAccessAPIClient
}

// NewMockTrustedDataAccessAPIClient creates a new mock instance.
func NewMockTrustedDataAccessAPIClient(ctrl *gomock.Controller) *MockTrustedDataAccessAPIClient {
	mock := &MockTrustedDataAccessAPIClient{ctrl: ctrl}
	mock.recorder = &MockTrustedDataAccessAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrustedDataAccessAPIClient) EXPECT() *MockTrustedDataAccessAPIClientMockRecorder {
	return m.recorder
}

// DataAccess mocks base method.
func (m *MockTrustedDataAccessAPIClient) DataAccess(ctx context.Context, in *tdav1beta1.DataAccessRequest, opts ...grpc.CallOption) (*tdav1beta1.DataAccessResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DataAccess", varargs...)
	ret0, _ := ret[0].(*tdav1beta1.DataAccessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DataAccess indicates an expected call of DataAccess.
func (mr *MockTrustedDataAccessAPIClientMockRecorder) DataAccess(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataAccess", reflect.TypeOf((*MockTrustedDataAccessAPIClient)(nil).DataAccess), varargs...)
}

// GrantConsent mocks base method.
func (m *MockTrustedDataAccessAPIClient) GrantConsent(ctx context.Context, in *tdav1beta1.GrantConsentRequest, opts ...grpc.CallOption) (*tdav1beta1.GrantConsentResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GrantConsent", varargs...)
	ret0, _ := ret[0].(*tdav1beta1.GrantConsentResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GrantConsent indicates an expected call of GrantConsent.
func (mr *MockTrustedDataAccessAPIClientMockRecorder) GrantConsent(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GrantConsent", reflect.TypeOf((*MockTrustedDataAccessAPIClient)(nil).GrantConsent), varargs...)
}

// ListConsents mocks base method.
func (m *MockTrustedDataAccessAPIClient) ListConsents(ctx context.Context, in *tdav1beta1.ListConsentsRequest, opts ...grpc.CallOption) (*tdav1beta1.ListConsentsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListConsents", varargs...)
	ret0, _ := ret[0].(*tdav1beta1.ListConsentsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListConsents indicates an expected call of ListConsents.
func (mr *MockTrustedDataAccessAPIClientMockRecorder) ListConsents(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListConsents", reflect.TypeOf((*MockTrustedDataAccessAPIClient)(nil).ListConsents), varargs...)
}

// RevokeConsent mocks base method.
func (m *MockTrustedDataAccessAPIClient) RevokeConsent(ctx context.Context, in *tdav1beta1.RevokeConsentRequest, opts ...grpc.CallOption) (*tdav1beta1.RevokeConsentResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RevokeConsent", varargs...)
	ret0, _ := ret[0].(*tdav1beta1.RevokeConsentResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevokeConsent indicates an expected call of RevokeConsent.
func (mr *MockTrustedDataAccessAPIClientMockRecorder) RevokeConsent(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeConsent", reflect.TypeOf((*MockTrustedDataAccessAPIClient)(nil).RevokeConsent), varargs...)
}
