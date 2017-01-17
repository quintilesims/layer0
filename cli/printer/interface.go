package printer

import (
	"github.com/quintilesims/layer0/common/models"
)

type Printer interface {
	StartSpinner(message string)
	StopSpinner()
	PrintDeploys(deploys ...*models.Deploy) error
	PrintDeploySummaries(deploys ...*models.DeploySummary) error
	PrintEnvironments(environments ...*models.Environment) error
	PrintJobs(jobs ...*models.Job) error
	PrintLoadBalancers(loadBalancers ...*models.LoadBalancer) error
	PrintLogs(logs ...*models.LogFile) error
	PrintServices(services ...*models.Service) error
	PrintServiceSummaries(services ...*models.ServiceSummary) error
	PrintTasks(tasks ...*models.Task) error
	Printf(format string, tokens ...interface{})
	Fatalf(code int64, format string, tokens ...interface{})
}
