provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

module "environment" {
  source           = "environment"
  num_environments = "${var.num_environments}"
}

module "deploy" {
  source         = "deploy"
  num_deploys    = "${var.num_deploys}"
  deploy_command = "${var.deploy_command}"
}

module "service" {
  source          = "service"
  num_services    = "${var.num_services}"
  environment_ids = "${module.environment.environment_ids}"
  deploy_ids      = "${module.deploy.deploy_ids}"
}

module "loadbalancer" {
  source            = "loadbalancer"
  num_loadbalancers = "${var.num_loadbalancers}"
  environment_ids   = "${module.environment.environment_ids}"
}
