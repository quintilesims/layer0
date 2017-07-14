package logic

import (
	"fmt"

	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type LoadBalancerLogic interface {
	ListLoadBalancers() ([]*models.LoadBalancerSummary, error)
	GetLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error)
	DeleteLoadBalancer(loadBalancerID string) error
	CreateLoadBalancer(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error)
	UpdateLoadBalancerPorts(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error)
	UpdateLoadBalancerHealthCheck(loadBalancerID string, healthCheck models.HealthCheck) (*models.LoadBalancer, error)
}

type L0LoadBalancerLogic struct {
	Logic
}

func NewL0LoadBalancerLogic(logic Logic) *L0LoadBalancerLogic {
	return &L0LoadBalancerLogic{
		Logic: logic,
	}
}

func (l *L0LoadBalancerLogic) ListLoadBalancers() ([]*models.LoadBalancerSummary, error) {
	loadBalancers, err := l.Backend.ListLoadBalancers()
	if err != nil {
		return nil, err
	}

	summaries := make([]*models.LoadBalancerSummary, len(loadBalancers))
	for i, loadBalancer := range loadBalancers {
		if err := l.populateModel(loadBalancer); err != nil {
			return nil, err
		}

		summaries[i] = &models.LoadBalancerSummary{
			LoadBalancerID:   loadBalancer.LoadBalancerID,
			LoadBalancerName: loadBalancer.LoadBalancerName,
			EnvironmentID:    loadBalancer.EnvironmentID,
			EnvironmentName:  loadBalancer.EnvironmentName,
		}
	}

	return summaries, nil
}

func (l *L0LoadBalancerLogic) GetLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error) {
	loadBalancer, err := l.Backend.GetLoadBalancer(loadBalancerID)
	if err != nil {
		return nil, err
	}

	if err := l.populateModel(loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (l *L0LoadBalancerLogic) DeleteLoadBalancer(loadBalancerID string) error {
	if err := l.Backend.DeleteLoadBalancer(loadBalancerID); err != nil {
		return err
	}

	if err := l.deleteEntityTags("load_balancer", loadBalancerID); err != nil {
		return err
	}

	return nil
}

func (l *L0LoadBalancerLogic) CreateLoadBalancer(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error) {
	if req.EnvironmentID == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentID not specified")
	}

	if req.LoadBalancerName == "" {
		return nil, errors.Newf(errors.MissingParameter, "LoadBalancerName not specified")
	}

	exists, err := l.doesLoadBalancerTagExist(req.EnvironmentID, req.LoadBalancerName)
	if err != nil {
		return nil, err
	}

	if exists {
		err := fmt.Errorf("LoadBalancer with name '%s' already exists in Environment '%s'", req.LoadBalancerName, req.EnvironmentID)
		return nil, errors.New(errors.InvalidLoadBalancerName, err)
	}

	loadBalancer, err := l.Backend.CreateLoadBalancer(
		req.LoadBalancerName,
		req.EnvironmentID,
		req.IsPublic,
		req.Ports,
		req.HealthCheck,
	)

	if err != nil {
		return loadBalancer, err
	}

	loadBalancerID := loadBalancer.LoadBalancerID
	if err := l.TagStore.Insert(models.Tag{EntityID: loadBalancerID, EntityType: "load_balancer", Key: "name", Value: req.LoadBalancerName}); err != nil {
		return loadBalancer, err
	}

	environmentID := loadBalancer.EnvironmentID
	if err := l.TagStore.Insert(models.Tag{EntityID: loadBalancerID, EntityType: "load_balancer", Key: "environment_id", Value: environmentID}); err != nil {
		return loadBalancer, err
	}

	if err := l.populateModel(loadBalancer); err != nil {
		return loadBalancer, err
	}

	return loadBalancer, nil
}

func (l *L0LoadBalancerLogic) UpdateLoadBalancerPorts(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error) {
	loadBalancer, err := l.Backend.UpdateLoadBalancerPorts(loadBalancerID, ports)
	if err != nil {
		return nil, err
	}

	if err := l.populateModel(loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (l *L0LoadBalancerLogic) UpdateLoadBalancerHealthCheck(loadBalancerID string, healthCheck models.HealthCheck) (*models.LoadBalancer, error) {
	loadBalancer, err := l.Backend.UpdateLoadBalancerHealthCheck(loadBalancerID, healthCheck)
	if err != nil {
		return nil, err
	}

	if err := l.populateModel(loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (l *L0LoadBalancerLogic) doesLoadBalancerTagExist(environmentID, name string) (bool, error) {
	tags, err := l.TagStore.SelectByType("load_balancer")
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

func (l *L0LoadBalancerLogic) populateModel(model *models.LoadBalancer) error {
	tags, err := l.TagStore.SelectByTypeAndID("load_balancer", model.LoadBalancerID)
	if err != nil {
		return err
	}

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		model.EnvironmentID = tag.Value
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		model.LoadBalancerName = tag.Value
	}

	if model.EnvironmentID != "" {
		tags, err := l.TagStore.SelectByTypeAndID("environment", model.EnvironmentID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			model.EnvironmentName = tag.Value
		}
	}

	tags, err = l.TagStore.SelectByType("service")
	if err != nil {
		return err
	}

	if tag, ok := tags.WithKey("load_balancer_id").WithValue(model.LoadBalancerID).First(); ok {
		model.ServiceID = tag.EntityID

		serviceTags, err := l.TagStore.SelectByTypeAndID("service", model.ServiceID)
		if err != nil {
			return err
		}

		if tag, ok := serviceTags.WithKey("name").First(); ok {
			model.ServiceName = tag.Value
		}
	}

	return nil
}
