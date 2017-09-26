package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

// todo: catch 'EntityDoesNotExist' errors
func (l *LoadBalancerProvider) Read(loadBalancerID string) (*models.LoadBalancer, error) {
	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), loadBalancerID)
	loadBalancer, err := l.describeLoadBalancer(fqLoadBalancerID)
	if err != nil {
		return nil, err
	}

	model, err := l.newModel(loadBalancerID)
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
			split := strings.SplitN(certificateARN, "/", -1)
			certificateName := split[len(split)-1]
			port.CertificateName = certificateName
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

func (l *LoadBalancerProvider) describeLoadBalancer(loadBalancerName string) (*elb.LoadBalancerDescription, error) {
	input := &elb.DescribeLoadBalancersInput{}
	input.SetLoadBalancerNames([]*string{aws.String(loadBalancerName)})
	input.SetPageSize(1)

	output, err := l.AWS.ELB.DescribeLoadBalancers(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "LoadBalancerNotFound" {
			return nil, errors.Newf(errors.LoadBalancerDoesNotExist, "LoadBalancer '%s' does not exist", loadBalancerName)
		}

		return nil, err
	}

	if len(output.LoadBalancerDescriptions) != 1 {
		return nil, errors.Newf(errors.LoadBalancerDoesNotExist, "LoadBalancer '%s' does not exist", loadBalancerName)
	}

	return output.LoadBalancerDescriptions[0], nil
}

func (l *LoadBalancerProvider) newModel(loadBalancerID string) (*models.LoadBalancer, error) {
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
