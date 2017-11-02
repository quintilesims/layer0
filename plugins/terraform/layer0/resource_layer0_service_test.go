package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceServiceCreateRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	req := models.CreateServiceRequest{
		ServiceName:    "svc_name",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		DeployID:       "dpl_id",
		Scale:          3,
	}

	mockClient.EXPECT().
		CreateService(req).
		Return("job_id", nil)

	job := &models.Job{
		Status: job.Completed.String(),
		Result: "svc_id",
	}

	mockClient.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	service := &models.Service{
		ServiceID:      "svc_id",
		ServiceName:    "svc_name",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		DesiredCount:   3,
		Deployments: []models.Deployment{
			{DeployID: "dpl_id", Status: "PRIMARY"},
		},
	}

	mockClient.EXPECT().
		ReadService("svc_id").
		Return(service, nil)

	serviceResource := Provider().(*schema.Provider).ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":          "svc_name",
		"environment":   "env_id",
		"load_balancer": "lb_id",
		"deploy":        "dpl_id",
		"scale":         3,
	})

	if err := resourceLayer0ServiceCreate(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "svc_id", d.Id())
	assert.Equal(t, "svc_name", d.Get("name").(string))
	assert.Equal(t, "env_id", d.Get("environment").(string))
	assert.Equal(t, "lb_id", d.Get("load_balancer").(string))
	assert.Equal(t, "dpl_id", d.Get("deploy").(string))
	assert.Equal(t, 3, d.Get("scale").(int))
}

func TestResourceServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	mockClient.EXPECT().
		DeleteService("svc_id").
		Return("job_id", nil)

	mockClient.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{Status: job.Completed.String()}, nil)

	serviceResource := Provider().(*schema.Provider).ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{})
	d.SetId("svc_id")

	if err := resourceLayer0ServiceDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
