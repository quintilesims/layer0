package system

import (
	"log"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/testutils"
)

// Test Resources:
// This test creates an environment named 'tp'
// and a deploy named 'alpine'
func TestTaskStress(t *testing.T) {
	t.Parallel()

	s := NewStressTest(t, "cases/task_stress", nil)
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
		go func(taskName string) {
			log.Printf("[DEBUG] Creating task %s (copies: %d)", taskName, copies)
			s.Layer0.CreateTask(taskName, environmentID, deployID, nil)
		}(taskName)
	}

	testutils.WaitFor(t, time.Second*30, time.Minute*10, func() bool {
		log.Printf("[DEBUG] Waiting for all tasks tun run")
		var numTasks int
		for _, taskSummary := range s.Layer0.ListTasks() {
			if taskSummary.EnvironmentID == environmentID {
				numTasks++
			}
		}

		log.Printf("[DEBUG] %d/100 tasks have run", numTasks)
		return numTasks >= 100
	})

	// each task sleeps for 10 seconds
	// wait for all of them to complete
	time.Sleep(time.Second * 10)

	log.Printf("[DEBUG] Checking task exit codes")
	for _, taskSummary := range s.Layer0.ListTasks() {
		if taskSummary.EnvironmentID == environmentID {
			task := s.Layer0.ReadTask(taskSummary.TaskID)
			container := task.Containers[0]

			if container.ExitCode != 0 {
				t.Fatalf("Task %s has unexpected exit code: %#v", task.TaskID, container)
			}
		}
	}
}
