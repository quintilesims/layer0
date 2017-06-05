package logic

import (
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/api/logic/mock_logic"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/common/types"
	"testing"
)

func TestJobSelectByID(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddJobs(t, []*models.Job{
		{JobID: "j1"},
	})

	jobLogic := NewL0JobLogic(testLogic.Logic(), nil, nil)
	job, err := jobLogic.GetJob("j1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, job.JobID, "j1")
}

func TestJobSelectAll(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddJobs(t, []*models.Job{
		{JobID: "j1"},
		{JobID: "j2"},
	})

	jobLogic := NewL0JobLogic(testLogic.Logic(), nil, nil)
	jobs, err := jobLogic.ListJobs()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(jobs), 2)
	testutils.AssertEqual(t, jobs[0].JobID, "j1")
	testutils.AssertEqual(t, jobs[1].JobID, "j2")
}

func TestJobDelete(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	taskLogic := mock_logic.NewMockTaskLogic(ctrl)
	deployLogic := mock_logic.NewMockDeployLogic(ctrl)
	defer ctrl.Finish()

	testLogic.AddJobs(t, []*models.Job{
		{JobID: "j1", TaskID: "t1"},
		{JobID: "extra"},
	})

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "j1", EntityType: "job", Key: "task_id", Value: "t1"},
		{EntityID: "extra", EntityType: "job", Key: "name", Value: "extra"},
	})

	taskLogic.EXPECT().
		DeleteTask("t1").
		Return(nil)

	jobLogic := NewL0JobLogic(testLogic.Logic(), taskLogic, deployLogic)
	if err := jobLogic.Delete("j1"); err != nil {
		t.Fatal(err)
	}

	tags, err := testLogic.TagStore.SelectByType("job")
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)

	jobs, err := testLogic.JobStore.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' job is the only one left
	testutils.AssertEqual(t, len(jobs), 1)
}

func TestCreateJob(t *testing.T) {
	tmp := id.GenerateHashedEntityID
	id.GenerateHashedEntityID = func(name string) string { return "j1" }
	defer func() { id.GenerateHashedEntityID = tmp }()

	testLogic, ctrl := NewTestLogic(t)
	taskLogic := mock_logic.NewMockTaskLogic(ctrl)
	deployLogic := mock_logic.NewMockDeployLogic(ctrl)
	defer ctrl.Finish()

	deployLogic.EXPECT().
		CreateDeploy(gomock.Any()).
		Return(&models.Deploy{DeployID: "d1"}, nil)

	taskLogic.EXPECT().
		CreateTask(gomock.Any()).
		Return(&models.Task{TaskID: "t1"}, nil)

	jobLogic := NewL0JobLogic(testLogic.Logic(), taskLogic, deployLogic)
	job, err := jobLogic.CreateJob(types.DeleteEnvironmentJob, "e1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, job.TaskID, "t1")
	testLogic.AssertTagExists(t, models.Tag{EntityID: "j1", EntityType: "job", Key: "task_id", Value: "t1"})

	if _, err := testLogic.JobStore.SelectByID("j1"); err != nil {
		t.Fatal(err)
	}
}
