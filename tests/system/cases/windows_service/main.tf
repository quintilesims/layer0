provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "ws" {
  name             = "ws"
  os               = "windows"
  scale            = 1
  instance_type    = "m3.large"
}

module "windows" {
  source         = "../modules/windows"
  environment_id = "${layer0_environment.ws.id}"
}
