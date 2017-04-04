output "environment_id" {
  value = "${layer0_environment.ws.id}"
}

output "service_id" {
  value = "${module.windows.service_id}"
}

output "service_url" {
  value = "http://${module.windows.load_balancer_url}"
}
