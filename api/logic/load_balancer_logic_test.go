package logic

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
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

	expected := []*models.LoadBalancer{
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

	tags, err := testLogic.TagStore.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCreateLoadBalancer(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retLoadBalancer := &models.LoadBalancer{
		LoadBalancerID: "l1",
		EnvironmentID:  "e1",
		IsPublic:       true,
		Ports:          []models.Port{},
	}

	testLogic.Backend.EXPECT().
		CreateLoadBalancer("name", "e1", true, []models.Port{}).
		Return(retLoadBalancer, nil)

	request := models.CreateLoadBalancerRequest{
		LoadBalancerName: "name",
		EnvironmentID:    "e1",
		IsPublic:         true,
		Ports:            []models.Port{},
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
		"Missing EnvironmentID": models.CreateLoadBalancerRequest{
			LoadBalancerName: "name",
		},
		"Missing LoadBalancerName": models.CreateLoadBalancerRequest{
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

func TestUpdateLoadBalancer(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retLoadBalancer := &models.LoadBalancer{
		LoadBalancerID: "l1",
		EnvironmentID:  "e1",
		Ports:          []models.Port{},
	}

	testLogic.Backend.EXPECT().
		UpdateLoadBalancer("l1", []models.Port{}).
		Return(retLoadBalancer, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
		{EntityID: "extra", EntityType: "load_balancer", Key: "name", Value: "extra"},
	})

	loadBalancerLogic := NewL0LoadBalancerLogic(testLogic.Logic())
	received, err := loadBalancerLogic.UpdateLoadBalancer("l1", []models.Port{})
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
