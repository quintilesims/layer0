package entity

import (
	"github.com/quintilesims/layer0/cli/printer/table"
	"github.com/quintilesims/layer0/common/models"
)

type Deploy models.Deploy

func NewDeploy(model *models.Deploy) *Deploy {
	deploy := Deploy(*model)
	return &deploy
}

func (this *Deploy) Table() table.Table {
	table := []table.Column{
		table.NewSingleRowColumn("DEPLOY ID", this.DeployID),
		table.NewSingleRowColumn("DEPLOY NAME", this.DeployName),
		table.NewSingleRowColumn("VERSION", this.Version),
	}

	return table
}
