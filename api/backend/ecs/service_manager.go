package ecsbackend

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/waitutils"
)

const MAX_SERVICE_CREATE_RETRIES = 10

type ECSServiceManager struct {
	ECS            ecs.Provider
	EC2            ec2.Provider
	CloudWatchLogs cloudwatchlogs.Provider
	Backend        backend.Backend
	Clock          waitutils.Clock
}

func NewECSServiceManager(
	ecsProvider ecs.Provider,
	ec2Provider ec2.Provider,
	cloudWatchLogsProvider cloudwatchlogs.Provider,
	backend backend.Backend,
) *ECSServiceManager {
	return &ECSServiceManager{
		ECS:            ecsProvider,
		EC2:            ec2Provider,
		CloudWatchLogs: cloudWatchLogsProvider,
		Backend:        backend,
		Clock:          waitutils.RealClock{},
	}
}

func (this *ECSServiceManager) ListServices() ([]id.ECSServiceID, error) {
	clusterNames, err := this.Backend.ListEnvironments()
	if err != nil {
		return nil, err
	}

	serviceIDs := []id.ECSServiceID{}
	for _, clusterName := range clusterNames {
		clusterServiceIDs, err := this.ECS.ListClusterServiceNames(clusterName.String(), id.PREFIX)
		if err != nil {
			return nil, err
		}

		for _, serviceID := range clusterServiceIDs {
			serviceIDs = append(serviceIDs, id.ECSServiceID(serviceID))
		}
	}

	return serviceIDs, nil
}

func (this *ECSServiceManager) GetService(environmentID, serviceID string) (*models.Service, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsServiceID := id.L0ServiceID(serviceID).ECSServiceID()

	description, err := this.ECS.DescribeService(ecsEnvironmentID.String(), ecsServiceID.String())
	if err != nil {
		if ContainsErrMsg(err, "Service Not Found") {
			err := fmt.Errorf("Service with id '%s' does not exist", serviceID)
			return nil, errors.New(errors.ServiceDoesNotExist, err)
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

	service, err := this.GetService(environmentID, serviceID)
	if err != nil {
		return err
	}

	taskARNs := []*string{}
	for _, deployment := range service.Deployments {
		arns, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(deployment.DeploymentID))
		if err != nil {
			return err
		}

		taskARNs = append(taskARNs, arns...)
	}

	for _, arn := range taskARNs {
		if err := this.ECS.StopTask(ecsEnvironmentID.String(), "Service deleted by User", *arn); err != nil {
			log.Warnf("Stop Task for Service '%s' had error: %v", serviceID, err)
		}
	}

	if err := this.ECS.UpdateService(ecsEnvironmentID.String(), ecsServiceID.String(), nil, &desiredCount); err != nil {
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

	count64 := int64(count)
	if pint64(service.DesiredCount) != count64 {
		if err := this.ECS.UpdateService(ecsEnvironmentID.String(), ecsServiceID.String(), nil, int64p(count64)); err != nil {
			return nil, err
		}
	}

	return this.GetService(environmentID, serviceID)
}

func (this *ECSServiceManager) GetServiceLogs(environmentID, serviceID, start, end string, tail int) ([]*models.LogFile, error) {
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

	return GetLogs(this.CloudWatchLogs, taskARNs, start, end, tail)
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
