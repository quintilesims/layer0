output "deploy_ids" {
  value = "${join(",", layer0_deploy.base.*.id)}"
}
