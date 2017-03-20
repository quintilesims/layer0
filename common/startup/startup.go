package startup

import (
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/aws/elb"
	"github.com/quintilesims/layer0/common/aws/iam"
	"github.com/quintilesims/layer0/common/aws/provider"
	"github.com/quintilesims/layer0/common/aws/s3"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/db/job_store"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/decorators"
	"github.com/quintilesims/layer0/common/waitutils"
)

func GetBackend(credProvider provider.CredProvider, region string) (backend.Backend, error) {
	s3Provider, err := s3.NewS3(credProvider, region)
	if err != nil {
		return nil, err
	}

	iamProvider, err := iam.NewIAM(credProvider, region)
	if err != nil {
		return nil, err
	}

	ec2Provider, err := ec2.NewEC2(credProvider, region)
	if err != nil {
		return nil, err
	}

	elbProvider, err := elb.NewELB(credProvider, region)
	if err != nil {
		return nil, err
	}

	cloudWatchLogsProvider, err := cloudwatchlogs.NewCloudWatchLogs(credProvider, region)
	if err != nil {
		return nil, err
	}

	tagStore, err := getNewTagStore()
	if err != nil {
		return nil, err
	}

	ec2Provider = wrapEC2(ec2Provider)
	elbProvider = wrapELB(elbProvider)
	cloudWatchLogsProvider = wrapCloudWatchLogs(cloudWatchLogsProvider)

	ecsProvider, err := GetECS(credProvider, region)
	if err != nil {
		return nil, err
	}

	autoscalingProvider, err := GetAutoscaling(credProvider, region)
	if err != nil {
		return nil, err
	}

	backend := ecsbackend.NewBackend(
		tagStore,
		s3Provider,
		iamProvider,
		ec2Provider,
		ecsProvider,
		elbProvider,
		autoscalingProvider,
		cloudWatchLogsProvider)

	return backend, nil
}

func GetECS(credProvider provider.CredProvider, region string) (ecs.Provider, error) {
	ecsProvider, err := ecs.NewECS(credProvider, region)
	if err != nil {
		return nil, err
	}

	return wrapECS(ecsProvider), nil
}

func GetAutoscaling(credProvider provider.CredProvider, region string) (autoscaling.Provider, error) {
	autoscalingProvider, err := autoscaling.NewAutoScaling(credProvider, region)
	if err != nil {
		return nil, err
	}

	return wrapAutoscaling(autoscalingProvider), nil
}

func GetLogic(backend backend.Backend) (*logic.Logic, error) {
	tagStore, err := getNewTagStore()
	if err != nil {
		return nil, err
	}

	jobStore, err := getNewJobStore()
	if err != nil {
		return nil, err
	}

	return logic.NewLogic(tagStore, jobStore, backend, nil), nil
}

func getNewTagStore() (tag_store.TagStore, error) {
	store := tag_store.NewMysqlTagStore(db.Config{
		Connection: config.DBConnection(),
		DBName:     config.DBName(),
	})

	if err := store.Init(); err != nil {
		return nil, err
	}

	return store, nil
}

func getNewJobStore() (job_store.JobStore, error) {
	store := job_store.NewMysqlJobStore(db.Config{
		Connection: config.DBConnection(),
		DBName:     config.DBName(),
	})

	if err := store.Init(); err != nil {
		return nil, err
	}

	return store, nil
}

func wrapECS(e ecs.Provider) ecs.Provider {
	retry := &decorators.Retry{
		Clock: waitutils.RealClock{},
	}

	wrap := &ecs.ProviderDecorator{
		Inner:     e,
		Decorator: retry.CallWithRetries,
	}

	return wrap
}

func wrapAutoscaling(a autoscaling.Provider) autoscaling.Provider {
	retry := &decorators.Retry{
		Clock: waitutils.RealClock{},
	}

	wrap := &autoscaling.ProviderDecorator{
		Inner:     a,
		Decorator: retry.CallWithRetries,
	}

	return wrap
}

func wrapEC2(e ec2.Provider) ec2.Provider {
	wrap := &ec2.ProviderDecorator{
		Inner:     e,
		Decorator: decorators.CallWithLogging,
	}

	return wrap
}

func wrapELB(e elb.Provider) elb.Provider {
	wrap := &elb.ProviderDecorator{
		Inner:     e,
		Decorator: decorators.CallWithLogging,
	}

	return wrap
}

func wrapCloudWatchLogs(c cloudwatchlogs.Provider) cloudwatchlogs.Provider {
	wrap := &cloudwatchlogs.ProviderDecorator{
		Inner:     c,
		Decorator: decorators.CallWithLogging,
	}

	retry := &decorators.Retry{
		Clock: waitutils.RealClock{},
	}

	wrap = &cloudwatchlogs.ProviderDecorator{
		Inner:     wrap,
		Decorator: retry.CallWithRetries,
	}

	return wrap
}
