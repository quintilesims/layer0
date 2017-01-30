package printer

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/types"
	"github.com/ryanuber/columnize"
	"os"
	"strings"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

type TextPrinter struct {
	spinner *spinner.Spinner
}

func (t *TextPrinter) StartSpinner(prefix string) {
	if t.spinner != nil {
		t.spinner.Stop()
	}

	t.spinner = spinner.New(spinner.CharSets[26], 1*time.Second)
	t.spinner.Prefix = prefix
	t.spinner.Start()
}

func (t *TextPrinter) StopSpinner() {
	if t.spinner != nil {
		t.spinner.Stop()
		fmt.Println()
	}
}

func (t *TextPrinter) Printf(format string, tokens ...interface{}) {
	t.StopSpinner()
	fmt.Printf(format, tokens...)
}

func (t *TextPrinter) Fatalf(code int64, format string, tokens ...interface{}) {
	t.Printf(format, tokens...)
	fmt.Println()
	os.Exit(1)
}

func (t *TextPrinter) PrintDeploys(deploys ...*models.Deploy) error {
	rows := []string{"DEPLOY ID | DEPLOY NAME | VERSION"}
	for _, d := range deploys {
		row := fmt.Sprintf("%s | %s |  %s", d.DeployID, d.DeployName, d.Version)
		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintDeploySummaries(deploys ...*models.DeploySummary) error {
	rows := []string{"DEPLOY ID | DEPLOY NAME | VERSION"}
	for _, d := range deploys {
		row := fmt.Sprintf("%s | %s |  %s", d.DeployID, d.DeployName, d.Version)
		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintEnvironments(environments ...*models.Environment) error {
	rows := []string{"ENVIRONMENT ID | ENVIRONMENT NAME | CLUSTER COUNT | INSTANCE SIZE"}
	for _, e := range environments {
		row := fmt.Sprintf("%s | %s | %d | %s",
			e.EnvironmentID,
			e.EnvironmentName,
			e.ClusterCount,
			e.InstanceSize)

		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintEnvironmentSummaries(environments ...*models.EnvironmentSummary) error {
	rows := []string{"ENVIRONMENT ID | ENVIRONMENT NAME"}
	for _, e := range environments {
		row := fmt.Sprintf("%s | %s", e.EnvironmentID, e.EnvironmentName)
		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintJobs(jobs ...*models.Job) error {
	getType := func(j *models.Job) string {
		jobType := types.JobType(j.JobType).String()
		return strings.Title(jobType)
	}

	getStatus := func(j *models.Job) string {
		jobStatus := types.JobStatus(j.JobStatus).String()
		return strings.Title(jobStatus)
	}

	rows := []string{"JOB ID | TASK ID | TYPE | STATUS | CREATED"}
	for _, j := range jobs {
		row := fmt.Sprintf("%s | %s | %s | %s | %s",
			j.JobID,
			j.TaskID,
			getType(j),
			getStatus(j),
			j.TimeCreated.Format(TIME_FORMAT))

		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
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

	rows := []string{"LOADBALANCER ID | LOADBALANCER NAME | ENVIRONMENT | SERVICE | PORTS | PUBLIC | URL "}
	for _, l := range loadBalancers {
		row := fmt.Sprintf("%s | %s | %s | %s | %s | %t | %s",
			l.LoadBalancerID,
			l.LoadBalancerName,
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

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintLoadBalancerSummaries(loadBalancers ...*models.LoadBalancerSummary) error {
	getEnvironment := func(l *models.LoadBalancerSummary) string {
		if l.EnvironmentName != "" {
			return l.EnvironmentName
		}

		return l.EnvironmentID
	}

	rows := []string{"LOADBALANCER ID | LOADBALANCER NAME | ENVIRONMENT"}
	for _, l := range loadBalancers {
		row := fmt.Sprintf("%s | %s | %s ",
			l.LoadBalancerID,
			l.LoadBalancerName,
			getEnvironment(l))

		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
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

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintLogs(logs ...*models.LogFile) error {
	for _, l := range logs {
		fmt.Println(l.Name)
		for i := 0; i < len(l.Name); i++ {
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
		if s.LoadBalancerName == "" {
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

	rows := []string{"SERVICE ID | SERVICE NAME | ENVIRONMENT | LOADBALANCER | DEPLOYMENTS | SCALE "}
	for _, s := range services {
		row := fmt.Sprintf("%s | %s | %s | %s | %s | %s",
			s.ServiceID,
			s.ServiceName,
			getEnvironment(s),
			getLoadBalancer(s),
			getDeployment(s, 0),
			getScale(s))

		rows = append(rows, row)

		// add the extra deployment rows
		for i := 1; i < len(s.Deployments); i++ {
			row := fmt.Sprintf(" | | | %s", getDeployment(s, i))
			rows = append(rows, row)
		}
	}

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintServiceSummaries(services ...*models.ServiceSummary) error {
	getEnvironment := func(s *models.ServiceSummary) string {
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

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintTasks(tasks ...*models.Task) error {
	getEnvironment := func(t *models.Task) string {
		if t.EnvironmentName != "" {
			return t.EnvironmentName
		}

		return t.EnvironmentID
	}

	getScale := func(t *models.Task) string {
		scale := fmt.Sprintf("%d/%d", t.RunningCount, t.DesiredCount)
		if t.PendingCount != 0 {
			scale = fmt.Sprintf("%s (%d)", scale, t.PendingCount)
		}

		return scale
	}

	getDeploy := func(t *models.Task) string {
		if t.DeployName != "" && t.DeployVersion != "" {
			return fmt.Sprintf("%s:%s", t.DeployName, t.DeployVersion)
		}

		return strings.Replace(t.DeployID, ".", ":", 1)
	}

	rows := []string{"TASK ID | TASK NAME | ENVIRONMENT | DEPLOY | SCALE "}
	for _, t := range tasks {
		row := fmt.Sprintf("%s | %s | %s | %s | %s",
			t.TaskID,
			t.TaskName,
			getEnvironment(t),
			getDeploy(t),
			getScale(t))

		rows = append(rows, row)
	}

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}

func (t *TextPrinter) PrintTaskSummaries(tasks ...*models.TaskSummary) error {
	getEnvironment := func(t *models.TaskSummary) string {
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

	fmt.Println(columnize.SimpleFormat(rows))
	return nil
}
