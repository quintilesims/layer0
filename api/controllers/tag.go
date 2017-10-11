package controllers

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	glob "github.com/ryanuber/go-glob"
	"github.com/zpatrick/fireball"
)

type TagController struct {
	TagStore tag.Store
}

func NewTagController(t tag.Store) *TagController {
	return &TagController{
		TagStore: t,
	}
}

func (t *TagController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/tag",
			Handlers: fireball.Handlers{
				"GET":    t.ListTags,
				"POST":   t.CreateTag,
				"DELETE": t.DeleteTag,
			},
		},
	}
}

func (t *TagController) DeleteTag(c *fireball.Context) (fireball.Response, error) {
	var tag models.Tag
	if err := json.NewDecoder(c.Request.Body).Decode(&tag); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := t.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, nil)
}

func (t *TagController) CreateTag(c *fireball.Context) (fireball.Response, error) {
	var tag models.Tag
	if err := json.NewDecoder(c.Request.Body).Decode(&tag); err != nil {
		return nil, errors.New(errors.InvalidRequest, err)
	}

	if err := t.TagStore.Insert(tag); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, tag)
}

func (t *TagController) ListTags(c *fireball.Context) (fireball.Response, error) {
	query := c.Request.URL.Query()

	tags, err := t.selectTagsFromQuery(query)
	if err != nil {
		return nil, err
	}

	entityTags := tags.GroupByID()
	if environmentID := query.Get(client.TagQueryParamEnvironmentID); environmentID != "" {
		filterTagsByEnvironmentID(entityTags, environmentID)
	}

	if name := query.Get(client.TagQueryParamName); name != "" {
		filterTagsByName(entityTags, name)
	}

	if fuzz := query.Get(client.TagQueryParamFuzz); fuzz != "" {
		filterTagsByIDOrNameGlob(entityTags, fuzz)
	}

	if version := query.Get(client.TagQueryParamVersion); version != "" {
		if version == "latest" {
			filterTagsByLatestVersion(entityTags)
		} else {
			filterTagsByVersion(entityTags, version)
		}
	}

	filteredTags := models.Tags{}
	for _, tags := range entityTags {
		filteredTags = append(filteredTags, tags...)
	}

	return fireball.NewJSONResponse(200, filteredTags)
}

func (t *TagController) selectTagsFromQuery(query url.Values) (models.Tags, error) {
	entityID := query.Get(client.TagQueryParamID)
	entityType := query.Get(client.TagQueryParamType)

	if entityType != "" && entityID != "" {
		return t.TagStore.SelectByTypeAndID(entityType, entityID)
	}

	if entityType != "" {
		return t.TagStore.SelectByType(entityType)
	}

	if entityID != "" {
		tags, err := t.TagStore.SelectAll()
		if err != nil {
			return nil, err
		}

		return tags.WithID(entityID), nil
	}

	return t.TagStore.SelectAll()
}

func filterTagsByEnvironmentID(entityTags map[string]models.Tags, environmentID string) {
	for entityID, tags := range entityTags {
		hasMatch := tags.Any(func(t models.Tag) bool {
			return t.Key == "environment_id" && t.Value == environmentID
		})

		if !hasMatch {
			delete(entityTags, entityID)
		}
	}
}

func filterTagsByName(entityTags map[string]models.Tags, name string) {
	for entityID, tags := range entityTags {
		hasMatch := tags.Any(func(t models.Tag) bool {
			return t.Key == "name" && t.Value == name
		})

		if !hasMatch {
			delete(entityTags, entityID)
		}
	}
}

func filterTagsByIDOrNameGlob(entityTags map[string]models.Tags, pattern string) {
	for entityID, tags := range entityTags {
		hasMatch := tags.Any(func(t models.Tag) bool {
			if glob.Glob(pattern, t.EntityID) {
				return true
			}

			return t.Key == "name" && glob.Glob(pattern, t.Value)
		})

		if !hasMatch {
			delete(entityTags, entityID)
		}
	}
}

func filterTagsByVersion(entityTags map[string]models.Tags, version string) {
	for entityID, tags := range entityTags {
		hasMatch := tags.Any(func(t models.Tag) bool {
			return t.Key == "version" && t.Value == version
		})

		if !hasMatch {
			delete(entityTags, entityID)
		}
	}
}

func filterTagsByLatestVersion(entityTags map[string]models.Tags) {
	var latest int

	// find the latest version
	for _, tags := range entityTags {
		if tag, ok := tags.WithKey("version").First(); ok {
			current, _ := strconv.Atoi(tag.Value)
			if current > latest {
				latest = current
			}
		}
	}

	// remove all but the latest version
	for entityID, tags := range entityTags {
		hasMatch := tags.Any(func(t models.Tag) bool {
			return t.Key == "version" && t.Value == strconv.Itoa(latest)
		})

		if !hasMatch {
			delete(entityTags, entityID)
		}
	}
}
