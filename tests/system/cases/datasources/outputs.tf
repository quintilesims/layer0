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

output "environment_scale" {
  value = "${data.layer0_environment.datasources.scale}"
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

output "stateless_deploy_id" {
  value = "${data.layer0_deploy.datasources_stateless.id}"
}

output "stateless_deploy_id_expected" {
  value = "${layer0_deploy.datasources_stateless.id}"
}

output "stateless_deploy_name" {
  value = "${data.layer0_deploy.datasources_stateless.version}"
}

output "stateless_deploy_name_expected" {
  value = "${layer0_deploy.datasources_stateless.version}"
}

output "stateless_deploy_version" {
  value = "${data.layer0_deploy.datasources_stateless.version}"
}

output "stateless_deploy_version_expected" {
  value = "${layer0_deploy.datasources_stateless.version}"
}

output "stateful_deploy_id" {
  value = "${data.layer0_deploy.datasources_stateful.id}"
}

output "stateful_deploy_id_expected" {
  value = "${layer0_deploy.datasources_stateful.id}"
}

output "stateful_deploy_name" {
  value = "${data.layer0_deploy.datasources_stateful.version}"
}

output "stateful_deploy_name_expected" {
  value = "${layer0_deploy.datasources_stateful.version}"
}

output "stateful_deploy_version" {
  value = "${data.layer0_deploy.datasources_stateful.version}"
}

output "stateful_deploy_version_expected" {
  value = "${layer0_deploy.datasources_stateful.version}"
}

# -- Load Balancer --

output "stateless_load_balancer_id" {
  value = "${data.layer0_load_balancer.datasources_stateless.id}"
}

output "stateless_load_balancer_id_expected" {
  value = "${layer0_load_balancer.datasources_stateless.id}"
}

output "stateless_load_balancer_name" {
  value = "${data.layer0_load_balancer.datasources_stateless.name}"
}

output "stateless_load_balancer_name_expected" {
  value = "${layer0_load_balancer.datasources_stateless.name}"
}

output "stateless_load_balancer_environment_name" {
  value = "${data.layer0_load_balancer.datasources_stateless.environment_name}"
}

output "stateless_load_balancer_environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "stateless_load_balancer_private" {
  value = "${data.layer0_load_balancer.datasources_stateless.private}"
}

output "stateless_load_balancer_private_expected" {
  value = "${layer0_load_balancer.datasources_stateless.private}"
}

output "stateless_load_balancer_url" {
  value = "${data.layer0_load_balancer.datasources_stateless.url}"
}

output "stateless_load_balancer_url_expected" {
  value = "${layer0_load_balancer.datasources_stateless.url}"
}

output "stateless_load_balancer_type" {
    value = "${data.layer0_load_balancer.datasources_stateless.type}"
}

output "stateless_load_balancer_type_expected" {
    value = "${layer0_load_balancer.datasources_stateless.type}"
}

output "stateful_load_balancer_id" {
  value = "${data.layer0_load_balancer.datasources_stateful.id}"
}

output "stateful_load_balancer_id_expected" {
  value = "${layer0_load_balancer.datasources_stateful.id}"
}

output "stateful_load_balancer_name" {
  value = "${data.layer0_load_balancer.datasources_stateful.name}"
}

output "stateful_load_balancer_name_expected" {
  value = "${layer0_load_balancer.datasources_stateful.name}"
}

output "stateful_load_balancer_environment_name" {
  value = "${data.layer0_load_balancer.datasources_stateful.environment_name}"
}

output "stateful_load_balancer_environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "stateful_load_balancer_private" {
  value = "${data.layer0_load_balancer.datasources_stateful.private}"
}

output "stateful_load_balancer_private_expected" {
  value = "${layer0_load_balancer.datasources_stateful.private}"
}

output "stateful_load_balancer_url" {
  value = "${data.layer0_load_balancer.datasources_stateful.url}"
}

output "stateful_load_balancer_url_expected" {
  value = "${layer0_load_balancer.datasources_stateful.url}"
}

output "stateful_load_balancer_type" {
    value = "${data.layer0_load_balancer.datasources_stateful.type}"
}

output "stateful_load_balancer_type_expected" {
    value = "${layer0_load_balancer.datasources_stateful.type}"
}

# -- Service --

output "stateless_service_id" {
  value = "${data.layer0_service.datasources_stateless.id}"
}

output "stateless_service_id_expected" {
  value = "${layer0_service.datasources_stateless.id}"
}

output "stateless_service_name" {
  value = "${data.layer0_service.datasources_stateless.name}"
}

output "stateless_service_name_expected" {
  value = "${layer0_service.datasources_stateless.name}"
}

output "stateless_service_environment_id" {
  value = "${data.layer0_service.datasources_stateless.environment_id}"
}

output "stateless_service_environment_id_expected" {
  value = "${layer0_environment.datasources.id}"
}

output "stateless_service_environment_name" {
  value = "${data.layer0_service.datasources_stateless.environment_name}"
}

output "stateless_service_environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "stateless_service_scale" {
  value = "${data.layer0_service.datasources_stateless.scale}"
}

output "stateless_service_scale_expected" {
  value = "${layer0_service.datasources_stateless.scale}"
}

output "stateless_service_stateful" {
  value = "${data.layer0_service.datasources_stateless.stateful}"
}

output "stateless_service_stateful_expected" {
  value = "${layer0_service.datasources_stateless.stateful}"
}

output "stateful_service_id" {
  value = "${data.layer0_service.datasources_stateful.id}"
}

output "stateful_service_id_expected" {
  value = "${layer0_service.datasources_stateful.id}"
}

output "stateful_service_name" {
  value = "${data.layer0_service.datasources_stateful.name}"
}

output "stateful_service_name_expected" {
  value = "${layer0_service.datasources_stateful.name}"
}

output "stateful_service_environment_id" {
  value = "${data.layer0_service.datasources_stateful.environment_id}"
}

output "stateful_service_environment_id_expected" {
  value = "${layer0_environment.datasources.id}"
}

output "stateful_service_environment_name" {
  value = "${data.layer0_service.datasources_stateful.environment_name}"
}

output "stateful_service_environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "stateful_service_scale" {
  value = "${data.layer0_service.datasources_stateful.scale}"
}

output "stateful_service_scale_expected" {
  value = "${layer0_service.datasources_stateful.scale}"
}

output "stateful_service_stateful" {
  value = "${data.layer0_service.datasources_stateful.stateful}"
}

output "stateful_service_stateful_expected" {
  value = "${layer0_service.datasources_stateful.stateful}"
}
