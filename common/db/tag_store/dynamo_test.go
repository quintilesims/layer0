package tag_store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func NewTestTagStore(t *testing.T) *DynamoTagStore {
	table := config.TestDynamoTableName()
	if table == "" {
		t.Skipf("Skipping test: %s not set", config.TEST_AWS_DYNAMO_TABLE)
	}

	creds := credentials.NewStaticCredentials(config.AWSAccessKey(), config.AWSSecretKey(), "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AWSRegion()),
	}

	session := session.New(awsConfig)
	store := NewDynamoTagStore(session, table)

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func TestDynamoTagStoreInsert(t *testing.T) {
	store := NewTestTagStore(t)

	tag := &models.Tag{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"}
	if err := store.Insert(tag); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoTagStoreDelete(t *testing.T) {
	store := NewTestTagStore(t)

	tag := &models.Tag{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"}
	if err := store.Insert(tag); err != nil {
		t.Fatal(err)
	}

	if err := store.Delete(tag.TagID); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoTagStoreSelectAll(t *testing.T) {
	store := NewTestTagStore(t)

	tags := []*models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"},
		{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env2"},
		{EntityID: "s1", EntityType: "service", Key: "name", Value: "svc1"},
		{EntityID: "s1", EntityType: "service", Key: "environment_id", Value: "env1"},
		{EntityID: "s2", EntityType: "service", Key: "name", Value: "svc2"},
		{EntityID: "s2", EntityType: "service", Key: "environment_id", Value: "env2"},
	}

	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	result, err := store.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Result: %#v", result[0])
	testutils.AssertEqual(t, len(result), len(tags))
}

func TestDynamoTagStoreSelectByQuery(t *testing.T) {
	store := NewTestTagStore(t)

	tags := models.Tags{
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

	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	cases := []struct {
		EntityID   string
		EntityType string
		Expected   models.Tags
	}{
		// filter by nothing
		{
			Expected: tags,
		},
		// filter by entityID
		{
			EntityID: "invalid",
			Expected: models.Tags{},
		},
		{
			EntityID: "e1",
			Expected: tags[0:1],
		},
		{
			EntityID: "e2",
			Expected: tags[1:2],
		},
		{
			EntityID: "s1",
			Expected: tags[2:4],
		},
		{
			EntityID: "s2",
			Expected: tags[4:6],
		},
		{
			EntityID: "l1",
			Expected: tags[6:8],
		},
		{
			EntityID: "l2",
			Expected: tags[8:10],
		},
		{
			EntityID: "d1",
			Expected: tags[10:12],
		},
		{
			EntityID: "d2",
			Expected: tags[12:14],
		},
		// filter by entityType
		{
			EntityType: "invalid",
			Expected:   models.Tags{},
		},
		{
			EntityType: "environment",
			Expected:   tags[0:2],
		},
		{
			EntityType: "service",
			Expected:   tags[2:6],
		},
		{
			EntityType: "load_balancer",
			Expected:   tags[6:10],
		},
		{
			EntityType: "deploy",
			Expected:   tags[10:14],
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
			Expected:   tags[0:1],
		},
		{
			EntityType: "service",
			EntityID:   "s1",
			Expected:   tags[2:4],
		},
		{
			EntityType: "load_balancer",
			EntityID:   "l1",
			Expected:   tags[6:8],
		},
		{
			EntityType: "deploy",
			EntityID:   "d1",
			Expected:   tags[10:12],
		},
	}

	for query, c := range cases {
		results, err := store.SelectByQuery(c.EntityType, c.EntityID)
		if err != nil {
			t.Fatalf("Query %d: %v", query, err)
		}

		if r, e := len(results), len(c.Expected); r != e {
			t.Fatalf("Query %d: expected %d tags, got %d", query, e, r)
		}

		for _, expected := range c.Expected {
			var isInResults bool

			for _, result := range results {
				if *expected == *result {
					isInResults = true
					break
				}
			}

			if !isInResults {
				t.Fatalf("Query %d: tag %#v is not in results", query, expected)
			}
		}
	}
}
