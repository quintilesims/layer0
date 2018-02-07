package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/quintilesims/layer0/common/models"
)

// Read returns a *models.LoadBalancer based on the provided loadBalancerID. The loadBalancerID
// is used when the DescribeLoadBalancers request is made to AWS.
func (l *LoadBalancerProvider) Read(loadBalancerID string) (*models.LoadBalancer, error) {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)
	loadBalancer, err := describeLoadBalancer(l.AWS.ELB, fqLoadBalancerID)
	if err != nil {
		return nil, err
	}

	model, err := l.makeLoadBalancerModel(loadBalancerID)
	if err != nil {
		return nil, err
	}

	model.Ports = make([]models.Port, len(loadBalancer.ListenerDescriptions))
	for i, description := range loadBalancer.ListenerDescriptions {
		port := models.Port{
			ContainerPort: aws.Int64Value(description.Listener.InstancePort),
			HostPort:      aws.Int64Value(description.Listener.LoadBalancerPort),
			Protocol:      aws.StringValue(description.Listener.Protocol),
		}

		if certificateARN := aws.StringValue(description.Listener.SSLCertificateId); certificateARN != "" {
			// certificate arn format:  arn:aws:iam:region:012345678910:certificate/path/to/name
			port.CertificateARN = certificateARN
		}

		model.Ports[i] = port
	}

	model.HealthCheck = models.HealthCheck{
		Target:             aws.StringValue(loadBalancer.HealthCheck.Target),
		Interval:           int(aws.Int64Value(loadBalancer.HealthCheck.Interval)),
		Timeout:            int(aws.Int64Value(loadBalancer.HealthCheck.Timeout)),
		HealthyThreshold:   int(aws.Int64Value(loadBalancer.HealthCheck.HealthyThreshold)),
		UnhealthyThreshold: int(aws.Int64Value(loadBalancer.HealthCheck.UnhealthyThreshold)),
	}

	model.IsPublic = aws.StringValue(loadBalancer.Scheme) == "internet-facing"
	model.URL = aws.StringValue(loadBalancer.DNSName)

	return model, nil
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
