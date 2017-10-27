package client

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

var (
	timeMultiplier time.Duration = 1
)

const (
	LogQueryParamTail  = "tail"
	LogQueryParamStart = "start"
	LogQueryParamEnd   = "end"
)

const (
	TagQueryParamEnvironmentID = "environment_id"
	TagQueryParamFuzz          = "fuzz"
	TagQueryParamID            = "id"
	TagQueryParamName          = "name"
	TagQueryParamType          = "type"
	TagQueryParamVersion       = "version"
)

func SetTimeMultiplier(v time.Duration) func() {
	timeMultiplier = v
	return func() { timeMultiplier = 1 }
}

func WaitForDeployment(client Client, serviceID string, timeout time.Duration) (*models.Service, error) {
	var consecutiveSuccesses int
	check := func(service *models.Service) bool {
		for _, deployment := range service.Deployments {
			if deployment.DesiredCount != deployment.RunningCount {
				consecutiveSuccesses = 0
				return false
			}
		}

		consecutiveSuccesses++
		return consecutiveSuccesses >= 3
	}

	sleep := newLinearBackoffSleeper(time.Second)
	for start := time.Now(); time.Since(start) < timeout; sleep() {
		service, err := client.ReadService(serviceID)
		if err != nil {
			return nil, err
		}

		if check(service) {
			return service, nil
		}
	}

	return nil, fmt.Errorf("Deployment of service '%s' has not completed after %v", serviceID, timeout)
}

func WaitForJob(client Client, jobID string, timeout time.Duration) (*models.Job, error) {
	sleep := newLinearBackoffSleeper(time.Second)
	for start := time.Now(); time.Since(start) < timeout; sleep() {
		j, err := client.ReadJob(jobID)
		if err != nil {
			return nil, err
		}

		switch job.Status(j.Status) {
		case job.Completed:
			return j, nil
		case job.Error:
			var se *errors.ServerError
			if err := json.Unmarshal([]byte(j.Error), &se); err != nil {
				log.Printf("[DEBUG] Failed to marshal job.Error into errors.ServerError: %v", err)
				return nil, fmt.Errorf("An error occurred during the job's execution: %s", j.Error)
			}

			return nil, se
		}
	}

	return nil, fmt.Errorf("Timeout: job has not completed after %v", timeout)
}

func newLinearBackoffSleeper(d time.Duration) func() {
	var i int
	return func() {
		i++
		time.Sleep(d * time.Duration(i) * timeMultiplier)
	}
}

func Contains(value string, values []string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
