package printer

import (
	"encoding/json"
	"fmt"
	"github.com/quintilesims/layer0/common/models"
	"os"
)

type JSONPrinter struct{}

func (j *JSONPrinter) StartSpinner(string) {}
func (j *JSONPrinter) StopSpinner()        {}

func (j *JSONPrinter) Printf(format string, tokens ...interface{}) {
	message := struct {
		Message string
	}{
		Message: fmt.Sprintf(format, tokens...),
	}

	if err := j.print(message); err != nil {
		fmt.Println(err)
	}
}

func (j *JSONPrinter) Fatalf(code int64, format string, tokens ...interface{}) {
	message := struct {
		Code    int64
		Message string
	}{
		Code:    code,
		Message: fmt.Sprintf(format, tokens...),
	}

	if err := j.print(message); err != nil {
		fmt.Println(err)
	}

	os.Exit(1)
}

func (j *JSONPrinter) print(obj interface{}) error {
	js, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(js))
	return nil
}

func (j *JSONPrinter) PrintDeploys(deploys ...*models.Deploy) error {
	return j.print(deploys)
}

func (j *JSONPrinter) PrintDeploySummaries(deploys ...*models.DeploySummary) error {
	return j.print(deploys)
}

func (j *JSONPrinter) PrintEnvironments(environments ...*models.Environment) error {
	return j.print(environments)
}

func (j *JSONPrinter) PrintEnvironmentSummaries(environments ...*models.EnvironmentSummary) error {
	return j.print(environments)
}

func (j *JSONPrinter) PrintJobs(jobs ...*models.Job) error {
	return j.print(jobs)
}

func (j *JSONPrinter) PrintLoadBalancers(loadBalancers ...*models.LoadBalancer) error {
	return j.print(loadBalancers)
}

func (j *JSONPrinter) PrintLogs(logs ...*models.LogFile) error {
	return j.print(logs)
}

func (j *JSONPrinter) PrintServices(services ...*models.Service) error {
	return j.print(services)
}

func (j *JSONPrinter) PrintTasks(tasks ...*models.Task) error {
	return j.print(tasks)
}
