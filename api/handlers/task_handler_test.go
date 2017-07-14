package handlers

import (
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
)

func TestListTasks(t *testing.T) {
	tasks := []*models.TaskSummary{
		{
			TaskID: "some_id_1",
		},
		{
			TaskID: "some_id_2",
		},
	}

	testCases := []HandlerTestCase{
		{
			Name:    "Should return tasks from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					ListTasks().
					Return(tasks, nil)

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.ListTasks(req, resp)

				var response []*models.TaskSummary
				read(&response)

				reporter.AssertEqual(response, tasks)
			},
		},
		{
			Name:    "Should propagate ListTasks error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					ListTasks().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.ListTasks(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestGetTask(t *testing.T) {
	task := &models.Task{
		TaskID: "some_id",
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call GetTask with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					GetTask("some_id").
					Return(task, nil)

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.GetTask(req, resp)
			},
		},
		{
			Name: "Should return task from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					GetTask(gomock.Any()).
					Return(task, nil)

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.GetTask(req, resp)

				var response *models.Task
				read(&response)

				reporter.AssertEqual(response, task)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.GetTask(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate GetTask error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					GetTask(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.GetTask(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestDeleteTask(t *testing.T) {
	testCases := []HandlerTestCase{
		{
			Name: "Should call DeleteTask with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					DeleteTask("some_id").
					Return(nil)

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.DeleteTask(req, resp)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.DeleteTask(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate DeleteTask error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockTaskLogic(ctrl)
				logicMock.EXPECT().
					DeleteTask(gomock.Any()).
					Return(errors.Newf(errors.UnexpectedError, "some error"))

				return NewTaskHandler(logicMock, nil)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TaskHandler)
				handler.DeleteTask(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestCreateTask(t *testing.T) {
	request := models.CreateTaskRequest{
		TaskName:      "tsk_name",
		DeployID:      "dply_id",
		EnvironmentID: "env_id",
		Copies:        2,
	}

	testCase := HandlerTestCase{
		Name: "Should call CreateJob with correct params",
		Request: &TestRequest{
			Body: request,
		},
		Setup: func(ctrl *gomock.Controller) interface{} {
			jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
			jobLogicMock.EXPECT().
				CreateJob(types.CreateTaskJob, request).
				Return(&models.Job{}, nil)

			return NewTaskHandler(nil, jobLogicMock)
		},
		Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
			handler := target.(*TaskHandler)
			handler.CreateTask(req, resp)
		},
	}

	RunHandlerTestCase(t, testCase)

}
