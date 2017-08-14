resource "layer0_service" "base" {
  name          = "ts-${var.num_services}-${count.index}"
  environment   = "${element(split(",", var.environment_ids), count.index)}"
  deploy        = "${element(split(",", var.deploy_ids), count.index)}"
  scale         = 1
  wait          = true
  count         = "${var.num_services}" 
}
