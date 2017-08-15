output "loadbalancer_ids" {
  value = "${join(",", layer0_load_balancer.base.*.id)}"
}
