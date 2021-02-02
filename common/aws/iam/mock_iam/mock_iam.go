// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/common/aws/iam (interfaces: Provider)

// Package mock_iam is a generated GoMock package.
package mock_iam

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	iam "github.com/quintilesims/layer0/common/aws/iam"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// CreateRole mocks base method.
func (m *MockProvider) CreateRole(arg0, arg1 string) (*iam.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", arg0, arg1)
	ret0, _ := ret[0].(*iam.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockProviderMockRecorder) CreateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockProvider)(nil).CreateRole), arg0, arg1)
}

// DeleteRole mocks base method.
func (m *MockProvider) DeleteRole(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRole", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRole indicates an expected call of DeleteRole.
func (mr *MockProviderMockRecorder) DeleteRole(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRole", reflect.TypeOf((*MockProvider)(nil).DeleteRole), arg0)
}

// DeleteRolePolicy mocks base method.
func (m *MockProvider) DeleteRolePolicy(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRolePolicy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRolePolicy indicates an expected call of DeleteRolePolicy.
func (mr *MockProviderMockRecorder) DeleteRolePolicy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRolePolicy", reflect.TypeOf((*MockProvider)(nil).DeleteRolePolicy), arg0, arg1)
}

// DeleteServerCertificate mocks base method.
func (m *MockProvider) DeleteServerCertificate(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteServerCertificate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteServerCertificate indicates an expected call of DeleteServerCertificate.
func (mr *MockProviderMockRecorder) DeleteServerCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteServerCertificate", reflect.TypeOf((*MockProvider)(nil).DeleteServerCertificate), arg0)
}

// GetAccountId mocks base method.
func (m *MockProvider) GetAccountId() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountId")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountId indicates an expected call of GetAccountId.
func (mr *MockProviderMockRecorder) GetAccountId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountId", reflect.TypeOf((*MockProvider)(nil).GetAccountId))
}

// GetRole mocks base method.
func (m *MockProvider) GetRole(arg0 string) (*iam.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", arg0)
	ret0, _ := ret[0].(*iam.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockProviderMockRecorder) GetRole(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockProvider)(nil).GetRole), arg0)
}

// GetUser mocks base method.
func (m *MockProvider) GetUser(arg0 *string) (*iam.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(*iam.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockProviderMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockProvider)(nil).GetUser), arg0)
}

// ListCertificates mocks base method.
func (m *MockProvider) ListCertificates() ([]*iam.ServerCertificateMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCertificates")
	ret0, _ := ret[0].([]*iam.ServerCertificateMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCertificates indicates an expected call of ListCertificates.
func (mr *MockProviderMockRecorder) ListCertificates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCertificates", reflect.TypeOf((*MockProvider)(nil).ListCertificates))
}

// ListRolePolicies mocks base method.
func (m *MockProvider) ListRolePolicies(arg0 string) ([]*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRolePolicies", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRolePolicies indicates an expected call of ListRolePolicies.
func (mr *MockProviderMockRecorder) ListRolePolicies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRolePolicies", reflect.TypeOf((*MockProvider)(nil).ListRolePolicies), arg0)
}

// ListRoles mocks base method.
func (m *MockProvider) ListRoles() ([]*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRoles")
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoles indicates an expected call of ListRoles.
func (mr *MockProviderMockRecorder) ListRoles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoles", reflect.TypeOf((*MockProvider)(nil).ListRoles))
}

// PutRolePolicy mocks base method.
func (m *MockProvider) PutRolePolicy(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutRolePolicy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutRolePolicy indicates an expected call of PutRolePolicy.
func (mr *MockProviderMockRecorder) PutRolePolicy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutRolePolicy", reflect.TypeOf((*MockProvider)(nil).PutRolePolicy), arg0, arg1)
}

// UploadServerCertificate mocks base method.
func (m *MockProvider) UploadServerCertificate(arg0, arg1, arg2, arg3 string, arg4 *string) (*iam.ServerCertificateMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadServerCertificate", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*iam.ServerCertificateMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadServerCertificate indicates an expected call of UploadServerCertificate.
func (mr *MockProviderMockRecorder) UploadServerCertificate(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadServerCertificate", reflect.TypeOf((*MockProvider)(nil).UploadServerCertificate), arg0, arg1, arg2, arg3, arg4)
}
