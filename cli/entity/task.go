package entity

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/cli/printer/table"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type Task models.Task

func NewTask(model *models.Task) *Task {
	task := Task(*model)
	return &task
}

func (this *Task) Table() table.Table {
	environment := this.EnvironmentName
	if environment == "" {
		environment = this.EnvironmentID
	}

	scale := fmt.Sprintf("%d/%d", this.RunningCount, this.DesiredCount)
	if this.PendingCount != 0 {
		scale = fmt.Sprintf("%s (%d)", scale, this.PendingCount)
	}

	deploy := fmt.Sprintf("%s:%s", this.DeployName, this.DeployVersion)

	table := []table.Column{
		table.NewSingleRowColumn("TASK ID", this.TaskID),
		table.NewSingleRowColumn("TASK NAME", this.TaskName),
		table.NewSingleRowColumn("ENVIRONMENT", environment),
		table.NewSingleRowColumn("DEPLOY", deploy),
		table.NewSingleRowColumn("SCALE", scale),
	}

	return table
}
