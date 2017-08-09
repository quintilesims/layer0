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
