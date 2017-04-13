package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestListDeploys(t *testing.T) {
	deploys := []*models.DeploySummary{
		{DeployID: "d1"},
		{DeployID: "d2"},
	}

	testCases := []HandlerTestCase{
		{
			Name:    "Should return deploys from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					ListDeploys().
					Return(deploys, nil)

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.ListDeploys(req, resp)

				var response []*models.DeploySummary
				read(&response)

				reporter.AssertEqual(response, deploys)
			},
		},
		{
			Name:    "Should propagate ListDeploys error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					ListDeploys().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.ListDeploys(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestGetDeploy(t *testing.T) {
	deploy := &models.Deploy{
		DeployID: "some_id",
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call GetDeploy with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					GetDeploy("some_id").
					Return(deploy, nil)

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.GetDeploy(req, resp)
			},
		},
		{
			Name: "Should return deploy from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					GetDeploy(gomock.Any()).
					Return(deploy, nil)

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.GetDeploy(req, resp)

				var response *models.Deploy
				read(&response)

				reporter.AssertEqual(response, deploy)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.GetDeploy(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate GetDeploy error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					GetDeploy(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.GetDeploy(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestDeleteDeploy(t *testing.T) {
	testCases := []HandlerTestCase{
		{
			Name: "Should call DeleteDeploy with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					DeleteDeploy("some_id").
					Return(nil)

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.DeleteDeploy(req, resp)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.DeleteDeploy(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate DeleteDeploy error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockDeployLogic(ctrl)
				logicMock.EXPECT().
					DeleteDeploy(gomock.Any()).
					Return(errors.Newf(errors.UnexpectedError, "some error"))

				return NewDeployHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.DeleteDeploy(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestCreateDeploy(t *testing.T) {
	request := models.CreateDeployRequest{
		DeployName: "dply_name",
		Dockerrun:  []byte("some dockerrun"),
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call CreateDeploy with correct params",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockDeploy.EXPECT().
					CreateDeploy(request).
					Return(&models.Deploy{}, nil)

				return NewDeployHandler(mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.CreateDeploy(req, resp)
			},
		},
		{
			Name: "Should propagate CreateDeploy error",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockDeploy := mock_logic.NewMockDeployLogic(ctrl)

				mockDeploy.EXPECT().
					CreateDeploy(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewDeployHandler(mockDeploy)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*DeployHandler)
				handler.CreateDeploy(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}
