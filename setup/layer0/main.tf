provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}

module "vpc" {
  source             = "github.com/terraform-community-modules/tf_aws_vpc"
  name               = "l0-${var.name}"
  cidr               = "10.100.0.0/16"
  private_subnets    = ["10.100.64.0/18", "10.100.192.0/18"]
  public_subnets     = ["10.100.0.0/18", "10.100.128.0/18"]
  azs                = ["${var.aws_region}a", "${var.aws_region}b"]
  enable_nat_gateway = true

  tags {
    "layer0" = "${var.name}"
  }
}

module "api" {
  source         = "./api"
  name           = "${var.name}"
  aws_region     = "${var.aws_region}"
  vpc            = "${module.vpc.vpc_id}"
  public_subnets = ["${module.vpc.public_subnets}"]

  # todo: don't hardcode
  ecs_role_arn     = "todo"
  s3_bucket        = "layer0-zpatrick-856306994068"
  key_pair         = "xfra-dev"
  instance_profile = "l0-zpatrick-ecs-instance-profile"
  ami              = "ami-84b44de4"

  tags {
    "layer0" = "${var.name}"
  }
}
