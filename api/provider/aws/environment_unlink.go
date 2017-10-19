package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/quintilesims/layer0/common/models"
)

func (e *EnvironmentProvider) Unlink(req models.DeleteEnvironmentLinkRequest) error {
	log.Printf("[DEBUG] unlink called")
	if err := models.EnvironmentLinkRequest(req).Validate(); err != nil {
		return err
	}

	fqSourceEnvID := addLayer0Prefix(e.Config.Instance(), req.SourceEnvironmentID)
	fqDestEnvID := addLayer0Prefix(e.Config.Instance(), req.DestEnvironmentID)

	sourceSecurityGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqSourceEnvID))
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "DoesNotExist" {
			log.Printf("[WARN] skipping environment unlink since security group '%s' does not exist\n", getEnvironmentSGName(fqSourceEnvID))
			return nil
		}

		return err
	}

	destSecurityGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "DoesNotExist" {
			log.Printf("[WARN] skipping environment unlink since security group '%s' does not exist\n", getEnvironmentSGName(fqDestEnvID))
			return nil
		}

		return err
	}

	sourceSecurityGroupID := aws.StringValue(sourceSecurityGroup.GroupId)
	destSecurityGroupID := aws.StringValue(destSecurityGroup.GroupId)

	if err := e.removeIngressRule(sourceSecurityGroup, destSecurityGroupID); err != nil {
		return err
	}

	if err := e.removeIngressRule(destSecurityGroup, sourceSecurityGroupID); err != nil {
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
				if _, err := e.AWS.EC2.RevokeSecurityGroupIngress(input); err != nil {
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
