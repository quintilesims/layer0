package aws

import (
	"fmt"
	"log"
	"strconv"

	"github.com/zpatrick/go-bytesize"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/provider"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

type EnvironmentScaler struct {
	Client              *awsc.Client
	EnvironmentProvider provider.EnvironmentProvider
	ServiceProvider     provider.ServiceProvider
	Config              config.APIConfig
	deployCache         map[string][]models.ResourceConsumer
}

func NewEnvironmentScaler(a *awsc.Client, e provider.EnvironmentProvider, s provider.ServiceProvider, c config.APIConfig) *EnvironmentScaler {
	return &EnvironmentScaler{
		Client:              a,
		EnvironmentProvider: e,
		ServiceProvider:     s,
		Config:              c,
	}
}

func (e *EnvironmentScaler) Scale(environmentID string) error {
	// get resource providers
	clusterName := addLayer0Prefix(e.Config.Instance(), environmentID)

	resourceProviders, err := e.getResourceProviders(clusterName)
	if err != nil {
		return err
	}

	// resourceConsumers, err := e.getResourceConsumers(clusterName)
	// if err != nil {
	// 	return err
	// }

	fmt.Println(resourceProviders)

	// calculate desired capacity

	// scale to new capacity if required

	return fmt.Errorf("EnvironmentScaler not implemented")
}

func (e *EnvironmentScaler) getResourceProviders(clusterName string) ([]models.ResourceProvider, error) {
	listContainerInstancesInput := &ecs.ListContainerInstancesInput{}
	listContainerInstancesInput.SetCluster(clusterName)

	containerInstanceARNs := []*string{}
	listContainerInstancesPagesFN := func(output *ecs.ListContainerInstancesOutput, lastPage bool) bool {
		containerInstanceARNs = append(containerInstanceARNs, output.ContainerInstanceArns...)

		return !lastPage
	}

	if err := e.Client.ECS.ListContainerInstancesPages(listContainerInstancesInput, listContainerInstancesPagesFN); err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] containerInstanceARNs:")
	for _, arn := range containerInstanceARNs {
		log.Printf("[DEBUG] %#v", aws.StringValue(arn))
	}

	describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{}
	describeContainerInstancesInput.SetCluster(clusterName)
	describeContainerInstancesInput.SetContainerInstances(containerInstanceARNs)

	output, err := e.Client.ECS.DescribeContainerInstances(describeContainerInstancesInput)
	if err != nil {
		return nil, err
	}

	result := []models.ResourceProvider{}
	if len(output.ContainerInstances) == 0 {
		return result, nil
	}

	for _, instance := range output.ContainerInstances {
		// todo: not sure if this is a permanent condition?
		// if !instance.AgentConnected {
		// 	continue
		// }

		// if aws.StringValue(instance.Status) != "ACTIVE" {
		// 	continue
		// }

		usedPorts := []int{}
		var availableMemory bytesize.Bytesize

		for _, resource := range instance.RemainingResources {
			switch aws.StringValue(resource.Name) {
			case "MEMORY":
				availableMemory = bytesize.MiB * bytesize.Bytesize(aws.Int64Value(resource.IntegerValue))

			case "PORTS":
				for _, port := range resource.StringSetValue {
					port, err := strconv.Atoi(aws.StringValue(port))
					if err != nil {
						//todo: log error
						continue
					}

					usedPorts = append(usedPorts, port)
				}
			}
		}

		inUse := aws.Int64Value(instance.PendingTasksCount)+aws.Int64Value(instance.RunningTasksCount) > 0
		r := models.ResourceProvider{
			ID:              aws.StringValue(instance.Ec2InstanceId),
			InUse:           inUse,
			UsedPorts:       usedPorts,
			AvailableMemory: fmt.Sprintf("%v", availableMemory), //todo: this is definitely not correct
		}

		result = append(result, r)
	}

	return result, nil
}

func (e *EnvironmentScaler) getResourceConsumers(clusterName string) ([]models.ResourceConsumer, error) {
	// from scaler in develop, there are funcs to aggregate three types of resources
	//   - getPendingServiceResources
	//   - getPendingTaskResourcesInECS
	//   - getPendingTaskResourcesInJobs

	// get pending service resource consumers
	input := &ecs.DescribeServicesInput{}
	input.SetCluster(clusterName)
	output, err := e.Client.ECS.DescribeServices(input)
	if err != nil {
		return nil, err
	}

	result := []models.ResourceConsumer{}

	for _, service := range output.Services {
		deployIDCopies := map[string]int64{}
		for _, d := range service.Deployments {
			desiredCount := aws.Int64Value(d.DesiredCount)
			runningCount := aws.Int64Value(d.RunningCount)
			pendingCount := aws.Int64Value(d.PendingCount)
			if numPending := desiredCount - (runningCount + pendingCount); numPending > 0 {
				deployIDCopies[aws.StringValue(d.Id)] = numPending
			}
		}

		if len(deployIDCopies) == 0 {
			continue
		}

		// iterate through deploys
		for deployID := range deployIDCopies {
			c, err := e.getContainerResourceFromDeploy(deployID)
			if err != nil {
				return nil, err
			}

			result = append(result, c...)
		}
	}

	// get pending task resource consumers in ECS
	// input := &ecs.DescribeTasksInput{}
	// input.SetCluster(clusterName)
	// output, err := e.Client.ECS.DescribeTasks(input)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, task := range output.Tasks {
	// }

	// get pending task resource consumers in jobs

	return nil, nil
}

func (e *EnvironmentScaler) getContainerResourceFromDeploy(deployID string) ([]models.ResourceConsumer, error) {
	// use some kind of deploy cache
	if consumers, ok := e.deployCache[deployID]; ok {
		return consumers, nil
	}

	input := &ecs.DescribeTaskDefinitionInput{}
	input.SetTaskDefinition(deployID)
	output, err := e.Client.ECS.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	consumers := make([]models.ResourceConsumer, len(output.TaskDefinition.ContainerDefinitions))
	for i, d := range output.TaskDefinition.ContainerDefinitions {
		var memory bytesize.Bytesize

		if aws.Int64Value(d.MemoryReservation) != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(aws.Int64Value(d.MemoryReservation))
		}

		if aws.Int64Value(d.Memory) != 0 {
			memory = bytesize.MiB * bytesize.Bytesize(aws.Int64Value(d.Memory))
		}

		ports := []int{}
		for _, p := range d.PortMappings {
			if aws.Int64Value(p.HostPort) != 0 {
				ports = append(ports, int(aws.Int64Value(p.HostPort)))
			}
		}

		consumers[i] = models.ResourceConsumer{
			ID:     "",
			Memory: fmt.Sprintf("%v", memory), //todo: translate bytesize.ByteSize to string correctly
			Ports:  ports,
		}
	}

	e.deployCache[deployID] = consumers
	return consumers, nil
}
