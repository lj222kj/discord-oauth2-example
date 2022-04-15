// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_discord is a generated GoMock package.
package discord

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockService) Auth(ctx context.Context, code string) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", ctx, code)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth.
func (mr *MockServiceMockRecorder) Auth(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockService)(nil).Auth), ctx, code)
}

// AuthCsrfUrl mocks base method.
func (m *MockService) AuthCsrfUrl() (string, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthCsrfUrl")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// AuthCsrfUrl indicates an expected call of AuthCsrfUrl.
func (mr *MockServiceMockRecorder) AuthCsrfUrl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthCsrfUrl", reflect.TypeOf((*MockService)(nil).AuthCsrfUrl))
}