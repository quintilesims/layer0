package handlers

import (
	"testing"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func TestListJobs(t *testing.T) {
	jobs := []*models.Job{
		{
			JobID:       "some_id_1",
			Meta:        map[string]string{},
			TimeCreated: time.Now(),
		},
		{
			JobID:       "some_id_2",
			Meta:        map[string]string{},
			TimeCreated: time.Now(),
		},
	}

	testCases := []HandlerTestCase{
		{
			Name:    "Should return jobs from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					ListJobs().
					Return(jobs, nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.ListJobs(req, resp)

				var response []*models.Job
				read(&response)

				// comparing time fields doesn't seem to work
				reporter.AssertEqual(len(response), 2)
				reporter.AssertEqual(response[0].JobID, jobs[0].JobID)
				reporter.AssertEqual(response[1].JobID, jobs[1].JobID)
			},
		},
		{
			Name:    "Should propagate ListJobs error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					ListJobs().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.ListJobs(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestGetJob(t *testing.T) {
	job := &models.Job{
		JobID: "some_id",
	}

	testCases := []HandlerTestCase{
		{
			Name: "Should call GetJob with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					GetJob("some_id").
					Return(job, nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.GetJob(req, resp)
			},
		},
		{
			Name: "Should return job from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					GetJob(gomock.Any()).
					Return(job, nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.GetJob(req, resp)

				var response *models.Job
				read(&response)

				// comparing time fields doesn't seem to work
				reporter.AssertEqual(response.JobID, job.JobID)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.GetJob(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate GetJob error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					GetJob(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.GetJob(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestDelete(t *testing.T) {
	testCases := []HandlerTestCase{
		{
			Name: "Should call Delete with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					Delete("some_id").
					Return(nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.Delete(req, resp)
			},
		},
		{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.Delete(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		{
			Name: "Should propagate Delete error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					Delete(gomock.Any()).
					Return(errors.Newf(errors.UnexpectedError, "some error"))

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.Delete(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}
