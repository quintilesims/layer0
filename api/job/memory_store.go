package job

import (
	"fmt"

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

func (m *MemoryStore) Init() error {
	return nil
}

func (m *MemoryStore) Insert(job models.Job) error {
	m.jobs = append(m.jobs, &job)
	return nil
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

	return nil, fmt.Errorf("Job with id '%d' does not exist", jobID)
}

func (m *MemoryStore) UpdateStatus(jobID string, status Status) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.JobStatus = string(status)
	return nil
}

func (m *MemoryStore) SetJobMeta(jobID string, meta map[string]string) error {
	job, err := m.SelectByID(jobID)
	if err != nil {
		return err
	}

	job.Meta = meta
	return nil
}
