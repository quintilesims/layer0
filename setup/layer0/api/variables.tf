# todo: descriptions

variable "name" {}

variable "region" {}

variable "version" {}

variable "vpc_id" {}

variable "username" {}

variable "password" {}

variable "ssh_key_pair" {}

variable "dockercfg" {}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}

# Current AMI: amzn-ami-2016.03.c-amazon-ecs-optimized
variable "linux_region_amis" {
  default = {
    us-west-1 = "ami-bb473cdb"
    us-west-2 = "ami-84b44de4"
    us-east-1 = "ami-8f7687e2"
    eu-west-1 = "ami-4e6ffe3d"
  }
}

# Current AMI: Microsoft Windows Server 2016 Base with Containers
variable "windows_region_amis" {
  default = {
    us-west-1 = "ami-4699ca26"
    us-west-2 = "ami-7a803d1a"
    us-east-1 = "ami-e7b755f1"
    eu-west-1 = "ami-eef4de9d"
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
  ]
}
