// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/quintilesims/layer0/client (interfaces: Client)

package mock_client

import (
	url "net/url"

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

func (_m *MockClient) CreateDeploy(_param0 models.CreateDeployRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateDeploy", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateDeploy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateDeploy", arg0)
}

func (_m *MockClient) CreateEnvironment(_param0 models.CreateEnvironmentRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateEnvironment", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateEnvironment", arg0)
}

func (_m *MockClient) CreateLink(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "CreateLink", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) CreateLink(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLink", arg0, arg1)
}

func (_m *MockClient) CreateLoadBalancer(_param0 models.CreateLoadBalancerRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateLoadBalancer", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateLoadBalancer", arg0)
}

func (_m *MockClient) CreateService(_param0 models.CreateServiceRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateService", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateService", arg0)
}

func (_m *MockClient) CreateTask(_param0 models.CreateTaskRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "CreateTask", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) CreateTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateTask", arg0)
}

func (_m *MockClient) DeleteDeploy(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "DeleteDeploy", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
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

func (_m *MockClient) DeleteJob(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteJob", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DeleteJob(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteJob", arg0)
}

func (_m *MockClient) DeleteLink(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteLink", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DeleteJob(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteJob", arg0)
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

func (_m *MockClient) DeleteTask(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "DeleteTask", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DeleteTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteTask", arg0)
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

func (_m *MockClient) ListTags(_param0 url.Values) (models.Tags, error) {
	ret := _m.ctrl.Call(_m, "ListTags", _param0)
	ret0, _ := ret[0].(models.Tags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListTags(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListTags", arg0)
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

func (_m *MockClient) ReadConfig() (*models.APIConfig, error) {
	ret := _m.ctrl.Call(_m, "ReadConfig")
	ret0, _ := ret[0].(*models.APIConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadConfig() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadConfig")
}

func (_m *MockClient) ReadDeploy(_param0 string) (*models.Deploy, error) {
	ret := _m.ctrl.Call(_m, "ReadDeploy", _param0)
	ret0, _ := ret[0].(*models.Deploy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadDeploy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadDeploy", arg0)
}

func (_m *MockClient) ReadEnvironment(_param0 string) (*models.Environment, error) {
	ret := _m.ctrl.Call(_m, "ReadEnvironment", _param0)
	ret0, _ := ret[0].(*models.Environment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadEnvironment", arg0)
}

func (_m *MockClient) ReadJob(_param0 string) (*models.Job, error) {
	ret := _m.ctrl.Call(_m, "ReadJob", _param0)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadJob(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadJob", arg0)
}

func (_m *MockClient) ReadLoadBalancer(_param0 string) (*models.LoadBalancer, error) {
	ret := _m.ctrl.Call(_m, "ReadLoadBalancer", _param0)
	ret0, _ := ret[0].(*models.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadLoadBalancer", arg0)
}

func (_m *MockClient) ReadService(_param0 string) (*models.Service, error) {
	ret := _m.ctrl.Call(_m, "ReadService", _param0)
	ret0, _ := ret[0].(*models.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadService", arg0)
}

func (_m *MockClient) ReadServiceLogs(_param0 string, _param1 url.Values) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "ReadServiceLogs", _param0, _param1)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadServiceLogs(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadServiceLogs", arg0, arg1)
}

func (_m *MockClient) ReadTask(_param0 string) (*models.Task, error) {
	ret := _m.ctrl.Call(_m, "ReadTask", _param0)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadTask(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadTask", arg0)
}

func (_m *MockClient) ReadTaskLogs(_param0 string, _param1 url.Values) ([]*models.LogFile, error) {
	ret := _m.ctrl.Call(_m, "ReadTaskLogs", _param0, _param1)
	ret0, _ := ret[0].([]*models.LogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ReadTaskLogs(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadTaskLogs", arg0, arg1)
}

func (_m *MockClient) UpdateEnvironment(_param0 models.UpdateEnvironmentRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "UpdateEnvironment", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateEnvironment(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateEnvironment", arg0)
}

func (_m *MockClient) UpdateLoadBalancer(_param0 models.UpdateLoadBalancerRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "UpdateLoadBalancer", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateLoadBalancer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateLoadBalancer", arg0)
}

func (_m *MockClient) UpdateService(_param0 models.UpdateServiceRequest) (string, error) {
	ret := _m.ctrl.Call(_m, "UpdateService", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) UpdateService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateService", arg0)
}
