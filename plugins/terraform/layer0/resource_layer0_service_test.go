package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceServiceCreateRead_stateless(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	defer client.SetTimeMultiplier(0)()

	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		Scale:          3,
		ServiceName:    "svc_name",
		Stateful:       false,
	}

	mockClient.EXPECT().
		CreateService(req).
		Return("svc_id", nil)

	service := &models.Service{
		Deployments: []models.Deployment{
			{DeployID: "dpl_id", Status: "PRIMARY", DesiredCount: 1, RunningCount: 1},
		},
		DesiredCount:   3,
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceID:      "svc_id",
		ServiceName:    "svc_name",
		Stateful:       false,
	}

	mockClient.EXPECT().
		ReadService("svc_id").
		Return(service, nil).
		AnyTimes()

	serviceResource := Provider().(*schema.Provider).ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":          "svc_name",
		"environment":   "env_id",
		"load_balancer": "lb_id",
		"deploy":        "dpl_id",
		"scale":         3,
		"stateful":      false,
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
	assert.Equal(t, false, d.Get("stateful").(bool))
}

func TestResourceServiceCreateRead_stateful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	defer client.SetTimeMultiplier(0)()

	req := models.CreateServiceRequest{
		DeployID:       "dpl_id",
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		Scale:          3,
		ServiceName:    "svc_name",
		Stateful:       true,
	}

	mockClient.EXPECT().
		CreateService(req).
		Return("svc_id", nil)

	service := &models.Service{
		Deployments: []models.Deployment{
			{DeployID: "dpl_id", Status: "PRIMARY", DesiredCount: 1, RunningCount: 1},
		},
		DesiredCount:   3,
		EnvironmentID:  "env_id",
		LoadBalancerID: "lb_id",
		ServiceID:      "svc_id",
		ServiceName:    "svc_name",
		Stateful:       true,
	}

	mockClient.EXPECT().
		ReadService("svc_id").
		Return(service, nil).
		AnyTimes()

	serviceResource := Provider().(*schema.Provider).ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{
		"name":          "svc_name",
		"environment":   "env_id",
		"load_balancer": "lb_id",
		"deploy":        "dpl_id",
		"scale":         3,
		"stateful":      true,
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
	assert.Equal(t, true, d.Get("stateful").(bool))
}

func TestResourceServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	mockClient.EXPECT().
		DeleteService("svc_id").
		Return(nil)

	serviceResource := Provider().(*schema.Provider).ResourcesMap["layer0_service"]
	d := schema.TestResourceDataRaw(t, serviceResource.Schema, map[string]interface{}{})
	d.SetId("svc_id")

	if err := resourceLayer0ServiceDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
