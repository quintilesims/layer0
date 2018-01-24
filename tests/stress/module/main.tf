provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "random_pet" "environment_names" {
  length = 1
  count  = "${var.num_environments}"
}

resource "layer0_environment" "te" {
  name          = "${element(random_pet.environment_names.*.id, count.index)}"
  instance_type = "t2.micro"
  count         = "${var.num_environments}"
}

resource "random_pet" "load_balancer_names" {
  length = 1
  count  = "${var.num_load_balancers}"
}

resource "layer0_load_balancer" "tlb" {
  name        = "${element(random_pet.load_balancer_names.*.id, count.index)}"
  environment = "${element(layer0_environment.te.*.id, count.index)}"
  count       = "${var.num_load_balancers}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "random_pet" "deploy_families" {
  length = 1
  count  = "${var.num_deploy_families == 0 ? 1 : var.num_deploy_families}"
}

resource "layer0_deploy" "td" {
  name    = "${element(random_pet.deploy_families.*.id, count.index)}"
  content = "${file("${path.module}/Dockerrun.aws.json")}"
  count   = "${var.num_deploys}"
}

resource "random_pet" "service_names" {
  length = 1
  count  = "${var.num_services}"
}

resource "layer0_service" "ts" {
  name        = "${element(random_pet.service_names.*.id, count.index)}"
  environment = "${element(layer0_environment.te.*.id, count.index)}"
  deploy      = "${element(layer0_deploy.td.*.id, count.index)}"
  scale       = 1
  count       = "${var.num_services}"
}
