package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/quintilesims/layer0/api/tag"
)

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

func listClusterNames(ecsapi ecsiface.ECSAPI, instance string) ([]string, error) {
	clusterNames := []string{}
	fn := func(output *ecs.ListClustersOutput, lastPage bool) bool {
		for _, arn := range output.ClusterArns {
			// cluster arn format: arn:aws:ecs:region:012345678910:cluster/name
			clusterName := strings.Split(aws.StringValue(arn), "/")[1]

			if hasLayer0Prefix(instance, clusterName) {
				clusterNames = append(clusterNames, clusterName)
			}
		}

		return !lastPage
	}

	if err := ecsapi.ListClustersPages(&ecs.ListClustersInput{}, fn); err != nil {
		return nil, err
	}

	return clusterNames, nil
}
