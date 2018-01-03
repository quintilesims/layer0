package ecs

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/aws/provider"
)

const MAX_DESCRIBE_SERVICE_IDS = 10

type Provider interface {
	CreateCluster(clusterName string) (*Cluster, error)
	CreateService(cluster, serviceName, taskDefinition string, desiredCount int64, loadBalancers []*LoadBalancer, loadBalancerRole *string) (*Service, error)

	DeleteCluster(cluster string) error
	DeleteService(cluster, service string) error
	DeleteTaskDefinition(familyAndRevision string) error

	DescribeContainerInstances(clusterName string, instances []*string) ([]*ContainerInstance, error)
	DescribeCluster(clusterName string) (*Cluster, error)

	Helper_DescribeClusters() ([]*Cluster, error)
	DescribeService(cluster, service string) (*Service, error)

	DescribeServices(cluster string, service []*string) ([]*Service, error)
	Helper_DescribeServices(prefix string) ([]*Service, error)

	DescribeTaskDefinition(familyAndRevision string) (*TaskDefinition, error)
	Helper_DescribeTaskDefinitions(prefix string) ([]*TaskDefinition, error)

	DescribeTask(cluster string, taskARN string) (*Task, error)
	DescribeTasks(clusterName string, taskARNs []*string) ([]*Task, error)

	ListClusters() ([]*string, error)
	ListClusterNames(prefix string) ([]string, error)
	ListContainerInstances(clusterName string) ([]*string, error)
	ListServices(clusterName string) ([]*string, error)
	Helper_ListServices(prefix string) ([]*string, error)

	ListClusterTaskARNs(clusterName, startedBy string) ([]string, error)
	ListClusterServiceNames(clusterName, prefix string) ([]string, error)
	ListTasks(clusterName string, serviceName, desiredStatus, startedBy, containerInstance *string) ([]*string, error)

	ListTaskDefinitions(familyName string, nextToken *string) ([]*string, *string, error)
	Helper_ListTaskDefinitions(prefix string) ([]*string, error)
	ListTaskDefinitionsPages(familyName string) ([]*string, error)

	ListTaskDefinitionFamilies(prefix string, nextToken *string) ([]*string, *string, error)
	ListTaskDefinitionFamiliesPages(prefix string) ([]*string, error)

	RegisterTaskDefinition(family string, roleARN string, networkMode string, containerDefinitions []*ContainerDefinition, volumes []*Volume, placementConstraints []*PlacementConstraint) (*TaskDefinition, error)
	RunTask(clusterName, taskDefinition, startedBy string, overrides []*ContainerOverride) (*Task, error)
	StartTask(cluster, taskDefinition string, overrides *TaskOverride, containerInstanceIDs []*string, startedBy *string) error
	StopTask(clusterName, taskARN, reason string) error

	UpdateService(cluster, service string, taskDefinition *string, desiredCount *int64) error
}

type ECS struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (ECSInternal, error)
}

type ECSInternal interface {
	CreateCluster(input *ecs.CreateClusterInput) (*ecs.CreateClusterOutput, error)
	CreateService(input *ecs.CreateServiceInput) (output *ecs.CreateServiceOutput, err error)
	DeleteCluster(input *ecs.DeleteClusterInput) (*ecs.DeleteClusterOutput, error)
	DeleteService(input *ecs.DeleteServiceInput) (*ecs.DeleteServiceOutput, error)
	DeregisterTaskDefinition(input *ecs.DeregisterTaskDefinitionInput) (*ecs.DeregisterTaskDefinitionOutput, error)
	DescribeClusters(input *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error)
	DescribeContainerInstances(input *ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error)
	DescribeServices(input *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error)
	DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error)

	DescribeTasks(input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error)

	ListTasksPages(input *ecs.ListTasksInput, fn func(*ecs.ListTasksOutput, bool) bool) error
	ListClustersPages(input *ecs.ListClustersInput, fn func(*ecs.ListClustersOutput, bool) bool) error
	ListContainerInstancesPages(input *ecs.ListContainerInstancesInput, fn func(*ecs.ListContainerInstancesOutput, bool) bool) error
	ListServicesPages(input *ecs.ListServicesInput, fn func(*ecs.ListServicesOutput, bool) bool) error

	ListTaskDefinitions(input *ecs.ListTaskDefinitionsInput) (*ecs.ListTaskDefinitionsOutput, error)
	ListTaskDefinitionsPages(input *ecs.ListTaskDefinitionsInput, fn func(p *ecs.ListTaskDefinitionsOutput, lastPage bool) (shouldContinue bool)) error

	ListTaskDefinitionFamilies(input *ecs.ListTaskDefinitionFamiliesInput) (*ecs.ListTaskDefinitionFamiliesOutput, error)
	ListTaskDefinitionFamiliesPages(input *ecs.ListTaskDefinitionFamiliesInput, fn func(p *ecs.ListTaskDefinitionFamiliesOutput, lastPage bool) (shouldContinue bool)) error

	RegisterTaskDefinition(input *ecs.RegisterTaskDefinitionInput) (output *ecs.RegisterTaskDefinitionOutput, err error)
	RunTask(input *ecs.RunTaskInput) (output *ecs.RunTaskOutput, err error)
	StartTask(input *ecs.StartTaskInput) (output *ecs.StartTaskOutput, err error)
	StopTask(input *ecs.StopTaskInput) (output *ecs.StopTaskOutput, err error)
	UpdateService(input *ecs.UpdateServiceInput) (output *ecs.UpdateServiceOutput, err error)
}

type ContainerInstance struct {
	*ecs.ContainerInstance
}

type Volume struct {
	*ecs.Volume
}

type Cluster struct {
	*ecs.Cluster
}

func NewCluster(name string) *Cluster {
	return &Cluster{&ecs.Cluster{
		ClusterName: &name,
	}}
}

type TaskOverride struct {
	*ecs.TaskOverride
}

type FailedTask struct {
	*ecs.Failure
}

type ContainerOverride struct {
	*ecs.ContainerOverride
}

type LoadBalancer struct {
	*ecs.LoadBalancer
}

type Service struct {
	*ecs.Service
}

func NewService(clusterARN, name string) *Service {
	return &Service{
		&ecs.Service{
			ClusterArn:   aws.String(clusterARN),
			ServiceName:  aws.String(name),
			DesiredCount: aws.Int64(0),
			RunningCount: aws.Int64(0),
			PendingCount: aws.Int64(0),
		},
	}
}

type TaskDefinition struct {
	*ecs.TaskDefinition
}

type Task struct {
	*ecs.Task
}

type Container struct {
	*ecs.Container
}

type PortMapping struct {
	*ecs.PortMapping
}

type ContainerDefinition struct {
	*ecs.ContainerDefinition
}

type PlacementConstraint struct {
	*ecs.TaskDefinitionPlacementConstraint
}

func NewLoadBalancer(containerName string, containerPort int64, loadBalancerName string) *LoadBalancer {
	return &LoadBalancer{
		&ecs.LoadBalancer{
			ContainerName:    &containerName,
			ContainerPort:    aws.Int64(containerPort),
			LoadBalancerName: &loadBalancerName,
		},
	}
}

func NewContainerInstance(agentConnected bool, cpuRegistered, memoryRegistered int, portsRegistered, udpPortsRegistered []*string) *ContainerInstance {
	return &ContainerInstance{
		&ecs.ContainerInstance{
			AgentConnected: &agentConnected,
			RemainingResources: []*ecs.Resource{
				newEcsIntResource("CPU", cpuRegistered),
				newEcsIntResource("MEMORY", memoryRegistered),
				newEcsStringSetResource("PORTS", portsRegistered),
				newEcsStringSetResource("PORTS_UDP", udpPortsRegistered),
			},
			PendingTasksCount: aws.Int64(0),
			RunningTasksCount: aws.Int64(0),
		},
	}
}

func NewDefaultContainerInstance() *ContainerInstance {
	return &ContainerInstance{
		&ecs.ContainerInstance{},
	}
}

func NewPortMapping(containerPort int64, hostPort *int64, protocol string) *PortMapping {
	return &PortMapping{
		&ecs.PortMapping{
			ContainerPort: aws.Int64(containerPort),
			HostPort:      hostPort,
			Protocol:      aws.String(protocol),
		},
	}
}

func NewContainerDefinition(name, imageName string, entryPoint []*string, command string, cpu, memory int64, essential bool, portMappings []*PortMapping) *ContainerDefinition {

	def := ecs.ContainerDefinition{
		Image: aws.String(imageName),
		Command: []*string{
			aws.String(command),
		},
		Cpu:        aws.Int64(cpu),
		Memory:     aws.Int64(memory),
		Name:       aws.String(name),
		Essential:  aws.Bool(essential),
		EntryPoint: entryPoint,
	}
	if portMappings != nil {
		def.PortMappings = []*ecs.PortMapping{}
		for _, mapping := range portMappings {
			def.PortMappings = append(
				def.PortMappings,
				mapping.PortMapping)
		}
	}
	return &ContainerDefinition{&def}
}

func NewTaskDefinition() *TaskDefinition {
	return &TaskDefinition{
		&ecs.TaskDefinition{},
	}
}

func NewTask(clusterARN, taskID, deployARN string) *Task {
	return &Task{
		&ecs.Task{
			ClusterArn:        aws.String(clusterARN),
			TaskDefinitionArn: aws.String(deployARN),
			StartedBy:         aws.String(taskID),
			LastStatus:        aws.String("RUNNING"),
		},
	}
}

func (this *TaskDefinition) AddContainerDefinition(def *ContainerDefinition) {
	if this.ContainerDefinitions == nil {
		this.ContainerDefinitions = []*ecs.ContainerDefinition{}
	}

	this.ContainerDefinitions = append(this.ContainerDefinitions, def.ContainerDefinition)
}

func NewContainerOverride(name string, envVars map[string]string) *ContainerOverride {
	environment := []*ecs.KeyValuePair{}
	for k, v := range envVars {
		name := k
		value := v
		environment = append(environment, &ecs.KeyValuePair{
			Name:  &name,
			Value: &value,
		})
	}
	return &ContainerOverride{
		&ecs.ContainerOverride{
			Name:        &name,
			Environment: environment,
		},
	}
}

func newEcsIntResource(name string, value int) *ecs.Resource {
	return &ecs.Resource{
		Name:         aws.String(name),
		Type:         aws.String("INTEGER"),
		IntegerValue: aws.Int64(int64(value)),
	}
}

func newEcsStringSetResource(name string, value []*string) *ecs.Resource {
	return &ecs.Resource{
		Name:           aws.String(name),
		Type:           aws.String("STRINGSET"),
		StringSetValue: value,
	}
}

func NewECS(credProvider provider.CredProvider, region string) (Provider, error) {
	ecs := ECS{
		credProvider,
		region,
		func() (ECSInternal, error) {
			return Connect(credProvider, region)
		},
	}

	if _, err := ecs.Connect(); err != nil {
		return nil, err
	}

	return &ecs, nil
}

func (t *Task) GetContainers() []*Container {
	ret := []*Container{}
	for _, c := range t.Containers {
		ret = append(ret, &Container{c})
	}
	return ret
}

func (d *TaskDefinition) GetContainerDefinitions() []*ContainerDefinition {
	ret := []*ContainerDefinition{}
	for _, c := range d.ContainerDefinitions {
		ret = append(ret, &ContainerDefinition{c})
	}
	return ret
}

func Connect(credProvider provider.CredProvider, region string) (ECSInternal, error) {
	connection, err := provider.GetECSConnection(credProvider, region)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (this *ECS) RegisterTaskDefinition(family string, roleARN string, networkMode string, containerDefinitions []*ContainerDefinition, volumes []*Volume, placementConstraints []*PlacementConstraint) (*TaskDefinition, error) {
	awsContainers := make([]*ecs.ContainerDefinition, 0)
	for _, val := range containerDefinitions {
		awsContainers = append(awsContainers, val.ContainerDefinition)
	}

	awsVolumes := make([]*ecs.Volume, 0)
	for _, val := range volumes {
		awsVolumes = append(awsVolumes, val.Volume)
	}

	var roleARNp *string
	if roleARN != "" {
		roleARNp = aws.String(roleARN)
	}

	var networkModep *string
	if networkMode != "" {
		networkModep = aws.String(networkMode)
	}

	awsPlacementConstraints := make([]*ecs.TaskDefinitionPlacementConstraint, 0)
	for _, val := range placementConstraints {
		awsPlacementConstraints = append(awsPlacementConstraints, val.TaskDefinitionPlacementConstraint)
	}

	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: awsContainers,
		Family:               aws.String(family),
		Volumes:              awsVolumes,
		TaskRoleArn:          roleARNp,
		NetworkMode:          networkModep,
		PlacementConstraints: awsPlacementConstraints,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}
	return &TaskDefinition{output.TaskDefinition}, err
}

func (this *ECS) RunTask(
	clusterName string,
	taskDefinition string,
	startedBy string,
	overrides []*ContainerOverride,
) (*Task, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	placementStrategy := []*ecs.PlacementStrategy{{
		Type:  aws.String(ecs.PlacementStrategyTypeBinpack),
		Field: aws.String("memory"),
	}}

	input := &ecs.RunTaskInput{
		Cluster:           aws.String(clusterName),
		TaskDefinition:    aws.String(taskDefinition),
		Count:             aws.Int64(1),
		StartedBy:         aws.String(startedBy),
		PlacementStrategy: placementStrategy,
	}

	if overrides != nil {
		cos := make([]*ecs.ContainerOverride, len(overrides))
		for i, o := range overrides {
			cos[i] = o.ContainerOverride
		}

		input.Overrides = &ecs.TaskOverride{ContainerOverrides: cos}
	}

	result, err := connection.RunTask(input)
	if err != nil {
		return nil, err
	}

	if len(result.Failures) > 0 {
		return nil, fmt.Errorf("Failed to start task: %s", aws.StringValue(result.Failures[0].Reason))
	}

	return &Task{result.Tasks[0]}, nil
}

func (this *ECS) StartTask(cluster, taskDefinition string, overrides *TaskOverride, containerInstanceIDs []*string, startedBy *string) (err error) {
	input := &ecs.StartTaskInput{
		Cluster:            aws.String(cluster),
		TaskDefinition:     aws.String(taskDefinition),
		ContainerInstances: containerInstanceIDs,
		StartedBy:          startedBy,
	}

	if overrides != nil {
		input.Overrides = overrides.TaskOverride
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.StartTask(input)
	return err
}

func (this *ECS) StopTask(clusterName, taskARN, reason string) error {
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	input := &ecs.StopTaskInput{
		Cluster: aws.String(clusterName),
		Task:    aws.String(taskARN),
		Reason:  aws.String(reason),
	}

	if _, err := connection.StopTask(input); err != nil {
		return err
	}

	return nil
}

func (this *ECS) CreateService(cluster, serviceName, taskDefinition string, desiredCount int64, loadBalancers []*LoadBalancer, loadBalancerRole *string) (*Service, error) {

	awsLoadBalancers := make([]*ecs.LoadBalancer, 0)
	for _, val := range loadBalancers {
		awsLoadBalancers = append(awsLoadBalancers, val.LoadBalancer)
	}

	input := &ecs.CreateServiceInput{
		Cluster:        aws.String(cluster),
		ServiceName:    aws.String(serviceName),
		TaskDefinition: aws.String(taskDefinition),
		DesiredCount:   aws.Int64(desiredCount),
		LoadBalancers:  awsLoadBalancers,
		Role:           loadBalancerRole,
		PlacementStrategy: []*ecs.PlacementStrategy{
			{
				Type:  aws.String(ecs.PlacementStrategyTypeBinpack),
				Field: aws.String("memory"),
			},
		},
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.CreateService(input)
	if err != nil {
		return nil, err
	}
	return &Service{output.Service}, nil
}

func (this *ECS) UpdateService(cluster, service string, taskDefinition *string, desiredCount *int64) error {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		DesiredCount:   desiredCount,
		Service:        aws.String(service),
		TaskDefinition: taskDefinition,
	}
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.UpdateService(input)
	return err
}

func (this *ECS) DeleteService(cluster, service string) error {
	input := &ecs.DeleteServiceInput{
		Cluster: aws.String(cluster),
		Service: aws.String(service),
	}
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteService(input)
	return err
}

func (this *ECS) DescribeService(cluster, service string) (*Service, error) {
	services, err := this.DescribeServices(cluster, []*string{&service})
	if err != nil {
		return nil, err
	}

	if len(services) > 0 {
		return services[0], nil
	}

	return nil, fmt.Errorf("Service Not Found")
}

func (this *ECS) DescribeServices(cluster string, serviceIDs []*string) ([]*Service, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	services := []*Service{}
	for i := len(serviceIDs); i > 0; i = len(serviceIDs) {
		if i > MAX_DESCRIBE_SERVICE_IDS {
			i = MAX_DESCRIBE_SERVICE_IDS
		}

		input := &ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster),
			Services: serviceIDs[:i],
		}

		output, err := connection.DescribeServices(input)
		if err != nil {
			return nil, err
		}

		for _, svc := range output.Services {
			services = append(services, &Service{svc})
		}

		serviceIDs = serviceIDs[i:]
	}

	return services, nil
}

func (this *ECS) ListClusters() ([]*string, error) {
	input := &ecs.ListClustersInput{}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output := []*string{}
	save := func(p *ecs.ListClustersOutput, lastPage bool) (shouldContinue bool) {
		output = append(output, p.ClusterArns...)
		return !lastPage
	}

	if err := connection.ListClustersPages(input, save); err != nil {
		return nil, err
	}

	return output, nil
}

func (this *ECS) Helper_ListServices(prefix string) ([]*string, error) {
	clusterARNs, err := this.ListClusters()
	if err != nil {
		return nil, err
	}

	serviceARNs := []*string{}
	for _, clusterARN := range clusterARNs {
		clusterName := strings.Split(*clusterARN, ":")[5]
		clusterName = strings.Replace(clusterName, "cluster/", "", -1)

		if !strings.HasPrefix(clusterName, prefix) {
			continue
		}

		arns, err := this.ListServices(*clusterARN)
		if err != nil {
			return nil, err
		}

		serviceARNs = append(serviceARNs, arns...)
	}

	return serviceARNs, nil
}

func (this *ECS) Helper_DescribeServices(prefix string) ([]*Service, error) {
	clusterARNs, err := this.ListClusters()
	if err != nil {
		return nil, err
	}

	clusterServices := map[string][]*string{}
	for _, clusterARN := range clusterARNs {
		clusterName := strings.Split(*clusterARN, ":")[5]
		clusterName = strings.Replace(clusterName, "cluster/", "", -1)

		if !strings.HasPrefix(clusterName, prefix) {
			continue
		}

		serviceARNs, err := this.ListServices(*clusterARN)
		if err != nil {
			return nil, err
		}

		clusterServices[*clusterARN] = serviceARNs
	}

	services := []*Service{}
	for clusterARN, serviceARNs := range clusterServices {
		svcs, err := this.DescribeServices(clusterARN, serviceARNs)
		if err != nil {
			return nil, err
		}

		services = append(services, svcs...)
	}

	return services, nil
}

func (this *ECS) CreateCluster(clusterName string) (*Cluster, error) {
	input := &ecs.CreateClusterInput{
		ClusterName: aws.String(clusterName),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.CreateCluster(input)
	if err != nil {
		return nil, err
	}

	return &Cluster{output.Cluster}, err
}

func (this *ECS) ListClusterNames(prefix string) ([]string, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	clusterNames := []string{}
	fn := func(output *ecs.ListClustersOutput, lastPage bool) bool {
		for _, arn := range output.ClusterArns {
			// cluster arn format: arn:aws:ecs:region:012345678910:cluster/name
			clusterName := strings.Split(aws.StringValue(arn), "/")[1]

			if strings.HasPrefix(clusterName, prefix) {
				clusterNames = append(clusterNames, clusterName)
			}
		}

		return !lastPage
	}

	if err := connection.ListClustersPages(&ecs.ListClustersInput{}, fn); err != nil {
		return nil, err
	}

	return clusterNames, nil
}

func (this *ECS) ListClusterTaskARNs(clusterName, startedBy string) ([]string, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	taskARNs := []string{}
	fn := func(output *ecs.ListTasksOutput, lastPage bool) bool {
		for _, taskARN := range output.TaskArns {
			taskARNs = append(taskARNs, aws.StringValue(taskARN))
		}

		return !lastPage
	}

	for _, status := range []string{ecs.DesiredStatusRunning, ecs.DesiredStatusStopped} {
		input := &ecs.ListTasksInput{}
		input.SetCluster(clusterName)
		input.SetDesiredStatus(status)
		input.SetStartedBy(startedBy)

		if err := connection.ListTasksPages(input, fn); err != nil {
			return nil, err
		}
	}

	return taskARNs, nil
}

func (this *ECS) ListClusterServiceNames(clusterName, prefix string) ([]string, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	var serviceNames []string
	fn := func(output *ecs.ListServicesOutput, lastPage bool) bool {
		for _, serviceARN := range output.ServiceArns {
			// sample service ARN:
			// arn:aws:ecs:us-west-2:856306994068:service/l0-tlakedev-guestbo80d9d
			serviceName := strings.Split(aws.StringValue(serviceARN), "/")[1]
			if strings.HasPrefix(serviceName, prefix) {
				serviceNames = append(serviceNames, serviceName)
			}
		}

		return !lastPage
	}

	input := &ecs.ListServicesInput{}
	input.SetCluster(clusterName)
	if err := connection.ListServicesPages(input, fn); err != nil {
		return nil, err
	}

	return serviceNames, nil
}

func (this *ECS) DescribeCluster(cluster string) (*Cluster, error) {
	clusters, err := this.DescribeClusters([]string{cluster})
	if err != nil {
		return nil, err
	}

	if len(clusters) > 0 {
		cluster := clusters[0]
		if *cluster.Status != "INACTIVE" {
			return cluster, nil
		}
	}

	return nil, fmt.Errorf("Cluster Not Found")
}

func (this *ECS) DescribeClusters(cluster []string) ([]*Cluster, error) {
	input_list := make([]*string, 0, len(cluster))
	for _, svc := range cluster {
		input_list = append(input_list, &svc)
	}

	input := &ecs.DescribeClustersInput{
		Clusters: input_list,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	descs, err := connection.DescribeClusters(input)
	if err != nil {
		return nil, err
	}

	result := []*Cluster{}
	for _, svc := range descs.Clusters {
		result = append(result, &Cluster{svc})
	}

	return result, nil
}

func (this *ECS) Helper_DescribeClusters() ([]*Cluster, error) {
	clusterARNs, err := this.ListClusters()
	if err != nil {
		return nil, err
	}

	input := &ecs.DescribeClustersInput{
		Clusters: clusterARNs,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.DescribeClusters(input)
	if err != nil {
		return nil, err
	}

	clusters := []*Cluster{}
	for _, description := range output.Clusters {
		if *description.Status != "INACTIVE" {
			clusters = append(clusters, &Cluster{description})
		}
	}

	return clusters, nil
}

func (this *ECS) DeleteCluster(clusterName string) error {
	input := &ecs.DeleteClusterInput{
		Cluster: aws.String(clusterName),
	}
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteCluster(input)
	if err != nil {
		return err
	}
	return nil
}

func (this *ECS) ListContainerInstances(clusterName string) ([]*string, error) {
	input := &ecs.ListContainerInstancesInput{
		Cluster: &clusterName,
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output := []*string{}
	save := func(p *ecs.ListContainerInstancesOutput, lastPage bool) bool {
		output = append(output, p.ContainerInstanceArns...)
		return !lastPage
	}
	err = connection.ListContainerInstancesPages(input, save)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (this *ECS) DescribeTaskDefinition(familyAndRevision string) (*TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &familyAndRevision,
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return &TaskDefinition{output.TaskDefinition}, nil
}

func (this *ECS) ListTaskDefinitionFamiliesPages(prefix string) ([]*string, error) {
	input := &ecs.ListTaskDefinitionFamiliesInput{
		FamilyPrefix: &prefix,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output := []*string{}
	save := func(p *ecs.ListTaskDefinitionFamiliesOutput, lastPage bool) bool {
		output = append(output, p.Families...)
		return !lastPage
	}

	if err := connection.ListTaskDefinitionFamiliesPages(input, save); err != nil {
		return nil, err
	}

	return output, nil
}

func (this *ECS) Helper_DescribeTaskDefinitions(prefix string) ([]*TaskDefinition, error) {
	familyNames, err := this.ListTaskDefinitionFamiliesPages(prefix)
	if err != nil {
		return nil, err
	}

	taskDefinitions := []*TaskDefinition{}
	for _, familyName := range familyNames {
		taskDefinitionARNs, err := this.ListTaskDefinitionsPages(*familyName)
		if err != nil {
			return nil, err
		}

		for _, taskDefinitionARN := range taskDefinitionARNs {
			taskDefinition, err := this.DescribeTaskDefinition(*taskDefinitionARN)
			if err != nil {
				return nil, err
			}

			td := &TaskDefinition{taskDefinition.TaskDefinition}
			taskDefinitions = append(taskDefinitions, td)
		}
	}

	return taskDefinitions, nil
}

func (this *ECS) Helper_ListTaskDefinitions(prefix string) ([]*string, error) {
	familyNames, err := this.ListTaskDefinitionFamiliesPages(prefix)
	if err != nil {
		return nil, err
	}

	taskDefinitionARNs := []*string{}
	for _, familyName := range familyNames {
		arns, err := this.ListTaskDefinitionsPages(*familyName)
		if err != nil {
			return nil, err
		}

		taskDefinitionARNs = append(taskDefinitionARNs, arns...)
	}

	return taskDefinitionARNs, nil
}

func (this *ECS) ListTaskDefinitions(familyName string, nextToken *string) ([]*string, *string, error) {
	input := &ecs.ListTaskDefinitionsInput{
		// The FamilyPrefix is misleading for this input. This must be the full family name
		FamilyPrefix: aws.String(familyName),
		NextToken:    nextToken,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, nil, err
	}

	output, err := connection.ListTaskDefinitions(input)
	if err != nil {
		return nil, nil, err
	}

	return output.TaskDefinitionArns, output.NextToken, nil
}

func (this *ECS) ListTaskDefinitionsPages(familyName string) ([]*string, error) {
	input := &ecs.ListTaskDefinitionsInput{
		// The FamilyPrefix is misleading for this input. This must be the full family name
		FamilyPrefix: aws.String(familyName),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output := []*string{}
	save := func(p *ecs.ListTaskDefinitionsOutput, lastPage bool) bool {
		output = append(output, p.TaskDefinitionArns...)
		return !lastPage
	}

	if err := connection.ListTaskDefinitionsPages(input, save); err != nil {
		return nil, err
	}

	return output, nil
}

func (this *ECS) ListTaskDefinitionFamilies(prefix string, nextToken *string) ([]*string, *string, error) {
	input := &ecs.ListTaskDefinitionFamiliesInput{
		FamilyPrefix: aws.String(prefix),
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, nil, err
	}

	output, err := connection.ListTaskDefinitionFamilies(input)
	if err != nil {
		return nil, nil, err
	}

	return output.Families, output.NextToken, nil
}

func (this *ECS) ListServices(clusterName string) ([]*string, error) {
	input := &ecs.ListServicesInput{
		Cluster: &clusterName,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output := []*string{}
	save := func(p *ecs.ListServicesOutput, lastPage bool) bool {
		output = append(output, p.ServiceArns...)
		return !lastPage
	}

	if err := connection.ListServicesPages(input, save); err != nil {
		return nil, err
	}

	return output, nil
}

func (this *ECS) ListTasks(clusterName string, serviceName, desiredStatus, startedBy, containerInstanceID *string) ([]*string, error) {
	input := &ecs.ListTasksInput{
		Cluster:           &clusterName,
		ServiceName:       serviceName,
		DesiredStatus:     desiredStatus,
		StartedBy:         startedBy,
		ContainerInstance: containerInstanceID,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output := []*string{}
	save := func(p *ecs.ListTasksOutput, lastPage bool) bool {
		output = append(output, p.TaskArns...)
		return !lastPage
	}
	err = connection.ListTasksPages(input, save)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (this *ECS) DescribeTasks(clusterName string, tasks []*string) ([]*Task, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   tasks,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.DescribeTasks(input)
	if err != nil {
		return nil, err
	}

	ret := []*Task{}
	for _, task := range output.Tasks {
		ret = append(ret, &Task{task})
	}

	return ret, nil
}

func (this *ECS) DescribeTask(clusterName string, taskARN string) (*Task, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(clusterName),
		Tasks:   []*string{aws.String(taskARN)},
	}

	output, err := connection.DescribeTasks(input)
	if err != nil {
		return nil, err
	}

	if len(output.Failures) > 0 {
		reason := aws.StringValue(output.Failures[0].Reason)
		if strings.Contains(reason, "MISSING") {
			return nil, fmt.Errorf("The specified task does not exist")
		}

		return nil, fmt.Errorf("Failed to describe task: %s", reason)
	}

	return &Task{output.Tasks[0]}, nil
}

func (this *ECS) DeleteTaskDefinition(taskName string) error {
	input := &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: &taskName,
	}
	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeregisterTaskDefinition(input)
	if err != nil {
		return err
	}

	return nil
}

func (this *ECS) DescribeContainerInstances(clusterName string, instances []*string) ([]*ContainerInstance, error) {
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            &clusterName,
		ContainerInstances: instances,
	}
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	output, err := connection.DescribeContainerInstances(input)
	if err != nil {
		return nil, err
	}

	ret := []*ContainerInstance{}
	for _, i := range output.ContainerInstances {
		ret = append(ret, &ContainerInstance{i})
	}

	errs := []string{}
	for _, fail := range output.Failures {
		errs = append(errs, fmt.Sprintf(
			"Encountered failure with container instance %v: %v", fail.Arn, fail.Reason))
	}

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, ", "))
	} else {
		err = nil
	}

	return ret, err
}
