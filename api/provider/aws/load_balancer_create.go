package aws

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

// Create is used to create an Classic or an Application Load Balancer using the
// specified Create Load Balancer Request. The Create Load Balancer Request contains
// the name of the Load Balancer, the Environment ID in which to create the Load Balancer,
// a flag to determine if the Load Balancer will be Internet-facing or internal,
// a list of ports to configure as the listeners, and a Health Check to determine
// if attached EC2 instances are in service or not. An IAM Role is created and
// an inline policy is attached to the Role that allows ECS to interact with the
// created Load Balancer. An EC2 Security Group is created and ingress rules are
// added based on the list of ports in the Create Load Balancer Request. The
// Security Group is then attached to the created Load Balancer.
func (l *LoadBalancerProvider) Create(req models.CreateLoadBalancerRequest) (string, error) {
	loadBalancerID := entityIDGenerator(req.LoadBalancerName)
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)
	fqEnvironmentID := addLayer0Prefix(l.Config.Instance(), req.EnvironmentID)

	environmentSGName := getEnvironmentSGName(fqEnvironmentID)
	environmentSG, err := readSG(l.AWS.EC2, environmentSGName)
	if err != nil {
		return "", err
	}

	scheme := "internal"
	subnets := l.Config.PrivateSubnets()
	securityGroupIDs := []string{aws.StringValue(environmentSG.GroupId)}

	if req.IsPublic {
		scheme = "internet-facing"
		subnets = l.Config.PublicSubnets()
	}

	loadBalancerSGName := getLoadBalancerSGName(fqLoadBalancerID)
	if err := createSG(
		l.AWS.EC2,
		loadBalancerSGName,
		fmt.Sprintf("SG for Layer0 load balancer %s", loadBalancerID),
		l.Config.VPC()); err != nil {
		return "", err
	}

	loadBalancerSG, err := readSG(l.AWS.EC2, loadBalancerSGName)
	if err != nil {
		return "", err
	}

	loadBalancerSGID := aws.StringValue(loadBalancerSG.GroupId)
	if len(req.Ports) == 0 {
		req.Ports = []models.Port{config.DefaultLoadBalancerPort()}
	}

	for _, port := range req.Ports {
		if err := l.authorizeSGIngressFromPort(loadBalancerSGID, int64(port.HostPort)); err != nil {
			return "", err
		}
	}

	securityGroupIDs = append(securityGroupIDs, aws.StringValue(loadBalancerSG.GroupId))

	roleName := getLoadBalancerRoleName(fqLoadBalancerID)
	if _, err := l.createRole(roleName, DefaultAssumeRolePolicy); err != nil {
		return "", err
	}

	policy, err := RenderLoadBalancerRolePolicy(
		l.Config.Region(),
		l.Config.AccountID(),
		fqLoadBalancerID,
		DefaultLBRolePolicyTemplate)
	if err != nil {
		return "", err
	}

	policyName := roleName
	if err := l.putRolePolicy(policyName, roleName, policy); err != nil {
		return "", err
	}

	if req.LoadBalancerType.Equals(models.ClassicLoadBalancerType) {
		listeners, err := l.portsToListeners(req.Ports)
		if err != nil {
			return "", err
		}

		if err := l.createLoadBalancer(
			fqLoadBalancerID,
			scheme,
			securityGroupIDs,
			subnets,
			listeners); err != nil {
			return "", err
		}

		if req.HealthCheck == (models.HealthCheck{}) {
			req.HealthCheck = config.DefaultLoadBalancerHealthCheck()
		}

		if err := l.updateELBHealthCheck(fqLoadBalancerID, req.HealthCheck); err != nil {
			return "", err
		}
	}

	if req.LoadBalancerType.Equals(models.ApplicationLoadBalancerType) {
		lb, err := l.createApplicationLoadBalancer(
			fqLoadBalancerID,
			scheme,
			securityGroupIDs,
			subnets)

		if err != nil {
			return "", err
		}

		if req.HealthCheck == (models.HealthCheck{}) {
			req.HealthCheck = config.DefaultLoadBalancerHealthCheck()
		}

		targetGroupName := fqLoadBalancerID
		tg, err := l.createTargetGroup(targetGroupName, req.HealthCheck)
		if err != nil {
			return "", err
		}

		if err := l.createListener(lb.LoadBalancerArn, tg.TargetGroupArn, req.Ports); err != nil {
			return "", err
		}
	}

	if err := l.createTags(loadBalancerID, req.LoadBalancerName, string(req.LoadBalancerType), req.EnvironmentID); err != nil {
		return "", err
	}

	return loadBalancerID, nil
}

func (l *LoadBalancerProvider) createTargetGroup(groupName string, healthCheck models.HealthCheck) (*alb.TargetGroup, error) {
	input := &alb.CreateTargetGroupInput{}
	input.SetName(groupName)
	input.SetPort(config.DefaultTargetGroupPort)
	input.SetProtocol(config.DefaultTargetGroupProtocol)
	input.SetVpcId(l.Config.VPC())
	input.SetTargetType(alb.TargetTypeEnumIp)

	// set health check
	input.SetHealthCheckPath(healthCheck.Path)
	input.SetHealthCheckIntervalSeconds(int64(healthCheck.Interval))
	input.SetHealthCheckTimeoutSeconds(int64(healthCheck.Timeout))
	input.SetHealthyThresholdCount(int64(healthCheck.HealthyThreshold))
	input.SetUnhealthyThresholdCount(int64(healthCheck.UnhealthyThreshold))

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := l.AWS.ALB.CreateTargetGroup(input)
	if err != nil {
		return nil, err
	}

	return output.TargetGroups[0], nil
}

func (l *LoadBalancerProvider) createListener(loadBalancerArn, targetGroupArn *string, ports []models.Port) error {
	for _, port := range ports {
		input := &alb.CreateListenerInput{}
		input.SetPort(port.HostPort)
		input.SetProtocol(port.Protocol)
		input.SetDefaultActions([]*alb.Action{
			{
				TargetGroupArn: targetGroupArn,
				Type:           aws.String(alb.ActionTypeEnumForward),
			},
		})
		input.LoadBalancerArn = loadBalancerArn

		if port.CertificateARN != "" {
			input.SetCertificates([]*alb.Certificate{
				{
					CertificateArn: aws.String(port.CertificateARN),
				},
			})
		}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := l.AWS.ALB.CreateListener(input); err != nil {
			return err
		}
	}

	return nil
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

func (l *LoadBalancerProvider) authorizeSGIngressFromPort(groupID string, hostPort int64) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{}
	input.SetGroupId(groupID)
	input.SetCidrIp("0.0.0.0/0")
	input.SetIpProtocol("TCP")
	input.SetFromPort(hostPort)
	input.SetToPort(hostPort)

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

func (l *LoadBalancerProvider) createApplicationLoadBalancer(
	loadBalancerName string,
	scheme string,
	securityGroupIDs []string,
	subnetIDs []string,
) (*alb.LoadBalancer, error) {
	securityGroups := make([]*string, len(securityGroupIDs))
	for i, securityGroupID := range securityGroupIDs {
		securityGroups[i] = aws.String(securityGroupID)
	}

	subnets := make([]*string, len(subnetIDs))
	for i, subnetID := range subnetIDs {
		subnets[i] = aws.String(subnetID)
	}

	input := &alb.CreateLoadBalancerInput{}
	input.SetName(loadBalancerName)
	input.SetScheme(scheme)
	input.SetSecurityGroups(securityGroups)
	input.SetSubnets(subnets)
	input.SetType(alb.LoadBalancerTypeEnumApplication)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	createLBOutput, err := l.AWS.ALB.CreateLoadBalancer(input)
	if err != nil {
		return nil, err
	}

	waitInput := &alb.DescribeLoadBalancersInput{}
	waitInput.SetLoadBalancerArns([]*string{
		createLBOutput.LoadBalancers[0].LoadBalancerArn,
	})

	if err := l.AWS.ALB.WaitUntilLoadBalancerExists(waitInput); err != nil {
		return nil, err
	}

	return createLBOutput.LoadBalancers[0], nil
}

func (l *LoadBalancerProvider) portsToListeners(ports []models.Port) ([]*elb.Listener, error) {
	listeners := make([]*elb.Listener, len(ports))
	for i, port := range ports {
		listener := &elb.Listener{}
		listener.SetProtocol(port.Protocol)
		listener.SetLoadBalancerPort(port.HostPort)
		listener.SetInstancePort(port.ContainerPort)

		if port.CertificateARN != "" {
			listener.SetSSLCertificateId(port.CertificateARN)
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

func (l *LoadBalancerProvider) createTags(loadBalancerID, loadBalancerName, loadBalancerType, environmentID string) error {
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
		{
			EntityID:   loadBalancerID,
			EntityType: "load_balancer",
			Key:        "type",
			Value:      loadBalancerType,
		},
	}

	for _, tag := range tags {
		if err := l.TagStore.Insert(tag); err != nil {
			return err
		}
	}

	return nil
}

func RenderLoadBalancerRolePolicy(region, accountID, loadBalancerID, rolePolicyTemplate string) (string, error) {
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
