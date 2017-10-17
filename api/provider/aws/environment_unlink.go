package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Unlink(req models.DeleteEnvironmentLinkRequest) error {
	fmt.Println(req)

	fqSourceEnvID := addLayer0Prefix(e.Config.Instance(), req.SourceEnvironmentID)
	fqDestEnvID := addLayer0Prefix(e.Config.Instance(), req.DestEnvironmentID)

	sourceGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqSourceEnvID))
	if err != nil {
		//todo: not sure if we want to still log these warning messages here
		//log.Warnf("Skipping environment unlink since security group '%s' does not exist", getEnvironmentSGName(fqDestEnvID))
		return err
	}

	destGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
	if err != nil {
		//log.Warnf("Skipping environment unlink since security group '%s' does not exist", getEnvironmentSGName(fqDestEnvID))
		return err
	}

	if err := e.removeIngressRule(sourceGroup, *destGroup.GroupId); err != nil {
		return err
	}

	if err := e.removeIngressRule(destGroup, *sourceGroup.GroupId); err != nil {
		return err
	}

	if err := e.deleteLinkTag(req.SourceEnvironmentID, req.DestEnvironmentID); err != nil {
		return err
	}

	if err := e.deleteLinkTag(req.DestEnvironmentID, req.SourceEnvironmentID); err != nil {
		return err
	}

	return nil
}

func (e *EnvironmentProvider) removeIngressRule(group *ec2.SecurityGroup, groupIDToRemove string) error {
	for _, permission := range group.IpPermissions {
		for _, pair := range permission.UserIdGroupPairs {
			if aws.StringValue(pair.GroupId) == groupIDToRemove {
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

func (e *EnvironmentProvider) deleteLinkTag(sourceEnvironmentID, destEnvironmentID string) error {
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
