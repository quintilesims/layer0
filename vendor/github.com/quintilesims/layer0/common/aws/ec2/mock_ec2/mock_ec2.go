// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/common/aws/ec2 (interfaces: Provider)

// Package mock_ec2 is a generated GoMock package.
package mock_ec2

import (
	gomock "github.com/golang/mock/gomock"
	ec2 "github.com/quintilesims/layer0/common/aws/ec2"
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

// AuthorizeSecurityGroupIngress mocks base method
func (m *MockProvider) AuthorizeSecurityGroupIngress(arg0 []*ec2.SecurityGroupIngress) error {
	ret := m.ctrl.Call(m, "AuthorizeSecurityGroupIngress", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthorizeSecurityGroupIngress indicates an expected call of AuthorizeSecurityGroupIngress
func (mr *MockProviderMockRecorder) AuthorizeSecurityGroupIngress(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizeSecurityGroupIngress", reflect.TypeOf((*MockProvider)(nil).AuthorizeSecurityGroupIngress), arg0)
}

// AuthorizeSecurityGroupIngressFromGroup mocks base method
func (m *MockProvider) AuthorizeSecurityGroupIngressFromGroup(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "AuthorizeSecurityGroupIngressFromGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthorizeSecurityGroupIngressFromGroup indicates an expected call of AuthorizeSecurityGroupIngressFromGroup
func (mr *MockProviderMockRecorder) AuthorizeSecurityGroupIngressFromGroup(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizeSecurityGroupIngressFromGroup", reflect.TypeOf((*MockProvider)(nil).AuthorizeSecurityGroupIngressFromGroup), arg0, arg1)
}

// CreateSecurityGroup mocks base method
func (m *MockProvider) CreateSecurityGroup(arg0, arg1, arg2 string) (*string, error) {
	ret := m.ctrl.Call(m, "CreateSecurityGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecurityGroup indicates an expected call of CreateSecurityGroup
func (mr *MockProviderMockRecorder) CreateSecurityGroup(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecurityGroup", reflect.TypeOf((*MockProvider)(nil).CreateSecurityGroup), arg0, arg1, arg2)
}

// DeleteSecurityGroup mocks base method
func (m *MockProvider) DeleteSecurityGroup(arg0 *ec2.SecurityGroup) error {
	ret := m.ctrl.Call(m, "DeleteSecurityGroup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecurityGroup indicates an expected call of DeleteSecurityGroup
func (mr *MockProviderMockRecorder) DeleteSecurityGroup(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecurityGroup", reflect.TypeOf((*MockProvider)(nil).DeleteSecurityGroup), arg0)
}

// DescribeInstance mocks base method
func (m *MockProvider) DescribeInstance(arg0 string) (*ec2.Instance, error) {
	ret := m.ctrl.Call(m, "DescribeInstance", arg0)
	ret0, _ := ret[0].(*ec2.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeInstance indicates an expected call of DescribeInstance
func (mr *MockProviderMockRecorder) DescribeInstance(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeInstance", reflect.TypeOf((*MockProvider)(nil).DescribeInstance), arg0)
}

// DescribeSecurityGroup mocks base method
func (m *MockProvider) DescribeSecurityGroup(arg0 string) (*ec2.SecurityGroup, error) {
	ret := m.ctrl.Call(m, "DescribeSecurityGroup", arg0)
	ret0, _ := ret[0].(*ec2.SecurityGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSecurityGroup indicates an expected call of DescribeSecurityGroup
func (mr *MockProviderMockRecorder) DescribeSecurityGroup(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSecurityGroup", reflect.TypeOf((*MockProvider)(nil).DescribeSecurityGroup), arg0)
}

// DescribeSubnet mocks base method
func (m *MockProvider) DescribeSubnet(arg0 string) (*ec2.Subnet, error) {
	ret := m.ctrl.Call(m, "DescribeSubnet", arg0)
	ret0, _ := ret[0].(*ec2.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSubnet indicates an expected call of DescribeSubnet
func (mr *MockProviderMockRecorder) DescribeSubnet(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSubnet", reflect.TypeOf((*MockProvider)(nil).DescribeSubnet), arg0)
}

// DescribeVPC mocks base method
func (m *MockProvider) DescribeVPC(arg0 string) (*ec2.VPC, error) {
	ret := m.ctrl.Call(m, "DescribeVPC", arg0)
	ret0, _ := ret[0].(*ec2.VPC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVPC indicates an expected call of DescribeVPC
func (mr *MockProviderMockRecorder) DescribeVPC(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVPC", reflect.TypeOf((*MockProvider)(nil).DescribeVPC), arg0)
}

// DescribeVPCByName mocks base method
func (m *MockProvider) DescribeVPCByName(arg0 string) (*ec2.VPC, error) {
	ret := m.ctrl.Call(m, "DescribeVPCByName", arg0)
	ret0, _ := ret[0].(*ec2.VPC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVPCByName indicates an expected call of DescribeVPCByName
func (mr *MockProviderMockRecorder) DescribeVPCByName(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVPCByName", reflect.TypeOf((*MockProvider)(nil).DescribeVPCByName), arg0)
}

// DescribeVPCGateways mocks base method
func (m *MockProvider) DescribeVPCGateways(arg0 string) ([]*ec2.InternetGateway, error) {
	ret := m.ctrl.Call(m, "DescribeVPCGateways", arg0)
	ret0, _ := ret[0].([]*ec2.InternetGateway)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVPCGateways indicates an expected call of DescribeVPCGateways
func (mr *MockProviderMockRecorder) DescribeVPCGateways(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVPCGateways", reflect.TypeOf((*MockProvider)(nil).DescribeVPCGateways), arg0)
}

// DescribeVPCRoutes mocks base method
func (m *MockProvider) DescribeVPCRoutes(arg0 string) ([]*ec2.RouteTable, error) {
	ret := m.ctrl.Call(m, "DescribeVPCRoutes", arg0)
	ret0, _ := ret[0].([]*ec2.RouteTable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVPCRoutes indicates an expected call of DescribeVPCRoutes
func (mr *MockProviderMockRecorder) DescribeVPCRoutes(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVPCRoutes", reflect.TypeOf((*MockProvider)(nil).DescribeVPCRoutes), arg0)
}

// DescribeVPCSubnets mocks base method
func (m *MockProvider) DescribeVPCSubnets(arg0 string) ([]*ec2.Subnet, error) {
	ret := m.ctrl.Call(m, "DescribeVPCSubnets", arg0)
	ret0, _ := ret[0].([]*ec2.Subnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVPCSubnets indicates an expected call of DescribeVPCSubnets
func (mr *MockProviderMockRecorder) DescribeVPCSubnets(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVPCSubnets", reflect.TypeOf((*MockProvider)(nil).DescribeVPCSubnets), arg0)
}

// RevokeSecurityGroupIngress mocks base method
func (m *MockProvider) RevokeSecurityGroupIngress(arg0 []*ec2.SecurityGroupIngress) error {
	ret := m.ctrl.Call(m, "RevokeSecurityGroupIngress", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeSecurityGroupIngress indicates an expected call of RevokeSecurityGroupIngress
func (mr *MockProviderMockRecorder) RevokeSecurityGroupIngress(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeSecurityGroupIngress", reflect.TypeOf((*MockProvider)(nil).RevokeSecurityGroupIngress), arg0)
}

// RevokeSecurityGroupIngressHelper mocks base method
func (m *MockProvider) RevokeSecurityGroupIngressHelper(arg0 string, arg1 ec2.IpPermission) error {
	ret := m.ctrl.Call(m, "RevokeSecurityGroupIngressHelper", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeSecurityGroupIngressHelper indicates an expected call of RevokeSecurityGroupIngressHelper
func (mr *MockProviderMockRecorder) RevokeSecurityGroupIngressHelper(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeSecurityGroupIngressHelper", reflect.TypeOf((*MockProvider)(nil).RevokeSecurityGroupIngressHelper), arg0, arg1)
}