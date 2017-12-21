package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/common/models"
)

// Update is used to update an Elastic Load Balancer using the specified Update
// Load Balancer Request. The Update Load Balancer Request contains the Load
// Balancer ID, a list of ports to configure as the listeners, and a Health
// Check to determine if attached EC2 instances are in service or not. If ports
// are included in the Update Load Balancer Request, all existing listeners and
// EC2 Security Group ingress rules are removed first and then new listeners and
// Security Group ingress rules are created based on the provided list of ports.
func (l *LoadBalancerProvider) Update(loadBalancerID string, req models.UpdateLoadBalancerRequest) error {
	fqLoadBalancerID := addLayer0Prefix(l.Context, loadBalancerID)
	loadBalancerName := fqLoadBalancerID

	if req.HealthCheck != nil {
		healthCheck := &elb.HealthCheck{
			Target:             aws.String(req.HealthCheck.Target),
			Interval:           aws.Int64(int64(req.HealthCheck.Interval)),
			Timeout:            aws.Int64(int64(req.HealthCheck.Timeout)),
			HealthyThreshold:   aws.Int64(int64(req.HealthCheck.HealthyThreshold)),
			UnhealthyThreshold: aws.Int64(int64(req.HealthCheck.UnhealthyThreshold)),
		}

		if err := l.updateHealthCheck(loadBalancerName, healthCheck); err != nil {
			return err
		}
	}

	if req.Ports != nil {
		listeners, err := l.portsToListeners(*req.Ports)
		if err != nil {
			return err
		}

		securityGroupName := getLoadBalancerSGName(fqLoadBalancerID)
		securityGroup, err := readSG(l.AWS.EC2, securityGroupName)
		if err != nil {
			return err
		}

		securityGroupID := aws.StringValue(securityGroup.GroupId)

		loadBalancerDescription, err := describeLoadBalancer(l.AWS.ELB, loadBalancerName)
		if err != nil {
			return err
		}

		// remove all of the current listeners and security group ingress rules from the
		// load balancer and its security group
		portNumbers := make([]int64, len(loadBalancerDescription.ListenerDescriptions))
		for i, listenerDescription := range loadBalancerDescription.ListenerDescriptions {
			portNumber := aws.Int64Value(listenerDescription.Listener.LoadBalancerPort)
			portNumbers[i] = portNumber

			if err := l.revokeSGIngressFromPort(securityGroupID, portNumber); err != nil {
				return err
			}
		}

		if err := l.removeListeners(loadBalancerName, portNumbers); err != nil {
			return err
		}

		// add all of the new listeners and security group ingress rules to the
		// load balancer and its security group
		if err := l.addListeners(loadBalancerName, listeners); err != nil {
			return err
		}

		for _, listener := range listeners {
			loadBalancerListenerPort := aws.Int64Value(listener.LoadBalancerPort)

			if err := l.authorizeSGIngressFromPort(securityGroupID, loadBalancerListenerPort); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *LoadBalancerProvider) updateHealthCheck(loadBalancerName string, healthCheck *elb.HealthCheck) error {
	input := &elb.ConfigureHealthCheckInput{}
	input.SetLoadBalancerName(loadBalancerName)
	input.SetHealthCheck(healthCheck)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.ELB.ConfigureHealthCheck(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) removeListeners(loadBalancerName string, portNumbers []int64) error {
	input := &elb.DeleteLoadBalancerListenersInput{}
	input.SetLoadBalancerName(loadBalancerName)

	ports := make([]*int64, len(portNumbers))
	for i, p := range portNumbers {
		ports[i] = aws.Int64(p)
	}
	input.SetLoadBalancerPorts(ports)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.ELB.DeleteLoadBalancerListeners(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) addListeners(loadBalancerName string, listeners []*elb.Listener) error {
	input := &elb.CreateLoadBalancerListenersInput{}
	input.SetLoadBalancerName(loadBalancerName)
	input.SetListeners(listeners)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.ELB.CreateLoadBalancerListeners(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) revokeSGIngressFromPort(groupID string, port int64) error {
	input := &ec2.RevokeSecurityGroupIngressInput{}
	input.SetGroupId(groupID)
	input.SetCidrIp("0.0.0.0/0")
	input.SetIpProtocol("TCP")
	input.SetFromPort(port)
	input.SetToPort(port)

	if _, err := l.AWS.EC2.RevokeSecurityGroupIngress(input); err != nil {
		return err
	}

	return nil
}
