output "load_balancer_url" {
  value = "${aws_elb.api.dns_name}"
}

output "token" {
  value = "${base64encode("${var.username}:${var.password}")}"
}

output "public_subnets" {
  value = "${join(",", data.aws_subnet_ids.public.ids)}"
}

output "private_subnets" {
  value = "${join(",", data.aws_subnet_ids.private.ids)}"
}

output "linux_service_ami" {
  value = "${lookup(var.linux_region_amis, var.region)}"
}

output "windows_service_ami" {
  value = "${lookup(var.windows_region_amis, var.region)}"
}

output "bucket_name" {
  value = "${aws_s3_bucket.mod.id}"
}

output "instance_profile" {
  value = "${aws_iam_instance_profile.ecs.id}"
}

output "iam_role" {
  value = "${aws_iam_role.ecs.id}"
}

output "log_group" {
  value = "${aws_cloudwatch_log_group.mod.id}"
}

output "user_access_key" {
  value = "${aws_iam_access_key.mod.id}"
}

output "user_secret_key" {
  value = "${aws_iam_access_key.mod.secret}"
}

output "dynamo_tag_table" {
  value = "${aws_dynamodb_table.tags.id}"
}

output "dynamo_job_table" {
  value = "${aws_dynamodb_table.jobs.id}"
}
