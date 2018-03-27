package client

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

// 'YYYY-MM-DD HH:MM' time layout as described by https://golang.org/src/time/format.go
const TimeLayout = "2006-01-02 15:04"

var timeMultiplier time.Duration = 1

func ParseLoggingQuery(query url.Values) (int, time.Time, time.Time, error) {
	var tail int
	if v := query.Get("tail"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil {
			return 0, time.Time{}, time.Time{}, fmt.Errorf("Tail must be an integer")
		}

		tail = t
	}

	parseTime := func(v string) (time.Time, error) {
		if v == "" {
			return time.Time{}, nil
		}

		return time.Parse(TimeLayout, v)
	}

	start, err := parseTime(query.Get("start"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("Invalid time: start must be in format YYYY-MM-DD HH:MM")
	}

	end, err := parseTime(query.Get("end"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("Invalid time: end must be in format YYYY-MM-DD HH:MM")
	}

	return tail, start, end, nil
}

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
