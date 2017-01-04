package logic

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retService := &models.Service{ServiceID: "s1"}

	testLogic.Backend.EXPECT().
		GetService("e1", "s1").
		Return(retService, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
		{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"},
		{EntityID: "extra", EntityType: "service", Key: "name", Value: "extra"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	received, err := serviceLogic.GetService("s1")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Service{
		ServiceID:      "s1",
		ServiceName:    "svc",
		EnvironmentID:  "e1",
		LoadBalancerID: "l1",
	}

	testutils.AssertEqual(t, received, expected)
}

func TestListServices(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retServices := []*models.Service{
		{ServiceID: "s1"},
		{ServiceID: "s2"},
	}

	testLogic.Backend.EXPECT().
		ListServices().
		Return(retServices, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc_1"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
		{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"},
		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc_2"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "e2"},
		{EntityID: "s2", EntityType: "service", Key: "load_balancer_id", Value: "l2"},
		{EntityID: "extra", EntityType: "service", Key: "name", Value: "extra"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	received, err := serviceLogic.ListServices()
	if err != nil {
		t.Fatal(err)
	}

	expected := []*models.Service{
		{
			ServiceID:      "s1",
			ServiceName:    "svc_1",
			EnvironmentID:  "e1",
			LoadBalancerID: "l1",
		},
		{
			ServiceID:      "s2",
			ServiceName:    "svc_2",
			EnvironmentID:  "e2",
			LoadBalancerID: "l2",
		},
	}

	testutils.AssertEqual(t, len(received), 2)
	testutils.AssertEqual(t, received[0], expected[0])
	testutils.AssertEqual(t, received[1], expected[1])
}

func TestDeleteService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		DeleteService("e1", "s1").
		Return(nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
		{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"},
		{EntityID: "extra", EntityType: "service", Key: "name", Value: "extra"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	if err := serviceLogic.DeleteService("s1"); err != nil {
		t.Fatal(err)
	}

	tags, err := testLogic.TagStore.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCreateService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retService := &models.Service{
		ServiceID:      "s1",
		EnvironmentID:  "e1",
		LoadBalancerID: "l1",
	}

	testLogic.Backend.EXPECT().
		CreateService("name", "e1", "d1", "l1").
		Return(retService, nil)

	request := models.CreateServiceRequest{
		ServiceName:    "name",
		EnvironmentID:  "e1",
		DeployID:       "d1",
		LoadBalancerID: "l1",
	}

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	received, err := serviceLogic.CreateService(request)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Service{
		ServiceID:      "s1",
		ServiceName:    "name",
		EnvironmentID:  "e1",
		LoadBalancerID: "l1",
	}

	testutils.AssertEqual(t, received, expected)
	testLogic.AssertTagExists(t, models.Tag{EntityID: "s1", EntityType: "service", Key: "name", Value: "name"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"})
}

func TestUpdateService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retService := &models.Service{
		ServiceID:      "s1",
		EnvironmentID:  "e1",
		LoadBalancerID: "l1",
	}

	testLogic.Backend.EXPECT().
		UpdateService("e1", "s1", "d1").
		Return(retService, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
		{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"},
		{EntityID: "extra", EntityType: "service", Key: "name", Value: "extra"},
	})

	request := models.UpdateServiceRequest{
		DeployID:      "d1",
	}

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	received, err := serviceLogic.UpdateService("s1", request)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Service{
		ServiceID:      "s1",
		ServiceName:    "svc",
		EnvironmentID:  "e1",
		LoadBalancerID: "l1",
	}

	testutils.AssertEqual(t, received, expected)
}

// todo: scale service
// todo: get service logs
