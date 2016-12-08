package ecsbackend

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/models"
	"sync"
	"time"
)

type TaskScheduler interface {
	GetTask(id.ECSTaskID) []*ecs.Task
	ListTasks() []*ecs.Task
	DeleteTask(id.ECSTaskID)
	AddTask(id.ECSTaskID, id.ECSDeployID, id.ECSEnvironmentID, int, []models.ContainerOverride)
}

// Unlike services, ECS doesn't provide a scheduling mechanism for tasks.
// The L0TaskScheduler implements the missing pieces:
//     - Retry calling ecs.RunTask when a task fails to start because of low cluster capacity
//     - Hold onto tasks that failed to start. ESC does not tack info about tasks that failed to start.
// Once a task has successfully started in ECS, it no longer needs to be tracked in this scheduler.

const SCHEDULED_TASK_TIMEOUT = time.Hour * 1

type ScheduledTask struct {
	Started       time.Time
	Quitc         chan bool
	TaskID        id.ECSTaskID
	DeployID      id.ECSDeployID
	EnvironmentID id.ECSEnvironmentID
	Overrides     []models.ContainerOverride
	ECSTask       *ecs.Task
}

type L0TaskScheduler struct {
	ECS   ecs.Provider
	mutex *sync.Mutex
	tasks map[id.ECSTaskID][]*ScheduledTask
}

func NewL0TaskScheduler(provider ecs.Provider) *L0TaskScheduler {
	scheduler := &L0TaskScheduler{
		ECS:   provider,
		mutex: &sync.Mutex{},
		tasks: map[id.ECSTaskID][]*ScheduledTask{},
	}

	// run the cleanup process forever async
	go func() {
		for {
			scheduler.cleanup()
			time.Sleep(time.Minute * 5)
		}
	}()

	return scheduler
}

func (this *L0TaskScheduler) cleanup() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for taskID := range this.tasks {
		copies := this.tasks[taskID]

		for i := 0; i < len(copies); i++ {
			task := copies[i]

			if time.Since(task.Started) > SCHEDULED_TASK_TIMEOUT {
				log.Infof("[TaskScheduler] Cleanup routine is removing old task '%s'", task.TaskID)
				close(task.Quitc)
				copies = append(copies[:i], copies[i+1:]...)
				i--
			}
		}

		this.tasks[taskID] = copies
	}
}

func (this *L0TaskScheduler) ListTasks() []*ecs.Task {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tasks := []*ecs.Task{}
	for _, copies := range this.tasks {
		for _, copy := range copies {
			tasks = append(tasks, copy.ECSTask)
		}
	}

	return tasks
}

func (this *L0TaskScheduler) GetTask(taskID id.ECSTaskID) []*ecs.Task {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	tasks := []*ecs.Task{}
	if copies, ok := this.tasks[taskID]; ok {
		for _, copy := range copies {
			tasks = append(tasks, copy.ECSTask)
		}
	}

	return tasks
}

func (this *L0TaskScheduler) DeleteTask(taskID id.ECSTaskID) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	copies, exists := this.tasks[taskID]
	if !exists {
		return
	}

	for _, copy := range copies {
		close(copy.Quitc)
	}

	delete(this.tasks, taskID)
}

func (this *L0TaskScheduler) deleteTaskCopy(copy *ScheduledTask) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	copies, exists := this.tasks[copy.TaskID]
	if !exists {
		return
	}

	for i := 0; i < len(copies); i++ {
		if copies[i] == copy {
			close(copy.Quitc)
			copies = append(copies[:i], copies[i+1:]...)
			break
		}
	}

	this.tasks[copy.TaskID] = copies
}

// adds the specified number of copies for the given taskID
func (this *L0TaskScheduler) AddTask(
	taskID id.ECSTaskID,
	deployID id.ECSDeployID,
	environmentID id.ECSEnvironmentID,
	numCopies int,
	overrides []models.ContainerOverride,
) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	copies, exists := this.tasks[taskID]
	if !exists {
		copies = []*ScheduledTask{}
	}

	for i := 0; i < numCopies; i++ {
		scheduledTask := &ScheduledTask{
			Started:       time.Now(),
			Quitc:         make(chan bool),
			TaskID:        taskID,
			DeployID:      deployID,
			EnvironmentID: environmentID,
			Overrides:     overrides,
			ECSTask:       ecsPendingTask(taskID, deployID, environmentID),
		}

		copies = append(copies, scheduledTask)
		go this.runTaskWithRetries(scheduledTask)
	}

	this.tasks[taskID] = copies
}

func (this *L0TaskScheduler) runTaskWithRetries(task *ScheduledTask) {
	shouldRetry := true

	for {
		if shouldRetry {
			shouldRetry = this.attemptRunTask(task)
		}

		select {
		case <-task.Quitc:
			log.Infof("[TaskScheduler] Quit signalled for task '%s'", task.TaskID)
			return

		case <-time.After(time.Minute * 1):
		}
	}
}

func (this *L0TaskScheduler) attemptRunTask(task *ScheduledTask) bool {
	clusterName := task.EnvironmentID.String()
	taskDefinition := task.DeployID.TaskDefinition()
	startedBy := stringp(task.TaskID.String())
	overrides := []*ecs.ContainerOverride{}

	for _, override := range task.Overrides {
		o := ecs.NewContainerOverride(override.ContainerName, override.EnvironmentOverrides)
		overrides = append(overrides, o)
	}

	log.Infof("[TaskScheduler] Attempting to start task '%s'", task.TaskID)
	_, failed, err := this.ECS.RunTask(clusterName, taskDefinition, 1, startedBy, overrides)
	if err != nil {
		if ContainsErrMsg(err, "No Container Instances were found in your cluster") {
			// do retry since we are just waiting for cluster to come online
			log.Infof("TaskScheduler] Waiting for increased cluster size for task '%s'...", task.TaskID)
			return true
		} else {
			// don't retry since we aren't sure what went wrong
			log.Infof("[TaskScheduler] Task '%s' failed to start: %v", task.TaskID, err)
			task.ECSTask = ecsFailedTask(task.TaskID, task.DeployID, task.EnvironmentID, err)
			return false
		}
	}

	if numFailed := len(failed); numFailed > 0 {
		// do retry since we assume failure to start is due to cluster size
		// todo: check failure.Reason to determine more specific course of action
		log.Infof("[TaskScheduler] Task '%s' returned a failure: %v", task.TaskID, failed)
		log.Infof("TaskScheduler] Waiting for increased cluster size for task '%s'...", task.TaskID)
		return true
	}

	// Tasks are tracked in ECS once they have been successfully started
	log.Infof("[TaskScheduler] Task '%s' has been started. Removing from scheduler", task.TaskID)

	// deleteTaskCopy will signal the quit channel along with removing the copy from the scheduler
	// this will stop the loop in runTaskWithRetries
	this.deleteTaskCopy(task)
	return false
}

func ecsPendingTask(taskID id.ECSTaskID, deployID id.ECSDeployID, environmentID id.ECSEnvironmentID) *ecs.Task {
	dummy := ecs.NewTask("", "", "")

	dummy.LastStatus = stringp("PENDING")
	dummy.DesiredStatus = stringp("RUNNING")
	dummy.StartedBy = stringp(taskID.String())
	dummy.StoppedReason = stringp("Waiting for cluster capacity to run")

	taskDefARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:task/%s", deployID.TaskDefinition())
	dummy.TaskDefinitionArn = stringp(taskDefARN)

	clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
	dummy.ClusterArn = stringp(clusterARN)

	return dummy
}

func ecsFailedTask(taskID id.ECSTaskID, deployID id.ECSDeployID, environmentID id.ECSEnvironmentID, err error) *ecs.Task {
	dummy := ecs.NewTask("", "", "")

	dummy.LastStatus = stringp("PENDING")
	dummy.DesiredStatus = stringp("RUNNING")
	dummy.StartedBy = stringp(taskID.String())
	dummy.StoppedReason = stringp(err.Error())

	taskDefARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:task/%s", deployID.TaskDefinition())
	dummy.TaskDefinitionArn = stringp(taskDefARN)

	clusterARN := fmt.Sprintf("arn:aws:ecs:region:aws_account_id:cluster/%s", environmentID.String())
	dummy.ClusterArn = stringp(clusterARN)

	return dummy
}
