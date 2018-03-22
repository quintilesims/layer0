package aws

import (
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
)

type genericLoadBalancer struct {
	CLB   *elb.LoadBalancerDescription
	ALB   *alb.LoadBalancer
	isALB bool
	isCLB bool
}

func (c genericLoadBalancer) GetScheme() *string {
	if c.isCLB {
		return c.CLB.Scheme
	} else if c.isALB {
		return c.ALB.Scheme
	}

	return nil
}

func (c genericLoadBalancer) GetDNSName() *string {
	if c.isCLB {
		return c.CLB.DNSName
	} else if c.isALB {
		return c.ALB.DNSName
	}

	return nil
}
