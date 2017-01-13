package printer

import (
	"encoding/json"
	"fmt"
	"github.com/quintilesims/layer0/common/models"
	"os"
)

type JSONPrinter struct{}

// don't use spinner for json output
func (j *JSONPrinter) StartSpinner(string) {}
func (j *JSONPrinter) StopSpinner()        {}

func (j *JSONPrinter) PrintDeploys(deploys ...*models.Deploy) error {
	return fmt.Errorf("Print not implemented")
}
func (j *JSONPrinter) PrintDeploySummaries(deploys ...*models.DeploySummary) error {
	return fmt.Errorf("Print not implemented")
}

func (j *JSONPrinter) PrintEnvironments(environments ...*models.Environment) error {
	return fmt.Errorf("Print not implemented")
}

func (j *JSONPrinter) PrintJobs(jobs ...*models.Job) error { return fmt.Errorf("Print not implemented") }

func (j *JSONPrinter) PrintLoadBalancers(loadBalancers ...*models.LoadBalancer) error {
	return fmt.Errorf("Print not implemented")
}

func (j *JSONPrinter) PrintLogs(logs ...*models.LogFile) error {
	js, err := json.MarshalIndent(logs, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(js))
	return nil
}

func (j *JSONPrinter) PrintServices(services ...*models.Service) error {
	return fmt.Errorf("Print not ipmelmtled")
}

func (j *JSONPrinter) PrintTasks(tasks ...*models.Task) error {
	return fmt.Errorf("Print not ipmelmtled")
}

type basicMessage struct {
	Message string
}

func (j *JSONPrinter) Printf(format string, tokens ...interface{}) {
	message := basicMessage{
		Message: fmt.Sprintf(format, tokens...),
	}

	j.printf(message)
}

type errorMessage struct {
	Message string
	Code    int64
}

func (j *JSONPrinter) Fatalf(code int64, format string, tokens ...interface{}) {
	message := errorMessage{
		Message: fmt.Sprintf(format, tokens...),
		Code:    code,
	}

	j.printf(message)
	os.Exit(1)
}

func (j *JSONPrinter) printf(output interface{}) {
	js, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		js = []byte(err.Error())
	}

	fmt.Println(string(js))
}
