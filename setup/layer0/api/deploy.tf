resource "aws_ecs_task_definition" "api" {
  family                = "l0-${var.name}-api"
  container_definitions = "${data.template_file.container_definitions.rendered}"
}

data "template_file" "container_definitions" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    api_auth_token       = "${base64encode("${var.username}:${var.password}")}"
    version              = "${var.version}"
    access_key           = "${aws_iam_access_key.mod.id}"
    secret_key           = "${aws_iam_access_key.mod.secret}"
    region               = "${var.region}"
    public_subnets       = "${join(",", data.aws_subnet_ids.public.ids)}"
    private_subnets      = "${join(",", data.aws_subnet_ids.private.ids)}"
    ecs_role             = "${aws_iam_role.mod.id}"
    ecs_instance_profile = "${aws_iam_instance_profile.mod.id}"
    vpc_id               = "${var.vpc_id}"
    s3_bucket            = "${aws_s3_bucket.mod.id}"
    linux_service_ami    = "${lookup(var.linux_region_amis, var.region)}"
    windows_service_ami  = "${lookup(var.windows_region_amis, var.region)}"
    l0_prefix            = "${var.name}"
    account_id           = "${data.aws_caller_identity.current.account_id}"
    ssh_key_pair         = "${var.ssh_key_pair}"
    log_group_name       = "${aws_cloudwatch_log_group.mod.id}"
    dynamo_tag_table     = "${aws_dynamodb_table.tags.id}"
    dynamo_job_table     = "${aws_dynamodb_table.jobs.id}"
  }
}
