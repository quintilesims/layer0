package config

import (
	"os"
)

const (
	TEST_PREFIX                          = "l0"
	TEST_AWS_PRIVATE_SUBNETS             = "subnet-12345678,subnet-87654321"
	TEST_AWS_PUBLIC_SUBNETS              = "subnet-11111111,subnet-22222222"
	TEST_AWS_VPC_ID                      = "vpc-12345678"
	TEST_AWS_ECS_INSTANCE_PROFILE        = "l0-test-vpc-ECSInstanceProfile-123456789ABC"
	TEST_AWS_S3_BUCKET                   = "layer0-l0-123456789ABC"
	TEST_AWS_SERVICE_AMI                 = "ami-abc123"
	TEST_AWS_ECS_ROLE                    = "role-abc123"
	TEST_AWS_KEY_PAIR                    = "test-key-pair"
	TEST_AWS_ECS_AGENT_SECURITY_GROUP_ID = "agent-abc123"
)

func SetTestConfig() {
	os.Setenv(PREFIX, TEST_PREFIX)
	os.Setenv(AWS_PRIVATE_SUBNETS, TEST_AWS_PRIVATE_SUBNETS)
	os.Setenv(AWS_PUBLIC_SUBNETS, TEST_AWS_PUBLIC_SUBNETS)
	os.Setenv(AWS_VPC_ID, TEST_AWS_VPC_ID)
	os.Setenv(AWS_ECS_INSTANCE_PROFILE, TEST_AWS_ECS_INSTANCE_PROFILE)
	os.Setenv(AWS_S3_BUCKET, TEST_AWS_S3_BUCKET)
	os.Setenv(AWS_LINUX_SERVICE_AMI, TEST_AWS_SERVICE_AMI)
	os.Setenv(AWS_WINDOWS_SERVICE_AMI, TEST_AWS_SERVICE_AMI)
	os.Setenv(AWS_ECS_ROLE, TEST_AWS_ECS_ROLE)
	os.Setenv(AWS_SSH_KEY_PAIR, TEST_AWS_KEY_PAIR)
	os.Setenv(AWS_ECS_AGENT_SECURITY_GROUP_ID, TEST_AWS_ECS_AGENT_SECURITY_GROUP_ID)
}
