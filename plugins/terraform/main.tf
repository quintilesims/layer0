variable "endpoint" {}

variable "token" {}

provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "test" {
  name = "test"
}
