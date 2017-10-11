package resolver

import (
	"net/url"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestTagResolverQueryParams(t *testing.T) {
	testCases := []struct {
		Name            string
		EntityType      string
		Target          string
		ExpectedQueries []url.Values
	}{
		{
			Name:       "unqualified deploy",
			EntityType: "deploy",
			Target:     "d*",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType: []string{"deploy"},
				client.TagQueryParamFuzz: []string{"d*"},
			}},
		},
		{
			Name:       "fully-qualified deploy",
			EntityType: "deploy",
			Target:     "d:1",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType:    []string{"deploy"},
				client.TagQueryParamFuzz:    []string{"d"},
				client.TagQueryParamVersion: []string{"1"},
			}},
		},
		{
			Name:       "environment",
			EntityType: "environment",
			Target:     "e*",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType: []string{"environment"},
				client.TagQueryParamFuzz: []string{"e*"},
			}},
		},
		{
			Name:       "job",
			EntityType: "job",
			Target:     "j*",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType: []string{"job"},
				client.TagQueryParamFuzz: []string{"j*"},
			}},
		},
		{
			Name:       "unqualified load_balancer",
			EntityType: "load_balancer",
			Target:     "l*",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType: []string{"load_balancer"},
				client.TagQueryParamFuzz: []string{"l*"},
			}},
		},
		{
			Name:       "fully-qualified load_balancer",
			EntityType: "load_balancer",
			Target:     "e*:l*",
			ExpectedQueries: []url.Values{
				{
					client.TagQueryParamType: []string{"environment"},
					client.TagQueryParamFuzz: []string{"e*"},
				},
				{
					client.TagQueryParamType:          []string{"load_balancer"},
					client.TagQueryParamFuzz:          []string{"l*"},
					client.TagQueryParamEnvironmentID: []string{"eid"},
				},
			},
		},
		{
			Name:       "unqualified service",
			EntityType: "service",
			Target:     "s*",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType: []string{"service"},
				client.TagQueryParamFuzz: []string{"s*"},
			}},
		},
		{
			Name:       "fully-qualified service",
			EntityType: "service",
			Target:     "e*:s*",
			ExpectedQueries: []url.Values{
				{
					client.TagQueryParamType: []string{"environment"},
					client.TagQueryParamFuzz: []string{"e*"},
				},
				{
					client.TagQueryParamType:          []string{"service"},
					client.TagQueryParamFuzz:          []string{"s*"},
					client.TagQueryParamEnvironmentID: []string{"eid"},
				},
			},
		},
		{
			Name:       "unqualified task",
			EntityType: "task",
			Target:     "t*",
			ExpectedQueries: []url.Values{{
				client.TagQueryParamType: []string{"task"},
				client.TagQueryParamFuzz: []string{"t*"},
			}},
		},
		{
			Name:       "fully-qualified task",
			EntityType: "task",
			Target:     "e*:t*",
			ExpectedQueries: []url.Values{
				{
					client.TagQueryParamType: []string{"environment"},
					client.TagQueryParamFuzz: []string{"e*"},
				},
				{
					client.TagQueryParamType:          []string{"task"},
					client.TagQueryParamFuzz:          []string{"t*"},
					client.TagQueryParamEnvironmentID: []string{"eid"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// always return an environment entity tag for fully-qualified lookups
			tags := models.Tags{{EntityID: "eid"}}

			mockClient := mock_client.NewMockClient(ctrl)
			for _, expectedQuery := range testCase.ExpectedQueries {
				mockClient.EXPECT().
					ListTags(expectedQuery).
					Return(tags, nil)
			}

			resolver := NewTagResolver(mockClient)
			resolver.Resolve(testCase.EntityType, testCase.Target)
		})
	}
}

func TestTagResolverExactIDMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tags := models.Tags{
		{EntityID: "id1"},
		{EntityID: "id2"},
	}

	mockClient := mock_client.NewMockClient(ctrl)
	mockClient.EXPECT().
		ListTags(gomock.Any()).
		Return(tags, nil)

	resolver := NewTagResolver(mockClient)
	result, err := resolver.Resolve("environment", "id1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"id1"}, result)
}

func TestTagResolverUniqueIDExtraction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tags := models.Tags{
		{EntityID: "id1"},
		{EntityID: "id1"},
		{EntityID: "id2"},
		{EntityID: "id2"},
		{EntityID: "id3"},
		{EntityID: "id3"},
	}

	mockClient := mock_client.NewMockClient(ctrl)
	mockClient.EXPECT().
		ListTags(gomock.Any()).
		Return(tags, nil)

	resolver := NewTagResolver(mockClient)
	result, err := resolver.Resolve("environment", "*")
	if err != nil {
		t.Fatal(err)
	}

	// order matters when doing assert.Equal on slices
	sort.Strings(result)
	assert.Equal(t, []string{"id1", "id2", "id3"}, result)
}
