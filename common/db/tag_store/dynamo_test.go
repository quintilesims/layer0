package tag_store

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
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

func NewTestTagStore(t *testing.T) *DynamoTagStore {
	table := config.TestDynamoTagTableName()
	if table == "" {
		t.Skipf("Skipping test: %s not set", config.TEST_AWS_TAG_DYNAMO_TABLE)
	}

	creds := credentials.NewStaticCredentials(config.AWSAccessKey(), config.AWSSecretKey(), "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AWSRegion()),
	}

	delay, err := time.ParseDuration(config.AWSTimeBetweenRequests())
	if err != nil {
		t.Fatalf("Error parsing time between requests: %v", err)
	}

	session := session.New(awsConfig)
	ticker := time.Tick(delay)
	session.Handlers.Send.PushBack(func(r *request.Request) {
		<-ticker
	})

	store := NewDynamoTagStore(session, table)

	if err := store.Clear(); err != nil {
		t.Fatalf("Error clearing table: %v", err)
	}

	return store
}

func TestDynamoTagStoreInsert(t *testing.T) {
	store := NewTestTagStore(t)

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

func TestDynamoTagStoreDelete(t *testing.T) {
	store := NewTestTagStore(t)

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

func TestDynamoTagStoreSelectByTypeAndID(t *testing.T) {
	store := NewTestTagStore(t)

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

func TestDynamoTagStoreSelectByType(t *testing.T) {
	store := NewTestTagStore(t)

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
