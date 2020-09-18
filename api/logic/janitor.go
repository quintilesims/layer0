package logic

import (
	"time"

	"github.com/quintilesims/layer0/common/db/job_store"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
	"github.com/quintilesims/layer0/common/waitutils"
)

const (
	JOB_LIFETIME           = time.Hour * 1
	JANITOR_SLEEP_DURATION = time.Minute * 10
)

var janitorLogger = logutils.NewStackTraceLogger("Janitor")

type Janitor struct {
	JobLogic  JobLogic
	TaskLogic TaskLogic
	TagStore  tag_store.TagStore
	JobStore  job_store.JobStore
	Clock     waitutils.Clock
}

func NewJanitor(jobLogic JobLogic, taskLogic TaskLogic, jobStore job_store.JobStore, tagStore tag_store.TagStore) *Janitor {
	return &Janitor{
		JobLogic:  jobLogic,
		TaskLogic: taskLogic,
		JobStore:  jobStore,
		TagStore:  tagStore,
		Clock:     waitutils.RealClock{},
	}
}

func (this *Janitor) Run() {
	go func() {
		for {
			janitorLogger.Info("Starting cleanup")
			this.pulse()
			janitorLogger.Infof("Finished cleanup")
			this.Clock.Sleep(JANITOR_SLEEP_DURATION)
		}
	}()
}

func getValidTags(jobs []*models.Job, tags *models.Tags) models.Tags {
	result := []models.Tag{}
	//loop job table
	for _, job := range jobs {
		//check meta task id
		if job.Meta != nil {
			if tags.Any(func(t models.Tag) bool {
				if t.EntityType == "task" && t.EntityID == job.Meta["task_id"] {
					result = append(result, t)
					return true
				}
				return false
			}) {
			}
		}
		//check task id
		if job.TaskID != "" {
			if tags.Any(func(t models.Tag) bool {
				if t.EntityType == "task" && t.EntityID == job.TaskID {
					result = append(result, t)
					return true
				}
				return false
			}) {
			}
		}
	}
	//loop tag table
	for _, tag := range *tags {
		//check job type
		if tag.EntityType == "job" && tag.Key == "task_id" {
			if tags.Any(func(t models.Tag) bool {
				if t.EntityType == "task" && t.EntityID == tag.Value {
					result = append(result, t)
					return true
				}
				return false
			}) {
			}
		}
		//check deploy type (job.xxx)
		if tag.EntityType == "deploy" && tag.Key == "job" {
			if tags.Any(func(t models.Tag) bool {
				if t.EntityType == "task" && t.Key == "deploy_id" && t.Value == tag.EntityID {
					result = append(result, t)
					return true
				}
				return false
			}) {
			}
		}

	}

	return result
}

func (this *Janitor) pulse() error {

	//clean jobs that is not running
	jobs, err := this.JobLogic.ListJobs()
	if err != nil {
		jobLogger.Errorf("Failed to list jobs: %v", err)
		return err
	}

	errs := []error{}

	for _, job := range jobs {
		timeSinceCreated := this.Clock.Since(job.TimeCreated)
		if job.JobStatus != int64(types.InProgress) || job.JobStatus != int64(types.Error) {
			if timeSinceCreated > JOB_LIFETIME || job.JobStatus == int64(types.Completed) {
				jobLogger.Infof("Deleting job '%s'", job.JobID)
				if err := this.JobLogic.Delete(job.JobID); err != nil {
					jobLogger.Errorf("Failed to delete job '%s': %v", job.JobID, err)
					errs = append(errs, err)
					if err := this.JobStore.Delete(job.JobID); err != nil {
						jobLogger.Errorf("Failed to delete job directly from dynamodb'%s': %v", job.JobID, err)
						errs = append(errs, err)
					}
				} else {
					jobLogger.Infof("Finished deleting job '%s'", job.JobID)
				}
			}
		}
	}

	// start clean tags table
	tasks, err := this.TaskLogic.ListTasks()
	if err != nil {
		tagLogger.Errorf("Failed to list tasks: %v", err)
		return err
	}

	taskExists := func(id string) bool {
		for _, t := range tasks {
			if t.TaskID == id {
				return true
			}
		}
		return false
	}

	//clean the task in db that not in ecs task
	tags, err := this.TagStore.SelectByType("task")
	if err != nil {
		tagLogger.Errorln("Could not query for tag store for task entity type - ", err.Error())
	}

	for _, tag := range tags {
		if !taskExists(tag.EntityID) {
			if err := this.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
				tagLogger.Errorf("Could not delete tag (%#v) -  %s\n", tag, err.Error())
				continue
			}

			tagLogger.Infof("Tag for task (%s) has been deleted\n", tag.EntityID)
		}
	}

	//clean up the tasks not linking to any running jobs
	tags, err = this.TagStore.SelectByType("task")
	jobTags, err := this.TagStore.SelectByType("job")
	deployTags, err := this.TagStore.SelectByType("deploy")
	tags = append(append(tags, jobTags...), deployTags...)
	keepTags := getValidTags(jobs, &tags)

	for _, tag := range tags {
		if !keepTags.Any(func(t models.Tag) bool {
			if t.EntityType == tag.EntityType && t.EntityID == tag.EntityID && t.Key == tag.Key {
				return true
			}
			return false
		}) {
			if err := this.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
				jobLogger.Errorf("Failed to delete tag '%s' %s : %v", tag.EntityType, tag.EntityID, err)
				errs = append(errs, err)
			}
		}

	}
	return errors.MultiError(errs)
}
