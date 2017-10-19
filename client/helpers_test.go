package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/client/mock_client"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestWaitForDeployment(t *testing.T) {
	defer SetTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	deployments := []models.Deployment{
		{
			DesiredCount: 1,
			RunningCount: 1,
		},
	}

	expected := &models.Service{
		Deployments:  deployments,
		DesiredCount: 1,
		RunningCount: 1,
	}

	client.EXPECT().
		ReadService("svc_id").
		Return(expected, nil).
		AnyTimes()

	result, err := WaitForDeployment(client, "svc_id", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestWaitForDeployment_error(t *testing.T) {
	defer SetTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	expected := fmt.Errorf("Error reading service")

	client.EXPECT().
		ReadService("svc_id").
		Return(nil, expected)

	if _, err := WaitForDeployment(client, "svc_id", time.Second); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestWaitForDeployment_timeout(t *testing.T) {
	defer SetTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	deployments := []models.Deployment{
		{
			DesiredCount: 1,
			RunningCount: 0,
		},
	}

	expected := &models.Service{
		Deployments: deployments,
	}

	client.EXPECT().
		ReadService("svc_id").
		AnyTimes().
		Return(expected, nil)

	if _, err := WaitForDeployment(client, "svc_id", time.Second); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestWaitForJob(t *testing.T) {
	defer SetTimeMultiplier(0)()

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
	defer SetTimeMultiplier(0)()

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
	defer SetTimeMultiplier(0)()

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
