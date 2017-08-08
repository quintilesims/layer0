// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/api/entity (interfaces: Environment)

package mock_entity

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
)

// Mock of Environment interface
type MockEnvironment struct {
	ctrl     *gomock.Controller
	recorder *_MockEnvironmentRecorder
}

// Recorder for MockEnvironment (not exported)
type _MockEnvironmentRecorder struct {
	mock *MockEnvironment
}

func NewMockEnvironment(ctrl *gomock.Controller) *MockEnvironment {
	mock := &MockEnvironment{ctrl: ctrl}
	mock.recorder = &_MockEnvironmentRecorder{mock}
	return mock
}

func (_m *MockEnvironment) EXPECT() *_MockEnvironmentRecorder {
	return _m.recorder
}

func (_m *MockEnvironment) Create(_param0 models.CreateEnvironmentRequest) error {
	ret := _m.ctrl.Call(_m, "Create", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEnvironmentRecorder) Create(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0)
}

func (_m *MockEnvironment) Delete() error {
	ret := _m.ctrl.Call(_m, "Delete")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEnvironmentRecorder) Delete() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete")
}

func (_m *MockEnvironment) Model() (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "Model")
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvironmentRecorder) Model() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Model")
}

func (_m *MockEnvironment) Summary() (*models.EnvironmentSummary, error) {
	ret := _m.ctrl.Call(_m, "Summary")
	ret0, _ := ret[0].(*models.EnvironmentSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvironmentRecorder) Summary() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Summary")
}
