// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/common/config (interfaces: APIConfig)

// Package mock_config is a generated GoMock package.
package mock_config

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockAPIConfig is a mock of APIConfig interface
type MockAPIConfig struct {
	ctrl     *gomock.Controller
	recorder *MockAPIConfigMockRecorder
}

// MockAPIConfigMockRecorder is the mock recorder for MockAPIConfig
type MockAPIConfigMockRecorder struct {
	mock *MockAPIConfig
}

// NewMockAPIConfig creates a new mock instance
func NewMockAPIConfig(ctrl *gomock.Controller) *MockAPIConfig {
	mock := &MockAPIConfig{ctrl: ctrl}
	mock.recorder = &MockAPIConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPIConfig) EXPECT() *MockAPIConfigMockRecorder {
	return m.recorder
}

// AccessKey mocks base method
func (m *MockAPIConfig) AccessKey() string {
	ret := m.ctrl.Call(m, "AccessKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// AccessKey indicates an expected call of AccessKey
func (mr *MockAPIConfigMockRecorder) AccessKey() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccessKey", reflect.TypeOf((*MockAPIConfig)(nil).AccessKey))
}

// AccountID mocks base method
func (m *MockAPIConfig) AccountID() string {
	ret := m.ctrl.Call(m, "AccountID")
	ret0, _ := ret[0].(string)
	return ret0
}

// AccountID indicates an expected call of AccountID
func (mr *MockAPIConfigMockRecorder) AccountID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountID", reflect.TypeOf((*MockAPIConfig)(nil).AccountID))
}

// DynamoLockTable mocks base method
func (m *MockAPIConfig) DynamoLockTable() string {
	ret := m.ctrl.Call(m, "DynamoLockTable")
	ret0, _ := ret[0].(string)
	return ret0
}

// DynamoLockTable indicates an expected call of DynamoLockTable
func (mr *MockAPIConfigMockRecorder) DynamoLockTable() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DynamoLockTable", reflect.TypeOf((*MockAPIConfig)(nil).DynamoLockTable))
}

// DynamoTagTable mocks base method
func (m *MockAPIConfig) DynamoTagTable() string {
	ret := m.ctrl.Call(m, "DynamoTagTable")
	ret0, _ := ret[0].(string)
	return ret0
}

// DynamoTagTable indicates an expected call of DynamoTagTable
func (mr *MockAPIConfigMockRecorder) DynamoTagTable() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DynamoTagTable", reflect.TypeOf((*MockAPIConfig)(nil).DynamoTagTable))
}

// Instance mocks base method
func (m *MockAPIConfig) Instance() string {
	ret := m.ctrl.Call(m, "Instance")
	ret0, _ := ret[0].(string)
	return ret0
}

// Instance indicates an expected call of Instance
func (mr *MockAPIConfigMockRecorder) Instance() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Instance", reflect.TypeOf((*MockAPIConfig)(nil).Instance))
}

// InstanceProfile mocks base method
func (m *MockAPIConfig) InstanceProfile() string {
	ret := m.ctrl.Call(m, "InstanceProfile")
	ret0, _ := ret[0].(string)
	return ret0
}

// InstanceProfile indicates an expected call of InstanceProfile
func (mr *MockAPIConfigMockRecorder) InstanceProfile() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceProfile", reflect.TypeOf((*MockAPIConfig)(nil).InstanceProfile))
}

// LinuxAMI mocks base method
func (m *MockAPIConfig) LinuxAMI() string {
	ret := m.ctrl.Call(m, "LinuxAMI")
	ret0, _ := ret[0].(string)
	return ret0
}

// LinuxAMI indicates an expected call of LinuxAMI
func (mr *MockAPIConfigMockRecorder) LinuxAMI() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LinuxAMI", reflect.TypeOf((*MockAPIConfig)(nil).LinuxAMI))
}

// LockExpiry mocks base method
func (m *MockAPIConfig) LockExpiry() time.Duration {
	ret := m.ctrl.Call(m, "LockExpiry")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// LockExpiry indicates an expected call of LockExpiry
func (mr *MockAPIConfigMockRecorder) LockExpiry() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockExpiry", reflect.TypeOf((*MockAPIConfig)(nil).LockExpiry))
}

// LogGroupName mocks base method
func (m *MockAPIConfig) LogGroupName() string {
	ret := m.ctrl.Call(m, "LogGroupName")
	ret0, _ := ret[0].(string)
	return ret0
}

// LogGroupName indicates an expected call of LogGroupName
func (mr *MockAPIConfigMockRecorder) LogGroupName() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogGroupName", reflect.TypeOf((*MockAPIConfig)(nil).LogGroupName))
}

// MaxRetries mocks base method
func (m *MockAPIConfig) MaxRetries() int {
	ret := m.ctrl.Call(m, "MaxRetries")
	ret0, _ := ret[0].(int)
	return ret0
}

// MaxRetries indicates an expected call of MaxRetries
func (mr *MockAPIConfigMockRecorder) MaxRetries() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaxRetries", reflect.TypeOf((*MockAPIConfig)(nil).MaxRetries))
}

// ParseAuthToken mocks base method
func (m *MockAPIConfig) ParseAuthToken() (string, string, error) {
	ret := m.ctrl.Call(m, "ParseAuthToken")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ParseAuthToken indicates an expected call of ParseAuthToken
func (mr *MockAPIConfigMockRecorder) ParseAuthToken() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAuthToken", reflect.TypeOf((*MockAPIConfig)(nil).ParseAuthToken))
}

// Port mocks base method
func (m *MockAPIConfig) Port() int {
	ret := m.ctrl.Call(m, "Port")
	ret0, _ := ret[0].(int)
	return ret0
}

// Port indicates an expected call of Port
func (mr *MockAPIConfigMockRecorder) Port() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Port", reflect.TypeOf((*MockAPIConfig)(nil).Port))
}

// PrivateSubnets mocks base method
func (m *MockAPIConfig) PrivateSubnets() []string {
	ret := m.ctrl.Call(m, "PrivateSubnets")
	ret0, _ := ret[0].([]string)
	return ret0
}

// PrivateSubnets indicates an expected call of PrivateSubnets
func (mr *MockAPIConfigMockRecorder) PrivateSubnets() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrivateSubnets", reflect.TypeOf((*MockAPIConfig)(nil).PrivateSubnets))
}

// PublicSubnets mocks base method
func (m *MockAPIConfig) PublicSubnets() []string {
	ret := m.ctrl.Call(m, "PublicSubnets")
	ret0, _ := ret[0].([]string)
	return ret0
}

// PublicSubnets indicates an expected call of PublicSubnets
func (mr *MockAPIConfigMockRecorder) PublicSubnets() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicSubnets", reflect.TypeOf((*MockAPIConfig)(nil).PublicSubnets))
}

// Region mocks base method
func (m *MockAPIConfig) Region() string {
	ret := m.ctrl.Call(m, "Region")
	ret0, _ := ret[0].(string)
	return ret0
}

// Region indicates an expected call of Region
func (mr *MockAPIConfigMockRecorder) Region() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Region", reflect.TypeOf((*MockAPIConfig)(nil).Region))
}

// S3Bucket mocks base method
func (m *MockAPIConfig) S3Bucket() string {
	ret := m.ctrl.Call(m, "S3Bucket")
	ret0, _ := ret[0].(string)
	return ret0
}

// S3Bucket indicates an expected call of S3Bucket
func (mr *MockAPIConfigMockRecorder) S3Bucket() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "S3Bucket", reflect.TypeOf((*MockAPIConfig)(nil).S3Bucket))
}

// SSHKeyPair mocks base method
func (m *MockAPIConfig) SSHKeyPair() string {
	ret := m.ctrl.Call(m, "SSHKeyPair")
	ret0, _ := ret[0].(string)
	return ret0
}

// SSHKeyPair indicates an expected call of SSHKeyPair
func (mr *MockAPIConfigMockRecorder) SSHKeyPair() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SSHKeyPair", reflect.TypeOf((*MockAPIConfig)(nil).SSHKeyPair))
}

// SecretKey mocks base method
func (m *MockAPIConfig) SecretKey() string {
	ret := m.ctrl.Call(m, "SecretKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// SecretKey indicates an expected call of SecretKey
func (mr *MockAPIConfigMockRecorder) SecretKey() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecretKey", reflect.TypeOf((*MockAPIConfig)(nil).SecretKey))
}

// TimeBetweenRequests mocks base method
func (m *MockAPIConfig) TimeBetweenRequests() time.Duration {
	ret := m.ctrl.Call(m, "TimeBetweenRequests")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// TimeBetweenRequests indicates an expected call of TimeBetweenRequests
func (mr *MockAPIConfigMockRecorder) TimeBetweenRequests() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeBetweenRequests", reflect.TypeOf((*MockAPIConfig)(nil).TimeBetweenRequests))
}

// VPC mocks base method
func (m *MockAPIConfig) VPC() string {
	ret := m.ctrl.Call(m, "VPC")
	ret0, _ := ret[0].(string)
	return ret0
}

// VPC indicates an expected call of VPC
func (mr *MockAPIConfigMockRecorder) VPC() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VPC", reflect.TypeOf((*MockAPIConfig)(nil).VPC))
}
