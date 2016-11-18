package handlers

import (
	"github.com/emicklei/go-restful"
	"gitlab.imshealth.com/xfra/layer0/api/data"
	"gitlab.imshealth.com/xfra/layer0/common/errors"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"net/http"
)

type TagHandler struct {
	TagData data.TagData
}

func NewTagHandler(tagData data.TagData) *TagHandler {
	return &TagHandler{
		TagData: tagData,
	}
}

func (this *TagHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/tag").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Param(service.HeaderParameter("Authorization", "Basic realm authentication token"))

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(this.FindTags).
		Doc("Lists tags, optionally filtered by the query parameters").
		Param(service.QueryParameter("name", "tag name to find").DataType("string")).
		Param(service.QueryParameter("type", "target type for the tag").DataType("string").AllowableValues(data.AllowedTagMap())).
		Returns(200, "OK", []models.EntityWithTags{}))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(this.CreateTag).
		Doc("Create a tag for a service, deploy, or environment").
		Reads(models.EntityTag{}).
		Returns(http.StatusCreated, "Created", models.EntityTag{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.EntityTag{}))

	service.Route(service.DELETE("/").
		Filter(basicAuthenticate).
		To(this.DeleteTag).
		Doc("Delete a tag").
		Reads(models.EntityTag{}).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (this *TagHandler) FindTags(request *restful.Request, response *restful.Response) {
	keys := make(map[string]string)
	for key, val := range request.Request.URL.Query() {
		keys[key] = val[0]
	}

	result, err := this.TagData.GetTags(keys)
	if err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteAsJson(result)
}

func (this *TagHandler) DeleteTag(request *restful.Request, response *restful.Response) {
	req := new(models.EntityTag)
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := this.TagData.Delete(*req); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (this *TagHandler) CreateTag(request *restful.Request, response *restful.Response) {
	req := new(models.EntityTag)
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := this.TagData.Make(*req); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
}
