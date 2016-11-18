package printer

import (
	"gitlab.imshealth.com/xfra/layer0/cli/entity"
	"gitlab.imshealth.com/xfra/layer0/common/models"
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
