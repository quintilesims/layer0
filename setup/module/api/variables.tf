# todo: descriptions

variable "name" {}

variable "region" {}

variable "version" {}

variable "vpc_id" {}

variable "username" {}

variable "password" {}

variable "ssh_key_pair" {}

variable "dockercfg" {}

variable "time_between_requests" {}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
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
