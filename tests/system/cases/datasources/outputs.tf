# -- Environment --

output "environment_name" {
  value = "${data.layer0_environment.datasources.name}"
}

output "environment_name_expected" {
  value = "${layer0_environment.datasources.name}"
}

output "environment_id" {
  value = "${data.layer0_environment.datasources.id}"
}

output "environment_id_expected" {
  value = "${layer0_environment.datasources.id}"
}

output "environment_instance_type" {
  value = "${data.layer0_environment.datasources.instance_type}"
}

output "environment_instance_type_expected" {
  value = "${layer0_environment.datasources.instance_type}"
}

output "environment_scale" {
  value = "${data.layer0_environment.datasources.scale}"
}

output "environment_scale_expected" {
  value = "${layer0_environment.datasources.scale}"
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

output "environment_security_group_id" {
  value = "${data.layer0_environment.datasources.security_group_id}"
}

output "environment_security_group_id_expected" {
  value = "${layer0_environment.datasources.security_group_id}"
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

output "load_balancer_environment_id" {
  value = "${data.layer0_load_balancer.datasources.environment_id}"
}

output "load_balancer_environment_id_expected" {
  value = "${layer0_environment.datasources.id}"
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
  value = "${layer0_load_balancer.datasources.private}"
}

output "load_balancer_url" {
  value = "${data.layer0_load_balancer.datasources.url}"
}

output "load_balancer_url_expected" {
  value = "${layer0_load_balancer.datasources.url}"
}

output "load_balancer_type" {
    value = "${data.layer0_load_balancer.datasources.type}"
}

output "load_balancer_type_expected" {
    value = "${layer0_load_balancer.datasources.type}"
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
