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
  name  = "${element(random_pet.environment_names.*.id, count.index)}"
  size  = "t2.micro"
  count = "${var.num_environments}"
}

resource "random_shuffle" "environments" {
  input = ["${layer0_environment.te.*.id}"]
  count = "${var.num_environments == 0 ? 0 : 1}"
}

resource "random_pet" "load_balancer_names" {
  length = 1
  count  = "${var.num_load_balancers}"
}

resource "layer0_load_balancer" "tlb" {
  name        = "${element(random_pet.load_balancer_names.*.id, count.index)}"
  environment = "${element(random_shuffle.environments.result, count.index)}"
  count       = "${var.num_load_balancers}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "random_shuffle" "load_balancers" {
  input = ["${layer0_load_balancer.tlb.*.id}"]
  count = "${var.num_load_balancers == 0 ? 0 : 1}"
}

resource "random_pet" "deploy_families" {
  length = 1
  count  = "${var.num_deploy_families == 0 ? 1 : var.num_deploy_families}"
}

resource "random_shuffle" "deploy_families" {
  input = ["${random_pet.deploy_families.*.id}"]
}

data "template_file" "deploy" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    deploy_command = "${var.deploy_command}"
  }
}

resource "layer0_deploy" "td" {
  name    = "${element(random_shuffle.deploy_families.result, count.index)}"
  content = "${data.template_file.deploy.rendered}"
  count   = "${var.num_deploys}"
}

resource "random_shuffle" "deploys" {
  input = ["${layer0_deploy.td.*.id}"]
  count = "${var.num_deploys == 0 ? 0 : 1}"
}

resource "random_pet" "service_names" {
  length = 1
  count  = "${var.num_services}"
}

resource "layer0_service" "ts" {
  name        = "${element(random_pet.service_names.*.id, count.index)}"
  environment = "${element(random_shuffle.environments.result, count.index)}"
  deploy      = "${element(random_shuffle.deploys.result, count.index)}"
  scale       = 1
  count       = "${var.num_services}"
}

resource "random_shuffle" "services" {
  input = ["${layer0_service.ts.*.id}"]
  count = "${var.num_services == 0 ? 0 : 1}"
}
