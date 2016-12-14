package logic

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.GetEnvironment with correct param",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetEnvironment("env_id").
					Return(&models.Environment{}, nil)

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)
				logic.GetEnvironment("env_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.GetEnvironment error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					GetEnvironment(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.GetEnvironment("env_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetEnvironment(gomock.Any()).
					Return(&models.Environment{EnvironmentID: "env_id"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "env_id",
					EntityType: "environment",
					Key:        "name",
					Value:      "some_name",
				})

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				environment, err := logic.GetEnvironment("env_id")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(environment.EnvironmentID, "env_id")
				reporter.AssertEqual(environment.EnvironmentName, "some_name")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetEnvironment(gomock.Any()).
					Return(&models.Environment{}, nil)

				mockLogic.Tag.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.GetEnvironment("env_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListEnvironments(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.ListEnvironments",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListEnvironments().
					Return([]*models.Environment{}, nil)

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)
				logic.ListEnvironments()
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.ListEnvironments error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListEnvironments().
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.ListEnvironments(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate models with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				environments := []*models.Environment{
					&models.Environment{EnvironmentID: "env_id_1"},
					&models.Environment{EnvironmentID: "env_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListEnvironments().
					Return(environments, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "env_id_1",
					EntityType: "environment",
					Key:        "name",
					Value:      "some_name_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "env_id_2",
					EntityType: "environment",
					Key:        "name",
					Value:      "some_name_2",
				})

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				environments, err := logic.ListEnvironments()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(environments[0].EnvironmentID, "env_id_1")
				reporter.AssertEqual(environments[0].EnvironmentName, "some_name_1")
				reporter.AssertEqual(environments[1].EnvironmentID, "env_id_2")
				reporter.AssertEqual(environments[1].EnvironmentName, "some_name_2")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				environments := []*models.Environment{
					&models.Environment{EnvironmentID: "env_id_1"},
					&models.Environment{EnvironmentID: "env_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListEnvironments().
					Return(environments, nil)

				mockLogic.Tag.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.ListEnvironments(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteEnvironment(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.DeleteEnvironment with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					DeleteEnvironment("env_id").
					Return(nil)

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)
				logic.DeleteEnvironment("env_id")
			},
		},
		testutils.TestCase{
			Name: "Should delete environment tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteEnvironment(gomock.Any()).
					Return(nil)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "env_id",
					EntityType: "environment",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "not_env_id",
					EntityType: "environment",
					Key:        "name",
					Value:      "some_name",
				})

				return map[string]interface{}{
					"target": NewL0EnvironmentLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0EnvironmentLogic)
				logic.DeleteEnvironment("env_id")

				sqlite := testMap["sqlite"].(*tag_store.TagStoreStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(1, len(tags))
				reporter.AssertEqual(tags[0].EntityID, "not_env_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.DeleteEnvironment error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteEnvironment(gomock.Any()).
					Return(fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if err := logic.DeleteEnvironment(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteEnvironment(gomock.Any()).
					Return(nil)

				mockLogic.Tag.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if err := logic.DeleteEnvironment(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateEnvironment(t *testing.T) {
	request := models.CreateEnvironmentRequest{
		EnvironmentName:  "some_name",
		InstanceSize:     "some_size",
		MinClusterCount:  2,
		UserDataTemplate: []byte("user_data"),
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should error if request.EnvironmentName is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.CreateEnvironment(models.CreateEnvironmentRequest{}); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should call backend.CreateEnvironment with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateEnvironment("some_name", "some_size", 2, []byte("user_data")).
					Return(&models.Environment{}, nil)

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)
				logic.CreateEnvironment(request)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.CreateEnvironment error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateEnvironment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.CreateEnvironment(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should add correct name tag in database",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateEnvironment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Environment{EnvironmentID: "env_id"}, nil)

				mockLogic.UseSQLite(t)

				return map[string]interface{}{
					"target": NewL0EnvironmentLogic(mockLogic.Logic()),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0EnvironmentLogic)
				if _, err := logic.CreateEnvironment(request); err != nil {
					reporter.Error(err)
				}

				sqlite := testMap["sqlite"].(*tag_store.TagStoreStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(len(tags), 1)
				reporter.AssertEqual(tags[0].EntityID, "env_id")
				reporter.AssertEqual(tags[0].EntityType, "environment")
				reporter.AssertEqual(tags[0].Key, "name")
				reporter.AssertEqual(tags[0].Value, "some_name")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateEnvironment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Environment{}, nil)

				mockLogic.Tag.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.CreateEnvironment(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateEnvironment(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.UpdateEnvironment with correct param",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					UpdateEnvironment("env_id", 2).
					Return(&models.Environment{}, nil)

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)
				logic.UpdateEnvironment("env_id", 2)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.UpdateEnvironment error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					UpdateEnvironment(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.UpdateEnvironment("env_id", 2); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					UpdateEnvironment(gomock.Any(), gomock.Any()).
					Return(&models.Environment{EnvironmentID: "env_id"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "env_id",
					EntityType: "environment",
					Key:        "name",
					Value:      "some_name",
				})

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				environment, err := logic.UpdateEnvironment("env_id", 2)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(environment.EnvironmentID, "env_id")
				reporter.AssertEqual(environment.EnvironmentName, "some_name")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					UpdateEnvironment(gomock.Any(), gomock.Any()).
					Return(&models.Environment{}, nil)

				mockLogic.Tag.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				return NewL0EnvironmentLogic(mockLogic.Logic())
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0EnvironmentLogic)

				if _, err := logic.UpdateEnvironment("env_id", 2); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
