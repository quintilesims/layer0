deps:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen

all: api common
api: entity
common: aws db

entity:
	mockgen github.com/quintilesims/layer0/api/entity Provider > ../api/entity/mock_entity/mock_provider.go &
	mockgen github.com/quintilesims/layer0/api/entity Environment > ../api/entity/mock_entity/mock_environment.go &
	mockgen github.com/quintilesims/layer0/api/entity Job > ../api/entity/mock_entity/mock_job.go &

aws:
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API > ../common/aws/mock_aws/mock_ec2.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI > ../common/aws/mock_aws/mock_ecs.go &
	mockgen -package mock_aws github.com/aws/aws-sdk-go/service/s3/s3iface S3API > ../common/aws/mock_aws/mock_s3.go &

db:
	mockgen github.com/quintilesims/layer0/common/db/job_store JobStore > ../common/db/job_store/mock_job_store/mock_job_store.go &

.PHONY: deps all api common entity aws db
