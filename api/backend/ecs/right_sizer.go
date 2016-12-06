package ecsbackend

import (
	"fmt"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/common/waitutils"
	"time"
)

const (
	RIGHTSIZER_SLEEP_DURATION = 60 * time.Minute
)

var rsLogger = logutils.NewStackTraceLogger("Right Sizer")

type ECSRightSizer struct {
	ECS           ecs.Provider
	Backend       backend.Backend
	Autoscaling   autoscaling.Provider
	ClusterScaler ClusterScaler
	Clock         waitutils.Clock
	lastRunTime   time.Time
}

func NewECSRightSizer(ecsprovider ecs.Provider, asg autoscaling.Provider, cluster ClusterScaler, b backend.Backend) *ECSRightSizer {
	return &ECSRightSizer{
		ECS:           ecsprovider,
		Backend:       b,
		Autoscaling:   asg,
		ClusterScaler: cluster,
		Clock:         waitutils.RealClock{},
	}
}

func (this *ECSRightSizer) StartRightSizer() {
	go func() {
		for {
			this.lastRunTime = this.Clock.Now()
			if err := this.run(); err != nil {
				rsLogger.Errorf("%v", err)
			}

			this.Clock.Sleep(RIGHTSIZER_SLEEP_DURATION)
		}
	}()
}

func (this *ECSRightSizer) run() error {
	environments, err := this.Backend.ListEnvironments()
	if err != nil {
		return err
	}

	for _, environment := range environments {
		environmentID := id.L0EnvironmentID(environment.EnvironmentID)
		rsLogger.Infof("Running on environment '%s'", environmentID)

		instanceDiff, err := this.optimizeCluster(environmentID.ECSEnvironmentID())
		if err != nil {
			rsLogger.Errorf("OptimizeCluster error on environment '%s': %v", environmentID, err)
			continue
		}

		if instanceDiff != 0 {
			rsLogger.Infof("OptimizeCluster on environment '%s' saw an instance differential of %v", environmentID, instanceDiff)
		}

		rsLogger.Infof("Finished on environment '%s'", environmentID)
	}

	return nil
}

func (this *ECSRightSizer) optimizeCluster(ecsEnvironmentID id.ECSEnvironmentID) (int, error) {
	newInstancesAdded, hasPendingTasks, err := this.ClusterScaler.TriggerScalingAlgorithm(ecsEnvironmentID, nil, 0)
	if err != nil {
		return 0, err
	}

	instanceDiff := newInstancesAdded
	if newInstancesAdded == 0 && !hasPendingTasks {
		// we can scale down
		// super simple for now -- scan for CIs that have no tasks running
		// this is good enough for smallish container workloads (say <100)
		// more than that and you'd want something more intelligent to get
		// better packing
		instanceARNs, err := this.ECS.ListContainerInstances(ecsEnvironmentID.String())
		if err != nil {
			return 0, err
		}

		if len(instanceARNs) > 0 {
			instances, err := this.ECS.DescribeContainerInstances(ecsEnvironmentID.String(), instanceARNs)
			if err != nil {
				return 0, err
			}

			for _, instance := range instances {
				deleted, err := this.terminateInstanceIfEmpty(ecsEnvironmentID, instance)
				if err != nil {
					return 0, err
				}

				if deleted {
					instanceDiff -= 1
				}
			}
		}
	}

	return instanceDiff, nil
}

func (this *ECSRightSizer) terminateInstanceIfEmpty(ecsEnvironmentID id.ECSEnvironmentID, instance *ecs.ContainerInstance) (bool, error) {
	if *instance.AgentConnected == false {
		rsLogger.Warningf("Cluster '%s' found container instance '%s' with disconnected ECS agent but %d running tasks",
			ecsEnvironmentID,
			*instance.Ec2InstanceId,
			*instance.RunningTasksCount)

		// don't decrement capacity since we want autoscaling to re-create the instance
		if _, err := this.Autoscaling.TerminateInstanceInAutoScalingGroup(*instance.Ec2InstanceId, false); err != nil {
			return false, err
		}

		return true, nil
	}

	if *instance.RunningTasksCount == 0 && *instance.PendingTasksCount == 0 {
		rsLogger.Infof("Cluster '%s' found extra instance '%s'", ecsEnvironmentID, *instance.Ec2InstanceId)

		asg, err := this.Autoscaling.DescribeAutoScalingGroup(ecsEnvironmentID.String())
		if err != nil {
			return false, err
		}

		if *asg.DesiredCapacity == *asg.MinSize {
			rsLogger.Warningf("Cluster '%s' is at minimum capacity. Keeping extra instance '%s'", ecsEnvironmentID, *instance.Ec2InstanceId)
			return false, nil
		}

		rsLogger.Infof("Terminating instance '%s' in cluster '%s'", *instance.Ec2InstanceId, ecsEnvironmentID)
		if _, err := this.Autoscaling.TerminateInstanceInAutoScalingGroup(*instance.Ec2InstanceId, true); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (this *ECSRightSizer) GetRightSizerHealth() (string, error) {
	timeSinceCompletion := this.Clock.Since(this.lastRunTime)

	if timeSinceCompletion > RIGHTSIZER_SLEEP_DURATION*2 {
		return "", fmt.Errorf("RightSizer hasn't completed in %v", timeSinceCompletion)
	}

	return fmt.Sprintf("RightSizer heartbeat %v ago", timeSinceCompletion), nil
}
