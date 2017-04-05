output "alpha_service_url" {
  value = "http://${module.sts_alpha.load_balancer_url}"
}

output "beta_service_url" {
  value = "http://${module.sts_beta.load_balancer_url}"
}
