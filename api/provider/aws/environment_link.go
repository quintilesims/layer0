package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Link(sourceEnvironmentID, destEnvironmentID string) error {
	fqSourceEnvID := addLayer0Prefix(e.Config.Instance(), sourceEnvironmentID)
	fqDestEnvID := addLayer0Prefix(e.Config.Instance(), destEnvironmentID)

	sourceGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqSourceEnvID))
	if err != nil {
		return err
	}

	destGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
	if err != nil {
		return err
	}

	createIngressInput := func(sourceGroupID, destGroupID *string) *ec2.AuthorizeSecurityGroupIngressInput {
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

	_, err = e.AWS.EC2.AuthorizeSecurityGroupIngress(createIngressInput(sourceGroup.GroupId, destGroup.GroupId))
	if err != nil {
		return err
	}

	_, err = e.AWS.EC2.AuthorizeSecurityGroupIngress(createIngressInput(destGroup.GroupId, sourceGroup.GroupId))
	if err != nil {
		return err
	}

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

func (e *EnvironmentProvider) Unlink(sourceEnvironmentID, destEnvironmentID string) error {
	fqSourceEnvID := addLayer0Prefix(e.Config.Instance(), sourceEnvironmentID)
	fqDestEnvID := addLayer0Prefix(e.Config.Instance(), destEnvironmentID)

	sourceGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqSourceEnvID))
	if err != nil {
		return err
	}

	if sourceGroup == nil {
		//log.Warnf("Skipping environment unlink since security group '%s' does not exist", getEnvironmentSGName(fqDestEnvID))
		return nil
	}

	destGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
	if err != nil {
		return err
	}

	if destGroup == nil {
		//log.Warnf("Skipping environment unlink since security group '%s' does not exist", getEnvironmentSGName(fqDestEnvID))
		return nil
	}

	removeIngressRule := func(group *ec2.SecurityGroup, groupIDToRemove string) error {
		for _, permission := range group.IpPermissions {
			for _, pair := range permission.UserIdGroupPairs {
				if *pair.GroupId == groupIDToRemove {
					groupPermission := &ec2.IpPermission{
						IpProtocol:       permission.IpProtocol,
						UserIdGroupPairs: []*ec2.UserIdGroupPair{pair},
					}

					input := &ec2.RevokeSecurityGroupIngressInput{
						GroupId:       group.GroupId,
						IpPermissions: []*ec2.IpPermission{groupPermission},
					}
					_, err := e.AWS.EC2.RevokeSecurityGroupIngress(input)
					if err != nil {
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

	deleteLinkTag := func(sourceEnvironmentID, destEnvironmentID string) error {
		tags, err := e.TagStore.SelectByTypeAndID("environment", sourceEnvironmentID)
		if err != nil {
			return err
		}

		for _, tag := range tags.WithKey("link").WithValue(destEnvironmentID) {
			if err := e.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
				return err
			}
		}

		return nil
	}

	if err := deleteLinkTag(sourceEnvironmentID, destEnvironmentID); err != nil {
		return err
	}

	if err := deleteLinkTag(destEnvironmentID, sourceEnvironmentID); err != nil {
		return err
	}

	return nil
}
