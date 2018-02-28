package printer

import (
	"github.com/quintilesims/layer0/common/models"
)

// testing stdout: https://blog.golang.org/examples

func ExampleTextPrintDeploys() {
	printer := &TextPrinter{}
	deploys := []*models.Deploy{
		{
			Compatibilities: []string{models.DeployCompatibilityStateless},
			DeployID:        "id1",
			DeployName:      "name1",
			Version:         "1",
		},
		{
			Compatibilities: []string{
				models.DeployCompatibilityStateless,
				models.DeployCompatibilityStateful,
			},
			DeployID:   "id2",
			DeployName: "name2",
			Version:    "2",
		},
	}

	printer.PrintDeploys(deploys...)
	// Output:
	// DEPLOY ID  DEPLOY NAME  VERSION  COMPATIBILITIES
	// id1        name1        1        stateless
	// id2        name2        2        stateless
	//                                  stateful
}

func ExampleTextPrintDeploySummaries() {
	printer := &TextPrinter{}
	deploys := []models.DeploySummary{
		{DeployID: "id1", DeployName: "name1", Version: "1"},
		{DeployID: "id2", DeployName: "name2", Version: "2"},
	}

	printer.PrintDeploySummaries(deploys...)
	// Output:
	// DEPLOY ID  DEPLOY NAME  VERSION
	// id1        name1        1
	// id2        name2        2
}

func ExampleTextPrintEnvironments() {
	printer := &TextPrinter{}
	environments := []*models.Environment{
		{
			EnvironmentID:   "id1",
			EnvironmentName: "name1",
			OperatingSystem: "linux",
			CurrentScale:    2,
			DesiredScale:    3,
			InstanceType:    "t2.small",
			Links:           []string{"id2"},
		},
		{
			EnvironmentID:   "id2",
			EnvironmentName: "name2",
			OperatingSystem: "windows",
			CurrentScale:    2,
			DesiredScale:    5,
			InstanceType:    "t2.small",
			Links:           []string{"id1", "api"},
		},
		{
			EnvironmentID:   "id3",
			EnvironmentName: "name3",
			OperatingSystem: "linux",
		},
	}

	printer.PrintEnvironments(environments...)
	// Output:
	// ENVIRONMENT ID  ENVIRONMENT NAME  OS       LINKS
	// id1             name1             linux    id2
	// id2             name2             windows  id1
	//                                            api
	// id3             name3             linux
}

func ExampleTextPrintEnvironmentSummaries() {
	printer := &TextPrinter{}
	environments := []models.EnvironmentSummary{
		{EnvironmentID: "id1", EnvironmentName: "name1", OperatingSystem: "linux"},
		{EnvironmentID: "id2", EnvironmentName: "name2", OperatingSystem: "linux"},
		{EnvironmentID: "id3", EnvironmentName: "name3", OperatingSystem: "windows"},
	}

	printer.PrintEnvironmentSummaries(environments...)
	// Output:
	// ENVIRONMENT ID  ENVIRONMENT NAME  OS
	// id1             name1             linux
	// id2             name2             linux
	// id3             name3             windows
}

func ExampleTextPrintLoadBalancers() {
	printer := &TextPrinter{}
	loadBalancers := []*models.LoadBalancer{
		{
			LoadBalancerID:   "id1",
			LoadBalancerName: "lb1",
			EnvironmentID:    "eid1",
			EnvironmentName:  "ename1",
			ServiceID:        "sid1",
			ServiceName:      "sname1",
			IsPublic:         true,
			URL:              "url1",
			Ports: []models.Port{
				{
					HostPort:      80,
					ContainerPort: 80,
					Protocol:      "http",
				},
			},
		},
		{
			LoadBalancerID:   "id2",
			LoadBalancerName: "lb2",
			EnvironmentID:    "eid2",
			ServiceID:        "sid2",
			IsPublic:         false,
			URL:              "url2",
			Ports: []models.Port{
				{
					HostPort:      443,
					ContainerPort: 80,
					Protocol:      "https",
				},
				{
					HostPort:      22,
					ContainerPort: 22,
					Protocol:      "tcp",
				},
			},
		},
	}

	printer.PrintLoadBalancers(loadBalancers...)
	// Output:
	// LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS         PUBLIC  URL
	// id1              lb1                ename1       sname1   80:80/HTTP    true    url1
	// id2              lb2                eid2         sid2     443:80/HTTPS  false   url2
	//                                                           22:22/TCP

}

func ExampleTextPrintLoadBalancerSummaries() {
	printer := &TextPrinter{}
	loadBalancers := []models.LoadBalancerSummary{
		{LoadBalancerID: "id1", LoadBalancerName: "lb1", EnvironmentID: "eid1", EnvironmentName: "ename1"},
		{LoadBalancerID: "id2", LoadBalancerName: "lb2", EnvironmentID: "eid2"},
	}

	printer.PrintLoadBalancerSummaries(loadBalancers...)
	// Output:
	// LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT
	// id1              lb1                ename1
	// id2              lb2                eid2
}

func ExampleTextPrintLoadBalancerHealthCheck() {
	printer := &TextPrinter{}
	loadBalancer1 := &models.LoadBalancer{
		LoadBalancerID:   "id1",
		LoadBalancerName: "lb1",
		EnvironmentID:    "eid1",
		EnvironmentName:  "ename1",
		HealthCheck: models.HealthCheck{
			Target:             "tcp:22",
			Interval:           5,
			Timeout:            30,
			HealthyThreshold:   10,
			UnhealthyThreshold: 2,
		},
	}

	loadBalancer2 := &models.LoadBalancer{
		LoadBalancerID:   "id2",
		LoadBalancerName: "lb2",
		EnvironmentID:    "eid1",
		HealthCheck: models.HealthCheck{
			Target:             "http:80/health",
			Interval:           6,
			Timeout:            10,
			HealthyThreshold:   3,
			UnhealthyThreshold: 5,
		},
	}

	printer.PrintLoadBalancerHealthCheck(loadBalancer1)
	printer.PrintLoadBalancerHealthCheck(loadBalancer2)
	// Output:
	// LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  TARGET  INTERVAL  TIMEOUT  HEALTHY THRESHOLD  UNHEALTHY THRESHOLD
	// id1              lb1                ename1       tcp:22  5         30       10                 2
	// LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  TARGET          INTERVAL  TIMEOUT  HEALTHY THRESHOLD  UNHEALTHY THRESHOLD
	// id2              lb2                eid1         http:80/health  6         10       3                  5
}

func ExampleTextPrintLogs() {
	printer := &TextPrinter{}
	logs := []models.LogFile{
		{ContainerName: "file1", Lines: []string{"line1", "line2", "line3"}},
		{ContainerName: "file2", Lines: []string{"lineA", "lineB", "lineC"}},
	}

	printer.PrintLogs(logs...)
	// Output:
	//file1
	//-----
	//line1
	//line2
	//line3
	//
	//file2
	//-----
	//lineA
	//lineB
	//lineC
}

func ExampleTextPrintServices() {
	printer := &TextPrinter{}
	services := []*models.Service{
		{
			Deployments: []models.Deployment{
				{DeployName: "d1", DeployVersion: "1"},
			},
			DesiredCount:     1,
			EnvironmentID:    "eid1",
			EnvironmentName:  "ename1",
			LoadBalancerID:   "lid1",
			LoadBalancerName: "lname1",
			RunningCount:     1,
			ServiceID:        "id1",
			ServiceName:      "svc1",
			ServiceType:      models.DeployCompatibilityStateless,
		},
		{
			Deployments: []models.Deployment{
				{DeployID: "d2.1"},
			},
			DesiredCount:   1,
			EnvironmentID:  "eid2",
			LoadBalancerID: "lid2",
			RunningCount:   1,
			ServiceID:      "id2",
			ServiceName:    "svc2",
			ServiceType:    models.DeployCompatibilityStateful,
		},
		{
			Deployments: []models.Deployment{
				{DeployID: "d3.1", RunningCount: 0, DesiredCount: 1},
			},
			DesiredCount:  1,
			EnvironmentID: "eid3",
			PendingCount:  1,
			RunningCount:  0,
			ServiceID:     "id3",
			ServiceName:   "svc3",
			ServiceType:   models.DeployCompatibilityStateless,
		},
		{
			Deployments: []models.Deployment{
				{DeployID: "d4.1", RunningCount: 1, DesiredCount: 2},
			},
			DesiredCount:  2,
			EnvironmentID: "eid4",
			PendingCount:  1,
			RunningCount:  1,
			ServiceID:     "id4",
			ServiceName:   "svc4",
			ServiceType:   models.DeployCompatibilityStateful,
		},
		{
			Deployments: []models.Deployment{
				{DeployID: "d5.1", RunningCount: 1, DesiredCount: 0},
				{DeployID: "d5.2", RunningCount: 0, DesiredCount: 1},
			},
			DesiredCount:  1,
			EnvironmentID: "eid5",
			RunningCount:  2,
			ServiceID:     "id5",
			ServiceName:   "svc5",
			ServiceType:   models.DeployCompatibilityStateful,
		},
	}

	printer.PrintServices(services...)
	// Output:
	// SERVICE ID  SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE    TYPE
	// id1         svc1          ename1       lname1        d1:1         1/1      stateless
	// id2         svc2          eid2         lid2          d2:1         1/1      stateful
	// id3         svc3          eid3                       d3:1*        0/1 (1)  stateless
	// id4         svc4          eid4                       d4:1*        1/2 (1)  stateful
	// id5         svc5          eid5                       d5:1*        2/1      stateful
	//                                                      d5:2*
}

func ExampleTextPrintServiceSummaries() {
	printer := &TextPrinter{}
	services := []models.ServiceSummary{
		{ServiceID: "id1", ServiceName: "svc1", EnvironmentID: "eid1", EnvironmentName: "ename1"},
		{ServiceID: "id2", ServiceName: "svc2", EnvironmentID: "eid2"},
	}

	printer.PrintServiceSummaries(services...)
	// Output:
	// SERVICE ID  SERVICE NAME  ENVIRONMENT
	// id1         svc1          ename1
	// id2         svc2          eid2
}

func ExampleTextPrintTasks() {
	printer := &TextPrinter{}
	tasks := []*models.Task{
		{
			DeployName:      "d1",
			DeployVersion:   "1",
			EnvironmentID:   "eid1",
			EnvironmentName: "ename1",
			Status:          "RUNNING",
			TaskID:          "id1",
			TaskName:        "tsk1",
			TaskType:        models.DeployCompatibilityStateless,
		},
		{
			DeployID:      "d2.1",
			EnvironmentID: "eid2",
			Status:        "RUNNING",
			TaskID:        "id2",
			TaskName:      "tsk2",
			TaskType:      models.DeployCompatibilityStateful,
		},
		{
			DeployID:      "d3.1",
			EnvironmentID: "eid3",
			Status:        "RUNNING",
			TaskID:        "id3",
			TaskName:      "tsk3",
			TaskType:      models.DeployCompatibilityStateless,
		},
		{
			DeployID:      "d4.1",
			EnvironmentID: "eid4",
			Status:        "RUNNING",
			TaskID:        "id4",
			TaskName:      "tsk4",
			TaskType:      models.DeployCompatibilityStateful,
		},
	}

	printer.PrintTasks(tasks...)
	// Output:
	// TASK ID  TASK NAME  ENVIRONMENT  DEPLOY  STATUS   TYPE
	// id1      tsk1       ename1       d1:1    RUNNING  stateless
	// id2      tsk2       eid2         d2:1    RUNNING  stateful
	// id3      tsk3       eid3         d3:1    RUNNING  stateless
	// id4      tsk4       eid4         d4:1    RUNNING  stateful
}

func ExampleTextPrintTaskSummaries() {
	printer := &TextPrinter{}
	tasks := []models.TaskSummary{
		{TaskID: "id1", TaskName: "tsk1", EnvironmentID: "eid1", EnvironmentName: "ename1"},
		{TaskID: "id2", TaskName: "tsk2", EnvironmentID: "eid2"},
	}

	printer.PrintTaskSummaries(tasks...)
	// Output:
	// TASK ID  TASK NAME  ENVIRONMENT
	// id1      tsk1       ename1
	// id2      tsk2       eid2
}
