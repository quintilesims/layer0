output "environment_id" {
  value = "${layer0_environment.dsr.id}"
}

output "stateless_service_id" {
  value = "${module.sts_stateless.service_id}"
}

output "stateless_service_url" {
  value = "http://${module.sts_stateless.load_balancer_url}"
}

output "stateful_service_id" {
  value = "${module.sts_stateful.service_id}"
}

output "stateful_service_url" {
  value = "http://${module.sts_stateful.load_balancer_url}"
}
