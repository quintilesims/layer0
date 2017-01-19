package tag_store

import (
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"reflect"
	"testing"
)

func getTestTags() models.Tags {
	return models.Tags{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"},
		{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env2"},
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc1"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "env1"},
		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc2"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "env2"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb1"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "environment_id", Value: "env1"},
		{EntityID: "l2", EntityType: "load_balancer", Key: "name", Value: "lb2"},
		{EntityID: "l2", EntityType: "load_balancer", Key: "environment_id", Value: "env2"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"},
		{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d2", EntityType: "deploy", Key: "version", Value: "2"},
	}
}

func NewTestTagStore(t *testing.T) *MysqlTagStore {
	store := NewMysqlTagStore(db.Config{
		Connection: config.DBConnection(),
		DBName:     config.DBName(),
	})

	if err := store.Init(); err != nil {
		t.Fatal(err)
	}

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func NewTestTagStoreWithTags(t *testing.T, tags models.Tags) *MysqlTagStore {
	store := NewTestTagStore(t)
	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	return store
}

func assertTagsMatch(t *testing.T, store *MysqlTagStore, expected models.Tags) {
	tags, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(tags), len(expected))
	for _, tag := range tags {
		testutils.AssertInSlice(t, tag, expected)
	}
}

func TestMysqlTagStoreInsert(t *testing.T) {
	store := NewTestTagStore(t)

	tags := getTestTags()
	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	assertTagsMatch(t, store, tags)
}

func TestMysqlTagStoreDelete(t *testing.T) {
	tags := getTestTags()
	store := NewTestTagStoreWithTags(t, tags)

	for _, tag := range tags[:2] {
		if err := store.Delete(tag); err != nil {
			t.Fatal(err)
		}
	}

	assertTagsMatch(t, store, tags[2:])
}

func TestMysqlTagStoreSelectAll(t *testing.T) {
	tags := getTestTags()
	store := NewTestTagStoreWithTags(t, tags)

	assertTagsMatch(t, store, tags)
}

func TestMysqlTagStoreSelectByQuery(t *testing.T) {
	testTags := getTestTags()
	store := NewTestTagStoreWithTags(t, testTags)

	cases := []struct {
		EntityID   string
		EntityType string
		Expected   models.Tags
	}{
		// filter by nothing
		{
			Expected: testTags,
		},
		// filter by entityID
		{
			EntityID: "invalid",
			Expected: models.Tags{},
		},
		{
			EntityID: "e1",
			Expected: testTags[0:1],
		},
		{
			EntityID: "e2",
			Expected: testTags[1:2],
		},
		{
			EntityID: "s1",
			Expected: testTags[2:4],
		},
		{
			EntityID: "s2",
			Expected: testTags[4:6],
		},
		{
			EntityID: "l1",
			Expected: testTags[6:8],
		},
		{
			EntityID: "l2",
			Expected: testTags[8:10],
		},
		{
			EntityID: "d1",
			Expected: testTags[10:12],
		},
		{
			EntityID: "d2",
			Expected: testTags[12:14],
		},
		// filter by entityType
		{
			EntityType: "invalid",
			Expected:   models.Tags{},
		},
		{
			EntityType: "environment",
			Expected:   testTags[0:2],
		},
		{
			EntityType: "service",
			Expected:   testTags[2:6],
		},
		{
			EntityType: "load_balancer",
			Expected:   testTags[6:10],
		},
		{
			EntityType: "deploy",
			Expected:   testTags[10:14],
		},
		// filter by both
		{
			EntityType: "invalid",
			EntityID:   "invalid",
			Expected:   models.Tags{},
		},
		{
			EntityType: "environment",
			EntityID:   "e1",
			Expected:   testTags[0:1],
		},
		{
			EntityType: "service",
			EntityID:   "s1",
			Expected:   testTags[2:4],
		},
		{
			EntityType: "load_balancer",
			EntityID:   "l1",
			Expected:   testTags[6:8],
		},
		{
			EntityType: "deploy",
			EntityID:   "d1",
			Expected:   testTags[10:12],
		},
	}

	for i, c := range cases {
		tags, err := store.SelectByQuery(c.EntityType, c.EntityID)
		if err != nil {
			t.Fatalf("Query %#v: %v", i, err)
		}

		if !reflect.DeepEqual(tags, c.Expected) {
			t.Errorf("Query %#v: got %v, expected: %v", i, tags, c.Expected)
		}
	}
}
