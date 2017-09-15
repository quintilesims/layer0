package job

import (
	"fmt"
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type MemoryStore struct {
	jobs []*models.Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: []*models.Job{},
	}
}

func (m *MemoryStore) Insert(jobType JobType, req string) (string, error) {
	job := &models.Job{
		JobID:   fmt.Sprintf("%v", time.Now().UnixNano()),
		Type:    string(jobType),
		Request: req,
		Status:  string(Pending),
		Created: time.Now(),
		Result:    "",
	}

	m.jobs = append(m.jobs, job)
	return job.JobID, nil
}

func (m *MemoryStore) AcquireJob(jobID string) (bool, error) {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return false, err
	}

	if Status(job.Status) != Pending {
		return false, nil
	}

	job.Status = string(InProgress)
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

func (m *MemoryStore) SetJobStatus(jobID string, status Status) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.Status = string(status)
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
