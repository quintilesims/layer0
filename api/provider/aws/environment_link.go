package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Link(req models.CreateEnvironmentLinkRequest) error {
	fqSourceEnvID := addLayer0Prefix(e.Config.Instance(), req.SourceEnvironmentID)
	fqDestEnvID := addLayer0Prefix(e.Config.Instance(), req.DestEnvironmentID)

	sourceGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqSourceEnvID))
	if err != nil {
		return err
	}

	destGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
	if err != nil {
		return err
	}

	sourceGroupID := sourceGroup.GroupId
	destGroupID := destGroup.GroupId

	if _, err := e.createIngressInput(sourceGroupID, destGroupID); err != nil {
		return err
	}

	if _, err := e.createIngressInput(destGroupID, sourceGroupID); err != nil {
		return err
	}

	if err := e.createLinkTags(req.SourceEnvironmentID, req.DestEnvironmentID); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createIngressInput(sourceGroupID, destGroupID *string) *ec2.AuthorizeSecurityGroupIngressInput {
	groupPair := &ec2.UserIdGroupPair{}
	groupPair.SetGroupId(aws.StringValue(destGroupID))

	permission := &ec2.IpPermission{}
	permission.SetIpProtocol("-1")
	permission.SetUserIdGroupPairs([]*ec2.UserIdGroupPair{groupPair})

	ingressInput := &ec2.AuthorizeSecurityGroupIngressInput{}
	ingressInput.SetGroupId(aws.StringValue(sourceGroupID))

	ingressInput.SetIpPermissions([]*ec2.IpPermission{permission})

	if _, err := e.AWS.EC2.AuthorizeSecurityGroupIngress(ingressInput); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createLinkTags(sourceEnvironmentID, destEnvironmentID string) error {
	tags := models.Tags{
		{
			EntityID:   sourceEnvironmentID,
			EntityType: "environment",
			Key:        "link",
			Value:      destEnvironmentID,
		},
		{
			EntityID:   destEnvironmentID,
			EntityType: "environment",
			Key:        "link",
			Value:      sourceEnvironmentID,
		},
	}

	for _, tag := range tags {
		if err := e.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}
