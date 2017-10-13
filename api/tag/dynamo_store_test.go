package tag

import (
	"os"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

var TestTags = models.Tags{
	{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
	{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"},
	{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl"},
	{EntityID: "d2", EntityType: "deploy", Key: "version", Value: "2"},

	{EntityID: "e1", EntityType: "environment", Key: "name", Value: "e1"},
	{EntityID: "e2", EntityType: "environment", Key: "name", Value: "e2"},

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

func newTestStore(t *testing.T) *DynamoStore {
	session := config.GetTestAWSSession()
	table := os.Getenv(config.ENVVAR_TEST_AWS_DYNAMO_TAG_TABLE)
	if table == "" {
		t.Skipf("Test table not set (envvar: %s)", config.ENVVAR_TEST_AWS_DYNAMO_TAG_TABLE)
	}

	store := NewDynamoStore(session, table)
	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func TestDynamoStoreInsert(t *testing.T) {
	store := newTestStore(t)

	tags := []models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"},
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc1"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "e1"},
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl1"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"},
	}

	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDynamoStoreDelete(t *testing.T) {
	store := newTestStore(t)

	tags := []models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl1"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"},
	}

	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	if err := store.Delete(tags[0].EntityType, tags[0].EntityID, tags[0].Key); err != nil {
		t.Fatal(err)
	}

	result, err := store.SelectByTypeAndID(tags[1].EntityType, tags[1].EntityID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result[0], tags[1])
}

func TestDynamoStoreSelectAll(t *testing.T) {
	store := newTestStore(t)
	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	for _, tag := range TestTags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, len(TestTags))
}

func TestDynamoStoreSelectByTypeAndID(t *testing.T) {
	store := newTestStore(t)

	for _, tag := range TestTags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	cases := []struct {
		EntityType string
		EntityID   string
		Expected   models.Tags
	}{
		{
			EntityType: "invalid",
			EntityID:   "invalid",
			Expected:   models.Tags{},
		},
		{
			EntityType: "deploy",
			EntityID:   "d1",
			Expected:   TestTags[0:2],
		},
		{
			EntityType: "environment",
			EntityID:   "e1",
			Expected:   TestTags[4:5],
		},
		{
			EntityType: "load_balancer",
			EntityID:   "l1",
			Expected:   TestTags[6:8],
		},
		{
			EntityType: "service",
			EntityID:   "s1",
			Expected:   TestTags[10:12],
		},
		{
			EntityType: "task",
			EntityID:   "t1",
			Expected:   TestTags[14:16],
		},
	}

	for query, c := range cases {
		results, err := store.SelectByTypeAndID(c.EntityType, c.EntityID)
		if err != nil {
			t.Fatalf("Query %d: %v", query, err)
		}

		assert.Len(t, results, len(c.Expected))
		for _, e := range c.Expected {
			assert.Contains(t, results, e)
		}
	}
}

func TestDynamoStoreSelectByType(t *testing.T) {
	store := newTestStore(t)

	for _, tag := range TestTags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	cases := []struct {
		EntityType string
		EntityID   string
		Expected   models.Tags
	}{
		{
			EntityType: "invalid",
			EntityID:   "invalid",
			Expected:   models.Tags{},
		},
		{
			EntityType: "deploy",
			Expected:   TestTags[0:4],
		},
		{
			EntityType: "environment",
			EntityID:   "e1",
			Expected:   TestTags[4:6],
		},
		{
			EntityType: "load_balancer",
			EntityID:   "l1",
			Expected:   TestTags[6:10],
		},
		{
			EntityType: "service",
			EntityID:   "s1",
			Expected:   TestTags[10:14],
		},
		{
			EntityType: "task",
			EntityID:   "t1",
			Expected:   TestTags[14:18],
		},
	}

	for query, c := range cases {
		results, err := store.SelectByType(c.EntityType)
		if err != nil {
			t.Fatalf("Query %d: %v", query, err)
		}

		assert.Len(t, results, len(c.Expected))
		for _, e := range c.Expected {
			assert.Contains(t, results, e)
		}
	}
}
