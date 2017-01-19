package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
	 "github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/config"
)

var testTags = []*models.Tag{
	{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},
	{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl_2"},
	{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
	{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env_2"},
	{EntityID: "j1", EntityType: "job", Key: "name", Value: "job_1"},
	{EntityID: "j2", EntityType: "job", Key: "name", Value: "job_2"},
	{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb_1"},
	{EntityID: "l2", EntityType: "load_balancer", Key: "name", Value: "lb_2"},
	{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc_1"},
	{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc_2"},
	{EntityID: "t1", EntityType: "task", Key: "name", Value: "tsk_1"},
	{EntityID: "t2", EntityType: "task", Key: "name", Value: "tsk_2"},
}

func getTestTagStore(t *testing.T, tags []*models.Tag) *tag_store.MysqlTagStore {
	tagStore := tag_store.NewMysqlTagStore(db.Config{
                Connection: config.DBConnection(),
                DBName:     config.DBName(),
        })

	if err := tagStore.Init(); err != nil {
		t.Fatal(err)
	}

	if err := tagStore.Clear(); err != nil {
		t.Fatal(err)
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	return tagStore
}

func TestFindTags_noParams(t *testing.T) {
	store := getTestTagStore(t, testTags)
	handler := NewTagHandler(store)

	cases := []HandlerTestCase{
		{
			Name:    "no params",
			Request: &TestRequest{},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), len(testTags))
			},
		},
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_byFuzz(t *testing.T) {
	store := getTestTagStore(t, testTags)
	handler := NewTagHandler(store)

	cases := []HandlerTestCase{}
	for _, f := range []string{"d", "e", "j", "l", "s", "t"} {
		fuzz := f

		cases = append(cases, HandlerTestCase{
			Name: fuzz,
			Request: &TestRequest{
				Query: fmt.Sprintf("fuzz=%s", fuzz),
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
			},
		})
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_byID(t *testing.T) {
	store := getTestTagStore(t, testTags)
	handler := NewTagHandler(store)

	cases := []HandlerTestCase{}
	for _, i := range []string{"d1", "d2", "e1", "e2", "j1", "j2", "l1", "l2", "s1", "s2", "t1", "t2"} {
		id := i

		cases = append(cases, HandlerTestCase{
			Name: id,
			Request: &TestRequest{
				Query: fmt.Sprintf("id=%s", id),
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, id)
			},
		})
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_byType(t *testing.T) {
	store := getTestTagStore(t, testTags)
	handler := NewTagHandler(store)

	cases := []HandlerTestCase{}
	for _, t := range []string{"deploy", "environment", "job", "load_balancer", "service", "task"} {
		typeName := t

		cases = append(cases, HandlerTestCase{
			Name: typeName,
			Request: &TestRequest{
				Query: fmt.Sprintf("type=%s", typeName),
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				for _, tag := range tags {
					r.AssertEqual(tag.EntityType, typeName)
				}
			},
		})
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_byVersion(t *testing.T) {
	store := getTestTagStore(t, []*models.Tag{
		{EntityID: "d.1", EntityType: "deploy", Key: "version", Value: "1"},
		{EntityID: "d.2", EntityType: "deploy", Key: "version", Value: "2"},
		{EntityID: "d.3", EntityType: "deploy", Key: "version", Value: "3"},
		{EntityID: "d.4", EntityType: "deploy", Key: "version", Value: "4"},
		{EntityID: "d.5", EntityType: "deploy", Key: "version", Value: "5"},
		{EntityID: "d.6", EntityType: "deploy", Key: "version", Value: "6"},
		{EntityID: "d.7", EntityType: "deploy", Key: "version", Value: "7"},
		{EntityID: "d.8", EntityType: "deploy", Key: "version", Value: "8"},
		{EntityID: "d.9", EntityType: "deploy", Key: "version", Value: "9"},
		{EntityID: "d.10", EntityType: "deploy", Key: "version", Value: "10"},
		{EntityID: "d.11", EntityType: "deploy", Key: "version", Value: "11"},
		{EntityID: "d.12", EntityType: "deploy", Key: "version", Value: "12"},
	})

	handler := NewTagHandler(store)

	callAndRead := func(req *restful.Request, resp *restful.Response, read Readf) []models.EntityWithTags {
		handler.FindTags(req, resp)

		var tags []models.EntityWithTags
		read(&tags)

		return tags
	}

	cases := []HandlerTestCase{
		{
			Name:    "2",
			Request: &TestRequest{Query: "version=2"},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				tags := callAndRead(req, resp, read)
				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "d.2")
			},
		},
		{
			Name:    "7",
			Request: &TestRequest{Query: "version=7"},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				tags := callAndRead(req, resp, read)
				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "d.7")
			},
		},
		{
			Name:    "11",
			Request: &TestRequest{Query: "version=11"},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				tags := callAndRead(req, resp, read)
				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "d.11")
			},
		},
		{
			Name:    "84",
			Request: &TestRequest{Query: "version=84"},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				tags := callAndRead(req, resp, read)
				r.AssertEqual(len(tags), 0)
			},
		},
		{
			Name:    "latest",
			Request: &TestRequest{Query: "version=latest"},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				tags := callAndRead(req, resp, read)
				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "d.12")
			},
		},
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_byArbitraryKeyVal(t *testing.T) {
	store := getTestTagStore(t, []*models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},
		{EntityID: "s1", EntityType: "service", Key: "load_balancer_id", Value: "l1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "e1"},
	})

	handler := NewTagHandler(store)

	cases := []HandlerTestCase{
		{
			Name: "name=dpl_1",
			Request: &TestRequest{
				Query: "name=dpl_1",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "d1")
			},
		},
		{
			Name: "load_balancer_id=l1",
			Request: &TestRequest{
				Query: "load_balancer_id=l1",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "s1")
			},
		},
		{
			Name: "environment_id=e1",
			Request: &TestRequest{
				Query: "environment_id=e1",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 1)
				r.AssertEqual(tags[0].EntityID, "l1")
			},
		},
	}

	RunHandlerTestCases(t, cases)
}

func TestFindTags_byTypeAndFuzz(t *testing.T) {
	store := getTestTagStore(t, testTags)
	handler := NewTagHandler(store)

	cases := []HandlerTestCase{
		{
			Name: "type=deploy&fuzz=d",
			Request: &TestRequest{
				Query: "type=deploy&fuzz=d",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				r.AssertAny(tags[0].EntityID, "d1", "d2")
				r.AssertAny(tags[1].EntityID, "d1", "d2")
			},
		},
		{
			Name: "type=environment&fuzz=e",
			Request: &TestRequest{
				Query: "type=environment&fuzz=e",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				r.AssertAny(tags[0].EntityID, "e1", "e2")
				r.AssertAny(tags[1].EntityID, "e1", "e2")
			},
		},
		{
			Name: "type=job&fuzz=j",
			Request: &TestRequest{
				Query: "type=job&fuzz=j",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				r.AssertAny(tags[0].EntityID, "j1", "j2")
				r.AssertAny(tags[1].EntityID, "j1", "j2")
			},
		},
		{
			Name: "type=load_balancer&fuzz=l",
			Request: &TestRequest{
				Query: "type=load_balancer&fuzz=l",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				r.AssertAny(tags[0].EntityID, "l1", "l2")
				r.AssertAny(tags[1].EntityID, "l1", "l2")
			},
		},
		{
			Name: "type=service&fuzz=s",
			Request: &TestRequest{
				Query: "type=service&fuzz=s",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				r.AssertAny(tags[0].EntityID, "s1", "s2")
				r.AssertAny(tags[1].EntityID, "s1", "s2")
			},
		},
		{
			Name: "type=task&fuzz=t",
			Request: &TestRequest{
				Query: "type=task&fuzz=t",
			},
			Run: func(r *testutils.Reporter, _ interface{}, req *restful.Request, resp *restful.Response, read Readf) {
				handler.FindTags(req, resp)

				var tags []models.EntityWithTags
				read(&tags)

				r.AssertEqual(len(tags), 2)
				r.AssertAny(tags[0].EntityID, "t1", "t2")
				r.AssertAny(tags[1].EntityID, "t1", "t2")
			},
		},
	}

	RunHandlerTestCases(t, cases)
}
