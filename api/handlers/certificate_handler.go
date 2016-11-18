package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"gitlab.imshealth.com/xfra/layer0/api/logic"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"net/http"
)

type CertificateHandler struct {
	Logic logic.CertificateLogic
}

func NewCertificateHandler(certificate logic.CertificateLogic) *CertificateHandler {
	return &CertificateHandler{
		Logic: certificate,
	}
}

func (this *CertificateHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/certificate").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	id := service.PathParameter("id", "identifier of the certificate").
		DataType("string")

	service.Route(service.GET("/").
		To(this.ListCertificates).
		Doc("List all certificates").
		Returns(200, "OK", []models.Certificate{}))

	service.Route(service.GET("/{id}").
		To(this.GetCertificate).
		Doc("Return a certificate").
		Param(id).
		Writes(models.Certificate{}))

	service.Route(service.DELETE("/{id}").
		To(this.DeleteCertificate).
		Filter(basicAuthenticate).
		Doc("Delete a certificate").
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(this.CreateCertificate).
		Doc("Upload a SSL certificate").
		Reads(models.CreateCertificateRequest{}).
		Returns(http.StatusCreated, "Created", models.Certificate{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Certificate{}))

	return service
}

func (this *CertificateHandler) ListCertificates(request *restful.Request, response *restful.Response) {
	certificates, err := this.Logic.ListCertificates()
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(certificates)
}

func (this *CertificateHandler) GetCertificate(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Parameter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	certificate, err := this.Logic.GetCertificate(id)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(certificate)
}

func (this *CertificateHandler) DeleteCertificate(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	if id == "" {
		err := fmt.Errorf("Paramter 'id' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	if err := this.Logic.DeleteCertificate(id); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(``)
}

func (this *CertificateHandler) CreateCertificate(request *restful.Request, response *restful.Response) {
	var req models.CreateCertificateRequest
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	certificate, err := this.Logic.CreateCertificate(req)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(certificate)
}
