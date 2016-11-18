// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.imshealth.com/xfra/layer0/api/logic (interfaces: TaskLogic)

package mock_logic

import (
	gomock "github.com/golang/mock/gomock"
	models "gitlab.imshealth.com/xfra/layer0/common/models"
)

// Mock of TaskLogic interface
type MockTaskLogic struct {
	ctrl     *gomock.Controller
	recorder *_MockTaskLogicRecorder
}

// Recorder for MockTaskLogic (not exported)
type _MockTaskLogicRecorder struct {
	mock *MockTaskLogic
}

func NewMockTaskLogic(ctrl *gomock.Controller) *MockTaskLogic {
	mock := &MockTaskLogic{ctrl: ctrl}
	mock.recorder = &_MockTaskLogicRecorder{mock}
	return mock
}

func (_m *MockTaskLogic) EXPECT() *_MockTaskLogicRecorder {
	return _m.recorder
}

func (_m *MockTaskLogic) CreateTask(_param0 models.CreateTaskRequest) (*models.Task, error) {
	ret := _m.ctrl.Call(_m, "CreateTask", _param0)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTaskLogicRecorder) CreateTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateTask", arg0)
}

func (_m *MockTaskLogic) DeleteTask(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteTask", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockTaskLogicRecorder) DeleteTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteTask", arg0)
}

func (_m *MockTaskLogic) GetTask(_param0 string) (*models.Task, error) {
	ret := _m.ctrl.Call(_m, "GetTask", _param0)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTaskLogicRecorder) GetTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTask", arg0)
}

func (_m *MockTaskLogic) GetTaskLogs(_param0 string, _param1 int) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "GetTaskLogs", _param0, _param1)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTaskLogicRecorder) GetTaskLogs(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTaskLogs", arg0, arg1)
}

func (_m *MockTaskLogic) ListTasks() ([]*models.Task, error) {
	ret := _m.ctrl.Call(_m, "ListTasks")
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTaskLogicRecorder) ListTasks() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListTasks")
}
