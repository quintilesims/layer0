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
  content = "${file("Dockerrun.aws.json")}"
}
