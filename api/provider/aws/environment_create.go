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

// todo: ensure envirnment name is unique
func (e *EnvironmentProvider) Create(req models.CreateEnvironmentRequest) (*models.Environment, error) {
	environmentID := generateEntityID(req.EnvironmentName)
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	instanceType := DEFAULT_INSTANCE_SIZE
	if req.InstanceSize != "" {
		instanceType = req.InstanceSize
	}

	var userDataTemplate []byte
	var amiID string

	switch strings.ToLower(req.OperatingSystem) {
	case "linux":
		userDataTemplate = []byte(DEFAULT_LINUX_USERDATA_TEMPLATE)
		amiID = e.Config.LinuxAMI()
	case "windows":
		userDataTemplate = []byte(DEFAULT_WINDOWS_USERDATA_TEMPLATE)
		amiID = e.Config.WindowsAMI()
	default:
		return nil, fmt.Errorf("Operating system '%s' is not recognized", req.OperatingSystem)
	}

	if req.AMIID != "" {
		amiID = req.AMIID
	}

	if len(req.UserDataTemplate) > 0 {
		userDataTemplate = req.UserDataTemplate
	}

	userData, err := renderUserData(fqEnvironmentID, e.Config.S3Bucket(), userDataTemplate)
	if err != nil {
		return nil, err
	}

	securityGroupName := getEnvironmentSGName(fqEnvironmentID)
	if err := createSG(
		e.AWS.EC2,
		securityGroupName,
		fmt.Sprintf("SG for Layer0 environment %s", environmentID),
		e.Config.VPC()); err != nil {
		return nil, err
	}

	securityGroup, err := readSG(e.AWS.EC2, securityGroupName)
	if err != nil {
		return nil, err
	}

	groupID := aws.StringValue(securityGroup.GroupId)
	if err := e.authorizeSGIngress(groupID); err != nil {
		return nil, err
	}

	launchConfigName := fqEnvironmentID
	if err := e.createLC(
		launchConfigName,
		aws.StringValue(securityGroup.GroupId),
		instanceType,
		e.Config.InstanceProfile(),
		amiID,
		userData); err != nil {
		return nil, err
	}

	autoScalingGroupName := fqEnvironmentID
	if err := e.createASG(
		autoScalingGroupName,
		launchConfigName,
		int64(req.MinClusterCount),
		e.Config.PrivateSubnets()); err != nil {
		return nil, err
	}

	clusterName := fqEnvironmentID
	if err := e.createCluster(clusterName); err != nil {
		return nil, err
	}

	if err := e.createTags(environmentID, req.EnvironmentName, req.OperatingSystem); err != nil {
		return nil, err
	}

	return e.Read(environmentID)
}

func (e *EnvironmentProvider) authorizeSGIngress(groupID string) error {
	groupPair := &ec2.UserIdGroupPair{}
	groupPair.SetGroupId(groupID)

	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetGroupId(groupID)

	ingressInput.SetIpPermissions([]*ec2.IpPermission{permission})

	if _, err := e.AWS.EC2.AuthorizeSecurityGroupIngress(ingressInput); err != nil {
		return err
	}

	return nil
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

func (e *EnvironmentProvider) createASG(autoScalingGroupName, launchConfigName string, minSize int64, privateSubnets []string) error {
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
	input.SetMinSize(minSize)
	input.SetMaxSize(minSize)
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
