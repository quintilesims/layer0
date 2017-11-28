package layer0

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceDeployCreateRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	req := models.CreateDeployRequest{
		DeployName: "dpl_name",
		DeployFile: []byte("some_content"),
	}

	mockClient.EXPECT().
		CreateDeploy(req).
		Return("job_id", nil)

	job := &models.Job{
		Status: job.Completed.String(),
		Result: "dpl_id",
	}

	mockClient.EXPECT().
		ReadJob("job_id").
		Return(job, nil)

	deploy := &models.Deploy{
		DeployID:   "dpl_id",
		DeployName: "dpl_name",
		Version:    "1",
		DeployFile: []byte("some_content"),
	}

	mockClient.EXPECT().
		ReadDeploy("dpl_id").
		Return(deploy, nil)

	deployResource := Provider().(*schema.Provider).ResourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{
		"name":    "dpl_name",
		"content": "some_content",
	})

	if err := resourceLayer0DeployCreate(d, mockClient); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "dpl_id", d.Id())
	assert.Equal(t, "dpl_name", d.Get("name").(string))
	assert.Equal(t, "1", d.Get("version").(string))
}

func TestResourceDeployDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_client.NewMockClient(ctrl)

	mockClient.EXPECT().
		DeleteDeploy("dpl_id").
		Return("job_id", nil)

	mockClient.EXPECT().
		ReadJob("job_id").
		Return(&models.Job{Status: job.Completed.String()}, nil)

	deployResource := Provider().(*schema.Provider).ResourcesMap["layer0_deploy"]
	d := schema.TestResourceDataRaw(t, deployResource.Schema, map[string]interface{}{})
	d.SetId("dpl_id")

	if err := resourceLayer0DeployDelete(d, mockClient); err != nil {
		t.Fatal(err)
	}
}
