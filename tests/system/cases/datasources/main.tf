provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "datasources" {
  name  = "dsrctest"
  scale = 2
}

resource "layer0_load_balancer" "datasources" {
  name        = "dsrctest"
  environment = "${layer0_environment.datasources.id}"
  type        = "application"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_deploy" "datasources" {
  name    = "dsrctest"
  content = "${file("${path.module}/stateless.dockerrun.aws.json")}"
}

resource "layer0_service" "datasources" {
  name          = "dsrctest"
  environment   = "${layer0_environment.datasources.id}"
  deploy        = "${layer0_deploy.datasources.id}"
  load_balancer = "${layer0_load_balancer.datasources.id}"
}

data "layer0_environment" "datasources" {
  depends_on = ["layer0_environment.datasources"]
  name       = "${layer0_environment.datasources.name}"
}

data "layer0_load_balancer" "datasources" {
  depends_on     = ["layer0_load_balancer.datasources", "layer0_service.datasources"]
  name           = "${layer0_load_balancer.datasources.name}"
  environment_id = "${layer0_environment.datasources.id}"
}

data "layer0_deploy" "datasources" {
  depends_on = ["layer0_deploy.datasources"]
  name       = "${layer0_deploy.datasources.name}"
  version    = "${layer0_deploy.datasources.version}"
}

data "layer0_service" "datasources" {
  depends_on     = ["layer0_service.datasources"]
  name           = "${layer0_service.datasources.name}"
  environment_id = "${layer0_environment.datasources.id}"
}
