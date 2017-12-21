package aws

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type AdminProvider struct {
	AWS      *awsc.Client
	TagStore tag.Store
	Context  *cli.Context
}

func NewAdminProvider(a *awsc.Client, t tag.Store, c *cli.Context) *AdminProvider {
	return &AdminProvider{
		AWS:      a,
		TagStore: t,
		Context:  c,
	}
}

func (a *AdminProvider) Init() error {
	log.Printf("[DEBUG] Adding API tags")

	fqAPI := addLayer0Prefix(a.Context, "api")
	service, err := readService(a.AWS.ECS, fqAPI, fqAPI)
	if err != nil {
		return err
	}

	// arn format: "arn:aws:ecs:region:123:task-definition/l0-instance-api:1"
	taskDefinitionARN := aws.StringValue(service.TaskDefinition)
	split := strings.Split(taskDefinitionARN, ":")
	taskDefinitionRevision := split[len(split)-1]

	tags := []models.Tag{
		{EntityID: "api", EntityType: "deploy", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "deploy", Key: "arn", Value: taskDefinitionARN},
		{EntityID: "api", EntityType: "deploy", Key: "version", Value: taskDefinitionRevision},
		{EntityID: "api", EntityType: "environment", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "environment", Key: "os", Value: "linux"},
		{EntityID: "api", EntityType: "load_balancer", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "load_balancer", Key: "environment_id", Value: "api"},
		{EntityID: "api", EntityType: "service", Key: "name", Value: "api"},
		{EntityID: "api", EntityType: "service", Key: "environment_id", Value: "api"},
	}

	for _, tag := range tags {
		if err := a.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
