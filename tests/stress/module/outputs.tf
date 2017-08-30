output "random_environment" {
  value = "${element(random_shuffle.environments.result, 0)}"
}

output "random_load_balancer" {
  value = "${element(random_shuffle.load_balancers.result, 0)}"
}

output "random_deploy" {
  value = "${element(random_shuffle.deploys.result, 0)}"
}

output "random_service" {
  value = "${element(random_shuffle.services.result, 0)}"
}
