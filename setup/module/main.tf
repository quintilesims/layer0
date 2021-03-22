provider "aws" {
  version    = ">= 2.0.0, <= 3.33.0"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "aws_caller_identity" "current" {}

module "vpc" {
  # todo: count_hack is workaround for https://github.com/hashicorp/terraform/issues/953
  count_hack = "${ var.vpc_id == "" ? 1 : 0 }"

  source = "./vpc"
  name   = "${var.name}"
  cidr   = "10.100.0.0/16"

  tags {
    "layer0" = "${var.name}"
  }
}

module "api" {
  source         = "./api"
  name           = "${var.name}"
  region         = "${var.region}"
  layer0_version = "${var.layer0_version}"
  username       = "${var.username}"
  password       = "${var.password}"
  docker_registry = "${var.docker_registry}"

  # todo: format hack is a workaround for https://github.com/hashicorp/terraform/issues/14399
  vpc_id = "${ var.vpc_id == "" ? format("%s", module.vpc.vpc_id) : var.vpc_id }"

  ssh_key_pair = "${var.ssh_key_pair}"
  dockercfg    = "${var.dockercfg}"

  tags {
    "layer0" = "${var.name}"
  }
}
