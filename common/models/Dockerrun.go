package models

import (
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Dockerrun struct {
	ContainerDefinitions []*ecs.ContainerDefinition               `json:"containerDefinitions,omitempty"`
	Volumes              []*ecs.Volume                            `json:"volumes,omitempty"`
	Family               string                                   `json:"family,omitempty"`
	NetworkMode          string                                   `json:"networkMode,omitempty"`
	TaskRoleARN          string                                   `json:"taskRoleArn,omitempty"`
	PlacementConstraints []*ecs.TaskDefinitionPlacementConstraint `json:"placementConstraints,omitempty"`
}
