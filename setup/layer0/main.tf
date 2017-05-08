provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

data "aws_caller_identity" "current" {}

module "vpc" {
  # todo: include once 'count' param supported: count = "${var.vpc_id == "" ? 1 : 0 }"

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

module "core" {
  source    = "./core"
  name      = "${var.name}"
  region    = "${var.region}"
  dockercfg = "${var.dockercfg}"
  vpc_id    = "${var.vpc_id == "" ? module.vpc.vpc_id : var.vpc_id}"
}

module "api" {
  source           = "./api"
  name             = "${var.name}"
  region           = "${var.region}"
  version          = "${var.version}"
  username         = "${var.username}"
  password         = "${var.password}"
  access_key       = "${module.core.user_access_key}"
  secret_key       = "${module.core.user_secret_key}"
  vpc_id           = "${var.vpc_id == "" ? module.vpc.vpc_id : var.vpc_id}"
  bucket_name      = "${module.core.bucket_name}"
  ssh_key_pair     = "${var.ssh_key_pair}"
  instance_profile = "${module.core.instance_profile}"
  iam_role         = "${module.core.iam_role}"
  log_group        = "${module.core.log_group}"

  tags {
    "layer0" = "${var.name}"
  }
}
