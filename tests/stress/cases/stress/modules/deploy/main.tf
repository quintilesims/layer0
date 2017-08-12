resource "layer0_deploy" "base" {
   name    = "td-${var.num_deploys}-${count.index}"
   content = "${file("Dockerrun.aws.json")}"

   count = "${var.num_deploys}"
}
