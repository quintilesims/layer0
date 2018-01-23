variable "name" {}

variable "layer0_version" {}

variable "access_key" {}

variable "secret_key" {}

variable "region" {}

variable "ssh_key_pair" {}

variable "dockercfg" {}

variable "username" {}

variable "password" {}

variable "time_between_requests" {
  default = "10ms"
}

variable "vpc_id" {
  description = "optional - use an empty string to provision a new vpc"
  type        = "string"
}
