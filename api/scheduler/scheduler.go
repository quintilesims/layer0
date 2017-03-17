package scheduler

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/logutils"
	"time"
)

const (
	SCALE_SLEEP_DURATION = time.Minute * 5
)

type Scheduler interface {
	Run()
}

type L0Scheduler struct {
	EnvironmentLogic  logic.EnvironmentLogic
	EnvironmentScaler EnvironmentScaler
	Logger            *logrus.Logger
}

func New(l logic.EnvironmentLogic, s EnvironmentScaler) *L0Scheduler {
	return &L0Scheduler{
		EnvironmentLogic:  l,
		EnvironmentScaler: s,
		Logger:            logutils.NewStandardLogger("Scheduler"),
	}
}

func (s *L0Scheduler) Run() {
	for {
		s.Logger.Infof("Scaling all Environments")
		if err := s.scaleEnvironments(); err != nil {
			s.Logger.Error(err)
		}

		time.Sleep(SCALE_SLEEP_DURATION)
	}
}

func (s *L0Scheduler) scaleEnvironments() error {
	environments, err := s.EnvironmentLogic.ListEnvironments()
	if err != nil {
		return fmt.Errorf("Failed to list environments: %v", err)
	}

	errs := []error{}
	for _, environment := range environments {
		s.Logger.Infof("Scaling Environment %s", environment.EnvironmentID)

		if _, err := s.EnvironmentScaler.Scale(environment.EnvironmentID); err != nil {
			err = fmt.Errorf("Failed to scale environment %s: %v", err)
			errs = append(errs, err)
			continue
		}
	}

	return errors.MultiError(errs)
}
