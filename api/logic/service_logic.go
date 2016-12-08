package logic

import (
	"fmt"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type ServiceLogic interface {
	ListServices() ([]*models.Service, error)
	GetService(serviceID string) (*models.Service, error)
	CreateService(req models.CreateServiceRequest) (*models.Service, error)
	DeleteService(serviceID string) error
	UpdateService(serviceID string, req models.UpdateServiceRequest) (*models.Service, error)
	ScaleService(serviceID string, size int) (*models.Service, error)
	GetServiceLogs(serviceID string, tail int) ([]*models.LogFile, error)
}

type L0ServiceLogic struct {
	Logic
	DeployLogic DeployLogic
}

func NewL0ServiceLogic(logic Logic, deployLogic DeployLogic) *L0ServiceLogic {
	return &L0ServiceLogic{
		Logic:       logic,
		DeployLogic: deployLogic,
	}
}

func (this *L0ServiceLogic) ListServices() ([]*models.Service, error) {
	services, err := this.Backend.ListServices()
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if err := this.populateModel(service); err != nil {
			return nil, err
		}
	}

	return services, nil
}

func (this *L0ServiceLogic) GetService(serviceID string) (*models.Service, error) {
	environmentID, err := this.getEnvironmentID(serviceID)
	if err != nil {
		return nil, err
	}

	service, err := this.Backend.GetService(environmentID, serviceID)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (this *L0ServiceLogic) DeleteService(serviceID string) error {
	environmentID, err := this.getEnvironmentID(serviceID)
	if err != nil {
		return err
	}

	if err := this.Backend.DeleteService(environmentID, serviceID); err != nil {
		return err
	}

	if err := this.deleteEntityTags(serviceID, "service"); err != nil {
		return err
	}

	return nil
}

func (this *L0ServiceLogic) ScaleService(serviceID string, size int) (*models.Service, error) {
	environmentID, err := this.getEnvironmentID(serviceID)
	if err != nil {
		return nil, err
	}

	service, err := this.Backend.ScaleService(environmentID, serviceID, size)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (this *L0ServiceLogic) UpdateService(serviceID string, req models.UpdateServiceRequest) (*models.Service, error) {
	environmentID, err := this.getEnvironmentID(serviceID)
	if err != nil {
		return nil, err
	}

	service, err := this.Backend.UpdateService(environmentID, serviceID, req.DeployID)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (this *L0ServiceLogic) CreateService(req models.CreateServiceRequest) (*models.Service, error) {
	if req.EnvironmentID == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentID not specified")
	}

	if req.ServiceName == "" {
		return nil, errors.Newf(errors.MissingParameter, "ServiceName not specified")
	}

	exists, err := this.doesServiceTagExist(req.EnvironmentID, req.ServiceName)
	if err != nil {
		return nil, err
	}

	if exists {
		err := fmt.Errorf("Service with name '%s' already exists in Environment '%s'", req.ServiceName, req.EnvironmentID)
		return nil, errors.New(errors.InvalidServiceName, err)
	}

	service, err := this.Backend.CreateService(
		req.ServiceName,
		req.EnvironmentID,
		req.DeployID,
		req.LoadBalancerID)
	if err != nil {
		return service, err
	}

	serviceID := service.ServiceID
	if err := this.upsertTagf(serviceID, "service", "name", req.ServiceName); err != nil {
		return service, err
	}

	environmentID := service.EnvironmentID
	if err := this.upsertTagf(serviceID, "service", "environment_id", environmentID); err != nil {
		return service, err
	}

	if loadBalancerID := req.LoadBalancerID; loadBalancerID != "" {
		if err := this.upsertTagf(serviceID, "service", "load_balancer_id", loadBalancerID); err != nil {
			return service, err
		}
	}

	if err := this.populateModel(service); err != nil {
		return service, err
	}

	return service, nil
}

func (this *L0ServiceLogic) GetServiceLogs(serviceID string, tail int) ([]*models.LogFile, error) {
	environmentID, err := this.getEnvironmentID(serviceID)
	if err != nil {
		return nil, err
	}

	logs, err := this.Backend.GetServiceLogs(environmentID, serviceID, tail)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (this *L0ServiceLogic) getEnvironmentID(serviceID string) (string, error) {
	filter := map[string]string{
		"type": "service",
		"id":   serviceID,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return "", err
	}

	for _, tag := range rangeTags(tags) {
		if tag.Key == "environment_id" {
			return tag.Value, nil
		}

	}

	return "", fmt.Errorf("Failed to find Environment ID for Service %s", serviceID)
}

func (this *L0ServiceLogic) doesServiceTagExist(environmentID, name string) (bool, error) {
	filter := map[string]string{
		"type":           "service",
		"environment_id": environmentID,
		"name":           name,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return false, err
	}

	return len(tags) > 0, nil
}

func (this *L0ServiceLogic) populateModel(model *models.Service) error {
	filter := map[string]string{
		"type": "service",
		"id":   model.ServiceID,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	for _, tag := range rangeTags(tags) {
		switch tag.Key {
		case "environment_id":
			model.EnvironmentID = tag.Value
		case "load_balancer_id":
			model.LoadBalancerID = tag.Value
		case "name":
			model.ServiceName = tag.Value
		}
	}

	// todo: make this errors for all environment-scoped entities
	if model.EnvironmentID != "" {
		filter := map[string]string{
			"type": "environment",
			"id":   model.EnvironmentID,
		}

		tags, err := this.TagData.GetTags(filter)
		if err != nil {
			return err
		}

		for _, tag := range rangeTags(tags) {
			if tag.Key == "name" {
				model.EnvironmentName = tag.Value
				break
			}
		}
	}

	// todo: lookupEntityName could be made generic
	if model.LoadBalancerID != "" {
		filter := map[string]string{
			"type": "load_balancer",
			"id":   model.LoadBalancerID,
		}

		tags, err := this.TagData.GetTags(filter)
		if err != nil {
			return err
		}

		for _, tag := range rangeTags(tags) {
			if tag.Key == "name" {
				model.LoadBalancerName = tag.Value
				break
			}
		}
	}

	deployments := []models.Deployment{}
	for _, deploy := range model.Deployments {
		filter = map[string]string{
			"type": "deploy",
			"id":   deploy.DeployID,
		}

		tags, err := this.TagData.GetTags(filter)
		if err != nil {
			return err
		}

		for _, tag := range rangeTags(tags) {
			switch tag.Key {
			case "name":
				deploy.DeployName = tag.Value
			case "version":
				deploy.DeployVersion = tag.Value
			}
		}

		deployments = append(deployments, deploy)
	}

	model.Deployments = deployments

	return nil
}
