package ecsbackend

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	awselb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/elb"
	"github.com/quintilesims/layer0/common/aws/iam"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/waitutils"
)

type ECSLoadBalancerManager struct {
	EC2     ec2.Provider
	ELB     elb.Provider
	IAM     iam.Provider
	Backend backend.Backend
	Clock   waitutils.Clock
}

func NewECSLoadBalancerManager(ec2 ec2.Provider, elb elb.Provider, iam iam.Provider, backend backend.Backend) *ECSLoadBalancerManager {
	return &ECSLoadBalancerManager{
		EC2:     ec2,
		ELB:     elb,
		IAM:     iam,
		Backend: backend,
		Clock:   waitutils.RealClock{},
	}
}

func (e *ECSLoadBalancerManager) ListLoadBalancers() ([]*models.LoadBalancer, error) {
	loadBalancerDescriptions, err := e.ELB.DescribeLoadBalancers()
	if err != nil {
		return nil, err
	}

	loadBalancers := []*models.LoadBalancer{}
	for _, description := range loadBalancerDescriptions {
		if name := *description.LoadBalancerName; strings.HasPrefix(name, id.PREFIX) {
			ecsLoadBalancerID := id.ECSLoadBalancerID(name)
			loadBalancer := &models.LoadBalancer{
				LoadBalancerID: ecsLoadBalancerID.L0LoadBalancerID(),
			}

			loadBalancers = append(loadBalancers, loadBalancer)
		}
	}

	return loadBalancers, nil
}

func (e *ECSLoadBalancerManager) GetLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error) {
	ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()

	loadBalancer, err := e.ELB.DescribeLoadBalancer(ecsLoadBalancerID.String())
	if err != nil {
		if ContainsErrCode(err, "LoadBalancerNotFound") {
			err := fmt.Errorf("LoadBalancer with id '%s' does not exist", loadBalancerID)
			return nil, errors.New(errors.LoadBalancerDoesNotExist, err)
		}

		return nil, err
	}

	lbAttributes, err := e.ELB.DescribeLoadBalancerAttributes(aws.StringValue(loadBalancer.LoadBalancerName))
	if err != nil {
		if ContainsErrCode(err, "LoadBalancerAttributeNotFound") {
			err := fmt.Errorf("LoadBalancer attributes cannot be described for LoadBalancer with id '%s'", loadBalancerID)
			return nil, errors.New(errors.LoadBalancerAttributeNotFound, err)
		}
	}

	return e.populateModel(loadBalancer, lbAttributes), nil
}

func (e *ECSLoadBalancerManager) DeleteLoadBalancer(loadBalancerID string) error {
	ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()
	roleName := ecsLoadBalancerID.RoleName()

	policyList, err := e.IAM.ListRolePolicies(roleName)
	if err != nil && !ContainsErrCode(err, "NoSuchEntity") {
		return err
	}

	for _, name := range policyList {
		policy := stringOrEmpty(name)
		if err := e.IAM.DeleteRolePolicy(roleName, policy); err != nil {
			return err
		}
	}

	if err := e.waitUntilRolePoliciesDeleted(roleName); err != nil {
		return err
	}

	if err := e.IAM.DeleteRole(roleName); err != nil {
		if !ContainsErrCode(err, "NoSuchEntity") {
			return err
		}
	}

	if err := e.waitUntilRoleDeleted(roleName); err != nil {
		return err
	}

	if err := e.ELB.DeleteLoadBalancer(ecsLoadBalancerID.String()); err != nil {
		if !ContainsErrCode(err, "NoSuchEntity") {
			return err
		}
	}

	securityGroup, err := e.EC2.DescribeSecurityGroup(ecsLoadBalancerID.SecurityGroupName())
	if err != nil {
		return err
	}

	// wait a couple minutes for the elb and it's network interfaces to delete
	// todo: using waiters seems pretty verbose - should find a way to clean this up
	if securityGroup != nil {
		check := func() (bool, error) {
			if err := e.EC2.DeleteSecurityGroup(securityGroup); err == nil {
				return true, nil
			}

			return false, nil
		}

		waiter := waitutils.Waiter{
			Name:    fmt.Sprintf("SecurityGroup delete for '%s'", securityGroup),
			Retries: 50,
			Delay:   time.Second * 10,
			Clock:   e.Clock,
			Check:   check,
		}

		if err := waiter.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func (e *ECSLoadBalancerManager) waitUntilRolePoliciesDeleted(roleName string) error {
	check := func() (bool, error) {
		policies, err := e.IAM.ListRolePolicies(roleName)
		if err != nil && !ContainsErrCode(err, "NoSuchEntity") {
			return false, err
		}

		return len(policies) == 0, nil
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("Wait for deleted role policies %s", roleName),
		Retries: 50,
		Delay:   time.Second * 5,
		Clock:   e.Clock,
		Check:   check,
	}

	return waiter.Wait()
}

func (e *ECSLoadBalancerManager) waitUntilRoleDeleted(roleName string) error {
	check := func() (bool, error) {
		policies, err := e.IAM.ListRolePolicies(roleName)
		if err != nil && !ContainsErrCode(err, "NoSuchEntity") {
			return false, err
		}

		return len(policies) == 0, nil
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("Wait for deleted role %s", roleName),
		Retries: 50,
		Delay:   time.Second * 5,
		Clock:   e.Clock,
		Check:   check,
	}

	return waiter.Wait()
}

func (e *ECSLoadBalancerManager) populateModel(description *elb.LoadBalancerDescription, lbAttributes *elb.LoadBalancerAttributes) *models.LoadBalancer {
	ecsLoadBalancerID := id.ECSLoadBalancerID(*description.LoadBalancerName)

	ports := []models.Port{}
	for _, listener := range description.ListenerDescriptions {
		ports = append(ports, e.listenerToPort(listener.Listener))
	}

	healthCheck := models.HealthCheck{
		Target:             *description.HealthCheck.Target,
		Interval:           int(*description.HealthCheck.Interval),
		Timeout:            int(*description.HealthCheck.Timeout),
		HealthyThreshold:   int(*description.HealthCheck.HealthyThreshold),
		UnhealthyThreshold: int(*description.HealthCheck.UnhealthyThreshold),
	}

	model := &models.LoadBalancer{
		LoadBalancerID: ecsLoadBalancerID.L0LoadBalancerID(),
		Ports:          ports,
		IsPublic:       stringOrEmpty(description.Scheme) == "internet-facing",
		URL:            stringOrEmpty(description.DNSName),
		HealthCheck:    healthCheck,
		IdleTimeout:    int(aws.Int64Value(lbAttributes.ConnectionSettings.IdleTimeout)),
	}

	return model
}

func (e *ECSLoadBalancerManager) listenerToPort(listener *awselb.Listener) models.Port {
	port := models.Port{
		ContainerPort: *listener.InstancePort,
		HostPort:      *listener.LoadBalancerPort,
		// assuming LB protocol == Instance protocol
		Protocol: *listener.Protocol,
	}

	if listener.SSLCertificateId != nil {
		port.CertificateARN = *listener.SSLCertificateId
		port.CertificateName = id.CertificateARNToName(*listener.SSLCertificateId)
	}

	return port
}

func (e *ECSLoadBalancerManager) getSecurityGroupIDByName(securityGroupName string) (string, error) {
	securityGroup, err := e.EC2.DescribeSecurityGroup(securityGroupName)
	if err != nil {
		return "", err
	}

	if securityGroup == nil {
		return "", fmt.Errorf("Security group '%s' does not exist!", securityGroupName)
	}

	return *securityGroup.GroupId, nil
}

func (e *ECSLoadBalancerManager) CreateLoadBalancer(
	loadBalancerName,
	environmentID string,
	isPublic bool,
	ports []models.Port,
	healthCheck models.HealthCheck,
	idleTimeout int,
) (*models.LoadBalancer, error) {
	// we generate a hashed id for load balancers since aws does not enforce unique load balancer names
	loadBalancerID := id.GenerateHashedEntityID(loadBalancerName)
	ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	if err := e.createLoadBalancer(ecsLoadBalancerID, ecsEnvironmentID, isPublic, ports); err != nil {
		return nil, err
	}

	// Once the loadbalancer has been created, we can update its healthcheck
	if err := e.updateHealthCheck(ecsLoadBalancerID, healthCheck); err != nil {
		return nil, err
	}

	// Then set the load balancer's idle timeout
	if err := e.setIdleTimeout(ecsLoadBalancerID, idleTimeout); err != nil {
		return nil, err
	}

	model := &models.LoadBalancer{
		LoadBalancerID:   ecsLoadBalancerID.L0LoadBalancerID(),
		LoadBalancerName: loadBalancerName,
		EnvironmentID:    ecsEnvironmentID.L0EnvironmentID(),
		IsPublic:         isPublic,
		Ports:            ports,
		HealthCheck:      healthCheck,
		IdleTimeout:      idleTimeout,
	}

	return model, nil
}

func (e *ECSLoadBalancerManager) createLoadBalancer(
	ecsLoadBalancerID id.ECSLoadBalancerID,
	ecsEnvironmentID id.ECSEnvironmentID,
	isPublic bool,
	ports []models.Port,
) error {
	listeners := []*elb.Listener{}
	for _, port := range ports {
		listener, err := e.portToListener(port)
		if err != nil {
			return err
		}

		listeners = append(listeners, listener)
	}

	roleName := ecsLoadBalancerID.RoleName()
	if _, err := e.IAM.CreateRole(roleName, "ecs.amazonaws.com"); err != nil {
		if !ContainsErrCode(err, "EntityAlreadyExists") {
			return err
		}
	}

	policy, err := e.generateRolePolicy(ecsLoadBalancerID)
	if err != nil {
		return err
	}

	if err := e.IAM.PutRolePolicy(roleName, policy); err != nil {
		return err
	}

	// only public load balancers get an additional security group
	securityGroupIDs := []*string{}
	if isPublic {
		securityGroup, err := e.upsertSecurityGroup(ecsLoadBalancerID, ports)
		if err != nil {
			return err
		}

		securityGroupIDs = append(securityGroupIDs, securityGroup.GroupId)
	}

	environmentSecurityGroupID, err := e.getSecurityGroupIDByName(ecsEnvironmentID.SecurityGroupName())
	if err != nil {
		return fmt.Errorf("Failed to find environment Security Group: %v", err)
	}

	securityGroupIDs = append(securityGroupIDs, &environmentSecurityGroupID)

	scheme := "internal"
	if isPublic {
		scheme = "internet-facing"
	}

	subnets, _, err := e.getSubnetsAndAvailZones(isPublic)
	if err != nil {
		return err
	}

	if _, err := e.ELB.CreateLoadBalancer(ecsLoadBalancerID.String(), scheme, securityGroupIDs, subnets, listeners); err != nil {
		return err
	}

	return nil
}

func (e *ECSLoadBalancerManager) UpdateLoadBalancerHealthCheck(loadBalancerID string, healthCheck models.HealthCheck) (*models.LoadBalancer, error) {
	ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()
	if err := e.updateHealthCheck(ecsLoadBalancerID, healthCheck); err != nil {
		return nil, err
	}

	return e.GetLoadBalancer(loadBalancerID)
}

func (e *ECSLoadBalancerManager) UpdateLoadBalancerIdleTimeout(loadBalancerID string, idleTimeout int) (*models.LoadBalancer, error) {
	ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()
	if err := e.setIdleTimeout(ecsLoadBalancerID, idleTimeout); err != nil {
		return nil, err
	}

	return e.GetLoadBalancer(loadBalancerID)
}

func (e *ECSLoadBalancerManager) updateHealthCheck(ecsLoadBalancerID id.ECSLoadBalancerID, healthCheck models.HealthCheck) error {
	elbHealthCheck := elb.NewHealthCheck(
		healthCheck.Target,
		int64(healthCheck.Interval),
		int64(healthCheck.Timeout),
		int64(healthCheck.HealthyThreshold),
		int64(healthCheck.UnhealthyThreshold),
	)

	return e.ELB.ConfigureHealthCheck(ecsLoadBalancerID.String(), elbHealthCheck)
}

func (e *ECSLoadBalancerManager) setIdleTimeout(ecsLoadBalancerID id.ECSLoadBalancerID, idleTimeout int) error {
	return e.ELB.SetIdleTimeout(ecsLoadBalancerID.String(), idleTimeout)
}

func (e *ECSLoadBalancerManager) UpdateLoadBalancerPorts(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error) {
	model, err := e.GetLoadBalancer(loadBalancerID)
	if err != nil {
		return nil, err
	}

	ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()
	updatedPorts, err := e.updatePorts(ecsLoadBalancerID, model.IsPublic, model.Ports, ports)
	if err != nil {
		return nil, err
	}

	model.Ports = updatedPorts
	return model, nil
}

func (e *ECSLoadBalancerManager) updatePorts(ecsLoadBalancerID id.ECSLoadBalancerID, isPublic bool, currentPorts, requestedPorts []models.Port) ([]models.Port, error) {
	if reflect.DeepEqual(currentPorts, requestedPorts) {
		return currentPorts, nil
	}

	// remove first so we don't duplicate host ports
	listenersToRemove := []*elb.Listener{}
	for _, port := range portDifference(currentPorts, requestedPorts) {
		listener, err := e.portToListener(port)
		if err != nil {
			return nil, err
		}

		listenersToRemove = append(listenersToRemove, listener)
	}

	if len(listenersToRemove) > 0 {
		if err := e.ELB.DeleteLoadBalancerListeners(ecsLoadBalancerID.String(), listenersToRemove); err != nil {
			return nil, err
		}
	}

	listenersToAdd := []*elb.Listener{}
	for _, port := range portDifference(requestedPorts, currentPorts) {
		listener, err := e.portToListener(port)
		if err != nil {
			return nil, err
		}

		listenersToAdd = append(listenersToAdd, listener)
	}

	if len(listenersToAdd) > 0 {
		if err := e.ELB.CreateLoadBalancerListeners(ecsLoadBalancerID.String(), listenersToAdd); err != nil {
			return nil, err
		}
	}

	// only public load balancers have an additional security group
	if isPublic {
		if _, err := e.upsertSecurityGroup(ecsLoadBalancerID, requestedPorts); err != nil {
			return nil, err
		}
	}

	return requestedPorts, nil
}

func (e *ECSLoadBalancerManager) portToListener(port models.Port) (*elb.Listener, error) {
	hostProtocol := strings.ToUpper(port.Protocol)
	if hostProtocol != "SSL" && hostProtocol != "TCP" && hostProtocol != "HTTP" && hostProtocol != "HTTPS" {
		return nil, fmt.Errorf("Protocol '%s' is not valid", port.Protocol)
	}

	containerProtocol := hostProtocol

	// terminate https, ssl connections at elb
	if containerProtocol == "HTTPS" {
		containerProtocol = "HTTP"
	} else if containerProtocol == "SSL" {
		containerProtocol = "TCP"
	}

	// use cert arn if specified by the user
	// otherwise, if name is specified, convert it to an arn
	certificateARN := port.CertificateARN
	if certificateARN == "" && port.CertificateName != "" {
		arn, err := e.getCertificateARN(port.CertificateName)
		if err != nil {
			return nil, err
		}

		certificateARN = arn
	}

	listener := elb.NewListener(port.ContainerPort, containerProtocol, port.HostPort, hostProtocol, certificateARN)
	return listener, nil
}

func (e *ECSLoadBalancerManager) getCertificateARN(name string) (string, error) {
	certificates, err := e.IAM.ListCertificates()
	if err != nil {
		return "", err
	}

	for _, c := range certificates {
		if *c.ServerCertificateName == name {
			return *c.Arn, nil
		}
	}

	return "", fmt.Errorf("Certificate with name '%s' does not exist. ", name)
}

func (e *ECSLoadBalancerManager) upsertSecurityGroup(ecsLoadBalancerID id.ECSLoadBalancerID, ports []models.Port) (*ec2.SecurityGroup, error) {
	securityGroupName := ecsLoadBalancerID.SecurityGroupName()

	securityGroup, err := e.EC2.DescribeSecurityGroup(securityGroupName)
	if err != nil {
		return nil, err
	}

	if securityGroup == nil {
		desc := "Auto-generated Layer0 Load Balancer Security Group"
		vpcID := config.AWSVPCID()

		if _, err = e.EC2.CreateSecurityGroup(securityGroupName, desc, vpcID); err != nil {
			return nil, err
		}

		check := func() (bool, error) {
			securityGroup, err = e.EC2.DescribeSecurityGroup(securityGroupName)
			if err != nil {
				return false, err
			}

			return securityGroup != nil, nil
		}

		waiter := waitutils.Waiter{
			Name:    fmt.Sprintf("SecurityGroup setup for '%s'", ecsLoadBalancerID),
			Retries: 50,
			Delay:   time.Second * 10,
			Clock:   e.Clock,
			Check:   check,
		}

		if err := waiter.Wait(); err != nil {
			return nil, err
		}
	}

	currentIngressPorts := []int64{}
	requestedIngressPorts := []int64{}

	for _, permission := range securityGroup.IpPermissions {
		currentIngressPorts = append(currentIngressPorts, *permission.FromPort)
	}

	for _, port := range ports {
		requestedIngressPorts = append(requestedIngressPorts, int64(port.HostPort))
	}

	ingressesToRemove := []*ec2.SecurityGroupIngress{}
	for _, ingressPort := range ingressPortDifference(currentIngressPorts, requestedIngressPorts) {
		ingress := ec2.NewSecurityGroupIngress(*securityGroup.GroupId, "0.0.0.0/0", "TCP", int(ingressPort), int(ingressPort))
		ingressesToRemove = append(ingressesToRemove, ingress)
	}

	if len(ingressesToRemove) > 0 {
		log.Debug("Removing Ports: ", *securityGroup.GroupId, ingressesToRemove)
		if err := e.EC2.RevokeSecurityGroupIngress(ingressesToRemove); err != nil {
			return nil, err
		}
	}

	ingressesToAdd := []*ec2.SecurityGroupIngress{}
	for _, ingressPort := range ingressPortDifference(requestedIngressPorts, currentIngressPorts) {
		ingress := ec2.NewSecurityGroupIngress(*securityGroup.GroupId, "0.0.0.0/0", "TCP", int(ingressPort), int(ingressPort))
		ingressesToAdd = append(ingressesToAdd, ingress)
	}

	if len(ingressesToAdd) > 0 {
		log.Debug("Adding ports: ", *securityGroup.GroupId, ingressesToAdd)
		if err := e.EC2.AuthorizeSecurityGroupIngress(ingressesToAdd); err != nil {
			return nil, err
		}
	}

	return securityGroup, nil
}

// returns ports in "requested" that aren't in "current"
func portDifference(requested, current []models.Port) []models.Port {
	difference := []models.Port{}
	for _, r := range requested {
		var exists bool
		for _, c := range current {
			if reflect.DeepEqual(r, c) {
				exists = true
				break
			}
		}

		if !exists {
			difference = append(difference, r)
		}
	}

	return difference
}

// returns ports in "requested" that aren't in "current"
func ingressPortDifference(requested, current []int64) []int64 {
	difference := []int64{}
	for _, r := range requested {
		var exists bool
		for _, c := range current {
			if r == c {
				exists = true
				break
			}
		}

		if !exists {
			difference = append(difference, r)
		}
	}

	return difference
}

// this is awkward, strongly assumes that PrivateSubnets will be distributed across AZs,
// using each at most once.  We error out on bad config for now, in the future we'll
// need to do something to calculate which subnets to use based on where the instance
// got provisioned.

func (e *ECSLoadBalancerManager) getSubnetsAndAvailZones(public bool) ([]*string, []*string, error) {

	// todo: the majority of this function can be taken out, we essentially just need to split
	// config.Subnets() and return []string. AWS Handles the overlap error check for us already
	var subnets string
	if public {
		subnets = config.AWSPublicSubnets()
	} else {
		subnets = config.AWSPrivateSubnets()
	}

	subnetIDs := []*string{}
	availZones := []*string{}
	for _, subnetID := range strings.Split(subnets, ",") {
		subnet := strings.TrimSpace(subnetID)
		subnetIDs = append(subnetIDs, &subnet)

		description, err := e.EC2.DescribeSubnet(subnetID)
		if err != nil {
			return nil, nil, err
		}

		for _, zone := range availZones {
			if *description.AvailabilityZone == *zone {
				if public {
					err = fmt.Errorf("Public Subnets an availability zone overlap: %s", *zone)
				} else {
					err = fmt.Errorf("Private Subnets an availability zone overlap: %s", *zone)
				}

				return nil, nil, err
			}
		}

		availZones = append(availZones, description.AvailabilityZone)
	}

	return subnetIDs, availZones, nil
}

func (e *ECSLoadBalancerManager) generateRolePolicy(ecsLoadBalancerID id.ECSLoadBalancerID) (string, error) {
	// the default policy includes "ec2:AuthorizeSecurityGroupIngress" which we exclude
	// because we don't know why it's there
	policy := `
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:Describe*",
                "ec2:Describe*"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
                "elasticloadbalancing:RegisterInstancesWithLoadBalancer"
            ],
            "Resource": [
                "arn:aws:elasticloadbalancing:%s:%s:loadbalancer/%s"
            ]
        }

    ]
}`
	awsAccountID, err := e.IAM.GetAccountId()
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf(policy, config.AWSRegion(), awsAccountID, ecsLoadBalancerID.String())
	out = strings.Replace(out, "\n", "", -1) // AWS API requires no newlines
	return out, nil
}
