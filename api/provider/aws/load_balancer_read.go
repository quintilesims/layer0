package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

// Read returns a *models.LoadBalancer based on the provided loadBalancerID. The loadBalancerID
// is used when the DescribeLoadBalancers request is made to AWS.
func (l *LoadBalancerProvider) Read(loadBalancerID string) (*models.LoadBalancer, error) {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)
	loadBalancer, err := describeLoadBalancer(l.AWS.ELB, l.AWS.ALB, fqLoadBalancerID)
	if err != nil {
		return nil, err
	}

	model, err := l.makeLoadBalancerModel(loadBalancerID)
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
			Path:               config.DefaultLoadBalancerHealthCheck().Path,
			Interval:           int(aws.Int64Value(loadBalancer.CLB.HealthCheck.Interval)),
			Timeout:            int(aws.Int64Value(loadBalancer.CLB.HealthCheck.Timeout)),
			HealthyThreshold:   int(aws.Int64Value(loadBalancer.CLB.HealthCheck.HealthyThreshold)),
			UnhealthyThreshold: int(aws.Int64Value(loadBalancer.CLB.HealthCheck.UnhealthyThreshold)),
		}
	}

	if loadBalancer.isALB {
		listeners, err := l.readListeners(loadBalancer.ALB.LoadBalancerArn)
		if err != nil {
			return nil, err
		}

		model.Ports = make([]models.Port, len(listeners))
		for i, listener := range listeners {
			port := models.Port{
				ContainerPort: aws.Int64Value(listener.Port),
				HostPort:      aws.Int64Value(listener.Port),
				Protocol:      aws.StringValue(listener.Protocol),
			}

			if len(listener.Certificates) > 0 {
				port.CertificateARN = aws.StringValue(listener.Certificates[0].CertificateArn)
			}

			model.Ports[i] = port
		}

		targetGroupName := fqLoadBalancerID
		targetGroup, err := readTargetGroup(l.AWS.ALB, aws.String(targetGroupName), nil)
		if err != nil {
			return nil, err
		}

		model.HealthCheck = models.HealthCheck{
			Target:             config.DefaultLoadBalancerHealthCheck().Target,
			Path:               aws.StringValue(targetGroup.HealthCheckPath),
			Interval:           int(aws.Int64Value(targetGroup.HealthCheckIntervalSeconds)),
			Timeout:            int(aws.Int64Value(targetGroup.HealthCheckTimeoutSeconds)),
			HealthyThreshold:   int(aws.Int64Value(targetGroup.HealthyThresholdCount)),
			UnhealthyThreshold: int(aws.Int64Value(targetGroup.UnhealthyThresholdCount)),
		}
	}

	model.IsPublic = loadBalancer.Scheme() == "internet-facing"
	model.URL = loadBalancer.DNSName()

	return model, nil
}

func (l *LoadBalancerProvider) readListeners(loadBalancerArn *string) ([]*alb.Listener, error) {
	input := &alb.DescribeListenersInput{}
	input.LoadBalancerArn = loadBalancerArn

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := l.AWS.ALB.DescribeListeners(input)
	if err != nil {
		return nil, err
	}

	return output.Listeners, nil
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
		model.LoadBalancerType = models.LoadBalancerType(tag.Value)
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
