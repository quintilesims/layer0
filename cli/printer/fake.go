package printer

import (
	"github.com/quintilesims/layer0/cli/entity"
	"github.com/quintilesims/layer0/common/models"
)

type FakePrinter struct{}

func (this *FakePrinter) StartSpinner(string) {}

func (this *FakePrinter) StopSpinner() {}

func (this *FakePrinter) PrintEntity(entity.Entity) error {
	return nil
}

func (this *FakePrinter) PrintEntities([]entity.Entity) error {
	return nil
}

func (this *FakePrinter) PrintLogs([]*models.LogFile) error {
	return nil
}

func (this *FakePrinter) Printf(string, ...interface{}) {}

func (this *FakePrinter) Fatalf(int64, string, ...interface{}) {}
