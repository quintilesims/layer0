resource "aws_ecs_service" "api" {
  name                               = "l0-${var.layer0_instance_name}-api"
  cluster                            = "${aws_ecs_cluster.api.id}"
  task_definition                    = "${aws_ecs_task_definition.api.arn}"
  desired_count                      = 1
  deployment_minimum_healthy_percent = 0
  iam_role                           = "${aws_iam_role.ecs.arn}"
  depends_on                         = ["aws_iam_role_policy.ecs"]

  load_balancer {
    elb_name       = "${aws_elb.api.name}"
    container_name = "l0-api"
    container_port = 9090
  }

  tags {
    "layer0" = "${var.layer0_instance_name}"
  }
}

resource "aws_cloudwatch_log_group" "l0" {
  name = "l0-${var.layer0_instance_name}"

   tags {
    "layer0" = "${var.layer0_instance_name}"
  }
}
