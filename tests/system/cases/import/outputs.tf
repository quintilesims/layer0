output "environment_id" {
  value = "${layer0_environment.import.id}"
}

output "service_id" {
  value = "${module.sts.stateless_service_id}"
}

output "deploy_name" {
  value = "${module.sts.stateless_deploy_name}"
}

output "load_balancer_id" {
  value = "${module.sts.stateless_load_balancer_id}"
}
