package client

import (
	"fmt"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/waitutils"
	"time"
)

const REQUIRED_SUCCESS_WAIT_COUNT = 3

func (c *APIClient) CreateService(name, environmentID, deployID, serviceID string) (*models.Service, error) {
	req := models.CreateServiceRequest{
		ServiceName:    name,
		EnvironmentID:  environmentID,
		DeployID:       deployID,
		LoadBalancerID: serviceID,
	}

	var service *models.Service
	if err := c.Execute(c.Sling("service/").Post("").BodyJSON(req), &service); err != nil {
		return nil, err
	}

	return service, nil
}

func (c *APIClient) DeleteService(id string) (string, error) {
	jobID, err := c.ExecuteWithJob(c.Sling("service/").Delete(id))
	if err != nil {
		return "", err
	}

	return jobID, nil
}

func (c *APIClient) UpdateService(serviceID, deployID string) (*models.Service, error) {
	request := models.UpdateServiceRequest{
		DeployID: deployID,
	}

	var service *models.Service
	if err := c.Execute(c.Sling("service/").Put(serviceID+"/deploy").BodyJSON(request), &service); err != nil {
		return nil, err
	}

	return service, nil
}

func (c *APIClient) GetService(id string) (*models.Service, error) {
	var service *models.Service
	if err := c.Execute(c.Sling("service/").Get(id), &service); err != nil {
		return nil, err
	}

	return service, nil
}

func (c *APIClient) GetServiceLogs(id string, tail int) ([]*models.LogFile, error) {
	url := id + "/logs"
	if tail > 0 {
		url = fmt.Sprintf("%s?tail=%d", url, tail)
	}

	var logFiles []*models.LogFile
	if err := c.Execute(c.Sling("service/").Get(url), &logFiles); err != nil {
		return nil, err
	}

	return logFiles, nil
}

func (c *APIClient) ListServices() ([]*models.ServiceSummary, error) {
	var services []*models.ServiceSummary
	if err := c.Execute(c.Sling("service/").Get(""), &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (c *APIClient) ScaleService(id string, count int) (*models.Service, error) {
	request := models.ScaleServiceRequest{
		DesiredCount: int64(count),
	}

	var service *models.Service
	if err := c.Execute(c.Sling("service/").Put(id+"/scale").BodyJSON(request), &service); err != nil {
		return nil, err
	}

	return service, nil
}

func (c *APIClient) WaitForDeployment(serviceID string, timeout time.Duration) (*models.Service, error) {
	var successCount int

	waiter := waitutils.Waiter{
		Name:    "WaitForDeployment",
		Timeout: timeout,
		Delay:   time.Second * 5,
		Clock:   c.Clock,
		Check: func() (bool, error) {
			service, err := c.GetService(serviceID)
			if err != nil {
				return false, err
			}

			for _, deploy := range service.Deployments {
				if deploy.DesiredCount != deploy.RunningCount {
					return false, nil
				}
			}

			successCount++
			return successCount >= REQUIRED_SUCCESS_WAIT_COUNT, nil
		},
	}

	if err := waiter.Wait(); err != nil {
		return nil, err
	}

	return c.GetService(serviceID)
}
