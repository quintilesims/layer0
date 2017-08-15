output "environment_ids" {
  value = "${join(",", layer0_environment.base.*.id)}"
}
