provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "import" {
  name = "import"
}

module "sts" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.import.id}"
}
