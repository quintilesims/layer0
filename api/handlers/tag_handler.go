package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"net/http"
)

type TagHandler struct {
	TagStore tag_store.TagStore
}

func NewTagHandler(tagData tag_store.TagStore) *TagHandler {
	return &TagHandler{
		TagStore: tagData,
	}
}

func (t *TagHandler) Routes() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/tag").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Param(service.HeaderParameter("Authorization", "Basic realm authentication token"))

	service.Route(service.GET("/").
		Filter(basicAuthenticate).
		To(t.FindTags).
		Doc("Lists tags, optionally filtered by the query parameters").
		Param(service.QueryParameter("name", "tag name to find").DataType("string")).
		Param(service.QueryParameter("type", "target type for the tag").DataType("string")).
		Returns(200, "OK", []models.EntityWithTags{}))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(t.CreateTag).
		Doc("Create a tag for a service, deploy, or environment").
		Reads(models.Tag{}).
		Returns(http.StatusCreated, "Created", models.Tag{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Tag{}))

	service.Route(service.DELETE("/").
		Filter(basicAuthenticate).
		To(t.DeleteTag).
		Doc("Delete a tag").
		Reads(models.Tag{}).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (t *TagHandler) FindTags(request *restful.Request, response *restful.Response) {
	keys := make(map[string]string)
	for key, val := range request.Request.URL.Query() {
		keys[key] = val[0]
	}

	// todo: re-create t

	response.WriteAsJson(keys)
}

func (t *TagHandler) DeleteTag(request *restful.Request, response *restful.Response) {
	req := new(models.Tag)
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := t.TagStore.Delete(req); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (t *TagHandler) CreateTag(request *restful.Request, response *restful.Response) {
	req := new(models.Tag)
	if err := request.ReadEntity(&req); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := t.TagStore.Insert(req); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
}
