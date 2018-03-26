// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/provider (interfaces: AdminProvider)

// Package mock_provider is a generated GoMock package.
package mock_provider

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	reflect "reflect"
	time "time"
)

// MockAdminProvider is a mock of AdminProvider interface
type MockAdminProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAdminProviderMockRecorder
}

// MockAdminProviderMockRecorder is the mock recorder for MockAdminProvider
type MockAdminProviderMockRecorder struct {
	mock *MockAdminProvider
}

// NewMockAdminProvider creates a new mock instance
func NewMockAdminProvider(ctrl *gomock.Controller) *MockAdminProvider {
	mock := &MockAdminProvider{ctrl: ctrl}
	mock.recorder = &MockAdminProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAdminProvider) EXPECT() *MockAdminProviderMockRecorder {
	return m.recorder
}

// Logs mocks base method
func (m *MockAdminProvider) Logs(arg0 int, arg1, arg2 time.Time) ([]models.LogFile, error) {
	ret := m.ctrl.Call(m, "Logs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logs indicates an expected call of Logs
func (mr *MockAdminProviderMockRecorder) Logs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logs", reflect.TypeOf((*MockAdminProvider)(nil).Logs), arg0, arg1, arg2)
}
