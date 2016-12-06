package ecsbackend

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gitlab.imshealth.com/xfra/layer0/api/backend"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ec2"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/waitutils"
	"strings"
	"time"
)

const MAX_SERVICE_CREATE_RETRIES = 10

type ECSServiceManager struct {
	ECS            ecs.Provider
	EC2            ec2.Provider
	CloudWatchLogs cloudwatchlogs.Provider
	Backend        backend.Backend
	ClusterScaler  ClusterScaler
	Clock          waitutils.Clock
}

func NewECSServiceManager(
	ecsProvider ecs.Provider,
	ec2Provider ec2.Provider,
	cloudWatchLogsProvider cloudwatchlogs.Provider,
	clusterScaler ClusterScaler,
	backend backend.Backend,
) *ECSServiceManager {
	return &ECSServiceManager{
		ECS:            ecsProvider,
		EC2:            ec2Provider,
		CloudWatchLogs: cloudWatchLogsProvider,
		Backend:        backend,
		ClusterScaler:  clusterScaler,
		Clock:          waitutils.RealClock{},
	}
}

func (this *ECSServiceManager) ListServices() ([]*models.Service, error) {
	descriptions, err := this.ECS.Helper_DescribeServices(id.PREFIX)
	if err != nil {
		return nil, err
	}

	models := []*models.Service{}
	for _, description := range descriptions {
		if name := *description.ServiceName; strings.HasPrefix(name, id.PREFIX) {
			models = append(models, this.populateModel(description))
		}
	}

	return models, nil
}

func (this *ECSServiceManager) GetService(environmentID, serviceID string) (*models.Service, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsServiceID := id.L0ServiceID(serviceID).ECSServiceID()

	description, err := this.ECS.DescribeService(ecsEnvironmentID.String(), ecsServiceID.String())
	if err != nil {
		if ContainsErrMsg(err, "Service Not Found") {
			err := fmt.Errorf("Service with id '%s' does not exist", serviceID)
			return nil, errors.New(errors.InvalidServiceID, err)
		}

		return nil, err
	}

	return this.populateModel(description), nil
}

func (this *ECSServiceManager) UpdateService(
	environmentID string,
	serviceID string,
	deployID string,
) (*models.Service, error) {
	if err := this.updateService(environmentID, serviceID, deployID); err != nil {
		return nil, err
	}

	return this.GetService(environmentID, serviceID)
}

func (this *ECSServiceManager) updateService(
	environmentID string,
	serviceID string,
	deployID string,
) error {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsServiceID := id.L0ServiceID(serviceID).ECSServiceID()
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

	// trigger the scaling algorithm first or the service we are about to create gets
	// included in the pending count of the cluster
	if _, _, err := this.ClusterScaler.TriggerScalingAlgorithm(ecsEnvironmentID, &ecsDeployID, 1); err != nil {
		return err
	}

	if err := this.ECS.UpdateService(
		ecsEnvironmentID.String(),
		ecsServiceID.String(),
		stringp(ecsDeployID.TaskDefinition()),
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (this *ECSServiceManager) DeleteService(environmentID, serviceID string) error {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsServiceID := id.L0ServiceID(serviceID).ECSServiceID()
	desiredCount := int64(0)

	if err := this.ECS.UpdateService(ecsEnvironmentID.String(), ecsServiceID.String(), nil, &desiredCount); err != nil {
		return err
	}

	log.Debugf("Waiting for service to stop")
	if err := this.waitUntilServiceStopped(ecsEnvironmentID, ecsServiceID); err != nil {
		return err
	}

	if err := this.ECS.DeleteService(ecsEnvironmentID.String(), ecsServiceID.String()); err != nil {
		return err
	}

	return nil
}

func (this *ECSServiceManager) CreateService(
	serviceName,
	environmentID,
	deployID,
	loadBalancerID string,
) (*models.Service, error) {

	// we generate a hashed id for services since aws does not enforce unique service names
	serviceID := id.GenerateHashedEntityID(serviceName)

	var loadBalancerContainers []*ecs.LoadBalancer
	var loadBalancerRole *string
	if loadBalancerID != "" {
		ecsLoadBalancerID := id.L0LoadBalancerID(loadBalancerID).ECSLoadBalancerID()
		ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

		loadBalancerContainer, err := this.getLoadBalancerContainer(ecsLoadBalancerID, ecsDeployID)
		if err != nil {
			return nil, err
		}

		loadBalancerContainers = []*ecs.LoadBalancer{loadBalancerContainer}
		loadBalancerRole = stringp(ecsLoadBalancerID.RoleName())
	}

	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsServiceID := id.L0ServiceID(serviceID).ECSServiceID()
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()
	desiredCount := 1

	// trigger the scaling algorithm first or the service we are about to create gets
	// included in the pending count of the cluster
	if _, _, err := this.ClusterScaler.TriggerScalingAlgorithm(ecsEnvironmentID, &ecsDeployID, desiredCount); err != nil {
		return nil, err
	}

	var service *ecs.Service
	var attempts int
	check := func() (bool, error) {
		attempts++

		svc, err := this.ECS.CreateService(
			ecsEnvironmentID.String(),
			ecsServiceID.String(),
			ecsDeployID.TaskDefinition(),
			int64(desiredCount),
			loadBalancerContainers,
			loadBalancerRole,
		)

		switch {
		case err == nil:
			service = svc
			return true, nil

		case ContainsErrCode(err, "ClusterNotFoundException"):
			return false, errors.Newf(errors.InvalidEnvironmentID, "Environment with id '%s' was not found", environmentID)

		case ContainsErrMsg(err, "Creation of service was not idempotent"):
			return false, errors.Newf(errors.InvalidServiceID, "Service with name '%s' already exists", serviceName)

		case ContainsErrMsg(err, "unable to assume role and validate the listeners configured on your load balancer"):
			// must be a real loadbalancer-deploy mismatch, return this error
			// instead of the waiter's max retry attempt error
			if attempts == MAX_SERVICE_CREATE_RETRIES {
				return false, err
			}

			// loadbalancer's iam role is probably still propagating
			// log error and try again
			log.Warning(err)
			return false, nil

		default:
			return false, err
		}
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("Wait for service create %s", ecsServiceID),
		Retries: MAX_SERVICE_CREATE_RETRIES,
		Delay:   time.Second * 5,
		Clock:   this.Clock,
		Check:   check,
	}

	if err := waiter.Wait(); err != nil {
		return nil, err
	}

	return this.populateModel(service), nil
}

func (this *ECSServiceManager) getLoadBalancerContainer(ecsLoadBalancerID id.ECSLoadBalancerID, ecsDeployID id.ECSDeployID) (*ecs.LoadBalancer, error) {
	loadBalancer, err := this.Backend.GetLoadBalancer(ecsLoadBalancerID.L0LoadBalancerID())
	if err != nil {
		return nil, err
	}

	deploy, err := this.ECS.DescribeTaskDefinition(ecsDeployID.TaskDefinition())
	if err != nil {
		return nil, err
	}

	for _, container := range deploy.ContainerDefinitions {
		for _, containerPortMap := range container.PortMappings {
			for _, lbPort := range loadBalancer.Ports {
				if *containerPortMap.HostPort == lbPort.ContainerPort {
					loadBalancerContainer := ecs.NewLoadBalancer(
						*container.Name,
						*containerPortMap.ContainerPort,
						ecsLoadBalancerID.String())

					return loadBalancerContainer, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("No containers defined that listen on a port that is mapped by the load balancer")
}

func (this *ECSServiceManager) ScaleService(environmentID string, serviceID string, count int) (*models.Service, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsServiceID := id.L0ServiceID(serviceID).ECSServiceID()

	service, err := this.ECS.DescribeService(ecsEnvironmentID.String(), ecsServiceID.String())
	if err != nil {
		return nil, err
	}

	newTasksNeeded := count - int(*service.DesiredCount)
	if newTasksNeeded > 0 {
		// only trigger scaling when we need new tasks
		// count on the RightSizer to scale down the cluster next time it runs

		ecsDeployID := id.TaskDefinitionToECSDeployID(*service.TaskDefinition)

		if _, _, err := this.ClusterScaler.TriggerScalingAlgorithm(
			ecsEnvironmentID,
			&ecsDeployID,
			newTasksNeeded,
		); err != nil {
			return nil, err
		}
	}

	desiredCount := int64(count)
	if err := this.ECS.UpdateService(ecsEnvironmentID.String(), ecsServiceID.String(), nil, &desiredCount); err != nil {
		return nil, err
	}

	return this.GetService(environmentID, serviceID)
}

func (this *ECSServiceManager) GetServiceLogs(environmentID, serviceID string, tail int) ([]*models.LogFile, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()

	service, err := this.GetService(environmentID, serviceID)
	if err != nil {
		return nil, err
	}

	taskARNs := []*string{}
	for _, deployment := range service.Deployments {
		arns, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(deployment.DeploymentID))
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, arns...)
	}

	return GetLogs(this.CloudWatchLogs, taskARNs, tail)
}

func (this *ECSServiceManager) waitUntilServiceStopped(ecsEnvironmentID id.ECSEnvironmentID, ecsServiceID id.ECSServiceID) error {
	check := func() (bool, error) {
		service, err := this.ECS.DescribeService(ecsEnvironmentID.String(), ecsServiceID.String())
		if err != nil {
			return false, err
		}

		return int64OrZero(service.RunningCount) == 0, nil
	}

	waiter := waitutils.Waiter{
		Name:    fmt.Sprintf("Stop Service %s", ecsServiceID),
		Retries: 10,
		Delay:   time.Second * 4,
		Clock:   this.Clock,
		Check:   check,
	}

	return waiter.Wait()
}

func (this *ECSServiceManager) populateModel(service *ecs.Service) *models.Service {
	ecsEnvironmentID := id.ClusterARNToECSEnvironmentID(*service.ClusterArn)

	deployments := []models.Deployment{}
	for _, deployment := range service.Deployments {
		ecsDeployID := id.TaskDefinitionARNToECSDeployID(*deployment.TaskDefinition)

		model := models.Deployment{
			DeploymentID: *deployment.Id,
			Created:      *deployment.CreatedAt,
			Updated:      *deployment.UpdatedAt,
			Status:       *deployment.Status,
			PendingCount: *deployment.PendingCount,
			RunningCount: *deployment.RunningCount,
			DesiredCount: *deployment.DesiredCount,
			DeployID:     ecsDeployID.L0DeployID(),
		}

		deployments = append(deployments, model)
	}

	var loadBalancerID string
	if len(service.LoadBalancers) > 0 {
		ecsLoadBalancerName := *service.LoadBalancers[0].LoadBalancerName
		loadBalancerID = id.ECSLoadBalancerID(ecsLoadBalancerName).L0LoadBalancerID()
	}

	return &models.Service{
		ServiceID:      id.ECSServiceID(*service.ServiceName).L0ServiceID(),
		EnvironmentID:  ecsEnvironmentID.L0EnvironmentID(),
		LoadBalancerID: loadBalancerID,
		DesiredCount:   *service.DesiredCount,
		RunningCount:   *service.RunningCount,
		PendingCount:   *service.PendingCount,
		Deployments:    deployments,
	}
}
