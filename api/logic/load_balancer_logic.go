package logic

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type LoadBalancerLogic interface {
	ListLoadBalancers() ([]*models.LoadBalancer, error)
	GetLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error)
	DeleteLoadBalancer(loadBalancerID string) error
	CreateLoadBalancer(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error)
	UpdateLoadBalancer(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error)
}

type L0LoadBalancerLogic struct {
	Logic
}

func NewL0LoadBalancerLogic(logic Logic) *L0LoadBalancerLogic {
	return &L0LoadBalancerLogic{
		Logic: logic,
	}
}

func (this *L0LoadBalancerLogic) ListLoadBalancers() ([]*models.LoadBalancer, error) {
	loadBalancers, err := this.Backend.ListLoadBalancers()
	if err != nil {
		return nil, err
	}

	for _, loadBalancer := range loadBalancers {
		if err := this.populateModel(loadBalancer); err != nil {
			return nil, err
		}
	}

	return loadBalancers, nil
}

func (this *L0LoadBalancerLogic) GetLoadBalancer(loadBalancerID string) (*models.LoadBalancer, error) {
	loadBalancer, err := this.Backend.GetLoadBalancer(loadBalancerID)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (this *L0LoadBalancerLogic) DeleteLoadBalancer(loadBalancerID string) error {
	if err := this.Backend.DeleteLoadBalancer(loadBalancerID); err != nil {
		return err
	}

	if err := this.deleteEntityTags(loadBalancerID, "load_balancer"); err != nil {
		return err
	}

	return nil
}

func (this *L0LoadBalancerLogic) CreateLoadBalancer(req models.CreateLoadBalancerRequest) (*models.LoadBalancer, error) {
	if req.EnvironmentID == "" {
		return nil, errors.Newf(errors.MissingParameter, "EnvironmentID not specified")
	}

	if req.LoadBalancerName == "" {
		return nil, errors.Newf(errors.MissingParameter, "LoadBalancerName not specified")
	}

	exists, err := this.doesLoadBalancerTagExist(req.EnvironmentID, req.LoadBalancerName)
	if err != nil {
		return nil, err
	}

	if exists {
		err := fmt.Errorf("LoadBalancer with name '%s' already exists in Environment '%s'", req.LoadBalancerName, req.EnvironmentID)
		return nil, errors.New(errors.InvalidLoadBalancerName, err)
	}

	loadBalancer, err := this.Backend.CreateLoadBalancer(
		req.LoadBalancerName,
		req.EnvironmentID,
		req.IsPublic,
		req.Ports)
	if err != nil {
		return loadBalancer, err
	}

	loadBalancerID := loadBalancer.LoadBalancerID
	if err := this.upsertTagf(loadBalancerID, "load_balancer", "name", req.LoadBalancerName); err != nil {
		return loadBalancer, err
	}

	environmentID := loadBalancer.EnvironmentID
	if err := this.upsertTagf(loadBalancerID, "load_balancer", "environment_id", environmentID); err != nil {
		return loadBalancer, err
	}

	if err := this.populateModel(loadBalancer); err != nil {
		return loadBalancer, err
	}

	return loadBalancer, nil
}

func (this *L0LoadBalancerLogic) UpdateLoadBalancer(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error) {
	loadBalancer, err := this.Backend.UpdateLoadBalancer(loadBalancerID, ports)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(loadBalancer); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (this *L0LoadBalancerLogic) doesLoadBalancerTagExist(environmentID, name string) (bool, error) {
	filter := map[string]string{
		"type":           "load_balancer",
		"environment_id": environmentID,
		"name":           name,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return false, err
	}

	return len(tags) > 0, nil
}

func (this *L0LoadBalancerLogic) populateModel(model *models.LoadBalancer) error {
	filter := map[string]string{
		"type": "load_balancer",
		"id":   model.LoadBalancerID,
	}

	tags, err := this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	for _, tag := range rangeTags(tags) {
		switch tag.Key {
		case "environment_id":
			model.EnvironmentID = tag.Value
		case "name":
			model.LoadBalancerName = tag.Value
		}
	}

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

	filter = map[string]string{
		"type":             "service",
		"load_balancer_id": model.LoadBalancerID,
	}

	tags, err = this.TagData.GetTags(filter)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		model.Services = append(model.Services, tag.EntityID)
	}

	return nil
}
