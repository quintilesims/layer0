resource "layer0_load_balancer" "base" {
  name        = "tl-${var.num_loadbalancers}-${count.index}"
  environment = "${element(split(",", var.environment_ids), count.index)}"
  private     = true
  count       = "${var.num_loadbalancers}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}
