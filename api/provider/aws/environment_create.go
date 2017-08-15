package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Create(req models.CreateEnvironmentRequest) (*models.Environment, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// todo: use config to get instance name
	environmentID := generateEntityID(req.EnvironmentName)
	fqEnvironmentID := addLayer0Prefix("INSTANCE", environmentID)

	// todo: use config to get vpc id
	securityGroup, err := e.createSG(fqEnvironmentID, "VPC")
	if err != nil {
		return nil, err
	}

	instanceType := DEFAULT_INSTANCE_SIZE
	if req.InstanceSize != "" {
		instanceType = req.InstanceSize
	}

	// todo: use config.AMIID()
	amiID := "AMIID"
	if req.AMIID != "" {
		amiID = req.AMIID
	}

	userDataTemplate := DEFAULT_USER_DATA_TEMPLATE
	if len(req.UserDataTemplate) > 0 {
		userDataTemplate = req.UserDataTemplate
	}

	// todo: use config.S3Bucket()
	userData, err := renderUserData(fqEnvironmentID, "BUCKET", userDataTemplate)
	if err != nil {
		return nil, err
	}

	// todo: use config.InstanceProfile()
	// todo: operating systems
	if err := e.createLC(
		fqEnvironmentID,
		aws.StringValue(securityGroup.GroupId),
		instanceType,
		"INSTANCEPROFILE",
		amiID,
		userData); err != nil {
		return nil, err
	}

	// todo: use private subnets
	// launchConfig name is same as environmentID
	if err := e.createASG(
		fqEnvironmentID,
		fqEnvironmentID,
		[]string{}); err != nil {
		return nil, err
	}

	if err := e.createCluster(fqEnvironmentID); err != nil {
		return nil, err
	}

	if err := e.createTags(environmentID, req.EnvironmentName, req.OperatingSystem); err != nil {
		return nil, err
	}

	return e.Read(environmentID)
}

func (e *EnvironmentProvider) createSG(environmentID, vpcID string) (*ec2.SecurityGroup, error) {
	input := &ec2.CreateSecurityGroupInput{}
	input.SetGroupName(environmentID)
	input.SetDescription(fmt.Sprintf("SG for Layer0 environment %s", environmentID))
	input.SetVpcId(vpcID)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if _, err := e.AWS.EC2.CreateSecurityGroup(input); err != nil {
		return nil, err
	}

	securityGroup, err := e.readSG(environmentID)
	if err != nil {
		return nil, err
	}

	groupPair := &ec2.UserIdGroupPair{
		GroupId: securityGroup.GroupId,
	}

	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: securityGroup.GroupId,
	}

	ingressInput.SetIpPermissions([]*ec2.IpPermission{permission})

	if _, err := e.AWS.EC2.AuthorizeSecurityGroupIngress(ingressInput); err != nil {
		return nil, err
	}

	return securityGroup, nil
}

func (e *EnvironmentProvider) createLC(
	launchConfigName string,
	securityGroupID string,
	instanceType string,
	instanceProfile string,
	amiID string,
	userData string,
) error {
	input := &autoscaling.CreateLaunchConfigurationInput{}
	input.SetLaunchConfigurationName(launchConfigName)
	input.SetSecurityGroups([]*string{aws.String(securityGroupID)})
	input.SetInstanceType(instanceType)
	input.SetIamInstanceProfile(instanceProfile)
	input.SetImageId(amiID)
	input.SetUserData(userData)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.AutoScaling.CreateLaunchConfiguration(input); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createASG(autoScalingGroupName, launchConfigName string, privateSubnets []string) error {
	tag := &autoscaling.Tag{}
	tag.SetKey("Name")
	tag.SetValue(autoScalingGroupName)
	tag.SetPropagateAtLaunch(true)

	var subnetIdentifier string
	for _, subnet := range privateSubnets {
		subnetIdentifier = fmt.Sprintf("%s%s,", subnetIdentifier, subnet)
	}

	input := &autoscaling.CreateAutoScalingGroupInput{}
	input.SetAutoScalingGroupName(autoScalingGroupName)
	input.SetLaunchConfigurationName(launchConfigName)
	input.SetVPCZoneIdentifier(subnetIdentifier)
	input.SetMinSize(0)
	input.SetMaxSize(0)
	input.SetTags([]*autoscaling.Tag{tag})

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.AutoScaling.CreateAutoScalingGroup(input); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createCluster(clusterName string) error {
	input := &ecs.CreateClusterInput{}
	input.SetClusterName(clusterName)

	if _, err := e.AWS.ECS.CreateCluster(input); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createTags(environmentID, environmentName, operatingSystem string) error {
	tags := []models.Tag{
		{
			EntityID:   environmentID,
			EntityType: "environment",
			Key:        "name",
			Value:      environmentName,
		},
		{
			EntityID:   environmentID,
			EntityType: "environment",
			Key:        "os",
			Value:      strings.ToLower(operatingSystem),
		},
	}

	for _, tag := range tags {
		if err := e.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}

func renderUserData(environmentID, s3Bucket string, userDataTemplate []byte) (string, error) {
	tmpl, err := template.New("").Parse(string(userDataTemplate))
	if err != nil {
		return "", fmt.Errorf("Failed to parse user data: %v", err)
	}

	context := struct {
		ECSEnvironmentID string
		S3Bucket         string
	}{
		ECSEnvironmentID: environmentID,
		S3Bucket:         s3Bucket,
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, context); err != nil {
		return "", fmt.Errorf("Failed to render user data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(rendered.Bytes()), nil
}
