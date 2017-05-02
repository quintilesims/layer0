variable "name" {}

variable "region" {}

variable "dockercfg" {}

variable "vpc_id" {}

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
