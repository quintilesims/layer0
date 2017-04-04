output "load_balancer_url_alpha" {
        value = "${{module.sts_alpha.load_balancer_url}"
}

output "load_balancer_url_beta" {
        value = "${module.sts_beta.load_balancer_url}"
}
