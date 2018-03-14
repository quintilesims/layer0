package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	albiface "github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/errors"
)

// searches for the load balancer name as both classic and application load balancers and returns
// the first found result or an error if the neither classic or application lb could be found for
// the give lb name
func describeLoadBalancer(elbapi elbiface.ELBAPI, albapi albiface.ELBV2API, loadBalancerName string) (*genericLoadBalancer, error) {
	// search for classic ELB
	searchELB := func() (*genericLoadBalancer, error) {
		input := &elb.DescribeLoadBalancersInput{}
		input.SetLoadBalancerNames([]*string{aws.String(loadBalancerName)})
		input.SetPageSize(1)

		output, err := elbapi.DescribeLoadBalancers(input)
		if err != nil {
			return nil, err
		}

		if len(output.LoadBalancerDescriptions) != 1 {
			return nil, errors.Newf(errors.LoadBalancerDoesNotExist, "LoadBalancer '%s' does not exist", loadBalancerName)
		}

		return &genericLoadBalancer{
			ELB:   output.LoadBalancerDescriptions[0],
			isELB: true,
		}, nil
	}

	result, err := searchELB()
	if err == nil {
		return result, nil
	}

	// return error is it isn't a 'LoadBalancerNotFound' error
	if err, ok := err.(awserr.Error); ok && err.Code() != "LoadBalancerNotFound" {
		return nil, err
	}

	// search application load balancers
	input := &alb.DescribeLoadBalancersInput{}
	input.SetNames([]*string{aws.String(loadBalancerName)})

	output, err := albapi.DescribeLoadBalancers(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "LoadBalancerNotFound" {
			return nil, errors.Newf(errors.LoadBalancerDoesNotExist, "LoadBalancer '%s' does not exist", loadBalancerName)
		}
		return nil, err
	}

	if len(output.LoadBalancers) != 1 {
		return nil, errors.Newf(errors.LoadBalancerDoesNotExist, "LoadBalancer '%s' does not exist", loadBalancerName)
	}

	return &genericLoadBalancer{
		ALB:   output.LoadBalancers[0],
		isALB: true,
	}, nil

}

func describeTaskDefinition(ecsapi ecsiface.ECSAPI, taskDefinitionARN string) (*ecs.TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{}
	input.SetTaskDefinition(taskDefinitionARN)

	output, err := ecsapi.DescribeTaskDefinition(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && strings.Contains(err.Message(), "Unable to describe task definition") {
			return nil, errors.Newf(errors.DeployDoesNotExist, "Deploy '%s' does not exist", taskDefinitionARN)
		}

		return nil, err
	}

	return output.TaskDefinition, nil
}

func lookupDeployIDFromTaskDefinitionARN(store tag.Store, taskDefinitionARN string) (string, error) {
	tags, err := store.SelectByType("deploy")
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("arn").WithValue(taskDefinitionARN).First(); ok {
		return tag.EntityID, nil
	}

	return "", errors.Newf(errors.DeployDoesNotExist, "Failed to find deploy with ARN '%s'", taskDefinitionARN)
}

func lookupTaskDefinitionARNFromDeployID(store tag.Store, deployID string) (string, error) {
	tags, err := store.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return "", err
	}

	if tag, ok := tags.WithKey("arn").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Could not resolve task definition ARN for deploy '%s'", deployID)
}

func lookupEntityEnvironmentID(store tag.Store, entityType, entityID string) (string, error) {
	tags, err := store.SelectByTypeAndID(entityType, entityID)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", errors.NewEntityDoesNotExistError(entityType, entityID)
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		return tag.Value, nil
	}

	return "", fmt.Errorf("Could not resolve environment ID for %s '%s'", entityType, entityID)
}

func lookupDeployNameAndVersion(store tag.Store, deployID string) (string, string, error) {
	tags, err := store.SelectByTypeAndID("deploy", deployID)
	if err != nil {
		return "", "", err
	}

	if len(tags) == 0 {
		return "", "", errors.Newf(errors.DeployDoesNotExist, "Deploy '%s' does not exist", deployID)
	}

	nameTag, ok := tags.WithKey("name").First()
	if !ok {
		return "", "", fmt.Errorf("Could not resolve name for deploy '%s'", deployID)
	}

	versionTag, ok := tags.WithKey("version").First()
	if !ok {
		return "", "", fmt.Errorf("Could not resolve version for deploy '%s'", deployID)
	}

	deployName := nameTag.Value
	deployVersion := versionTag.Value

	return deployName, deployVersion, nil
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

	message := fmt.Sprintf("Security group '%s' does not exist", groupName)
	return nil, awserr.New("DoesNotExist", message, nil)
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

func listClusterTaskARNs(ecsapi ecsiface.ECSAPI, clusterName, startedBy, status string) ([]string, error) {
	taskARNs := []string{}
	fn := func(output *ecs.ListTasksOutput, lastPage bool) bool {
		for _, taskARN := range output.TaskArns {
			taskARNs = append(taskARNs, aws.StringValue(taskARN))
		}

		return !lastPage
	}

	input := &ecs.ListTasksInput{}
	input.SetCluster(clusterName)
	input.SetDesiredStatus(status)
	input.SetStartedBy(startedBy)

	if err := ecsapi.ListTasksPages(input, fn); err != nil {
		return nil, err
	}

	return taskARNs, nil
}

func readService(ecsapi ecsiface.ECSAPI, clusterName, serviceID string) (*ecs.Service, error) {
	input := &ecs.DescribeServicesInput{}
	input.SetCluster(clusterName)
	input.SetServices([]*string{
		aws.String(serviceID),
	})

	output, err := ecsapi.DescribeServices(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ServiceNotFoundException" {
			return nil, errors.Newf(errors.ServiceDoesNotExist, "Service '%s' does not exist", serviceID)
		}

		return nil, err
	}

	if len(output.Services) == 0 {
		return nil, errors.Newf(errors.ServiceDoesNotExist, "Service '%s' does not exist", serviceID)
	}

	return output.Services[0], nil
}

func getTargetGroupArn(albapi albiface.ELBV2API, targetGroupName string) (*alb.TargetGroup, error) {
	input := &alb.DescribeTargetGroupsInput{}
	input.SetNames([]*string{aws.String(targetGroupName)})

	output, err := albapi.DescribeTargetGroups(input)
	if err != nil {
		return nil, err
	}

	if len(output.TargetGroups) == 0 {
		return nil, fmt.Errorf("target group with name '%s' does not exist", targetGroupName)
	}

	return output.TargetGroups[0], nil
}
