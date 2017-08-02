output "environment_id" {
  value = "${layer0_environment.tp.id}"
}

output "deploy_id" {
  value = "${layer0_deploy.alpine.id}"
}
