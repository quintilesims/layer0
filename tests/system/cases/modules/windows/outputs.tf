output "service_id" {
  value = "${layer0_service.windows.id}"
}

output "load_balancer_id" {
  value = "${layer0_load_balancer.windows.id}"
}

output "load_balancer_url" {
  value = "${layer0_load_balancer.windows.url}"
}

output "deploy_id" {
  value = "${layer0_deploy.windows.id}"
}
