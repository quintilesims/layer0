// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/api/provider (interfaces: ServiceProvider)

package mock_provider

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	time "time"
)

// Mock of ServiceProvider interface
type MockServiceProvider struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceProviderRecorder
}

// Recorder for MockServiceProvider (not exported)
type _MockServiceProviderRecorder struct {
	mock *MockServiceProvider
}

func NewMockServiceProvider(ctrl *gomock.Controller) *MockServiceProvider {
	mock := &MockServiceProvider{ctrl: ctrl}
	mock.recorder = &_MockServiceProviderRecorder{mock}
	return mock
}

func (_m *MockServiceProvider) EXPECT() *_MockServiceProviderRecorder {
	return _m.recorder
}

func (_m *MockServiceProvider) Create(_param0 models.CreateServiceRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "Create", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceProviderRecorder) Create(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0)
}

func (_m *MockServiceProvider) Delete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceProviderRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockServiceProvider) List() ([]models.ServiceSummary, error) {
	ret := _m.ctrl.Call(_m, "List")
	ret0, _ := ret[0].([]models.ServiceSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceProviderRecorder) List() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "List")
}

func (_m *MockServiceProvider) Logs(_param0 string, _param1 int, _param2 time.Time, _param3 time.Time) ([]models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "Logs", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceProviderRecorder) Logs(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Logs", arg0, arg1, arg2, arg3)
}

func (_m *MockServiceProvider) Read(_param0 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "Read", _param0)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceProviderRecorder) Read(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Read", arg0)
}

func (_m *MockServiceProvider) Update(_param0 string, _param1 models.UpdateServiceRequest) error {
	ret := _m.ctrl.Call(_m, "Update", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceProviderRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0, arg1)
}
