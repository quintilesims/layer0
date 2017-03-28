variable "endpoint" {}

variable "token" {}

provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "ss" {
  name = "ss"
}

module "sts" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.ss.id}"
}

output "environment_id" {
	value = "${layer0_environment.ss.id}"
}

output "service_id" {
	value = "${module.sts.service_id}"
}
