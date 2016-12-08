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

func TestGetTask(t *testing.T) {
	entityIDTag := models.EntityTag{
		EntityID:   "tsk_id",
		EntityType: "task",
		Key:        "environment_id",
		Value:      "env_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.GetTask with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					GetTask("env_id", "tsk_id").
					Return(&models.Task{}, nil)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)
				logic.GetTask("tsk_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.GetTask error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					GetTask(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if _, err := logic.GetTask("tsk_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate model with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					GetTask(gomock.Any(), gomock.Any()).
					Return(&models.Task{TaskID: "tsk_id"}, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "environment_id",
					Value:      "env_id",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "deploy_id",
					Value:      "dpl_id",
				})

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				task, err := logic.GetTask("tsk_id")
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(task.TaskID, "tsk_id")
				reporter.AssertEqual(task.TaskName, "some_name")
				reporter.AssertEqual(task.EnvironmentID, "env_id")
				reporter.AssertEqual(task.DeployID, "dpl_id")
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
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if _, err := logic.GetTask("tsk_id"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListTasks(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.ListTasks",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListTasks().
					Return([]*models.Task{}, nil)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)
				logic.ListTasks()
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.ListTasks error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					ListTasks().
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if _, err := logic.ListTasks(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should populate models with correct tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				tasks := []*models.Task{
					&models.Task{TaskID: "some_id_1"},
					&models.Task{TaskID: "some_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListTasks().
					Return(tasks, nil)

				mockLogic.UseSQLite(t)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_1",
					EntityType: "task",
					Key:        "name",
					Value:      "some_name_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_1",
					EntityType: "task",
					Key:        "environment_id",
					Value:      "some_env_1",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_2",
					EntityType: "task",
					Key:        "name",
					Value:      "some_name_2",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "some_id_2",
					EntityType: "task",
					Key:        "environment_id",
					Value:      "some_env_2",
				})

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				tasks, err := logic.ListTasks()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(tasks[0].TaskID, "some_id_1")
				reporter.AssertEqual(tasks[0].TaskName, "some_name_1")
				reporter.AssertEqual(tasks[0].EnvironmentID, "some_env_1")

				reporter.AssertEqual(tasks[1].TaskID, "some_id_2")
				reporter.AssertEqual(tasks[1].TaskName, "some_name_2")
				reporter.AssertEqual(tasks[1].EnvironmentID, "some_env_2")
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				tasks := []*models.Task{
					&models.Task{TaskID: "some_id_1"},
					&models.Task{TaskID: "some_id_2"},
				}

				mockLogic.Backend.EXPECT().
					ListTasks().
					Return(tasks, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if _, err := logic.ListTasks(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteTask(t *testing.T) {
	entityIDTag := models.EntityTag{
		EntityID:   "tsk_id",
		EntityType: "task",
		Key:        "environment_id",
		Value:      "env_id",
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should call backend.DeleteTask with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					DeleteTask("env_id", "tsk_id").
					Return(nil)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)
				logic.DeleteTask("tsk_id")
			},
		},
		testutils.TestCase{
			Name: "Should delete task tags",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Return(nil)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "name",
					Value:      "some_name",
				})

				addTag(t, mockLogic.SQLite, models.EntityTag{
					EntityID:   "not_tsk_id",
					EntityType: "task",
					Key:        "name",
					Value:      "some_name",
				})

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return map[string]interface{}{
					"target": NewL0TaskLogic(mockLogic.Logic(), mockDeploy),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0TaskLogic)
				logic.DeleteTask("tsk_id")

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				reporter.AssertEqual(1, len(tags))
				reporter.AssertEqual(tags[0].EntityID, "not_tsk_id")
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.DeleteTask error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.UseSQLite(t)
				addTag(t, mockLogic.SQLite, entityIDTag)

				mockLogic.Backend.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if err := logic.DeleteTask("tsk_id"); err == nil {
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
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if err := logic.DeleteTask(""); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateTask(t *testing.T) {
	request := models.CreateTaskRequest{
		EnvironmentID: "env_id",
		TaskName:      "tsk_name",
		DeployID:      "dpl_id",
		Copies:        1,
	}

	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should error if request.EnvironmentID is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				req := request
				req.EnvironmentID = ""
				if _, err := logic.CreateTask(req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should error if request.TaskName is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				req := request
				req.TaskName = ""
				if _, err := logic.CreateTask(req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should error if request.DeployID is empty",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				req := request
				req.DeployID = ""
				if _, err := logic.CreateTask(req); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should call backend.CreateTask with correct params",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockLogic.Backend.EXPECT().
					// ignore container overrides for now
					CreateTask("env_id", "tsk_name", "dpl_id", 1, gomock.Any()).
					Return(&models.Task{}, nil)

				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)
				logic.CreateTask(request)
			},
		},
		testutils.TestCase{
			Name: "Should propagate backend.CreateTask error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)
				mockLogic.StubTagMock()

				mockLogic.Backend.EXPECT().
					CreateTask(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if _, err := logic.CreateTask(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
		testutils.TestCase{
			Name: "Should add correct name tag in database",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				task := &models.Task{
					EnvironmentID: "env_id",
					TaskID:        "tsk_id",
					DeployID:      "dpl_id",
				}

				mockLogic.Backend.EXPECT().
					CreateTask(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(task, nil)

				mockLogic.UseSQLite(t)

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return map[string]interface{}{
					"target": NewL0TaskLogic(mockLogic.Logic(), mockDeploy),
					"sqlite": mockLogic.SQLite,
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				testMap := target.(map[string]interface{})

				logic := testMap["target"].(*L0TaskLogic)
				if _, err := logic.CreateTask(request); err != nil {
					reporter.Error(err)
				}

				sqlite := testMap["sqlite"].(*data.TagDataStoreSQLite)
				tags, err := sqlite.Select()
				if err != nil {
					reporter.Error(err)
				}

				nameTag := models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "name",
					Value:      "tsk_name",
				}

				environmentTag := models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "environment_id",
					Value:      "env_id",
				}

				deployTag := models.EntityTag{
					EntityID:   "tsk_id",
					EntityType: "task",
					Key:        "deploy_id",
					Value:      "dpl_id",
				}

				reporter.AssertInSlice(nameTag, tags)
				reporter.AssertInSlice(environmentTag, tags)
				reporter.AssertInSlice(deployTag, tags)
			},
		},
		testutils.TestCase{
			Name: "Should propagate tag data error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockLogic := NewMockLogic(ctrl)

				mockLogic.Backend.EXPECT().
					CreateTask(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&models.Task{}, nil)

				mockLogic.Tag.EXPECT().
					GetTags(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)
				return NewL0TaskLogic(mockLogic.Logic(), mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				logic := target.(*L0TaskLogic)

				if _, err := logic.CreateTask(request); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
