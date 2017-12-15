provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "aws_caller_identity" "current" {}

module "vpc" {
  # todo: count_hack is workaround for https://github.com/hashicorp/terraform/issues/953
  count_hack = "${ var.vpc_id == "" ? 1 : 0 }"

  source          = "./vpc"
  name            = "${var.name}"
  cidr            = "10.100.0.0/16"
  private_subnets = ["10.100.1.0/24", "10.100.2.0/24", "10.100.3.0/24"]
  public_subnets  = ["10.100.101.0/24", "10.100.102.0/24", "10.100.103.0/24"]
  azs             = ["${var.region}a", "${var.region}b", "${var.region}c"]

  tags {
    "layer0" = "${var.name}"
  }
}

module "api" {
  source   = "./api"
  name     = "${var.name}"
  region   = "${var.region}"
  version  = "${var.version}"
  username = "${var.username}"
  password = "${var.password}"

  # todo: format hack is a workaround for https://github.com/hashicorp/terraform/issues/14399
  vpc_id = "${ var.vpc_id == "" ? format("%s", module.vpc.vpc_id) : var.vpc_id }"

  ssh_key_pair  = "${var.ssh_key_pair}"
  dockercfg     = "${var.dockercfg}"
  request_delay = "${var.request_delay}"

  tags {
    "layer0" = "${var.name}"
  }
}
