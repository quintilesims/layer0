package system

import (
	"strconv"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/testutils"
)

func BenchmarkEnvironmentStress10(b *testing.B) { benchmarkEnvironmentStress(10, b) }
func BenchmarkEnvironmentStress50(b *testing.B) { benchmarkEnvironmentStress(50, b) }

// func BenchmarkEnvironmentStress100(b *testing.B) { benchmarkEnvironmentStress(100, b) }
// func BenchmarkEnvironmentStress250(b *testing.B) { benchmarkEnvironmentStress(250, b) }
// func BenchmarkEnvironmentStress500(b *testing.B) { benchmarkEnvironmentStress(500, b) }

func benchmarkEnvironmentStress(i int, b *testing.B) {
	tfvars := map[string]string{
		"num_environments": strconv.Itoa(i),
	}

	log.Debugf("Creating %v environments", i)

	s := NewStressTest(b, "cases/environment_stress", tfvars)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	log.Debugf("Benchmarking list operations for %v environments", i)

	for n := 0; n < b.N; n++ {
		s.Layer0.ListEnvironments()
	}
}

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
		go func(taskName string, copies int) {
			log.Debugf("Creating task %s (copies: %d)", taskName, copies)
			s.Layer0.CreateTask(taskName, environmentID, deployID, copies, nil)
		}(taskName, copies)
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

	// each task sleeps for 10 seconds
	// wait for all of them to complete
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
