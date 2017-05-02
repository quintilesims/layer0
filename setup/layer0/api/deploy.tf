resource "aws_ecs_task_definition" "api" {
  family                = "l0-${var.name}-api"
  container_definitions = "${data.template_file.container_definitions.rendered}"
}

data "template_file" "container_definitions" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    api_auth_token          = "todo"
    api_docker_image        = "todo"
    access_key              = "todo"
    secret_key              = "todo"
    region                  = "todo"
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
    key_pair                = "todo"
    log_group_name          = "todo"
    dynamo_tag_table        = "todo"
    dynamo_job_table        = "todo"
  }
}
