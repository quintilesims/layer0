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
  environment_id = "${layer0_environment.dsr.id}"
}
