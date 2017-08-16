deps:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen

all: api common
api: provider scheduler
common: aws db

provider:
	mockgen github.com/quintilesims/layer0/api/provider EnvironmentProvider > ../api/provider/mock_provider/mock_environment.go &
	mockgen github.com/quintilesims/layer0/api/provider JobProvider > ../api/provider/mock_provider/mock_job.go &
	mockgen github.com/quintilesims/layer0/api/provider DeployProvider > ../api/provider/mock_provider/mock_deploy.go &
	mockgen github.com/quintilesims/layer0/api/provider LoadBalancerProvider > ../api/provider/mock_provider/mock_loadbalancer.go &

scheduler:
	mockgen github.com/quintilesims/layer0/api/scheduler JobScheduler > ../api/scheduler/mock_scheduler/mock_job_scheduler.go &

aws:
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface AutoScalingAPI > ../common/aws/mock_aws/mock_autoscaling.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API > ../common/aws/mock_aws/mock_ec2.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI > ../common/aws/mock_aws/mock_ecs.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/s3/s3iface S3API > ../common/aws/mock_aws/mock_s3.go &

db:
	mockgen github.com/quintilesims/layer0/common/db/job_store JobStore > ../common/db/job_store/mock_job_store/mock_job_store.go &

.PHONY: deps all api common provider scheduler aws db
