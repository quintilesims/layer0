data "aws_caller_identity" "current" {}

resource "aws_s3_bucket" "core" {
  bucket        = "layer0-${var.name}-${data.aws_caller_identity.current.account_id}"
  region        = "${var.region}"
  force_destroy = true
}

resource "aws_s3_bucket_object" "dockercfg" {
  bucket  = "${aws_s3_bucket.core.id}"
  key     = "bootstrap/dockercfg"
  content = "${var.dockercfg}"
}
