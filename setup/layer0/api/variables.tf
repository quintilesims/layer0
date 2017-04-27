# todo: descriptions

variable "name" {}

variable "version" {}

variable "vpc_id" {}

variable "bucket_name" {}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}
