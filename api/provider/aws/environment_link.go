package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

	return nil
}

func (e *EnvironmentProvider) UnLink(sourceEnvironmentID, destEnvironmentID string) error {
	// fqSourceEnvID := addLayer0Prefix(e.Config.Instance(), sourceEnvironmentID)
	// fqDestEnvID := addLayer0Prefix(e.Config.Instance(), destEnvironmentID)

	// sourceGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqSourceEnvID))
	// if err != nil {
	// 	return err
	// }

	// if sourceGroup == nil {
	// 	log.Warnf("Skipping environment unlink since security group '%s' does not exist", getEnvironmentSGName(fqDestEnvID))
	// 	return nil
	// }

	// destGroup, err := readSG(e.AWS.EC2, getEnvironmentSGName(fqDestEnvID))
	// if err != nil {
	// 	return err
	// }

	// if destGroup == nil {
	// 	log.Warnf("Skipping environment unlink since security group '%s' does not exist", getEnvironmentSGName(fqDestEnvID))
	// 	return nil
	// }

	// removeIngressRule := func(group *ec2.SecurityGroup, groupIDToRemove string) error {
	// 	for _, permission := range group.IpPermissions {
	// 		for _, pair := range permission.UserIdGroupPairs {
	// 			if *pair.GroupId == groupIDToRemove {
	// 				groupPermission := ec2.IpPermission{
	// 					&awsec2.IpPermission{
	// 						IpProtocol:       permission.IpProtocol,
	// 						UserIdGroupPairs: []*awsec2.UserIdGroupPair{pair},
	// 					},
	// 				}

	// 				if err := e.EC2.RevokeSecurityGroupIngressHelper(*group.GroupId, groupPermission); err != nil {
	// 					return err
	// 				}
	// 			}
	// 		}
	// 	}

	// 	return nil
	// }

	// if err := removeIngressRule(sourceGroup, *destGroup.GroupId); err != nil {
	// 	return err
	// }

	// if err := removeIngressRule(destGroup, *sourceGroup.GroupId); err != nil {
	// 	return err
	// }

	return nil
}
