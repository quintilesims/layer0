// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/api/provider (interfaces: LoadBalancerProvider)

package mock_provider

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
)

// Mock of LoadBalancerProvider interface
type MockLoadBalancerProvider struct {
	ctrl     *gomock.Controller
	recorder *_MockLoadBalancerProviderRecorder
}

// Recorder for MockLoadBalancerProvider (not exported)
type _MockLoadBalancerProviderRecorder struct {
	mock *MockLoadBalancerProvider
}

func NewMockLoadBalancerProvider(ctrl *gomock.Controller) *MockLoadBalancerProvider {
	mock := &MockLoadBalancerProvider{ctrl: ctrl}
	mock.recorder = &_MockLoadBalancerProviderRecorder{mock}
	return mock
}

func (_m *MockLoadBalancerProvider) EXPECT() *_MockLoadBalancerProviderRecorder {
	return _m.recorder
}

func (_m *MockLoadBalancerProvider) Create(_param0 models.CreateLoadBalancerRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "Create", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockLoadBalancerProviderRecorder) Create(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0)
}

func (_m *MockLoadBalancerProvider) Delete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockLoadBalancerProviderRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockLoadBalancerProvider) List() ([]models.LoadBalancerSummary, error) {
	ret := _m.ctrl.Call(_m, "List")
	ret0, _ := ret[0].([]models.LoadBalancerSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockLoadBalancerProviderRecorder) List() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "List")
}

func (_m *MockLoadBalancerProvider) Read(_param0 string) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "Read", _param0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockLoadBalancerProviderRecorder) Read(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Read", arg0)
}

func (_m *MockLoadBalancerProvider) Update(_param0 string, _param1 models.UpdateLoadBalancerRequest) error {
	ret := _m.ctrl.Call(_m, "Update", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockLoadBalancerProviderRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0, arg1)
}
