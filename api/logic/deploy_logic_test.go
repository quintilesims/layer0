package logic

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetDeploy(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retDeploy := &models.Deploy{DeployID: "d1"}

	testLogic.Backend.EXPECT().
		GetDeploy("d1").
		Return(retDeploy, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "2"},
		{EntityID: "extra", EntityType: "deploy", Key: "name", Value: "extra"},
	})

	deployLogic := NewL0DeployLogic(testLogic.Logic())
	received, err := deployLogic.GetDeploy("d1")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Deploy{DeployID: "d1", DeployName: "dpl", Version: "2"}
	testutils.AssertEqual(t, received, expected)
}

func TestListDeploys(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retDeploys := []*models.Deploy{
		{DeployID: "d1"},
		{DeployID: "d2"},
	}

	testLogic.Backend.EXPECT().
		ListDeploys().
		Return(retDeploys, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl_1"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "2"},
		{EntityID: "d2", EntityType: "deploy", Key: "name", Value: "dpl_2"},
		{EntityID: "d2", EntityType: "deploy", Key: "version", Value: "3"},
		{EntityID: "extra", EntityType: "deploy", Key: "name", Value: "extra"},
	})

	deployLogic := NewL0DeployLogic(testLogic.Logic())
	received, err := deployLogic.ListDeploys()
	if err != nil {
		t.Fatal(err)
	}

	expected := []*models.DeploySummary{
		{DeployID: "d1", DeployName: "dpl_1", Version: "2"},
		{DeployID: "d2", DeployName: "dpl_2", Version: "3"},
	}

	testutils.AssertEqual(t, received, expected)
}

func TestDeleteDeploy(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		DeleteDeploy("d1").
		Return(nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "dpl"},
		{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "2"},
		{EntityID: "extra", EntityType: "deploy", Key: "name", Value: "extra"},
	})

	deployLogic := NewL0DeployLogic(testLogic.Logic())
	if err := deployLogic.DeleteDeploy("d1"); err != nil {
		t.Fatal(err)
	}

	tags, err := testLogic.TagStore.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCreateDeploy(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retDeploy := &models.Deploy{DeployID: "d1", Version: "1"}

	testLogic.Backend.EXPECT().
		CreateDeploy("name", []byte("dockerrun")).
		Return(retDeploy, nil)

	request := models.CreateDeployRequest{
		DeployName: "name",
		Dockerrun:  []byte("dockerrun"),
	}

	deployLogic := NewL0DeployLogic(testLogic.Logic())
	received, err := deployLogic.CreateDeploy(request)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Deploy{
		DeployID:   "d1",
		DeployName: "name",
		Version:    "1",
	}

	testutils.AssertEqual(t, received, expected)
	testLogic.AssertTagExists(t, models.Tag{EntityID: "d1", EntityType: "deploy", Key: "name", Value: "name"})
	testLogic.AssertTagExists(t, models.Tag{EntityID: "d1", EntityType: "deploy", Key: "version", Value: "1"})
}

func TestCreateDeployError_missingRequiredParams(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	deployLogic := NewL0DeployLogic(testLogic.Logic())

	cases := map[string]models.CreateDeployRequest{
		"Missing DeployName": models.CreateDeployRequest{},
	}

	for name, request := range cases {
		if _, err := deployLogic.CreateDeploy(request); err == nil {
			t.Errorf("Case %s: error was nil!", name)
		}
	}
}
