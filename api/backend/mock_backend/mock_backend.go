// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/api/backend (interfaces: Backend)

// Package mock_backend is a generated GoMock package.
package mock_backend

import (
	gomock "github.com/golang/mock/gomock"
	id "github.com/quintilesims/layer0/api/backend/ecs/id"
	models "github.com/quintilesims/layer0/common/models"
	reflect "reflect"
)

// MockBackend is a mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// CreateDeploy mocks base method
func (m *MockBackend) CreateDeploy(arg0 string, arg1 []byte) (*models.Deploy, error) {
	ret := m.ctrl.Call(m, "CreateDeploy", arg0, arg1)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDeploy indicates an expected call of CreateDeploy
func (mr *MockBackendMockRecorder) CreateDeploy(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeploy", reflect.TypeOf((*MockBackend)(nil).CreateDeploy), arg0, arg1)
}

// CreateEnvironment mocks base method
func (m *MockBackend) CreateEnvironment(arg0, arg1, arg2, arg3 string, arg4 int, arg5 []byte) (*models.Environment, error) {
	ret := m.ctrl.Call(m, "CreateEnvironment", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEnvironment indicates an expected call of CreateEnvironment
func (mr *MockBackendMockRecorder) CreateEnvironment(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEnvironment", reflect.TypeOf((*MockBackend)(nil).CreateEnvironment), arg0, arg1, arg2, arg3, arg4, arg5)
}

// CreateEnvironmentLink mocks base method
func (m *MockBackend) CreateEnvironmentLink(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "CreateEnvironmentLink", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEnvironmentLink indicates an expected call of CreateEnvironmentLink
func (mr *MockBackendMockRecorder) CreateEnvironmentLink(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEnvironmentLink", reflect.TypeOf((*MockBackend)(nil).CreateEnvironmentLink), arg0, arg1)
}

// CreateLoadBalancer mocks base method
func (m *MockBackend) CreateLoadBalancer(arg0, arg1 string, arg2 bool, arg3 []models.Port, arg4 models.HealthCheck) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "CreateLoadBalancer", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoadBalancer indicates an expected call of CreateLoadBalancer
func (mr *MockBackendMockRecorder) CreateLoadBalancer(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoadBalancer", reflect.TypeOf((*MockBackend)(nil).CreateLoadBalancer), arg0, arg1, arg2, arg3, arg4)
}

// CreateService mocks base method
func (m *MockBackend) CreateService(arg0, arg1, arg2, arg3 string) (*models.Service, error) {
	ret := m.ctrl.Call(m, "CreateService", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateService indicates an expected call of CreateService
func (mr *MockBackendMockRecorder) CreateService(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateService", reflect.TypeOf((*MockBackend)(nil).CreateService), arg0, arg1, arg2, arg3)
}

// CreateTask mocks base method
func (m *MockBackend) CreateTask(arg0, arg1 string, arg2 []models.ContainerOverride) (string, error) {
	ret := m.ctrl.Call(m, "CreateTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask
func (mr *MockBackendMockRecorder) CreateTask(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockBackend)(nil).CreateTask), arg0, arg1, arg2)
}

// DeleteDeploy mocks base method
func (m *MockBackend) DeleteDeploy(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteDeploy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDeploy indicates an expected call of DeleteDeploy
func (mr *MockBackendMockRecorder) DeleteDeploy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDeploy", reflect.TypeOf((*MockBackend)(nil).DeleteDeploy), arg0)
}

// DeleteEnvironment mocks base method
func (m *MockBackend) DeleteEnvironment(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteEnvironment", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEnvironment indicates an expected call of DeleteEnvironment
func (mr *MockBackendMockRecorder) DeleteEnvironment(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEnvironment", reflect.TypeOf((*MockBackend)(nil).DeleteEnvironment), arg0)
}

// DeleteEnvironmentLink mocks base method
func (m *MockBackend) DeleteEnvironmentLink(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteEnvironmentLink", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEnvironmentLink indicates an expected call of DeleteEnvironmentLink
func (mr *MockBackendMockRecorder) DeleteEnvironmentLink(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEnvironmentLink", reflect.TypeOf((*MockBackend)(nil).DeleteEnvironmentLink), arg0, arg1)
}

// DeleteLoadBalancer mocks base method
func (m *MockBackend) DeleteLoadBalancer(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteLoadBalancer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLoadBalancer indicates an expected call of DeleteLoadBalancer
func (mr *MockBackendMockRecorder) DeleteLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLoadBalancer", reflect.TypeOf((*MockBackend)(nil).DeleteLoadBalancer), arg0)
}

// DeleteService mocks base method
func (m *MockBackend) DeleteService(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteService", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteService indicates an expected call of DeleteService
func (mr *MockBackendMockRecorder) DeleteService(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteService", reflect.TypeOf((*MockBackend)(nil).DeleteService), arg0, arg1)
}

// DeleteTask mocks base method
func (m *MockBackend) DeleteTask(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask
func (mr *MockBackendMockRecorder) DeleteTask(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockBackend)(nil).DeleteTask), arg0, arg1)
}

// GetDeploy mocks base method
func (m *MockBackend) GetDeploy(arg0 string) (*models.Deploy, error) {
	ret := m.ctrl.Call(m, "GetDeploy", arg0)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeploy indicates an expected call of GetDeploy
func (mr *MockBackendMockRecorder) GetDeploy(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeploy", reflect.TypeOf((*MockBackend)(nil).GetDeploy), arg0)
}

// GetEnvironment mocks base method
func (m *MockBackend) GetEnvironment(arg0 string) (*models.Environment, error) {
	ret := m.ctrl.Call(m, "GetEnvironment", arg0)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnvironment indicates an expected call of GetEnvironment
func (mr *MockBackendMockRecorder) GetEnvironment(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnvironment", reflect.TypeOf((*MockBackend)(nil).GetEnvironment), arg0)
}

// GetLoadBalancer mocks base method
func (m *MockBackend) GetLoadBalancer(arg0 string) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "GetLoadBalancer", arg0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoadBalancer indicates an expected call of GetLoadBalancer
func (mr *MockBackendMockRecorder) GetLoadBalancer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoadBalancer", reflect.TypeOf((*MockBackend)(nil).GetLoadBalancer), arg0)
}

// GetService mocks base method
func (m *MockBackend) GetService(arg0, arg1 string) (*models.Service, error) {
	ret := m.ctrl.Call(m, "GetService", arg0, arg1)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetService indicates an expected call of GetService
func (mr *MockBackendMockRecorder) GetService(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetService", reflect.TypeOf((*MockBackend)(nil).GetService), arg0, arg1)
}

// GetServiceLogs mocks base method
func (m *MockBackend) GetServiceLogs(arg0, arg1, arg2, arg3 string, arg4 int) ([]*models.LogFile, error) {
	ret := m.ctrl.Call(m, "GetServiceLogs", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceLogs indicates an expected call of GetServiceLogs
func (mr *MockBackendMockRecorder) GetServiceLogs(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceLogs", reflect.TypeOf((*MockBackend)(nil).GetServiceLogs), arg0, arg1, arg2, arg3, arg4)
}

// GetTask mocks base method
func (m *MockBackend) GetTask(arg0, arg1 string) (*models.Task, error) {
	ret := m.ctrl.Call(m, "GetTask", arg0, arg1)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask
func (mr *MockBackendMockRecorder) GetTask(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockBackend)(nil).GetTask), arg0, arg1)
}

// GetTaskLogs mocks base method
func (m *MockBackend) GetTaskLogs(arg0, arg1, arg2, arg3 string, arg4 int) ([]*models.LogFile, error) {
	ret := m.ctrl.Call(m, "GetTaskLogs", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskLogs indicates an expected call of GetTaskLogs
func (mr *MockBackendMockRecorder) GetTaskLogs(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskLogs", reflect.TypeOf((*MockBackend)(nil).GetTaskLogs), arg0, arg1, arg2, arg3, arg4)
}

// ListDeploys mocks base method
func (m *MockBackend) ListDeploys() ([]*models.Deploy, error) {
	ret := m.ctrl.Call(m, "ListDeploys")
	ret0, _ := ret[0].([]*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeploys indicates an expected call of ListDeploys
func (mr *MockBackendMockRecorder) ListDeploys() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeploys", reflect.TypeOf((*MockBackend)(nil).ListDeploys))
}

// ListEnvironments mocks base method
func (m *MockBackend) ListEnvironments() ([]id.ECSEnvironmentID, error) {
	ret := m.ctrl.Call(m, "ListEnvironments")
	ret0, _ := ret[0].([]id.ECSEnvironmentID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEnvironments indicates an expected call of ListEnvironments
func (mr *MockBackendMockRecorder) ListEnvironments() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEnvironments", reflect.TypeOf((*MockBackend)(nil).ListEnvironments))
}

// ListLoadBalancers mocks base method
func (m *MockBackend) ListLoadBalancers() ([]*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "ListLoadBalancers")
	ret0, _ := ret[0].([]*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLoadBalancers indicates an expected call of ListLoadBalancers
func (mr *MockBackendMockRecorder) ListLoadBalancers() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLoadBalancers", reflect.TypeOf((*MockBackend)(nil).ListLoadBalancers))
}

// ListServices mocks base method
func (m *MockBackend) ListServices() ([]id.ECSServiceID, error) {
	ret := m.ctrl.Call(m, "ListServices")
	ret0, _ := ret[0].([]id.ECSServiceID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListServices indicates an expected call of ListServices
func (mr *MockBackendMockRecorder) ListServices() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServices", reflect.TypeOf((*MockBackend)(nil).ListServices))
}

// ListTasks mocks base method
func (m *MockBackend) ListTasks() ([]string, error) {
	ret := m.ctrl.Call(m, "ListTasks")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks
func (mr *MockBackendMockRecorder) ListTasks() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockBackend)(nil).ListTasks))
}

// ScaleService mocks base method
func (m *MockBackend) ScaleService(arg0, arg1 string, arg2 int) (*models.Service, error) {
	ret := m.ctrl.Call(m, "ScaleService", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScaleService indicates an expected call of ScaleService
func (mr *MockBackendMockRecorder) ScaleService(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleService", reflect.TypeOf((*MockBackend)(nil).ScaleService), arg0, arg1, arg2)
}

// UpdateEnvironment mocks base method
func (m *MockBackend) UpdateEnvironment(arg0 string, arg1 int) (*models.Environment, error) {
	ret := m.ctrl.Call(m, "UpdateEnvironment", arg0, arg1)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEnvironment indicates an expected call of UpdateEnvironment
func (mr *MockBackendMockRecorder) UpdateEnvironment(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEnvironment", reflect.TypeOf((*MockBackend)(nil).UpdateEnvironment), arg0, arg1)
}

// UpdateLoadBalancerHealthCheck mocks base method
func (m *MockBackend) UpdateLoadBalancerHealthCheck(arg0 string, arg1 models.HealthCheck) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "UpdateLoadBalancerHealthCheck", arg0, arg1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLoadBalancerHealthCheck indicates an expected call of UpdateLoadBalancerHealthCheck
func (mr *MockBackendMockRecorder) UpdateLoadBalancerHealthCheck(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLoadBalancerHealthCheck", reflect.TypeOf((*MockBackend)(nil).UpdateLoadBalancerHealthCheck), arg0, arg1)
}

// UpdateLoadBalancerPorts mocks base method
func (m *MockBackend) UpdateLoadBalancerPorts(arg0 string, arg1 []models.Port) (*models.LoadBalancer, error) {
	ret := m.ctrl.Call(m, "UpdateLoadBalancerPorts", arg0, arg1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLoadBalancerPorts indicates an expected call of UpdateLoadBalancerPorts
func (mr *MockBackendMockRecorder) UpdateLoadBalancerPorts(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLoadBalancerPorts", reflect.TypeOf((*MockBackend)(nil).UpdateLoadBalancerPorts), arg0, arg1)
}

// UpdateService mocks base method
func (m *MockBackend) UpdateService(arg0, arg1, arg2 string) (*models.Service, error) {
	ret := m.ctrl.Call(m, "UpdateService", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateService indicates an expected call of UpdateService
func (mr *MockBackendMockRecorder) UpdateService(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateService", reflect.TypeOf((*MockBackend)(nil).UpdateService), arg0, arg1, arg2)
}
