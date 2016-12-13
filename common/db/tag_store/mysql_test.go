package tag_store

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func getTestTags() models.Tags {
	return models.Tags{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"},
		{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env2"},
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc1"},
		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc2"},
		{EntityID: "l1", EntityType: "load_balancer", Key: "name", Value: "lb1"},
		{EntityID: "l2", EntityType: "load_balancer", Key: "name", Value: "lb2"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"},
		{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d2", EntityType: "deploy", Key: "version", Value: "2"},
	}
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
	defer store.Close()

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
	defer store.Close()

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
	defer store.Close()

	assertTagsMatch(t, store, tags)
}

func TestMysqlTagStoreSelectByEntityID(t *testing.T) {
	tags := getTestTags()
	store := NewTestTagStoreWithTags(t, tags)
	defer store.Close()

	r1, err := store.SelectByEntityID("e1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(r1), 1)
	testutils.AssertEqual(t, r1[0].EntityID, "e1")

	r2, err := store.SelectByEntityID("d1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(r2), 2)
	testutils.AssertEqual(t, r2[0].EntityID, "d1")
	testutils.AssertEqual(t, r2[1].EntityID, "d1")
}

func TestMysqlTagStoreSelectByEntityType(t *testing.T) {
	tags := getTestTags()
	store := NewTestTagStoreWithTags(t, tags)
	defer store.Close()

	// environemnt, service, loadbalnccer, deplo
	environments, err := store.SelectByEntityType("environment")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(environments), 2)
	for _, tag := range environments {
		testutils.AssertEqual(t, tag.EntityType, "environment")
	}

	services, err := store.SelectByEntityType("service")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(services), 2)
	for _, tag := range services {
		testutils.AssertEqual(t, tag.EntityType, "service")
	}

	loadBalancers, err := store.SelectByEntityType("load_balancer")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(loadBalancers), 2)
	for _, tag := range loadBalancers {
		testutils.AssertEqual(t, tag.EntityType, "load_balancer")
	}

	deploys, err := store.SelectByEntityType("deploy")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(deploys), 4)
	for _, tag := range deploys {
		testutils.AssertEqual(t, tag.EntityType, "deploy")
	}
}
