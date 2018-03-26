package provider

import (
	"time"

	"github.com/quintilesims/layer0/common/models"
)

type AdminProvider interface {
	Logs(tail int, start, end time.Time) ([]models.LogFile, error)
}
