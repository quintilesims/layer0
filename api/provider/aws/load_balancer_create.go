package aws

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/quintilesims/layer0/common/models"
)

func (l *LoadBalancerProvider) Create(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error) {
	loadBalancerID := generateEntityID(req.LoadBalancerName)
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)
	fqEnvironmentID := addLayer0Prefix(l.Config.Instance(), req.EnvironmentID)

	environmentSGName := getEnvironmentSGName(fqEnvironmentID)
	environmentSG, err := readSG(l.AWS.EC2, environmentSGName)
	if err != nil {
		return nil, err
	}

	scheme := "internal"
	subnets := l.Config.PrivateSubnets()
	securityGroupIDs := []string{aws.StringValue(environmentSG.GroupId)}

	if req.IsPublic {
		scheme = "internet-facing"
		subnets = l.Config.PublicSubnets()

		loadBalancerSGName := getLoadBalancerSGName(fqLoadBalancerID)
		if err := createSG(
			l.AWS.EC2,
			loadBalancerSGName,
			fmt.Sprintf("SG for Layer0 load balancer %s", loadBalancerID),
			l.Config.VPC()); err != nil {
			return nil, err
		}

		loadBalancerSG, err := readSG(l.AWS.EC2, loadBalancerSGName)
		if err != nil {
			return nil, err
		}

		loadBalancerSGID := aws.StringValue(loadBalancerSG.GroupId)
		for _, port := range req.Ports {
			if err := l.authorizeSGIngressFromPort(loadBalancerSGID, int64(port.HostPort)); err != nil {
				return nil, err
			}
		}

		securityGroupIDs = append(securityGroupIDs, aws.StringValue(loadBalancerSG.GroupId))
	}

	roleName := getLoadBalancerRoleName(fqLoadBalancerID)
	if _, err := l.createRole(roleName, DEFAULT_ASSUME_ROLE_POLICY); err != nil {
		return nil, err
	}

	policy, err := renderLoadBalancerRolePolicy(
		l.Config.Region(),
		l.Config.AccountID(),
		fqLoadBalancerID,
		DEFAULT_LB_ROLE_POLICY_TEMPLATE)
	if err != nil {
		return nil, err
	}

	policyName := roleName
	if err := l.putRolePolicy(policyName, roleName, policy); err != nil {
		return nil, err
	}

	listeners, err := l.portsToListeners(req.Ports)
	if err != nil {
		return nil, err
	}

	if err := l.createLoadBalancer(
		fqLoadBalancerID,
		scheme,
		securityGroupIDs,
		subnets,
		listeners); err != nil {
		return nil, err
	}

	if err := l.createTags(loadBalancerID, req.LoadBalancerName, req.EnvironmentID); err != nil {
		return nil, err
	}

	return l.Read(loadBalancerID)
}

func (l *LoadBalancerProvider) createRole(roleName, policy string) (*iam.Role, error) {
	input := &iam.CreateRoleInput{}
	input.SetRoleName(roleName)
	input.SetAssumeRolePolicyDocument(policy)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := l.AWS.IAM.CreateRole(input)
	if err != nil {
		return nil, err
	}

	return output.Role, nil
}

func (l *LoadBalancerProvider) putRolePolicy(policyName, roleName, policy string) error {
	input := &iam.PutRolePolicyInput{}
	input.SetPolicyName(policyName)
	input.SetRoleName(roleName)
	input.SetPolicyDocument(policy)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.IAM.PutRolePolicy(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) authorizeSGIngressFromPort(groupID string, port int64) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{}
	input.SetGroupId(groupID)
	input.SetCidrIp("0.0.0.0/0")
	input.SetIpProtocol("TCP")
	input.SetFromPort(port)
	input.SetToPort(port)

	if _, err := l.AWS.EC2.AuthorizeSecurityGroupIngress(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) createLoadBalancer(
	loadBalancerName string,
	scheme string,
	securityGroupIDs []string,
	subnetIDs []string,
	listeners []*elb.Listener,
) error {
	securityGroups := make([]*string, len(securityGroupIDs))
	for i, securityGroupID := range securityGroupIDs {
		securityGroups[i] = aws.String(securityGroupID)
	}

	subnets := make([]*string, len(subnetIDs))
	for i, subnetID := range subnetIDs {
		subnets[i] = aws.String(subnetID)
	}

	input := &elb.CreateLoadBalancerInput{}
	input.SetLoadBalancerName(loadBalancerName)
	input.SetScheme(scheme)
	input.SetSecurityGroups(securityGroups)
	input.SetSubnets(subnets)
	input.SetListeners(listeners)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.ELB.CreateLoadBalancer(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) portsToListeners(ports []models.Port) ([]*elb.Listener, error) {
	listeners := make([]*elb.Listener, len(ports))
	for i, port := range ports {
		listener := &elb.Listener{}
		listener.SetProtocol(port.Protocol)
		listener.SetLoadBalancerPort(port.HostPort)
		listener.SetInstancePort(port.ContainerPort)

		if port.CertificateName != "" {
			certificateARN, err := l.lookupCertificateARN(port.CertificateName)
			if err != nil {
				return nil, err
			}

			listener.SetSSLCertificateId(certificateARN)
		}

		// terminate ssl/https on load balancer
		switch strings.ToLower(port.Protocol) {
		case "http", "https":
			listener.SetInstanceProtocol("http")
		case "tcp", "ssl":
			listener.SetInstanceProtocol("tcp")
		default:
			return nil, fmt.Errorf("Unrecognized procotol '%s'", port.Protocol)
		}

		listeners[i] = listener
	}

	return listeners, nil
}

func (l *LoadBalancerProvider) lookupCertificateARN(certificateName string) (string, error) {
	output, err := l.AWS.IAM.ListServerCertificates(&iam.ListServerCertificatesInput{})
	if err != nil {
		return "", err
	}

	for _, meta := range output.ServerCertificateMetadataList {
		if aws.StringValue(meta.ServerCertificateName) == certificateName {
			return aws.StringValue(meta.Arn), nil
		}
	}

	return "", fmt.Errorf("Certificate with name '%s' does not exist", certificateName)
}

func (l *LoadBalancerProvider) createTags(loadBalancerID, loadBalancerName, environmentID string) error {
	tags := []models.Tag{
		{
			EntityID:   loadBalancerID,
			EntityType: "load_balancer",
			Key:        "name",
			Value:      loadBalancerName,
		},
		{
			EntityID:   loadBalancerID,
			EntityType: "load_balancer",
			Key:        "environment_id",
			Value:      environmentID,
		},
	}

	for _, tag := range tags {
		if err := l.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}

func renderLoadBalancerRolePolicy(region, accountID, loadBalancerID, rolePolicyTemplate string) (string, error) {
	tmpl, err := template.New("").Parse(rolePolicyTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse role policy template: %v", err)
	}

	context := struct {
		Region         string
		AccountID      string
		LoadBalancerID string
	}{
		Region:         region,
		AccountID:      accountID,
		LoadBalancerID: loadBalancerID,
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, context); err != nil {
		return "", fmt.Errorf("Failed to render role policy: %v", err)
	}

	return rendered.String(), nil
}
