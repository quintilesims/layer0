resource "layer0_load_balancer" "windows" {
  name        = "windows"
  environment = "${var.environment_id}"
  type        = "classic"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "windows" {
  name          = "windows"
  environment   = "${var.environment_id}"
  deploy        = "${layer0_deploy.windows.id}"
  load_balancer = "${layer0_load_balancer.windows.id}"
  scale         = "${var.scale}"
  stateful      = true
}

resource "layer0_deploy" "windows" {
  name    = "windows"
  content = "${file("${path.module}/Dockerrun.aws.json")}"
}
