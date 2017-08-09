provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "te" {
  name = "te-${var.num_environments}-${count.index}"
  size = "t2.micro"

  count = "${var.num_environments}"
}

resource "layer0_service" "ts" {
  name          = "ts-${var.num_services}-${count.index}"
  environment   = "${layer0_environment.te.id}"
  deploy        = "${layer0_deploy.td.id}"
  scale         = 1
  wait          = true

  count = "${var.num_services}"
}

resource "layer0_deploy" "td" {
   name    = "td-${var.num_deploys}-${count.index}"
   content = "${file("Dockerrun.aws.json")}"

   count = "${var.num_deploys}"
}
