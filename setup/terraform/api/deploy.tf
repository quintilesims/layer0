resource "aws_ecs_task_definition" "api" {
  family                = "l0-${var.layer0_instance_name}-api"
  container_definitions = "${data.template_file.api.rendered}"

  tags {
    "layer0" = "${var.layer0_instance_name}"
  }
}

data "template_file" "api" {
  template = "${file("templates/Dockerrun.aws.json")}"
  vars     = {}
}
