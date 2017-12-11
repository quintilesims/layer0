package aws

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/common/models"
)

// Update is used to update an ECS Cluster using the specified Update Environment
// Request. The Update Environment Request contains the Environment ID and the
// minimum size of the Cluster's Auto Scaling Group. The Cluster's Auto Scaling
// Group size is updated by making an UpdateAutoScalingGroup request to AWS.
func (e *EnvironmentProvider) Update(environmentID string, req models.UpdateEnvironmentRequest) error {
	fqEnvironmentID := addLayer0Prefix(e.Config.Instance(), environmentID)

	if req.MinScale != nil || req.MaxScale != nil {
		autoScalingGroupName := fqEnvironmentID
		asg, err := e.readASG(autoScalingGroupName)
		if err != nil {
			return err
		}

		minSize := aws.Int64Value(asg.MinSize)
		if req.MinScale != nil {
			minSize = int64(*req.MinScale)
		}

		maxSize := aws.Int64Value(asg.MaxSize)
		if req.MaxScale != nil {
			maxSize = int64(*req.MaxScale)
		}

		if err := e.updateASGSize(autoScalingGroupName, minSize, maxSize); err != nil {
			return err
		}
	}

	if req.Links != nil {
		sourceEnvSG, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqEnvironmentID))
		if err != nil {
			return err
		}

		// check which environments are currently linked
		currentTags, err := e.TagStore.SelectByTypeAndID("environment", environmentID)
		if err != nil {
			return err
		}

		currentLinkedEnvIDs := []string{}
		if tag, ok := currentTags.WithKey("link").First(); ok {
			currentLinkedEnvIDs = strings.Split(tag.Value, ",")
		}

		// add - check if current link state contains links in request
		for _, destEnvironmentID := range *req.Links {
			actualLinks := models.LinkTags(currentLinkedEnvIDs)
			if actualLinks.Contains(destEnvironmentID) {
				continue
			}

			fqDestEnvID := addLayer0Prefix(e.Config.Instance(), destEnvironmentID)
			sg, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
			if err != nil {
				return err
			}

			sourceSGID := aws.StringValue(sourceEnvSG.GroupId)
			destSGID := aws.StringValue(sg.GroupId)
			if err := e.createIngressInput(sourceSGID, destSGID); err != nil {
				return err
			}
		}

		// remove - check if links in request are missing from current links
		for _, destEnvironmentID := range currentLinkedEnvIDs {
			desiredLinks := models.LinkTags(*req.Links)
			if desiredLinks.Contains(destEnvironmentID) {
				continue
			}

			fqDestEnvID := addLayer0Prefix(e.Config.Instance(), destEnvironmentID)
			sg, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "DoesNotExist" {
					log.Printf("[WARN] skipping environment unlink since security group '%s' does not exist\n", getEnvironmentSGName(fqEnvironmentID))
					continue
				}

				return err
			}

			sourceSGID := aws.StringValue(sourceEnvSG.GroupId)
			destSGID := aws.StringValue(sg.GroupId)
			if err := e.removeIngressRule(sourceSGID, destSGID); err != nil {
				return err
			}
		}

		if err := e.setLinkTags(environmentID, *req.Links); err != nil {
			return err
		}
	}

	return nil
}

func (e *EnvironmentProvider) updateASGSize(autoScalingGroupName string, minSize, maxSize int64) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{}
	input.SetAutoScalingGroupName(autoScalingGroupName)
	input.SetMinSize(minSize)
	input.SetMaxSize(maxSize)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := e.AWS.AutoScaling.UpdateAutoScalingGroup(input); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createIngressInput(sourceGroupID, destGroupID string) error {
	groupPair := &ec2.UserIdGroupPair{}
	groupPair.SetGroupId(destGroupID)

	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetGroupId(sourceGroupID)

	ingressInput.SetIpPermissions([]*ec2.IpPermission{permission})

	if _, err := e.AWS.EC2.AuthorizeSecurityGroupIngress(ingressInput); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "InvalidPermission.Duplicate" {
			log.Printf("[WARN] skipping environment link since rule for security group '%s' already exists\n", destGroupID)
			return nil
		}

		return err
	}

	return nil
}

func (e *EnvironmentProvider) setLinkTags(environmentID string, links []string) error {
	l := models.LinkTags(links)
	newTags := models.Tags{
		{
			EntityID:   environmentID,
			EntityType: "environment",
			Key:        "link",
			Value:      strings.Join(l.Distinct(), ","),
		},
	}

	for _, tag := range newTags {
		if err := e.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}

func (e *EnvironmentProvider) removeIngressRule(groupID, groupIDToRemove string) error {
	groupPair := &ec2.UserIdGroupPair{}
	groupPair.SetGroupId(groupIDToRemove)

	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	input := &ec2.RevokeSecurityGroupIngressInput{}
	input.SetGroupId(groupID)
	input.SetIpPermissions([]*ec2.IpPermission{permission})

	if _, err := e.AWS.EC2.RevokeSecurityGroupIngress(input); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "InvalidPermission.NotFound" {
			log.Println("[DEBUG] skipping ingressRule deletion as the rule doesn't seem to exist")
			return nil
		}

		return err
	}

	return nil
}
