output "instance" {
  value = "${var.name}"
}

output "aws-account-id" {
  value = "${data.aws_caller_identity.current.account_id}"
}

output "endpoint" {
  value = "https://${module.api.load_balancer_url}"
}

output "token" {
  value = "${module.api.token}"
}

output "aws-s3-bucket" {
  value = "${module.api.bucket_name}"
}

output "aws-access-key" {
  value = "${module.api.user_access_key}"
}

output "aws-secret-key" {
  value = "${module.api.user_secret_key}"
}

output "aws-vpc" {
  value = "${ var.vpc_id == "" ? module.vpc.vpc_id : var.vpc_id }"
}

output "aws-public-subnets" {
  value = "${module.api.public_subnets}"
}

output "aws-private-subnets" {
  value = "${module.api.private_subnets}"
}

output "aws-ecs-role" {
  value = "${module.api.iam_role}"
}

output "aws-ssh-key" {
  value = "${var.ssh_key_pair}"
}

output "aws-instance-profile" {
  value = "${module.api.instance_profile}"
}

output "aws-linux-ami" {
  value = "${module.api.linux_service_ami}"
}

output "aws-windows-ami" {
  value = "${module.api.windows_service_ami}"
}

output "aws-tag-table" {
  value = "${module.api.dynamo_tag_table}"
}

output "aws-job-table" {
  value = "${module.api.dynamo_job_table}"
}

output "aws-lock-table" {
  value = "${module.api.dynamo_lock_table}"
}

output "aws-log-group" {
  value = "${module.api.log_group}"
}
