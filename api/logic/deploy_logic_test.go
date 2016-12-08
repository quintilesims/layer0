package logic

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/data"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetDeploy(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.GetDeploy with correct param",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetDeploy("dpl_id").
					Return(&models.Deploy{}, nil)

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)
				logic.GetDeploy("dpl_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.GetDeploy error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetDeploy(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if _, err := logic.GetDeploy("dpl_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetDeploy(gomock.Any()).
					Return(&models.Deploy{DeployID: "dpl_id"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "dpl_id",
					EntityType: "deploy",
					Key:        "name",
					Value:      "some_name",
				})

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				deploy, err := logic.GetDeploy("dpl_id")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(deploy.DeployID, "dpl_id")
				reporter.AssertEqual(deploy.DeployName, "some_name")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetDeploy(gomock.Any()).
					Return(&models.Deploy{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if _, err := logic.GetDeploy("dpl_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListDeploys(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.ListDeploys",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListDeploys().
					Return([]*models.Deploy{}, nil)

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)
				logic.ListDeploys()
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.ListDeploys error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListDeploys().
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if _, err := logic.ListDeploys(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate models with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				deploys := []*models.Deploy{
					&models.Deploy{DeployID: "dpl_id_1"},
					&models.Deploy{DeployID: "dpl_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListDeploys().
					Return(deploys, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "dpl_id_1",
					EntityType: "deploy",
					Key:        "name",
					Value:      "some_name_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "dpl_id_2",
					EntityType: "deploy",
					Key:        "name",
					Value:      "some_name_2",
				})

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				deploys, err := logic.ListDeploys()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(deploys[0].DeployID, "dpl_id_1")
				reporter.AssertEqual(deploys[0].DeployName, "some_name_1")
				reporter.AssertEqual(deploys[1].DeployID, "dpl_id_2")
				reporter.AssertEqual(deploys[1].DeployName, "some_name_2")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				deploys := []*models.Deploy{
					&models.Deploy{DeployID: "dpl_id_1"},
					&models.Deploy{DeployID: "dpl_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListDeploys().
					Return(deploys, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if _, err := logic.ListDeploys(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteDeploy(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.DeleteDeploy with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					DeleteDeploy("dpl_id").
					Return(nil)

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)
				logic.DeleteDeploy("dpl_id")
			},
		},
		testutils.TestCase{
			Name: "Should delete deploy tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteDeploy(gomock.Any()).
					Return(nil)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "dpl_id",
					EntityType: "deploy",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "not_dpl_id",
					EntityType: "deploy",
					Key:        "name",
					Value:      "some_name",
				})

				return map[string]interface{}{
					"target": NewL0DeployLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0DeployLogic)
				logic.DeleteDeploy("dpl_id")

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(1, len(tags))
				reporter.AssertEqual(tags[0].EntityID, "not_dpl_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.DeleteDeploy error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteDeploy(gomock.Any()).
					Return(fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if err := logic.DeleteDeploy(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteDeploy(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if err := logic.DeleteDeploy(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateDeploy(t *testing.T) {
	request := models.CreateDeployRequest{
		DeployName: "some_name",
		Dockerrun:  []byte{},
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should error if request.DeployName is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				req := models.CreateDeployRequest{}
				if _, err := logic.CreateDeploy(req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should call backend.CreateDeploy with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateDeploy("some_name", []byte{}).
					Return(&models.Deploy{}, nil)

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)
				logic.CreateDeploy(request)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.CreateDeploy error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateDeploy(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if _, err := logic.CreateDeploy(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should add correct name and version tag in database",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				model := &models.Deploy{
					DeployID: "dpl_id",
					Version:  "1",
				}

				mockLogic.Backend.EXPECT().
					CreateDeploy(gomock.Any(), gomock.Any()).
					Return(model, nil)

				mockLogic.UseSQLite(t)

				return map[string]interface{}{
					"target": NewL0DeployLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0DeployLogic)
				if _, err := logic.CreateDeploy(request); err != nil {
					reporter.Error(err)
				}

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				nameTag := models.EntityTag{
					EntityID:   "dpl_id",
					EntityType: "deploy",
					Key:        "name",
					Value:      "some_name",
				}

				versionTag := models.EntityTag{
					EntityID:   "dpl_id",
					EntityType: "deploy",
					Key:        "version",
					Value:      "1",
				}

				reporter.AssertEqual(len(tags), 2)
				reporter.AssertInSlice(nameTag, tags)
				reporter.AssertInSlice(versionTag, tags)
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateDeploy(gomock.Any(), gomock.Any()).
					Return(&models.Deploy{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0DeployLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0DeployLogic)

				if _, err := logic.CreateDeploy(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
