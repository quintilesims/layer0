package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
	"time"
)

func TestSelectAll(t *testing.T) {
	jobs := []*models.Job{
		&models.Job{
			JobID:       "some_id_1",
			Meta:        map[string]string{},
			TimeCreated: time.Now(),
			LastUpdated: time.Now(),
		},
		&models.Job{
			JobID:       "some_id_2",
			Meta:        map[string]string{},
			TimeCreated: time.Now(),
			LastUpdated: time.Now(),
		},
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name:    "Should return jobs from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					SelectAll().
					Return(jobs, nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.SelectAll(req, resp)

				var response []*models.Job
				read(&response)

				// comparing time fields doesn't seem to work
				reporter.AssertEqual(len(response), 2)
				reporter.AssertEqual(response[0].JobID, jobs[0].JobID)
				reporter.AssertEqual(response[1].JobID, jobs[1].JobID)
			},
		},
		HandlerTestCase{
			Name:    "Should propogate SelectAll error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					SelectAll().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.SelectAll(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestSelectByID(t *testing.T) {
	job := &models.Job{
		JobID: "some_id",
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call SelectByID with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					SelectByID("some_id").
					Return(job, nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.SelectByID(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should return job from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					SelectByID(gomock.Any()).
					Return(job, nil)

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.SelectByID(req, resp)

				var response *models.Job
				read(&response)

				// comparing time fields doesn't seem to work
				reporter.AssertEqual(response.JobID, job.JobID)
			},
		},
		HandlerTestCase{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.SelectByID(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		HandlerTestCase{
			Name: "Should propagate SelectByID error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockJobLogic(ctrl)
				logicMock.EXPECT().
					SelectByID(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewJobHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*JobHandler)
				handler.SelectByID(req, resp)

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
		HandlerTestCase{
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
		HandlerTestCase{
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
		HandlerTestCase{
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
