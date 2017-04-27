# todo: descriptions

variable "name" {}

variable "region" {}

variable "version" {}

variable "vpc_id" {}

variable "bucket_name" {}

variable "ssh_key_pair" {}

variable "instance_profile" {}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}

# amzn-ami-2016.03.c-amazon-ecs-optimized
variable "service_amis" {
  default = {
    us-west-2 = "ami-84b44de4"
    us-west-1 = "ami-bb473cdb"
    us-east-1 = "ami-8f7687e2"
    eu-west-1 = "ami-4e6ffe3d"
  }
}
