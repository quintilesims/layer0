package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/quintilesims/layer0/common/models"
)

// Read returns a *models.LoadBalancer based on the provided loadBalancerID. The loadBalancerID
// is used when the DescribeLoadBalancers request is made to AWS.
func (l *LoadBalancerProvider) Read(loadBalancerID string) (*models.LoadBalancer, error) {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)

	model, err := l.makeLoadBalancerModel(loadBalancerID)
	if err != nil {
		return nil, err
	}

	loadBalancer, err := describeLoadBalancer(l.AWS.ELB, l.AWS.ALB, fqLoadBalancerID)
	if err != nil {
		return nil, err
	}

	if loadBalancer.isCLB {
		model.Ports = make([]models.Port, len(loadBalancer.CLB.ListenerDescriptions))
		for i, description := range loadBalancer.CLB.ListenerDescriptions {
			port := models.Port{
				ContainerPort: aws.Int64Value(description.Listener.InstancePort),
				HostPort:      aws.Int64Value(description.Listener.LoadBalancerPort),
				Protocol:      aws.StringValue(description.Listener.Protocol),
			}

			if certificateARN := aws.StringValue(description.Listener.SSLCertificateId); certificateARN != "" {
				port.CertificateARN = certificateARN
			}

			model.Ports[i] = port
		}

		model.HealthCheck = models.HealthCheck{
			Target:             aws.StringValue(loadBalancer.CLB.HealthCheck.Target),
			Interval:           int(aws.Int64Value(loadBalancer.CLB.HealthCheck.Interval)),
			Timeout:            int(aws.Int64Value(loadBalancer.CLB.HealthCheck.Timeout)),
			HealthyThreshold:   int(aws.Int64Value(loadBalancer.CLB.HealthCheck.HealthyThreshold)),
			UnhealthyThreshold: int(aws.Int64Value(loadBalancer.CLB.HealthCheck.UnhealthyThreshold)),
		}
	}

	if loadBalancer.isALB {
		securityGroupName := getLoadBalancerSGName(fqLoadBalancerID)
		securityGroup, err := readSG(l.AWS.EC2, securityGroupName)
		if err != nil {
			return nil, err
		}

		model.Ports = make([]models.Port, len(securityGroup.IpPermissions))
		for i, p := range securityGroup.IpPermissions {
			port := models.Port{
				// container port isn't used for ALBs
				ContainerPort: aws.Int64Value(p.FromPort),
				HostPort:      aws.Int64Value(p.FromPort),
				Protocol:      aws.StringValue(p.IpProtocol),
			}

			model.Ports[i] = port
		}

		targetGroupID := fqLoadBalancerID
		targetGroup, err := l.readTargetGroup(targetGroupID)

		if err != nil {
			return nil, err
		}

		model.HealthCheck = models.HealthCheck{
			Path:               aws.StringValue(targetGroup.HealthCheckPath),
			Interval:           int(aws.Int64Value(targetGroup.HealthCheckIntervalSeconds)),
			Timeout:            int(aws.Int64Value(targetGroup.HealthCheckTimeoutSeconds)),
			HealthyThreshold:   int(aws.Int64Value(targetGroup.HealthyThresholdCount)),
			UnhealthyThreshold: int(aws.Int64Value(targetGroup.UnhealthyThresholdCount)),
		}

		cert, err := l.readCertificate(loadBalancer.ALB.LoadBalancerArn)
		if err != nil {
			return nil, err
		}

		if cert != nil {
			for _, p := range model.Ports {
				if aws.Int64Value(targetGroup.Port) == p.HostPort {
					p.CertificateARN = aws.StringValue(cert.CertificateArn)
					break
				}
			}
		}
	}

	model.IsPublic = aws.StringValue(loadBalancer.GetScheme()) == "internet-facing"
	model.URL = aws.StringValue(loadBalancer.GetDNSName())

	return model, nil
}

func (l *LoadBalancerProvider) readTargetGroup(targetGroupID string) (*alb.TargetGroup, error) {
	input := &alb.DescribeTargetGroupsInput{}
	input.SetNames([]*string{aws.String(targetGroupID)})

	output, err := l.AWS.ALB.DescribeTargetGroups(input)
	if err != nil {
		return nil, err
	}

	return output.TargetGroups[0], nil
}

func (l *LoadBalancerProvider) readCertificate(loadBalancerArn *string) (*alb.Certificate, error) {
	input := &alb.DescribeListenersInput{}
	input.LoadBalancerArn = loadBalancerArn

	output, err := l.AWS.ALB.DescribeListeners(input)
	if err != nil {
		return nil, err
	}

	for _, l := range output.Listeners {
		for _, c := range l.Certificates {
			return c, nil
		}
	}

	return nil, nil
}

func (l *LoadBalancerProvider) makeLoadBalancerModel(loadBalancerID string) (*models.LoadBalancer, error) {
	model := &models.LoadBalancer{
		LoadBalancerID: loadBalancerID,
	}

	tags, err := l.TagStore.SelectByTypeAndID("load_balancer", loadBalancerID)
	if err != nil {
		return nil, err
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.LoadBalancerName = tag.Value
	}

	if tag, ok := tags.WithKey("type").First(); ok {
		model.LoadBalancerType = tag.Value
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		model.EnvironmentID = tag.Value

		environmentTags, err := l.TagStore.SelectByTypeAndID("environment", tag.Value)
		if err != nil {
			return nil, err
		}

		if tag, ok := environmentTags.WithKey("name").First(); ok {
			model.EnvironmentName = tag.Value
		}
	}

	allServiceTags, err := l.TagStore.SelectByType("service")
	if err != nil {
		return nil, err
	}

	if tag, ok := allServiceTags.WithKey("load_balancer_id").WithValue(loadBalancerID).First(); ok {
		model.ServiceID = tag.EntityID

		if t, ok := allServiceTags.WithID(tag.EntityID).WithKey("name").First(); ok {
			model.ServiceName = t.Value
		}
	}

	return model, nil
}
