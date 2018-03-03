#!/usr/bin/env bats

@test "create" {
  l0 environment create env_name
  l0 deploy create ./common/Service_stateful.Dockerrun.aws.json dpl_name_stateful1
  l0 deploy create ./common/Service_stateful.Dockerrun.aws.json dpl_name_stateful2
  l0 deploy create ./common/Service_stateless.Dockerrun.aws.json dpl_name_stateless
  l0 loadbalancer create --port 80:80/http env_name lb_name1
  l0 service create --stateful --loadbalancer lb_name1 --scale 2 env_name svc_name_stateful dpl_name_stateful1
  l0 service create env_name svc_name_stateless dpl_name_stateless
}

@test "get" {
  l0 service get svc_name_stateful
  l0 service get svc_name_stateless
  l0 service get svc_name*
}

@test "list" {
  l0 service list
}

@test "scale" {
  l0 service scale svc_name_stateful 1
}

@test "update" {
  l0 service update svc_name_stateful dpl_name_stateful2
}

@test "logs" {
  l0 service logs svc_name_stateful
  l0 service logs --tail 100 --start '2001-01-01 01:01' --end '2012-12-12 12:12' svc_name_stateful
}

@test "delete" {
  l0 service delete svc_name_stateful
  l0 service delete svc_name_stateless
  l0 loadbalancer delete lb_name1
  l0 deploy delete dpl_name_stateful1
  l0 deploy delete dpl_name_stateful2
  l0 deploy delete dpl_name_stateless
  l0 environment delete env_name
}
