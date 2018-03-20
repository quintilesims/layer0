package aws

import (
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
)

type genericLoadBalancer struct {
	ELB   *elb.LoadBalancerDescription
	ALB   *alb.LoadBalancer
	isALB bool
	isELB bool
}

func (c genericLoadBalancer) GetScheme() *string {
	if c.isELB {
		return c.ELB.Scheme
	} else if c.isALB {
		return c.ALB.Scheme
	}

	return nil
}

func (c genericLoadBalancer) GetDNSName() *string {
	if c.isELB {
		return c.ELB.DNSName
	} else if c.isALB {
		return c.ALB.DNSName
	}

	return nil
}
