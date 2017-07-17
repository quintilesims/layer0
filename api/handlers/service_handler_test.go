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

func TestListServices(t *testing.T) {
	services := []*models.ServiceSummary{
		{
			ServiceID: "some_id_1",
		},
		{
			ServiceID: "some_id_2",
		},
	}

	testCases := []HandlerTestCase{
		{
			Name:    "Should return services from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				svcLogicMock.EXPECT().
					ListServices().
					Return(services, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.ListServices(req, resp)

				var response []*models.ServiceSummary
				read(&response)

				reporter.AssertEqual(response, services)
			},
		},
		{
			Name:    "Should propagate ListServices error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				svcLogicMock.EXPECT().
					ListServices().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.ListServices(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestGetService(t *testing.T) {
	service := &models.Service{
		ServiceID: "some_id",
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call GetService with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				svcLogicMock.EXPECT().
					GetService("some_id").
					Return(service, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.GetService(req, resp)
			},
		},
		{
			Name: "Should return service from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				svcLogicMock.EXPECT().
					GetService(gomock.Any()).
					Return(service, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.GetService(req, resp)

				var response *models.Service
				read(&response)

				reporter.AssertEqual(response, service)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewServiceHandler(svcLogicMock, jobLogicMock)

			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.GetService(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate GetService error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				svcLogicMock.EXPECT().
					GetService(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.GetService(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestDeleteService(t *testing.T) {
	testCases := []HandlerTestCase{
		{
			Name: "Should call CreateJob with correct params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				jobLogicMock.EXPECT().
					CreateJob(types.DeleteServiceJob, "some_id").
					Return(&models.Job{}, nil)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.DeleteService(req, resp)
			},
		},
		{
			Name: "Should set Location and X-Jobid headers",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				jobLogicMock.EXPECT().
					CreateJob(gomock.Any(), gomock.Any()).
					Return(&models.Job{JobID: "job_id"}, nil)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.DeleteService(req, resp)

				header := resp.Header()
				reporter.AssertInSlice("/job/job_id", header["Location"])
				reporter.AssertInSlice("job_id", header["X-Jobid"])
			},
		},
		{
			Name: "Should propagate CreateJob error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				jobLogicMock.EXPECT().
					CreateJob(gomock.Any(), gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.DeleteService(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				svcLogicMock := mock_logic.NewMockServiceLogic(ctrl)
				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)

				return NewServiceHandler(svcLogicMock, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.DeleteService(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.MissingParameter), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestCreateService(t *testing.T) {
	request := models.CreateServiceRequest{
		EnvironmentID: "env_id",
		ServiceName:   "svc_name",
		DeployID:      "dply_id",
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call CreateService with correct params",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockService := mock_logic.NewMockServiceLogic(ctrl)

				mockService.EXPECT().
					CreateService(request).
					Return(&models.Service{}, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewServiceHandler(mockService, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.CreateService(req, resp)
			},
		},
		{
			Name: "Should propagate CreateService error",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockService := mock_logic.NewMockServiceLogic(ctrl)

				mockService.EXPECT().
					CreateService(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewServiceHandler(mockService, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.CreateService(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestScaleService(t *testing.T) {
	request := models.ScaleServiceRequest{
		DesiredCount: int64(2),
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call ScaleService with correct params",
			Request: &TestRequest{
				Body:       request,
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockService := mock_logic.NewMockServiceLogic(ctrl)

				mockService.EXPECT().
					ScaleService("some_id", 2).
					Return(&models.Service{}, nil)

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewServiceHandler(mockService, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.ScaleService(req, resp)
			},
		},
		{
			Name: "Should propagate ScaleService error",
			Request: &TestRequest{
				Body:       request,
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockService := mock_logic.NewMockServiceLogic(ctrl)

				mockService.EXPECT().
					ScaleService(gomock.Any(), gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				jobLogicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewServiceHandler(mockService, jobLogicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*ServiceHandler)
				handler.ScaleService(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}
