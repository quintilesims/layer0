package system

import (
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

		s.Layer0.CreateTask(taskName, environmentID, deployID, copies, overrides)
	}

	// create 100 tasks
	createTask("TaskA", 50, "sleep 10")
	createTask("TaskB", 25, "sleep 20")
	createTask("TaskC", 15, "sleep 30")
	createTask("TaskD", 1, "sleep 00")
	createTask("TaskE", 1, "sleep 10")
	createTask("TaskF", 1, "sleep 20")
	createTask("TaskG", 1, "sleep 30")
	createTask("TaskH", 1, "sleep 40")
	createTask("TaskI", 1, "sleep 50")
	createTask("TaskJ", 1, "sleep 60")
	createTask("TaskK", 1, "sleep 70")
	createTask("TaskL", 1, "sleep 80")
	createTask("TaskM", 1, "sleep 90")

	testutils.WaitFor(t, time.Minute*5, func() bool {
		logrus.Printf("Waiting for tasks to be created")
		tasks := s.Layer0.ListTasks()
		// todo: ensure 100 tasks are in the specified environment
		return len(tasks) > 100
	})

	// run list, get on all
}
