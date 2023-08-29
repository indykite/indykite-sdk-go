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
// Source: github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta1 (interfaces: IdentityKnowledgeAPIClient)

// Package knowledge is a generated GoMock package.
package knowledge

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"

	knowledgev1beta1 "github.com/indykite/indykite-sdk-go/gen/indykite/knowledge/v1beta1"
)

// MockIdentityKnowledgeAPIClient is a mock of IdentityKnowledgeAPIClient interface.
type MockIdentityKnowledgeAPIClient struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityKnowledgeAPIClientMockRecorder
}

// MockIdentityKnowledgeAPIClientMockRecorder is the mock recorder for MockIdentityKnowledgeAPIClient.
type MockIdentityKnowledgeAPIClientMockRecorder struct {
	mock *MockIdentityKnowledgeAPIClient
}

// NewMockIdentityKnowledgeAPIClient creates a new mock instance.
func NewMockIdentityKnowledgeAPIClient(ctrl *gomock.Controller) *MockIdentityKnowledgeAPIClient {
	mock := &MockIdentityKnowledgeAPIClient{ctrl: ctrl}
	mock.recorder = &MockIdentityKnowledgeAPIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIdentityKnowledgeAPIClient) EXPECT() *MockIdentityKnowledgeAPIClientMockRecorder {
	return m.recorder
}

// IdentityKnowledge mocks base method.
func (m *MockIdentityKnowledgeAPIClient) IdentityKnowledge(arg0 context.Context, arg1 *knowledgev1beta1.IdentityKnowledgeRequest, arg2 ...grpc.CallOption) (*knowledgev1beta1.IdentityKnowledgeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IdentityKnowledge", varargs...)
	ret0, _ := ret[0].(*knowledgev1beta1.IdentityKnowledgeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IdentityKnowledge indicates an expected call of IdentityKnowledge.
func (mr *MockIdentityKnowledgeAPIClientMockRecorder) IdentityKnowledge(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdentityKnowledge", reflect.TypeOf((*MockIdentityKnowledgeAPIClient)(nil).IdentityKnowledge), varargs...)
}