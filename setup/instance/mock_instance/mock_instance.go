// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/setup/instance (interfaces: Instance)

// Package mock_instance is a generated GoMock package.
package mock_instance

import (
	gomock "github.com/golang/mock/gomock"
	s3iface "github.com/aws/aws-sdk-go/service/s3/s3iface"
	reflect "reflect"
)

// MockInstance is a mock of Instance interface
type MockInstance struct {
	ctrl     *gomock.Controller
	recorder *MockInstanceMockRecorder
}

// MockInstanceMockRecorder is the mock recorder for MockInstance
type MockInstanceMockRecorder struct {
	mock *MockInstance
}

// NewMockInstance creates a new mock instance
func NewMockInstance(ctrl *gomock.Controller) *MockInstance {
	mock := &MockInstance{ctrl: ctrl}
	mock.recorder = &MockInstanceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInstance) EXPECT() *MockInstanceMockRecorder {
	return m.recorder
}

// Apply mocks base method
func (m *MockInstance) Apply(arg0 bool) error {
	ret := m.ctrl.Call(m, "Apply", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Apply indicates an expected call of Apply
func (mr *MockInstanceMockRecorder) Apply(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockInstance)(nil).Apply), arg0)
}

// Destroy mocks base method
func (m *MockInstance) Destroy(arg0 bool) error {
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockInstanceMockRecorder) Destroy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockInstance)(nil).Destroy), arg0)
}

// Init mocks base method
func (m *MockInstance) Init(arg0 string, arg1 map[string]interface{}) error {
	ret := m.ctrl.Call(m, "Init", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockInstanceMockRecorder) Init(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockInstance)(nil).Init), arg0, arg1)
}

// Output mocks base method
func (m *MockInstance) Output(arg0 string) (string, error) {
	ret := m.ctrl.Call(m, "Output", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Output indicates an expected call of Output
func (mr *MockInstanceMockRecorder) Output(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockInstance)(nil).Output), arg0)
}

// Plan mocks base method
func (m *MockInstance) Plan() error {
	ret := m.ctrl.Call(m, "Plan")
	ret0, _ := ret[0].(error)
	return ret0
}

// Plan indicates an expected call of Plan
func (mr *MockInstanceMockRecorder) Plan() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Plan", reflect.TypeOf((*MockInstance)(nil).Plan))
}

// Pull mocks base method
func (m *MockInstance) Pull(arg0 s3iface.S3API) error {
	ret := m.ctrl.Call(m, "Pull", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pull indicates an expected call of Pull
func (mr *MockInstanceMockRecorder) Pull(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pull", reflect.TypeOf((*MockInstance)(nil).Pull), arg0)
}

// Push mocks base method
func (m *MockInstance) Push(arg0 s3iface.S3API) error {
	ret := m.ctrl.Call(m, "Push", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push
func (mr *MockInstanceMockRecorder) Push(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockInstance)(nil).Push), arg0)
}

// Set mocks base method
func (m *MockInstance) Set(arg0 map[string]interface{}) error {
	ret := m.ctrl.Call(m, "Set", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockInstanceMockRecorder) Set(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockInstance)(nil).Set), arg0)
}

// Unset mocks base method
func (m *MockInstance) Unset(arg0 string) error {
	ret := m.ctrl.Call(m, "Unset", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unset indicates an expected call of Unset
func (mr *MockInstanceMockRecorder) Unset(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unset", reflect.TypeOf((*MockInstance)(nil).Unset), arg0)
}

// Upgrade mocks base method
func (m *MockInstance) Upgrade(arg0 string, arg1 bool) error {
	ret := m.ctrl.Call(m, "Upgrade", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upgrade indicates an expected call of Upgrade
func (mr *MockInstanceMockRecorder) Upgrade(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upgrade", reflect.TypeOf((*MockInstance)(nil).Upgrade), arg0, arg1)
}
