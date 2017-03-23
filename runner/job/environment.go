package job

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

var DeleteEnvironmentSteps = []Step{
	Step{
		Name:    "Delete Dependencies",
		Timeout: time.Minute * 15,
		Action:  Fold(DeleteEnvironmentLoadBalancers, DeleteEnvironmentServices, DeleteEnvironmentTasks),
	},
	Step{
		Name:    "Delete Environment",
		Timeout: time.Minute * 10,
		Action:  DeleteEnvironment,
	},
}

func DeleteEnvironment(quit chan bool, context *JobContext) error {
	log.Infof("Running Action: DeleteEnvironment")
	environmentID := context.Request()

	return runAndRetry(quit, time.Second*10, func() error {
		log.Infof("Running Action: DeleteEnvironment on '%s'", environmentID)
		return context.EnvironmentLogic.DeleteEnvironment(environmentID)
	})
}

func DeleteEnvironmentLoadBalancers(quit chan bool, context *JobContext) error {
	log.Infof("Running Action: DeleteEnvironmentLoadBalancers")
	environmentID := context.Request()

	loadBalancers, err := context.LoadBalancerLogic.ListLoadBalancers()
	if err != nil {
		return err
	}

	for i := 0; i < len(loadBalancers); i++ {
		if loadBalancers[i].EnvironmentID != environmentID {
			loadBalancers = append(loadBalancers[:i], loadBalancers[i+1:]...)
			i--
		}
	}

	actions := []Action{}
	for _, loadBalancer := range loadBalancers {
		loadBalancerContext := context.CreateCopyWithNewRequest(loadBalancer.LoadBalancerID)
		action := func(chan bool, *JobContext) error {
			return DeleteLoadBalancer(quit, loadBalancerContext)
		}

		actions = append(actions, action)
	}

	runAll := Fold(actions...)
	return runAll(quit, nil)
}

func DeleteEnvironmentServices(quit chan bool, context *JobContext) error {
	log.Infof("Running Action: DeleteEnvironmentServices")
	environmentID := context.Request()

	services, err := context.ServiceLogic.ListServices()
	if err != nil {
		return err
	}

	for i := 0; i < len(services); i++ {
		if services[i].EnvironmentID != environmentID {
			services = append(services[:i], services[i+1:]...)
			i--
		}
	}

	actions := []Action{}
	for _, service := range services {
		serviceContext := context.CreateCopyWithNewRequest(service.ServiceID)
		action := func(chan bool, *JobContext) error {
			return DeleteService(quit, serviceContext)
		}

		actions = append(actions, action)
	}

	runAll := Fold(actions...)
	return runAll(quit, nil)
}

func DeleteEnvironmentTasks(quit chan bool, context *JobContext) error {
	log.Infof("Running Action: DeleteEnvironmentTasks")
	environmentID := context.Request()

	tasks, err := context.TaskLogic.ListTasks()
	if err != nil {
		return err
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].EnvironmentID != environmentID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			i--
		}
	}

	actions := []Action{}
	for _, task := range tasks {
		taskContext := context.CreateCopyWithNewRequest(task.TaskID)
		action := func(chan bool, *JobContext) error {
			return DeleteTask(quit, taskContext)
		}

		actions = append(actions, action)
	}

	runAll := Fold(actions...)
	return runAll(quit, nil)
}
