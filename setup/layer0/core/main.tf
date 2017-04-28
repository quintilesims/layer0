data "aws_caller_identity" "current" {}

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

resource "aws_iam_instance_profile" "mod" {
  name = "l0-${var.name}"
  path = "/l0/l0-${var.name}/"
  role = "${aws_iam_role.mod.name}"
}

data "template_file" "assume_role_policy" {
  template = "${file("${path.module}/policies/assume_role_policy.json")}"
}

resource "aws_iam_role" "mod" {
  name               = "l0-${var.name}"
  path               = "/l0/l0-${var.name}/"
  assume_role_policy = "${data.template_file.assume_role_policy.rendered}"
}
