provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "ws" {
  name = "ws"
  os   = "windows"
  instance_type = "m3.large"
  environment_type = "static"
}

module "windows" {
  source         = "../modules/windows"
  environment_id = "${layer0_environment.ws.id}"
}
