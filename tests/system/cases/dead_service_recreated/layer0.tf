variable "endpoint" {}

variable "token" {}

provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "dsr" {
  name = "dsr"
}

resource "layer0_load_balancer" "baxter" {
  name        = "baxter"
  environment = "${layer0_environment.dsr.id}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "baxter" {
  name          = "baxter"
  environment   = "${layer0_environment.dsr.id}"
  deploy        = "${layer0_deploy.baxter.id}"
  load_balancer = "${layer0_load_balancer.baxter.id}"
}

resource "layer0_deploy" "baxter" {
  name    = "baxter"
  content = "${data.template_file.baxter.rendered}"
}

data "template_file" "baxter" {
  template = "${file("Dockerrun.aws.json")}"
}
