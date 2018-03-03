#!/usr/bin/env bats

@test "create" {
  l0 deploy create ./common/Service_stateful.Dockerrun.aws.json dpl_name_stateful
  l0 deploy create ./common/Service_stateless.Dockerrun.aws.json dpl_name_stateless
}

@test "get" {
  l0 deploy get dpl_name_stateful
  l0 deploy get dpl_name_stateless
  l0 deploy get dpl_name*
}

@test "list" {
  l0 deploy list
  l0 deploy list --all
}

@test "delete" {
  l0 deploy delete dpl_name_stateful
  l0 deploy delete dpl_name_stateless
}
