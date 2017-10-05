// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/provider (interfaces: LoadBalancerProvider)

// Package mock_provider is a generated GoMock package.
package mock_provider

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	reflect "reflect"
)

// MockLoadBalancerProvider is a mock of LoadBalancerProvider interface
type MockLoadBalancerProvider struct {
	ctrl     *gomock.Controller
	recorder *MockLoadBalancerProviderMockRecorder
}

// MockLoadBalancerProviderMockRecorder is the mock recorder for MockLoadBalancerProvider
type MockLoadBalancerProviderMockRecorder struct {
	mock *MockLoadBalancerProvider
}

// NewMockLoadBalancerProvider creates a new mock instance
func NewMockLoadBalancerProvider(ctrl *gomock.Controller) *MockLoadBalancerProvider {
	mock := &MockLoadBalancerProvider{ctrl: ctrl}
	mock.recorder = &MockLoadBalancerProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLoadBalancerProvider) EXPECT() *MockLoadBalancerProviderMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockLoadBalancerProvider) Create(arg0 models.CreateLoadBalancerRequest) (string, error) {
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockLoadBalancerProviderMockRecorder) Create(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLoadBalancerProvider)(nil).Create), arg0)
}

// Delete mocks base method
func (m *MockLoadBalancerProvider) Delete(arg0 string) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockLoadBalancerProviderMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLoadBalancerProvider)(nil).Delete), arg0)
}

// List mocks base method
func (m *MockLoadBalancerProvider) List() ([]models.LoadBalancerSummary, error) {
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]models.LoadBalancerSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockLoadBalancerProviderMockRecorder) List() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockLoadBalancerProvider)(nil).List))
}

// Read mocks base method
func (m *MockLoadBalancerProvider) Read(arg0 string) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockLoadBalancerProviderMockRecorder) Read(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockLoadBalancerProvider)(nil).Read), arg0)
}

// Update mocks base method
func (m *MockLoadBalancerProvider) Update(arg0 models.UpdateLoadBalancerRequest) error {
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockLoadBalancerProviderMockRecorder) Update(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLoadBalancerProvider)(nil).Update), arg0)
}
