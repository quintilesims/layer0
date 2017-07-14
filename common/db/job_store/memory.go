package job_store

import (
	"fmt"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
)

type MemoryJobStore struct {
	jobs []*models.Job
}

func NewMemoryJobStore() *MemoryJobStore {
	return &MemoryJobStore{
		jobs: []*models.Job{},
	}
}

func (m *MemoryJobStore) Init() error {
	return nil
}

func (m *MemoryJobStore) Insert(job *models.Job) error {
	m.jobs = append(m.jobs, job)
	return nil
}

func (m *MemoryJobStore) Delete(jobID string) error {
	for i := 0; i < len(m.jobs); i++ {
		if m.jobs[i].JobID == jobID {
			m.jobs = append(m.jobs[:i], m.jobs[i+1:]...)
			i--
		}
	}

	return nil
}

func (m *MemoryJobStore) SelectAll() ([]*models.Job, error) {
	return m.jobs, nil
}

func (m *MemoryJobStore) SelectByID(jobID string) (*models.Job, error) {
	for _, job := range m.jobs {
		if job.JobID == jobID {
			return job, nil
		}
	}

	return nil, fmt.Errorf("Job with id '%d' does not exist", jobID)
}

func (m *MemoryJobStore) UpdateJobStatus(jobID string, status types.JobStatus) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.JobStatus = int64(status)
	return nil
}

func (m *MemoryJobStore) SetJobMeta(jobID string, meta map[string]string) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.Meta = meta
	return nil
}
