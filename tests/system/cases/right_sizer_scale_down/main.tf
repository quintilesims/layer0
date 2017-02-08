variable "endpoint" {}

variable "token" {}

provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "rssd" {
  name = "rssd"
}

module "sts" {
  source         = "../modules/sts"
  endpoint       = "${var.endpoint}"
  token          = "${var.token}"
  environment_id = "${layer0_environment.rssd.id}"
  scale          = 3
}
