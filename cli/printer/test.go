package printer

import (
	"github.com/quintilesims/layer0/common/models"
)

// we use a TestPrinter instead of a mock because gomock doesn't support
// using gomock.Any() for variadic functions
type TestPrinter struct{}

func (t *TestPrinter) StartSpinner(string)                                             {}
func (t *TestPrinter) StopSpinner()                                                    {}
func (t *TestPrinter) Printf(string, ...interface{})                                   {}
func (t *TestPrinter) Fatalf(int64, string, ...interface{})                            {}
func (t *TestPrinter) PrintDeploys(...*models.Deploy) error                            { return nil }
func (t *TestPrinter) PrintDeploySummaries(...*models.DeploySummary) error             { return nil }
func (t *TestPrinter) PrintEnvironments(...*models.Environment) error                  { return nil }
func (t *TestPrinter) PrintEnvironmentSummaries(...*models.EnvironmentSummary) error   { return nil }
func (t *TestPrinter) PrintJobs(...*models.Job) error                                  { return nil }
func (t *TestPrinter) PrintLoadBalancers(...*models.LoadBalancer) error                { return nil }
func (t *TestPrinter) PrintLoadBalancerSummaries(...*models.LoadBalancerSummary) error { return nil }
func (t *TestPrinter) PrintLoadBalancerHealthCheck(*models.LoadBalancer) error         { return nil }
func (t *TestPrinter) PrintLoadBalancerIdleTimeout(*models.LoadBalancer) error         { return nil }
func (t *TestPrinter) PrintLogs(...*models.LogFile) error                              { return nil }
func (t *TestPrinter) PrintScalerRunInfo(*models.ScalerRunInfo) error                  { return nil }
func (t *TestPrinter) PrintServices(...*models.Service) error                          { return nil }
func (t *TestPrinter) PrintServiceSummaries(...*models.ServiceSummary) error           { return nil }
func (t *TestPrinter) PrintTasks(...*models.Task) error                                { return nil }
func (t *TestPrinter) PrintTaskSummaries(...*models.TaskSummary) error                 { return nil }
