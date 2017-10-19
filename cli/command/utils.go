package command

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
)

var (
	timeMultiplier time.Duration = 1
)

func extractArgs(received []string, names ...string) (map[string]string, error) {
	args := map[string]string{}
	for i, name := range names {
		if len(received)-1 < i {
			return nil, fmt.Errorf("Argument %s is required", name)
		}

		args[name] = received[i]
	}

	return args, nil
}

func buildLogQueryHelper(start, end string, tail int) url.Values {
	query := url.Values{}

	if tail > 0 {
		query.Set(client.LogQueryParamTail, strconv.Itoa(tail))
	}

	if start != "" {
		query.Set(client.LogQueryParamStart, start)
	}

	if end != "" {
		query.Set(client.LogQueryParamEnd, end)
	}

	return query
}

func WaitForDeployment(client client.Client, serviceID string, timeout time.Duration) (*models.Service, error) {
	successCount := 0
	requiredSuccessCount := 3

	check := func() (bool, error) {
		service, err := client.ReadService(serviceID)
		if err != nil {
			return false, err
		}

		for _, deployment := range service.Deployments {
			if deployment.DesiredCount != deployment.RunningCount {
				return false, nil
			}
		}

		successCount++
		return successCount >= requiredSuccessCount, nil
	}

	delay := 5 * time.Second * timeMultiplier

	for start := time.Now(); time.Since(start) < timeout; time.Sleep(delay) {
		finished, err := check()
		if err != nil {
			return nil, err
		}

		if finished {
			return client.ReadService(serviceID)
		}
	}

	return nil, fmt.Errorf("Deployment of service '%s' has not completed within the timeout '%v'", serviceID, timeout)
}
