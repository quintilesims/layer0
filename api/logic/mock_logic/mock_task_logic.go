// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/logic (interfaces: TaskLogic)

// Package mock_logic is a generated GoMock package.
package mock_logic

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	reflect "reflect"
)

// MockTaskLogic is a mock of TaskLogic interface
type MockTaskLogic struct {
	ctrl     *gomock.Controller
	recorder *MockTaskLogicMockRecorder
}

// MockTaskLogicMockRecorder is the mock recorder for MockTaskLogic
type MockTaskLogicMockRecorder struct {
	mock *MockTaskLogic
}

// NewMockTaskLogic creates a new mock instance
func NewMockTaskLogic(ctrl *gomock.Controller) *MockTaskLogic {
	mock := &MockTaskLogic{ctrl: ctrl}
	mock.recorder = &MockTaskLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskLogic) EXPECT() *MockTaskLogicMockRecorder {
	return m.recorder
}

// CreateTask mocks base method
func (m *MockTaskLogic) CreateTask(arg0 models.CreateTaskRequest) (string, error) {
	ret := m.ctrl.Call(m, "CreateTask", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask
func (mr *MockTaskLogicMockRecorder) CreateTask(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskLogic)(nil).CreateTask), arg0)
}

// DeleteTask mocks base method
func (m *MockTaskLogic) DeleteTask(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteTask", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask
func (mr *MockTaskLogicMockRecorder) DeleteTask(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskLogic)(nil).DeleteTask), arg0)
}

// GetEnvironmentTasks mocks base method
func (m *MockTaskLogic) GetEnvironmentTasks(arg0 string) ([]*models.Task, error) {
	ret := m.ctrl.Call(m, "GetEnvironmentTasks", arg0)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnvironmentTasks indicates an expected call of GetEnvironmentTasks
func (mr *MockTaskLogicMockRecorder) GetEnvironmentTasks(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnvironmentTasks", reflect.TypeOf((*MockTaskLogic)(nil).GetEnvironmentTasks), arg0)
}

// GetTask mocks base method
func (m *MockTaskLogic) GetTask(arg0 string) (*models.Task, error) {
	ret := m.ctrl.Call(m, "GetTask", arg0)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask
func (mr *MockTaskLogicMockRecorder) GetTask(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskLogic)(nil).GetTask), arg0)
}

// GetTaskLogs mocks base method
func (m *MockTaskLogic) GetTaskLogs(arg0, arg1, arg2 string, arg3 int) ([]*models.LogFile, error) {
	ret := m.ctrl.Call(m, "GetTaskLogs", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskLogs indicates an expected call of GetTaskLogs
func (mr *MockTaskLogicMockRecorder) GetTaskLogs(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskLogs", reflect.TypeOf((*MockTaskLogic)(nil).GetTaskLogs), arg0, arg1, arg2, arg3)
}

// ListTasks mocks base method
func (m *MockTaskLogic) ListTasks() ([]*models.TaskSummary, error) {
	ret := m.ctrl.Call(m, "ListTasks")
	ret0, _ := ret[0].([]*models.TaskSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks
func (mr *MockTaskLogicMockRecorder) ListTasks() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockTaskLogic)(nil).ListTasks))
}
