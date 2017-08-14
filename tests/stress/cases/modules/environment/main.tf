resource "layer0_environment" "base" {
  size = "t2.micro"
  name = "te-${var.num_environments}-${count.index}"
  count = "${var.num_environments}"
}
