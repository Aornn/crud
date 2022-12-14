// Code generated by MockGen. DO NOT EDIT.
// Source: put.go

// Package put_test is a generated GoMock package.
package put_test

import (
	domain "crud/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	zap "go.uber.org/zap"
)

// MockiPutUser is a mock of iPutUser interface.
type MockiPutUser struct {
	ctrl     *gomock.Controller
	recorder *MockiPutUserMockRecorder
}

// MockiPutUserMockRecorder is the mock recorder for MockiPutUser.
type MockiPutUserMockRecorder struct {
	mock *MockiPutUser
}

// NewMockiPutUser creates a new mock instance.
func NewMockiPutUser(ctrl *gomock.Controller) *MockiPutUser {
	mock := &MockiPutUser{ctrl: ctrl}
	mock.recorder = &MockiPutUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockiPutUser) EXPECT() *MockiPutUserMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockiPutUser) Process(l *zap.Logger, id string, newdata domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", l, id, newdata)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockiPutUserMockRecorder) Process(l, id, newdata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockiPutUser)(nil).Process), l, id, newdata)
}
