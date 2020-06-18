package testutils

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Recorder interface
type Recorder struct {
	ctrl     *gomock.Controller
	recorder *_RecorderRecorder
}

// Recorder for Recorder (not exported)
type _RecorderRecorder struct {
	mock *Recorder
}

func NewRecorder(ctrl *gomock.Controller) *Recorder {
	mock := &Recorder{ctrl: ctrl}
	mock.recorder = &_RecorderRecorder{mock}
	return mock
}

func (_m *Recorder) EXPECT() *_RecorderRecorder {
	return _m.recorder
}

func (_m *Recorder) Call(_param0 string) {
	_m.ctrl.Call(_m, "Call", _param0)
	return
}

func (_mr *_RecorderRecorder) Call(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Call", arg0)
}
