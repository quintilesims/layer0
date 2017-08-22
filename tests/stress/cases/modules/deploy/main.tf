resource "layer0_deploy" "base" {
  name    = "td-${var.num_deploys}-${count.index}"
  count   = "${var.num_deploys}"
  content = "${data.template_file.container.rendered}"
}

data "template_file" "container" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    deploy_command = "${var.deploy_command}"
  }
}
