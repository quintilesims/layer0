// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/common/aws/elb (interfaces: Provider)

// Package mock_elb is a generated GoMock package.
package mock_elb

import (
	gomock "github.com/golang/mock/gomock"
	elb "github.com/quintilesims/layer0/common/aws/elb"
	reflect "reflect"
)

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// ConfigureHealthCheck mocks base method
func (m *MockProvider) ConfigureHealthCheck(arg0 string, arg1 *elb.HealthCheck) error {
	ret := m.ctrl.Call(m, "ConfigureHealthCheck", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfigureHealthCheck indicates an expected call of ConfigureHealthCheck
func (mr *MockProviderMockRecorder) ConfigureHealthCheck(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureHealthCheck", reflect.TypeOf((*MockProvider)(nil).ConfigureHealthCheck), arg0, arg1)
}

// CreateLoadBalancer mocks base method
func (m *MockProvider) CreateLoadBalancer(arg0, arg1 string, arg2, arg3 []*string, arg4 []*elb.Listener) (*string, error) {
	ret := m.ctrl.Call(m, "CreateLoadBalancer", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoadBalancer indicates an expected call of CreateLoadBalancer
func (mr *MockProviderMockRecorder) CreateLoadBalancer(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoadBalancer", reflect.TypeOf((*MockProvider)(nil).CreateLoadBalancer), arg0, arg1, arg2, arg3, arg4)
}

// CreateLoadBalancerListeners mocks base method
func (m *MockProvider) CreateLoadBalancerListeners(arg0 string, arg1 []*elb.Listener) error {
	ret := m.ctrl.Call(m, "CreateLoadBalancerListeners", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLoadBalancerListeners indicates an expected call of CreateLoadBalancerListeners
func (mr *MockProviderMockRecorder) CreateLoadBalancerListeners(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoadBalancerListeners", reflect.TypeOf((*MockProvider)(nil).CreateLoadBalancerListeners), arg0, arg1)
}

// DeleteLoadBalancer mocks base method
func (m *MockProvider) DeleteLoadBalancer(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteLoadBalancer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLoadBalancer indicates an expected call of DeleteLoadBalancer
func (mr *MockProviderMockRecorder) DeleteLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLoadBalancer", reflect.TypeOf((*MockProvider)(nil).DeleteLoadBalancer), arg0)
}

// DeleteLoadBalancerListeners mocks base method
func (m *MockProvider) DeleteLoadBalancerListeners(arg0 string, arg1 []*elb.Listener) error {
	ret := m.ctrl.Call(m, "DeleteLoadBalancerListeners", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLoadBalancerListeners indicates an expected call of DeleteLoadBalancerListeners
func (mr *MockProviderMockRecorder) DeleteLoadBalancerListeners(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLoadBalancerListeners", reflect.TypeOf((*MockProvider)(nil).DeleteLoadBalancerListeners), arg0, arg1)
}

// DeregisterInstancesFromLoadBalancer mocks base method
func (m *MockProvider) DeregisterInstancesFromLoadBalancer(arg0 string, arg1 []string) error {
	ret := m.ctrl.Call(m, "DeregisterInstancesFromLoadBalancer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeregisterInstancesFromLoadBalancer indicates an expected call of DeregisterInstancesFromLoadBalancer
func (mr *MockProviderMockRecorder) DeregisterInstancesFromLoadBalancer(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterInstancesFromLoadBalancer", reflect.TypeOf((*MockProvider)(nil).DeregisterInstancesFromLoadBalancer), arg0, arg1)
}

// DescribeInstanceHealth mocks base method
func (m *MockProvider) DescribeInstanceHealth(arg0 string) ([]*elb.InstanceState, error) {
	ret := m.ctrl.Call(m, "DescribeInstanceHealth", arg0)
	ret0, _ := ret[0].([]*elb.InstanceState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeInstanceHealth indicates an expected call of DescribeInstanceHealth
func (mr *MockProviderMockRecorder) DescribeInstanceHealth(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeInstanceHealth", reflect.TypeOf((*MockProvider)(nil).DescribeInstanceHealth), arg0)
}

// DescribeLoadBalancer mocks base method
func (m *MockProvider) DescribeLoadBalancer(arg0 string) (*elb.LoadBalancerDescription, error) {
	ret := m.ctrl.Call(m, "DescribeLoadBalancer", arg0)
	ret0, _ := ret[0].(*elb.LoadBalancerDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeLoadBalancer indicates an expected call of DescribeLoadBalancer
func (mr *MockProviderMockRecorder) DescribeLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeLoadBalancer", reflect.TypeOf((*MockProvider)(nil).DescribeLoadBalancer), arg0)
}

// DescribeLoadBalancers mocks base method
func (m *MockProvider) DescribeLoadBalancers() ([]*elb.LoadBalancerDescription, error) {
	ret := m.ctrl.Call(m, "DescribeLoadBalancers")
	ret0, _ := ret[0].([]*elb.LoadBalancerDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeLoadBalancers indicates an expected call of DescribeLoadBalancers
func (mr *MockProviderMockRecorder) DescribeLoadBalancers() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeLoadBalancers", reflect.TypeOf((*MockProvider)(nil).DescribeLoadBalancers))
}

// RegisterInstancesWithLoadBalancer mocks base method
func (m *MockProvider) RegisterInstancesWithLoadBalancer(arg0 string, arg1 []string) error {
	ret := m.ctrl.Call(m, "RegisterInstancesWithLoadBalancer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterInstancesWithLoadBalancer indicates an expected call of RegisterInstancesWithLoadBalancer
func (mr *MockProviderMockRecorder) RegisterInstancesWithLoadBalancer(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterInstancesWithLoadBalancer", reflect.TypeOf((*MockProvider)(nil).RegisterInstancesWithLoadBalancer), arg0, arg1)
}
