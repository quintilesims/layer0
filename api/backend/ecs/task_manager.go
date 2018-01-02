package ecsbackend

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

var ClusterCapacityReason = "Waiting for cluster capacity to run"

type ECSTaskManager struct {
	ECS            ecs.Provider
	CloudWatchLogs cloudwatchlogs.Provider
	Backend        backend.Backend
}

func NewECSTaskManager(
	ecsProvider ecs.Provider,
	cloudWatchLogsProvider cloudwatchlogs.Provider,
	backend backend.Backend,
) *ECSTaskManager {
	return &ECSTaskManager{
		ECS:            ecsProvider,
		CloudWatchLogs: cloudWatchLogsProvider,
		Backend:        backend,
	}
}

func (this *ECSTaskManager) ListTasks() ([]string, error) {
	clusterNames, err := this.Backend.ListEnvironments()
	if err != nil {
		return nil, err
	}

	taskARNs := []string{}
	for _, clusterName := range clusterNames {
		clusterTaskARNs, err := this.ECS.ListClusterTaskARNs(clusterName.String(), id.PREFIX)
		if err != nil {
			return nil, err
		}

		taskARNs = append(taskARNs, clusterTaskARNs...)
	}

	return taskARNs, nil
}

func (this *ECSTaskManager) GetTask(environmentID, taskARN string) (*models.Task, error) {
	clusterName := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	task, err := this.ECS.DescribeTask(clusterName.String(), taskARN)
	if err != nil {
		return nil, err
	}

	return modelFromTasks([]*ecs.Task{task})
}

func (this *ECSTaskManager) DeleteTask(environmentID, taskID string) error {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	taskARNs, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(ecsTaskID.String()))
	if err != nil {
		return err
	}

	// This stops the task, later reaping by AWS will prevent it from being returned.
	reason := "Task stopped by User"

	for _, taskARN := range taskARNs {
		if err := this.ECS.StopTask(ecsEnvironmentID.String(), reason, *taskARN); err != nil {
			return err
		}
	}

	return nil
}

func (this *ECSTaskManager) CreateTask(
	environmentID string,
	taskName string,
	deployID string,
	overrides []models.ContainerOverride,
) (*models.Task, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

	//taskID := id.GenerateHashedEntityID(taskName)
	//ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	ecsOverrides := []*ecs.ContainerOverride{}
	for _, override := range overrides {
		o := ecs.NewContainerOverride(override.ContainerName, override.EnvironmentOverrides)
		ecsOverrides = append(ecsOverrides, o)
	}

	startedBy := id.PREFIX
	tasks, failed, err := this.ECS.RunTask(ecsEnvironmentID.String(), ecsDeployID.TaskDefinition(), 1, stringp(startedBy), ecsOverrides)
	if err != nil {
		return nil, err
	}

	if len(failed) > 0 {
		return nil, fmt.Errorf("ECS failed to start the task!")
	}

	return modelFromTasks(tasks)
}

func (this *ECSTaskManager) GetTaskLogs(environmentID, taskID, start, end string, tail int) ([]*models.LogFile, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	taskARNs, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(ecsTaskID.String()))
	if err != nil {
		return nil, err
	}

	return GetLogs(this.CloudWatchLogs, taskARNs, start, end, tail)
}

// Assumes the tasks are all of the same type
func modelFromTasks(tasks []*ecs.Task) (*models.Task, error) {
	if len(tasks) == 0 {
		return nil, errors.Newf(errors.TaskDoesNotExist, "The specified task does not exist")
	}

	var pendingCount, runningCount int64
	copies := []models.TaskCopy{}
	for _, task := range tasks {
		switch status := aws.StringValue(task.LastStatus); status {
		case "RUNNING":
			runningCount = runningCount + 1
		case "PENDING":
			pendingCount = pendingCount + 1
		}

		details := []models.TaskDetail{}
		for _, container := range task.Containers {
			detail := models.TaskDetail{
				ContainerName: aws.StringValue(container.Name),
				LastStatus:    aws.StringValue(container.LastStatus),
				Reason:        stringOrEmpty(container.Reason),
				ExitCode:      int64OrZero(container.ExitCode),
			}

			details = append(details, detail)
		}

		copy := models.TaskCopy{
			Details:    details,
			Reason:     stringOrEmpty(task.StoppedReason),
			TaskCopyID: stringOrEmpty(task.TaskArn),
		}

		copies = append(copies, copy)
	}

	model := &models.Task{
		RunningCount: runningCount,
		PendingCount: pendingCount,
		Copies:       copies,
	}

	return model, nil
}

func (this *ECSTaskManager) describeTasks(ecsEnvironmentID id.ECSEnvironmentID, taskARNs []*string) ([]*ecs.Task, error) {
	ret := []*ecs.Task{}
	for i := len(taskARNs); i > 0; i = len(taskARNs) {
		if i > MAX_TASK_IDS {
			i = MAX_TASK_IDS
		}

		output, err := this.ECS.DescribeTasks(ecsEnvironmentID.String(), taskARNs[:i])
		if err != nil {
			return nil, err
		}

		ret = append(ret, output...)
		taskARNs = taskARNs[i:]
	}

	return ret, nil
}
