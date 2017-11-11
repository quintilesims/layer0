package aws

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Delete deregisters an ECS Task Definition using the specified deployID. The deployID is used
// to look up the associated Task Definition ARN. The Task Definition ARN is subsequently used
// when the DeregisterTaskDefinition request is made to AWS.
func (d *DeployProvider) Delete(deployID string) error {
	taskARN, err := d.lookupTaskDefinitionARN(deployID)
	if err != nil {
		return err
	}

	if err := d.deleteDeploy(taskARN); err != nil {
		return err
	}

	if err := deleteEntityTags(d.TagStore, "deploy", deployID); err != nil {
		return err
	}

	return nil
}

func (d *DeployProvider) deleteDeploy(taskARN string) error {
	input := &ecs.DeregisterTaskDefinitionInput{}
	input.SetTaskDefinition(taskARN)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := d.AWS.ECS.DeregisterTaskDefinition(input); err != nil {
		log.Printf("[WARN] Deploy not found\n")
		fmt.Println(err)
		if err, ok := err.(awserr.Error); ok && strings.Contains(err.Message(), "does not exist") {
			return nil
		}

		return err
	}

	return nil
}
