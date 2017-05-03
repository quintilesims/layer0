resource "aws_ecs_task_definition" "api" {
  family                = "l0-${var.name}-api"
  container_definitions = "${data.template_file.container_definitions.rendered}"
}

data "template_file" "container_definitions" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    api_auth_token          = "${base64encode("${var.username}:${var.password}")}"
    version                 = "${var.version}"
    access_key              = "todo"
    secret_key              = "todo"
    region                  = "${var.region}"
    public_subnets          = "todo"
    private_subnets         = "todo"
    ecs_role                = "todo"
    ecs_instance_profile    = "todo"
    vpc_id                  = "todo"
    s3_bucket               = "todo"
    linux_service_ami       = "todo"
    windows_service_ami     = "todo"
    l0_prefix               = "todo"
    agent_securitygroupid   = "todo??"
    runner_docker_image_tag = "todo"
    account_id              = "todo"
    ssh_key_pair            = "todo"
    log_group_name          = "${var.log_group}"
    dynamo_tag_table        = "${aws_dynamodb_table.tags.id}"
    dynamo_job_table        = "${aws_dynamodb_table.jobs.id}"
  }
}
