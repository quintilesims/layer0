package aws

/*
func TestTask_populateTaskSummaries(t *testing.T) {
	tagStore := tag.NewMemoryStore()
	task := NewTaskProvider(nil, tagStore, nil)

	tags := models.Tags{
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk1",
		},
		{
			EntityID:   "t1",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "e1",
		},
		{
			EntityID:   "e1",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename1",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "name",
			Value:      "tsk2",
		},
		{
			EntityID:   "t2",
			EntityType: "task",
			Key:        "environment_id",
			Value:      "e2",
		},
		{
			EntityID:   "e2",
			EntityType: "environment",
			Key:        "name",
			Value:      "ename2",
		},
		{
			EntityID:   "someid",
			EntityType: "task",
			Key:        "name",
			Value:      "badname",
		},
		{
			EntityID:   "someid",
			EntityType: "task",
			Key:        "bad_env_key",
			Value:      "env1",
		},
		{
			EntityID:   "env1",
			EntityType: "service",
			Key:        "name",
			Value:      "servicename",
		},
	}

	for _, tag := range tags {
		if err := tagStore.Insert(tag); err != nil {
			t.Fatal(err)
		}
	}

	results := []models.TaskSummary{
		{TaskID: "t1"},
		{TaskID: "t2"},
	}

	if err := task.populateTaskSummaries(results); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, results, 2)
	assert.Equal(t, "tsk1", results[0].TaskName)
	assert.Equal(t, "e1", results[0].EnvironmentID)
	assert.Equal(t, "ename1", results[0].EnvironmentName)
	assert.Equal(t, "tsk2", results[1].TaskName)
	assert.Equal(t, "e2", results[1].EnvironmentID)
	assert.Equal(t, "ename2", results[1].EnvironmentName)
}
*/
