package aws

import (
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
)

func TestService_deleteService(t *testing.T) {

	tagStore := tag.NewMemoryStore()
	serviceProvider := NewServiceProvider(nil, tagStore, nil)
	_ = serviceProvider

	tags := models.Tags{
		{
			EntityID:   "sID",
			EntityType: "service",
			Key:        "name",
			Value:      "sName",
		},
		{
			EntityID:   "eID",
			EntityType: "environment",
			Key:        "name",
			Value:      "eName",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	// TODO: Add delete service test.
}
