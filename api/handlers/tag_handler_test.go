package handlers

import (
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func getTagStore(t *testing.T) *tag_store.MysqlTagStore {
	config := testutils.GetDBConfig()
	tagStore := tag_store.NewMysqlTagStore(config)
	if err := tagStore.Init(); err != nil {
		t.Fatal(err)
	}

	if err := tagStore.Clear(); err != nil {
		t.Fatal(err)
	}

	return tagStore
}

func addTags(t *testing.T, store *tag_store.MysqlTagStore, tags []*models.Tag) {
	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}
}

// by type
// by id
// by fuzz
// by version
func TestFindTags_byType(t *testing.T) {
	store := getTagStore(t)

	addTags(t, store, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc_1"},
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb_1"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},

		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc_2"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "e2"},
		{EntityID: "s2", EntityType: "service", Key: "load_balancer_id", Value: "l2"},
	})

	handler := NewTagHandler(store)
	print(handler)
}

func TestFindTags_environmentQuery(t *testing.T) {
	store := getTagStore(t)
	handler := NewTagHandler(store)

	addTags(t, store, []*models.Tag{
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc_1"},
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb_1"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},

		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc_2"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "e2"},
		{EntityID: "s2", EntityType: "service", Key: "load_balancer_id", Value: "l2"},
	})


	cases := []HandlerTestCase{
		{
			Name:    "Should return jobs from logic layer",
			Request: &TestRequest{},
			Run: func(reporter *testutils.Reporter, target interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var response []models.EntityWithTags
				read(&response)
			},
		},
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_jobQuery(t *testing.T) {

}

func TestFindTags_loadBalancerQuery(t *testing.T) {

}

func TestFindTags_serviceQuery(t *testing.T) {

}

func TestFindTags_taskQuery(t *testing.T) {

}
