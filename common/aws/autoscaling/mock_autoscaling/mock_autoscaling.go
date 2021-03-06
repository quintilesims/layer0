// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/common/aws/autoscaling (interfaces: Provider)

// Package mock_autoscaling is a generated GoMock package.
package mock_autoscaling

import (
	gomock "github.com/golang/mock/gomock"
	autoscaling "github.com/quintilesims/layer0/common/aws/autoscaling"
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

// AttachLoadBalancer mocks base method
func (m *MockProvider) AttachLoadBalancer(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "AttachLoadBalancer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttachLoadBalancer indicates an expected call of AttachLoadBalancer
func (mr *MockProviderMockRecorder) AttachLoadBalancer(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachLoadBalancer", reflect.TypeOf((*MockProvider)(nil).AttachLoadBalancer), arg0, arg1)
}

// CreateAutoScalingGroup mocks base method
func (m *MockProvider) CreateAutoScalingGroup(arg0, arg1, arg2 string, arg3, arg4 int) error {
	ret := m.ctrl.Call(m, "CreateAutoScalingGroup", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAutoScalingGroup indicates an expected call of CreateAutoScalingGroup
func (mr *MockProviderMockRecorder) CreateAutoScalingGroup(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAutoScalingGroup", reflect.TypeOf((*MockProvider)(nil).CreateAutoScalingGroup), arg0, arg1, arg2, arg3, arg4)
}

// CreateLaunchConfiguration mocks base method
func (m *MockProvider) CreateLaunchConfiguration(arg0, arg1, arg2, arg3, arg4, arg5 *string, arg6 []*string, arg7 map[string]int) error {
	ret := m.ctrl.Call(m, "CreateLaunchConfiguration", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLaunchConfiguration indicates an expected call of CreateLaunchConfiguration
func (mr *MockProviderMockRecorder) CreateLaunchConfiguration(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLaunchConfiguration", reflect.TypeOf((*MockProvider)(nil).CreateLaunchConfiguration), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

// DeleteAutoScalingGroup mocks base method
func (m *MockProvider) DeleteAutoScalingGroup(arg0 *string) error {
	ret := m.ctrl.Call(m, "DeleteAutoScalingGroup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAutoScalingGroup indicates an expected call of DeleteAutoScalingGroup
func (mr *MockProviderMockRecorder) DeleteAutoScalingGroup(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAutoScalingGroup", reflect.TypeOf((*MockProvider)(nil).DeleteAutoScalingGroup), arg0)
}

// DeleteLaunchConfiguration mocks base method
func (m *MockProvider) DeleteLaunchConfiguration(arg0 *string) error {
	ret := m.ctrl.Call(m, "DeleteLaunchConfiguration", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLaunchConfiguration indicates an expected call of DeleteLaunchConfiguration
func (mr *MockProviderMockRecorder) DeleteLaunchConfiguration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLaunchConfiguration", reflect.TypeOf((*MockProvider)(nil).DeleteLaunchConfiguration), arg0)
}

// DescribeAutoScalingGroup mocks base method
func (m *MockProvider) DescribeAutoScalingGroup(arg0 string) (*autoscaling.Group, error) {
	ret := m.ctrl.Call(m, "DescribeAutoScalingGroup", arg0)
	ret0, _ := ret[0].(*autoscaling.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAutoScalingGroup indicates an expected call of DescribeAutoScalingGroup
func (mr *MockProviderMockRecorder) DescribeAutoScalingGroup(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAutoScalingGroup", reflect.TypeOf((*MockProvider)(nil).DescribeAutoScalingGroup), arg0)
}

// DescribeAutoScalingGroups mocks base method
func (m *MockProvider) DescribeAutoScalingGroups(arg0 []*string) ([]*autoscaling.Group, error) {
	ret := m.ctrl.Call(m, "DescribeAutoScalingGroups", arg0)
	ret0, _ := ret[0].([]*autoscaling.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAutoScalingGroups indicates an expected call of DescribeAutoScalingGroups
func (mr *MockProviderMockRecorder) DescribeAutoScalingGroups(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAutoScalingGroups", reflect.TypeOf((*MockProvider)(nil).DescribeAutoScalingGroups), arg0)
}

// DescribeLaunchConfiguration mocks base method
func (m *MockProvider) DescribeLaunchConfiguration(arg0 string) (*autoscaling.LaunchConfiguration, error) {
	ret := m.ctrl.Call(m, "DescribeLaunchConfiguration", arg0)
	ret0, _ := ret[0].(*autoscaling.LaunchConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeLaunchConfiguration indicates an expected call of DescribeLaunchConfiguration
func (mr *MockProviderMockRecorder) DescribeLaunchConfiguration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeLaunchConfiguration", reflect.TypeOf((*MockProvider)(nil).DescribeLaunchConfiguration), arg0)
}

// DescribeLaunchConfigurations mocks base method
func (m *MockProvider) DescribeLaunchConfigurations(arg0 []*string) ([]*autoscaling.LaunchConfiguration, error) {
	ret := m.ctrl.Call(m, "DescribeLaunchConfigurations", arg0)
	ret0, _ := ret[0].([]*autoscaling.LaunchConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeLaunchConfigurations indicates an expected call of DescribeLaunchConfigurations
func (mr *MockProviderMockRecorder) DescribeLaunchConfigurations(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeLaunchConfigurations", reflect.TypeOf((*MockProvider)(nil).DescribeLaunchConfigurations), arg0)
}

// SetDesiredCapacity mocks base method
func (m *MockProvider) SetDesiredCapacity(arg0 string, arg1 int) error {
	ret := m.ctrl.Call(m, "SetDesiredCapacity", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDesiredCapacity indicates an expected call of SetDesiredCapacity
func (mr *MockProviderMockRecorder) SetDesiredCapacity(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDesiredCapacity", reflect.TypeOf((*MockProvider)(nil).SetDesiredCapacity), arg0, arg1)
}

// TerminateInstanceInAutoScalingGroup mocks base method
func (m *MockProvider) TerminateInstanceInAutoScalingGroup(arg0 string, arg1 bool) (*autoscaling.Activity, error) {
	ret := m.ctrl.Call(m, "TerminateInstanceInAutoScalingGroup", arg0, arg1)
	ret0, _ := ret[0].(*autoscaling.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TerminateInstanceInAutoScalingGroup indicates an expected call of TerminateInstanceInAutoScalingGroup
func (mr *MockProviderMockRecorder) TerminateInstanceInAutoScalingGroup(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstanceInAutoScalingGroup", reflect.TypeOf((*MockProvider)(nil).TerminateInstanceInAutoScalingGroup), arg0, arg1)
}

// UpdateAutoScalingGroupMaxSize mocks base method
func (m *MockProvider) UpdateAutoScalingGroupMaxSize(arg0 string, arg1 int) error {
	ret := m.ctrl.Call(m, "UpdateAutoScalingGroupMaxSize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAutoScalingGroupMaxSize indicates an expected call of UpdateAutoScalingGroupMaxSize
func (mr *MockProviderMockRecorder) UpdateAutoScalingGroupMaxSize(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAutoScalingGroupMaxSize", reflect.TypeOf((*MockProvider)(nil).UpdateAutoScalingGroupMaxSize), arg0, arg1)
}

// UpdateAutoScalingGroupMinSize mocks base method
func (m *MockProvider) UpdateAutoScalingGroupMinSize(arg0 string, arg1 int) error {
	ret := m.ctrl.Call(m, "UpdateAutoScalingGroupMinSize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAutoScalingGroupMinSize indicates an expected call of UpdateAutoScalingGroupMinSize
func (mr *MockProviderMockRecorder) UpdateAutoScalingGroupMinSize(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAutoScalingGroupMinSize", reflect.TypeOf((*MockProvider)(nil).UpdateAutoScalingGroupMinSize), arg0, arg1)
}
