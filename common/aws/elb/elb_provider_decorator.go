// Generated by go-decorator, DO NOT EDIT
package elb

import ()

type ProviderDecorator struct {
	Inner     Provider
	Decorator func(name string, call func() error) error
}

func (this *ProviderDecorator) CreateLoadBalancer(p0 string, p1 string, p2 []*string, p3 []*string, p4 []*Listener) (v0 *string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.CreateLoadBalancer(p0, p1, p2, p3, p4)
		return err
	}
	err = this.Decorator("CreateLoadBalancer", call)
	return v0, err
}
func (this *ProviderDecorator) ConfigureHealthCheck(p0 string, p1 *HealthCheck) (err error) {
	call := func() error {
		var err error
		err = this.Inner.ConfigureHealthCheck(p0, p1)
		return err
	}
	err = this.Decorator("ConfigureHealthCheck", call)
	return err
}
func (this *ProviderDecorator) DescribeLoadBalancer(p0 string) (v0 *LoadBalancerDescription, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeLoadBalancer(p0)
		return err
	}
	err = this.Decorator("DescribeLoadBalancer", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeLoadBalancers() (v0 []*LoadBalancerDescription, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeLoadBalancers()
		return err
	}
	err = this.Decorator("DescribeLoadBalancers", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeInstanceHealth(p0 string) (v0 []*InstanceState, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeInstanceHealth(p0)
		return err
	}
	err = this.Decorator("DescribeInstanceHealth", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeLoadBalancerAttributes(p0 string) (v0 *LoadBalancerAttributes, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeLoadBalancerAttributes(p0)
		return err
	}
	err = this.Decorator("DescribeLoadBalancerAttributes", call)
	return v0, err
}
func (this *ProviderDecorator) DeleteLoadBalancer(p0 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeleteLoadBalancer(p0)
		return err
	}
	err = this.Decorator("DeleteLoadBalancer", call)
	return err
}
func (this *ProviderDecorator) RegisterInstancesWithLoadBalancer(p0 string, p1 []string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.RegisterInstancesWithLoadBalancer(p0, p1)
		return err
	}
	err = this.Decorator("RegisterInstancesWithLoadBalancer", call)
	return err
}
func (this *ProviderDecorator) DeregisterInstancesFromLoadBalancer(p0 string, p1 []string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeregisterInstancesFromLoadBalancer(p0, p1)
		return err
	}
	err = this.Decorator("DeregisterInstancesFromLoadBalancer", call)
	return err
}
func (this *ProviderDecorator) CreateLoadBalancerListeners(p0 string, p1 []*Listener) (err error) {
	call := func() error {
		var err error
		err = this.Inner.CreateLoadBalancerListeners(p0, p1)
		return err
	}
	err = this.Decorator("CreateLoadBalancerListeners", call)
	return err
}
func (this *ProviderDecorator) DeleteLoadBalancerListeners(p0 string, p1 []*Listener) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeleteLoadBalancerListeners(p0, p1)
		return err
	}
	err = this.Decorator("DeleteLoadBalancerListeners", call)
	return err
}
func (this *ProviderDecorator) SetIdleTimeout(p0 string, p1 int) (err error) {
	call := func() error {
		var err error
		err = this.Inner.SetIdleTimeout(p0, p1)
		return err
	}
	err = this.Decorator("SetIdleTimeout", call)
	return err
}

