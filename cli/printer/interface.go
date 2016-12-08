package printer

import (
	"github.com/quintilesims/layer0/cli/entity"
	"github.com/quintilesims/layer0/common/models"
)

type Printer interface {
	StartSpinner(string)
	StopSpinner()
	PrintEntity(entity.Entity) error
	PrintEntities([]entity.Entity) error
	PrintLogs([]*models.LogFile) error
	Printf(string, ...interface{})
	Fatalf(int64, string, ...interface{})
}
