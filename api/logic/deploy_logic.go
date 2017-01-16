package logic

import (
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

type DeployLogic interface {
	ListDeploys() ([]*models.DeploySummary, error)
	GetDeploy(deployID string) (*models.Deploy, error)
	DeleteDeploy(deployID string) error
	CreateDeploy(model models.CreateDeployRequest) (*models.Deploy, error)
}

type L0DeployLogic struct {
	Logic
}

func NewL0DeployLogic(lgc Logic) *L0DeployLogic {
	return &L0DeployLogic{lgc}
}

func (this *L0DeployLogic) ListDeploys() ([]*models.DeploySummary, error) {
	deploys, err := this.Backend.ListDeploys()
	if err != nil {
		return nil, err
	}

	summaries := make([]*models.DeploySummary, len(deploys))
	for i, deploy := range deploys {
		if err := this.populateModel(deploy); err != nil {
			return nil, err
		}

		summaries[i] = &models.DeploySummary{
			DeployID:   deploy.DeployID,
			DeployName: deploy.DeployName,
			Version:    deploy.Version,
		}
	}

	return summaries, nil
}

func (this *L0DeployLogic) GetDeploy(deployID string) (*models.Deploy, error) {
	deploy, err := this.Backend.GetDeploy(deployID)
	if err != nil {
		return nil, err
	}

	if err := this.populateModel(deploy); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (this *L0DeployLogic) DeleteDeploy(deployID string) error {
	if err := this.Backend.DeleteDeploy(deployID); err != nil {
		return err
	}

	if err := this.deleteEntityTags("deploy", deployID); err != nil {
		return err
	}

	return nil
}

func (this *L0DeployLogic) CreateDeploy(req models.CreateDeployRequest) (*models.Deploy, error) {
	if req.DeployName == "" {
		return nil, errors.Newf(errors.MissingParameter, "DeployName is required")
	}

	deploy, err := this.Backend.CreateDeploy(req.DeployName, req.Dockerrun)
	if err != nil {
		return deploy, err
	}

	if err := this.upsertTagf(deploy.DeployID, "deploy", "name", req.DeployName); err != nil {
		return deploy, err
	}

	if err := this.upsertTagf(deploy.DeployID, "deploy", "version", deploy.Version); err != nil {
		return deploy, err
	}

	if err := this.populateModel(deploy); err != nil {
		return deploy, err
	}

	return deploy, nil
}

func (this *L0DeployLogic) populateModel(model *models.Deploy) error {
	tags, err := this.TagStore.SelectByQuery("deploy", model.DeployID)
	if err != nil {
		return err
	}

	if tag := tags.WithKey("name").First(); tag != nil {
		model.DeployName = tag.Value
	}

	if tag := tags.WithKey("version").First(); tag != nil {
		model.Version = tag.Value
	}

	return nil
}
