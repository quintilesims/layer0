output "stateless_service_id" {
  value = "${layer0_service.sts_stateless.*.id}"
}

output "stateless_load_balancer_id" {
  value = "${layer0_load_balancer.sts_stateless.*.id}"
}

output "stateless_load_balancer_url" {
  value = "${layer0_load_balancer.sts_stateless.*.url}"
}

output "stateless_deploy_id" {
  value = "${layer0_deploy.sts_stateless.*.id}"
}

output "stateless_deploy_name" {
  value = "${layer0_deploy.sts_stateless.*.name}"
}

output "stateful_service_id" {
  value = "${layer0_service.sts_stateful.*.id}"
}

output "stateful_load_balancer_id" {
  value = "${layer0_load_balancer.sts_stateful.*.id}"
}

output "stateful_load_balancer_url" {
  value = "${layer0_load_balancer.sts_stateful.*.url}"
}

output "stateful_deploy_id" {
  value = "${layer0_deploy.sts_stateful.*.id}"
}

output "stateful_deploy_name" {
  value = "${layer0_deploy.sts_stateful.*.name}"
}
