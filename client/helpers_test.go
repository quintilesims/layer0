package client

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func setTimeMultiplier(v time.Duration) func() {
	timeMultiplier = v
	return func() { timeMultiplier = 1 }
}

func TestWaitForJob(t *testing.T) {
	defer setTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	expected := &models.Job{
		JobID:  "jid",
		Status: job.Completed.String(),
	}

	client.EXPECT().
		ReadJob("jid").
		Return(expected, nil)

	result, err := WaitForJob(client, "jid", time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestWaitForJobError(t *testing.T) {
	defer setTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	expected := &models.Job{
		JobID:  "jid",
		Status: job.Error.String(),
	}

	client.EXPECT().
		ReadJob("jid").
		Return(expected, nil)

	if _, err := WaitForJob(client, "jid", time.Millisecond); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestWaitForJobTimeout(t *testing.T) {
	defer setTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	expected := &models.Job{
		JobID:  "jid",
		Status: job.InProgress.String(),
	}

	client.EXPECT().
		ReadJob("jid").
		Return(expected, nil).
		AnyTimes()

	if _, err := WaitForJob(client, "jid", time.Millisecond); err == nil {
		t.Fatal("Error was nil!")
	}
}
