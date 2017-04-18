resource "aws_ecs_cluster" "api" {
  name = "l0-${var.name}-api"
}

/*
resource "aws_ecs_service" "api" {
  name            = "l0-${var.name}-api"
  cluster         = "${aws_ecs_cluster.api.id}"
  task_definition = "${aws_ecs_task_definition.api.arn}"
  desired_count   = 1
  iam_role        = "${var.ecs_role_arn}"
  # TODO: docs have: depends_on      = ["aws_iam_role_policy.ecs"]

  deployment_minimum_healthy_percent = 0
  deployment_minimum_healthy_percent = 200

  load_balancer {
    elb_name       = "${aws_elb.api.name}"
    container_name = "api"
    container_port = 9090
  }
}
*/

resource "aws_ecs_task_definition" "api" {
  family                = "l0-${var.name}-api"
  container_definitions = "${data.template_file.container_definitions.rendered}"
}

data "template_file" "container_definitions" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    name           = "${var.name}"
    region         = "${var.aws_region}"
    api_image_tag  = "todo"
    log_group_name = "todo"
  }
}
