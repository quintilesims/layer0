provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

module "environment" {
  source = "modules/environment"

  num_environments = "${var.num_environments}"
}


module "deploy" {
  source = "modules/deploy"

  num_deploys = "${var.num_deploys}"
}

module "service" {
  source = "modules/service"
  num_services = "${var.num_services}"
  environment_ids = "${module.environment.environment_ids}"
  deploy_ids = "${module.deploy.deploy_ids}"
}

