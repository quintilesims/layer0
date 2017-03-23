package ecsbackend

import (
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/aws/elb"
	"github.com/quintilesims/layer0/common/aws/iam"
	"github.com/quintilesims/layer0/common/aws/s3"
	"github.com/quintilesims/layer0/common/db/tag_store"
)

// todo: this is an awkward design pattern - we don't need to split the ECSBackend
// object into many parts with self-references to 'this.Backend'
// we should just make the backend one single object with many functions
type ECSBackend struct {
	*ECSEnvironmentManager
	*ECSServiceManager
	*ECSDeployManager
	*ECSLoadBalancerManager
	*ECSTaskManager
}

func NewBackend(
	tagData tag_store.TagStore,
	s3 s3.Provider,
	iam iam.Provider,
	ec2 ec2.Provider,
	ecs ecs.Provider,
	elb elb.Provider,
	autoscaling autoscaling.Provider,
	cloudWatchLogs cloudwatchlogs.Provider,
) *ECSBackend {

	backend := &ECSBackend{}

	backend.ECSEnvironmentManager = NewECSEnvironmentManager(ecs, ec2, autoscaling, backend)
	backend.ECSServiceManager = NewECSServiceManager(ecs, ec2, cloudWatchLogs, backend)
	backend.ECSLoadBalancerManager = NewECSLoadBalancerManager(ec2, elb, iam, backend)
	backend.ECSDeployManager = NewECSDeployManager(ecs)
	backend.ECSTaskManager = NewECSTaskManager(ecs, cloudWatchLogs, backend)

	return backend
}
