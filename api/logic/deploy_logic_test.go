package logic

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetDeploy(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		GetDeploy("d1").
		Return(&models.Deploy{DeployID: "d1"}, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "2"},
	})

	deployLogic := NewL0DeployLogic(testLogic.Logic())
	deploy, err := deployLogic.GetDeploy("d1")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, deploy.DeployID, "d1")
	testutils.AssertEqual(t, deploy.DeployName, "dpl")
	testutils.AssertEqual(t, deploy.Version, "2")
}

func TestListDeploys(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	expected := []*models.Deploy{
		{DeployID: "d1"},
		{DeployID: "d2"},
	}

	testLogic.Backend.EXPECT().
		ListDeploys().
		Return(expected, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "2"},
		{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl_2"},
		{EntityID: "d2", EntityType: "deploy", Key: "version", Value: "3"},
	})

	deployLogic := NewL0DeployLogic(testLogic.Logic())
	deploys, err := deployLogic.ListDeploys()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(deploys), len(expected))
	testutils.AssertEqual(t, deploys[0].DeployID, "d1")
	testutils.AssertEqual(t, deploys[0].DeployName, "dpl_1")
	testutils.AssertEqual(t, deploys[0].Version, "2")
	testutils.AssertEqual(t, deploys[1].DeployID, "d2")
        testutils.AssertEqual(t, deploys[1].DeployName, "dpl_2")
        testutils.AssertEqual(t, deploys[1].Version, "3")
}
