package printer

import (
	"encoding/json"
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/cli/entity"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"os"
)

type JSONPrinter struct{}

// don't use spinner for json output
func (this *JSONPrinter) StartSpinner(string) {}
func (this *JSONPrinter) StopSpinner()        {}

func (this *JSONPrinter) PrintEntity(e entity.Entity) error {
	return this.PrintEntities([]entity.Entity{e})
}

func (this *JSONPrinter) PrintEntities(entities []entity.Entity) error {
	js, err := json.MarshalIndent(entities, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(js))
	return nil
}

func (this *JSONPrinter) PrintLogs(logFiles []*models.LogFile) error {
	js, err := json.MarshalIndent(logFiles, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(js))
	return nil
}

type basicMessage struct {
	Message string
}

func (this *JSONPrinter) Printf(format string, tokens ...interface{}) {
	message := basicMessage{
		Message: fmt.Sprintf(format, tokens...),
	}

	this.printf(message)
}

type errorMessage struct {
	Message string
	Code    int64
}

func (this *JSONPrinter) Fatalf(code int64, format string, tokens ...interface{}) {
	message := errorMessage{
		Message: fmt.Sprintf(format, tokens...),
		Code:    code,
	}

	this.printf(message)
	os.Exit(1)
}

func (this *JSONPrinter) printf(output interface{}) {
	js, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		js = []byte(err.Error())
	}

	fmt.Println(string(js))
}
