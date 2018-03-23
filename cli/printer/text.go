package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/quintilesims/layer0/common/models"
	"github.com/ryanuber/columnize"
)

type TextPrinter struct{}

func (t *TextPrinter) Printf(format string, tokens ...interface{}) {
	fmt.Printf(format, tokens...)
}

func (t *TextPrinter) Println(tokens ...interface{}) {
	fmt.Println(tokens...)
}

func (t *TextPrinter) Fatalf(code int64, format string, tokens ...interface{}) {
	t.Printf(format, tokens...)
	fmt.Println()
	os.Exit(1)
}

func (t *TextPrinter) PrintDeploys(deploys ...*models.Deploy) error {
	getCompatibilities := func(d *models.Deploy, i int) string {
		if i > len(d.Compatibilities)-1 {
			return ""
		}

		return d.Compatibilities[i]
	}

	rows := []string{"DEPLOY ID | DEPLOY NAME | VERSION | COMPATIBILITIES"}
	for _, d := range deploys {
		row := fmt.Sprintf("%s | %s | %s | %s",
			d.DeployID,
			d.DeployName,
			d.Version,
			getCompatibilities(d, 0))

		rows = append(rows, row)

		// add the extra compatibility rows
		for i := 1; i < len(d.Compatibilities); i++ {
			row := fmt.Sprintf(" | | | %s", getCompatibilities(d, i))
			rows = append(rows, row)
		}
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintDeploySummaries(deploys ...models.DeploySummary) error {
	getCompatibilities := func(d models.DeploySummary, i int) string {
		if i > len(d.Compatibilities)-1 {
			return ""
		}

		return d.Compatibilities[i]
	}

	rows := []string{"DEPLOY ID | DEPLOY NAME | VERSION | COMPATIBILITIES"}
	for _, d := range deploys {
		row := fmt.Sprintf("%s | %s | %s | %s", d.DeployID, d.DeployName, d.Version, getCompatibilities(d, 0))
		rows = append(rows, row)

		// add the extra compatibility rows
		for i := 1; i < len(d.Compatibilities); i++ {
			row := fmt.Sprintf(" | | | %s", getCompatibilities(d, i))
			rows = append(rows, row)
		}
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintEnvironments(environments ...*models.Environment) error {
	getLink := func(e *models.Environment, i int) string {
		if i > len(e.Links)-1 {
			return ""
		}

		return e.Links[i]
	}

	rows := []string{"ENVIRONMENT ID | ENVIRONMENT NAME | OS | LINKS"}
	for _, e := range environments {
		row := fmt.Sprintf("%s | %s | %s | %s",
			e.EnvironmentID,
			e.EnvironmentName,
			e.OperatingSystem,
			getLink(e, 0))

		rows = append(rows, row)

		// add the extra link rows
		for i := 1; i < len(e.Links); i++ {
			row := fmt.Sprintf(" | | | %s", getLink(e, i))
			rows = append(rows, row)
		}
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintEnvironmentSummaries(environments ...models.EnvironmentSummary) error {
	rows := []string{"ENVIRONMENT ID | ENVIRONMENT NAME | OS "}
	for _, e := range environments {
		row := fmt.Sprintf("%s | %s | %s", e.EnvironmentID, e.EnvironmentName, e.OperatingSystem)
		rows = append(rows, row)
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintLoadBalancers(loadBalancers ...*models.LoadBalancer) error {
	getEnvironment := func(l *models.LoadBalancer) string {
		if l.EnvironmentName != "" {
			return l.EnvironmentName
		}

		return l.EnvironmentID
	}

	getService := func(l *models.LoadBalancer) string {
		if l.ServiceName != "" {
			return l.ServiceName
		}

		return l.ServiceID
	}

	getPort := func(l *models.LoadBalancer, i int) string {
		if i > len(l.Ports)-1 {
			return ""
		}

		p := l.Ports[i]
		return fmt.Sprintf("%d:%d/%s", p.HostPort, p.ContainerPort, strings.ToUpper(p.Protocol))
	}

	rows := []string{"LOADBALANCER ID | LOADBALANCER NAME | TYPE | ENVIRONMENT | SERVICE | PORTS | PUBLIC | URL "}
	for _, l := range loadBalancers {
		row := fmt.Sprintf("%s | %s | %s | %s | %s | %s | %t | %s",
			l.LoadBalancerID,
			l.LoadBalancerName,
			l.LoadBalancerType,
			getEnvironment(l),
			getService(l),
			getPort(l, 0),
			l.IsPublic,
			l.URL)

		rows = append(rows, row)

		// add the extra port rows
		for i := 1; i < len(l.Ports); i++ {
			row := fmt.Sprintf(" | | | | %s | |", getPort(l, i))
			rows = append(rows, row)
		}
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintLoadBalancerSummaries(loadBalancers ...models.LoadBalancerSummary) error {
	getEnvironment := func(l models.LoadBalancerSummary) string {
		if l.EnvironmentName != "" {
			return l.EnvironmentName
		}

		return l.EnvironmentID
	}

	rows := []string{"LOADBALANCER ID | LOADBALANCER NAME | TYPE | ENVIRONMENT"}
	for _, l := range loadBalancers {
		row := fmt.Sprintf("%s | %s | %s | %s ",
			l.LoadBalancerID,
			l.LoadBalancerName,
			l.LoadBalancerType,
			getEnvironment(l))

		rows = append(rows, row)
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintLoadBalancerHealthCheck(loadBalancer *models.LoadBalancer) error {
	getEnvironment := func(l *models.LoadBalancer) string {
		if l.EnvironmentName != "" {
			return l.EnvironmentName
		}

		return l.EnvironmentID
	}

	rows := []string{"LOADBALANCER ID | LOADBALANCER NAME | ENVIRONMENT | TARGET | INTERVAL | TIMEOUT | HEALTHY THRESHOLD | UNHEALTHY THRESHOLD "}
	row := fmt.Sprintf("%s | %s | %s | %s | %d | %d | %d | %d",
		loadBalancer.LoadBalancerID,
		loadBalancer.LoadBalancerName,
		getEnvironment(loadBalancer),
		loadBalancer.HealthCheck.Target,
		loadBalancer.HealthCheck.Interval,
		loadBalancer.HealthCheck.Timeout,
		loadBalancer.HealthCheck.HealthyThreshold,
		loadBalancer.HealthCheck.UnhealthyThreshold)

	rows = append(rows, row)

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintLogs(logs ...models.LogFile) error {
	for _, l := range logs {
		fmt.Println(l.ContainerName)
		for i := 0; i < len(l.ContainerName); i++ {
			fmt.Printf("-")
		}

		fmt.Println()
		for _, line := range l.Lines {
			fmt.Println(line)
		}
		fmt.Println()
	}

	return nil
}

func (t *TextPrinter) PrintServices(services ...*models.Service) error {
	getEnvironment := func(s *models.Service) string {
		if s.EnvironmentName != "" {
			return s.EnvironmentName
		}

		return s.EnvironmentID
	}

	getLoadBalancer := func(s *models.Service) string {
		if s.LoadBalancerName != "" {
			return s.LoadBalancerName
		}

		return s.LoadBalancerID
	}

	getDeployment := func(s *models.Service, i int) string {
		if i > len(s.Deployments)-1 {
			return ""
		}

		deployment := s.Deployments[i]
		display := strings.Replace(deployment.DeployID, ".", ":", 1)
		if deployment.DeployName != "" && deployment.DeployVersion != "" {
			display = fmt.Sprintf("%s:%s", deployment.DeployName, deployment.DeployVersion)
		}

		if deployment.RunningCount != deployment.DesiredCount {
			display = fmt.Sprintf("%s*", display)
		}

		return display
	}

	getScale := func(s *models.Service) string {
		scale := fmt.Sprintf("%d/%d", s.RunningCount, s.DesiredCount)
		if s.PendingCount != 0 {
			scale = fmt.Sprintf("%s (%d)", scale, s.PendingCount)
		}

		return scale
	}

	rows := []string{"SERVICE ID | SERVICE NAME | ENVIRONMENT | LOADBALANCER | DEPLOYMENTS | SCALE | STATEFUL"}
	for _, s := range services {
		row := fmt.Sprintf("%s | %s | %s | %s | %s | %s | %t",
			s.ServiceID,
			s.ServiceName,
			getEnvironment(s),
			getLoadBalancer(s),
			getDeployment(s, 0),
			getScale(s),
			s.Stateful)

		rows = append(rows, row)

		// add the extra deployment rows
		for i := 1; i < len(s.Deployments); i++ {
			row := fmt.Sprintf(" | | | | %s | ", getDeployment(s, i))
			rows = append(rows, row)
		}
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintServiceSummaries(services ...models.ServiceSummary) error {
	getEnvironment := func(s models.ServiceSummary) string {
		if s.EnvironmentName != "" {
			return s.EnvironmentName
		}

		return s.EnvironmentID
	}

	rows := []string{"SERVICE ID | SERVICE NAME | ENVIRONMENT"}
	for _, s := range services {
		row := fmt.Sprintf("%s | %s | %s ",
			s.ServiceID,
			s.ServiceName,
			getEnvironment(s))

		rows = append(rows, row)
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintTasks(tasks ...*models.Task) error {
	getEnvironment := func(t *models.Task) string {
		if t.EnvironmentName != "" {
			return t.EnvironmentName
		}

		return t.EnvironmentID
	}

	getDeploy := func(t *models.Task) string {
		if t.DeployName != "" && t.DeployVersion != "" {
			return fmt.Sprintf("%s:%s", t.DeployName, t.DeployVersion)
		}

		return strings.Replace(t.DeployID, ".", ":", 1)
	}

	rows := []string{"TASK ID | TASK NAME | ENVIRONMENT | DEPLOY | STATUS | STATEFUL "}
	for _, t := range tasks {
		row := fmt.Sprintf("%s | %s | %s | %s | %s | %t",
			t.TaskID,
			t.TaskName,
			getEnvironment(t),
			getDeploy(t),
			t.Status,
			t.Stateful)

		rows = append(rows, row)
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintTaskSummaries(tasks ...models.TaskSummary) error {
	getEnvironment := func(t models.TaskSummary) string {
		if t.EnvironmentName != "" {
			return t.EnvironmentName
		}

		return t.EnvironmentID
	}

	rows := []string{"TASK ID | TASK NAME | ENVIRONMENT"}
	for _, t := range tasks {
		row := fmt.Sprintf("%s | %s | %s",
			t.TaskID,
			t.TaskName,
			getEnvironment(t))

		rows = append(rows, row)
	}

	t.Println(columnize.SimpleFormat(rows))
	return nil
}
