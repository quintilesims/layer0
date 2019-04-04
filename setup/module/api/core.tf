data "aws_caller_identity" "current" {}

data "aws_subnet_ids" "public" {
  vpc_id = "${var.vpc_id}"

  tags {
    Tier = "Public"
  }
}

data "aws_subnet_ids" "private" {
  vpc_id = "${var.vpc_id}"

  tags {
    Tier = "Private"
  }
}

resource "aws_s3_bucket" "mod" {
  bucket        = "layer0-${var.name}-${data.aws_caller_identity.current.account_id}"
  region        = "${var.region}"
  force_destroy = true
}

resource "aws_s3_bucket_object" "dockercfg" {
  bucket  = "${aws_s3_bucket.mod.id}"
  key     = "bootstrap/dockercfg"
  content = "${var.dockercfg}"
}

resource "aws_cloudwatch_log_group" "mod" {
  name = "l0-${var.name}"
}

data "template_file" "ecs_assume_role_policy" {
  template = "${file("${path.module}/policies/ecs_assume_role_policy.json")}"
}

resource "aws_iam_role" "ecs" {
  name               = "l0-${var.name}-ecs-role"
  path               = "/l0/l0-${var.name}/"
  assume_role_policy = "${data.template_file.ecs_assume_role_policy.rendered}"
}

data "template_file" "ecs_role_policy" {
  template = "${file("${path.module}/policies/ecs_role_policy.json")}"

  vars {
    name       = "${var.name}"
    region     = "${var.region}"
    s3_bucket  = "${aws_s3_bucket.mod.id}"
    account_id = "${data.aws_caller_identity.current.account_id}"
  }
}

resource "aws_iam_role_policy" "ecs" {
  name   = "l0-${var.name}-ecs-role-policy"
  role   = "${aws_iam_role.ecs.id}"
  policy = "${data.template_file.ecs_role_policy.rendered}"
}

resource "aws_iam_instance_profile" "ecs" {
  name = "l0-${var.name}-ecs-instance-profile"
  path = "/l0/l0-${var.name}/"
  role = "${aws_iam_role.ecs.name}"
}

resource "aws_iam_user" "mod" {
  name = "l0-${var.name}-user"
  path = "/l0/l0-${var.name}/"
}

resource "aws_iam_access_key" "mod" {
  user = "${aws_iam_user.mod.name}"
}

resource "aws_iam_group_membership" "mod" {
  name  = "l0-${var.name}-group-membership"
  group = "${aws_iam_group.mod.name}"
  users = ["${aws_iam_user.mod.name}"]
}

resource "aws_iam_group" "mod" {
  name = "l0-${var.name}"
  path = "/l0/l0-${var.name}/"
}

data "template_file" "group_policy" {
  count    = "${length(var.group_policies)}"
  template = "${file("${path.module}/policies/${var.group_policies[count.index]}_group_policy.json")}"

  vars {
    name       = "${var.name}"
    region     = "${var.region}"
    account_id = "${data.aws_caller_identity.current.account_id}"
    s3_bucket  = "${aws_s3_bucket.mod.id}"
    vpc_id     = "${var.vpc_id}"
  }
}

resource "aws_iam_group_policy" "mod" {
  count  = "${length(var.group_policies)}"
  name   = "l0-${var.name}-${var.group_policies[count.index]}"
  group  = "${aws_iam_group.mod.id}"
  policy = "${element(data.template_file.group_policy.*.rendered, count.index)}"
}

data "aws_ami" "linux" {
  owners = ["amazon"]
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn-ami-2017.09.d-amazon-ecs-optimized"]
  }
}

data "aws_ami" "windows" {
  owners = ["amazon"]
  most_recent = true

  filter {
    name   = "name"
    values = ["Windows_Server-2016-English-Full-ECS_Optimized-2017.11.24"]
  }
}
