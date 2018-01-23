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
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

// Create is used to create an ECS Cluster using the specified Create Environment
// Request. The Create Environment Request contains the name of the Environment,
// the instance type and user data for the Launch Configuration, the minimum size
// of the Cluster's Auto Scaling Group, and the Operating System and EC2 AMI ID
// used in the Launch Configuration. The EC2 Launch Configuration, Auto Scaling
// Group, and Security Group are created before the Cluster is created.
func (e *EnvironmentProvider) Create(req models.CreateEnvironmentRequest) (string, error) {
	// TODO: Ensure environment name is unique
	environmentID := entityIDGenerator(req.EnvironmentName)
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	instanceType := config.DefaultEnvironmentInstanceType
	if req.InstanceType != "" {
		instanceType = req.InstanceType
	}

	var userDataTemplate []byte
	var amiID string

	if req.OperatingSystem == "" {
		req.OperatingSystem = config.DefaultEnvironmentOS
	}

	securityGroupName := getEnvironmentSGName(fqEnvironmentID)
	if err := createSG(
		e.AWS.EC2,
		securityGroupName,
		fmt.Sprintf("SG for Layer0 environment %s", environmentID),
		e.Config.VPC()); err != nil {
		return "", err
	}

	securityGroup, err := readSG(e.AWS.EC2, securityGroupName)
	if err != nil {
		return "", err
	}

	groupID := aws.StringValue(securityGroup.GroupId)
	if err := e.authorizeSGSelfIngress(groupID); err != nil {
		return "", err
	}

	// creating asg, lc isn't required for dynamic environments
	if strings.ToLower(req.EnvironmentType) == models.EnvironmentTypeStatic {
		switch strings.ToLower(req.OperatingSystem) {
		case models.LinuxOS:
			userDataTemplate = []byte(DefaultLinuxUserdataTemplate)
			amiID = e.Config.LinuxAMI()
		case models.WindowsOS:
			userDataTemplate = []byte(DefaultWindowsUserdataTemplate)
			amiID = e.Config.WindowsAMI()
		default:
			return "", fmt.Errorf("Operating system '%s' is not recognized", req.OperatingSystem)
		}

		if req.AMIID != "" {
			amiID = req.AMIID
		}

		if len(req.UserDataTemplate) > 0 {
			userDataTemplate = req.UserDataTemplate
		}

		userData, err := RenderUserData(fqEnvironmentID, e.Config.S3Bucket(), userDataTemplate)
		if err != nil {
			return "", err
		}

		launchConfigName := fqEnvironmentID
		if err := e.createLC(
			launchConfigName,
			aws.StringValue(securityGroup.GroupId),
			instanceType,
			e.Config.InstanceProfile(),
			e.Config.SSHKeyPair(),
			amiID,
			userData); err != nil {
			return "", err
		}

		autoScalingGroupName := fqEnvironmentID
		if err := e.createASG(
			autoScalingGroupName,
			launchConfigName,
			int64(req.Scale),
			int64(req.Scale),
			e.Config.PrivateSubnets()); err != nil {
			return "", err
		}
	}

	clusterName := fqEnvironmentID
	if err := e.createCluster(clusterName); err != nil {
		return "", err
	}

	if err := e.createTags(environmentID, req.EnvironmentName, req.EnvironmentType, req.OperatingSystem); err != nil {
		return "", err
	}

	return environmentID, nil
}

func (e *EnvironmentProvider) authorizeSGSelfIngress(groupID string) error {
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
	keyPairName string,
	amiID string,
	userData string,
) error {
	input := &autoscaling.CreateLaunchConfigurationInput{}
	input.SetLaunchConfigurationName(launchConfigName)
	input.SetSecurityGroups([]*string{aws.String(securityGroupID)})
	input.SetInstanceType(instanceType)
	input.SetIamInstanceProfile(instanceProfile)
	input.SetKeyName(keyPairName)
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

func (e *EnvironmentProvider) createASG(
	autoScalingGroupName string,
	launchConfigName string,
	minSize int64,
	maxSize int64,
	privateSubnets []string,
) error {
	tag := &autoscaling.Tag{}
	tag.SetKey("Name")
	tag.SetValue(autoScalingGroupName)
	tag.SetPropagateAtLaunch(true)

	subnetIdentifier := strings.Join(privateSubnets, ",")

	input := &autoscaling.CreateAutoScalingGroupInput{}
	input.SetAutoScalingGroupName(autoScalingGroupName)
	input.SetLaunchConfigurationName(launchConfigName)
	input.SetVPCZoneIdentifier(subnetIdentifier)
	input.SetMinSize(minSize)
	input.SetMaxSize(maxSize)
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

func (e *EnvironmentProvider) createTags(environmentID, environmentName, environmentType, operatingSystem string) error {
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
		{
			EntityID:   environmentID,
			EntityType: "environment",
			Key:        "type",
			Value:      strings.ToLower(environmentType),
		},
	}

	for _, tag := range tags {
		if err := e.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}

func RenderUserData(environmentID, s3Bucket string, userDataTemplate []byte) (string, error) {
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
