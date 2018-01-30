output "name" {
  value = "${var.name}"
}

output "account_id" {
  value = "${data.aws_caller_identity.current.account_id}"
}

output "endpoint" {
  value = "https://${module.api.load_balancer_url}"
}

output "token" {
  value = "${module.api.token}"
}

output "s3_bucket" {
  value = "${module.api.bucket_name}"
}

output "access_key" {
  value = "${module.api.user_access_key}"
}

output "secret_key" {
  value = "${module.api.user_secret_key}"
}

output "vpc_id" {
  value = "${ var.vpc_id == "" ? module.vpc.vpc_id : var.vpc_id }"
}

output "public_subnets" {
  value = "${module.api.public_subnets}"
}

output "private_subnets" {
  value = "${module.api.private_subnets}"
}

output "ecs_role" {
  value = "${module.api.iam_role}"
}

output "ssh_key_pair" {
  value = "${var.ssh_key_pair}"
}

output "ecs_agent_instance_profile" {
  value = "${module.api.instance_profile}"
}

output "linux_service_ami" {
  value = "${module.api.linux_service_ami}"
}

output "windows_service_ami" {
  value = "${module.api.windows_service_ami}"
}

output "dynamo_tag_table" {
  value = "${module.api.dynamo_tag_table}"
}

output "dynamo_lock_table" {
  value = "${module.api.dynamo_lock_table}"
}

output "log_group_name" {
  value = "${module.api.log_group}"
}
