output "environment_id" {
  value = "${layer0_environment.dsr.id}"
}

output "service_id" {
  value = "${module.sts.service_id}"
}

output "service_url" {
  value = "http://${module.sts.load_balancer_url}"
}
