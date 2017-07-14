package logic

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestServicePopulateModel(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc_1"},
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb_1"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},

		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc_2"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "e2"},
		{EntityID: "s2", EntityType: "service", Key: "load_balancer_id", Value: "l2"},
	})

	cases := map[*models.Service]func(*models.Service){
		{
			ServiceID:      "s1",
			EnvironmentID:  "e1",
			LoadBalancerID: "l1",
			Deployments: []models.Deployment{
				{DeployID: "d1"},
			},
		}: func(m *models.Service) {
			testutils.AssertEqual(t, m.ServiceName, "svc_1")
			testutils.AssertEqual(t, m.EnvironmentName, "env_1")
			testutils.AssertEqual(t, m.LoadBalancerName, "lb_1")
			testutils.AssertEqual(t, m.Deployments[0].DeployName, "dpl_1")
		},
		{
			ServiceID: "s2",
		}: func(m *models.Service) {
			testutils.AssertEqual(t, m.ServiceName, "svc_2")
			testutils.AssertEqual(t, m.EnvironmentID, "e2")
			testutils.AssertEqual(t, m.LoadBalancerID, "l2")
		},
	}

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	for model, fn := range cases {
		if err := serviceLogic.populateModel(model); err != nil {
			t.Fatal(err)
		}

		fn(model)
	}
}

func TestGetService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		GetService("e1", "s1").
		Return(&models.Service{ServiceID: "s1"}, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	service, err := serviceLogic.GetService("s1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "s1")
	testutils.AssertEqual(t, service.EnvironmentID, "e1")
}

func TestListServices(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		ListServices().
		Return([]*models.Service{
			{ServiceID: "s1"},
			{ServiceID: "s2"},
		}, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "e2"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	services, err := serviceLogic.ListServices()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(services), 2)
	testutils.AssertEqual(t, services[0].ServiceID, "s1")
	testutils.AssertEqual(t, services[0].EnvironmentID, "e1")
	testutils.AssertEqual(t, services[1].ServiceID, "s2")
	testutils.AssertEqual(t, services[1].EnvironmentID, "e2")
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

	tags, err := testLogic.TagStore.SelectByType("service")
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCreateService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		CreateService("name", "e1", "d1", "l1").
		Return(&models.Service{ServiceID: "s1"}, nil)

	testLogic.Scaler.EXPECT().
		ScheduleRun("e1", gomock.Any())

	request := models.CreateServiceRequest{
		ServiceName:    "name",
		EnvironmentID:  "e1",
		DeployID:       "d1",
		LoadBalancerID: "l1",
	}

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	service, err := serviceLogic.CreateService(request)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "s1")
	testutils.AssertEqual(t, service.EnvironmentID, "e1")

	testLogic.AssertTagExists(t, models.Tag{EntityID: "s1", EntityType: "service", Key: "name", Value: "name"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"})
}

func TestCreateServiceError_missingRequiredParams(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())

	cases := map[string]models.CreateServiceRequest{
		"Missing EnvironmentID": {
			ServiceName: "name",
			DeployID:    "d1",
		},
		"Missing ServiceName": {
			EnvironmentID: "e1",
			DeployID:      "d1",
		},
		"Missing DeployID": {
			EnvironmentID: "e1",
			ServiceName:   "name",
		},
	}

	for name, request := range cases {
		if _, err := serviceLogic.CreateService(request); err == nil {
			t.Errorf("Case %s: error was nil!", name)
		}
	}
}

func TestCreateServiceError_duplicateName(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
	})

	request := models.CreateServiceRequest{
		EnvironmentID: "e1",
		ServiceName:   "svc",
		DeployID:      "d1",
	}

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	if _, err := serviceLogic.CreateService(request); err == nil {
		t.Errorf("Error was nil!")
	}
}

func TestUpdateService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		UpdateService("e1", "s1", "d1").
		Return(&models.Service{ServiceID: "s1"}, nil)

	testLogic.Scaler.EXPECT().
		ScheduleRun("e1", gomock.Any())

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
	})

	request := models.UpdateServiceRequest{
		DeployID: "d1",
	}

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	service, err := serviceLogic.UpdateService("s1", request)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "s1")
	testutils.AssertEqual(t, service.EnvironmentID, "e1")
}

func TestScaleService(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		ScaleService("e1", "s1", 2).
		Return(&models.Service{ServiceID: "s1"}, nil)

	testLogic.Scaler.EXPECT().
		ScheduleRun("e1", gomock.Any())

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	service, err := serviceLogic.ScaleService("s1", 2)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, service.ServiceID, "s1")
	testutils.AssertEqual(t, service.EnvironmentID, "e1")
}

func TestGetServiceLogs(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	logs := []*models.LogFile{
		{Name: "alpha", Lines: []string{"first", "second"}},
		{Name: "beta", Lines: []string{"first", "second", "third"}},
	}

	testLogic.Backend.EXPECT().
		GetServiceLogs("e1", "s1", 100).
		Return(logs, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
	})

	serviceLogic := NewL0ServiceLogic(testLogic.Logic())
	received, err := serviceLogic.GetServiceLogs("s1", 100)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, received, logs)
}
