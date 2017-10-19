package command

import (
	"fmt"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func TestWaitForDeployment(t *testing.T) {
	defer SetTimeMultiplier(0)()

	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

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

	base.Client.EXPECT().
		ReadService("svc_id").
		Times(4).
		Return(expected, nil)

	result, err := WaitForDeployment(base.Client, "svc_id", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestWaitForDeployment_error(t *testing.T) {
	defer SetTimeMultiplier(0)()

	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	expected := fmt.Errorf("Error reading service")

	base.Client.EXPECT().
		ReadService("svc_id").
		Return(nil, expected)

	if _, err := WaitForDeployment(base.Client, "svc_id", time.Second); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestWaitForDeployment_timeout(t *testing.T) {
	defer SetTimeMultiplier(0)()

	base, ctrl := newTestCommand(t)
	defer ctrl.Finish()

	deployments := []models.Deployment{
		{
			DesiredCount: 1,
			RunningCount: 0,
		},
	}

	expected := &models.Service{
		Deployments: deployments,
	}

	base.Client.EXPECT().
		ReadService("svc_id").
		AnyTimes().
		Return(expected, nil)

	if _, err := WaitForDeployment(base.Client, "svc_id", time.Second); err == nil {
		t.Fatal("Error was nil!")
	}
}
