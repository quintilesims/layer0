output "public_service_url" {
  value = "http://${module.sts_public.load_balancer_url}"
}

output "private_service_url" {
  value = "http://${module.sts_private.load_balancer_url}"
}
