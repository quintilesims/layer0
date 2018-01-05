# todo: descriptions

variable "name" {}

variable "region" {}

variable "layer0_version" {}

variable "vpc_id" {}

variable "username" {}

variable "password" {}

variable "ssh_key_pair" {}

variable "dockercfg" {}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}

# Current AMI: 	amzn-ami-2016.09.g-amazon-ecs-optimized
variable "linux_region_amis" {
  default = {
    us-west-1 = "ami-689bc208"
    us-west-2 = "ami-62d35c02"
    us-east-1 = "ami-62745007"
    eu-west-1 = "ami-95f8d2f3"
  }
}

# Current AMI: Microsoft Windows Server 2016 Base with Containers
variable "windows_region_amis" {
  default = {
    us-west-1 = "ami-7c2b0e1c"
    us-west-2 = "ami-7729b917"
    us-east-1 = "ami-9667ef80"
    eu-west-1 = "ami-b9fac0df"
  }
}

variable "group_policies" {
  default = [
    "autoscaling",
    "dynamodb",
    "ec2",
    "ecs",
    "elb",
    "iam",
    "logs",
    "s3",
  ]
}
