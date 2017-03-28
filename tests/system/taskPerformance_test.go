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
	t.Skip("Task Performance updates are WIP")
	t.Parallel()

	s := NewSystemTest(t, "cases/task_performance", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	environmentID := s.Terraform.Output("environment_id")
	deployID := s.Terraform.Output("deploy_id")

	// create 100 tasks
	taskNameCopies := map[string]int{
		"TaskA": 10,
		"TaskB": 10,
		"TaskC": 10,
		"TaskD": 10,
		"TaskE": 10,
		"TaskF": 10,
		"TaskG": 10,
		"TaskH": 10,
		"TaskI": 10,
		"TaskJ": 5,
		"TaskK": 1,
		"TaskL": 1,
		"TaskM": 1,
		"TaskN": 1,
		"TaskO": 1,
	}

	for taskName, copies := range taskNameCopies {
		log.Debugf("Creating task %s (copies: %d)", taskName, copies)
		s.Layer0.CreateTask(taskName, environmentID, deployID, copies, nil)
	}

	testutils.WaitFor(t, time.Second*30, time.Minute*10, func() bool {
		log.Debugf("Waiting for tasks to run")
		currentTaskNameCopies := map[string]int{}
		for _, taskSummary := range s.Layer0.ListTasks() {
			if taskSummary.EnvironmentID == environmentID {
				task := s.Layer0.GetTask(taskSummary.TaskID)
				currentTaskNameCopies[task.TaskName] = int(task.DesiredCount)
			}
		}

		ok := true
		for taskName, expectedCopies := range taskNameCopies {
			currentCopies := currentTaskNameCopies[taskName]
			log.Debugf("Task '%s' has %d/%d copies", taskName, currentCopies, expectedCopies)
			if currentCopies != expectedCopies {
				ok = false
			}
		}

		return ok
	})
}
