package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"net/http"
	"strconv"
	"strings"
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
		Param(service.QueryParameter("type", "Require the EntityType field match the specified parameter").DataType("string")).
		Param(service.QueryParameter("id", "Require the EntityID field match the specified parameter").DataType("string")).
		Param(service.QueryParameter("fuzz", "Require the prefix of the EntityID field or 'name' tag match the specified parameter").DataType("string")).
		Param(service.QueryParameter("version", "Require the 'version' tag match the specified parameter. If 'latest' is used, only the latest will be returned").DataType("string")).
		Param(service.QueryParameter("environment_id", "Require the 'environment_id' tag match the specified parameter").DataType("string")).
		Returns(200, "OK", []models.EntityWithTags{}))

	service.Route(service.POST("/").
		Filter(basicAuthenticate).
		To(t.CreateTag).
		Doc("Create a tag for a service, deploy, or environment").
		Reads(models.Tag{}).
		Returns(http.StatusCreated, "Created", models.Tag{}).
		Returns(400, "Invalid request", models.ServerError{}).
		Writes(models.Tag{}))

	id := service.PathParameter("id", "identifier of the tag").
		DataType("integer")

	service.Route(service.DELETE("/").
		Filter(basicAuthenticate).
		To(t.DeleteTag).
		Doc("Delete a tag").
		Reads(models.Tag{}).
		Param(id).
		Returns(http.StatusNoContent, "Deleted", nil))

	return service
}

func (t *TagHandler) FindTags(request *restful.Request, response *restful.Response) {
	params := make(map[string]string)
	for key, val := range request.Request.URL.Query() {
		params[key] = val[0]
	}

	var entityType string
	var entityID string
	var fuzz string
	var latestVersion bool

	// break out special filter params so we don't filter
	// them by tag.Key and tag.Value
	if val, ok := params["type"]; ok {
		entityType = val
		delete(params, "type")
	}

	if val, ok := params["id"]; ok {
		entityID = val
		delete(params, "id")
	}

	if val, ok := params["fuzz"]; ok {
		fuzz = val
		delete(params, "fuzz")
	}

	if val, ok := params["version"]; ok && val == "latest" {
		latestVersion = true
		delete(params, "version")
	}

	if entityType == "" {
		err := fmt.Errorf("Parameter 'type' is required")
		BadRequest(response, errors.MissingParameter, err)
		return
	}

	var query func() (models.Tags, error)
	if entityID == "" {
		query = func() (models.Tags, error) { return t.TagStore.SelectByType(entityType) }
	} else {
		query = func() (models.Tags, error) { return t.TagStore.SelectByTypeAndID(entityType, entityID) }
	}

	tags, err := query()
	if err != nil {
		ReturnError(response, err)
		return
	}

	// filter the non-special params the by tag.Name and tag.Value (e.g. environment_id, version)
	ewts := tags.GroupByEntity()
	for key, val := range params {
		ewts = ewts.WithKey(key).WithValue(val)
	}

	if fuzz != "" {
		ewts = ewts.RemoveIf(func(e models.EntityWithTags) bool {
			// don't remove if the EntityID matches the fuzz prefix
			if strings.HasPrefix(e.EntityID, fuzz) {
				return false
			}

			// don't remove if the name tag matches the fuzz prefix
			if tag, ok := e.Tags.WithKey("name").First(); ok {
				return !strings.HasPrefix(tag.Value, fuzz)
			}

			return true
		})
	}

	if latestVersion {
		indexOfLatestVersion := -1
		latestVersion := -1

		for i, ewt := range ewts {
			if current, ok := ewt.Tags.WithKey("version").First(); ok {
				currentVersion, err := strconv.Atoi(current.Value)
				if err != nil {
					ReturnError(response, err)
					return
				}

				if currentVersion > latestVersion {
					latestVersion = currentVersion
					indexOfLatestVersion = i
				}
			}
		}

		if latestVersion == -1 {
			ewts = models.EntitiesWithTags{}
		} else {
			ewts = models.EntitiesWithTags{ewts[indexOfLatestVersion]}
		}
	}

	response.WriteAsJson(ewts)
}

func (t *TagHandler) DeleteTag(request *restful.Request, response *restful.Response) {
	var tag models.Tag
	if err := request.ReadEntity(&tag); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := t.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (t *TagHandler) CreateTag(request *restful.Request, response *restful.Response) {
	var tag models.Tag
	if err := request.ReadEntity(&tag); err != nil {
		BadRequest(response, errors.InvalidJSON, err)
		return
	}

	if err := t.TagStore.Insert(tag); err != nil {
		ReturnError(response, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
}
