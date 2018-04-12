package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	alb "github.com/aws/aws-sdk-go/service/elbv2"
)

func newGenericLoadBalancer(CLB *elb.LoadBalancerDescription, ALB *alb.LoadBalancer) *genericLoadBalancer {
	if CLB != nil {
		return &genericLoadBalancer{
			CLB:   CLB,
			isCLB: true,
		}
	}

	if ALB != nil {
		return &genericLoadBalancer{
			ALB:   ALB,
			isALB: true,
		}
	}

	return nil
}

type genericLoadBalancer struct {
	CLB   *elb.LoadBalancerDescription
	ALB   *alb.LoadBalancer
	isALB bool
	isCLB bool
}

func (c genericLoadBalancer) Scheme() string {
	if c.isCLB {
		return aws.StringValue(c.CLB.Scheme)
	} else if c.isALB {
		return aws.StringValue(c.ALB.Scheme)
	}

	return ""
}

func (c genericLoadBalancer) DNSName() string {
	if c.isCLB {
		return aws.StringValue(c.CLB.DNSName)
	} else if c.isALB {
		return aws.StringValue(c.ALB.DNSName)
	}

	return ""
}
