// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/logic (interfaces: JobLogic)

// Package mock_logic is a generated GoMock package.
package mock_logic

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	types "github.com/quintilesims/layer0/common/types"
	reflect "reflect"
)

// MockJobLogic is a mock of JobLogic interface
type MockJobLogic struct {
	ctrl     *gomock.Controller
	recorder *MockJobLogicMockRecorder
}

// MockJobLogicMockRecorder is the mock recorder for MockJobLogic
type MockJobLogicMockRecorder struct {
	mock *MockJobLogic
}

// NewMockJobLogic creates a new mock instance
func NewMockJobLogic(ctrl *gomock.Controller) *MockJobLogic {
	mock := &MockJobLogic{ctrl: ctrl}
	mock.recorder = &MockJobLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobLogic) EXPECT() *MockJobLogicMockRecorder {
	return m.recorder
}

// CreateJob mocks base method
func (m *MockJobLogic) CreateJob(arg0 types.JobType, arg1 interface{}) (*models.Job, error) {
	ret := m.ctrl.Call(m, "CreateJob", arg0, arg1)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJob indicates an expected call of CreateJob
func (mr *MockJobLogicMockRecorder) CreateJob(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJob", reflect.TypeOf((*MockJobLogic)(nil).CreateJob), arg0, arg1)
}

// Delete mocks base method
func (m *MockJobLogic) Delete(arg0 string) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockJobLogicMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockJobLogic)(nil).Delete), arg0)
}

// GetJob mocks base method
func (m *MockJobLogic) GetJob(arg0 string) (*models.Job, error) {
	ret := m.ctrl.Call(m, "GetJob", arg0)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJob indicates an expected call of GetJob
func (mr *MockJobLogicMockRecorder) GetJob(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJob", reflect.TypeOf((*MockJobLogic)(nil).GetJob), arg0)
}

// ListJobs mocks base method
func (m *MockJobLogic) ListJobs() ([]*models.Job, error) {
	ret := m.ctrl.Call(m, "ListJobs")
	ret0, _ := ret[0].([]*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListJobs indicates an expected call of ListJobs
func (mr *MockJobLogicMockRecorder) ListJobs() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListJobs", reflect.TypeOf((*MockJobLogic)(nil).ListJobs))
}
