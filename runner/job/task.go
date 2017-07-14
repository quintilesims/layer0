package job

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/models"
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
	var createTaskRequest models.CreateTaskRequest
	if err := json.Unmarshal([]byte(context.Request()), &createTaskRequest); err != nil {
		return err
	}

	for i := 0; i < createTaskRequest.Copies; i++ {
		if err := runAndRetry(quit, time.Second*10, func() error {
			log.Infof("Running Action: CreateTask '%s', copy %d", createTaskRequest.TaskName, i)
			task, err := context.TaskLogic.CreateTask(createTaskRequest)
			if err != nil {
				log.Infof("Failed CreateTask '%s', copy %d", createTaskRequest.TaskName, i)
				return err
			}

			return runAndRetry(quit, time.Second*10, func() error {
				key := fmt.Sprintf("task_%d", i)
				return context.AddJobMeta(key, task.TaskID)
			})
		}); err != nil {
			return err
		}
	}

	return nil
}
