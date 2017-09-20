package job

import "github.com/quintilesims/layer0/common/models"

type Runner interface {
	Run(job models.Job) (string, error)
}

type RunnerFunc func(job models.Job) (string, error)

func (r RunnerFunc) Run(job models.Job) (string, error) {
	return r(job)
}
