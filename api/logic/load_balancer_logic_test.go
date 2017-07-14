package logic

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestGetLoadBalancer(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retLoadBalancer := &models.LoadBalancer{LoadBalancerID: "l1"}

	testLogic.Backend.EXPECT().
		GetLoadBalancer("l1").
		Return(retLoadBalancer, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
		{EntityID: "extra", EntityType: "load_balancer", Key: "name", Value: "extra"},
	})

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	received, err := loadBalancerLogic.GetLoadBalancer("l1")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.LoadBalancer{
		LoadBalancerID:   "l1",
		LoadBalancerName: "lb",
		EnvironmentID:    "e1",
	}

	testutils.AssertEqual(t, received, expected)
}

func TestListLoadBalancers(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retLoadBalancers := []*models.LoadBalancer{
		{LoadBalancerID: "l1"},
		{LoadBalancerID: "l2"},
	}

	testLogic.Backend.EXPECT().
		ListLoadBalancers().
		Return(retLoadBalancers, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb_1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
		{EntityID: "l2", EntityType: "load_balancer", Key: "name", Value: "lb_2"},
		{EntityID: "l2", EntityType: "load_balancer", Key: "environment_id", Value: "e2"},
		{EntityID: "extra", EntityType: "load_balancer", Key: "name", Value: "extra"},
	})

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	received, err := loadBalancerLogic.ListLoadBalancers()
	if err != nil {
		t.Fatal(err)
	}

	expected := []*models.LoadBalancerSummary{
		{
			LoadBalancerID:   "l1",
			LoadBalancerName: "lb_1",
			EnvironmentID:    "e1",
		},
		{
			LoadBalancerID:   "l2",
			LoadBalancerName: "lb_2",
			EnvironmentID:    "e2",
		},
	}

	testutils.AssertEqual(t, received, expected)
}

func TestDeleteLoadBalancer(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		DeleteLoadBalancer("l1").
		Return(nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
		{EntityID: "extra", EntityType: "load_balancer", Key: "name", Value: "extra"},
	})

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	if err := loadBalancerLogic.DeleteLoadBalancer("l1"); err != nil {
		t.Fatal(err)
	}

	tags, err := testLogic.TagStore.SelectByType("load_balancer")
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCreateLoadBalancer(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	healthCheck := models.HealthCheck{
		Target:             "TCP:80:",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	retLoadBalancer := &models.LoadBalancer{
		LoadBalancerID: "l1",
		EnvironmentID:  "e1",
		IsPublic:       true,
		Ports:          []models.Port{},
		HealthCheck:    healthCheck,
	}

	testLogic.Backend.EXPECT().
		CreateLoadBalancer("name", "e1", true, []models.Port{}, healthCheck).
		Return(retLoadBalancer, nil)

	request := models.CreateLoadBalancerRequest{
		LoadBalancerName: "name",
		EnvironmentID:    "e1",
		IsPublic:         true,
		Ports:            []models.Port{},
		HealthCheck:      healthCheck,
	}

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	received, err := loadBalancerLogic.CreateLoadBalancer(request)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.LoadBalancer{
		LoadBalancerID:   "l1",
		LoadBalancerName: "name",
		EnvironmentID:    "e1",
		IsPublic:         true,
		Ports:            []models.Port{},
		HealthCheck:      healthCheck,
	}

	testutils.AssertEqual(t, received, expected)
	testLogic.AssertTagExists(t, models.Tag{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "name"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"})
}

func TestCreateLoadBalancerError_missingRequiredParams(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())

	cases := map[string]models.CreateLoadBalancerRequest{
		"Missing EnvironmentID": {
			LoadBalancerName: "name",
		},
		"Missing LoadBalancerName": {
			EnvironmentID: "e1",
		},
	}

	for name, request := range cases {
		if _, err := loadBalancerLogic.CreateLoadBalancer(request); err == nil {
			t.Errorf("Case %s: error was nil!", name)
		}
	}
}

func TestCreateLoadBalancerError_duplicateName(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb_1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
	})

	request := models.CreateLoadBalancerRequest{
		EnvironmentID:    "e1",
		LoadBalancerName: "lb_1",
	}

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	if _, err := loadBalancerLogic.CreateLoadBalancer(request); err == nil {
		t.Errorf("Error was nil!")
	}
}

func TestUpdateLoadBalancerPorts(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retLoadBalancer := &models.LoadBalancer{
		LoadBalancerID: "l1",
		EnvironmentID:  "e1",
		Ports:          []models.Port{},
	}

	testLogic.Backend.EXPECT().
		UpdateLoadBalancerPorts("l1", []models.Port{}).
		Return(retLoadBalancer, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
		{EntityID: "extra", EntityType: "load_balancer", Key: "name", Value: "extra"},
	})

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	received, err := loadBalancerLogic.UpdateLoadBalancerPorts("l1", []models.Port{})
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.LoadBalancer{
		LoadBalancerID:   "l1",
		LoadBalancerName: "lb",
		EnvironmentID:    "e1",
		Ports:            []models.Port{},
	}

	testutils.AssertEqual(t, received, expected)
}

func TestUpdateLoadBalancerHealthCheck(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	healthCheck := models.HealthCheck{
		Target:             "TCP:80",
		Interval:           30,
		Timeout:            5,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	testLogic.Backend.EXPECT().
		UpdateLoadBalancerHealthCheck("lb_id", healthCheck).
		Return(&models.LoadBalancer{HealthCheck: healthCheck}, nil)

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	received, err := loadBalancerLogic.UpdateLoadBalancerHealthCheck("lb_id", healthCheck)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, received.HealthCheck, healthCheck)
}
