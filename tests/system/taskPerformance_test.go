package system

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
	"time"
)

// Test Resources:
// This test creates an environment named 'tp'
// and a deploy named 'apline'
func TestTaskPerformance(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/task_performance", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	environmentID := s.Terraform.Output("environment_id")
	deployID := s.Terraform.Output("deploy_id")

	createTask := func(taskName string, copies int, command string) {
		overrides := []models.ContainerOverride{{
			ContainerName: "alpine",
			EnvironmentOverrides: map[string]string{
				"COMMAND": command,
			},
		}}

		logrus.Printf("Creating task %s (copies: %d)", taskName, copies)
		s.Layer0.CreateTask(taskName, environmentID, deployID, copies, overrides)
	}

	// create 100 tasks
	createTask("TaskA", 10, "sleep 10")
	createTask("TaskB", 10, "sleep 20")
	createTask("TaskC", 10, "sleep 30")
	createTask("TaskD", 10, "sleep 30")
	createTask("TaskE", 10, "sleep 50")
	createTask("TaskF", 10, "sleep 60")
	createTask("TaskG", 10, "sleep 70")
	createTask("TaskH", 10, "sleep 80")
	createTask("TaskI", 10, "sleep 90")
	createTask("TaskJ", 5, "sleep 95")
	createTask("TaskK", 1, "sleep 96")
	createTask("TaskL", 1, "sleep 97")
	createTask("TaskM", 1, "sleep 98")
	createTask("TaskN", 1, "sleep 99")
	createTask("Task0", 1, "sleep 00")

	forEachTaskDetail := func(do func(*models.Task, models.TaskDetail)) {
		for _, taskSummary := range s.Layer0.ListTasks() {
			if taskSummary.EnvironmentID == environmentID {
				task := s.Layer0.GetTask(taskSummary.TaskID)
				for _, copy := range task.Copies {
					for _, detail := range copy.Details {
						do(task, detail)
					}
				}
			}
		}
	}

	testutils.WaitFor(t, time.Second*30, time.Minute*10, func() bool {
		logrus.Printf("Waiting for tasks to finish running")

		count := 0
		forEachTaskDetail(func(task *models.Task, detail models.TaskDetail) {
			if detail.LastStatus == "STOPPED" {
				count++
			}
		})

		logrus.Printf("%d/100 tasks have finished running", count)
		return count >= 100
	})

	logrus.Printf("Checking task exit codes")
	forEachTaskDetail(func(task *models.Task, detail models.TaskDetail) {
		if detail.ExitCode != 0 {
			id := fmt.Sprintf("Task %s Container %s", task.TaskID, detail.ContainerName)
			t.Fatalf("%s has non-zero exit code %d! %s", id, detail.ExitCode, detail.Reason)
		}
	})
}
