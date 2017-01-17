package ecsbackend

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
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

func (this *ECSEnvironmentManager) ListEnvironments() ([]*models.Environment, error) {
	clusters, err := this.ECS.Helper_DescribeClusters()
	if err != nil {
		return nil, err
	}

	environments := make([]*models.Environment, len(clusters))
	for i, cluster := range clusters {
		if strings.HasPrefix(*cluster.ClusterName, id.PREFIX) {
			ecsEnvironmentID := id.ECSEnvironmentID(*cluster.ClusterName)
			environments[i] = &models.Environment{
				EnvironmentID: ecsEnvironmentID.L0EnvironmentID(),
			}
		}
	}

	return environments, nil
}

func (this *ECSEnvironmentManager) GetEnvironment(environmentID string) (*models.Environment, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	cluster, err := this.ECS.DescribeCluster(ecsEnvironmentID.String())
	if err != nil {
		if ContainsErrCode(err, "ClusterNotFoundException") {
			return nil, errors.Newf(errors.InvalidEnvironmentID, "Environment with id '%s' was not found", environmentID)
		}

		return nil, err
	}

	return this.populateModel(cluster)
}

func (this *ECSEnvironmentManager) populateModel(cluster *ecs.Cluster) (*models.Environment, error) {
	// assuming id.ECSEnvironmentID == ECSEnvironmentID.ClusterName()
	ecsEnvironmentID := id.ECSEnvironmentID(*cluster.ClusterName)

	var clusterCount int
	var instanceSize string

	asg, err := this.describeAutoscalingGroup(ecsEnvironmentID)
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
			launchConfig, err := this.AutoScaling.DescribeLaunchConfiguration(*asg.LaunchConfigurationName)
			if err != nil {
				if ContainsErrMsg(err, "not found") {
					log.Errorf("Launch Config for environment '%s' not found", ecsEnvironmentID)
				} else {
					return nil, err
				}
			}

			if launchConfig != nil {
				instanceSize = *launchConfig.InstanceType
			}
		}
	}

	var securityGroupID string
	securityGroup, err := this.EC2.DescribeSecurityGroup(ecsEnvironmentID.SecurityGroupName())
	if err != nil {
		return nil, err
	}

	if securityGroup != nil {
		securityGroupID = *securityGroup.SecurityGroup.GroupId
	}
	model := &models.Environment{
		EnvironmentID:   ecsEnvironmentID.L0EnvironmentID(),
		ClusterCount:    clusterCount,
		InstanceSize:    instanceSize,
		SecurityGroupID: securityGroupID,
	}

	return model, nil
}

func (this *ECSEnvironmentManager) describeAutoscalingGroup(ecsEnvironmentID id.ECSEnvironmentID) (*autoscaling.Group, error) {
	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
	asg, err := this.AutoScaling.DescribeAutoScalingGroup(autoScalingGroupName)
	if err != nil {
		return nil, err
	}

	return asg, nil
}

func (this *ECSEnvironmentManager) CreateEnvironment(
	environmentName string,
	instanceSize string,
	minClusterCount int,
	userDataTemplate []byte,
) (*models.Environment, error) {
	environmentID := id.GenerateHashedEntityID(environmentName)
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	if len(userDataTemplate) == 0 {
		userDataTemplate = defaultUserDataTemplate
	}

	userData, err := renderUserData(ecsEnvironmentID, userDataTemplate)
	if err != nil {
		return nil, err
	}

	cluster, err := this.ECS.CreateCluster(ecsEnvironmentID.String())
	if err != nil {
		return nil, err
	}

	description := "Auto-generated Layer0 Environment Security Group"
	vpcID := config.AWSVPCID()

	groupID, err := this.EC2.CreateSecurityGroup(ecsEnvironmentID.SecurityGroupName(), description, vpcID)
	if err != nil {
		return nil, err
	}

	// wait for security group to propagate
	this.Clock.Sleep(time.Second * 2)
	if err := this.EC2.AuthorizeSecurityGroupIngressFromGroup(groupID, groupID); err != nil {
		return nil, err
	}

	agentGroupID := config.AWSAgentGroupID()
	securityGroups := []*string{groupID, &agentGroupID}
	serviceAMI := config.AWSServiceAMI()
	ecsRole := config.AWSECSInstanceProfile()
	keyPair := config.AWSKeyPair()
	launchConfigurationName := ecsEnvironmentID.LaunchConfigurationName()

	if err := this.AutoScaling.CreateLaunchConfiguration(
		&launchConfigurationName,
		&serviceAMI,
		&ecsRole,
		&instanceSize,
		&keyPair,
		&userData,
		securityGroups,
	); err != nil {
		return nil, err
	}

	maxClusterCount := 0
	if minClusterCount > 0 {
		maxClusterCount = minClusterCount
	}

	if err := this.AutoScaling.CreateAutoScalingGroup(
		ecsEnvironmentID.AutoScalingGroupName(),
		launchConfigurationName,
		config.AWSPrivateSubnets(),
		minClusterCount,
		maxClusterCount,
	); err != nil {
		return nil, err
	}

	return this.populateModel(cluster)
}

func (this *ECSEnvironmentManager) UpdateEnvironment(environmentID string, minClusterCount int) (*models.Environment, error) {
	model, err := this.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}

	if err := this.updateEnvironmentMinCount(model, minClusterCount); err != nil {
		return nil, err
	}

	return model, nil
}

func (this *ECSEnvironmentManager) updateEnvironmentMinCount(model *models.Environment, minClusterCount int) error {
	ecsEnvironmentID := id.L0EnvironmentID(model.EnvironmentID).ECSEnvironmentID()
	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()

	asg, err := this.describeAutoscalingGroup(ecsEnvironmentID)
	if err != nil {
		return err
	}

	if int(*asg.MaxSize) < minClusterCount {
		if err := this.AutoScaling.UpdateAutoScalingGroupMaxSize(autoScalingGroupName, minClusterCount); err != nil {
			return err
		}
	}

	if err := this.AutoScaling.UpdateAutoScalingGroupMinSize(autoScalingGroupName, minClusterCount); err != nil {
		return err
	}

	return nil
}

func (this *ECSEnvironmentManager) DeleteEnvironment(environmentID string) error {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
	if err := this.AutoScaling.UpdateAutoScalingGroupMinSize(autoScalingGroupName, 0); err != nil {
		if !ContainsErrMsg(err, "name not found") && !ContainsErrMsg(err, "is pending delete") {
			return err
		}
	}

	if err := this.AutoScaling.UpdateAutoScalingGroupMaxSize(autoScalingGroupName, 0); err != nil {
		if !ContainsErrMsg(err, "name not found") && !ContainsErrMsg(err, "is pending delete") {
			return err
		}
	}

	if err := this.AutoScaling.DeleteAutoScalingGroup(&autoScalingGroupName); err != nil {
		if !ContainsErrMsg(err, "name not found") {
			return err
		}
	}

	launchConfigurationName := ecsEnvironmentID.LaunchConfigurationName()
	if err := this.AutoScaling.DeleteLaunchConfiguration(&launchConfigurationName); err != nil {
		if !ContainsErrMsg(err, "name not found") {
			return err
		}
	}

	if err := this.waitForAutoScalingGroupInactive(ecsEnvironmentID); err != nil {
		return err
	}

	securityGroup, err := this.EC2.DescribeSecurityGroup(ecsEnvironmentID.SecurityGroupName())
	if err != nil {
		return err
	}

	if securityGroup != nil {
		if err := this.waitForSecurityGroupDeleted(securityGroup); err != nil {
			return err
		}
	}

	if err := this.ECS.DeleteCluster(ecsEnvironmentID.String()); err != nil {
		if !ContainsErrCode(err, "ClusterNotFoundException") {
			return err
		}
	}

	return nil
}

func (this *ECSEnvironmentManager) waitForAutoScalingGroupInactive(ecsEnvironmentID id.ECSEnvironmentID) error {
	autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()

	check := func() (bool, error) {
		group, err := this.AutoScaling.DescribeAutoScalingGroup(autoScalingGroupName)
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
		Retries: 30,
		Delay:   time.Second * 10,
		Clock:   this.Clock,
		Check:   check,
	}

	return waiter.Wait()
}

func (this *ECSEnvironmentManager) waitForSecurityGroupDeleted(securityGroup *ec2.SecurityGroup) error {
	check := func() (bool, error) {
		if err := this.EC2.DeleteSecurityGroup(securityGroup); err == nil {
			return true, nil
		}

		return false, nil
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("SecurityGroup delete for '%v'", securityGroup),
		Retries: 30,
		Delay:   time.Second * 10,
		Clock:   this.Clock,
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

var defaultUserDataTemplate = []byte(
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
