provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "datasources" {
  name = "datasources"
}

module "sts" {
  source         = "../modules/sts"
  environment_id = "${layer0_environment.datasources.id}"
}

data "layer0_environment" "datasources" {
  depends_on = ["layer0_environment.datasources"]
  name       = "datasources"
}

data "layer0_deploy" "datasources" {
  depends_on = ["module.sts"]
  name       = "${module.sts.deploy_name}"
  version    = "1"
}

data "layer0_service" "datasources" {
  depends_on     = ["module.sts"]
  name           = "${module.sts.service_id}"
  environment_id = "${layer0_environment.datasources.id}"
}

data "layer0_load_balancer" "datasources" {
  depends_on     = ["module.sts"]
  name           = "${module.sts.load_balancer_id}"
  environment_id = "${layer0_environment.datasources.id}"
}
