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
