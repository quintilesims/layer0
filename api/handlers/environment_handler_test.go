package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
	"testing"
)

func TestListEnvironments(t *testing.T) {
	environments := []*models.EnvironmentSummary{
		{
			EnvironmentID: "some_id_1",
		},
		{
			EnvironmentID: "some_id_2",
		},
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name:    "Should return environments from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				envLogicMock.EXPECT().
					ListEnvironments().
					Return(environments, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.ListEnvironments(req, resp)

				var response []*models.EnvironmentSummary
				read(&response)

				reporter.AssertEqual(response, environments)
			},
		},
		HandlerTestCase{
			Name:    "Should propogate ListEnvironments error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				envLogicMock.EXPECT().
					ListEnvironments().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.ListEnvironments(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestGetEnvironment(t *testing.T) {
	environment := &models.Environment{
		EnvironmentID: "some_id",
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call GetEnvironment with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				envLogicMock.EXPECT().
					GetEnvironment("some_id").
					Return(environment, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.GetEnvironment(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should return environment from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				envLogicMock.EXPECT().
					GetEnvironment(gomock.Any()).
					Return(environment, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.GetEnvironment(req, resp)

				var response *models.Environment
				read(&response)

				reporter.AssertEqual(response, environment)
			},
		},
		HandlerTestCase{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.GetEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		HandlerTestCase{
			Name: "Should propagate GetEnvironment error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				envLogicMock.EXPECT().
					GetEnvironment(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.GetEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestDeleteEnvironment(t *testing.T) {
	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call CreateJob with correct params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				jobLogicMock.EXPECT().
					CreateJob(types.DeleteEnvironmentJob, "some_id").
					Return(&models.Job{}, nil)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.DeleteEnvironment(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should set Location and X-Jobid headers",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				jobLogicMock.EXPECT().
					CreateJob(gomock.Any(), gomock.Any()).
					Return(&models.Job{JobID: "job_id"}, nil)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.DeleteEnvironment(req, resp)

				header := resp.Header()
				reporter.AssertInSlice("/job/job_id", header["Location"])
				reporter.AssertInSlice("job_id", header["X-Jobid"])
			},
		},
		HandlerTestCase{
			Name: "Should propagate CreateJob error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				jobLogicMock.EXPECT().
					CreateJob(gomock.Any(), gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.DeleteEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
		HandlerTestCase{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				envLogicMock := mock_logic.NewMockEnvironmentLogic(ctrl)
				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewEnvironmentHandler(envLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.DeleteEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.MissingParameter), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestCreateEnvironment(t *testing.T) {
	request := models.CreateEnvironmentRequest{
		EnvironmentName: "env_name",
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call CanCreateEnvironment and CreateEnvironment with correct params",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockEnvironment := mock_logic.NewMockEnvironmentLogic(ctrl)
				mockEnvironment.EXPECT().
					CanCreateEnvironment(request).
					Return(true, nil)

				mockEnvironment.EXPECT().
					CreateEnvironment(request).
					Return(&models.Environment{}, nil)

				mockJob := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(mockEnvironment, mockJob)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.CreateEnvironment(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should return error if CanCreateEnvironment returns false",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockEnvironment := mock_logic.NewMockEnvironmentLogic(ctrl)
				mockEnvironment.EXPECT().
					CanCreateEnvironment(request).
					Return(false, nil)

				mockJob := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(mockEnvironment, mockJob)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.CreateEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.InvalidEnvironmentID), response.ErrorCode)
			},
		},
		HandlerTestCase{
			Name: "Should propagate CanCreateEnvironment error",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockEnvironment := mock_logic.NewMockEnvironmentLogic(ctrl)
				mockEnvironment.EXPECT().
					CanCreateEnvironment(gomock.Any()).
					Return(false, fmt.Errorf("some error"))

				mockJob := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(mockEnvironment, mockJob)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.CreateEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
		HandlerTestCase{
			Name: "Should propagate CreateEnvironment error",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockEnvironment := mock_logic.NewMockEnvironmentLogic(ctrl)
				mockEnvironment.EXPECT().
					CanCreateEnvironment(request).
					Return(true, nil)

				mockEnvironment.EXPECT().
					CreateEnvironment(gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockJob := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(mockEnvironment, mockJob)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.CreateEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestUpdateEnvironment(t *testing.T) {
	request := models.UpdateEnvironmentRequest{
		MinClusterCount: 2,
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call UpdateEnvironment with correct params",
			Request: &TestRequest{
				Body:       request,
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockEnvironment := mock_logic.NewMockEnvironmentLogic(ctrl)

				mockEnvironment.EXPECT().
					UpdateEnvironment("some_id", 2).
					Return(&models.Environment{}, nil)

				mockJob := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(mockEnvironment, mockJob)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.UpdateEnvironment(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should propagate UpdateEnvironment error",
			Request: &TestRequest{
				Body:       request,
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockEnvironment := mock_logic.NewMockEnvironmentLogic(ctrl)

				mockEnvironment.EXPECT().
					UpdateEnvironment(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("some error"))

				mockJob := mock_logic.NewMockJobLogic(ctrl)
				return NewEnvironmentHandler(mockEnvironment, mockJob)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*EnvironmentHandler)
				handler.UpdateEnvironment(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}
