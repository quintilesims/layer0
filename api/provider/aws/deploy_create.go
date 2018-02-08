package aws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

// Create registers an ECS Task Definition using the specified Create Deploy Request.
// The Create Deploy Request contains the name of the Deploy and the JSON
// representation of the Task Definition to create.
func (d *DeployProvider) Create(req models.CreateDeployRequest) (string, error) {
	deployID := entityIDGenerator(req.DeployName)
	familyName := addLayer0Prefix(d.Config.Instance(), req.DeployName)

	renderedTaskDefinition, err := d.renderTaskDefinition(req.DeployFile, familyName)
	if err != nil {
		return "", err
	}

	taskDefinition, err := d.createTaskDefinition(renderedTaskDefinition)
	if err != nil {
		return "", err
	}

	version := int(aws.Int64Value(taskDefinition.Revision))
	taskDefinitionARN := aws.StringValue(taskDefinition.TaskDefinitionArn)
	if err := d.createTags(deployID, req.DeployName, strconv.Itoa(version), taskDefinitionARN); err != nil {
		return "", err
	}

	return deployID, nil
}

func (d *DeployProvider) createTaskDefinition(taskDefinition *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {
	input := &ecs.RegisterTaskDefinitionInput{}
	input.SetFamily(aws.StringValue(taskDefinition.Family))
	input.SetTaskRoleArn(aws.StringValue(taskDefinition.TaskRoleArn))
	input.SetContainerDefinitions(taskDefinition.ContainerDefinitions)
	input.SetVolumes(taskDefinition.Volumes)
	input.SetPlacementConstraints(taskDefinition.PlacementConstraints)
	if nm := aws.StringValue(taskDefinition.NetworkMode); nm != "" {
		input.SetNetworkMode(nm)
	}

	// We'll be explicit here and set any and all compatibilities that a user specifies in the task definition.
	// At the moment, that should only be "FARGATE" and "EC2".
	requiresCompatibilities := []*string{}
	for _, compatibility := range taskDefinition.RequiresCompatibilities {
		requiresCompatibilities = append(requiresCompatibilities, compatibility)

		// There are some additional requirements for a task definition to be considered Fargate-compatible:
		// https://github.com/aws/aws-sdk-go/blob/master/service/ecs/api.go#L8172
		// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html#task_size
		//
		// For the most part, there are four things that should be specified in a Fargate task definition
		// that aren't required for an EC2 task definition:
		//   "requiresCompatibilities": [ "FARGATE" ],
		//   "networkMode": "awsvpc",
		//   "cpu": "<an appropriate value>",
		//   "memory": "<an appropriate value>"
		if aws.StringValue(compatibility) == ecs.LaunchTypeFargate {
			cpu := aws.StringValue(taskDefinition.Cpu)
			memory := aws.StringValue(taskDefinition.Memory)

			// Question: Do we want to be proactive here and error for certain configurations of a task
			// definition that we know to be bad? They would come out of AWS when the user tried to use the
			// deploy in a service or a task, but it might be nice UX to save them a step.
			if cpu == "" || memory == "" {
				return nil, fmt.Errorf("Fargate task definitions require 'cpu' and 'memory' values")
			}

			input.SetCpu(cpu)
			input.SetMemory(memory)

			// hard-coding an ARN isn't usually the best thing, but we know exactly what this ARN will look like
			// and the only variable in it is the account ID
			accountID := d.Config.AccountID()
			ecsTaskExecutionRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/ecsTaskExecutionRole", accountID)
			input.SetExecutionRoleArn(ecsTaskExecutionRoleARN)
		}
	}

	input.SetRequiresCompatibilities(requiresCompatibilities)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := d.AWS.ECS.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

func (d *DeployProvider) renderTaskDefinition(body []byte, familyName string) (*ecs.TaskDefinition, error) {
	var taskDefinition *ecs.TaskDefinition

	if err := json.Unmarshal(body, &taskDefinition); err != nil {
		return nil, fmt.Errorf("Failed to decode deploy: %s", err.Error())
	}

	if len(taskDefinition.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("Deploy must have at least one container definition")
	}

	taskDefinition.SetFamily(familyName)
	for _, container := range taskDefinition.ContainerDefinitions {
		if container.LogConfiguration == nil {
			logConfig := &ecs.LogConfiguration{}
			logConfig.SetLogDriver("awslogs")
			logConfig.SetOptions(map[string]*string{
				"awslogs-group":         aws.String(d.Config.LogGroupName()),
				"awslogs-region":        aws.String(d.Config.Region()),
				"awslogs-stream-prefix": aws.String("l0"),
			})

			container.SetLogConfiguration(logConfig)
		}
	}

	return taskDefinition, nil
}

func (d *DeployProvider) createTags(deployID, deployName, deployVersion, taskDefinitionARN string) error {
	tags := []models.Tag{
		{
			EntityID:   deployID,
			EntityType: "deploy",
			Key:        "name",
			Value:      deployName,
		},
		{
			EntityID:   deployID,
			EntityType: "deploy",
			Key:        "version",
			Value:      deployVersion,
		},
		{
			EntityID:   deployID,
			EntityType: "deploy",
			Key:        "arn",
			Value:      taskDefinitionARN,
		},
	}

	for _, tag := range tags {
		if err := d.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
