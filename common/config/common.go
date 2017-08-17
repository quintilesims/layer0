package config

const (
	FLAG_PORT                 = "port"
	FLAG_DEBUG                = "debug"
	FLAG_INSTANCE             = "instance"
	FLAG_AWS_ACCOUNT_ID       = "account-id"
	FLAG_AWS_ACCESS_KEY       = "access-key"
	FLAG_AWS_SECRET_KEY       = "secret-key"
	FLAG_AWS_REGION           = "region"
	FLAG_AWS_VPC              = "vpc"
	FLAG_AWS_LINUX_AMI        = "linux-ami"
	FLAG_AWS_WINDOWS_AMI      = "windows-ami"
	FLAG_AWS_S3_BUCKET        = "s3-bucket"
	FLAG_AWS_INSTANCE_PROFILE = "instance-profile"
	FLAG_AWS_PUBLIC_SUBNETS   = "public-subnets"
	FLAG_AWS_PRIVATE_SUBNETS  = "private-subnets"
	FLAG_AWS_DYNAMO_TAG_TABLE = "tag-table"
)

const (
	ENVVAR_PORT                 = "LAYER0_PORT"
	ENVVAR_DEBUG                = "LAYER0_DEBUG"
	ENVVAR_INSTANCE             = "LAYER0_PREFIX"
	ENVVAR_AWS_ACCOUNT_ID       = "LAYER0_AWS_ACCOUNT_ID"
	ENVVAR_AWS_ACCESS_KEY       = "LAYER0_AWS_ACCESS_KEY_ID"
	ENVVAR_AWS_SECRET_KEY       = "LAYER0_AWS_SECRET_ACCESS_KEY"
	ENVVAR_AWS_REGION           = "LAYER0_AWS_REGION"
	ENVVAR_AWS_VPC              = "LAYER0_AWS_VPC_ID"
	ENVVAR_AWS_LINUX_AMI        = "LAYER0_AWS_LINUX_SERVICE_AMI"
	ENVVAR_AWS_WINDOWS_AMI      = "LAYER0_AWS_WINDOWS_SERVICE_AMI"
	ENVVAR_AWS_S3_BUCKET        = "LAYER0_AWS_S3_BUCKET"
	ENVVAR_AWS_INSTANCE_PROFILE = "LAYER0_AWS_ECS_INSTANCE_PROFILE"
	ENVVAR_AWS_PUBLIC_SUBNETS   = "LAYER0_AWS_PUBLIC_SUBNETS"
	ENVVAR_AWS_PRIVATE_SUBNETS  = "LAYER0_AWS_PRIVATE_SUBNETS"
	ENVVAR_AWS_DYNAMO_TAG_TABLE = "LAYER0_AWS_DYNAMO_TAG_TABLE"
)
