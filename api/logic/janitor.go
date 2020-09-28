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
	JANITOR_SLEEP_DURATION = time.Minute * 5
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
		if job.JobStatus != int64(types.InProgress) {
			if timeSinceCreated > JOB_LIFETIME && job.JobStatus == int64(types.Completed) {
				jobLogger.Infof("Deleting job '%s'", job.JobID)
				if err := this.JobLogic.Delete(job.JobID); err != nil {
					jobLogger.Errorf("Failed to delete job '%s': %v", job.JobID, err)
					errs = append(errs, err)
					//If by any reasons the job can not be deleted (data corruption), because it already pass the Job life time, just delete the item from dynamodb.
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
	taskTags, err := this.TagStore.SelectByType("task")
	if err != nil {
		tagLogger.Errorln("Could not query for tag store for task entity type - ", err.Error())
	}

	for _, tag := range taskTags {
		if !taskExists(tag.EntityID) {
			if err := this.TagStore.Delete(tag.EntityType, tag.EntityID, tag.Key); err != nil {
				tagLogger.Errorf("Could not delete tag (%#v) -  %s\n", tag, err.Error())
				continue
			}
			//to do
			tagLogger.Infof("Tag for task (%s) has been deleted\n", tag.EntityID)
		}
	}

	//clean up orphan api deploy
	deployTags, err := this.TagStore.SelectByType("deploy")
	for _, dtag := range deployTags.WithValue("job") {
		if !taskTags.WithKey("deploy_id").Any(func(t models.Tag) bool {
			if t.Value == dtag.EntityID {
				return true
			}
			return false
		}) {
			if err := this.TagStore.Delete("deploy", dtag.EntityID, "name"); err != nil {
				jobLogger.Errorf("Failed to delete tag '%s' %s : %v", dtag.EntityType, dtag.EntityID, err)
				errs = append(errs, err)
			}
			tagLogger.Infof("orphan deploy record (%s) has been deleted\n", dtag.EntityID)
		}
	}
	return errors.MultiError(errs)
}
