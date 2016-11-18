package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/api/logic/mock_logic"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"testing"
)

func TestListCertificates(t *testing.T) {
	certificates := []*models.Certificate{
		&models.Certificate{
			CertificateID: "some_id_1",
		},
		&models.Certificate{
			CertificateID: "some_id_2",
		},
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name:    "Should return certificates from logic layer",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					ListCertificates().
					Return(certificates, nil)

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.ListCertificates(req, resp)

				var response []*models.Certificate
				read(&response)

				reporter.AssertEqual(response, certificates)
			},
		},
		HandlerTestCase{
			Name:    "Should propogate ListCertificates error",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					ListCertificates().
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.ListCertificates(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestGetCertificate(t *testing.T) {
	certificate := &models.Certificate{
		CertificateID: "some_id",
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call GetCertificate with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					GetCertificate("some_id").
					Return(certificate, nil)

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.GetCertificate(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should return certificate from logic layer",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					GetCertificate(gomock.Any()).
					Return(certificate, nil)

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.GetCertificate(req, resp)

				var response *models.Certificate
				read(&response)

				reporter.AssertEqual(response, certificate)
			},
		},
		HandlerTestCase{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.GetCertificate(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		HandlerTestCase{
			Name: "Should propagate GetCertificate error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					GetCertificate(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.GetCertificate(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestDeleteCertificate(t *testing.T) {
	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call DeleteCertificate with proper params",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					DeleteCertificate("some_id").
					Return(nil)

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.DeleteCertificate(req, resp)
			},
		},
		HandlerTestCase{
			Name:    "Should return MissingParameter error with no id",
			Request: &TestRequest{},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.DeleteCertificate(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.MissingParameter))
			},
		},
		HandlerTestCase{
			Name: "Should propagate DeleteCertificate error",
			Request: &TestRequest{
				Parameters: map[string]string{"id": "some_id"},
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				logicMock := mock_logic.NewMockCertificateLogic(ctrl)
				logicMock.EXPECT().
					DeleteCertificate(gomock.Any()).
					Return(errors.Newf(errors.UnexpectedError, "some error"))

				return NewCertificateHandler(logicMock)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.DeleteCertificate(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(response.ErrorCode, int64(errors.UnexpectedError))
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}

func TestCreateCertificate(t *testing.T) {
	request := models.CreateCertificateRequest{
		CertificateName:  "cert_name",
		IntermediateCert: "intermed",
		PrivateKey:       "private",
		PublicCert:       "public",
	}

	testCases := []HandlerTestCase{
		HandlerTestCase{
			Name: "Should call CreateCertificate with correct params",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockCertificate := mock_logic.NewMockCertificateLogic(ctrl)

				mockCertificate.EXPECT().
					CreateCertificate(request).
					Return(&models.Certificate{}, nil)

				return NewCertificateHandler(mockCertificate)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.CreateCertificate(req, resp)
			},
		},
		HandlerTestCase{
			Name: "Should propagate CreateCertificate error",
			Request: &TestRequest{
				Body: request,
			},
			Setup: func(ctrl *gomock.Controller) interface{} {
				mockCertificate := mock_logic.NewMockCertificateLogic(ctrl)

				mockCertificate.EXPECT().
					CreateCertificate(gomock.Any()).
					Return(nil, errors.Newf(errors.UnexpectedError, "some error"))

				return NewCertificateHandler(mockCertificate)
			},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler := target.(*CertificateHandler)
				handler.CreateCertificate(req, resp)

				var response *models.ServerError
				read(&response)

				reporter.AssertEqual(int64(errors.UnexpectedError), response.ErrorCode)
			},
		},
	}

	RunHandlerTestCases(t, testCases)
}
