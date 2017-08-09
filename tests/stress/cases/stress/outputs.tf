output "environment_ids" {
  value = "${join(",", layer0_environment.te.*.id)}"
}

output "deploy_ids" {
  value = "${join(",", layer0_deploy.td.*.id)}"
}

output "service_ids" {
  value = "${join(",", layer0_service.ts.*.id)}"
}
