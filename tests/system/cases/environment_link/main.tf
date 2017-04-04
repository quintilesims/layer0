provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "el_alpha" {
  name = "el_alpha"
}

resource "layer0_environment" "el_beta" {
  name = "el_beta"
}

module "sts_alpha" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.el_alpha.id}"
}

module "sts_beta" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.el_beta.id}"
  private = true
}
