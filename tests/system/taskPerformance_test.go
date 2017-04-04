package system

import (
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
	"time"
)

// Test Resources:
// This test creates an environment named 'tp'
// and a deploy named 'alpine'
func TestTaskPerformance(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/task_performance", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	environmentID := s.Terraform.Output("environment_id")
	deployID := s.Terraform.Output("deploy_id")

	// create 100 tasks
	taskNameCopies := map[string]int{
		"TaskA": 50,
		"TaskB": 25,
		"TaskC": 10,
		"TaskD": 5,
		"TaskE": 5,
		"TaskF": 1,
		"TaskG": 1,
		"TaskH": 1,
		"TaskI": 1,
		"TaskJ": 1,
	}

	for taskName, copies := range taskNameCopies {
		log.Debugf("Creating task %s (copies: %d)", taskName, copies)
		s.Layer0.CreateTask(taskName, environmentID, deployID, copies, nil)
	}

	testutils.WaitFor(t, time.Second*30, time.Minute*10, func() bool {
		log.Debugf("Waiting for all tasks to run")

		var numTasks int
		for _, taskSummary := range s.Layer0.ListTasks() {
			if taskSummary.EnvironmentID == environmentID {
				numTasks++
			}
		}

		log.Debugf("%d/100 tasks have ran", numTasks)
		return numTasks >= 100
	})

	// sleep for 10 seconds since thats the longest a task will take
	time.Sleep(time.Second * 10)

	log.Debugf("Checking task exit codes")
	for _, taskSummary := range s.Layer0.ListTasks() {
		if taskSummary.EnvironmentID == environmentID {
			task := s.Layer0.GetTask(taskSummary.TaskID)
			detail := task.Copies[0].Details[0]

			if detail.ExitCode != 0 {
				t.Fatalf("Task %s has unexpected exit code: %#v", task.TaskID, detail)
			}
		}
	}
}
