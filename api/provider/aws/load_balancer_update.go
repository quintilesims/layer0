package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/quintilesims/layer0/common/models"
)

// Update is used to update Classic and Application Load Balancers using the
// specified Update Load Balancer Request. The Update Load Balancer Request
// contains the Load Balancer ID, a list of ports to configure as the listeners,
// and a Health Check to determine the state of the attached EC2 instances in the
// case of Classic ELBs or registered Targets in the case of ALBs. If ports are
// included in the Update Load Balancer Request, all existing listeners and EC2
// Security Group ingress rules are removed first and then new listeners and
// Security Group ingress rules are created based on the provided list of ports.
func (l *LoadBalancerProvider) Update(loadBalancerID string, req models.UpdateLoadBalancerRequest) error {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)
	loadBalancerName := fqLoadBalancerID

	model, err := l.makeLoadBalancerModel(loadBalancerID)
	if err != nil {
		return err
	}

	isClassicELB := strings.EqualFold(model.LoadBalancerType, models.ClassicLoadBalancerType)
	isAppLB := strings.EqualFold(model.LoadBalancerType, models.ApplicationLoadBalancerType)

	if req.HealthCheck != nil {
		if isClassicELB {
			if err := l.updateELBHealthCheck(loadBalancerName, *req.HealthCheck); err != nil {
				return err
			}
		}

		if isAppLB {
			if err := l.updateALBHealthCheck(loadBalancerName, *req.HealthCheck); err != nil {
				return err
			}
		}
	}

	if req.Ports != nil {
		loadBalancer, err := describeLoadBalancer(l.AWS.ELB, l.AWS.ALB, loadBalancerName)
		if err != nil {
			return err
		}

		if isClassicELB {
			if err := l.updateELBListeners(*req.Ports, loadBalancer.CLB.ListenerDescriptions, loadBalancerName); err != nil {
				return err
			}
		}

		if isAppLB {
			targetGroupName := loadBalancerName
			if err := l.updateALBListeners(*req.Ports, targetGroupName, loadBalancer.ALB.LoadBalancerArn); err != nil {
				return err
			}
		}

		// update ingress and egress rules of the loadbalancer security group
		securityGroupName := getLoadBalancerSGName(fqLoadBalancerID)
		securityGroup, err := readSG(l.AWS.EC2, securityGroupName)
		if err != nil {
			return err
		}

		containsPort := func(port int64, ports []models.Port) bool {
			for _, p := range ports {
				if port == p.HostPort {
					return true
				}
			}

			return false
		}

		// remove permissions for ports not in the request
		for _, p := range securityGroup.IpPermissions {
			if !containsPort(aws.Int64Value(p.FromPort), *req.Ports) {
				l.revokeSGIngressFromPort(aws.StringValue(securityGroup.GroupId), aws.Int64Value(p.FromPort))
			}
		}

		permissionsContainsPort := func(port int64, permissions []*ec2.IpPermission) bool {
			for _, p := range permissions {
				if port == aws.Int64Value(p.FromPort) {
					return true
				}
			}

			return false
		}

		// add permission for request ports that don't exist in the security group
		for _, p := range *req.Ports {
			if !permissionsContainsPort(p.HostPort, securityGroup.IpPermissions) {
				l.authorizeSGIngressFromPort(aws.StringValue(securityGroup.GroupId), p.HostPort)
			}
		}
	}

	return nil
}

func (l *LoadBalancerProvider) updateELBListeners(ports []models.Port, listenerDescriptions []*elb.ListenerDescription, loadBalancerName string) error {
	// remove listener not in ports
	var listenersToRemove []*int64
	for _, ld := range listenerDescriptions {
		removeListener := true
		for _, p := range ports {
			if p.HostPort == aws.Int64Value(ld.Listener.LoadBalancerPort) {
				removeListener = false
				break
			}
		}

		if removeListener {
			listenersToRemove = append(listenersToRemove, ld.Listener.LoadBalancerPort)
		}
	}

	if len(listenersToRemove) > 0 {
		input := &elb.DeleteLoadBalancerListenersInput{}
		input.SetLoadBalancerName(loadBalancerName)
		input.SetLoadBalancerPorts(listenersToRemove)

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := l.AWS.ELB.DeleteLoadBalancerListeners(input); err != nil {
			return err
		}
	}

	// add listener which doesn't exist in ports
	var listenersToAdd []*elb.Listener
	for _, p := range ports {
		addListener := true
		for _, ld := range listenerDescriptions {
			if p.HostPort == aws.Int64Value(ld.Listener.LoadBalancerPort) {
				addListener = false
				break
			}
		}

		if addListener {
			newListener, err := l.portsToListeners([]models.Port{p})
			if err != nil {
				return err
			}

			listenersToAdd = append(listenersToAdd, newListener[0])
		}
	}

	if len(listenersToAdd) > 0 {
		input := &elb.CreateLoadBalancerListenersInput{}
		input.SetLoadBalancerName(loadBalancerName)
		input.SetListeners(listenersToAdd)

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := l.AWS.ELB.CreateLoadBalancerListeners(input); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (l *LoadBalancerProvider) updateALBListeners(ports []models.Port, targetGroupName string, loadBalancerArn *string) error {
	targetGroup, err := readTargetGroup(l.AWS.ALB, targetGroupName)
	if err != nil {
		return err
	}

	targetGroupArn := targetGroup.TargetGroupArn
	var listeners []alb.Listener

	descListenersInput := &alb.DescribeListenersInput{}
	descListenersInput.LoadBalancerArn = loadBalancerArn
	descListenersInput.SetPageSize(10)
	fnPage := func(output *alb.DescribeListenersOutput, lastPage bool) bool {
		for _, l := range output.Listeners {
			listeners = append(listeners, *l)
		}

		return !lastPage
	}

	if err := l.AWS.ALB.DescribeListenersPages(descListenersInput, fnPage); err != nil {
		return err
	}

	// remove listeners
	for _, listener := range listeners {
		removeListener := true
		for _, p := range ports {
			if aws.Int64Value(listener.Port) == p.HostPort {
				removeListener = false
				break
			}
		}

		if removeListener {
			removeListenerInput := &alb.DeleteListenerInput{}
			removeListenerInput.ListenerArn = listener.ListenerArn

			if err := removeListenerInput.Validate(); err != nil {
				return err
			}

			if _, err := l.AWS.ALB.DeleteListener(removeListenerInput); err != nil {
				return err
			}
		}
	}

	// add listeners
	for _, p := range ports {
		addListener := true
		for _, listener := range listeners {
			if aws.Int64Value(listener.Port) == p.HostPort {
				addListener = false
				break
			}
		}

		if addListener {
			createListenerInput := &alb.CreateListenerInput{}
			createListenerInput.SetPort(p.HostPort)
			createListenerInput.SetProtocol(p.Protocol)
			createListenerInput.LoadBalancerArn = loadBalancerArn
			createListenerInput.SetDefaultActions([]*alb.Action{
				{
					TargetGroupArn: targetGroupArn,
					Type:           aws.String(alb.ActionTypeEnumForward),
				},
			})

			if p.CertificateARN != "" {
				createListenerInput.SetCertificates([]*alb.Certificate{
					{
						CertificateArn: aws.String(p.CertificateARN),
					},
				})
			}

			if err := createListenerInput.Validate(); err != nil {
				return err
			}

			if _, err := l.AWS.ALB.CreateListener(createListenerInput); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *LoadBalancerProvider) updateELBHealthCheck(loadBalancerName string, healthCheck models.HealthCheck) error {
	hc := &elb.HealthCheck{
		Target:             aws.String(healthCheck.Target),
		Interval:           aws.Int64(int64(healthCheck.Interval)),
		Timeout:            aws.Int64(int64(healthCheck.Timeout)),
		HealthyThreshold:   aws.Int64(int64(healthCheck.HealthyThreshold)),
		UnhealthyThreshold: aws.Int64(int64(healthCheck.UnhealthyThreshold)),
	}

	input := &elb.ConfigureHealthCheckInput{}
	input.SetLoadBalancerName(loadBalancerName)
	input.SetHealthCheck(hc)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.ELB.ConfigureHealthCheck(input); err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProvider) updateALBHealthCheck(loadBalancerName string, healthCheck models.HealthCheck) error {
	targetGroupName := loadBalancerName
	targetGroup, err := readTargetGroup(l.AWS.ALB, targetGroupName)
	if err != nil {
		return err
	}

	input := &alb.ModifyTargetGroupInput{}
	input.SetHealthCheckIntervalSeconds(int64(healthCheck.Interval))
	input.SetHealthCheckPath(healthCheck.Path)
	input.SetHealthCheckTimeoutSeconds(int64(healthCheck.Timeout))
	input.SetHealthyThresholdCount(int64(healthCheck.HealthyThreshold))
	input.SetUnhealthyThresholdCount(int64(healthCheck.UnhealthyThreshold))
	input.SetTargetGroupArn(aws.StringValue(targetGroup.TargetGroupArn))

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := l.AWS.ALB.ModifyTargetGroup(input); err != nil {
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
