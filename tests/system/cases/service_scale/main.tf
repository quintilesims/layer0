provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "ss" {
  name = "ss"
  scale = 2
}

module "sts" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.ss.id}"
  stateful       = true
}
