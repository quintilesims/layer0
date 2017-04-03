package job

import (
	"github.com/quintilesims/layer0/api/logic"
)

type JobContext struct {
	jobID             string
	request           string
	Logic             *logic.Logic
	LoadBalancerLogic logic.LoadBalancerLogic
	ServiceLogic      logic.ServiceLogic
	TaskLogic         logic.TaskLogic
	EnvironmentLogic  logic.EnvironmentLogic
}

func NewJobContext(jobID string, lgc *logic.Logic, request string) *JobContext {
	return &JobContext{
		jobID:             jobID,
		request:           request,
		Logic:             lgc,
		LoadBalancerLogic: logic.NewL0LoadBalancerLogic(*lgc),
		ServiceLogic:      logic.NewL0ServiceLogic(*lgc),
		TaskLogic:         logic.NewL0TaskLogic(*lgc),
		EnvironmentLogic:  logic.NewL0EnvironmentLogic(*lgc),
	}
}

func (j *JobContext) CreateCopyWithNewRequest(request string) *JobContext {
	return &JobContext{
		jobID:             j.jobID,
		request:           request,
		Logic:             j.Logic,
		LoadBalancerLogic: j.LoadBalancerLogic,
		ServiceLogic:      j.ServiceLogic,
		TaskLogic:         j.TaskLogic,
		EnvironmentLogic:  j.EnvironmentLogic,
	}
}

func (j *JobContext) SetJobMeta(meta map[string]string) error {
	return j.Logic.JobStore.SetJobMeta(j.jobID, meta)
}

func (j *JobContext) AddJobMeta(key, val string) error {
	job, err := j.Logic.JobStore.SelectByID(j.jobID)
	if err != nil {
		return err
	}

	job.Meta[key] = val
	return j.SetJobMeta(job.Meta)
}

func (j *JobContext) Request() string {
	return j.request
}
