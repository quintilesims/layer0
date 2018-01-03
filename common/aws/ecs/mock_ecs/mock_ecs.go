// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quintilesims/layer0/common/aws/ecs (interfaces: Provider)

// Package mock_ecs is a generated GoMock package.
package mock_ecs

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ecs "github.com/quintilesims/layer0/common/aws/ecs"
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

// CreateCluster mocks base method
func (m *MockProvider) CreateCluster(arg0 string) (*ecs.Cluster, error) {
	ret := m.ctrl.Call(m, "CreateCluster", arg0)
	ret0, _ := ret[0].(*ecs.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCluster indicates an expected call of CreateCluster
func (mr *MockProviderMockRecorder) CreateCluster(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockProvider)(nil).CreateCluster), arg0)
}

// CreateService mocks base method
func (m *MockProvider) CreateService(arg0, arg1, arg2 string, arg3 int64, arg4 []*ecs.LoadBalancer, arg5 *string) (*ecs.Service, error) {
	ret := m.ctrl.Call(m, "CreateService", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*ecs.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateService indicates an expected call of CreateService
func (mr *MockProviderMockRecorder) CreateService(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateService", reflect.TypeOf((*MockProvider)(nil).CreateService), arg0, arg1, arg2, arg3, arg4, arg5)
}

// DeleteCluster mocks base method
func (m *MockProvider) DeleteCluster(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteCluster", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCluster indicates an expected call of DeleteCluster
func (mr *MockProviderMockRecorder) DeleteCluster(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockProvider)(nil).DeleteCluster), arg0)
}

// DeleteService mocks base method
func (m *MockProvider) DeleteService(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteService", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteService indicates an expected call of DeleteService
func (mr *MockProviderMockRecorder) DeleteService(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteService", reflect.TypeOf((*MockProvider)(nil).DeleteService), arg0, arg1)
}

// DeleteTaskDefinition mocks base method
func (m *MockProvider) DeleteTaskDefinition(arg0 string) error {
	ret := m.ctrl.Call(m, "DeleteTaskDefinition", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTaskDefinition indicates an expected call of DeleteTaskDefinition
func (mr *MockProviderMockRecorder) DeleteTaskDefinition(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTaskDefinition", reflect.TypeOf((*MockProvider)(nil).DeleteTaskDefinition), arg0)
}

// DescribeCluster mocks base method
func (m *MockProvider) DescribeCluster(arg0 string) (*ecs.Cluster, error) {
	ret := m.ctrl.Call(m, "DescribeCluster", arg0)
	ret0, _ := ret[0].(*ecs.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeCluster indicates an expected call of DescribeCluster
func (mr *MockProviderMockRecorder) DescribeCluster(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeCluster", reflect.TypeOf((*MockProvider)(nil).DescribeCluster), arg0)
}

// DescribeContainerInstances mocks base method
func (m *MockProvider) DescribeContainerInstances(arg0 string, arg1 []*string) ([]*ecs.ContainerInstance, error) {
	ret := m.ctrl.Call(m, "DescribeContainerInstances", arg0, arg1)
	ret0, _ := ret[0].([]*ecs.ContainerInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeContainerInstances indicates an expected call of DescribeContainerInstances
func (mr *MockProviderMockRecorder) DescribeContainerInstances(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeContainerInstances", reflect.TypeOf((*MockProvider)(nil).DescribeContainerInstances), arg0, arg1)
}

// DescribeService mocks base method
func (m *MockProvider) DescribeService(arg0, arg1 string) (*ecs.Service, error) {
	ret := m.ctrl.Call(m, "DescribeService", arg0, arg1)
	ret0, _ := ret[0].(*ecs.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeService indicates an expected call of DescribeService
func (mr *MockProviderMockRecorder) DescribeService(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeService", reflect.TypeOf((*MockProvider)(nil).DescribeService), arg0, arg1)
}

// DescribeServices mocks base method
func (m *MockProvider) DescribeServices(arg0 string, arg1 []*string) ([]*ecs.Service, error) {
	ret := m.ctrl.Call(m, "DescribeServices", arg0, arg1)
	ret0, _ := ret[0].([]*ecs.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeServices indicates an expected call of DescribeServices
func (mr *MockProviderMockRecorder) DescribeServices(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeServices", reflect.TypeOf((*MockProvider)(nil).DescribeServices), arg0, arg1)
}

// DescribeTask mocks base method
func (m *MockProvider) DescribeTask(arg0, arg1 string) (*ecs.Task, error) {
	ret := m.ctrl.Call(m, "DescribeTask", arg0, arg1)
	ret0, _ := ret[0].(*ecs.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTask indicates an expected call of DescribeTask
func (mr *MockProviderMockRecorder) DescribeTask(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTask", reflect.TypeOf((*MockProvider)(nil).DescribeTask), arg0, arg1)
}

// DescribeTaskDefinition mocks base method
func (m *MockProvider) DescribeTaskDefinition(arg0 string) (*ecs.TaskDefinition, error) {
	ret := m.ctrl.Call(m, "DescribeTaskDefinition", arg0)
	ret0, _ := ret[0].(*ecs.TaskDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTaskDefinition indicates an expected call of DescribeTaskDefinition
func (mr *MockProviderMockRecorder) DescribeTaskDefinition(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTaskDefinition", reflect.TypeOf((*MockProvider)(nil).DescribeTaskDefinition), arg0)
}

// DescribeTasks mocks base method
func (m *MockProvider) DescribeTasks(arg0 string, arg1 []*string) ([]*ecs.Task, error) {
	ret := m.ctrl.Call(m, "DescribeTasks", arg0, arg1)
	ret0, _ := ret[0].([]*ecs.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeTasks indicates an expected call of DescribeTasks
func (mr *MockProviderMockRecorder) DescribeTasks(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeTasks", reflect.TypeOf((*MockProvider)(nil).DescribeTasks), arg0, arg1)
}

// Helper_DescribeClusters mocks base method
func (m *MockProvider) Helper_DescribeClusters() ([]*ecs.Cluster, error) {
	ret := m.ctrl.Call(m, "Helper_DescribeClusters")
	ret0, _ := ret[0].([]*ecs.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Helper_DescribeClusters indicates an expected call of Helper_DescribeClusters
func (mr *MockProviderMockRecorder) Helper_DescribeClusters() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper_DescribeClusters", reflect.TypeOf((*MockProvider)(nil).Helper_DescribeClusters))
}

// Helper_DescribeServices mocks base method
func (m *MockProvider) Helper_DescribeServices(arg0 string) ([]*ecs.Service, error) {
	ret := m.ctrl.Call(m, "Helper_DescribeServices", arg0)
	ret0, _ := ret[0].([]*ecs.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Helper_DescribeServices indicates an expected call of Helper_DescribeServices
func (mr *MockProviderMockRecorder) Helper_DescribeServices(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper_DescribeServices", reflect.TypeOf((*MockProvider)(nil).Helper_DescribeServices), arg0)
}

// Helper_DescribeTaskDefinitions mocks base method
func (m *MockProvider) Helper_DescribeTaskDefinitions(arg0 string) ([]*ecs.TaskDefinition, error) {
	ret := m.ctrl.Call(m, "Helper_DescribeTaskDefinitions", arg0)
	ret0, _ := ret[0].([]*ecs.TaskDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Helper_DescribeTaskDefinitions indicates an expected call of Helper_DescribeTaskDefinitions
func (mr *MockProviderMockRecorder) Helper_DescribeTaskDefinitions(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper_DescribeTaskDefinitions", reflect.TypeOf((*MockProvider)(nil).Helper_DescribeTaskDefinitions), arg0)
}

// Helper_ListServices mocks base method
func (m *MockProvider) Helper_ListServices(arg0 string) ([]*string, error) {
	ret := m.ctrl.Call(m, "Helper_ListServices", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Helper_ListServices indicates an expected call of Helper_ListServices
func (mr *MockProviderMockRecorder) Helper_ListServices(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper_ListServices", reflect.TypeOf((*MockProvider)(nil).Helper_ListServices), arg0)
}

// Helper_ListTaskDefinitions mocks base method
func (m *MockProvider) Helper_ListTaskDefinitions(arg0 string) ([]*string, error) {
	ret := m.ctrl.Call(m, "Helper_ListTaskDefinitions", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Helper_ListTaskDefinitions indicates an expected call of Helper_ListTaskDefinitions
func (mr *MockProviderMockRecorder) Helper_ListTaskDefinitions(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper_ListTaskDefinitions", reflect.TypeOf((*MockProvider)(nil).Helper_ListTaskDefinitions), arg0)
}

// ListClusterNames mocks base method
func (m *MockProvider) ListClusterNames(arg0 string) ([]string, error) {
	ret := m.ctrl.Call(m, "ListClusterNames", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusterNames indicates an expected call of ListClusterNames
func (mr *MockProviderMockRecorder) ListClusterNames(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusterNames", reflect.TypeOf((*MockProvider)(nil).ListClusterNames), arg0)
}

// ListClusterServiceNames mocks base method
func (m *MockProvider) ListClusterServiceNames(arg0, arg1 string) ([]string, error) {
	ret := m.ctrl.Call(m, "ListClusterServiceNames", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusterServiceNames indicates an expected call of ListClusterServiceNames
func (mr *MockProviderMockRecorder) ListClusterServiceNames(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusterServiceNames", reflect.TypeOf((*MockProvider)(nil).ListClusterServiceNames), arg0, arg1)
}

// ListClusterTaskARNs mocks base method
func (m *MockProvider) ListClusterTaskARNs(arg0, arg1 string) ([]string, error) {
	ret := m.ctrl.Call(m, "ListClusterTaskARNs", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusterTaskARNs indicates an expected call of ListClusterTaskARNs
func (mr *MockProviderMockRecorder) ListClusterTaskARNs(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusterTaskARNs", reflect.TypeOf((*MockProvider)(nil).ListClusterTaskARNs), arg0, arg1)
}

// ListClusters mocks base method
func (m *MockProvider) ListClusters() ([]*string, error) {
	ret := m.ctrl.Call(m, "ListClusters")
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusters indicates an expected call of ListClusters
func (mr *MockProviderMockRecorder) ListClusters() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusters", reflect.TypeOf((*MockProvider)(nil).ListClusters))
}

// ListContainerInstances mocks base method
func (m *MockProvider) ListContainerInstances(arg0 string) ([]*string, error) {
	ret := m.ctrl.Call(m, "ListContainerInstances", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListContainerInstances indicates an expected call of ListContainerInstances
func (mr *MockProviderMockRecorder) ListContainerInstances(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListContainerInstances", reflect.TypeOf((*MockProvider)(nil).ListContainerInstances), arg0)
}

// ListServices mocks base method
func (m *MockProvider) ListServices(arg0 string) ([]*string, error) {
	ret := m.ctrl.Call(m, "ListServices", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListServices indicates an expected call of ListServices
func (mr *MockProviderMockRecorder) ListServices(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServices", reflect.TypeOf((*MockProvider)(nil).ListServices), arg0)
}

// ListTaskDefinitionFamilies mocks base method
func (m *MockProvider) ListTaskDefinitionFamilies(arg0 string, arg1 *string) ([]*string, *string, error) {
	ret := m.ctrl.Call(m, "ListTaskDefinitionFamilies", arg0, arg1)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(*string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListTaskDefinitionFamilies indicates an expected call of ListTaskDefinitionFamilies
func (mr *MockProviderMockRecorder) ListTaskDefinitionFamilies(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaskDefinitionFamilies", reflect.TypeOf((*MockProvider)(nil).ListTaskDefinitionFamilies), arg0, arg1)
}

// ListTaskDefinitionFamiliesPages mocks base method
func (m *MockProvider) ListTaskDefinitionFamiliesPages(arg0 string) ([]*string, error) {
	ret := m.ctrl.Call(m, "ListTaskDefinitionFamiliesPages", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTaskDefinitionFamiliesPages indicates an expected call of ListTaskDefinitionFamiliesPages
func (mr *MockProviderMockRecorder) ListTaskDefinitionFamiliesPages(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaskDefinitionFamiliesPages", reflect.TypeOf((*MockProvider)(nil).ListTaskDefinitionFamiliesPages), arg0)
}

// ListTaskDefinitions mocks base method
func (m *MockProvider) ListTaskDefinitions(arg0 string, arg1 *string) ([]*string, *string, error) {
	ret := m.ctrl.Call(m, "ListTaskDefinitions", arg0, arg1)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(*string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListTaskDefinitions indicates an expected call of ListTaskDefinitions
func (mr *MockProviderMockRecorder) ListTaskDefinitions(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaskDefinitions", reflect.TypeOf((*MockProvider)(nil).ListTaskDefinitions), arg0, arg1)
}

// ListTaskDefinitionsPages mocks base method
func (m *MockProvider) ListTaskDefinitionsPages(arg0 string) ([]*string, error) {
	ret := m.ctrl.Call(m, "ListTaskDefinitionsPages", arg0)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTaskDefinitionsPages indicates an expected call of ListTaskDefinitionsPages
func (mr *MockProviderMockRecorder) ListTaskDefinitionsPages(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaskDefinitionsPages", reflect.TypeOf((*MockProvider)(nil).ListTaskDefinitionsPages), arg0)
}

// ListTasks mocks base method
func (m *MockProvider) ListTasks(arg0 string, arg1, arg2, arg3, arg4 *string) ([]*string, error) {
	ret := m.ctrl.Call(m, "ListTasks", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks
func (mr *MockProviderMockRecorder) ListTasks(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockProvider)(nil).ListTasks), arg0, arg1, arg2, arg3, arg4)
}

// RegisterTaskDefinition mocks base method
func (m *MockProvider) RegisterTaskDefinition(arg0, arg1, arg2 string, arg3 []*ecs.ContainerDefinition, arg4 []*ecs.Volume, arg5 []*ecs.PlacementConstraint) (*ecs.TaskDefinition, error) {
	ret := m.ctrl.Call(m, "RegisterTaskDefinition", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*ecs.TaskDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterTaskDefinition indicates an expected call of RegisterTaskDefinition
func (mr *MockProviderMockRecorder) RegisterTaskDefinition(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTaskDefinition", reflect.TypeOf((*MockProvider)(nil).RegisterTaskDefinition), arg0, arg1, arg2, arg3, arg4, arg5)
}

// RunTask mocks base method
func (m *MockProvider) RunTask(arg0, arg1, arg2 string, arg3 []*ecs.ContainerOverride) (*ecs.Task, error) {
	ret := m.ctrl.Call(m, "RunTask", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*ecs.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunTask indicates an expected call of RunTask
func (mr *MockProviderMockRecorder) RunTask(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunTask", reflect.TypeOf((*MockProvider)(nil).RunTask), arg0, arg1, arg2, arg3)
}

// StartTask mocks base method
func (m *MockProvider) StartTask(arg0, arg1 string, arg2 *ecs.TaskOverride, arg3 []*string, arg4 *string) error {
	ret := m.ctrl.Call(m, "StartTask", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartTask indicates an expected call of StartTask
func (mr *MockProviderMockRecorder) StartTask(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTask", reflect.TypeOf((*MockProvider)(nil).StartTask), arg0, arg1, arg2, arg3, arg4)
}

// StopTask mocks base method
func (m *MockProvider) StopTask(arg0, arg1, arg2 string) error {
	ret := m.ctrl.Call(m, "StopTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopTask indicates an expected call of StopTask
func (mr *MockProviderMockRecorder) StopTask(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopTask", reflect.TypeOf((*MockProvider)(nil).StopTask), arg0, arg1, arg2)
}

// UpdateService mocks base method
func (m *MockProvider) UpdateService(arg0, arg1 string, arg2 *string, arg3 *int64) error {
	ret := m.ctrl.Call(m, "UpdateService", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateService indicates an expected call of UpdateService
func (mr *MockProviderMockRecorder) UpdateService(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateService", reflect.TypeOf((*MockProvider)(nil).UpdateService), arg0, arg1, arg2, arg3)
}
