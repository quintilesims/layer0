package logic

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/data"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetService(t *testing.T) {
	entityIDTag := models.EntityTag{
		EntityID:   "svc_id",
		EntityType: "service",
		Key:        "environment_id",
		Value:      "env_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.GetService with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					GetService("env_id", "svc_id").
					Return(&models.Service{}, nil)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)
				logic.GetService("svc_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.GetService error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					GetService(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.GetService("svc_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetService(gomock.Any(), gomock.Any()).
					Return(&models.Service{ServiceID: "svc_id"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "environment_id",
					Value:      "env_id",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "load_balancer_id",
					Value:      "lb_id",
				})

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				service, err := logic.GetService("svc_id")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(service.ServiceID, "svc_id")
				reporter.AssertEqual(service.ServiceName, "some_name")
				reporter.AssertEqual(service.EnvironmentID, "env_id")
				reporter.AssertEqual(service.LoadBalancerID, "lb_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.GetService("svc_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteService(t *testing.T) {
	entityIDTag := models.EntityTag{
		EntityID:   "svc_id",
		EntityType: "service",
		Key:        "environment_id",
		Value:      "env_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.DeleteService with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					DeleteService("env_id", "svc_id").
					Return(nil)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)
				logic.DeleteService("svc_id")
			},
		},
		testutils.TestCase{
			Name: "Should delete service tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteService(gomock.Any(), gomock.Any()).
					Return(nil)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "not_svc_id",
					EntityType: "service",
					Key:        "name",
					Value:      "some_name",
				})

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return map[string]interface{}{
					"target": NewL0ServiceLogic(mockLogic.Logic(), mockDeploy),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0ServiceLogic)
				logic.DeleteService("svc_id")

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(1, len(tags))
				reporter.AssertEqual(tags[0].EntityID, "not_svc_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.DeleteService error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					DeleteService(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if err := logic.DeleteService("svc_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if err := logic.DeleteService(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateService(t *testing.T) {
	request := models.CreateServiceRequest{
		EnvironmentID:  "env_id",
		ServiceName:    "svc_name",
		LoadBalancerID: "lb_id",
		DeployID:       "dpl_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should error if request.EnvironmentID is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				req := request
				req.EnvironmentID = ""
				if _, err := logic.CreateService(req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should error if request.ServiceName is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				req := request
				req.ServiceName = ""
				if _, err := logic.CreateService(req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should error if name/environment_id tags already exist",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "svcid",
					EntityType: "service",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "svcid",
					EntityType: "service",
					Key:        "environment_id",
					Value:      "envid",
				})

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				request := models.CreateServiceRequest{
					EnvironmentID: "envid",
					ServiceName:   "some_name",
				}

				if _, err := logic.CreateService(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should call backend.CreateService with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateService("svc_name", "env_id", "dpl_id", "lb_id").
					Return(&models.Service{}, nil)

				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)
				logic.CreateService(request)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.CreateService error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateService(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.CreateService(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should add correct name tag in database",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				service := &models.Service{
					EnvironmentID:  "env_id",
					ServiceID:      "svc_id",
					LoadBalancerID: "lb_id",
				}

				mockLogic.Backend.EXPECT().
					CreateService(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(service, nil)

				mockLogic.UseSQLite(t)

				return map[string]interface{}{
					"target": NewL0ServiceLogic(mockLogic.Logic(), mockDeploy),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0ServiceLogic)
				if _, err := logic.CreateService(request); err != nil {
					reporter.Error(err)
				}

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				nameTag := models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "name",
					Value:      "svc_name",
				}

				environmentTag := models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "environment_id",
					Value:      "env_id",
				}

				loadBalancerTag := models.EntityTag{
					EntityID:   "svc_id",
					EntityType: "service",
					Key:        "load_balancer_id",
					Value:      "lb_id",
				}

				reporter.AssertInSlice(nameTag, tags)
				reporter.AssertInSlice(environmentTag, tags)
				reporter.AssertInSlice(loadBalancerTag, tags)
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.CreateService(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestScaleService(t *testing.T) {
	entityIDTag := models.EntityTag{
		EntityID:   "svc_id",
		EntityType: "service",
		Key:        "environment_id",
		Value:      "env_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.ScaleService with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					ScaleService("env_id", "svc_id", 2).
					Return(&models.Service{}, nil)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)
				logic.ScaleService("svc_id", 2)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.ScaleService error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					ScaleService(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.ScaleService("svc_id", 2); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.ScaleService("svc_id", 2); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateService(t *testing.T) {
	entityIDTag := models.EntityTag{
		EntityID:   "svc_id",
		EntityType: "service",
		Key:        "environment_id",
		Value:      "env_id",
	}

	req := models.UpdateServiceRequest{
		DeployID: "dpl_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.UpdateService with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					UpdateService("env_id", "svc_id", "dpl_id").
					Return(&models.Service{}, nil)

				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)
				logic.UpdateService("svc_id", req)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.UpdateService error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					UpdateService(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.UpdateService("svc_id", req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0ServiceLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0ServiceLogic)

				if _, err := logic.UpdateService("svc_id", req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
