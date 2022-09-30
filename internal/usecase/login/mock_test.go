// Code generated by MockGen. DO NOT EDIT.
// Source: login.go

// Package login_test is a generated GoMock package.
package login_test

import (
	domain "crud/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockiGetUserFromDatabase is a mock of iGetUserFromDatabase interface.
type MockiGetUserFromDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockiGetUserFromDatabaseMockRecorder
}

// MockiGetUserFromDatabaseMockRecorder is the mock recorder for MockiGetUserFromDatabase.
type MockiGetUserFromDatabaseMockRecorder struct {
	mock *MockiGetUserFromDatabase
}

// NewMockiGetUserFromDatabase creates a new mock instance.
func NewMockiGetUserFromDatabase(ctrl *gomock.Controller) *MockiGetUserFromDatabase {
	mock := &MockiGetUserFromDatabase{ctrl: ctrl}
	mock.recorder = &MockiGetUserFromDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockiGetUserFromDatabase) EXPECT() *MockiGetUserFromDatabaseMockRecorder {
	return m.recorder
}

// GetUser mocks base method.
func (m *MockiGetUserFromDatabase) GetUser(id string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockiGetUserFromDatabaseMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockiGetUserFromDatabase)(nil).GetUser), id)
}
