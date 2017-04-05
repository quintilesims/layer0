// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/api/backend (interfaces: Backend)

package mock_backend

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
)

// Mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *_MockBackendRecorder
}

// Recorder for MockBackend (not exported)
type _MockBackendRecorder struct {
	mock *MockBackend
}

func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &_MockBackendRecorder{mock}
	return mock
}

func (_m *MockBackend) EXPECT() *_MockBackendRecorder {
	return _m.recorder
}

func (_m *MockBackend) CreateDeploy(_param0 string, _param1 []byte) (*models.Deploy, error) {
	ret := _m.ctrl.Call(_m, "CreateDeploy", _param0, _param1)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) CreateDeploy(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateDeploy", arg0, arg1)
}

func (_m *MockBackend) CreateEnvironment(_param0 string, _param1 string, _param2 string, _param3 string, _param4 int, _param5 []byte) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "CreateEnvironment", _param0, _param1, _param2, _param3, _param4, _param5)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) CreateEnvironment(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateEnvironment", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockBackend) CreateEnvironmentLink(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "CreateEnvironmentLink", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) CreateEnvironmentLink(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateEnvironmentLink", arg0, arg1)
}

func (_m *MockBackend) CreateLoadBalancer(_param0 string, _param1 string, _param2 bool, _param3 []models.Port, _param4 models.HealthCheck) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "CreateLoadBalancer", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) CreateLoadBalancer(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLoadBalancer", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockBackend) CreateService(_param0 string, _param1 string, _param2 string, _param3 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "CreateService", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) CreateService(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateService", arg0, arg1, arg2, arg3)
}

func (_m *MockBackend) CreateTask(_param0 string, _param1 string, _param2 string, _param3 []models.ContainerOverride) (*models.Task, error) {
	ret := _m.ctrl.Call(_m, "CreateTask", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) CreateTask(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateTask", arg0, arg1, arg2, arg3)
}

func (_m *MockBackend) DeleteDeploy(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteDeploy", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) DeleteDeploy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteDeploy", arg0)
}

func (_m *MockBackend) DeleteEnvironment(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteEnvironment", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) DeleteEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteEnvironment", arg0)
}

func (_m *MockBackend) DeleteEnvironmentLink(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteEnvironmentLink", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) DeleteEnvironmentLink(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteEnvironmentLink", arg0, arg1)
}

func (_m *MockBackend) DeleteLoadBalancer(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteLoadBalancer", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) DeleteLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteLoadBalancer", arg0)
}

func (_m *MockBackend) DeleteService(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteService", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) DeleteService(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteService", arg0, arg1)
}

func (_m *MockBackend) DeleteTask(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteTask", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBackendRecorder) DeleteTask(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteTask", arg0, arg1)
}

func (_m *MockBackend) GetDeploy(_param0 string) (*models.Deploy, error) {
	ret := _m.ctrl.Call(_m, "GetDeploy", _param0)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetDeploy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDeploy", arg0)
}

func (_m *MockBackend) GetEnvironment(_param0 string) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "GetEnvironment", _param0)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnvironment", arg0)
}

func (_m *MockBackend) GetLoadBalancer(_param0 string) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "GetLoadBalancer", _param0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLoadBalancer", arg0)
}

func (_m *MockBackend) GetService(_param0 string, _param1 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "GetService", _param0, _param1)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetService(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetService", arg0, arg1)
}

func (_m *MockBackend) GetServiceLogs(_param0 string, _param1 string, _param2 int) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "GetServiceLogs", _param0, _param1, _param2)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetServiceLogs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServiceLogs", arg0, arg1, arg2)
}

func (_m *MockBackend) GetTask(_param0 string, _param1 string) (*models.Task, error) {
	ret := _m.ctrl.Call(_m, "GetTask", _param0, _param1)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetTask(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTask", arg0, arg1)
}

func (_m *MockBackend) GetTaskLogs(_param0 string, _param1 string, _param2 int) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "GetTaskLogs", _param0, _param1, _param2)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetTaskLogs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTaskLogs", arg0, arg1, arg2)
}

func (_m *MockBackend) ListDeploys() ([]*models.Deploy, error) {
	ret := _m.ctrl.Call(_m, "ListDeploys")
	ret0, _ := ret[0].([]*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) ListDeploys() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListDeploys")
}

func (_m *MockBackend) ListEnvironments() ([]*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "ListEnvironments")
	ret0, _ := ret[0].([]*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) ListEnvironments() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListEnvironments")
}

func (_m *MockBackend) ListLoadBalancers() ([]*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "ListLoadBalancers")
	ret0, _ := ret[0].([]*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) ListLoadBalancers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListLoadBalancers")
}

func (_m *MockBackend) ListServices() ([]*models.Service, error) {
	ret := _m.ctrl.Call(_m, "ListServices")
	ret0, _ := ret[0].([]*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) ListServices() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListServices")
}

func (_m *MockBackend) ListTasks() ([]*models.Task, error) {
	ret := _m.ctrl.Call(_m, "ListTasks")
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) ListTasks() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListTasks")
}

func (_m *MockBackend) ScaleService(_param0 string, _param1 string, _param2 int) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "ScaleService", _param0, _param1, _param2)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) ScaleService(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ScaleService", arg0, arg1, arg2)
}

func (_m *MockBackend) UpdateEnvironment(_param0 string, _param1 int) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "UpdateEnvironment", _param0, _param1)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) UpdateEnvironment(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateEnvironment", arg0, arg1)
}

func (_m *MockBackend) UpdateLoadBalancerHealthCheck(_param0 string, _param1 models.HealthCheck) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "UpdateLoadBalancerHealthCheck", _param0, _param1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) UpdateLoadBalancerHealthCheck(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateLoadBalancerHealthCheck", arg0, arg1)
}

func (_m *MockBackend) UpdateLoadBalancerPorts(_param0 string, _param1 []models.Port) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "UpdateLoadBalancerPorts", _param0, _param1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) UpdateLoadBalancerPorts(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateLoadBalancerPorts", arg0, arg1)
}

func (_m *MockBackend) UpdateService(_param0 string, _param1 string, _param2 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "UpdateService", _param0, _param1, _param2)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) UpdateService(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateService", arg0, arg1, arg2)
}
