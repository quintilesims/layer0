package aws

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/db/job_store/mock_job_store"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestJobRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobModel := &models.Job{}

	mockJobStore := mock_job_store.NewMockJobStore(ctrl)
	mockJobStore.EXPECT().
		SelectByID("j1").
		Return(jobModel, nil)

	job := NewJob(nil, mockJobStore, "j1")
	result, err := job.Model()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, jobModel, result)
}

func TestJobDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobStore := mock_job_store.NewMockJobStore(ctrl)
	mockJobStore.EXPECT().
		Delete("j1").
		Return(nil)

	job := NewJob(nil, mockJobStore, "j1")
	if err := job.Delete(); err != nil {
		t.Fatal(err)
	}
}
