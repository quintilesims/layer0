resource "aws_ecs_service" "api" {
  name                               = "l0-${var.name}-api"
  cluster                            = "${aws_ecs_cluster.api.id}"
  task_definition                    = "${aws_ecs_task_definition.api.arn}"
  desired_count                      = 1
  iam_role                           = "${aws_iam_role.ecs.arn}"
  depends_on                         = ["aws_iam_role_policy.ecs"]
  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  load_balancer {
    elb_name       = "${aws_elb.api.name}"
    container_name = "api"
    container_port = 9090
  }
}
