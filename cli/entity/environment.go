package entity

import (
	"github.com/quintilesims/layer0/cli/printer/table"
	"github.com/quintilesims/layer0/common/models"
	"strconv"
)

type Environment models.Environment

func NewEnvironment(model *models.Environment) *Environment {
	environment := Environment(*model)
	return &environment
}

func (this *Environment) Table() table.Table {
	table := []table.Column{
		table.NewSingleRowColumn("ENVIRONMENT ID", this.EnvironmentID),
		table.NewSingleRowColumn("ENVIRONMENT NAME", this.EnvironmentName),
		table.NewSingleRowColumn("CLUSTER COUNT", strconv.Itoa(this.ClusterCount)),
		table.NewSingleRowColumn("INSTANCE SIZE", this.InstanceSize),
	}

	return table
}
