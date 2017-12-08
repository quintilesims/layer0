provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "el_public" {
  name = "el_public"
}

resource "layer0_environment" "el_private" {
  name = "el_private"
}

resource "layer0_environment_link" "public_private" {
  source = "${layer0_environment.el_public.id}"
  dest   = "${layer0_environment.el_private.id}"
}

module "sts_public" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.el_public.id}"
}

module "sts_private" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.el_private.id}"
  private        = true
}
