output "environment_id" {
  value = "${data.layer0_environment.datasources.id}"
}

output "environment_size" {
  value = "${data.layer0_environment.datasources.size}"
}

output "environment_min_count" {
  value = "${data.layer0_environment.datasources.min_count}"
}

output "environment_os" {
  value = "${data.layer0_environment.datasources.os}"
}

output "environment_ami" {
  value = "${data.layer0_environment.datasources.ami}"
}

output "deploy_id" {
  value = "${data.layer0_deploy.datasources.id}"
}

output "load_balancer_id" {
  value = "${data.layer0_load_balancer.datasources.id}"
}

output "load_balancer_name" {
  value = "${data.layer0_load_balancer.datasources.name}"
}

output "load_balancer_environment_name" {
  value = "${data.layer0_load_balancer.datasources.environment_name}"
}

output "load_balancer_private" {
  value = "${data.layer0_load_balancer.datasources.private}"
}

output "load_balancer_url" {
  value = "${data.layer0_load_balancer.datasources.url}"
}

output "load_balancer_service_id" {
  value = "${data.layer0_load_balancer.datasources.service_id}"
}

output "load_balancer_service_name" {
  value = "${data.layer0_load_balancer.datasources.service_name}"
}

output "service_id" {
  value = "${data.layer0_service.datasources.id}"
}

output "service_environment_name" {
  value = "${data.layer0_service.datasources.environment_name}"
}

output "service_lb_name" {
  value = "${data.layer0_service.datasources.load_balancer_name}"
}

output "service_lb_id" {
  value = "${data.layer0_service.datasources.load_balancer_id}"
}

output "service_scale" {
  value = "${data.layer0_service.datasources.scale}"
}
