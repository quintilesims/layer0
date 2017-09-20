package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/errors"
)

func lookupEntityEnvironmentID(store tag.Store, entityType, entityID string) (string, error) {
	tags, err := store.SelectByTypeAndID(entityType, entityID)
	if err != nil {
		return "", err
	}

	tag, ok := tags.WithKey("environment_id").First()
	if !ok {
		// is this canonical error handling?
		return "", fmt.Errorf("Cannot resolve environment_id for %s '%s'", entityType, entityID)
	}

	return tag.Value, nil
}

func lookupTaskDefinitionFamily(store tag.Store, deployID string) (string, error) {
	tags, err := store.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		return tag.Value, nil
	}

	return "", errors.Newf(errors.DeployDoesNotExist, "Could not resolve task definition family for deploy '%s'", deployID)
}

func lookupTaskDefinitionRevision(store tag.Store, deployID string) (string, error) {
	tags, err := store.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("version").First(); ok {
		return tag.Value, err
	}

	return "", errors.Newf(errors.DeployDoesNotExist, "Could not resolve task definition revision for deploy '%s'", deployID)
}

func getEnvironmentSGName(environmentID string) string {
	return fmt.Sprintf("%s-env", environmentID)
}

func getLoadBalancerSGName(loadBalancerID string) string {
	return fmt.Sprintf("%s-lb", loadBalancerID)
}

func getLoadBalancerRoleName(loadBalancerID string) string {
	return fmt.Sprintf("%s-lb", loadBalancerID)
}

func createSG(ec2api ec2iface.EC2API, groupName, description, vpcID string) error {
	input := &ec2.CreateSecurityGroupInput{}
	input.SetGroupName(groupName)
	input.SetDescription(description)
	input.SetVpcId(vpcID)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := ec2api.CreateSecurityGroup(input); err != nil {
		return err
	}

	return nil
}

func readSG(ec2api ec2iface.EC2API, groupName string) (*ec2.SecurityGroup, error) {
	filter := &ec2.Filter{}
	filter.SetName("group-name")
	filter.SetValues([]*string{aws.String(groupName)})

	input := &ec2.DescribeSecurityGroupsInput{}
	input.SetFilters([]*ec2.Filter{filter})

	output, err := ec2api.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	for _, group := range output.SecurityGroups {
		if aws.StringValue(group.GroupName) == groupName {
			return group, nil
		}
	}

	// todo: this should be a wrapped error: 'errors.MissingResource' or something
	return nil, fmt.Errorf("Security group '%s' does not exist", groupName)
}

func deleteSG(ec2api ec2iface.EC2API, securityGroupID string) error {
	input := &ec2.DeleteSecurityGroupInput{}
	input.SetGroupId(securityGroupID)

	if _, err := ec2api.DeleteSecurityGroup(input); err != nil {
		return err
	}

	return nil
}

func deleteEntityTags(tagStore tag.Store, entityType, entityID string) error {
	tags, err := tagStore.SelectByTypeAndID(entityType, entityID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if err := tagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
			return err
		}
	}

	return nil
}
