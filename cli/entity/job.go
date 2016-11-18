package entity

import (
	"gitlab.imshealth.com/xfra/layer0/cli/printer/table"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/types"
	"strings"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

type Job models.Job

func NewJob(model *models.Job) *Job {
	job := Job(*model)
	return &job
}

func (this *Job) Table() table.Table {
	jobType := types.JobType(this.JobType).String()
	jobType = strings.Title(jobType)

	jobStatus := types.JobStatus(this.JobStatus).String()
	jobStatus = strings.Title(jobStatus)

	table := []table.Column{
		table.NewSingleRowColumn("JOB ID", this.JobID),
		table.NewSingleRowColumn("TASK ID", this.TaskID),
		table.NewSingleRowColumn("TYPE", jobType),
		table.NewSingleRowColumn("STATUS", jobStatus),
		table.NewSingleRowColumn("CREATED", this.TimeCreated.Format(TIME_FORMAT)),
	}

	return table
}
