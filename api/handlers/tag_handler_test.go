package handlers

import (
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

var TestTags = models.Tags{
	{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
	{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"},
	{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl"},
	{EntityID: "d2", EntityType: "deploy", Key: "version", Value: "2"},

	{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"},
	{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env2"},

	{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb1"},
	{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
	{EntityID: "l2", EntityType: "load_balancer", Key: "name", Value: "lb2"},
	{EntityID: "l2", EntityType: "load_balancer", Key: "environment_id", Value: "e2"},

	{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc1"},
	{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
	{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc2"},
	{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "e2"},

	{EntityID: "t1", EntityType: "task", Key: "name", Value: "tsk1"},
	{EntityID: "t1", EntityType: "task", Key: "environment_id", Value: "e1"},
	{EntityID: "t2", EntityType: "task", Key: "name", Value: "tsk2"},
	{EntityID: "t2", EntityType: "task", Key: "environment_id", Value: "e2"},
}

func getTestTagStore(t *testing.T, tags models.Tags) tag_store.TagStore {
	tagStore := tag_store.NewMemoryTagStore()
	if err := tagStore.Init(); err != nil {
		t.Fatal(err)
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	return tagStore
}

func TestFindTags(t *testing.T) {
	store := getTestTagStore(t, TestTags)
	handler := NewTagHandler(store)

	cases := []HandlerTestCase{
		{
			Name: "type=environment",
			Request: &TestRequest{
				Query: "type=environment",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				for _, tag := range tags {
					r.AssertEqual(tag.EntityType, "environment")
				}
			},
		},
		{
			Name: "type=service&id=s1",
			Request: &TestRequest{
				Query: "type=service&id=s1",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityType, "service")
				r.AssertEqual(tags[0].EntityID, "s1")
			},
		},
		{
			Name: "type=task&environment_id=e1",
			Request: &TestRequest{
				Query: "type=task&environment_id=e1",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityType, "task")
				r.AssertEqual(tags[0].EntityID, "t1")
			},
		},
		{
			Name: "type=deploy&fuzz=dpl",
			Request: &TestRequest{
				Query: "type=deploy&fuzz=dpl",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				for _, tag := range tags {
					r.AssertEqual(tag.EntityType, "deploy")
				}
			},
		},
		{
			Name: "type=deploy&fuzz=dpl&version=1",
			Request: &TestRequest{
				Query: "type=deploy&fuzz=dpl&version=1",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityType, "deploy")
				r.AssertEqual(tags[0].EntityID, "d1")
			},
		},
		{
			Name: "type=deploy&fuzz=dpl&version=latest",
			Request: &TestRequest{
				Query: "type=deploy&fuzz=dpl&version=latest",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityType, "deploy")
				r.AssertEqual(tags[0].EntityID, "d2")
			},
		},
	}

	RunHandlerTestCases(t, cases)
}
