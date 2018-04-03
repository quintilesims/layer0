output "public_environment_id" {
  value = "${layer0_environment.el_public.id}"
}

output "private_environment_id" {
  value = "${layer0_environment.el_private.id}"
}

output "public_service_url" {
  value = "http://${module.sts_public.stateless_load_balancer_url}"
}

output "private_service_url" {
  value = "http://${module.sts_private.stateful_load_balancer_url}"
}
