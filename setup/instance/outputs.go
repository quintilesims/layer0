package instance

import (
	"github.com/quintilesims/layer0/setup/terraform"
)

const (
	OUTPUT_NAME                        = "name"
	OUTPUT_ENDPOINT                    = "endpoint"
	OUTPUT_TOKEN                       = "token"
	OUTPUT_S3_BUCKET                   = "s3_bucket"
	OUTPUT_ACCOUNT_ID                  = "account_id"
	OUTPUT_ACCESS_KEY                  = "access_key"
	OUTPUT_SECRET_KEY                  = "secret_key"
	OUTPUT_VPC_ID                      = "vpc_id"
	OUTPUT_PRIVATE_SUBNETS             = "private_subnets"
	OUTPUT_PUBLIC_SUBNETS              = "public_subnets"
	OUTPUT_ECS_ROLE                    = "ecs_role"
	OUTPUT_SSH_KEY_PAIR                = "ssh_key_pair"
	OUTPUT_ECS_AGENT_SECURITY_GROUP_ID = "ecs_agent_security_group_id"
	OUTPUT_ECS_INSTANCE_PROFILE        = "ecs_agent_instance_profile"
	OUTPUT_AWS_LINUX_SERVICE_AMI       = "linux_service_ami"
	OUTPUT_WINDOWS_SERVICE_AMI         = "windows_service_ami"
)

// todo: fill these out
var Layer0ModuleOutputs = map[string]terraform.Output{
	OUTPUT_NAME:                        {Value: "${module.layer0.name}"},
	OUTPUT_ENDPOINT:                    {Value: "TODO!"},
	OUTPUT_TOKEN:                       {Value: "TODO!"},
	OUTPUT_S3_BUCKET:                   {Value: "TODO!"},
	OUTPUT_ACCOUNT_ID:                  {Value: "TODO!"},
	OUTPUT_ACCESS_KEY:                  {Value: "TODO!"},
	OUTPUT_SECRET_KEY:                  {Value: "TODO!"},
	OUTPUT_VPC_ID:                      {Value: "TODO!"},
	OUTPUT_PRIVATE_SUBNETS:             {Value: "TODO!"},
	OUTPUT_PUBLIC_SUBNETS:              {Value: "TODO!"},
	OUTPUT_ECS_ROLE:                    {Value: "TODO!"},
	OUTPUT_SSH_KEY_PAIR:                {Value: "TODO!"},
	OUTPUT_ECS_AGENT_SECURITY_GROUP_ID: {Value: "TODO!"},
	OUTPUT_ECS_INSTANCE_PROFILE:        {Value: "TODO!"},
	OUTPUT_AWS_LINUX_SERVICE_AMI:       {Value: "TODO!"},
	OUTPUT_WINDOWS_SERVICE_AMI:         {Value: "TODO!"},
}
