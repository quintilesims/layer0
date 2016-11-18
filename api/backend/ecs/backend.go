package ecsbackend

import (
	"gitlab.imshealth.com/xfra/layer0/api/data"
	"gitlab.imshealth.com/xfra/layer0/common/aws/autoscaling"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ec2"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/elb"
	"gitlab.imshealth.com/xfra/layer0/common/aws/iam"
	"gitlab.imshealth.com/xfra/layer0/common/aws/s3"
)

// todo: this is an awkward design pattern - we don't need to split the ECSBackend
// object into many parts with self-references to 'this.Backend'
// we should just make the backend one single object with many functions
type ECSBackend struct {
	*ECSEnvironmentManager
	*ECSServiceManager
	*ECSDeployManager
	*ECSCertificateManager
	*ECSLoadBalancerManager
	*ECSTaskManager
	*ECSRightSizer
}

func NewBackend(
	sqlAdmin data.SQLAdmin,
	tagData data.TagData,
	s3 s3.Provider,
	iam iam.Provider,
	ec2 ec2.Provider,
	ecs ecs.Provider,
	elb elb.Provider,
	autoscaling autoscaling.Provider,
	cloudWatchLogs cloudwatchlogs.Provider,
) *ECSBackend {

	backend := &ECSBackend{}

	cluster := NewECSClusterScaler(ecs, autoscaling, backend)

	backend.ECSEnvironmentManager = NewECSEnvironmentManager(ecs, ec2, autoscaling, backend)
	backend.ECSServiceManager = NewECSServiceManager(ecs, ec2, cloudWatchLogs, cluster, backend)
	backend.ECSLoadBalancerManager = NewECSLoadBalancerManager(ec2, elb, iam, backend)
	backend.ECSCertificateManager = NewECSCertificateManager(iam)
	backend.ECSDeployManager = NewECSDeployManager(ecs, cluster)
	backend.ECSTaskManager = NewECSTaskManager(ecs, cloudWatchLogs, backend, cluster)
	backend.ECSRightSizer = NewECSRightSizer(ecs, autoscaling, cluster, backend)

	return backend
}
