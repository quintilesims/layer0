provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

module "vpc" {
  # todo: include once 'count' param supported: count = "${var.vpc_id == "" ? 1 : 0 }"

  source             = "./vpc"
  name               = "l0-${var.name}"
  cidr               = "10.100.0.0/16"
  private_subnets    = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  azs                = ["${var.region}a", "${var.region}b", "${var.region}c"]
  enable_nat_gateway = "true"

  tags {
    "layer0" = "${var.name}"
  }
}

module "core" {
  source = "./core"
  name = "${var.name}"
  region = "${var.region}"
  dockercfg = "${var.dockercfg}"
}

module "api" {
  source  = "./api"
  name    = "${var.name}"
  version = "${var.version}"
  vpc_id  = "${var.vpc_id == "" ? module.vpc.vpc_id : var.vpc_id }"
  bucket_name = "${module.core.bucket_name}"

  tags {
    "layer0" = "${var.name}"
  }
}
