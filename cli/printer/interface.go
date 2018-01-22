package printer

import (
	"github.com/quintilesims/layer0/common/models"
)

type Printer interface {
	PrintDeploys(deploys ...*models.Deploy) error
	PrintDeploySummaries(deploys ...models.DeploySummary) error
	PrintEnvironments(environments ...*models.Environment) error
	PrintEnvironmentSummaries(environments ...models.EnvironmentSummary) error
	PrintLoadBalancers(loadBalancers ...*models.LoadBalancer) error
	PrintLoadBalancerSummaries(loadBalancers ...models.LoadBalancerSummary) error
	PrintLoadBalancerHealthCheck(loadBalancer *models.LoadBalancer) error
	PrintLogs(logs ...*models.LogFile) error
	PrintServices(services ...*models.Service) error
	PrintServiceSummaries(services ...models.ServiceSummary) error
	PrintTasks(tasks ...*models.Task) error
	PrintTaskSummaries(tasks ...models.TaskSummary) error
	Printf(format string, tokens ...interface{})
	Fatalf(code int64, format string, tokens ...interface{})
}
