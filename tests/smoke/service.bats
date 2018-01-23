#!/usr/bin/env bats

@test "create" {
  l0 environment create env_name
  l0 deploy create ./common/Service.Dockerrun.aws.json dpl_name1
  l0 deploy create ./common/Service.Dockerrun.aws.json dpl_name2
  l0 loadbalancer create --port 80:80/http env_name lb_name1
  l0 service create env_name svc_name1 dpl_name1
  l0 service create --loadbalancer lb_name1 --scale 2 env_name svc_name2 dpl_name1
}

@test "get" {
  l0 service get svc_name1
  l0 service get svc_name2
  l0 service get svc_name*
}

@test "list" {
  l0 service list
}

@test "scale" {
  l0 service scale svc_name1 2
}

@test "update" {
  l0 service update svc_name1 dpl_name2
}

@test "logs" {
  l0 service logs svc_name1
  l0 service logs --tail 100 --start '2001-01-01 01:01' --end '2012-12-12 12:12' svc_name1
}

@test "delete" {
  l0 service delete svc_name1
  l0 service delete svc_name2
  l0 loadbalancer delete lb_name1
  l0 deploy delete dpl_name1
  l0 deploy delete dpl_name2
  l0 environment delete env_name
}
