package logic

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
	"github.com/zpatrick/go-bytesize"
)

func newTestEnvironmentResourceGetter(t *testing.T) (*TestEnvironmentResourceGetter, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	crg := &TestEnvironmentResourceGetter{
		ServiceLogic: mock_logic.NewMockServiceLogic(ctrl),
		TaskLogic:    mock_logic.NewMockTaskLogic(ctrl),
		DeployLogic:  mock_logic.NewMockDeployLogic(ctrl),
		JobLogic:     mock_logic.NewMockJobLogic(ctrl),
	}

	return crg, ctrl
}

type TestEnvironmentResourceGetter struct {
	ServiceLogic *mock_logic.MockServiceLogic
	TaskLogic    *mock_logic.MockTaskLogic
	DeployLogic  *mock_logic.MockDeployLogic
	JobLogic     *mock_logic.MockJobLogic
}

func (c *TestEnvironmentResourceGetter) EnvironmentResourceGetter() *EnvironmentResourceGetter {
	return NewEnvironmentResourceGetter(c.ServiceLogic, c.TaskLogic, c.DeployLogic, c.JobLogic)
}

func requestToString(t *testing.T, r models.CreateTaskRequest) string {
	bytes, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	return string(bytes)
}

var deployWithOneContainer []byte = []byte(`
{
  "containerDefinitions": [
    {
      "name": "one",
      "memory": 500,
      "portMappings": [
        {
          "hostPort": 80,
          "containerPort": 80,
          "protocol": "tcp"
        },
        {
          "hostPort": 22,
          "containerPort": 22,
          "protocol": "tcp"
        }
      ]
    }
  ]
}
`)

var deployWithTwoContainers []byte = []byte(`
{
  "containerDefinitions": [
    {
      "name": "one",
      "memory": 500,
      "portMappings": [
        {
          "hostPort": 80,
          "containerPort": 80,
          "protocol": "tcp"
        }
      ]
    },
    {
      "name": "two",
      "memory": 1000,
      "portMappings": [
        {
          "hostPort": 8000,
          "containerPort": 8000,
          "protocol": "tcp"
        }
      ]
    }
  ]
}
`)

func TestGetPendingTaskResourcesInJobs(t *testing.T) {
	crg, ctrl := newTestEnvironmentResourceGetter(t)
	defer ctrl.Finish()

	jobs := []*models.Job{
		{
			JobID:     "j1",
			JobType:   int64(types.CreateTaskJob),
			JobStatus: int64(types.Pending),
			Request: requestToString(t, models.CreateTaskRequest{
				TaskName:      "t1",
				DeployID:      "d1",
				EnvironmentID: "e1",
				Copies:        2,
			}),
		},
		{
			JobID:     "j2",
			JobType:   int64(types.CreateTaskJob),
			JobStatus: int64(types.InProgress),
			Request: requestToString(t, models.CreateTaskRequest{
				TaskName:      "t2",
				DeployID:      "d2",
				EnvironmentID: "e1",
				Copies:        1,
			}),
		},
		{
			JobID:     "j3",
			JobType:   int64(types.CreateTaskJob),
			JobStatus: int64(types.InProgress),
			Request: requestToString(t, models.CreateTaskRequest{
				EnvironmentID: "e3",
			}),
		},
		{
			JobID:     "j4",
			JobType:   int64(types.DeleteEnvironmentJob),
			JobStatus: int64(types.InProgress),
		},
	}

	crg.JobLogic.EXPECT().
		ListJobs().
		Return(jobs, nil)

	crg.DeployLogic.EXPECT().
		GetDeploy("d1").
		Return(&models.Deploy{Dockerrun: deployWithOneContainer}, nil)

	crg.DeployLogic.EXPECT().
		GetDeploy("d2").
		Return(&models.Deploy{Dockerrun: deployWithTwoContainers}, nil)

	resources, err := crg.EnvironmentResourceGetter().getPendingTaskResourcesInJobs("e1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(resources), 4)

	// job1, deploy1, container1, copy1
	testutils.AssertEqual(t, resources[0].Ports, []int{80, 22})
	testutils.AssertEqual(t, resources[0].Memory, bytesize.MiB*500)

	// job1, deploy1, container1, copy2
	testutils.AssertEqual(t, resources[1].Ports, []int{80, 22})
	testutils.AssertEqual(t, resources[1].Memory, bytesize.MiB*500)

	// job2, deploy2, container1, copy1
	testutils.AssertEqual(t, resources[2].Ports, []int{80})
	testutils.AssertEqual(t, resources[2].Memory, bytesize.MiB*500)

	// job2, deploy2, container2, copy1
	testutils.AssertEqual(t, resources[3].Ports, []int{8000})
	testutils.AssertEqual(t, resources[3].Memory, bytesize.MiB*1000)
}

func TestGetPendingTaskResourcesInECS(t *testing.T) {
	crg, ctrl := newTestEnvironmentResourceGetter(t)
	defer ctrl.Finish()

	tasks := []*models.Task{
		{
			TaskID:        "t1",
			EnvironmentID: "e1",
			DeployID:      "d1",
			PendingCount:  2,
		},
		{
			TaskID:        "t2",
			EnvironmentID: "e1",
			DeployID:      "d2",
			PendingCount:  1,
		},
		{
			TaskID:        "t3",
			EnvironmentID: "e1",
			DeployID:      "d3",
			PendingCount:  0,
		},
	}

	crg.TaskLogic.EXPECT().
		GetEnvironmentTasks("e1").
		Return(tasks, nil)

	crg.DeployLogic.EXPECT().
		GetDeploy("d1").
		Return(&models.Deploy{Dockerrun: deployWithOneContainer}, nil)

	crg.DeployLogic.EXPECT().
		GetDeploy("d2").
		Return(&models.Deploy{Dockerrun: deployWithTwoContainers}, nil)

	resources, err := crg.EnvironmentResourceGetter().getPendingTaskResourcesInECS("e1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(resources), 4)

	// task1, deploy1, container1, copy1
	testutils.AssertEqual(t, resources[0].Ports, []int{80, 22})
	testutils.AssertEqual(t, resources[0].Memory, bytesize.MiB*500)

	// task1, deploy1, container1, copy2
	testutils.AssertEqual(t, resources[1].Ports, []int{80, 22})
	testutils.AssertEqual(t, resources[1].Memory, bytesize.MiB*500)

	// task2, deploy2, container1, copy1
	testutils.AssertEqual(t, resources[2].Ports, []int{80})
	testutils.AssertEqual(t, resources[2].Memory, bytesize.MiB*500)

	// task2, deploy2, container2, copy1
	testutils.AssertEqual(t, resources[3].Ports, []int{8000})
	testutils.AssertEqual(t, resources[3].Memory, bytesize.MiB*1000)
}

func TestGetPendingServiceResources(t *testing.T) {
	crg, ctrl := newTestEnvironmentResourceGetter(t)
	defer ctrl.Finish()

	services := []*models.Service{
		{
			ServiceID:     "s1",
			EnvironmentID: "e1",
			Deployments: []models.Deployment{
				{
					DeployID:     "d1",
					DesiredCount: 2,
					RunningCount: 0,
				},
			},
		},
		{
			ServiceID:     "s2",
			EnvironmentID: "e1",
			Deployments: []models.Deployment{
				{
					DeployID:     "d2",
					DesiredCount: 1,
					RunningCount: 0,
				},
			},
		},
		{
			ServiceID:     "s3",
			EnvironmentID: "e1",
			Deployments: []models.Deployment{
				{
					DeployID:     "d1",
					DesiredCount: 1,
					RunningCount: 1,
				},
			},
		},
	}

	crg.ServiceLogic.EXPECT().
		GetEnvironmentServices("e1").
		Return(services, nil)

	crg.DeployLogic.EXPECT().
		GetDeploy("d1").
		Return(&models.Deploy{Dockerrun: deployWithOneContainer}, nil)

	crg.DeployLogic.EXPECT().
		GetDeploy("d2").
		Return(&models.Deploy{Dockerrun: deployWithTwoContainers}, nil)

	resources, err := crg.EnvironmentResourceGetter().getPendingServiceResources("e1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(resources), 4)

	// service1, deploy1, container1, copy1
	testutils.AssertEqual(t, resources[0].Ports, []int{80, 22})
	testutils.AssertEqual(t, resources[0].Memory, bytesize.MiB*500)

	// service1, deploy1, container1, copy2
	testutils.AssertEqual(t, resources[1].Ports, []int{80, 22})
	testutils.AssertEqual(t, resources[1].Memory, bytesize.MiB*500)

	// service2, deploy2, container1, copy1
	testutils.AssertEqual(t, resources[2].Ports, []int{80})
	testutils.AssertEqual(t, resources[2].Memory, bytesize.MiB*500)

	// service2, deploy2, container2, copy1
	testutils.AssertEqual(t, resources[3].Ports, []int{8000})
	testutils.AssertEqual(t, resources[3].Memory, bytesize.MiB*1000)
}
