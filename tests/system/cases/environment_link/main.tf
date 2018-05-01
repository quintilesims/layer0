provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "el_public" {
  name  = "el_public"
  scale = 1
}

resource "layer0_environment" "el_private" {
  name  = "el_private"
  scale = 1
}

resource "layer0_environment_links" "public_private" {
  environment_id = "${layer0_environment.el_public.id}"
  links          = ["${layer0_environment.el_private.id}"]
}

resource "layer0_environment_links" "private_public" {
  environment_id = "${layer0_environment.el_private.id}"
  links          = ["${layer0_environment.el_public.id}"]
}

module "sts_public" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.el_public.id}"
  name           = "sts_public"
}

module "sts_private" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.el_private.id}"
  name           = "sts_private"
  private        = true
}
