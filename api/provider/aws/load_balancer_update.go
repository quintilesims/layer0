package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/quintilesims/layer0/common/models"
)

func (l *LoadBalancerProvider) Update(req models.UpdateLoadBalancerRequest) (*models.LoadBalancer, error) {
	ports := req.Ports
	healthCheck := req.HealthCheck

	fqLoadBalancerID := addLayer0Prefix(l.Config.Instance(), req.LoadBalancerID)

	loadBalancer, err := l.describeLoadBalancer(fqLoadBalancerID)
	if err != nil {
		return nil, err
	}

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

		port = (*ports)[i]
	}

	model := &models.LoadBalancer{
		LoadBalancerID: req.LoadBalancerID,
		IsPublic:       aws.StringValue(loadBalancer.Scheme) == "internet-facing",
		URL:            aws.StringValue(loadBalancer.DNSName),
		Ports:          (*ports),
		HealthCheck:    healthCheck,
	}

	if err := l.populateModelTagss(req.LoadBalancerID, model); err != nil {
		return nil, err
	}

	return model, nil
}
