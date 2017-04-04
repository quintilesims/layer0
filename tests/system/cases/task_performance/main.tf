provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}

resource "layer0_environment" "tp" {
  name = "tp"
  size = "t2.micro"
}

resource "layer0_deploy" "alpine" {
  name    = "alpine"
  content = "${data.template_file.alpine.rendered}"
}

data "template_file" "alpine" {
  template = "${file("Dockerrun.aws.json")}"
}
