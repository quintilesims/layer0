// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/api/job (interfaces: Store)

package mock_job

import (
	gomock "github.com/golang/mock/gomock"
	job "github.com/quintilesims/layer0/api/job"
	models "github.com/quintilesims/layer0/common/models"
)

// Mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *_MockStoreRecorder
}

// Recorder for MockStore (not exported)
type _MockStoreRecorder struct {
	mock *MockStore
}

func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &_MockStoreRecorder{mock}
	return mock
}

func (_m *MockStore) EXPECT() *_MockStoreRecorder {
	return _m.recorder
}

func (_m *MockStore) AcquireJob(_param0 string) (bool, error) {
	ret := _m.ctrl.Call(_m, "AcquireJob", _param0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) AcquireJob(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AcquireJob", arg0)
}

func (_m *MockStore) Delete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockStore) Insert(_param0 job.JobType, _param1 string) (string, error) {
	ret := _m.ctrl.Call(_m, "Insert", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Insert", arg0, arg1)
}

func (_m *MockStore) SelectAll() ([]*models.Job, error) {
	ret := _m.ctrl.Call(_m, "SelectAll")
	ret0, _ := ret[0].([]*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) SelectAll() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SelectAll")
}

func (_m *MockStore) SelectByID(_param0 string) (*models.Job, error) {
	ret := _m.ctrl.Call(_m, "SelectByID", _param0)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockStoreRecorder) SelectByID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SelectByID", arg0)
}

func (_m *MockStore) SetJobError(_param0 string, _param1 error) error {
	ret := _m.ctrl.Call(_m, "SetJobError", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) SetJobError(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetJobError", arg0, arg1)
}

func (_m *MockStore) SetJobResult(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "SetJobResult", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) SetJobResult(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetJobResult", arg0, arg1)
}

func (_m *MockStore) SetJobStatus(_param0 string, _param1 job.Status) error {
	ret := _m.ctrl.Call(_m, "SetJobStatus", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockStoreRecorder) SetJobStatus(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetJobStatus", arg0, arg1)
}
