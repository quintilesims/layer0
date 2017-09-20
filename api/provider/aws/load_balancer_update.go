package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/common/models"
)

func (l *LoadBalancerProvider) Update(req models.UpdateLoadBalancerRequest) error {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), req.LoadBalancerID)
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

		loadBalancerDescription, err := l.describeLoadBalancer(loadBalancerName)
		if err != nil {
			return err
		}

		portNumbers := make([]int64, len(loadBalancerDescription.ListenerDescriptions))
		for i, listenerDescription := range loadBalancerDescription.ListenerDescriptions {
			// todo: unsure if it's loadbalancer port we need
			portNumbers[i] = aws.Int64Value(listenerDescription.Listener.LoadBalancerPort)
		}

		if err := l.removeListeners(loadBalancerName, portNumbers); err != nil {
			return err
		}

		// todo: remove ingress from sg

		if err := l.addListeners(loadBalancerName, listeners); err != nil {
			return err
		}

		// todo: authorize ingress to sg

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
