variable "name" {
  description = "Name of your Layer0 instance"
}

variable "aws_region" {
  description = "AWS region"
}

variable "vpc" {
  description = "ID of your Layer0 VPC"
}

variable "ecs_role_arn" {
  description = "ARN of your Layer0 ECS role"
}

variable "s3_bucket" {
  description = "S3 bucket of your Layer0 instance"
}

variable "key_pair" {
  description = "Key pair for your Layer0 instance"
}

variable "instance_profile" {
  description = ""
}

variable "ami" {
  description = ""
}

variable "public_subnets" {
  type        = "list"
  description = "A list of public subnets inside the VPC"
  default     = []
}

variable "tags" {
  type        = "map"
  description = "A map of tags to add to all resources"
  default     = {}
}
