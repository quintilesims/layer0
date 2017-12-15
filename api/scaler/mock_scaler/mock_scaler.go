// Code generated by MockGen. DO NOT EDIT.
// Source: scaler.go

// Package mock_scaler is a generated GoMock package.
package mock_scaler

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockScaler is a mock of Scaler interface
type MockScaler struct {
	ctrl     *gomock.Controller
	recorder *MockScalerMockRecorder
}

// MockScalerMockRecorder is the mock recorder for MockScaler
type MockScalerMockRecorder struct {
	mock *MockScaler
}

// NewMockScaler creates a new mock instance
func NewMockScaler(ctrl *gomock.Controller) *MockScaler {
	mock := &MockScaler{ctrl: ctrl}
	mock.recorder = &MockScalerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScaler) EXPECT() *MockScalerMockRecorder {
	return m.recorder
}

// Scale mocks base method
func (m *MockScaler) Scale(environmentID string) error {
	ret := m.ctrl.Call(m, "Scale", environmentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scale indicates an expected call of Scale
func (mr *MockScalerMockRecorder) Scale(environmentID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scale", reflect.TypeOf((*MockScaler)(nil).Scale), environmentID)
}