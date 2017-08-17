// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/cli/client (interfaces: Client)

package mock_client

import (
	time "time"

	gomock "github.com/golang/mock/gomock"
	models "github.com/quintilesims/layer0/common/models"
)

// Mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *_MockClientRecorder
}

// Recorder for MockClient (not exported)
type _MockClientRecorder struct {
	mock *MockClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &_MockClientRecorder{mock}
	return mock
}

func (_m *MockClient) EXPECT() *_MockClientRecorder {
	return _m.recorder
}

func (_m *MockClient) CreateDeploy(_param0 string, _param1 []byte) (*models.Deploy, error) {
	ret := _m.ctrl.Call(_m, "CreateDeploy", _param0, _param1)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateDeploy(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateDeploy", arg0, arg1)
}

func (_m *MockClient) CreateEnvironment(_param0 string, _param1 string, _param2 int, _param3 []byte, _param4 string, _param5 string) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "CreateEnvironment", _param0, _param1, _param2, _param3, _param4, _param5)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateEnvironment(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateEnvironment", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockClient) CreateLink(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "CreateLink", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) CreateLink(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLink", arg0, arg1)
}

func (_m *MockClient) CreateLoadBalancer(_param0 string, _param1 string, _param2 models.HealthCheck, _param3 []models.Port, _param4 bool) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "CreateLoadBalancer", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateLoadBalancer(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLoadBalancer", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockClient) CreateService(_param0 string, _param1 string, _param2 string, _param3 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "CreateService", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateService(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateService", arg0, arg1, arg2, arg3)
}

func (_m *MockClient) CreateTask(_param0 string, _param1 string, _param2 string, _param3 int, _param4 []models.ContainerOverride) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateTask", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateTask(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateTask", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockClient) Delete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockClient) DeleteDeploy(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteDeploy", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DeleteDeploy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteDeploy", arg0)
}

func (_m *MockClient) DeleteEnvironment(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "DeleteEnvironment", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DeleteEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteEnvironment", arg0)
}

func (_m *MockClient) DeleteLink(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteLink", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DeleteLink(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteLink", arg0, arg1)
}

func (_m *MockClient) DeleteLoadBalancer(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "DeleteLoadBalancer", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DeleteLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteLoadBalancer", arg0)
}

func (_m *MockClient) DeleteService(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "DeleteService", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DeleteService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteService", arg0)
}

func (_m *MockClient) DeleteTask(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteTask", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DeleteTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteTask", arg0)
}

func (_m *MockClient) GetConfig() (*models.APIConfig, error) {
	ret := _m.ctrl.Call(_m, "GetConfig")
	ret0, _ := ret[0].(*models.APIConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetConfig() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfig")
}

func (_m *MockClient) GetDeploy(_param0 string) (*models.Deploy, error) {
	ret := _m.ctrl.Call(_m, "GetDeploy", _param0)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetDeploy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDeploy", arg0)
}

func (_m *MockClient) GetEnvironment(_param0 string) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "GetEnvironment", _param0)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnvironment", arg0)
}

func (_m *MockClient) GetJob(_param0 string) (*models.Job, error) {
	ret := _m.ctrl.Call(_m, "GetJob", _param0)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetJob(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetJob", arg0)
}

func (_m *MockClient) GetLoadBalancer(_param0 string) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "GetLoadBalancer", _param0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLoadBalancer", arg0)
}

func (_m *MockClient) GetService(_param0 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "GetService", _param0)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetService", arg0)
}

func (_m *MockClient) GetServiceLogs(_param0 string, _param1 string, _param2 string, _param3 int) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "GetServiceLogs", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetServiceLogs(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServiceLogs", arg0, arg1, arg2, arg3)
}

func (_m *MockClient) GetTask(_param0 string) (*models.Task, error) {
	ret := _m.ctrl.Call(_m, "GetTask", _param0)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTask", arg0)
}

func (_m *MockClient) GetTaskLogs(_param0 string, _param1 string, _param2 string, _param3 int) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "GetTaskLogs", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetTaskLogs(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetTaskLogs", arg0, arg1, arg2, arg3)
}

func (_m *MockClient) GetVersion() (string, error) {
	ret := _m.ctrl.Call(_m, "GetVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetVersion() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVersion")
}

func (_m *MockClient) ListDeploys() ([]*models.DeploySummary, error) {
	ret := _m.ctrl.Call(_m, "ListDeploys")
	ret0, _ := ret[0].([]*models.DeploySummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListDeploys() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListDeploys")
}

func (_m *MockClient) ListEnvironments() ([]*models.EnvironmentSummary, error) {
	ret := _m.ctrl.Call(_m, "ListEnvironments")
	ret0, _ := ret[0].([]*models.EnvironmentSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListEnvironments() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListEnvironments")
}

func (_m *MockClient) ListJobs() ([]*models.Job, error) {
	ret := _m.ctrl.Call(_m, "ListJobs")
	ret0, _ := ret[0].([]*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListJobs() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListJobs")
}

func (_m *MockClient) ListLoadBalancers() ([]*models.LoadBalancerSummary, error) {
	ret := _m.ctrl.Call(_m, "ListLoadBalancers")
	ret0, _ := ret[0].([]*models.LoadBalancerSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListLoadBalancers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListLoadBalancers")
}

func (_m *MockClient) ListServices() ([]*models.ServiceSummary, error) {
	ret := _m.ctrl.Call(_m, "ListServices")
	ret0, _ := ret[0].([]*models.ServiceSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListServices() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListServices")
}

func (_m *MockClient) ListTasks() ([]*models.TaskSummary, error) {
	ret := _m.ctrl.Call(_m, "ListTasks")
	ret0, _ := ret[0].([]*models.TaskSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListTasks() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListTasks")
}

func (_m *MockClient) RunScaler(_param0 string) (*models.ScalerRunInfo, error) {
	ret := _m.ctrl.Call(_m, "RunScaler", _param0)
	ret0, _ := ret[0].(*models.ScalerRunInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) RunScaler(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RunScaler", arg0)
}

func (_m *MockClient) ScaleService(_param0 string, _param1 int) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "ScaleService", _param0, _param1)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ScaleService(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ScaleService", arg0, arg1)
}

func (_m *MockClient) SelectByQuery(_param0 map[string]string) ([]*models.EntityWithTags, error) {
	ret := _m.ctrl.Call(_m, "SelectByQuery", _param0)
	ret0, _ := ret[0].([]*models.EntityWithTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) SelectByQuery(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SelectByQuery", arg0)
}

func (_m *MockClient) UpdateEnvironment(_param0 string, _param1 int) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "UpdateEnvironment", _param0, _param1)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateEnvironment(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateEnvironment", arg0, arg1)
}

func (_m *MockClient) UpdateLoadBalancerHealthCheck(_param0 string, _param1 models.HealthCheck) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "UpdateLoadBalancerHealthCheck", _param0, _param1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateLoadBalancerHealthCheck(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateLoadBalancerHealthCheck", arg0, arg1)
}

func (_m *MockClient) UpdateLoadBalancerPorts(_param0 string, _param1 []models.Port) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "UpdateLoadBalancerPorts", _param0, _param1)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateLoadBalancerPorts(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateLoadBalancerPorts", arg0, arg1)
}

func (_m *MockClient) UpdateSQL() error {
	ret := _m.ctrl.Call(_m, "UpdateSQL")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) UpdateSQL() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateSQL")
}

func (_m *MockClient) UpdateService(_param0 string, _param1 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "UpdateService", _param0, _param1)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateService(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateService", arg0, arg1)
}

func (_m *MockClient) WaitForDeployment(_param0 string, _param1 time.Duration) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "WaitForDeployment", _param0, _param1)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) WaitForDeployment(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WaitForDeployment", arg0, arg1)
}

func (_m *MockClient) WaitForJob(_param0 string, _param1 time.Duration) error {
	ret := _m.ctrl.Call(_m, "WaitForJob", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) WaitForJob(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WaitForJob", arg0, arg1)
}