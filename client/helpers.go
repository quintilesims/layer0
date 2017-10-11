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
