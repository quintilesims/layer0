// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/api/scheduler (interfaces: EnvironmentScaler)

package mock_scheduler

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	time "time"
)

// Mock of EnvironmentScaler interface
type MockEnvironmentScaler struct {
	ctrl     *gomock.Controller
	recorder *_MockEnvironmentScalerRecorder
}

// Recorder for MockEnvironmentScaler (not exported)
type _MockEnvironmentScalerRecorder struct {
	mock *MockEnvironmentScaler
}

func NewMockEnvironmentScaler(ctrl *gomock.Controller) *MockEnvironmentScaler {
	mock := &MockEnvironmentScaler{ctrl: ctrl}
	mock.recorder = &_MockEnvironmentScalerRecorder{mock}
	return mock
}

func (_m *MockEnvironmentScaler) EXPECT() *_MockEnvironmentScalerRecorder {
	return _m.recorder
}

func (_m *MockEnvironmentScaler) Scale(_param0 string) (*models.ScalerRunInfo, error) {
	ret := _m.ctrl.Call(_m, "Scale", _param0)
	ret0, _ := ret[0].(*models.ScalerRunInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvironmentScalerRecorder) Scale(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Scale", arg0)
}

func (_m *MockEnvironmentScaler) ScheduleRun(_param0 string, _param1 time.Duration) {
	_m.ctrl.Call(_m, "ScheduleRun", _param0, _param1)
}

func (_mr *_MockEnvironmentScalerRecorder) ScheduleRun(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ScheduleRun", arg0, arg1)
}
