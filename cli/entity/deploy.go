package entity

import (
	"github.com/quintilesims/layer0/cli/printer/table"
	"github.com/quintilesims/layer0/common/models"
	"strconv"
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

// sorting functions
type ByVersion []*Deploy

func (d ByVersion) Len() int {
	return len(d)
}
func (d ByVersion) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
func (d ByVersion) Less(i, j int) bool {
	vLeft, errLeft := strconv.Atoi(d[i].Version)
	vRight, errRight := strconv.Atoi(d[j].Version)

	// these are not numbers
	if errLeft != nil || errRight != nil {
		return d[i].Version > d[j].Version
	}

	// sort descending
	return vLeft > vRight
}
