variable "endpoint" {}

variable "token" {}

variable "environment_id" {}

variable "scale" {
  default = 1
}

provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_load_balancer" "sts" {
  name        = "sts"
  environment = "${var.environment_id}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "layer0_service" "sts" {
  name          = "sts"
  environment   = "${var.environment_id}"
  deploy        = "${layer0_deploy.sts.id}"
  load_balancer = "${layer0_load_balancer.sts.id}"
  scale         = "${var.scale}"
  wait          = true
}

resource "layer0_deploy" "sts" {
  name    = "sts"
  content = "${data.template_file.sts.rendered}"
}

data "template_file" "sts" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"
}
