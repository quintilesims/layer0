package entity

import (
	"github.com/quintilesims/layer0/cli/printer/table"
	"github.com/quintilesims/layer0/common/models"
)

type DeploySummary models.DeploySummary

func NewDeploySummary(model *models.DeploySummary) *DeploySummary {
	summary := DeploySummary(*model)
	return &summary
}

func (this *DeploySummary) Table() table.Table {
	table := []table.Column{
		table.NewSingleRowColumn("DEPLOY ID", this.DeployID),
		table.NewSingleRowColumn("DEPLOY NAME", this.DeployName),
		table.NewSingleRowColumn("VERSION", this.Version),
	}

	return table
}
