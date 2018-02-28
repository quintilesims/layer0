package printer

import (
	"github.com/quintilesims/layer0/common/models"
)

// testing stdout: https://blog.golang.org/examples

func ExampleTextPrintDeploys() {
	printer := &TextPrinter{}
	deploys := []*models.Deploy{
		{DeployID: "id1", DeployName: "name1", Version: "1"},
		{DeployID: "id2", DeployName: "name2", Version: "2"},
	}

	printer.PrintDeploys(deploys...)
	// Output:
	// DEPLOY ID  DEPLOY NAME  VERSION
	// id1        name1        1
	// id2        name2        2
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
			EnvironmentType: "static",
			OperatingSystem: "linux",
			CurrentScale:    2,
			DesiredScale:    3,
			InstanceType:    "t2.small",
			Links:           []string{"id2"},
		},
		{
			EnvironmentID:   "id2",
			EnvironmentName: "name2",
			EnvironmentType: "static",
			OperatingSystem: "windows",
			CurrentScale:    2,
			DesiredScale:    5,
			InstanceType:    "t2.small",
			Links:           []string{"id1", "api"},
		},
		{
			EnvironmentID:   "id3",
			EnvironmentName: "name3",
			EnvironmentType: "dynamic",
			OperatingSystem: "linux",
		},
	}

	printer.PrintEnvironments(environments...)
	// Output:
	// ENVIRONMENT ID  ENVIRONMENT NAME  TYPE     OS       LINKS
	// id1             name1             static   linux    id2
	// id2             name2             static   windows  id1
	//                                                     api
	// id3             name3             dynamic  linux
}

func ExampleTextPrintEnvironmentSummaries() {
	printer := &TextPrinter{}
	environments := []models.EnvironmentSummary{
		{EnvironmentID: "id1", EnvironmentName: "name1", OperatingSystem: "linux", EnvironmentType: "static"},
		{EnvironmentID: "id2", EnvironmentName: "name2", OperatingSystem: "linux", EnvironmentType: "dynamic"},
		{EnvironmentID: "id3", EnvironmentName: "name3", OperatingSystem: "windows", EnvironmentType: "static"},
	}

	printer.PrintEnvironmentSummaries(environments...)
	// Output:
	// ENVIRONMENT ID  ENVIRONMENT NAME  TYPE     OS
	// id1             name1             static   linux
	// id2             name2             dynamic  linux
	// id3             name3             static   windows
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
			IdleTimeout: 80,
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
			IdleTimeout: 90,
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

func ExampleTextPrintLoadBalancerIdleTimeout() {
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
		IdleTimeout: 75,
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
		IdleTimeout: 85,
	}

	printer.PrintLoadBalancerIdleTimeout(loadBalancer1)
	printer.PrintLoadBalancerIdleTimeout(loadBalancer2)
	// Output:
	// LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  IDLE TIMEOUT
	// id1              lb1                ename1       75
	// LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  IDLE TIMEOUT
	// id2              lb2                eid1         85
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
			ServiceID:        "id1",
			ServiceName:      "svc1",
			EnvironmentID:    "eid1",
			EnvironmentName:  "ename1",
			LoadBalancerID:   "lid1",
			LoadBalancerName: "lname1",
			RunningCount:     1,
			DesiredCount:     1,
			Deployments: []models.Deployment{
				{DeployName: "d1", DeployVersion: "1"},
			},
		},
		{
			ServiceID:      "id2",
			ServiceName:    "svc2",
			EnvironmentID:  "eid2",
			LoadBalancerID: "lid2",
			RunningCount:   1,
			DesiredCount:   1,
			Deployments: []models.Deployment{
				{DeployID: "d2.1"},
			},
		},
		{
			ServiceID:     "id3",
			ServiceName:   "svc3",
			EnvironmentID: "eid3",
			RunningCount:  0,
			DesiredCount:  1,
			PendingCount:  1,
			Deployments: []models.Deployment{
				{DeployID: "d3.1", RunningCount: 0, DesiredCount: 1},
			},
		},
		{
			ServiceID:     "id4",
			ServiceName:   "svc4",
			EnvironmentID: "eid4",
			RunningCount:  1,
			DesiredCount:  2,
			PendingCount:  1,
			Deployments: []models.Deployment{
				{DeployID: "d4.1", RunningCount: 1, DesiredCount: 2},
			},
		},
		{
			ServiceID:     "id5",
			ServiceName:   "svc5",
			EnvironmentID: "eid5",
			RunningCount:  2,
			DesiredCount:  1,
			Deployments: []models.Deployment{
				{DeployID: "d5.1", RunningCount: 1, DesiredCount: 0},
				{DeployID: "d5.2", RunningCount: 0, DesiredCount: 1},
			},
		},
	}

	printer.PrintServices(services...)
	// Output:
	// SERVICE ID  SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE
	// id1         svc1          ename1       lname1        d1:1         1/1
	// id2         svc2          eid2         lid2          d2:1         1/1
	// id3         svc3          eid3                       d3:1*        0/1 (1)
	// id4         svc4          eid4                       d4:1*        1/2 (1)
	// id5         svc5          eid5                       d5:1*        2/1
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
			TaskID:          "id1",
			TaskName:        "tsk1",
			EnvironmentID:   "eid1",
			EnvironmentName: "ename1",
			DeployName:      "d1",
			DeployVersion:   "1",
			Status:          "RUNNING",
		},
		{
			TaskID:        "id2",
			TaskName:      "tsk2",
			EnvironmentID: "eid2",
			DeployID:      "d2.1",
			Status:        "RUNNING",
		},
		{
			TaskID:        "id3",
			TaskName:      "tsk3",
			EnvironmentID: "eid3",
			DeployID:      "d3.1",
			Status:        "RUNNING",
		},
		{
			TaskID:        "id4",
			TaskName:      "tsk4",
			EnvironmentID: "eid4",
			DeployID:      "d4.1",
			Status:        "RUNNING",
		},
	}

	printer.PrintTasks(tasks...)
	// Output:
	// TASK ID  TASK NAME  ENVIRONMENT  DEPLOY  STATUS
	// id1      tsk1       ename1       d1:1    RUNNING
	// id2      tsk2       eid2         d2:1    RUNNING
	// id3      tsk3       eid3         d3:1    RUNNING
	// id4      tsk4       eid4         d4:1    RUNNING
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
