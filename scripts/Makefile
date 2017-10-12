deps:
	go get -u github.com/golang/mock/mockgen
	go install github.com/golang/mock/mockgen

all: api cli common
api: provider job
cli: command client
common: aws

provider:
	mockgen github.com/quintilesims/layer0/api/provider EnvironmentProvider > ../api/provider/mock_provider/mock_environment.go &
	mockgen github.com/quintilesims/layer0/api/provider DeployProvider > ../api/provider/mock_provider/mock_deploy.go &
	mockgen github.com/quintilesims/layer0/api/provider LoadBalancerProvider > ../api/provider/mock_provider/mock_loadbalancer.go &
	mockgen github.com/quintilesims/layer0/api/provider ServiceProvider > ../api/provider/mock_provider/mock_service.go &
	mockgen github.com/quintilesims/layer0/api/provider TaskProvider > ../api/provider/mock_provider/mock_task.go &

job:
	mockgen github.com/quintilesims/layer0/api/job Store > ../api/job/mock_job/mock_store.go &

client:
	mockgen github.com/quintilesims/layer0/client Client > ../client/mock_client/mock_client.go &

command:
	mockgen github.com/quintilesims/layer0/cli/resolver Resolver > ../cli/resolver/mock_resolver/mock_resolver.go &

aws:
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface AutoScalingAPI > ../common/aws/mock_aws/mock_autoscaling.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface CloudWatchLogsAPI > ../common/aws/mock_aws/mock_cloudwatchlogs.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API > ../common/aws/mock_aws/mock_ec2.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI > ../common/aws/mock_aws/mock_ecs.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI > ../common/aws/mock_aws/mock_elb.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/iam/iamiface IAMAPI > ../common/aws/mock_aws/mock_iam.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/s3/s3iface S3API > ../common/aws/mock_aws/mock_s3.go &

.PHONY: deps all api common provider job aws
