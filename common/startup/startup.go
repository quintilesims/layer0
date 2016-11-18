package startup

import (
	"gitlab.imshealth.com/xfra/layer0/api/backend"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs"
	"gitlab.imshealth.com/xfra/layer0/api/data"
	"gitlab.imshealth.com/xfra/layer0/api/logic"
	"gitlab.imshealth.com/xfra/layer0/common/aws/autoscaling"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ec2"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/elb"
	"gitlab.imshealth.com/xfra/layer0/common/aws/iam"
	"gitlab.imshealth.com/xfra/layer0/common/aws/provider"
	"gitlab.imshealth.com/xfra/layer0/common/aws/s3"
	"gitlab.imshealth.com/xfra/layer0/common/decorators"
	"gitlab.imshealth.com/xfra/layer0/common/waitutils"
	"os"
)

func GetBackend(credProvider provider.CredProvider, region, mysqlConnection, adminConnection string) (backend.Backend, error) {
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

	ecsProvider, err := ecs.NewECS(credProvider, region)
	if err != nil {
		return nil, err
	}

	elbProvider, err := elb.NewELB(credProvider, region)
	if err != nil {
		return nil, err
	}

	autoscalingProvider, err := autoscaling.NewAutoScaling(credProvider, region)
	if err != nil {
		return nil, err
	}

	cloudWatchLogsProvider, err := cloudwatchlogs.NewCloudWatchLogs(credProvider, region)
	if err != nil {
		return nil, err
	}

	// todo: data stores should take connection strings as inputs
	// todo: all data logic should use same data store
	sqlAdmin, err := getSQLAdmin(adminConnection)
	if err != nil {
		return nil, err
	}

	tagData, err := getTagData(mysqlConnection)
	if err != nil {
		return nil, err
	}

	ec2Provider = wrapEC2(ec2Provider)
	elbProvider = wrapELB(elbProvider)
	ecsProvider = wrapECS(ecsProvider)
	autoscalingProvider = wrapAutoscaling(autoscalingProvider)
	cloudWatchLogsProvider = wrapCloudWatchLogs(cloudWatchLogsProvider)

	backend := ecsbackend.NewBackend(
		sqlAdmin,
		tagData,
		s3Provider,
		iamProvider,
		ec2Provider,
		ecsProvider,
		elbProvider,
		autoscalingProvider,
		cloudWatchLogsProvider)

	return backend, nil
}

func GetLogic(backend backend.Backend, mysqlConnection, adminConnection string) (*logic.Logic, error) {
	sqlAdmin, err := getSQLAdmin(adminConnection)
	if err != nil {
		return nil, err
	}

	tagData, err := getTagData(mysqlConnection)
	if err != nil {
		return nil, err
	}

	jobData, err := getJobData(mysqlConnection)
	if err != nil {
		return nil, err
	}

	return logic.NewLogic(sqlAdmin, tagData, jobData, backend), nil
}

func getSQLAdmin(adminConnection string) (data.SQLAdmin, error) {
	var getAdminDataStore func() (data.AdminDataStore, error)
	if adminConnection != "" {
		os.Setenv("LAYER0_MYSQL_ADMIN_CONNECTION", adminConnection)
		getAdminDataStore = func() (data.AdminDataStore, error) { return data.NewMySQLAdmin() }
	} else {
		getAdminDataStore = func() (data.AdminDataStore, error) { return data.NewSQLiteAdminDataStore() }
	}

	adminDataStore, err := getAdminDataStore()
	if err != nil {
		return nil, err
	}

	return data.NewSQLAdminLayer(adminDataStore), nil
}

func getTagData(mysqlConnection string) (data.TagData, error) {
	var getTagDataStore func() (data.TagDataStore, error)
	if mysqlConnection != "" {
		os.Setenv("LAYER0_MYSQL_CONNECTION", mysqlConnection)
		getTagDataStore = func() (data.TagDataStore, error) { return data.NewTagMySQLDataStore() }
	} else {
		getTagDataStore = func() (data.TagDataStore, error) { return data.NewTagSQLiteDataStore() }
	}

	tagDataStore, err := getTagDataStore()
	if err != nil {
		return nil, err
	}

	return data.NewTagLogicLayer(tagDataStore), nil
}

func getJobData(mysqlConnection string) (data.JobData, error) {
	var getJobDataStore func() (data.JobDataStore, error)
	if mysqlConnection != "" {
		os.Setenv("LAYER0_MYSQL_CONNECTION", mysqlConnection)
		getJobDataStore = func() (data.JobDataStore, error) { return data.NewJobMySQLDataStore() }
	} else {
		getJobDataStore = func() (data.JobDataStore, error) { return data.NewJobSQLiteDataStore() }
	}

	jobDataStore, err := getJobDataStore()
	if err != nil {
		return nil, err
	}

	return data.NewJobLogicLayer(jobDataStore), nil
}

func wrapECS(e ecs.Provider) ecs.Provider {
	wrap := &ecs.ProviderDecorator{
		Inner:     e,
		Decorator: decorators.CallWithLogging,
	}

	retry := &decorators.Retry{
		Clock: waitutils.RealClock{},
	}

	wrap = &ecs.ProviderDecorator{
		Inner:     wrap,
		Decorator: retry.CallWithRetries,
	}

	return wrap
}

func wrapAutoscaling(a autoscaling.Provider) autoscaling.Provider {
	wrap := &autoscaling.ProviderDecorator{
		Inner:     a,
		Decorator: decorators.CallWithLogging,
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
