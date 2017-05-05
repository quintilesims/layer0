package tag_store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"testing"
	"time"
	"math/rand"
)

func TestBatchInsert(t *testing.T) {
	t.Skip("Table already populated")

	store := NewTestTagStore(t)
	entityGroups := []struct {
		EntityType string
		Count      int
		ID         func(i int) string
		Name       func(i int) string
	}{
		{
			EntityType: "environment",
			Count:      50,
			ID:         func(i int) string { return fmt.Sprintf("e%d", i) },
			Name:       func(i int) string { return fmt.Sprintf("env_%d", i) },
		},
		{
			EntityType: "load_balancer",
			Count:      250,
			ID:         func(i int) string { return fmt.Sprintf("l%d", i) },
			Name:       func(i int) string { return fmt.Sprintf("lbr_%d", i) },
		},
		{
			EntityType: "service",
			Count:      250,
			ID:         func(i int) string { return fmt.Sprintf("s%d", i) },
			Name:       func(i int) string { return fmt.Sprintf("svc_%d", i) },
		},
		{
			EntityType: "task",
			Count:      5000,
			ID:         func(i int) string { return fmt.Sprintf("t%d", i) },
			Name:       func(i int) string { return fmt.Sprintf("tsk_%d", i) },
		},
		{
			EntityType: "deploy",
			Count:      50000,
			ID:         func(i int) string { return fmt.Sprintf("d%d", i) },
			Name:       func(i int) string { return fmt.Sprintf("dpl_%d", i) },
		},
	}

	for _, eg := range entityGroups {
		for i := 0; i < eg.Count; i++ {
			tag := models.Tag{
				EntityType: eg.EntityType,
				EntityID:   eg.ID(i),
				Key:        "name",
				Value:      eg.Name(i),
			}

			fmt.Printf("Adding tag %v %v\n", eg.EntityType, i)
			if err := store.Insert(tag); err != nil {
				t.Fatal(err)
			}
		}
	}
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

	session := session.New(awsConfig)
	store := NewDynamoTagStore(session, table)

	if err := store.Clear(); err != nil {
		t.Fatal(err)
	}

	return store
}

func TestDynamoTagStoreInsert(t *testing.T) {
	t.Skip("skipping for vbatch stuff")

	store := NewTestTagStore(t)

	tags := []models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env1"},
		{EntityID: "e1", EntityType: "environment", Key: "k1", Value: "v1"},
	}

	for _, tag := range tags {
		if err := store.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDynamoTagStoreSelectByType(t *testing.T) {
	store := NewTestTagStore(t)

	entityGroups := []struct {
		EntityType string
		Count      int
	}{
		{
			EntityType: "environment",
			Count:      5,
		},
		{
			EntityType: "deploy",
			Count:      5000,
		},
		{
			EntityType: "load_balancer",
			Count:      25,
		},
		{
			EntityType: "service",
			Count:      25,
		},
		{
			EntityType: "task",
			Count:      500,
		},
	}

	for start := time.Now(); time.Since(start) < time.Minute*180; {
		sleep := time.Second*time.Duration(rand.Intn(10))
		fmt.Printf("Sleeping for %v\n", sleep)

		time.Sleep(sleep)

		for _, eg := range entityGroups {
			now := time.Now()
			 fmt.Printf("running %s query\n", eg.EntityType)
			_, err := store.SelectByType(eg.EntityType)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("query took %v\n", time.Since(now))
			/*
			if result, expected := len(tags), eg.Count; result != expected {
				t.Errorf("%s: Length was %d, expected %d", eg.EntityType, result, expected)
			}

			for _, tag := range tags {
				if result, expected := tag.EntityType, eg.EntityType; result != expected {
					t.Errorf("EntityType was '%s', expected '%s'", result, expected)
				}
			}
			*/
		}
	}
}
