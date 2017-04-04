provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "ws" {
  name = "ws"
  os   = "windows"
  size = "m3.large"
}

module "windows" {
  source         = "../modules/windows"
  environment_id = "${layer0_environment.ws.id}"
}
