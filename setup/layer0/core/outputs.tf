output "bucket_name" {
  value = "${aws_s3_bucket.mod.id}"
}

output "instance_profile" {
  value = "${aws_iam_instance_profile.mod.id}"
}

output "iam_role" {
  value = "${aws_iam_role.mod.id}"
}

output "log_group" {
  value = "${aws_cloudwatch_log_group.mod.id}"
}
