provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "datasources" {
  name = "dsrctest"
}

resource "layer0_load_balancer" "datasources_stateless" {
  name        = "dsrctest_stateless"
  environment = "${layer0_environment.datasources.id}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "datasources_stateless" {
  name          = "dsrctest_stateless"
  environment   = "${layer0_environment.datasources.id}"
  deploy        = "${layer0_deploy.datasources_stateless.id}"
  load_balancer = "${layer0_load_balancer.datasources_stateless.id}"
}

resource "layer0_deploy" "datasources_stateless" {
  name    = "dsrctest_stateless"
  content = "${file("${path.module}/stateless.dockerrun.aws.json")}"
}

resource "layer0_load_balancer" "datasources_stateful" {
  name        = "dsrctest_stateful"
  environment = "${layer0_environment.datasources.id}"
  type        = "classic"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "datasources_stateful" {
  name          = "dsrctest_stateful"
  environment   = "${layer0_environment.datasources.id}"
  deploy        = "${layer0_deploy.datasources_stateful.id}"
  load_balancer = "${layer0_load_balancer.datasources_stateful.id}"
  stateful      = true
}

resource "layer0_deploy" "datasources_stateful" {
  name    = "dsrctest_stateful"
  content = "${file("${path.module}/stateful.dockerrun.aws.json")}"
}

data "layer0_environment" "datasources" {
  depends_on = ["layer0_environment.datasources"]
  name       = "${layer0_environment.datasources.name}"
}

data "layer0_deploy" "datasources_stateless" {
  depends_on = ["layer0_deploy.datasources_stateless"]
  name       = "${layer0_deploy.datasources_stateless.name}"
  version    = "${layer0_deploy.datasources_stateless.version}"
}

data "layer0_service" "datasources_stateless" {
  depends_on     = ["layer0_service.datasources_stateless"]
  name           = "${layer0_service.datasources_stateless.name}"
  environment_id = "${layer0_environment.datasources_stateless.id}"
}

data "layer0_load_balancer" "datasources_stateless" {
  depends_on     = ["layer0_load_balancer.datasources_stateless", "layer0_service.datasources_stateless"]
  name           = "${layer0_load_balancer.datasources_stateless.name}"
  environment_id = "${layer0_environment.datasources_stateless.id}"
}

data "layer0_deploy" "datasources_stateful" {
  depends_on = ["layer0_deploy.datasources_stateful"]
  name       = "${layer0_deploy.datasources_stateful.name}"
  version    = "${layer0_deploy.datasources_stateful.version}"
}

data "layer0_service" "datasources_stateful" {
  depends_on     = ["layer0_service.datasources_stateful"]
  name           = "${layer0_service.datasources_stateful.name}"
  environment_id = "${layer0_environment.datasources_stateful.id}"
}

data "layer0_load_balancer" "datasources_stateful" {
  depends_on     = ["layer0_load_balancer.datasources_stateful", "layer0_service.datasources_stateful"]
  name           = "${layer0_load_balancer.datasources_stateful.name}"
  environment_id = "${layer0_environment.datasources_stateful.id}"
}
