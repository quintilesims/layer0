package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/db/mock_data"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestSelectByQuery(t *testing.T) {
	tags := []models.EntityWithTags{
		models.EntityWithTags{
			EntityID:   "some_id",
			EntityType: "some_type",
			Tags:       []models.EntityTag{},
		},
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call SelectByQuery with correct params",
			Request: &TestRequest{
				Query: "key1=val1&key2=val2",
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				tagDataMock := mock_data.NewMockTagStore(ctrl)

				params := map[string]string{
					"key1": "val1",
					"key2": "val2",
				}

				tagDataMock.EXPECT().
					SelectByQuery(params).
					Return(tags, nil)

				return NewTagHandler(tagDataMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TagHandler)
				handler.FindTags(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should return tags from logic layer",
			Request: &TestRequest{
				Query: "key1=val1&key2=val2",
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				tagDataMock := mock_data.NewMockTagStore(ctrl)

				tagDataMock.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(tags, nil)

				return NewTagHandler(tagDataMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TagHandler)
				handler.FindTags(req, resp)

				var response []models.EntityWithTags
				read(&response)

				reporter.AssertEqual(response, tags)
			},
		},
		HandlerTestCase{
			Name: "Should propogate SelectByQuery error",
			Request: &TestRequest{
				Query: "key1=val1&key2=val2",
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				tagDataMock := mock_data.NewMockTagStore(ctrl)

				tagDataMock.EXPECT().
					SelectByQuery(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewTagHandler(tagDataMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*TagHandler)
				handler.FindTags(req, resp)

				var response models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}
