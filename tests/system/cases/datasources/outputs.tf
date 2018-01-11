# -- Environment --

output "environment_id" {
  value = "${data.layer0_environment.datasources.id}"
}

output "environment_id_expected" {
  value = "${layer0_environment.datasources.id}"
}

output "environment_name" {
  value = "${data.layer0_environment.datasources.name}"
}

output "environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "environment_size" {
  value = "${data.layer0_environment.datasources.size}"
}

output "environment_size_expected" {
  value = "${layer0_environment.datasources.size}"
}

output "environment_min_count" {
  value = "${data.layer0_environment.datasources.min_count}"
}

output "environment_min_count_expected" {
  value = "${layer0_environment.datasources.cluster_count}"
}

output "environment_os" {
  value = "${data.layer0_environment.datasources.os}"
}

output "environment_os_expected" {
  value = "${layer0_environment.datasources.os}"
}

output "environment_ami" {
  value = "${data.layer0_environment.datasources.ami}"
}

output "environment_ami_expected" {
  value = "${layer0_environment.datasources.ami}"
}

# -- Deploy --

output "deploy_id" {
  value = "${data.layer0_deploy.datasources.id}"
}

output "deploy_id_expected" {
  value = "${layer0_deploy.datasources.id}"
}

output "deploy_name" {
  value = "${data.layer0_deploy.datasources.version}"
}

output "deploy_name_expected" {
  value = "${layer0_deploy.datasources.version}"
}

output "deploy_version" {
  value = "${data.layer0_deploy.datasources.version}"
}

output "deploy_version_expected" {
  value = "${layer0_deploy.datasources.version}"
}

# -- Load Balancer --

output "load_balancer_id" {
  value = "${data.layer0_load_balancer.datasources.id}"
}

output "load_balancer_id_expected" {
  value = "${layer0_load_balancer.datasources.id}"
}

output "load_balancer_name" {
  value = "${data.layer0_load_balancer.datasources.name}"
}

output "load_balancer_name_expected" {
  value = "${layer0_load_balancer.datasources.name}"
}

output "load_balancer_environment_name" {
  value = "${data.layer0_load_balancer.datasources.environment_name}"
}

output "load_balancer_environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "load_balancer_private" {
  value = "${data.layer0_load_balancer.datasources.private}"
}

output "load_balancer_private_expected" {
  value = "false"
}

output "load_balancer_url" {
  value = "${data.layer0_load_balancer.datasources.url}"
}

output "load_balancer_url_expected" {
  value = "${layer0_load_balancer.datasources.url}"
}

# -- Service --

output "service_id" {
  value = "${data.layer0_service.datasources.id}"
}

output "service_id_expected" {
  value = "${layer0_service.datasources.id}"
}

output "service_name" {
  value = "${data.layer0_service.datasources.name}"
}

output "service_name_expected" {
  value = "${layer0_service.datasources.name}"
}

output "service_environment_id" {
  value = "${data.layer0_service.datasources.environment_id}"
}

output "service_environment_id_expected" {
  value = "${layer0_environment.datasources.id}"
}

output "service_environment_name" {
  value = "${data.layer0_service.datasources.environment_name}"
}

output "service_environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "service_scale" {
  value = "${data.layer0_service.datasources.scale}"
}

output "service_scale_expected" {
  value = "${layer0_service.datasources.scale}"
}
