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

module "sts" {
  source         = "../modules/sts"
  endpoint       = "${var.endpoint}"
  token          = "${var.token}"
  environment_id = "${layer0_environment.dsr.id}"
}

output "environment_id" {
        value = "${layer0_environment.dsr.id}"
}

output "service_id" {
        value = "${module.sts.service_id}"
}

output "service_url" {
        value = "http://${module.sts.load_balancer_url}"
}
