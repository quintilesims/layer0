package elb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/common/aws/provider"
)

type Provider interface {
	CreateLoadBalancer(loadBalancerName, scheme string, securityGroups, subnets []*string, listeners []*Listener) (*string, error)
	ConfigureHealthCheck(loadBalancerName string, check *HealthCheck) error
	DescribeLoadBalancer(loadBalancerName string) (*LoadBalancerDescription, error)
	DescribeLoadBalancers() ([]*LoadBalancerDescription, error)
	DescribeInstanceHealth(loadBalancerName string) ([]*InstanceState, error)
	DeleteLoadBalancer(loadBalancerName string) error
	RegisterInstancesWithLoadBalancer(loadBalancerName string, instanceIDs []string) error
	DeregisterInstancesFromLoadBalancer(loadBalancerName string, instanceIDs []string) error
	CreateLoadBalancerListeners(loadBalancerName string, listeners []*Listener) error
	DeleteLoadBalancerListeners(loadBalancerName string, listeners []*Listener) error
}

type Listener struct {
	*elb.Listener
}

type ListenerDescription struct {
	*elb.ListenerDescription
}

func NewListenerDescription(listener *Listener) *ListenerDescription {
	return &ListenerDescription{
		&elb.ListenerDescription{
			Listener: listener.Listener,
		},
	}
}

type HealthCheck struct {
	*elb.HealthCheck
}

type LoadBalancerDescription struct {
	*elb.LoadBalancerDescription
}

func NewLoadBalancerDescription(name, scheme string, listeners []*Listener) *LoadBalancerDescription {
	listenerDescriptions := []*elb.ListenerDescription{}
	for _, l := range listeners {
		listenerDescriptions = append(listenerDescriptions, NewListenerDescription(l).ListenerDescription)
	}

	return &LoadBalancerDescription{
		&elb.LoadBalancerDescription{
			LoadBalancerName:     aws.String(name),
			DNSName:              aws.String(name),
			Scheme:               aws.String(scheme),
			ListenerDescriptions: listenerDescriptions,
			HealthCheck: &elb.HealthCheck{
				Target:             aws.String("TCP:80"),
				Interval:           aws.Int64(30),
				Timeout:            aws.Int64(5),
				HealthyThreshold:   aws.Int64(2),
				UnhealthyThreshold: aws.Int64(2),
			},
		},
	}
}

type InstanceState struct {
	*elb.InstanceState
}

func NewInstanceState() *InstanceState {
	return &InstanceState{&elb.InstanceState{}}
}

func NewListener(instancePort int64, instanceProtocol string, lbPort int64, lbProtocol, certificate string) *Listener {
	listener := &Listener{
		&elb.Listener{
			InstancePort:     aws.Int64(instancePort),
			InstanceProtocol: aws.String(instanceProtocol),
			LoadBalancerPort: aws.Int64(lbPort),
			Protocol:         aws.String(lbProtocol),
		},
	}

	if certificate != "" {
		listener.SSLCertificateId = aws.String(certificate)
	}

	return listener
}

func NewHealthCheck(target string, interval, timeout, healthyThresh, unhealthyThresh int64) *HealthCheck {
	return &HealthCheck{
		&elb.HealthCheck{
			Target:             aws.String(target),
			Interval:           aws.Int64(interval),
			Timeout:            aws.Int64(timeout),
			UnhealthyThreshold: aws.Int64(unhealthyThresh),
			HealthyThreshold:   aws.Int64(healthyThresh),
		},
	}
}

type Tag struct {
	elb.Tag
}

type ELB struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (ELBInternal, error)
}

type ELBInternal interface {
	AttachLoadBalancerToSubnets(*elb.AttachLoadBalancerToSubnetsInput) (*elb.AttachLoadBalancerToSubnetsOutput, error)
	CreateLoadBalancer(input *elb.CreateLoadBalancerInput) (output *elb.CreateLoadBalancerOutput, err error)
	ConfigureHealthCheck(input *elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error)
	DescribeLoadBalancers(input *elb.DescribeLoadBalancersInput) (*elb.DescribeLoadBalancersOutput, error)
	DescribeInstanceHealth(input *elb.DescribeInstanceHealthInput) (*elb.DescribeInstanceHealthOutput, error)
	DeleteLoadBalancer(input *elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error)
	RegisterInstancesWithLoadBalancer(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
	DeregisterInstancesFromLoadBalancer(input *elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error)
	CreateLoadBalancerListeners(input *elb.CreateLoadBalancerListenersInput) (*elb.CreateLoadBalancerListenersOutput, error)
	DeleteLoadBalancerListeners(input *elb.DeleteLoadBalancerListenersInput) (*elb.DeleteLoadBalancerListenersOutput, error)
}

func NewELB(credProvider provider.CredProvider, region string) (Provider, error) {
	elb := ELB{
		credProvider,
		region,
		func() (ELBInternal, error) {
			return Connect(credProvider, region)
		},
	}

	_, err := elb.Connect()
	if err != nil {
		return nil, err
	}

	return &elb, nil
}

func Connect(credProvider provider.CredProvider, region string) (ELBInternal, error) {
	connection, err := provider.GetELBConnection(credProvider, region)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (this *ELB) CreateLoadBalancer(loadBalancerName, scheme string, securityGroups, subnets []*string, listeners []*Listener) (*string, error) {
	if len(subnets) == 0 {
		return nil, fmt.Errorf("Must specify at least 1 subnet")
	}

	awsListeners := make([]*elb.Listener, 0)
	for _, val := range listeners {
		awsListeners = append(awsListeners, val.Listener)
	}

	input := &elb.CreateLoadBalancerInput{
		Listeners:        awsListeners,
		LoadBalancerName: aws.String(loadBalancerName),
		Scheme:           aws.String(scheme),
		SecurityGroups:   securityGroups,
		Subnets:          subnets,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.CreateLoadBalancer(input)
	if err != nil {
		return nil, err
	}

	return out.DNSName, nil
}

func (this *ELB) DescribeLoadBalancers() ([]*LoadBalancerDescription, error) {
	input := &elb.DescribeLoadBalancersInput{}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.DescribeLoadBalancers(input)
	if err != nil {
		return nil, err
	}

	descriptions := []*LoadBalancerDescription{}
	for _, desc := range out.LoadBalancerDescriptions {
		descriptions = append(descriptions, &LoadBalancerDescription{desc})
	}

	return descriptions, nil
}

func (this *ELB) DescribeLoadBalancer(loadBalancerName string) (*LoadBalancerDescription, error) {
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{&loadBalancerName},
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.DescribeLoadBalancers(input)
	var desc *LoadBalancerDescription
	if err == nil {
		if len(out.LoadBalancerDescriptions) > 0 {
			desc = &LoadBalancerDescription{out.LoadBalancerDescriptions[0]}
		}
	}

	return desc, err
}

func (this *ELB) DescribeInstanceHealth(loadBalancerName string) ([]*InstanceState, error) {
	input := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: &loadBalancerName,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	out, err := connection.DescribeInstanceHealth(input)
	states := []*InstanceState{}
	if err == nil {
		for _, state := range out.InstanceStates {
			states = append(states, &InstanceState{state})
		}
	}

	return states, err
}

func (this *ELB) ConfigureHealthCheck(loadBalancerName string, check *HealthCheck) error {
	input := &elb.ConfigureHealthCheckInput{
		HealthCheck:      check.HealthCheck,
		LoadBalancerName: aws.String(loadBalancerName),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.ConfigureHealthCheck(input)
	return err
}

func (this *ELB) DeleteLoadBalancer(loadBalancerName string) error {
	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(loadBalancerName),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	_, err = connection.DeleteLoadBalancer(input)
	return err
}

func (this *ELB) RegisterInstancesWithLoadBalancer(loadBalancerName string, instances []string) error {
	awsInstances := make([]*elb.Instance, 0)
	for _, val := range instances {
		awsInstances = append(awsInstances, &elb.Instance{InstanceId: aws.String(val)})
	}

	input := &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(loadBalancerName),
		Instances:        awsInstances,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.RegisterInstancesWithLoadBalancer(input); err != nil {
		return err
	}

	return nil
}

func (this *ELB) DeregisterInstancesFromLoadBalancer(loadBalancerName string, instances []string) error {
	awsInstances := make([]*elb.Instance, 0)
	for _, val := range instances {
		awsInstances = append(awsInstances, &elb.Instance{InstanceId: aws.String(val)})
	}

	input := &elb.DeregisterInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(loadBalancerName),
		Instances:        awsInstances,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.DeregisterInstancesFromLoadBalancer(input); err != nil {
		return err
	}

	return nil
}

func (this *ELB) CreateLoadBalancerListeners(loadBalancerName string, listeners []*Listener) error {
	awsListeners := []*elb.Listener{}
	for _, listener := range listeners {
		awsListeners = append(awsListeners, listener.Listener)
	}

	input := &elb.CreateLoadBalancerListenersInput{
		LoadBalancerName: aws.String(loadBalancerName),
		Listeners:        awsListeners,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.CreateLoadBalancerListeners(input); err != nil {
		return err
	}

	return nil
}

func (this *ELB) DeleteLoadBalancerListeners(loadBalancerName string, listeners []*Listener) error {
	awsPorts := []*int64{}
	for _, listener := range listeners {
		awsPorts = append(awsPorts, listener.LoadBalancerPort)
	}

	input := &elb.DeleteLoadBalancerListenersInput{
		LoadBalancerName:  aws.String(loadBalancerName),
		LoadBalancerPorts: awsPorts,
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.DeleteLoadBalancerListeners(input); err != nil {
		return err
	}

	return nil
}
