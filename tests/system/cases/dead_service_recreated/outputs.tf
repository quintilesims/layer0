output "environment_id" {
  value = "${layer0_environment.dsr.id}"
}

output "stateless_service_id" {
  value = "${module.sts.stateless_service_id}"
}

output "stateless_service_url" {
  value = "http://${module.sts.stateless_load_balancer_url}"
}

output "stateful_service_id" {
  value = "${module.sts.stateful_service_id}"
}

output "stateful_service_url" {
  value = "http://${module.sts.stateful_load_balancer_url}"
}
