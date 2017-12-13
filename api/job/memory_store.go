package job

import (
	"fmt"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type MemoryStore struct {
	jobs       []*models.Job
	insertHook func(jobID string)
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs:       []*models.Job{},
		insertHook: func(string) {},
	}
}

func (m *MemoryStore) Jobs() []*models.Job {
	return m.jobs
}

func (m *MemoryStore) Insert(jobType models.JobType, req string) (string, error) {
	job := &models.Job{
		JobID:   fmt.Sprintf("%v", time.Now().UnixNano()),
		Type:    jobType,
		Request: req,
		Status:  models.PendingJobStatus,
		Created: time.Now(),
		Result:  "",
	}

	m.jobs = append(m.jobs, job)
	m.insertHook(job.JobID)
	return job.JobID, nil
}

func (m *MemoryStore) SetInsertHook(hook func(jobID string)) {
	m.insertHook = hook
}

func (m *MemoryStore) AcquireJob(jobID string) (bool, error) {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return false, err
	}

	if job.Status != models.PendingJobStatus {
		return false, nil
	}

	job.Status = models.InProgressJobStatus
	return true, nil
}

func (m *MemoryStore) Delete(jobID string) error {
	for i := 0; i < len(m.jobs); i++ {
		if m.jobs[i].JobID == jobID {
			m.jobs = append(m.jobs[:i], m.jobs[i+1:]...)
			i--
		}
	}

	return nil
}

func (m *MemoryStore) SelectAll() ([]*models.Job, error) {
	return m.jobs, nil
}

func (m *MemoryStore) SelectByID(jobID string) (*models.Job, error) {
	for _, job := range m.jobs {
		if job.JobID == jobID {
			return job, nil
		}
	}

	return nil, fmt.Errorf("Job with id '%s' does not exist", jobID)
}

func (m *MemoryStore) SetJobStatus(jobID string, status models.JobStatus) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.Status = status
	return nil
}

func (m *MemoryStore) SetJobResult(jobID, result string) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.Result = result
	return nil
}

func (m *MemoryStore) SetJobError(jobID string, err error) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.Error = err.Error()
	return nil
}
