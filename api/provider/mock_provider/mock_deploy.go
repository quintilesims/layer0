// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/provider (interfaces: DeployProvider)

// Package mock_provider is a generated GoMock package.
package mock_provider

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
	reflect "reflect"
)

// MockDeployProvider is a mock of DeployProvider interface
type MockDeployProvider struct {
	ctrl     *gomock.Controller
	recorder *MockDeployProviderMockRecorder
}

// MockDeployProviderMockRecorder is the mock recorder for MockDeployProvider
type MockDeployProviderMockRecorder struct {
	mock *MockDeployProvider
}

// NewMockDeployProvider creates a new mock instance
func NewMockDeployProvider(ctrl *gomock.Controller) *MockDeployProvider {
	mock := &MockDeployProvider{ctrl: ctrl}
	mock.recorder = &MockDeployProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDeployProvider) EXPECT() *MockDeployProviderMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockDeployProvider) Create(arg0 models.CreateDeployRequest) (*models.Deploy, error) {
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockDeployProviderMockRecorder) Create(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDeployProvider)(nil).Create), arg0)
}

// Delete mocks base method
func (m *MockDeployProvider) Delete(arg0 string) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockDeployProviderMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeployProvider)(nil).Delete), arg0)
}

// List mocks base method
func (m *MockDeployProvider) List() ([]models.DeploySummary, error) {
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]models.DeploySummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockDeployProviderMockRecorder) List() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDeployProvider)(nil).List))
}

// Read mocks base method
func (m *MockDeployProvider) Read(arg0 string) (*models.Deploy, error) {
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockDeployProviderMockRecorder) Read(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockDeployProvider)(nil).Read), arg0)
}
