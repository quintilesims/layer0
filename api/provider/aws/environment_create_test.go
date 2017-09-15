package aws

import (
	"encoding/base64"
	"testing"

	"github.com/quintilesims/layer0/api/tag"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_createTags(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	environment := NewEnvironmentProvider(nil, tagStore, nil)

	if err := environment.createTags("env_id", "env_name", "env_os"); err != nil {
		t.Fatal(err)
	}

	expectedTags := models.Tags{
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "name",
			Value:      "env_name",
		},
		{
			EntityID:   "env_id",
			EntityType: "environment",
			Key:        "os",
			Value:      "env_os",
		},
	}

	for _, tag := range expectedTags {
		assert.Contains(t, tagStore.Tags(), tag)
	}
}

func TestEnvironment_renderUserData(t *testing.T) {
	template := "{{ .ECSEnvironmentID }} {{ .S3Bucket }}"

	encoded, err := renderUserData("env_id", "bucket", []byte(template))
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "env_id bucket", string(decoded))
}
