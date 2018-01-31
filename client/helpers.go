package client

import (
	"fmt"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

var timeMultiplier time.Duration = 1

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

func newLinearBackoffSleeper(d time.Duration) func() {
	var i int
	return func() {
		i++
		time.Sleep(d * time.Duration(i) * timeMultiplier)
	}
}
