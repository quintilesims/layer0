output "environment_id" {
  value = "${layer0_environment.import.id}"
}

output "service_id" {
  value = "${module.sts.service_id}"
}

output "deploy_name" {
  value = "${module.sts.deploy_name}"
}

output "load_balancer_id" {
  value = "${module.sts.load_balancer_id}"
}
