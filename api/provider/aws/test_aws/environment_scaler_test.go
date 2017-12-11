package test_aws

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job/mock_job"
	provider "github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config/mock_config"

	"github.com/quintilesims/layer0/common/models"
	bytesize "github.com/zpatrick/go-bytesize"
)

type EnvironmentScalerUnitTest struct {
	ExpectedScale     int
	MemoryPerProvider bytesize.Bytesize
	ResourceProviders []*models.Resource
	ResourceConsumers []models.Resource
}

// Add Old Test Verbatim from Dev PR

func (e *EnvironmentScalerUnitTest) Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAWS := awsc.NewMockClient(ctrl)
	// tagStore := tag.NewMemoryStore()
	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	mockServiceProvider := mock_provider.NewMockServiceProvider(ctrl)
	mockTaskProvider := mock_provider.NewMockTaskProvider(ctrl)
	mockJobStore := mock_job.NewMockStore(ctrl)
	mockConfig := mock_config.NewMockAPIConfig(ctrl)

	mockConfig.EXPECT().Instance().Return("test").AnyTimes()
	clusterName := "l0-test-env_id"

	// GetResouceProvider
	listContainerInstancesInput := &ecs.ListContainerInstancesInput{}
	listContainerInstancesInput.SetCluster(clusterName)
	listContainerInstancesInput.SetStatus("ACTIVE")

	// TODO: Change to strings to container
	containerInstanceARNs := make([]*string, len(e.ResourceProviders))
	for i := range e.ResourceProviders {
		containerInstanceARNs[i] = aws.String("arn:aws:ecs:region:012345678910:cluster/l0-test-env_id" + strconv.Itoa(i))
	}

	listContainerInstancesPagesFN := func(input *ecs.ListContainerInstancesInput, fn func(output *ecs.ListContainerInstancesOutput, lastPage bool) bool) {
		output := &ecs.ListContainerInstancesOutput{}
		output.SetContainerInstanceArns(containerInstanceARNs)
		fn(output, true)
	}

	mockAWS.ECS.EXPECT().
		ListContainerInstancesPages(listContainerInstancesInput, gomock.Any()).
		Do(listContainerInstancesPagesFN).
		Return(nil)

	describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{}
	describeContainerInstancesInput.SetCluster(clusterName)
	describeContainerInstancesInput.SetContainerInstances(containerInstanceARNs)

	// Assign Provider Models

	containerInstances := make([]*ecs.ContainerInstance, len(containerInstanceARNs))
	for i := 0; i < len(containerInstanceARNs); i++ {
		containerInstance := &ecs.ContainerInstance{}

		// containerInstance.SetRegisteredResources
		containerInstances[i] = containerInstance
	}

	describeContainerInstancesOutput := &ecs.DescribeContainerInstancesOutput{}
	describeContainerInstancesOutput.SetContainerInstances(containerInstances)

	mockAWS.ECS.EXPECT().
		DescribeContainerInstances(describeContainerInstancesInput).
		Return(describeContainerInstancesOutput, nil)

	// GetResouceConsumer
	listServicesInput := &ecs.ListServicesInput{}
	listServicesInput.SetCluster(clusterName)

	// Assign Consumer Models
	serviceARNs := []*string{
		aws.String("arn:aws:ecs:region:012345678910:cluster/l0-test-env_id1"),
	}

	listServicesPagesFN := func(input *ecs.ListServicesInput, fn func(output *ecs.ListServicesOutput, lastPage bool) bool) {
		output := &ecs.ListServicesOutput{}
		output.SetServiceArns(serviceARNs)

		fn(output, true)
	}

	mockAWS.ECS.EXPECT().
		ListServicesPages(listServicesInput, gomock.Any()).
		Do(listServicesPagesFN).
		Return(nil)

	if len(serviceARNs) > 0 {
		// for i := 0; i < len(serviceARNs); i += 10 {
		// 		// end := i + 10
		describeServicesInput := &ecs.DescribeServicesInput{}
		describeServicesInput.SetCluster(clusterName)
		describeServicesInput.SetServices(serviceARNs[0:len(serviceARNs)])

		deployments := []*ecs.Deployment{}
		deployment := &ecs.Deployment{}
		deployment.SetDesiredCount(1)
		deployment.SetId("test")
		deployments = append(deployments, deployment)

		services := []*ecs.Service{}
		service := &ecs.Service{}
		service.SetDeployments(deployments)

		services = append(services, service)
		describeServicesOutput := &ecs.DescribeServicesOutput{}
		describeServicesOutput.SetServices(services)

		mockAWS.ECS.EXPECT().
			DescribeServices(describeServicesInput).
			Return(describeServicesOutput, nil)
		// containerInstances := make([]*ecs.ContainerInstance, len(containerInstanceARNs))
		// for i, _ := range containerInstances {
		// 	containerInstance := &ecs.ContainerInstance{}

		// 	containerInstances[i] = containerInstance
		// }

		// }
	}

	// GetContainerResourceFromDeploy
	describeTaskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	describeTaskDefinitionInput.SetTaskDefinition("test")

	taskDefinition := &ecs.TaskDefinition{}

	containerDefinitions := make([]*ecs.ContainerDefinition, len(e.ResourceConsumers))

	for i, consumer := range e.ResourceConsumers {
		containerDefinition := &ecs.ContainerDefinition{}
		containerDefinition.SetMemory(int64(consumer.Memory))

		containerDefinitions[i] = containerDefinition
	}

	taskDefinition.SetContainerDefinitions(containerDefinitions)

	describeTaskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	describeTaskDefinitionOutput.SetTaskDefinition(taskDefinition)

	mockAWS.ECS.EXPECT().
		DescribeTaskDefinition(describeTaskDefinitionInput).
		Return(describeTaskDefinitionOutput, nil)

	// GET PENDING TASK RESOURCE CONSUMERS IN ECS

	// TODO: Change to strings to container
	// taskARNs := []*string{
	// 	aws.String("arn:aws:ecs:region:012345678910:task/arn"),
	// }

	listTaskFN := func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) {
		output := &ecs.ListTasksOutput{}
		// output.SetTaskArns(taskARNs)

		fn(output, true)
	}

	for _, status := range []string{ecs.DesiredStatusRunning, ecs.DesiredStatusStopped} {
		listTasksInput := &ecs.ListTasksInput{}
		listTasksInput.SetStartedBy("test")
		listTasksInput.SetDesiredStatus(status)
		listTasksInput.SetCluster(clusterName)

		mockAWS.ECS.EXPECT().
			ListTasksPages(listTasksInput, gomock.Any()).
			Do(listTaskFN).
			Return(nil)

	}

	// // Expect Task
	// for i := range taskARNs {
	// 	task := &models.Task{
	// 		DeployID: "test",
	// 	}

	// 	mockTaskProvider.EXPECT().
	// 		Read(*taskARNs[i]).
	// 		Return(task, nil)
	// }

	// GetContainerResourceFromDeploy
	// describeTaskDefinitionInput := &ecs.DescribeTaskDefinitionInput{}
	// describeTaskDefinitionInput.SetTaskDefinition("test")

	// taskDefinition := &ecs.TaskDefinition{}

	// containerDefinitions := make([]*ecs.ContainerDefinition, len(e.ResourceConsumers))

	// for i, consumer := range e.ResourceConsumers {
	// 	containerDefinition := &ecs.ContainerDefinition{}
	// 	containerDefinition.SetMemory(int64(consumer.Memory))

	// 	containerDefinitions[i] = containerDefinition
	// }

	// taskDefinition.SetContainerDefinitions(containerDefinitions)

	// describeTaskDefinitionOutput := &ecs.DescribeTaskDefinitionOutput{}
	// describeTaskDefinitionOutput.SetTaskDefinition(taskDefinition)

	// mockAWS.ECS.EXPECT().
	// 	DescribeTaskDefinition(describeTaskDefinitionInput).
	// 	Return(describeTaskDefinitionOutput, nil)

	// GET PENDING TASK RESOURCE CONSUMERS IN JOBS

	// jobs := []*models.Job{}
	// jobs = append(jobs, &models.Job{
	// 	Type:    models.CreateDeployJob,
	// 	Status:  models.PendingJobStatus,
	// 	Request: "{request}",
	// })

	mockJobStore.EXPECT().
		SelectAll().
		Return(nil, nil)

	// AutoScaling
	describeAutoScalingGroupsInput := &autoscaling.DescribeAutoScalingGroupsInput{}
	describeAutoScalingGroupsInput.SetAutoScalingGroupNames([]*string{&clusterName})

	group := &autoscaling.Group{}
	group.SetDesiredCapacity(0)

	describeAutoScalingGroup := &autoscaling.DescribeAutoScalingGroupsOutput{}
	describeAutoScalingGroup.SetAutoScalingGroups([]*autoscaling.Group{group})

	// describeAutoScalingGroups := []*autoscaling.DescribeAutoScalingGroupsOutput{}

	// describeAutoScalingGroups = append(describeAutoScalingGroups, describeAutoScalingGroup)

	mockAWS.AutoScaling.EXPECT().
		DescribeAutoScalingGroups(describeAutoScalingGroupsInput).
		Return(describeAutoScalingGroup, nil)

		///CalculateNewProvider
	env := &models.Environment{}
	env.InstanceSize = "100MB"

	for i := 0; i < e.ExpectedScale; i++ {
		mockEnvironmentProvider.EXPECT().
			Read(clusterName).
			Return(env, nil)

		mockEnvironmentProvider.EXPECT().
			Read(clusterName).
			Return(env, nil)
	}

	// Scaleup

	// Make Call to Scale
	// mockScaler.EXPECT().
	// 	Scale("env_id").
	// 	Return(nil)

	environmentScaler := provider.NewEnvironmentScaler(mockAWS.Client(), mockEnvironmentProvider, mockServiceProvider, mockTaskProvider, mockJobStore, mockConfig)

	if err := environmentScaler.Scale("env_id"); err != nil {
		t.Fatal(err)
	}
}

func TestResourceManagerScaleUp_noProviders(t *testing.T) {
	// there are 0 providers in the cluster
	// there is 1 consumer
	// we should scale up to size 1
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     1,
		MemoryPerProvider: bytesize.GB,
		ResourceProviders: []*models.Resource{},
		ResourceConsumers: []models.Resource{
			{Memory: 1024},
		},
	}

	test.Run(t)
}

func TestResourceManagerScaleUp_notEnoughPorts(t *testing.T) {
	// there is 1 provider in the cluster that has port 80 being used
	// there is 1 consumer that needs port 80
	// we should scale up to size 2
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     1,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*models.Resource{
			{ID: "",
				Ports:  []int{80},
				Memory: 1,
				InUse:  true}},
		ResourceConsumers: []models.Resource{
			{Ports: []int{80}},
		},
	}

	test.Run(t)
}

func TestResourceManagerNoScale_noResourceConsumers(t *testing.T) {
	// there are 2 providers in the cluster that are in use
	// there are 0 consumers
	// we should stay at size 2
	test := EnvironmentScalerUnitTest{
		ExpectedScale:     0,
		MemoryPerProvider: bytesize.MB,
		ResourceProviders: []*models.Resource{
			&models.Resource{
				ID:     "",
				InUse:  true,
				Ports:  []int{80},
				Memory: 512,
				CPU:    1024,
			},
			&models.Resource{
				ID:     "",
				InUse:  true,
				Ports:  nil,
				Memory: 1024,
				CPU:    1024,
			},
		},
		ResourceConsumers: []models.Resource{},
	}

	test.Run(t)
}
