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

	ingressInput := e.createIngressInput(sourceGroup.GroupId, destGroup.GroupId)
	if _, err := e.AWS.EC2.AuthorizeSecurityGroupIngress(ingressInput); err != nil {
		return err
	}

	ingressInput = e.createIngressInput(destGroup.GroupId, sourceGroup.GroupId)
	if _, err := e.AWS.EC2.AuthorizeSecurityGroupIngress(ingressInput); err != nil {
		return err
	}

	if err := e.createLinkTags(req.SourceEnvironmentID, req.DestEnvironmentID); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) createIngressInput(sourceGroupID, destGroupID *string) *ec2.AuthorizeSecurityGroupIngressInput {
	return &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: sourceGroupID,
		IpPermissions: []*ec2.IpPermission{
			{
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: destGroupID,
					},
				},
				IpProtocol: aws.String("-1"),
			},
		},
	}
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
