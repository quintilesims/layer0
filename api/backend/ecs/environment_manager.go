package ecsbackend

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/waitutils"
)

type ECSEnvironmentManager struct {
	ECS         ecs.Provider
	EC2         ec2.Provider
	AutoScaling autoscaling.Provider
	Backend     backend.Backend
	Clock       waitutils.Clock
}

func NewECSEnvironmentManager(
	ecsprovider ecs.Provider,
	ec2 ec2.Provider,
	asg autoscaling.Provider,
	backend backend.Backend) *ECSEnvironmentManager {

	return &ECSEnvironmentManager{
		ECS:         ecsprovider,
		EC2:         ec2,
		AutoScaling: asg,
		Backend:     backend,
		Clock:       waitutils.RealClock{},
	}
}

func (e *ECSEnvironmentManager) ListEnvironments() ([]id.ECSEnvironmentID, error) {
	clusterNames, err := e.ECS.ListClusterNames(id.PREFIX)
	if err != nil {
		return nil, err
	}

	ecsEnvironmentIDs := make([]id.ECSEnvironmentID, len(clusterNames))
	for i, clusterName := range clusterNames {
		ecsEnvironmentIDs[i] = id.ECSEnvironmentID(clusterName)
	}

	return ecsEnvironmentIDs, nil
}

func (e *ECSEnvironmentManager) GetEnvironment(environmentID string) (*models.Environment, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	cluster, err := e.ECS.DescribeCluster(ecsEnvironmentID.String())
	if err != nil {
		if ContainsErrCode(err, "ClusterNotFoundException") || ContainsErrMsg(err, "cluster not found") {
			return nil, errors.Newf(errors.EnvironmentDoesNotExist, "Environment with id '%s' does not exist", environmentID)
		}

		return nil, err
	}

	return e.populateModel(cluster)
}

func (e *ECSEnvironmentManager) populateModel(cluster *ecs.Cluster) (*models.Environment, error) {
	// assuming id.ECSEnvironmentID == ECSEnvironmentID.ClusterName()
	ecsEnvironmentID := id.ECSEnvironmentID(*cluster.ClusterName)

	var clusterCount int
	var instanceSize string
	var amiID string

	asg, err := e.describeAutoscalingGroup(ecsEnvironmentID)
	if err != nil {
		if ContainsErrMsg(err, "not found") {
			log.Errorf("Autoscaling Group for environment '%s' not found", ecsEnvironmentID)
		} else {
			return nil, err
		}
	}

	if asg != nil {
		clusterCount = len(asg.Instances)

		if asg.LaunchConfigurationName != nil {
			launchConfig, err := e.AutoScaling.DescribeLaunchConfiguration(*asg.LaunchConfigurationName)
			if err != nil {
				if ContainsErrMsg(err, "not found") {
					log.Errorf("Launch Config for environment '%s' not found", ecsEnvironmentID)
				} else {
					return nil, err
				}
			}

			if launchConfig != nil {
				instanceSize = *launchConfig.InstanceType
				amiID = *launchConfig.ImageId
			}
		}
	}

	var securityGroupID string
	securityGroup, err := e.EC2.DescribeSecurityGroup(ecsEnvironmentID.SecurityGroupName())
	if err != nil {
		return nil, err
	}

	if securityGroup != nil {
		securityGroupID = pstring(securityGroup.GroupId)
	}

	model := &models.Environment{
		EnvironmentID:   ecsEnvironmentID.L0EnvironmentID(),
		ClusterCount:    clusterCount,
		InstanceSize:    instanceSize,
		SecurityGroupID: securityGroupID,
		AMIID:           amiID,
		MinCount:        int(*asg.MinSize),
		MaxCount:        int(*asg.MaxSize),
	}

	return model, nil
}

func (e *ECSEnvironmentManager) describeAutoscalingGroup(ecsEnvironmentID id.ECSEnvironmentID) (*autoscaling.Group, error) {
	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
	asg, err := e.AutoScaling.DescribeAutoScalingGroup(autoScalingGroupName)
	if err != nil {
		return nil, err
	}

	return asg, nil
}

func (e *ECSEnvironmentManager) CreateEnvironment(
	environmentName string,
	instanceSize string,
	operatingSystem string,
	amiID string,
	minClusterCount int,
	maxClusterCount int,
	targetCapSize int,
	userDataTemplate []byte,
) (*models.Environment, error) {

	var defaultUserDataTemplate []byte
	var serviceAMI string
	switch strings.ToLower(operatingSystem) {
	case "linux":
		defaultUserDataTemplate = defaultLinuxUserDataTemplate
		serviceAMI = config.AWSLinuxServiceAMI()
	case "windows":
		defaultUserDataTemplate = defaultWindowsUserDataTemplate
		serviceAMI = config.AWSWindowsServiceAMI()
	default:
		return nil, fmt.Errorf("Operating system '%s' is not recognized", operatingSystem)
	}

	environmentID := id.GenerateHashedEntityID(environmentName)
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	if len(userDataTemplate) == 0 {
		userDataTemplate = defaultUserDataTemplate
	}

	if amiID != "" {
		serviceAMI = amiID
	}

	userData, err := renderUserData(ecsEnvironmentID, userDataTemplate)
	if err != nil {
		return nil, err
	}

	description := "Auto-generated Layer0 Environment Security Group"
	vpcID := config.AWSVPCID()

	groupID, err := e.EC2.CreateSecurityGroup(ecsEnvironmentID.SecurityGroupName(), description, vpcID)
	if err != nil {
		return nil, err
	}

	// wait for security group to propagate
	e.Clock.Sleep(time.Second * 2)
	if err := e.EC2.AuthorizeSecurityGroupIngressFromGroup(*groupID, *groupID); err != nil {
		return nil, err
	}

	securityGroups := []*string{groupID}
	ecsRole := config.AWSECSInstanceProfile()
	keyPair := config.AWSKeyPair()
	launchConfigurationName := ecsEnvironmentID.LaunchConfigurationName()
	volSizes := make(map[string]int)
	if operatingSystem == "linux" {	
		volSizes["/dev/xvda"] = 30;	
	} else {
		volSizes["/dev/sda1"] = 200
	}

	if err := e.AutoScaling.CreateLaunchConfiguration(
		&launchConfigurationName,
		&serviceAMI,
		&ecsRole,
		&instanceSize,
		&keyPair,
		&userData,
		securityGroups,
		volSizes,
	); err != nil {
		return nil, err
	}

	//set the default value
	if minClusterCount > maxClusterCount {
		maxClusterCount = minClusterCount
	}
	if targetCapSize == 0 {
		targetCapSize = 100
	}

	if err := e.AutoScaling.CreateAutoScalingGroup(
		ecsEnvironmentID.AutoScalingGroupName(),
		launchConfigurationName,
		config.AWSPrivateSubnets(),
		minClusterCount,
		maxClusterCount,
	); err != nil {
		return nil, err
	}

	asg, err := e.describeAutoscalingGroup(ecsEnvironmentID)
	cluster, err := e.ECS.CreateCluster(ecsEnvironmentID.String(), *asg.AutoScalingGroupARN, maxClusterCount, minClusterCount, targetCapSize)
	if err != nil {
		return nil, err
	}

	// wait for cluster
	e.Clock.Sleep(time.Second * 40)

	return e.populateModel(cluster)
}

func (e *ECSEnvironmentManager) UpdateEnvironment(environmentID string, minClusterCount int) (*models.Environment, error) {
	model, err := e.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	if err := e.updateEnvironmentMinCount(model, minClusterCount); err != nil {
		return nil, err
	}

	return model, nil
}

func (e *ECSEnvironmentManager) updateEnvironmentMinCount(model *models.Environment, minClusterCount int) error {
	ecsEnvironmentID := id.L0EnvironmentID(model.EnvironmentID).ECSEnvironmentID()
	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()

	asg, err := e.describeAutoscalingGroup(ecsEnvironmentID)
	if err != nil {
		return err
	}

	if int(*asg.MaxSize) < minClusterCount {
		if err := e.AutoScaling.UpdateAutoScalingGroupMaxSize(autoScalingGroupName, minClusterCount); err != nil {
			return err
		}
	}

	if err := e.AutoScaling.UpdateAutoScalingGroupMinSize(autoScalingGroupName, minClusterCount); err != nil {
		return err
	}

	return nil
}

func (e *ECSEnvironmentManager) DeleteEnvironment(environmentID string) error {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
	if err := e.AutoScaling.UpdateAutoScalingGroupMinSize(autoScalingGroupName, 0); err != nil {
		if !ContainsErrMsg(err, "name not found") && !ContainsErrMsg(err, "is pending delete") {
			return err
		}
	}

	if err := e.AutoScaling.UpdateAutoScalingGroupMaxSize(autoScalingGroupName, 0); err != nil {
		if !ContainsErrMsg(err, "name not found") && !ContainsErrMsg(err, "is pending delete") {
			return err
		}
	}

	if err := e.AutoScaling.DeleteAutoScalingGroup(&autoScalingGroupName); err != nil {
		if !ContainsErrMsg(err, "name not found") {
			return err
		}
	}

	launchConfigurationName := ecsEnvironmentID.LaunchConfigurationName()
	if err := e.AutoScaling.DeleteLaunchConfiguration(&launchConfigurationName); err != nil {
		if !ContainsErrMsg(err, "name not found") {
			return err
		}
	}

	if err := e.waitForAutoScalingGroupInactive(ecsEnvironmentID); err != nil {
		return err
	}

	securityGroup, err := e.EC2.DescribeSecurityGroup(ecsEnvironmentID.SecurityGroupName())
	if err != nil {
		return err
	}

	if securityGroup != nil {
		if err := e.waitForSecurityGroupDeleted(securityGroup); err != nil {
			return err
		}
	}

	if err := e.ECS.DeleteCluster(ecsEnvironmentID.String()); err != nil {
		if !ContainsErrCode(err, "ClusterNotFoundException") {
			return err
		}
	}

	return nil
}

func (e *ECSEnvironmentManager) CreateEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error {
	sourceECSID := id.L0EnvironmentID(sourceEnvironmentID).ECSEnvironmentID()
	destECSID := id.L0EnvironmentID(destEnvironmentID).ECSEnvironmentID()

	sourceGroup, err := e.getEnvironmentSecurityGroup(sourceECSID)
	if err != nil {
		return err
	}

	destGroup, err := e.getEnvironmentSecurityGroup(destECSID)
	if err != nil {
		return err
	}

	if err := e.EC2.AuthorizeSecurityGroupIngressFromGroup(*sourceGroup.GroupId, *destGroup.GroupId); err != nil {
		if !ContainsErrCode(err, "InvalidPermission.Duplicate") {
			return err
		}
	}

	if err := e.EC2.AuthorizeSecurityGroupIngressFromGroup(*destGroup.GroupId, *sourceGroup.GroupId); err != nil {
		if !ContainsErrCode(err, "InvalidPermission.Duplicate") {
			return err
		}
	}

	return nil
}

func (e *ECSEnvironmentManager) DeleteEnvironmentLink(sourceEnvironmentID, destEnvironmentID string) error {
	sourceECSID := id.L0EnvironmentID(sourceEnvironmentID).ECSEnvironmentID()
	destECSID := id.L0EnvironmentID(destEnvironmentID).ECSEnvironmentID()

	sourceGroup, err := e.EC2.DescribeSecurityGroup(sourceECSID.SecurityGroupName())
	if err != nil {
		return err
	}

	if sourceGroup == nil {
		log.Warnf("Skipping environment unlink since security group '%s' does not exist", sourceECSID.SecurityGroupName())
		return nil
	}

	destGroup, err := e.EC2.DescribeSecurityGroup(destECSID.SecurityGroupName())
	if err != nil {
		return err
	}

	if destGroup == nil {
		log.Warnf("Skipping environment unlink since security group '%s' does not exist", destECSID.SecurityGroupName())
		return nil
	}

	removeIngressRule := func(group *ec2.SecurityGroup, groupIDToRemove string) error {
		for _, permission := range group.IpPermissions {
			for _, pair := range permission.UserIdGroupPairs {
				if *pair.GroupId == groupIDToRemove {
					groupPermission := ec2.IpPermission{
						&awsec2.IpPermission{
							IpProtocol:       permission.IpProtocol,
							UserIdGroupPairs: []*awsec2.UserIdGroupPair{pair},
						},
					}

					if err := e.EC2.RevokeSecurityGroupIngressHelper(*group.GroupId, groupPermission); err != nil {
						return err
					}
				}
			}
		}

		return nil
	}

	if err := removeIngressRule(sourceGroup, *destGroup.GroupId); err != nil {
		return err
	}

	if err := removeIngressRule(destGroup, *sourceGroup.GroupId); err != nil {
		return err
	}

	return nil
}

func (e *ECSEnvironmentManager) getEnvironmentSecurityGroup(environmentID id.ECSEnvironmentID) (*ec2.SecurityGroup, error) {
	group, err := e.EC2.DescribeSecurityGroup(environmentID.SecurityGroupName())
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, fmt.Errorf("Security group for environment '%s' does not exist", environmentID.L0EnvironmentID())
	}

	return group, nil
}

func (e *ECSEnvironmentManager) waitForAutoScalingGroupInactive(ecsEnvironmentID id.ECSEnvironmentID) error {
	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()

	check := func() (bool, error) {
		group, err := e.AutoScaling.DescribeAutoScalingGroup(autoScalingGroupName)
		if err != nil {
			if ContainsErrMsg(err, "not found") {
				return true, nil
			}

			return false, err
		}

		log.Debugf("Waiting for ASG %s to delete (status: '%s')", autoScalingGroupName, pstring(group.Status))
		return false, nil
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("Stop Autoscaling %s", autoScalingGroupName),
		Retries: 50,
		Delay:   time.Second * 10,
		Clock:   e.Clock,
		Check:   check,
	}

	return waiter.Wait()
}

func (e *ECSEnvironmentManager) waitForSecurityGroupDeleted(securityGroup *ec2.SecurityGroup) error {
	check := func() (bool, error) {
		if err := e.EC2.DeleteSecurityGroup(securityGroup); err == nil {
			return true, nil
		}

		return false, nil
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("SecurityGroup delete for '%v'", securityGroup),
		Retries: 50,
		Delay:   time.Second * 10,
		Clock:   e.Clock,
		Check:   check,
	}

	return waiter.Wait()
}

func renderUserData(ecsEnvironmentID id.ECSEnvironmentID, userData []byte) (string, error) {
	tmpl, err := template.New("").Parse(string(userData))
	if err != nil {
		return "", fmt.Errorf("Failed to parse user data: %v", err)
	}

	context := struct {
		ECSEnvironmentID string
		S3Bucket         string
	}{
		ECSEnvironmentID: ecsEnvironmentID.String(),
		S3Bucket:         config.AWSS3Bucket(),
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, context); err != nil {
		return "", fmt.Errorf("Failed to render user data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(rendered.Bytes()), nil
}

var defaultWindowsUserDataTemplate = []byte(
	`<powershell>
# Set agent env variables for the Machine context (durable)
$clusterName = "{{ .ECSEnvironmentID }}"
Write-Host Cluster name set as: $clusterName -foreground green

[Environment]::SetEnvironmentVariable("ECS_CLUSTER", $clusterName, "Machine")
[Environment]::SetEnvironmentVariable("ECS_ENABLE_TASK_IAM_ROLE", "false", "Machine")
$agentVersion = 'v1.15.2'
$agentZipUri = "https://s3.amazonaws.com/amazon-ecs-agent/ecs-agent-windows-$agentVersion.zip"
$agentZipMD5Uri = "$agentZipUri.md5"

# Configure docker auth
Read-S3Object -BucketName {{ .S3Bucket }} -Key bootstrap/dockercfg -File dockercfg.json
$dockercfgContent = [IO.File]::ReadAllText("dockercfg.json")
[Environment]::SetEnvironmentVariable("ECS_ENGINE_AUTH_DATA", $dockercfgContent, "Machine")
[Environment]::SetEnvironmentVariable("ECS_ENGINE_AUTH_TYPE", "dockercfg", "Machine")

### --- Nothing user configurable after this point ---
$ecsExeDir = "$env:ProgramFiles\Amazon\ECS"
$zipFile = "$env:TEMP\ecs-agent.zip"
$md5File = "$env:TEMP\ecs-agent.zip.md5"

### Get the files from S3
Invoke-RestMethod -OutFile $zipFile -Uri $agentZipUri
Invoke-RestMethod -OutFile $md5File -Uri $agentZipMD5Uri

## MD5 Checksum
$expectedMD5 = (Get-Content $md5File)
$md5 = New-Object -TypeName System.Security.Cryptography.MD5CryptoServiceProvider
$actualMD5 = [System.BitConverter]::ToString($md5.ComputeHash([System.IO.File]::ReadAllBytes($zipFile))).replace('-', '')

if($expectedMD5 -ne $actualMD5) {
    echo "Download doesn't match hash."
    echo "Expected: $expectedMD5 - Got: $actualMD5"
    exit 1
}

## Put the executables in the executable directory.
Expand-Archive -Path $zipFile -DestinationPath $ecsExeDir -Force

## Start the agent script in the background.
$jobname = "ECS-Agent-Init"
$script =  "cd '$ecsExeDir'; .\amazon-ecs-agent.ps1"
$repeat = (New-TimeSpan -Minutes 1)

$jobpath = $env:LOCALAPPDATA + "\Microsoft\Windows\PowerShell\ScheduledJobs\$jobname\ScheduledJobDefinition.xml"
if($(Test-Path -Path $jobpath)) {
  echo "Job definition already present"
  exit 0

}

$scriptblock = [scriptblock]::Create("$script")
$trigger = New-JobTrigger -At (Get-Date).Date -RepeatIndefinitely -RepetitionInterval $repeat -Once
$options = New-ScheduledJobOption -RunElevated -ContinueIfGoingOnBattery -StartIfOnBattery
Register-ScheduledJob -Name $jobname -ScriptBlock $scriptblock -Trigger $trigger -ScheduledJobOption $options -RunNow
Add-JobTrigger -Name $jobname -Trigger (New-JobTrigger -AtStartup -RandomDelay 00:1:00)
</powershell>
<persist>true</persist>
`)

var defaultLinuxUserDataTemplate = []byte(
	`#!/bin/bash
    echo ECS_CLUSTER={{ .ECSEnvironmentID }} >> /etc/ecs/ecs.config
    echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
    yum install -y aws-cli awslogs jq
    aws s3 cp s3://{{ .S3Bucket }}/bootstrap/dockercfg dockercfg
    cfg=$(cat dockercfg)
    echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
    docker pull amazon/amazon-ecs-agent:latest
    start ecs
`)
