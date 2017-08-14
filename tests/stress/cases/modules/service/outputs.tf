output "service_ids" {
  value = "${join(",", layer0_service.base.*.id)}"
}
