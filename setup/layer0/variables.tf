# todo: descriptions

variable "name" {}

variable "version" {}

variable "access_key" {}

variable "secret_key" {}

variable "region" {}

variable "ssh_key_pair" {}

variable "dockercfg" {}

variable "vpc_id" {
  description = "optional - use an empty string to provision a new vpc"
}
