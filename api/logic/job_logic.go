package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

type JobLogic interface {
	ListJobs() ([]*models.Job, error)
	GetJob(string) (*models.Job, error)
	CreateJob(types.JobType, interface{}) (*models.Job, error)
	Delete(string) error
}

type L0JobLogic struct {
	Logic
	TaskLogic   TaskLogic
	DeployLogic DeployLogic
}

func NewL0JobLogic(logic Logic, taskLogic TaskLogic, deployLogic DeployLogic) *L0JobLogic {
	return &L0JobLogic{
		Logic:       logic,
		TaskLogic:   taskLogic,
		DeployLogic: deployLogic,
	}
}

func (this *L0JobLogic) ListJobs() ([]*models.Job, error) {
	jobs, err := this.JobStore.SelectAll()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (this *L0JobLogic) GetJob(jobID string) (*models.Job, error) {
	job, err := this.JobStore.SelectByID(jobID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (this *L0JobLogic) Delete(jobID string) error {
	job, err := this.GetJob(jobID)
	if err != nil {
		return err
	}

	if err := this.TaskLogic.DeleteTask(job.TaskID); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code != errors.InvalidTaskID {
			return err
		}
	}

	if err := this.JobStore.Delete(jobID); err != nil {
		return err
	}

	if err := this.deleteEntityTags("job", jobID); err != nil {
		return err
	}

	return nil
}

func (this *L0JobLogic) CreateJob(jobType types.JobType, request interface{}) (*models.Job, error) {
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// remove extra "" added by json marshalling
	reqStr := string(bytes)
	reqStr = strings.TrimPrefix(reqStr, "\"")
	reqStr = strings.TrimSuffix(reqStr, "\"")

	jobID := id.GenerateHashedEntityID(string(jobType.String()))

	deploy, err := this.createJobDeploy(jobID)
	if err != nil {
		return nil, err
	}

	taskID, err := this.createJobTask(jobID, deploy.DeployID)
	if err != nil {
		return nil, err
	}
	job := &models.Job{
		JobID:       jobID,
		TaskID:      taskID,
		JobStatus:   int64(types.Pending),
		JobType:     int64(jobType),
		Request:     reqStr,
		TimeCreated: time.Now(),
	}

	if err := this.JobStore.Insert(job); err != nil {
		return nil, err
	}

	if err := this.TagStore.Insert(models.Tag{EntityID: jobID, EntityType: "job", Key: "task_id", Value: taskID}); err != nil {
		return nil, err
	}

	if jobType == types.CreateTaskJob {
		req, ok := request.(models.CreateTaskRequest)
		if !ok {
			return nil, fmt.Errorf("Unexpected request type for 'CreateTask' job type!")
		}

		this.Logic.Scaler.ScheduleRun(req.EnvironmentID, time.Second*10)
	}

	return job, nil
}

func (this *L0JobLogic) createJobTask(jobID, deployID string) (string, error) {
	taskRequest := models.CreateTaskRequest{
		DeployID:      deployID,
		EnvironmentID: config.API_ENVIRONMENT_ID,
		TaskName:      jobID,
	}

	taskID, err := this.TaskLogic.CreateTask(taskRequest)
	if err != nil {
		return "", err
	}

	return taskID, nil
}

func (this *L0JobLogic) createJobDeploy(jobID string) (*models.Deploy, error) {
	tmpl, err := template.New("").Parse(jobDockerrun)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse template: %v", err)
	}

	context := struct {
		RunnerVersionTag string
		Variables        []struct{ Key, Val string }
	}{
		RunnerVersionTag: config.RunnerVersionTag(),
		Variables: []struct{ Key, Val string }{
			{
				Key: config.JOB_ID,
				Val: jobID,
			},
			{
				Key: config.AWS_DYNAMO_TAG_TABLE,
				Val: config.DynamoTagTableName(),
			},
			{
				Key: config.AWS_DYNAMO_JOB_TABLE,
				Val: config.DynamoJobTableName(),
			},
			{
				Key: config.AWS_ACCESS_KEY_ID,
				Val: config.AWSAccessKey(),
			},
			{
				Key: config.AWS_SECRET_ACCESS_KEY,
				Val: config.AWSSecretKey(),
			},
			{
				Key: config.PREFIX,
				Val: config.Prefix(),
			},
			{
				Key: config.AWS_REGION,
				Val: config.AWSRegion(),
			},
			{
				Key: config.AWS_VPC_ID,
				Val: config.AWSVPCID(),
			},
			{
				Key: config.AWS_PUBLIC_SUBNETS,
				Val: config.AWSPublicSubnets(),
			},
			{
				Key: config.AWS_PRIVATE_SUBNETS,
				Val: config.AWSPrivateSubnets(),
			},
			{
				Key: config.RUNNER_LOG_LEVEL,
				Val: config.RunnerLogLevel(),
			},
		},
	}

	var dockerrun bytes.Buffer
	if err := tmpl.Execute(&dockerrun, context); err != nil {
		return nil, fmt.Errorf("Failed to write template: %v", err)
	}

	deployRequest := models.CreateDeployRequest{
		DeployName: "job",
		Dockerrun:  dockerrun.Bytes(),
	}

	deploy, err := this.DeployLogic.CreateDeploy(deployRequest)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

var jobDockerrun string = `
{
    "AWSEBDockerrunVersion": 2,
    "containerDefinitions": [
        {
            "name": "l0-job",
            "image": "quintilesims/l0-runner:{{ .RunnerVersionTag }}",
            "essential": true,
            "memory": 64,
            "environment": [
		{{ range $i, $v := .Variables }}{{ if $i }}, {{ end }}
                {
                    "name":  "{{ .Key }}",
                    "value": "{{ .Val }}"
                }{{ end }}
            ]
        }
    ]
}
`
