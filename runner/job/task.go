package job

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend/ecs"
	"github.com/quintilesims/layer0/common/models"
	"time"
)

var DeleteTaskSteps = []Step{
	{
		Name:    "Delete Task",
		Timeout: time.Minute * 10,
		Action:  DeleteTask,
	},
}

var CreateTaskSteps = []Step{
	{
		Name:    "Create Task",
		Timeout: time.Hour * 24,
		Action:  CreateTask,
	},
}

func DeleteTask(quit chan bool, context *JobContext) error {
	taskID := context.Request()

	return runAndRetry(quit, time.Second*10, func() error {
		log.Infof("Running Action: DeleteTask on '%s'", taskID)
		return context.TaskLogic.DeleteTask(taskID)
	})
}

func CreateTask(quit chan bool, context *JobContext) error {
	var taskID string
	var createTaskRequest models.CreateTaskRequest
	if err := json.Unmarshal([]byte(context.Request()), &createTaskRequest); err != nil {
		return err
	}

	attempt := func() (*models.Task, error) {
		return context.TaskLogic.CreateTask(createTaskRequest)
	}

	if err := runAndRetry(quit, time.Second*10, func() error {
		log.Infof("Running Action: CreateTask on '%s'", createTaskRequest.TaskName)
		task, err := attempt()
		if err != nil {
			log.Printf("Failed to create task %s: %v\n", createTaskRequest.TaskName, err)
			if err, ok := err.(*ecsbackend.PartialCreateTaskFailure); ok {
				attempt = err.Retry
			}

			return err
		}

		taskID = task.TaskID
		return nil
	}); err != nil {
		return err
	}

	meta := map[string]string{"task_id": taskID}
	return context.SetJobMeta(meta)
}
