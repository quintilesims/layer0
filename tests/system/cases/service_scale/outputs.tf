output "environment_id" {
  value = "${layer0_environment.ss.id}"
}

output "service_id" {
  value = "${module.sts.stateful_service_id}"
}
