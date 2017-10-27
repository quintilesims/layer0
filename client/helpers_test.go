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

func TestWaitForDeployment_standard(t *testing.T) {
	defer SetTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	pending := &models.Service{
		Deployments: []models.Deployment{
			{
				DesiredCount: 1,
				RunningCount: 0,
			},
		},
	}

	finished := &models.Service{
		Deployments: []models.Deployment{
			{
				DesiredCount: 1,
				RunningCount: 1,
			},
		},
	}

	gomock.InOrder(
		// simulate usual lengthy run
		client.EXPECT().
			ReadService("svc_id").
			Return(pending, nil).
			Times(3),
		// simulate flapping service
		client.EXPECT().
			ReadService("svc_id").
			Return(finished, nil).
			Times(2),
		client.EXPECT().
			ReadService("svc_id").
			Return(pending, nil).
			Times(2),
		// simulate finished deployment
		client.EXPECT().
			ReadService("svc_id").
			Return(finished, nil).
			AnyTimes(),
	)

	result, err := WaitForDeployment(client, "svc_id", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, finished, result)
}

func TestWaitForDeployment_error(t *testing.T) {
	defer SetTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_client.NewMockClient(ctrl)

	client.EXPECT().
		ReadService("svc_id").
		Return(nil, fmt.Errorf("Error reading service"))

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
		Return(expected, nil).
		AnyTimes()

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

func TestContains(t *testing.T) {
	tests := []struct {
		InputValue string
		InputSlice []string
		Output     bool
	}{
		{
			InputValue: "a",
			InputSlice: []string{"a", "b", "c"},
			Output:     true,
		},
		{
			InputValue: "c",
			InputSlice: []string{"a", "b", "c"},
			Output:     true,
		},
		{
			InputValue: "d",
			InputSlice: []string{"a", "b", "c"},
			Output:     false,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("is \"%s\" in %q", tc.InputValue, tc.InputSlice), func(t *testing.T) {
			result := Contains(tc.InputValue, tc.InputSlice)
			if result != tc.Output {
				t.Fatalf("Expected %v but got %v", tc.Output, result)
			}
		})
	}
}
