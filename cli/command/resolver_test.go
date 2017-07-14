package command

import (
	"testing"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func fuzzParams(entityType, target string) map[string]string {
	return map[string]string{
		"type": entityType,
		"fuzz": target,
	}
}

func tagsWithIDs(ids ...string) []*models.EntityWithTags {
	tags := []*models.EntityWithTags{}
	for _, id := range ids {
		tags = append(tags, &models.EntityWithTags{EntityID: id})
	}

	return tags
}

func tagWithName(id, name string) *models.EntityWithTags {
	return &models.EntityWithTags{
		EntityID: id,
		Tags: models.Tags{
			{
				EntityID: id,
				Key:      "name",
				Value:    name,
			},
		},
	}
}

func TestResolveExactIDMatch(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	resolver := NewTagResolver(tc.Command().Client)

	tc.Client.EXPECT().
		SelectByQuery(fuzzParams("environment", "id")).
		Return(tagsWithIDs("id"), nil)

	ids, err := resolver.Resolve("environment", "id")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(ids), 1)
	testutils.AssertEqual(t, ids[0], "id")
}

func TestResolveExactNameMatch(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	resolver := NewTagResolver(tc.Command().Client)

	tag := tagWithName("id", "name")
	tc.Client.EXPECT().
		SelectByQuery(fuzzParams("environment", "name")).
		Return([]*models.EntityWithTags{tag}, nil)

	ids, err := resolver.Resolve("environment", "name")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(ids), 1)
	testutils.AssertEqual(t, ids[0], "id")
}

func TestResolveWildcardMatch(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	resolver := NewTagResolver(tc.Command().Client)

	tags := append(
		tagsWithIDs("nid"),
		tagWithName("id1", "name1"),
		tagWithName("id2", "name2"),
	)

	tc.Client.EXPECT().
		SelectByQuery(fuzzParams("environment", "n")).
		Return(tags, nil)

	ids, err := resolver.Resolve("environment", "n*")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(ids), 3)
	testutils.AssertEqual(t, ids[0], "nid")
	testutils.AssertEqual(t, ids[1], "id1")
	testutils.AssertEqual(t, ids[2], "id2")
}

func TestResolveEnvironmentScopedEntity(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	resolver := NewTagResolver(tc.Command().Client)

	tc.Client.EXPECT().
		SelectByQuery(fuzzParams("environment", "envid")).
		Return(tagsWithIDs("envid"), nil)

	params := fuzzParams("service", "svcid")
	params["environment_id"] = "envid"

	tc.Client.EXPECT().
		SelectByQuery(params).
		Return(tagsWithIDs("svcid"), nil)

	ids, err := resolver.Resolve("service", "envid:svcid")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(ids), 1)
	testutils.AssertEqual(t, ids[0], "svcid")
}

func TestResolveDeploy(t *testing.T) {
	tc, ctrl := newTestCommand(t)
	defer ctrl.Finish()
	resolver := NewTagResolver(tc.Command().Client)

	params := fuzzParams("deploy", "id")
	params["version"] = "version"

	tc.Client.EXPECT().
		SelectByQuery(params).
		Return(tagsWithIDs("id"), nil)

	ids, err := resolver.Resolve("deploy", "id:version")
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, len(ids), 1)
	testutils.AssertEqual(t, ids[0], "id")
}
