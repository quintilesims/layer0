package job

import "github.com/quintilesims/layer0/common/models"

type Runner interface {
	Run(job models.Job) error
}

type RunnerFunc func(job models.Job) error

func (r RunnerFunc) Run(job models.Job) error {
	return r(job)
}
