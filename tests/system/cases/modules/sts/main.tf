# -- stateless deployment --
resource "layer0_load_balancer" "sts_stateless" {
  count = "${var.stateless ? 1 : 0}"

  name        = "sts_stateless"
  environment = "${var.environment_id}"
  private     = "${var.private}"
  type        = "application"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "sts_stateless" {
  count = "${var.stateless ? 1 : 0}"

  name          = "sts_stateless"
  environment   = "${var.environment_id}"
  deploy        = "${layer0_deploy.sts_stateless.id}"
  load_balancer = "${layer0_load_balancer.sts_stateless.id}"
  stateful      = false
}

resource "layer0_deploy" "sts_stateless" {
  count = "${var.stateless ? 1 : 0}"

  name    = "sts_stateless"
  content = "${data.template_file.sts_stateless.rendered}"
}

data "template_file" "sts_stateless" {
  template = "${file("${path.module}/stateless.Dockerrun.aws.json")}"
}

# -- stateful deployment --
resource "layer0_load_balancer" "sts_stateful" {
  count = "${var.stateful ? 1 : 0}"

  name        = "sts_stateful"
  environment = "${var.environment_id}"
  private     = "${var.private}"
  type        = "classic"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "sts_stateful" {
  count = "${var.stateful ? 1 : 0}"

  name          = "sts_stateful"
  environment   = "${var.environment_id}"
  deploy        = "${layer0_deploy.sts_stateful.id}"
  load_balancer = "${layer0_load_balancer.sts_stateful.id}"
  scale         = "${var.scale}"
  stateful      = true
}

resource "layer0_deploy" "sts_stateful" {
  count = "${var.stateful ? 1 : 0}"

  name    = "sts_stateful"
  content = "${data.template_file.sts_stateful.rendered}"
}

data "template_file" "sts_stateful" {
  template = "${file("${path.module}/stateful.Dockerrun.aws.json")}"
}
