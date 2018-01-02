// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/logic (interfaces: LoadBalancerLogic)

// Package mock_logic is a generated GoMock package.
package mock_logic

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	reflect "reflect"
)

// MockLoadBalancerLogic is a mock of LoadBalancerLogic interface
type MockLoadBalancerLogic struct {
	ctrl     *gomock.Controller
	recorder *MockLoadBalancerLogicMockRecorder
}

// MockLoadBalancerLogicMockRecorder is the mock recorder for MockLoadBalancerLogic
type MockLoadBalancerLogicMockRecorder struct {
	mock *MockLoadBalancerLogic
}

// NewMockLoadBalancerLogic creates a new mock instance
func NewMockLoadBalancerLogic(ctrl *gomock.Controller) *MockLoadBalancerLogic {
	mock := &MockLoadBalancerLogic{ctrl: ctrl}
	mock.recorder = &MockLoadBalancerLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLoadBalancerLogic) EXPECT() *MockLoadBalancerLogicMockRecorder {
	return m.recorder
}

// CreateLoadBalancer mocks base method
func (m *MockLoadBalancerLogic) CreateLoadBalancer(arg0 models.CreateLoadBalancerRequest) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "CreateLoadBalancer", arg0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoadBalancer indicates an expected call of CreateLoadBalancer
func (mr *MockLoadBalancerLogicMockRecorder) CreateLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoadBalancer", reflect.TypeOf((*MockLoadBalancerLogic)(nil).CreateLoadBalancer), arg0)
}

// DeleteLoadBalancer mocks base method
func (m *MockLoadBalancerLogic) DeleteLoadBalancer(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteLoadBalancer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLoadBalancer indicates an expected call of DeleteLoadBalancer
func (mr *MockLoadBalancerLogicMockRecorder) DeleteLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLoadBalancer", reflect.TypeOf((*MockLoadBalancerLogic)(nil).DeleteLoadBalancer), arg0)
}

// GetLoadBalancer mocks base method
func (m *MockLoadBalancerLogic) GetLoadBalancer(arg0 string) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "GetLoadBalancer", arg0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoadBalancer indicates an expected call of GetLoadBalancer
func (mr *MockLoadBalancerLogicMockRecorder) GetLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoadBalancer", reflect.TypeOf((*MockLoadBalancerLogic)(nil).GetLoadBalancer), arg0)
}

// ListLoadBalancers mocks base method
func (m *MockLoadBalancerLogic) ListLoadBalancers() ([]*models.LoadBalancerSummary, error) {
	ret := m.ctrl.Call(m, "ListLoadBalancers")
	ret0, _ := ret[0].([]*models.LoadBalancerSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLoadBalancers indicates an expected call of ListLoadBalancers
func (mr *MockLoadBalancerLogicMockRecorder) ListLoadBalancers() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLoadBalancers", reflect.TypeOf((*MockLoadBalancerLogic)(nil).ListLoadBalancers))
}

// UpdateLoadBalancerHealthCheck mocks base method
func (m *MockLoadBalancerLogic) UpdateLoadBalancerHealthCheck(arg0 string, arg1 models.HealthCheck) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "UpdateLoadBalancerHealthCheck", arg0, arg1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLoadBalancerHealthCheck indicates an expected call of UpdateLoadBalancerHealthCheck
func (mr *MockLoadBalancerLogicMockRecorder) UpdateLoadBalancerHealthCheck(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLoadBalancerHealthCheck", reflect.TypeOf((*MockLoadBalancerLogic)(nil).UpdateLoadBalancerHealthCheck), arg0, arg1)
}

// UpdateLoadBalancerPorts mocks base method
func (m *MockLoadBalancerLogic) UpdateLoadBalancerPorts(arg0 string, arg1 []models.Port) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "UpdateLoadBalancerPorts", arg0, arg1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLoadBalancerPorts indicates an expected call of UpdateLoadBalancerPorts
func (mr *MockLoadBalancerLogicMockRecorder) UpdateLoadBalancerPorts(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLoadBalancerPorts", reflect.TypeOf((*MockLoadBalancerLogic)(nil).UpdateLoadBalancerPorts), arg0, arg1)
}
