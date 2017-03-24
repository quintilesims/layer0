package logic

import (
	"fmt"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"time"
)

type ServiceLogic interface {
	ListServices() ([]*models.ServiceSummary, error)
	GetService(serviceID string) (*models.Service, error)
	CreateService(req models.CreateServiceRequest) (*models.Service, error)
	DeleteService(serviceID string) error
	UpdateService(serviceID string, req models.UpdateServiceRequest) (*models.Service, error)
	ScaleService(serviceID string, size int) (*models.Service, error)
	GetServiceLogs(serviceID string, tail int) ([]*models.LogFile, error)
}

type L0ServiceLogic struct {
	Logic
}

func NewL0ServiceLogic(logic Logic) *L0ServiceLogic {
	return &L0ServiceLogic{
		Logic: logic,
	}
}

func (this *L0ServiceLogic) ListServices() ([]*models.ServiceSummary, error) {
	services, err := this.Backend.ListServices()
	if err != nil {
		return nil, err
	}

	summaries := make([]*models.ServiceSummary, len(services))
	for i, service := range services {
		if err := this.populateModel(service); err != nil {
			return nil, err
		}

		summaries[i] = &models.ServiceSummary{
			ServiceID:       service.ServiceID,
			ServiceName:     service.ServiceName,
			EnvironmentID:   service.EnvironmentID,
			EnvironmentName: service.EnvironmentName,
		}
	}

	return summaries, nil
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

	if err := this.deleteEntityTags("service", serviceID); err != nil {
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

	this.Logic.Scaler.ScheduleRun(service.EnvironmentID, time.Second*10)
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

	this.Logic.Scaler.ScheduleRun(service.EnvironmentID, time.Second*10)

	return service, nil
}

func (this *L0ServiceLogic) CreateService(req models.CreateServiceRequest) (*models.Service, error) {
	if req.EnvironmentID == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentID not specified")
	}

	if req.ServiceName == "" {
		return nil, errors.Newf(errors.MissingParameter, "ServiceName not specified")
	}

	if req.DeployID == "" {
		return nil, errors.Newf(errors.MissingParameter, "DeployID not specified")
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

	environmentID := req.EnvironmentID
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

	this.Logic.Scaler.ScheduleRun(service.EnvironmentID, time.Second*10)

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
	tags, err := this.TagStore.SelectByQuery("service", serviceID)
	if err != nil {
		return "", err
	}

	if tag := tags.WithKey("environment_id").First(); tag != nil {
		return tag.Value, nil
	}

	services, err := this.ListServices()
	if err != nil {
		return "", err
	}

	for _, service := range services {
		if service.ServiceID == serviceID {
			return service.EnvironmentID, nil
		}
	}

	return "", errors.Newf(errors.InvalidServiceID, "Service %s does not exist", serviceID)
}

func (this *L0ServiceLogic) doesServiceTagExist(environmentID, name string) (bool, error) {
	tags, err := this.TagStore.SelectByQuery("service", "")
	if err != nil {
		return false, err
	}

	ewts := tags.GroupByEntity().
		WithKey("environment_id").
		WithValue(environmentID).
		WithKey("name").
		WithValue(name)

	return len(ewts) > 0, nil
}

func (this *L0ServiceLogic) populateModel(model *models.Service) error {
	tags, err := this.TagStore.SelectByQuery("service", model.ServiceID)
	if err != nil {
		return err
	}

	if tag := tags.WithKey("environment_id").First(); tag != nil {
		model.EnvironmentID = tag.Value
	}

	if tag := tags.WithKey("load_balancer_id").First(); tag != nil {
		model.LoadBalancerID = tag.Value
	}

	if tag := tags.WithKey("name").First(); tag != nil {
		model.ServiceName = tag.Value
	}

	if model.EnvironmentID != "" {
		tags, err := this.TagStore.SelectByQuery("environment", model.EnvironmentID)
		if err != nil {
			return err
		}

		if tag := tags.WithKey("name").First(); tag != nil {
			model.EnvironmentName = tag.Value
		}
	}

	if model.LoadBalancerID != "" {
		tags, err := this.TagStore.SelectByQuery("load_balancer", model.LoadBalancerID)
		if err != nil {
			return err
		}

		if tag := tags.WithKey("name").First(); tag != nil {
			model.LoadBalancerName = tag.Value
		}
	}

	deployments := []models.Deployment{}
	for _, deploy := range model.Deployments {
		tags, err := this.TagStore.SelectByQuery("deploy", deploy.DeployID)
		if err != nil {
			return err
		}

		if tag := tags.WithKey("name").First(); tag != nil {
			deploy.DeployName = tag.Value
		}

		if tag := tags.WithKey("version").First(); tag != nil {
			deploy.DeployVersion = tag.Value
		}

		deployments = append(deployments, deploy)
	}

	model.Deployments = deployments

	return nil
}
