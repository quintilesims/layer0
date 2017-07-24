output "service_id" {
  value = "${layer0_service.sts.id}"
}

output "load_balancer_id" {
  value = "${layer0_load_balancer.sts.id}"
}

output "load_balancer_url" {
  value = "${layer0_load_balancer.sts.url}"
}

output "deploy_id" {
  value = "${layer0_deploy.sts.id}"
}

output "deploy_name" {
  value = "${layer0_deploy.sts.name}"
}

