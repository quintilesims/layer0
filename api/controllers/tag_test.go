package controllers

import (
	"net/url"
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
	glob "github.com/ryanuber/go-glob"
	"github.com/stretchr/testify/assert"
)

func TestListTags(t *testing.T) {
	tags := models.Tags{
		{EntityID: "eid1", EntityType: "environment", Key: "name", Value: "ename1"},
		{EntityID: "eid2", EntityType: "environment", Key: "name", Value: "ename2"},
		{EntityID: "sid1", EntityType: "service", Key: "environment_id", Value: "eid1"},
		{EntityID: "sid1", EntityType: "service", Key: "name", Value: "sname1"},
		{EntityID: "sid2", EntityType: "service", Key: "environment_id", Value: "eid2"},
		{EntityID: "sid2", EntityType: "service", Key: "name", Value: "sname2"},
		{EntityID: "did1", EntityType: "deploy", Key: "name", Value: "dname1"},
		{EntityID: "did1", EntityType: "deploy", Key: "version", Value: "1"},
		{EntityID: "did2", EntityType: "deploy", Key: "name", Value: "dname2"},
		{EntityID: "did2", EntityType: "deploy", Key: "version", Value: "2"},
	}

	store := tag.NewMemoryStore()
	for _, tag := range tags {
		store.Insert(tag)
	}

	testCases := []struct {
		Name        string
		Query       url.Values
		CheckResult func(t *testing.T, result models.Tags)
	}{
		{
			Name:  "no params",
			Query: url.Values{},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.Len(t, result, len(tags))
			},
		},
		{
			Name: "type",
			Query: url.Values{
				client.TagQueryParamType: []string{"environment"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return tag.EntityType != "environment"
				}))
			},
		},
		{
			Name: "id",
			Query: url.Values{
				client.TagQueryParamID: []string{"eid1"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return tag.EntityID != "eid1"
				}))
			},
		},
		{
			Name: "name",
			Query: url.Values{
				client.TagQueryParamName: []string{"ename1"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return tag.Key == "name" && tag.Value != "ename1"
				}))
			},
		},
		{
			Name: "fuzz",
			Query: url.Values{
				client.TagQueryParamName: []string{"e*"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return !glob.Glob("e*", tag.EntityID) || tag.Key == "name" && !glob.Glob("e*", tag.Value)
				}))
			},
		},
		{
			Name: "version=1",
			Query: url.Values{
				client.TagQueryParamVersion: []string{"1"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return tag.Key == "version" && tag.Value != "1"
				}))
			},
		},
		{
			Name: "version=latest",
			Query: url.Values{
				client.TagQueryParamVersion: []string{"latest"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return tag.Key == "version" && tag.Value != "2"
				}))
			},
		},
		{
			Name: "environment_id",
			Query: url.Values{
				client.TagQueryParamEnvironmentID: []string{"eid1"},
			},
			CheckResult: func(t *testing.T, result models.Tags) {
				assert.False(t, result.Any(func(tag models.Tag) bool {
					return tag.Key == "environment_id" && tag.Value != "eid1"
				}))
			},
		},
		 {
                        Name: "fuzz+type",
                        Query: url.Values{
                                client.TagQueryParamType: []string{"environment"},
				client.TagQueryParamFuzz: []string{"e*"},
                        },
                        CheckResult: func(t *testing.T, result models.Tags) {
                                assert.False(t, result.Any(func(tag models.Tag) bool {
					wrongType := tag.EntityType != "environment"
					noFuzzMatch := !glob.Glob("e*", tag.EntityID) || tag.Key == "name" && !glob.Glob("e*", tag.Value)
					return wrongType || noFuzzMatch
                                }))
                        },
                },
		   {
                        Name: "fuzz+type+environment_id",
                        Query: url.Values{
                                client.TagQueryParamType: []string{"service"},
                                client.TagQueryParamFuzz: []string{"s*"},
				client.TagQueryParamEnvironmentID: []string{"eid1"},
                        },
                        CheckResult: func(t *testing.T, result models.Tags) {
                                assert.False(t, result.Any(func(tag models.Tag) bool {
                                        wrongType := tag.EntityType != "service"
					noFuzzMatch := !glob.Glob("s*", tag.EntityID) || tag.Key == "name" && !glob.Glob("s*", tag.Value)
					wrongEnvironment := tag.Key == "environment_id" && tag.Value != "eid1"
					return wrongType || noFuzzMatch || wrongEnvironment
                                }))
                        },
                },
	}

	controller := NewTagController(store)
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			c := newFireballContext(t, nil, nil)
			c.Request.URL.RawQuery = testCase.Query.Encode()

			resp, err := controller.ListTags(c)
			if err != nil {
				t.Fatal(err)
			}

			var response models.Tags
			recorder := unmarshalBody(t, resp, &response)

			response.Sort()
			assert.Equal(t, 200, recorder.Code)
			testCase.CheckResult(t, response)
		})
	}
}
